package i18n

import (
	"net/http"
	"testing"
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
		From, To,
		ToggleSelect, ClearSelection, SelectedCount, ExportSelected,
		NoEntriesMatch, NoEntriesYet, ExportSelEntries,
		Loading, AllLoaded,
		Format, HTML, PlainText, Cancel, Export, FailedQueryEntries,
		SelEntriesTitle, ExportedOn, NoEntriesSel, TxTHeader,
		TitleLabel, FeedLabel, URLLabel, DateLabel,
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
