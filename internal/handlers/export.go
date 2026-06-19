package handlers

import (
	"database/sql"
	"encoding/csv"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jlettori/piponews/internal/auth"
	"github.com/jlettori/piponews/internal/db"
	"github.com/jlettori/piponews/internal/i18n"
	"github.com/jlettori/piponews/internal/templates"
	texttemplate "text/template"
)

type ExportsHandler struct {
	DB *db.DB
}

type exportEntry struct {
	ID          int64
	FeedID      int64
	FeedTitle   string
	GUID        string
	Title       string
	URL         string
	Summary     string
	PublishedAt *time.Time
	IsSelected  bool
}

var exportHTMLTemplate = template.Must(
	template.New("export.html").Funcs(templateFuncs()).ParseFS(templates.FS, "export.html"),
)

var exportTxtTemplate = texttemplate.Must(
	texttemplate.New("export.txt").Funcs(textTemplateFuncs()).ParseFS(templates.FS, "export.txt"),
)

func (h *ExportsHandler) Export(w http.ResponseWriter, r *http.Request) {
	locale := detectLocale(r)
	userID := auth.GetUserID(r)

	format := "html"
	if strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		var signals struct {
			ExportFormat string `json:"exportFormat"`
		}
		if err := parseSignals(r, &signals); err == nil {
			format = signals.ExportFormat
		}
	} else {
		format = r.FormValue("exportFormat")
	}

	format = strings.TrimSpace(format)

	rows, err := h.DB.Query(`
		SELECT e.id, e.feed_id, e.guid, e.title, e.url, e.summary, e.published_at, e.is_selected, f.title as feed_title
		FROM entries e
		JOIN feeds f ON f.id = e.feed_id
		WHERE e.is_selected = 1 AND f.user_id = ?
		ORDER BY e.published_at DESC
	`, userID)
	if err != nil {
		ErrResponse(w, i18n.T(locale, i18n.FailedQueryEntries), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var entries []exportEntry
	for rows.Next() {
		var e exportEntry
		var pubAt sql.NullTime
		if err := rows.Scan(&e.ID, &e.FeedID, &e.GUID, &e.Title, &e.URL, &e.Summary, &pubAt, &e.IsSelected, &e.FeedTitle); err != nil {
			continue
		}
		if pubAt.Valid {
			e.PublishedAt = &pubAt.Time
		}
		entries = append(entries, e)
	}

	date := time.Now().Format("Jan 2, 2006 15:04")

	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	if format == "csv" {
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Content-Disposition", "attachment; filename=entries.csv")
		wr := csv.NewWriter(w)
		wr.Write([]string{"ID", "Feed ID", "Feed Title", "GUID", "Title", "URL", "Published At", "Summary", "Is Selected"})
		for _, e := range entries {
			pubAt := ""
			if e.PublishedAt != nil {
				pubAt = e.PublishedAt.Format(time.RFC3339)
			}
			wr.Write([]string{
				strconv.FormatInt(e.ID, 10),
				strconv.FormatInt(e.FeedID, 10),
				e.FeedTitle,
				e.GUID,
				e.Title,
				e.URL,
				pubAt,
				stripTags(e.Summary),
				strconv.FormatBool(e.IsSelected),
			})
		}
		wr.Flush()
		if err := wr.Error(); err != nil {
			http.Error(w, i18n.T(locale, i18n.InternalError), http.StatusInternalServerError)
		}
		return
	}

	if format == "txt" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		exportTxtTemplate.Execute(w, map[string]any{
			"Locale":  locale,
			"Date":    date,
			"Entries": entries,
		})
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	exportHTMLTemplate.Execute(w, map[string]any{
		"Locale":  locale,
		"Date":    date,
		"Entries": entries,
	})
}
