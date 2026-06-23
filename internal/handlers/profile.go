package handlers

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/jlettori/piponews/internal/auth"
	"github.com/jlettori/piponews/internal/db"
	"github.com/jlettori/piponews/internal/i18n"
	"github.com/jlettori/piponews/internal/models"
)

var validEmail = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type ProfileHandler struct {
	DB       *db.DB
	Sessions *auth.Store
}

func (h *ProfileHandler) render(w http.ResponseWriter, r *http.Request, extra map[string]any) {
	locale := detectLocale(r)
	userID := auth.GetUserID(r)

	var user models.User
	h.DB.Get(&user, `
		SELECT id, username, password_hash, first_name, last_name, email, preferred_language, created_at
		FROM users
		WHERE id = ?
	`, userID)

	data := map[string]any{
		"Locale":            locale,
		"Username":          user.Username,
		"FirstName":         user.FirstName,
		"LastName":          user.LastName,
		"Email":             user.Email,
		"PreferredLanguage": user.PreferredLanguage,
	}
	for k, v := range extra {
		data[k] = v
	}

	t := parseTemplates("profile.html")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.ExecuteTemplate(w, "base.html", data); err != nil {
		log.Printf("template error: %v", err)
	}
}

func (h *ProfileHandler) GET(w http.ResponseWriter, r *http.Request) {
	h.render(w, r, nil)
}

func (h *ProfileHandler) POST(w http.ResponseWriter, r *http.Request) {
	locale := detectLocale(r)
	userID := auth.GetUserID(r)

	if err := r.ParseForm(); err != nil {
		h.render(w, r, map[string]any{"Error": i18n.T(locale, i18n.InvalidFormData)})
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	firstName := strings.TrimSpace(r.FormValue("first_name"))
	lastName := strings.TrimSpace(r.FormValue("last_name"))
	email := strings.TrimSpace(r.FormValue("email"))
	preferredLanguage := strings.TrimSpace(r.FormValue("preferred_language"))

	if username == "" {
		h.render(w, r, map[string]any{"Error": i18n.T(locale, i18n.UsernamePasswordRequired)})
		return
	}
	if len(username) < 2 || len(username) > 50 {
		h.render(w, r, map[string]any{"Error": i18n.T(locale, i18n.UsernameLength)})
		return
	}
	if !validUsername.MatchString(username) {
		h.render(w, r, map[string]any{"Error": i18n.T(locale, i18n.UsernameChars)})
		return
	}

	if email != "" && !validEmail.MatchString(email) {
		h.render(w, r, map[string]any{"Error": i18n.T(locale, i18n.EmailInvalid)})
		return
	}

	if preferredLanguage != "" && preferredLanguage != "en" && preferredLanguage != "fr" && preferredLanguage != "it" {
		preferredLanguage = ""
	}

	var existingID int64
	err := h.DB.Get(&existingID, "SELECT id FROM users WHERE username = ? AND id != ?", username, userID)
	if err == nil {
		h.render(w, r, map[string]any{"Error": i18n.T(locale, i18n.UsernameTaken)})
		return
	}

	_, err = h.DB.Exec(`
		UPDATE users SET username = ?, first_name = ?, last_name = ?, email = ?, preferred_language = ?
		WHERE id = ?
	`, username, firstName, lastName, email, preferredLanguage, userID)
	if err != nil {
		h.render(w, r, map[string]any{"Error": i18n.T(locale, i18n.ProfileUpdateFailed)})
		return
	}

	if err := h.Sessions.CreateSession(w, userID, username, auth.DefaultSessionDuration); err != nil {
		h.render(w, r, map[string]any{"Error": i18n.T(locale, i18n.InternalError)})
		return
	}

	if preferredLanguage != "" {
		http.SetCookie(w, &http.Cookie{
			Name:    i18n.LocaleCookie,
			Value:   preferredLanguage,
			Path:    "/",
			Expires: time.Now().Add(365 * 24 * time.Hour),
		})
	}

	http.Redirect(w, r, "/feeds", http.StatusSeeOther)
}
