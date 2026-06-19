package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/jlettori/piponews/internal/auth"
	"github.com/jlettori/piponews/internal/db"
	"github.com/jlettori/piponews/internal/i18n"
	"github.com/jlettori/piponews/internal/models"
	"github.com/jlettori/piponews/internal/templates"
)

const entryPageSize = 10

type EntriesHandler struct {
	DB *db.DB
}

var entriesTemplate = template.Must(
	template.New("entries_list.html").Funcs(template.FuncMap{
		"dict":       dict,
		"formatTime": formatTime,
		"safeHTML":   safeHTML,
		"T":          i18n.T,
		"version":    version,
	}).ParseFS(templates.FS, "entries_list.html", "entry_card.html"),
)

var moreEntriesTemplate = template.Must(
	template.New("more_entries.html").Funcs(template.FuncMap{
		"dict":       dict,
		"formatTime": formatTime,
		"safeHTML":   safeHTML,
		"T":          i18n.T,
		"version":    version,
	}).ParseFS(templates.FS, "more_entries.html", "entry_card.html"),
)

type entryRow struct {
	models.Entry
	FeedTitle string `db:"feed_title"`
}

// scanEntry scans an entry row from a query result
func scanEntry(scanner interface {
	Scan(dest ...any) error
}, e *entryRow) error {
	return scanner.Scan(
		&e.ID, &e.FeedID, &e.GUID, &e.Title, &e.Summary,
		&e.URL,
		&e.PublishedAt,
		&e.IsSelected,
		&e.FeedTitle,
	)
}

func (h *EntriesHandler) List(w http.ResponseWriter, r *http.Request) {
	locale := detectLocale(r)
	userID := auth.GetUserID(r)

	var feedFilter, dateFrom, dateTo, searchQuery string
	ds := r.URL.Query().Get("datastar")
	if ds != "" {
		var signals struct {
			FeedFilter  string `json:"feedFilter"`
			DateFrom    string `json:"dateFrom"`
			DateTo      string `json:"dateTo"`
			SearchQuery string `json:"searchQuery"`
		}
		if err := json.Unmarshal([]byte(ds), &signals); err == nil {
			feedFilter, dateFrom, dateTo, searchQuery = signals.FeedFilter, signals.DateFrom, signals.DateTo, signals.SearchQuery
		}
	} else {
		feedFilter = r.URL.Query().Get("feedFilter")
		dateFrom = r.URL.Query().Get("dateFrom")
		dateTo = r.URL.Query().Get("dateTo")
		searchQuery = r.URL.Query().Get("searchQuery")
	}

	entries := FetchEntries(h.DB, userID, feedFilter, dateFrom, dateTo, searchQuery, entryPageSize, 0)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	RenderEntriesHTML(w, entries, locale)
}

func (h *EntriesHandler) More(w http.ResponseWriter, r *http.Request) {
	locale := detectLocale(r)
	userID := auth.GetUserID(r)

	offset := 0
	var feedFilter, dateFrom, dateTo, searchQuery string
	ds := r.URL.Query().Get("datastar")
	if ds != "" {
		var signals struct {
			Offset      int    `json:"offset"`
			FeedFilter  string `json:"feedFilter"`
			DateFrom    string `json:"dateFrom"`
			DateTo      string `json:"dateTo"`
			SearchQuery string `json:"searchQuery"`
		}
		if err := json.Unmarshal([]byte(ds), &signals); err == nil {
			offset = signals.Offset
			feedFilter, dateFrom, dateTo, searchQuery = signals.FeedFilter, signals.DateFrom, signals.DateTo, signals.SearchQuery
		}
	}

	entries := FetchEntries(h.DB, userID, feedFilter, dateFrom, dateTo, searchQuery, entryPageSize, offset)
	hasMore := len(entries) >= entryPageSize

	var buf bytes.Buffer
	moreEntriesTemplate.Execute(&buf, map[string]any{
		"Entries": entries,
		"Locale":  locale,
	})

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")

	writeSSEEvent(w, "datastar-patch-elements", map[string]string{
		"selector": "#entries-list",
		"mode":     "append",
		"elements": buf.String(),
	})

	if hasMore {
		writeSSEEvent(w, "datastar-patch-signals", map[string]string{
			"signals": fmt.Sprintf(`{"offset":%d}`, offset+entryPageSize),
		})
	} else {
		writeSSEEvent(w, "datastar-patch-elements", map[string]string{
			"selector": "#entries-sentinel",
			"mode":     "remove",
		})
	}
}

func writeSSEEvent(w http.ResponseWriter, event string, data map[string]string) {
	fmt.Fprintf(w, "event: %s\n", event)
	for k, v := range data {
		for _, line := range bytes.Split([]byte(v), []byte("\n")) {
			fmt.Fprintf(w, "data: %s %s\n", k, line)
		}
	}
	fmt.Fprintf(w, "\n")
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

func (h *EntriesHandler) ToggleSelect(w http.ResponseWriter, r *http.Request) {
	locale := detectLocale(r)
	userID := auth.GetUserID(r)
	entryID := r.PathValue("id")
	h.DB.Exec(`
		UPDATE entries SET is_selected = CASE WHEN is_selected THEN 0 ELSE 1 END
		WHERE id = ? AND feed_id IN (SELECT id FROM feeds WHERE user_id = ?)
	`, entryID, userID)
	entries := FetchEntries(h.DB, userID, "", "", "", "", entryPageSize, 0)
	var feeds []models.Feed
	h.DB.Select(&feeds, "SELECT * FROM feeds WHERE user_id = ? ORDER BY title COLLATE NOCASE ASC", userID)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	selCount := 0
	for _, e := range entries {
		if e.IsSelected {
			selCount++
		}
	}
	filterBarTemplate.Execute(w, map[string]any{"Feeds": feeds, "Locale": locale, "SelectedCount": selCount})
	RenderEntriesHTML(w, entries, locale)
}

func (h *EntriesHandler) ClearSelection(w http.ResponseWriter, r *http.Request) {
	locale := detectLocale(r)
	userID := auth.GetUserID(r)
	h.DB.Exec(`
		UPDATE entries SET is_selected = 0
		WHERE feed_id IN (SELECT id FROM feeds WHERE user_id = ?)
	`, userID)
	entries := FetchEntries(h.DB, userID, "", "", "", "", entryPageSize, 0)
	var feeds []models.Feed
	h.DB.Select(&feeds, "SELECT * FROM feeds WHERE user_id = ? ORDER BY title COLLATE NOCASE ASC", userID)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	selCount := 0
	filterBarTemplate.Execute(w, map[string]any{"Feeds": feeds, "Locale": locale, "SelectedCount": selCount})
	RenderEntriesHTML(w, entries, locale)
}

// FetchEntries is shared between FeedsHandler and EntriesHandler
func FetchEntries(database *db.DB, userID int64, feedFilter, dateFrom, dateTo, searchQuery string, limit, offset int) []entryRow {
	var entries []entryRow

	query := `
		SELECT e.*, f.title as feed_title
		FROM entries e
		JOIN feeds f ON f.id = e.feed_id
		WHERE f.user_id = ?`
	args := []any{userID}

	if feedFilter != "" {
		query += " AND e.feed_id = ?"
		args = append(args, feedFilter)
	}
	if dateFrom != "" {
		query += " AND e.published_at >= ?"
		args = append(args, dateFrom)
	}
	if dateTo != "" {
		query += " AND e.published_at <= ?"
		args = append(args, dateTo+"T23:59:59")
	}
	if searchQuery != "" {
		query += " AND (e.title LIKE ? OR e.summary LIKE ?)"
		like := "%" + searchQuery + "%"
		args = append(args, like, like)
	}
	query += " ORDER BY e.published_at DESC"
	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}
	if offset > 0 {
		query += " OFFSET ?"
		args = append(args, offset)
	}

	rows, err := database.Query(query, args...)
	if err != nil {
		return entries
	}
	defer rows.Close()

	for rows.Next() {
		var e entryRow
		if err := scanEntry(rows, &e); err != nil {
			continue
		}
		entries = append(entries, e)
	}
	return entries
}

func formatTime(t *time.Time, locale i18n.Locale) string {
	if t == nil {
		return ""
	}
	switch locale {
	case i18n.Fr, i18n.It:
		return t.Format("2 Jan 2006")
	default:
		return t.Format("Jan 2, 2006")
	}
}

func RenderEntriesHTML(w http.ResponseWriter, entries []entryRow, locale i18n.Locale) {
	msg := i18n.T(locale, i18n.NoEntriesMatch)
	if len(entries) > 0 {
		msg = ""
	}
	hasMore := len(entries) >= entryPageSize
	selCount := 0
	for _, e := range entries {
		if e.IsSelected {
			selCount++
		}
	}
	entriesTemplate.Execute(w, map[string]any{
		"Entries":       entries,
		"EmptyMessage":  msg,
		"Locale":        locale,
		"HasMore":       hasMore,
		"SelectedCount": selCount,
	})
}
