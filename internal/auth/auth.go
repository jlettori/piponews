package auth

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type ctxKey string

const (
	ctxKeyUserID   ctxKey = "user_id"
	ctxKeyUsername ctxKey = "username"
	cookieName            = "session"
)

var hmacKey = func() []byte {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic("failed to generate HMAC key: " + err.Error())
	}
	return b
}()

type sessionClaims struct {
	UserID   int64  `json:"uid"`
	Username string `json:"un"`
	Expiry   int64  `json:"exp"`
}

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

func CheckPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func CreateSessionCookie(w http.ResponseWriter, userID int64, username string) error {
	claims := sessionClaims{
		UserID:   userID,
		Username: username,
		Expiry:   time.Now().Add(24 * time.Hour).Unix(),
	}
	return setCookie(w, claims)
}

func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
	})
}

func setCookie(w http.ResponseWriter, claims sessionClaims) error {
	b, err := json.Marshal(claims)
	if err != nil {
		return err
	}
	payload := base64.RawURLEncoding.EncodeToString(b)
	mac := computeHMAC(payload)
	token := payload + "." + mac
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(claims.Expiry, 0),
	})
	return nil
}

func computeHMAC(payload string) string {
	m := hmac.New(sha256.New, hmacKey)
	m.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(m.Sum(nil))
}

func verifyToken(token string) (*sessionClaims, error) {
	dot := strings.LastIndex(token, ".")
	if dot < 1 {
		return nil, fmt.Errorf("invalid token format")
	}
	payload := token[:dot]
	mac := token[dot+1:]
	expected := computeHMAC(payload)
	if !hmac.Equal([]byte(mac), []byte(expected)) {
		return nil, fmt.Errorf("invalid signature")
	}
	b, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return nil, err
	}
	var claims sessionClaims
	if err := json.Unmarshal(b, &claims); err != nil {
		return nil, err
	}
	if time.Now().Unix() > claims.Expiry {
		return nil, fmt.Errorf("session expired")
	}
	return &claims, nil
}

func GetSession(r *http.Request) (*sessionClaims, error) {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return nil, err
	}
	return verifyToken(c.Value)
}

func GetUserID(r *http.Request) int64 {
	v, _ := r.Context().Value(ctxKeyUserID).(int64)
	return v
}

func GetUsername(r *http.Request) string {
	v, _ := r.Context().Value(ctxKeyUsername).(string)
	return v
}

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(cookieName)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		claims, err := verifyToken(c.Value)
		if err != nil {
			ClearSessionCookie(w)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(), ctxKeyUserID, claims.UserID)
		ctx = context.WithValue(ctx, ctxKeyUsername, claims.Username)
		next(w, r.WithContext(ctx))
	}
}
