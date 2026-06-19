package handlers

import (
	"html"
	"html/template"
	"net/http"
	"regexp"
	"strings"

	"github.com/jlettori/piponews/internal/auth"
	"github.com/jlettori/piponews/internal/db"
	"github.com/jlettori/piponews/internal/i18n"
	"github.com/jlettori/piponews/internal/models"
	"github.com/jlettori/piponews/internal/templates"
)

var validEmail = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

var profileDialogTemplate = template.Must(
	template.New("profile_dialog.html").Funcs(templateFuncs()).ParseFS(templates.FS, "profile_dialog.html"),
)

type ProfileHandler struct {
	DB *db.DB
}

func (h *ProfileHandler) GET(w http.ResponseWriter, r *http.Request) {
	locale := detectLocale(r)
	userID := auth.GetUserID(r)

	var user models.User
	h.DB.Get(&user, "SELECT * FROM users WHERE id = ?", userID)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	profileDialogTemplate.Execute(w, map[string]any{
		"Locale":            locale,
		"Username":          user.Username,
		"FirstName":         user.FirstName,
		"LastName":          user.LastName,
		"Email":             user.Email,
		"PreferredLanguage": user.PreferredLanguage,
	})
}

func (h *ProfileHandler) POST(w http.ResponseWriter, r *http.Request) {
	locale := detectLocale(r)
	userID := auth.GetUserID(r)

	if err := r.ParseForm(); err != nil {
		w.Write(errorAlert(i18n.T(locale, i18n.InvalidFormData)))
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	firstName := strings.TrimSpace(r.FormValue("first_name"))
	lastName := strings.TrimSpace(r.FormValue("last_name"))
	email := strings.TrimSpace(r.FormValue("email"))
	preferredLanguage := strings.TrimSpace(r.FormValue("preferred_language"))

	if username == "" {
		w.Write(errorAlert(i18n.T(locale, i18n.UsernamePasswordRequired)))
		return
	}
	if len(username) < 2 || len(username) > 50 {
		w.Write(errorAlert(i18n.T(locale, i18n.UsernameLength)))
		return
	}
	if !validUsername.MatchString(username) {
		w.Write(errorAlert(i18n.T(locale, i18n.UsernameChars)))
		return
	}

	if email != "" && !validEmail.MatchString(email) {
		w.Write(errorAlert(i18n.T(locale, i18n.EmailInvalid)))
		return
	}

	if preferredLanguage != "" && preferredLanguage != "en" && preferredLanguage != "fr" && preferredLanguage != "it" {
		preferredLanguage = ""
	}

	var existingID int64
	err := h.DB.Get(&existingID, "SELECT id FROM users WHERE username = ? AND id != ?", username, userID)
	if err == nil {
		w.Write(errorAlert(i18n.T(locale, i18n.UsernameTaken)))
		return
	}

	_, err = h.DB.Exec(`
		UPDATE users SET username = ?, first_name = ?, last_name = ?, email = ?, preferred_language = ?
		WHERE id = ?
	`, username, firstName, lastName, email, preferredLanguage, userID)
	if err != nil {
		w.Write(errorAlert(i18n.T(locale, i18n.ProfileUpdateFailed)))
		return
	}

	auth.CreateSessionCookie(w, userID, username)
	http.Redirect(w, r, "/feeds", http.StatusSeeOther)
}

func errorAlert(msg string) []byte {
	return []byte(`<div id="profile-dialog-alert"><div class="alert alert-error">` + html.EscapeString(msg) + `</div></div>`)
}
