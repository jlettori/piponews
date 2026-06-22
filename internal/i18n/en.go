package i18n

var en = map[Key]string{
	// Auth
	SignIn:                   "Sign in",
	Register:                 "Register",
	Username:                 "Username",
	Password:                 "Password",
	DontHaveAccount:          "Don't have an account?",
	AlreadyHaveAccount:       "Already have an account?",
	InvalidFormData:          "Invalid form data",
	UsernamePasswordRequired: "Username and password are required",
	InvalidUsernamePassword:  "Invalid username or password",
	InternalError:            "Internal error",
	UsernameLength:           "Username must be between 2 and 50 characters",
	UsernameChars:            "Username may only contain letters, digits, underscores, and hyphens",
	PasswordLength:           "Password must be at least 4 characters",
	UsernameTaken:            "Username already taken",
	RememberMe:               "Keep me signed in",

	// Layout
	AppName:     "piponews",
	Logout:      "Logout",
	EditProfile: "Edit profile",
	Language:    "Language",
	LangEN:      "English",
	LangFR:      "Français",
	LangIT:      "Italiano",

	// Feeds
	Feeds:              "Feeds",
	FeedURLPlaceholder: "Feed URL...",
	AllFeeds:           "All feeds",
	NoFeedsYet:         "No feeds yet. Add one above!",
	Refresh:            "Refresh",
	Remove:             "Remove",
	RefreshAll:         "Refresh all",
	FeedAlreadyAdded:   "Feed already added",
	FailedToSaveFeed:   "Failed to save feed",
	URLRequired:        "URL is required",
	BadRequest:         "Bad request",
	FailedParseFeed:    "Failed to parse feed. Check the URL and try again.",
	InvalidFeedURL:     "Invalid feed URL: %s",

	// Filters
	Search: "Search...",
	From:   "From",
	To:     "To",

	// Entries
	ToggleSelect:     "Toggle select",
	ClearSelection:   "Clear selection",
	SelectedCount:    "%d selected entries",
	SelectedCountOne: "1 selected entry",
	ExportSelected:   "Export selected",
	NoEntriesMatch:   "No entries match the current filters.",
	NoEntriesYet:     "No entries yet. Add a feed!",
	ExportSelEntries: "Export selected entries",

	// Infinite scroll
	Loading:   "Loading...",
	AllLoaded: "All entries loaded.",

	// Export dialog
	CSV:                "CSV",
	Format:             "Format",
	HTML:               "HTML",
	PlainText:          "Plain text",
	Cancel:             "Cancel",
	Export:             "Export",
	FailedQueryEntries: "Failed to query entries",

	// Export output
	SelEntriesTitle: "%s — Selected Entries",
	ExportedOn:      "Exported on %s",
	NoEntriesSel:    "No entries selected.",
	TxTHeader:       "%s - Selected Entries",
	TitleLabel:      "Title:",
	FeedLabel:       "Feed:",
	URLLabel:        "URL:",
	DateLabel:       "Date:",

	// Profile
	Profile:             "Profile",
	FirstName:           "First name",
	LastName:            "Last name",
	Email:               "Email",
	Save:                "Save",
	ProfileUpdated:      "Profile updated",
	ProfileUpdateFailed: "Failed to update profile",
	EmailInvalid:        "Invalid email address",

	// Validation
	InvalidURL:    "invalid URL",
	HTTPHTTPSOnly: "only http and https URLs are allowed",
	URLHostReq:    "URL must have a host",
}
