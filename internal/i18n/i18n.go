package i18n

import (
	"fmt"
	"net/http"
)

type Locale string

const (
	En Locale = "en"
	Fr Locale = "fr"
	It Locale = "it"
)

type Key string

type Bundle struct {
	messages map[Locale]map[Key]string
}

func NewBundle() *Bundle {
	b := &Bundle{messages: map[Locale]map[Key]string{
		En: en,
		Fr: fr,
		It: it,
	}}
	return b
}

var Global = NewBundle()

func (b *Bundle) T(locale Locale, key Key, args ...any) string {
	m, ok := b.messages[locale]
	if !ok {
		m = b.messages[En]
	}
	pattern, ok := m[key]
	if !ok {
		pattern = string(key)
	}
	if len(args) > 0 {
		return fmt.Sprintf(pattern, args...)
	}
	return pattern
}

func T(locale Locale, key Key, args ...any) string {
	return Global.T(locale, key, args...)
}

const LocaleCookie = "lang"

func DetectLocale(r *http.Request) Locale {
	if c, err := r.Cookie(LocaleCookie); err == nil {
		switch c.Value {
		case "en":
			return En
		case "fr":
			return Fr
		case "it":
			return It
		}
	}
	return Global.DetectLocale(r.Header.Get("Accept-Language"))
}

func (b *Bundle) DetectLocale(acceptLang string) Locale {
	if len(acceptLang) >= 2 {
		switch acceptLang[:2] {
		case "fr":
			return Fr
		case "it":
			return It
		}
	}
	return En
}
