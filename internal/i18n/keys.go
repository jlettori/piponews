package i18n

// Auth
const (
	SignIn                   Key = "sign_in"
	Register                 Key = "register"
	Username                 Key = "username"
	Password                 Key = "password"
	DontHaveAccount          Key = "dont_have_account"
	AlreadyHaveAccount       Key = "already_have_account"
	InvalidFormData          Key = "invalid_form_data"
	UsernamePasswordRequired Key = "username_password_required"
	InvalidUsernamePassword  Key = "invalid_username_password"
	InternalError            Key = "internal_error"
	UsernameLength           Key = "username_length"
	UsernameChars            Key = "username_chars"
	PasswordLength           Key = "password_length"
	UsernameTaken            Key = "username_taken"
)

// Layout
const (
	AppName     Key = "app_name"
	Logout      Key = "logout"
	EditProfile Key = "edit_profile"
	Language    Key = "language"
	LangEN      Key = "lang_en"
	LangFR      Key = "lang_fr"
	LangIT      Key = "lang_it"
)

// Feeds
const (
	Feeds              Key = "feeds"
	FeedURLPlaceholder Key = "feed_url_placeholder"
	AllFeeds           Key = "all_feeds"
	NoFeedsYet         Key = "no_feeds_yet"
	Refresh            Key = "refresh"
	Remove             Key = "remove"
	RefreshAll         Key = "refresh_all"
	FeedAlreadyAdded   Key = "feed_already_added"
	FailedToSaveFeed   Key = "failed_to_save_feed"
	URLRequired        Key = "url_required"
	BadRequest         Key = "bad_request"
	FailedParseFeed    Key = "failed_parse_feed"
	InvalidFeedURL     Key = "invalid_feed_url"
)

// Filters
const (
	Search Key = "search"
	From   Key = "from"
	To     Key = "to"
)

// Entries
const (
	ToggleSelect     Key = "toggle_select"
	ClearSelection   Key = "clear_selection"
	SelectedCount    Key = "selected_count"
	SelectedCountOne Key = "selected_count_one"
	ExportSelected   Key = "export_selected"
	NoEntriesMatch   Key = "no_entries_match"
	NoEntriesYet     Key = "no_entries_yet"
	ExportSelEntries Key = "export_sel_entries"
)

// Infinite scroll
const (
	Loading   Key = "loading"
	AllLoaded Key = "all_loaded"
)

// Export dialog
const (
	CSV                Key = "csv"
	Format             Key = "format"
	HTML               Key = "html"
	PlainText          Key = "plain_text"
	Cancel             Key = "cancel"
	Export             Key = "export"
	FailedQueryEntries Key = "failed_query_entries"
)

// Export output
const (
	SelEntriesTitle Key = "sel_entries_title"
	ExportedOn      Key = "exported_on"
	NoEntriesSel    Key = "no_entries_sel"
	TxTHeader       Key = "txt_header"
	TitleLabel      Key = "title_label"
	FeedLabel       Key = "feed_label"

	URLLabel  Key = "url_label"
	DateLabel Key = "date_label"
)

// Profile
const (
	Profile             Key = "profile"
	FirstName           Key = "first_name"
	LastName            Key = "last_name"
	Email               Key = "email"
	Save                Key = "save"
	ProfileUpdated      Key = "profile_updated"
	ProfileUpdateFailed Key = "profile_update_failed"
	EmailInvalid        Key = "email_invalid"
)

// Validation
const (
	InvalidURL    Key = "invalid_url"
	HTTPHTTPSOnly Key = "http_https_only"
	URLHostReq    Key = "url_host_required"
)
