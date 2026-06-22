package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/jlettori/piponews/internal/db"
	"golang.org/x/crypto/bcrypt"
)

type ctxKey string

const (
	ctxKeyUserID   ctxKey = "user_id"
	ctxKeyUsername ctxKey = "username"
	cookieName            = "session"
)

type Session struct {
	UserID   int64
	Username string
}

type Store struct {
	db *db.DB
}

func NewStore(database *db.DB) *Store {
	return &Store{db: database}
}

const (
	DefaultSessionDuration    = 24 * time.Hour
	RememberedSessionDuration = 30 * 24 * time.Hour
)

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func CheckPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (s *Store) CreateSession(w http.ResponseWriter, userID int64, username string, duration time.Duration) error {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return err
	}
	tokenStr := base64.RawURLEncoding.EncodeToString(token)

	hash := sha256.Sum256([]byte(tokenStr))
	hashStr := base64.RawURLEncoding.EncodeToString(hash[:])

	expiresAt := time.Now().Add(duration)
	if _, err := s.db.Exec(
		"INSERT INTO user_sessions (user_id, token_hash, username, expires_at) VALUES (?, ?, ?, ?)",
		userID, hashStr, username, expiresAt,
	); err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    tokenStr,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  expiresAt,
		MaxAge:   int(duration.Seconds()),
	})
	return nil
}

func (s *Store) GetSession(r *http.Request) (*Session, error) {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256([]byte(c.Value))
	hashStr := base64.RawURLEncoding.EncodeToString(hash[:])

	var sess Session
	var expiresAt time.Time
	err = s.db.QueryRow(
		"SELECT user_id, username, expires_at FROM user_sessions WHERE token_hash = ?",
		hashStr,
	).Scan(&sess.UserID, &sess.Username, &expiresAt)
	if err != nil {
		return nil, fmt.Errorf("invalid session")
	}

	if time.Now().After(expiresAt) {
		s.db.Exec("DELETE FROM user_sessions WHERE token_hash = ?", hashStr)
		return nil, fmt.Errorf("session expired")
	}

	return &sess, nil
}

func (s *Store) ClearSession(w http.ResponseWriter, r *http.Request) {
	if c, err := r.Cookie(cookieName); err == nil {
		hash := sha256.Sum256([]byte(c.Value))
		s.db.Exec("DELETE FROM user_sessions WHERE token_hash = ?",
			base64.RawURLEncoding.EncodeToString(hash[:]))
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
}

func GetUserID(r *http.Request) int64 {
	v, _ := r.Context().Value(ctxKeyUserID).(int64)
	return v
}

func GetUsername(r *http.Request) string {
	v, _ := r.Context().Value(ctxKeyUsername).(string)
	return v
}

func (s *Store) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, err := s.GetSession(r)
		if err != nil {
			s.ClearSession(w, r)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), ctxKeyUserID, sess.UserID)
		ctx = context.WithValue(ctx, ctxKeyUsername, sess.Username)
		next(w, r.WithContext(ctx))
	}
}
