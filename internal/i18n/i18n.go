package i18n

import (
	"fmt"
	"net/http"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Locale string

const (
	En Locale = "en"
	Fr Locale = "fr"
	It Locale = "it"
)

type Key string

var matcher = language.NewMatcher([]language.Tag{
	language.English,
	language.French,
	language.Italian,
})

var pluralKeys = map[Key]bool{}

func pluralKeysInit() {
	pluralKeys[SelectedCount] = true
}

func localeToTag(l Locale) language.Tag {
	switch l {
	case Fr:
		return language.French
	case It:
		return language.Italian
	default:
		return language.English
	}
}

func tagToLocale(t language.Tag) Locale {
	base, _ := t.Base()
	switch base.String() {
	case "fr":
		return Fr
	case "it":
		return It
	default:
		return En
	}
}

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

func (b *Bundle) N(locale Locale, key Key, count int, args ...any) string {
	m, ok := b.messages[locale]
	if !ok {
		m = b.messages[En]
	}

	pattern, ok := m[key]
	if !ok {
		pattern = string(key)
	}

	if pluralKeys[key] && count == 1 {
		oneKey := Key(string(key) + "_one")
		if p, ok := m[oneKey]; ok {
			pattern = p
		}
	}

	tag := localeToTag(locale)
	p := message.NewPrinter(tag)
	if pluralKeys[key] {
		return p.Sprintf(pattern, append([]any{count}, args...)...)
	}
	return p.Sprintf(pattern, args...)
}

func N(locale Locale, key Key, count int, args ...any) string {
	return Global.N(locale, key, count, args...)
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
	if acceptLang == "" {
		return En
	}
	tag, _ := language.MatchStrings(matcher, acceptLang)
	return tagToLocale(tag)
}

func init() {
	pluralKeysInit()
}
