package handlers

import (
	"html/template"
	"net/http"
	"regexp"
	"strings"

	"github.com/jlettori/piponews/internal/auth"
	"github.com/jlettori/piponews/internal/db"
	"github.com/jlettori/piponews/internal/i18n"
	"github.com/jlettori/piponews/internal/templates"
)

var validUsername = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

type AuthHandler struct {
	DB       *db.DB
	Sessions *auth.Store
}

func (h *AuthHandler) AuthPage(w http.ResponseWriter, r *http.Request) {
	renderAuth(w, r, "login", nil)
}

func (h *AuthHandler) LoginGET(w http.ResponseWriter, r *http.Request) {
	if _, err := h.Sessions.GetSession(r); err == nil {
		http.Redirect(w, r, "/feeds", http.StatusSeeOther)
		return
	}
	renderAuth(w, r, "login", nil)
}

func (h *AuthHandler) LoginPOST(w http.ResponseWriter, r *http.Request) {
	locale := detectLocale(r)
	if err := r.ParseForm(); err != nil {
		renderAuth(w, r, "login", map[string]any{"Error": i18n.T(locale, i18n.InvalidFormData)})
		return
	}
	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")

	if username == "" || password == "" {
		renderAuth(w, r, "login", map[string]any{"Error": i18n.T(locale, i18n.UsernamePasswordRequired)})
		return
	}

	var user struct {
		ID           int64  `db:"id"`
		Username     string `db:"username"`
		PasswordHash string `db:"password_hash"`
	}
	err := h.DB.Get(&user, "SELECT id, username, password_hash FROM users WHERE username = ?", username)
	if err != nil {
		renderAuth(w, r, "login", map[string]any{"Error": i18n.T(locale, i18n.InvalidUsernamePassword)})
		return
	}

	if !auth.CheckPassword(password, user.PasswordHash) {
		renderAuth(w, r, "login", map[string]any{"Error": i18n.T(locale, i18n.InvalidUsernamePassword)})
		return
	}

	duration := auth.DefaultSessionDuration
	if r.FormValue("remember") == "on" {
		duration = auth.RememberedSessionDuration
	}
	if err := h.Sessions.CreateSession(w, user.ID, user.Username, duration); err != nil {
		renderAuth(w, r, "login", map[string]any{"Error": i18n.T(locale, i18n.InternalError)})
		return
	}
	http.Redirect(w, r, "/feeds", http.StatusSeeOther)
}

func (h *AuthHandler) RegisterGET(w http.ResponseWriter, r *http.Request) {
	if _, err := h.Sessions.GetSession(r); err == nil {
		http.Redirect(w, r, "/feeds", http.StatusSeeOther)
		return
	}
	renderAuth(w, r, "register", nil)
}

func (h *AuthHandler) RegisterPOST(w http.ResponseWriter, r *http.Request) {
	locale := detectLocale(r)
	if err := r.ParseForm(); err != nil {
		renderAuth(w, r, "register", map[string]any{"Error": i18n.T(locale, i18n.InvalidFormData)})
		return
	}
	username := strings.TrimSpace(r.FormValue("username"))
	password := r.FormValue("password")

	if username == "" || password == "" {
		renderAuth(w, r, "register", map[string]any{"Error": i18n.T(locale, i18n.UsernamePasswordRequired)})
		return
	}
	if len(username) < 2 || len(username) > 50 {
		renderAuth(w, r, "register", map[string]any{"Error": i18n.T(locale, i18n.UsernameLength)})
		return
	}
	if !validUsername.MatchString(username) {
		renderAuth(w, r, "register", map[string]any{"Error": i18n.T(locale, i18n.UsernameChars)})
		return
	}
	if len(password) < 4 {
		renderAuth(w, r, "register", map[string]any{"Error": i18n.T(locale, i18n.PasswordLength)})
		return
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		renderAuth(w, r, "register", map[string]any{"Error": i18n.T(locale, i18n.InternalError)})
		return
	}

	preferredLanguage := string(locale)
	res, err := h.DB.Exec("INSERT INTO users (username, password_hash, preferred_language) VALUES (?, ?, ?)", username, hash, preferredLanguage)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			renderAuth(w, r, "register", map[string]any{"Error": i18n.T(locale, i18n.UsernameTaken)})
			return
		}
		renderAuth(w, r, "register", map[string]any{"Error": i18n.T(locale, i18n.InternalError)})
		return
	}

	id, _ := res.LastInsertId()
	if err := h.Sessions.CreateSession(w, id, username, auth.DefaultSessionDuration); err != nil {
		renderAuth(w, r, "register", map[string]any{"Error": i18n.T(locale, i18n.InternalError)})
		return
	}
	http.Redirect(w, r, "/feeds", http.StatusSeeOther)
}

func (h *AuthHandler) LogoutPOST(w http.ResponseWriter, r *http.Request) {
	h.Sessions.ClearSession(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) CheckSession(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if _, err := h.Sessions.GetSession(r); err == nil {
		w.Write([]byte(`{"authenticated":true}`))
	} else {
		w.Write([]byte(`{"authenticated":false}`))
	}
}

func renderAuth(w http.ResponseWriter, r *http.Request, mode string, data map[string]any) {
	if data == nil {
		data = map[string]any{}
	}
	data["AuthMode"] = mode
	data["Locale"] = detectLocale(r)
	t := template.Must(template.New("base").Funcs(templateFuncs()).ParseFS(
		templates.FS,
		"base.html", "auth.html",
	))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t.ExecuteTemplate(w, "base.html", data)
}
