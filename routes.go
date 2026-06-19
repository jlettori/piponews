package main

import (
	"net/http"
	"time"

	"github.com/jlettori/piponews/internal/auth"
	"github.com/jlettori/piponews/internal/db"
	"github.com/jlettori/piponews/internal/handlers"
	"github.com/jlettori/piponews/internal/i18n"
)

func newRouter(database *db.DB) http.Handler {
	authH := &handlers.AuthHandler{DB: database}
	feedsH := &handlers.FeedsHandler{DB: database}
	entriesH := &handlers.EntriesHandler{DB: database}
	exportsH := &handlers.ExportsHandler{DB: database}

	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(staticFileSystem())))

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if _, err := auth.GetSession(r); err == nil {
			http.Redirect(w, r, "/feeds", http.StatusSeeOther)
			return
		}
		authH.AuthPage(w, r)
	})

	mux.HandleFunc("GET /login", authH.LoginGET)
	mux.HandleFunc("POST /login", authH.LoginPOST)
	mux.HandleFunc("GET /register", authH.RegisterGET)
	mux.HandleFunc("POST /register", authH.RegisterPOST)
	mux.Handle("POST /logout", auth.RequireAuth(authH.LogoutPOST))

	mux.HandleFunc("GET /lang", auth.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		locale := r.URL.Query().Get("locale")
		if locale != "" {
			http.SetCookie(w, &http.Cookie{
				Name:    i18n.LocaleCookie,
				Value:   locale,
				Path:    "/",
				Expires: time.Now().Add(365 * 24 * time.Hour),
			})
		}
		http.Redirect(w, r, "/feeds", http.StatusSeeOther)
	}))

	mux.Handle("GET /feeds", auth.RequireAuth(feedsH.List))
	mux.Handle("POST /feeds", auth.RequireAuth(feedsH.Create))
	mux.Handle("DELETE /feeds/{id}", auth.RequireAuth(feedsH.Delete))
	mux.Handle("POST /feeds/{id}/refresh", auth.RequireAuth(feedsH.Refresh))
	mux.Handle("POST /entries/refresh-all", auth.RequireAuth(feedsH.RefreshAll))

	mux.Handle("GET /entries", auth.RequireAuth(entriesH.List))
	mux.Handle("GET /entries/more", auth.RequireAuth(entriesH.More))
	mux.Handle("POST /entries/{id}/toggle-select", auth.RequireAuth(entriesH.ToggleSelect))
	mux.Handle("POST /entries/clear-selection", auth.RequireAuth(entriesH.ClearSelection))

	mux.Handle("POST /entries/export", auth.RequireAuth(exportsH.Export))

	profileH := &handlers.ProfileHandler{DB: database}
	mux.Handle("GET /profile", auth.RequireAuth(profileH.GET))
	mux.Handle("POST /profile", auth.RequireAuth(profileH.POST))

	return mux
}
