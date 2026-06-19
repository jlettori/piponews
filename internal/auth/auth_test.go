package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "mySecretPassword123!"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword() returned empty hash")
	}
	if !CheckPassword(password, hash) {
		t.Error("CheckPassword() = false; want true")
	}
	if CheckPassword("wrongPassword", hash) {
		t.Error("CheckPassword(wrong) = true; want false")
	}
}

func TestCreateAndVerifySession(t *testing.T) {
	w := httptest.NewRecorder()
	err := CreateSessionCookie(w, 42, "testuser")
	if err != nil {
		t.Fatalf("CreateSessionCookie() error = %v", err)
	}

	resp := w.Result()
	cookies := resp.Cookies()
	if len(cookies) != 1 {
		t.Fatalf("got %d cookies; want 1", len(cookies))
	}

	cookie := cookies[0]
	if cookie.Name != "session" {
		t.Errorf("cookie name = %q; want %q", cookie.Name, "session")
	}
	if cookie.Value == "" {
		t.Fatal("cookie value is empty")
	}
	if !cookie.HttpOnly {
		t.Error("cookie not HttpOnly")
	}

	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(cookie)
	claims, err := GetSession(r)
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}
	if claims.UserID != 42 {
		t.Errorf("claims.UserID = %d; want 42", claims.UserID)
	}
	if claims.Username != "testuser" {
		t.Errorf("claims.Username = %q; want %q", claims.Username, "testuser")
	}
}

func TestGetSession_NoCookie(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	_, err := GetSession(r)
	if err == nil {
		t.Fatal("GetSession() expected error for missing cookie")
	}
}

func TestGetSession_InvalidToken(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: "session", Value: "invalid-token"})
	_, err := GetSession(r)
	if err == nil {
		t.Fatal("GetSession() expected error for invalid token")
	}
}

func TestGetSession_TamperedPayload(t *testing.T) {
	w := httptest.NewRecorder()
	CreateSessionCookie(w, 1, "user")
	resp := w.Result()
	cookie := resp.Cookies()[0]

	parts := strings.Split(cookie.Value, ".")
	if len(parts) != 2 {
		t.Fatal("unexpected token format")
	}
	tampered := parts[0] + ".invalidsignature"

	r := httptest.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: "session", Value: tampered})
	_, err := GetSession(r)
	if err == nil {
		t.Fatal("GetSession() expected error for tampered token")
	}
}

func TestClearSessionCookie(t *testing.T) {
	w := httptest.NewRecorder()
	ClearSessionCookie(w)
	resp := w.Result()
	cookies := resp.Cookies()
	if len(cookies) != 1 {
		t.Fatalf("got %d cookies; want 1", len(cookies))
	}
	c := cookies[0]
	if c.Value != "" {
		t.Errorf("cookie value = %q; want empty", c.Value)
	}
	if !c.Expires.IsZero() && c.Expires.After(time.Now()) {
		t.Error("cookie should be expired")
	}
}

func TestRequireAuth_ValidSession(t *testing.T) {
	w := httptest.NewRecorder()
	CreateSessionCookie(w, 7, "alice")
	resp := w.Result()
	cookie := resp.Cookies()[0]

	r := httptest.NewRequest("GET", "/feeds", nil)
	r.AddCookie(cookie)

	var handled bool
	handler := RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		handled = true
		if id := GetUserID(r); id != 7 {
			t.Errorf("GetUserID() = %d; want 7", id)
		}
		if un := GetUsername(r); un != "alice" {
			t.Errorf("GetUsername() = %q; want %q", un, "alice")
		}
	})

	handler(w, r)
	if !handled {
		t.Error("handler was not called")
	}
}

func TestRequireAuth_NoSession(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/feeds", nil)

	var handled bool
	handler := RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		handled = true
	})

	handler(w, r)
	if handled {
		t.Error("handler was called despite no session")
	}
	if w.Code != http.StatusSeeOther {
		t.Errorf("status = %d; want %d", w.Code, http.StatusSeeOther)
	}
	loc := w.Header().Get("Location")
	if loc != "/login" {
		t.Errorf("Location = %q; want %q", loc, "/login")
	}
}
