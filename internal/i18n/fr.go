package i18n

var fr = map[Key]string{
	// Auth
	SignIn:                   "Connexion",
	Register:                 "S'inscrire",
	Username:                 "Nom d'utilisateur",
	Password:                 "Mot de passe",
	DontHaveAccount:          "Vous n'avez pas de compte ?",
	AlreadyHaveAccount:       "Vous avez déjà un compte ?",
	InvalidFormData:          "Données de formulaire invalides",
	UsernamePasswordRequired: "Nom d'utilisateur et mot de passe requis",
	InvalidUsernamePassword:  "Nom d'utilisateur ou mot de passe invalide",
	InternalError:            "Erreur interne",
	UsernameLength:           "Le nom d'utilisateur doit contenir entre 2 et 50 caractères",
	UsernameChars:            "Le nom d'utilisateur ne peut contenir que des lettres, chiffres, tirets et underscores",
	PasswordLength:           "Le mot de passe doit contenir au moins 4 caractères",
	UsernameTaken:            "Ce nom d'utilisateur est déjà pris",

	// Layout
	AppName:     "piponews",
	Logout:      "Déconnexion",
	EditProfile: "Modifier le profil",
	Language:    "Langue",
	LangEN:      "English",
	LangFR:      "Français",
	LangIT:      "Italiano",

	// Feeds
	Feeds:              "Flux",
	FeedURLPlaceholder: "URL du flux...",
	AllFeeds:           "Tous les flux",
	NoFeedsYet:         "Aucun flux. Ajoutez-en un ci-dessus !",
	Refresh:            "Actualiser",
	Remove:             "Supprimer",
	RefreshAll:         "Tout actualiser",
	FeedAlreadyAdded:   "Ce flux a déjà été ajouté",
	FailedToSaveFeed:   "Échec de l'enregistrement du flux",
	URLRequired:        "L'URL est requise",
	BadRequest:         "Requête invalide",
	FailedParseFeed:    "Échec de l'analyse du flux. Vérifiez l'URL et réessayez.",
	InvalidFeedURL:     "URL de flux invalide : %s",

	// Filters
	From: "Du",
	To:   "Au",

	// Entries
	ToggleSelect:     "Sélectionner/Désélectionner",
	ClearSelection:   "Effacer la sélection",
	SelectedCount:    "%d articles sélectionnés",
	ExportSelected:   "Exporter la sélection",
	NoEntriesMatch:   "Aucun article ne correspond aux filtres actuels.",
	NoEntriesYet:     "Aucun article. Ajoutez un flux !",
	ExportSelEntries: "Exporter les articles sélectionnés",

	// Infinite scroll
	Loading:   "Chargement...",
	AllLoaded: "Tous les articles chargés.",

	// Export dialog
	CSV:                "CSV",
	Format:             "Format",
	HTML:               "HTML",
	PlainText:          "Texte brut",
	Cancel:             "Annuler",
	Export:             "Exporter",
	FailedQueryEntries: "Échec de la requête des articles",

	// Export output
	SelEntriesTitle: "%s — Articles sélectionnés",
	ExportedOn:      "Exporté le %s",
	NoEntriesSel:    "Aucun article sélectionné.",
	TxTHeader:       "%s - Articles sélectionnés",
	TitleLabel:      "Titre :",
	FeedLabel:       "Flux :",
	URLLabel:        "URL :",
	DateLabel:       "Date :",

	// Profile
	Profile:             "Profil",
	FirstName:           "Prénom",
	LastName:            "Nom",
	Email:               "E-mail",
	Save:                "Enregistrer",
	ProfileUpdated:      "Profil mis à jour",
	ProfileUpdateFailed: "Échec de la mise à jour du profil",
	EmailInvalid:        "Adresse e-mail invalide",

	// Validation
	InvalidURL:    "URL invalide",
	HTTPHTTPSOnly: "seules les URLs http et https sont autorisées",
	URLHostReq:    "L'URL doit avoir un hôte",
}
