package i18n

import (
	"net/http"
	"testing"

	"golang.org/x/text/language"
)

func TestT_ExistingKey(t *testing.T) {
	got := T(En, SignIn)
	if got != "Sign in" {
		t.Errorf("T(En, SignIn) = %q; want %q", got, "Sign in")
	}
}

func TestT_French(t *testing.T) {
	got := T(Fr, SignIn)
	if got != "Connexion" {
		t.Errorf("T(Fr, SignIn) = %q; want %q", got, "Connexion")
	}
}

func TestT_Italian(t *testing.T) {
	got := T(It, SignIn)
	if got != "Accedi" {
		t.Errorf("T(It, SignIn) = %q; want %q", got, "Accedi")
	}
}

func TestT_MissingKeyReturnsKeyItself(t *testing.T) {
	got := T(En, Key("nonexistent_key"))
	if got != "nonexistent_key" {
		t.Errorf("T(En, nonexistent) = %q; want %q", got, "nonexistent_key")
	}
}

func TestT_UnknownLocaleFallsBackToEnglish(t *testing.T) {
	got := T(Locale("de"), SignIn)
	if got != "Sign in" {
		t.Errorf("T(de, SignIn) = %q; want %q", got, "Sign in")
	}
}

func TestT_WithArgs(t *testing.T) {
	got := T(En, InvalidFeedURL, "bad url")
	if got != "Invalid feed URL: bad url" {
		t.Errorf("T(En, InvalidFeedURL, ...) = %q; want %q", got, "Invalid feed URL: bad url")
	}
}

func TestT_FrenchWithArgs(t *testing.T) {
	got := T(Fr, InvalidFeedURL, "mauvaise url")
	if got != "URL de flux invalide : mauvaise url" {
		t.Errorf("T(Fr, InvalidFeedURL, ...) = %q; want %q", got, "URL de flux invalide : mauvaise url")
	}
}

func TestT_WithMultipleArgs(t *testing.T) {
	got := T(En, SelEntriesTitle, "My Feed")
	want := "My Feed — Selected Entries"
	if got != want {
		t.Errorf("T(En, SelEntriesTitle, ...) = %q; want %q", got, want)
	}
}

func TestT_WithNoArgsPreservesPattern(t *testing.T) {
	got := T(En, AppName)
	if got != "piponews" {
		t.Errorf("T(En, AppName) = %q; want %q", got, "piponews")
	}
}

func TestBundleT_ReturnsPatternForMissingKey(t *testing.T) {
	b := NewBundle()
	got := b.T(En, Key("not_in_map"))
	if got != "not_in_map" {
		t.Errorf("b.T(En, not_in_map) = %q; want %q", got, "not_in_map")
	}
}

func TestBundleT_FallsBackToEnglish(t *testing.T) {
	b := NewBundle()
	got := b.T(Locale("de"), SignIn)
	if got != "Sign in" {
		t.Errorf("b.T(de, SignIn) = %q; want %q", got, "Sign in")
	}
}

func TestBundleT_WithArgs(t *testing.T) {
	b := NewBundle()
	got := b.T(En, InvalidFeedURL, "bad")
	if got != "Invalid feed URL: bad" {
		t.Errorf("b.T(En, InvalidFeedURL, ...) = %q; want %q", got, "Invalid feed URL: bad")
	}
}

func TestGlobalIsPopulated(t *testing.T) {
	if Global == nil {
		t.Fatal("Global is nil")
	}
	got := Global.T(En, AppName)
	if got != "piponews" {
		t.Errorf("Global.T(En, AppName) = %q; want %q", got, "piponews")
	}
}

func TestN_EnglishSingular(t *testing.T) {
	got := N(En, SelectedCount, 1)
	if got != "1 selected entry" {
		t.Errorf("N(En, SelectedCount, 1) = %q; want %q", got, "1 selected entry")
	}
}

func TestN_EnglishPlural(t *testing.T) {
	got := N(En, SelectedCount, 5)
	if got != "5 selected entries" {
		t.Errorf("N(En, SelectedCount, 5) = %q; want %q", got, "5 selected entries")
	}
}

func TestN_FrenchSingular(t *testing.T) {
	got := N(Fr, SelectedCount, 1)
	if got != "1 article sélectionné" {
		t.Errorf("N(Fr, SelectedCount, 1) = %q; want %q", got, "1 article sélectionné")
	}
}

func TestN_FrenchPlural(t *testing.T) {
	got := N(Fr, SelectedCount, 3)
	if got != "3 articles sélectionnés" {
		t.Errorf("N(Fr, SelectedCount, 3) = %q; want %q", got, "3 articles sélectionnés")
	}
}

func TestN_ItalianSingular(t *testing.T) {
	got := N(It, SelectedCount, 1)
	if got != "1 articolo selezionato" {
		t.Errorf("N(It, SelectedCount, 1) = %q; want %q", got, "1 articolo selezionato")
	}
}

func TestN_ItalianPlural(t *testing.T) {
	got := N(It, SelectedCount, 7)
	if got != "7 articoli selezionati" {
		t.Errorf("N(It, SelectedCount, 7) = %q; want %q", got, "7 articoli selezionati")
	}
}

func TestN_EnglishZeroCount(t *testing.T) {
	got := N(En, SelectedCount, 0)
	if got != "0 selected entries" {
		t.Errorf("N(En, SelectedCount, 0) = %q; want %q", got, "0 selected entries")
	}
}

func TestN_NonPluralKeyReturnsPattern(t *testing.T) {
	got := N(En, Feeds, 5)
	if got != "Feeds" {
		t.Errorf("N(En, Feeds, 5) = %q; want %q", got, "Feeds")
	}
}

func TestN_MissingKeyReturnsKeyItself(t *testing.T) {
	got := N(En, Key("nonexistent_key"), 3)
	if got != "nonexistent_key" {
		t.Errorf("N(En, nonexistent_key, 3) = %q; want %q", got, "nonexistent_key")
	}
}

func TestN_UnknownLocaleFallsBackToEnglish(t *testing.T) {
	got := N(Locale("de"), SelectedCount, 1)
	if got != "1 selected entry" {
		t.Errorf("N(de, SelectedCount, 1) = %q; want %q", got, "1 selected entry")
	}
}

func TestN_WithFormatPattern(t *testing.T) {
	got := N(En, InvalidFeedURL, 3, "arg")
	want := "Invalid feed URL: arg"
	if got != want {
		t.Errorf("N(En, InvalidFeedURL, 3, ...) = %q; want %q", got, want)
	}
}

func TestBundleN_Singular(t *testing.T) {
	b := NewBundle()
	got := b.N(En, SelectedCount, 1)
	if got != "1 selected entry" {
		t.Errorf("b.N(En, SelectedCount, 1) = %q; want %q", got, "1 selected entry")
	}
}

func TestBundleN_Plural(t *testing.T) {
	b := NewBundle()
	got := b.N(En, SelectedCount, 10)
	if got != "10 selected entries" {
		t.Errorf("b.N(En, SelectedCount, 10) = %q; want %q", got, "10 selected entries")
	}
}

func TestBundleN_MissingKey(t *testing.T) {
	b := NewBundle()
	got := b.N(En, Key("missing"), 3)
	if got != "missing" {
		t.Errorf("b.N(En, missing, 3) = %q; want %q", got, "missing")
	}
}

func TestBundleN_UnknownLocale(t *testing.T) {
	b := NewBundle()
	got := b.N(Locale("de"), SelectedCount, 1)
	if got != "1 selected entry" {
		t.Errorf("b.N(de, SelectedCount, 1) = %q; want %q", got, "1 selected entry")
	}
}

func TestDetectLocale_French(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Language", "fr-FR,fr;q=0.9")
	got := DetectLocale(r)
	if got != Fr {
		t.Errorf("DetectLocale(fr-FR) = %q; want %q", got, Fr)
	}
}

func TestDetectLocale_Italian(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Language", "it-IT,it;q=0.9")
	got := DetectLocale(r)
	if got != It {
		t.Errorf("DetectLocale(it-IT) = %q; want %q", got, It)
	}
}

func TestDetectLocale_English(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Language", "en-US,en;q=0.9")
	got := DetectLocale(r)
	if got != En {
		t.Errorf("DetectLocale(en-US) = %q; want %q", got, En)
	}
}

func TestDetectLocale_Empty(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	got := DetectLocale(r)
	if got != En {
		t.Errorf("DetectLocale(empty) = %q; want %q", got, En)
	}
}

func TestDetectLocale_Unknown(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Language", "de-DE,de;q=0.9")
	got := DetectLocale(r)
	if got != En {
		t.Errorf("DetectLocale(de-DE) = %q; want %q", got, En)
	}
}

func TestDetectLocale_CookiePrecedesHeader(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: LocaleCookie, Value: "it"})
	r.Header.Set("Accept-Language", "fr-FR,fr;q=0.9")
	got := DetectLocale(r)
	if got != It {
		t.Errorf("DetectLocale with cookie=it and header=fr = %q; want %q", got, It)
	}
}

func TestDetectLocale_CookieInvalidValueFallsThrough(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: LocaleCookie, Value: "xx"})
	r.Header.Set("Accept-Language", "fr-FR,fr;q=0.9")
	got := DetectLocale(r)
	if got != Fr {
		t.Errorf("DetectLocale with cookie=xx and header=fr = %q; want %q", got, Fr)
	}
}

func TestDetectLocale_CookieItalian(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: LocaleCookie, Value: "it"})
	got := DetectLocale(r)
	if got != It {
		t.Errorf("DetectLocale with cookie=it = %q; want %q", got, It)
	}
}

func TestDetectLocale_CookieEnglish(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: LocaleCookie, Value: "en"})
	got := DetectLocale(r)
	if got != En {
		t.Errorf("DetectLocale with cookie=en = %q; want %q", got, En)
	}
}

func TestDetectLocale_CookieFrench(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)
	r.AddCookie(&http.Cookie{Name: LocaleCookie, Value: "fr"})
	got := DetectLocale(r)
	if got != Fr {
		t.Errorf("DetectLocale with cookie=fr = %q; want %q", got, Fr)
	}
}

func TestBundleDetectLocale_French(t *testing.T) {
	b := NewBundle()
	got := b.DetectLocale("fr-FR,fr;q=0.9")
	if got != Fr {
		t.Errorf("b.DetectLocale(fr-FR) = %q; want %q", got, Fr)
	}
}

func TestBundleDetectLocale_Empty(t *testing.T) {
	b := NewBundle()
	got := b.DetectLocale("")
	if got != En {
		t.Errorf("b.DetectLocale('') = %q; want %q", got, En)
	}
}

func TestBundleDetectLocale_Unknown(t *testing.T) {
	b := NewBundle()
	got := b.DetectLocale("de")
	if got != En {
		t.Errorf("b.DetectLocale(de) = %q; want %q", got, En)
	}
}

func TestLocaleToTag_English(t *testing.T) {
	tag := localeToTag(En)
	base, _ := tag.Base()
	if base.String() != "en" {
		t.Errorf("localeToTag(En) base = %q; want %q", base.String(), "en")
	}
}

func TestLocaleToTag_French(t *testing.T) {
	tag := localeToTag(Fr)
	base, _ := tag.Base()
	if base.String() != "fr" {
		t.Errorf("localeToTag(Fr) base = %q; want %q", base.String(), "fr")
	}
}

func TestLocaleToTag_Italian(t *testing.T) {
	tag := localeToTag(It)
	base, _ := tag.Base()
	if base.String() != "it" {
		t.Errorf("localeToTag(It) base = %q; want %q", base.String(), "it")
	}
}

func TestTagToLocale_English(t *testing.T) {
	got := tagToLocale(language.English)
	if got != En {
		t.Errorf("tagToLocale(English) = %q; want %q", got, En)
	}
}

func TestTagToLocale_French(t *testing.T) {
	got := tagToLocale(language.French)
	if got != Fr {
		t.Errorf("tagToLocale(French) = %q; want %q", got, Fr)
	}
}

func TestTagToLocale_Italian(t *testing.T) {
	got := tagToLocale(language.Italian)
	if got != It {
		t.Errorf("tagToLocale(Italian) = %q; want %q", got, It)
	}
}

func TestAllKeysPresentInAllLocales(t *testing.T) {
	allKeys := []Key{
		SignIn, Register, Username, Password,
		DontHaveAccount, AlreadyHaveAccount, InvalidFormData,
		UsernamePasswordRequired, InvalidUsernamePassword, InternalError,
		UsernameLength, UsernameChars, PasswordLength, UsernameTaken,
		AppName, Logout, EditProfile, Language, LangEN, LangFR, LangIT,
		Feeds, FeedURLPlaceholder, AllFeeds, NoFeedsYet,
		Refresh, Remove, RefreshAll, FeedAlreadyAdded,
		FailedToSaveFeed, URLRequired, BadRequest, FailedParseFeed, InvalidFeedURL,
		Search, From, To,
		ToggleSelect, ClearSelection, SelectedCount, SelectedCountOne, ExportSelected,
		NoEntriesMatch, NoEntriesYet, ExportSelEntries,
		Loading, AllLoaded,
		CSV, Format, HTML, PlainText, Cancel, Export, FailedQueryEntries,
		SelEntriesTitle, ExportedOn, NoEntriesSel, TxTHeader,
		TitleLabel, FeedLabel, URLLabel, DateLabel,
		Profile, FirstName, LastName, Email, Save,
		ProfileUpdated, ProfileUpdateFailed, EmailInvalid,
		InvalidURL, HTTPHTTPSOnly, URLHostReq,
	}
	locales := []Locale{En, Fr, It}
	for _, locale := range locales {
		for _, key := range allKeys {
			got := T(locale, key)
			if got == string(key) {
				t.Errorf("locale %q missing key %q", locale, key)
			}
		}
	}
}
