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

func TestN_FallsBackForNonPluralKey(t *testing.T) {
	got := N(En, Feeds, 5)
	if got != "Feeds" {
		t.Errorf("N(En, Feeds, 5) = %q; want %q", got, "Feeds")
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
		ToggleSelect, ClearSelection, SelectedCount, SelectedCountOne, ExportSelected,
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
