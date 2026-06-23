package i18n

var it = map[Key]string{
	// Auth
	SignIn:                   "Accedi",
	Register:                 "Registrati",
	Username:                 "Nome utente",
	Password:                 "Password",
	DontHaveAccount:          "Non hai un account?",
	AlreadyHaveAccount:       "Hai già un account?",
	InvalidFormData:          "Dati del modulo non validi",
	UsernamePasswordRequired: "Nome utente e password richiesti",
	InvalidUsernamePassword:  "Nome utente o password non validi",
	InternalError:            "Errore interno",
	UsernameLength:           "Il nome utente deve essere tra 2 e 50 caratteri",
	UsernameChars:            "Il nome utente può contenere solo lettere, numeri, trattini e underscore",
	PasswordLength:           "La password deve essere di almeno 4 caratteri",
	UsernameTaken:            "Nome utente già in uso",
	RememberMe:               "Mantieni l'accesso",

	// Layout
	AppName:     "piponews",
	Logout:      "Esci",
	EditProfile: "Modifica profilo",
	Language:    "Lingua",
	LangEN:      "Inglese",
	LangFR:      "Francese",
	LangIT:      "Italiano",

	// Feeds
	Feeds:              "Feed",
	FeedURLPlaceholder: "URL del feed...",
	AllFeeds:           "Tutti i feed",
	NoFeedsYet:         "Nessun feed. Aggiungine uno sopra!",
	Refresh:            "Aggiorna",
	Remove:             "Rimuovi",
	RefreshAll:         "Aggiorna tutto",
	FeedAlreadyAdded:   "Feed già aggiunto",
	FailedToSaveFeed:   "Impossibile salvare il feed",
	URLRequired:        "L'URL è richiesto",
	BadRequest:         "Richiesta non valida",
	FailedParseFeed:    "Impossibile analizzare il feed. Controlla l'URL e riprova.",
	InvalidFeedURL:     "URL del feed non valido: %s",

	// Filters
	Search: "Cerca...",
	From:   "Da",
	To:     "A",

	// Entries
	ToggleSelect:     "Seleziona/Deseleziona",
	ClearSelection:   "Cancella selezione",
	SelectedCount:    "%d articoli selezionati",
	SelectedCountOne: "1 articolo selezionato",
	ExportSelected:   "Esporta selezionati",
	NoEntriesMatch:   "Nessun articolo corrisponde ai filtri correnti.",
	NoEntriesYet:     "Nessun articolo. Aggiungi un feed!",
	ExportSelEntries: "Esporta articoli selezionati",

	// Infinite scroll
	Loading:   "Caricamento...",
	AllLoaded: "Tutti gli articoli caricati.",

	// Export dialog
	CSV:                "CSV",
	Format:             "Formato",
	HTML:               "HTML",
	PlainText:          "Testo semplice",
	Cancel:             "Annulla",
	Export:             "Esporta",
	FailedQueryEntries: "Impossibile interrogare gli articoli",

	// Export output
	SelEntriesTitle: "%s — Articoli selezionati",
	ExportedOn:      "Esportato il %s",
	NoEntriesSel:    "Nessun articolo selezionato.",
	TxTHeader:       "%s - Articoli selezionati",
	TitleLabel:      "Titolo:",
	FeedLabel:       "Feed:",
	URLLabel:        "URL:",
	DateLabel:       "Data:",

	// Profile
	Profile:             "Profilo",
	FirstName:           "Nome",
	LastName:            "Cognome",
	Email:               "Email",
	Save:                "Salva",
	ProfileUpdated:      "Profilo aggiornato",
	ProfileUpdateFailed: "Aggiornamento del profilo non riuscito",
	EmailInvalid:        "Indirizzo email non valido",

	// Validation
	InvalidURL:    "URL non valido",
	HTTPHTTPSOnly: "sono consentiti solo URL http e https",
	URLHostReq:    "L'URL deve avere un host",
}
