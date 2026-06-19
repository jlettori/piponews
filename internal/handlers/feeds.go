package handlers

import (
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jlettori/piponews/internal/auth"
	"github.com/jlettori/piponews/internal/db"
	"github.com/jlettori/piponews/internal/i18n"
	"github.com/jlettori/piponews/internal/models"
	"github.com/jlettori/piponews/internal/templates"

	"github.com/mmcdole/gofeed"
)

type FeedsHandler struct {
	DB *db.DB
}

var feedsTemplate = template.Must(
	template.New("feed_list.html").Funcs(template.FuncMap{"T": i18n.T, "version": version}).ParseFS(templates.FS, "feed_list.html"),
)

var filterBarTemplate = template.Must(
	template.New("filter_bar.html").Funcs(template.FuncMap{"T": i18n.T, "version": version}).ParseFS(templates.FS, "filter_bar.html"),
)

func (h *FeedsHandler) List(w http.ResponseWriter, r *http.Request) {
	locale := detectLocale(r)
	userID := auth.GetUserID(r)
	feeds := h.getFeedsWithCounts(userID)
	entries := FetchEntries(h.DB, userID, "", "", "", "", entryPageSize, 0)

	selCount := 0
	for _, e := range entries {
		if e.IsSelected {
			selCount++
		}
	}
	data := map[string]any{
		"Locale":        locale,
		"Username":      auth.GetUsername(r),
		"Feeds":         feeds,
		"TotalEntries":  totalEntries(feeds),
		"Entries":       entries,
		"HasMore":       len(entries) >= entryPageSize,
		"SelectedCount": selCount,
	}

	t := parseTemplates("feeds.html")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.ExecuteTemplate(w, "base.html", data); err != nil {
		log.Printf("template error: %v", err)
	}
}

func (h *FeedsHandler) Create(w http.ResponseWriter, r *http.Request) {
	locale := detectLocale(r)
	userID := auth.GetUserID(r)

	var signals struct {
		URL string `json:"url"`
	}
	if err := parseSignals(r, &signals); err != nil {
		ErrResponse(w, i18n.T(locale, i18n.BadRequest), http.StatusBadRequest)
		return
	}

	feedURL := strings.TrimSpace(signals.URL)
	if feedURL == "" {
		ErrResponse(w, i18n.T(locale, i18n.URLRequired), http.StatusBadRequest)
		return
	}

	if err := validateFeedURL(feedURL); err != nil {
		ErrResponse(w, fmt.Sprintf(i18n.T(locale, i18n.InvalidFeedURL), err.Error()), http.StatusBadRequest)
		return
	}

	fp := gofeed.NewParser()
	fp.Client = httpClientWithTimeout(30 * time.Second)
	parsed, err := fp.ParseURL(feedURL)
	if err != nil {
		ErrResponse(w, i18n.T(locale, i18n.FailedParseFeed), http.StatusBadRequest)
		return
	}

	title := html.UnescapeString(parsed.Title)
	if title == "" {
		title = feedURL
	}

	res, err := h.DB.Exec(
		"INSERT INTO feeds (user_id, title, url) VALUES (?, ?, ?)",
		userID, title, feedURL,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			ErrResponse(w, i18n.T(locale, i18n.FeedAlreadyAdded), http.StatusConflict)
			return
		}
		ErrResponse(w, i18n.T(locale, i18n.FailedToSaveFeed), http.StatusInternalServerError)
		return
	}

	feedID, _ := res.LastInsertId()
	h.saveFeedEntries(feedID, parsed)
	h.DB.Exec("UPDATE feeds SET last_fetched_at = ? WHERE id = ?", time.Now(), feedID)

	feeds := h.getFeedsWithCounts(userID)
	entries := FetchEntries(h.DB, userID, "", "", "", "", entryPageSize, 0)
	selCount := selectedCount(h.DB, userID)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	feedsTemplate.Execute(w, map[string]any{"Feeds": feeds, "TotalEntries": totalEntries(feeds), "Locale": locale})
	filterBarTemplate.Execute(w, map[string]any{"Feeds": feeds, "Locale": locale, "SelectedCount": selCount})
	RenderEntriesHTML(w, entries, locale)
}

func (h *FeedsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	locale := detectLocale(r)
	userID := auth.GetUserID(r)
	feedID := r.PathValue("id")

	h.DB.Exec("DELETE FROM feeds WHERE id = ? AND user_id = ?", feedID, userID)

	feeds := h.getFeedsWithCounts(userID)
	selCount := selectedCount(h.DB, userID)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	feedsTemplate.Execute(w, map[string]any{"Feeds": feeds, "TotalEntries": totalEntries(feeds), "Locale": locale})
	filterBarTemplate.Execute(w, map[string]any{"Feeds": feeds, "Locale": locale, "SelectedCount": selCount})
}

func (h *FeedsHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	locale := detectLocale(r)
	userID := auth.GetUserID(r)
	feedID := parseInt64(r.PathValue("id"))

	h.refreshFeed(feedID, userID)

	feeds := h.getFeedsWithCounts(userID)
	entries := FetchEntries(h.DB, userID, "", "", "", "", entryPageSize, 0)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	feedsTemplate.Execute(w, map[string]any{"Feeds": feeds, "TotalEntries": totalEntries(feeds), "Locale": locale})
	RenderEntriesHTML(w, entries, locale)
}

func (h *FeedsHandler) RefreshAll(w http.ResponseWriter, r *http.Request) {
	locale := detectLocale(r)
	userID := auth.GetUserID(r)

	var feeds []models.Feed
	h.DB.Select(&feeds, "SELECT id FROM feeds WHERE user_id = ?", userID)
	for _, f := range feeds {
		h.refreshFeed(f.ID, userID)
	}

	feedsWithUnread := h.getFeedsWithCounts(userID)
	entries := FetchEntries(h.DB, userID, "", "", "", "", entryPageSize, 0)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	feedsTemplate.Execute(w, map[string]any{"Feeds": feedsWithUnread, "TotalEntries": totalEntries(feedsWithUnread), "Locale": locale})
	RenderEntriesHTML(w, entries, locale)
}

func (h *FeedsHandler) getFeedsWithCounts(userID int64) []models.Feed {
	var feeds []models.Feed
	h.DB.Select(&feeds, "SELECT *, (SELECT COUNT(*) FROM entries WHERE feed_id = feeds.id) AS entry_count FROM feeds WHERE user_id = ? ORDER BY title COLLATE NOCASE ASC", userID)
	return feeds
}

func (h *FeedsHandler) refreshFeed(feedID, userID int64) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic refreshing feed %d: %v", feedID, r)
		}
	}()

	var feed models.Feed
	err := h.DB.Get(&feed, "SELECT * FROM feeds WHERE id = ? AND user_id = ?", feedID, userID)
	if err != nil {
		return
	}

	fp := gofeed.NewParser()
	fp.Client = httpClientWithTimeout(30 * time.Second)
	parsed, err := fp.ParseURL(feed.URL)
	if err != nil {
		log.Printf("refresh feed %d: %v", feedID, err)
		return
	}

	h.saveFeedEntries(feedID, parsed)
	h.DB.Exec("UPDATE feeds SET last_fetched_at = ? WHERE id = ?", time.Now(), feedID)
}

func (h *FeedsHandler) saveFeedEntries(feedID int64, parsed *gofeed.Feed) {
	for _, item := range parsed.Items {
		guid := item.GUID
		if guid == "" {
			guid = item.Link
		}
		if guid == "" {
			continue
		}
		var pubAt *time.Time
		if item.PublishedParsed != nil {
			pubAt = item.PublishedParsed
		} else if item.UpdatedParsed != nil {
			pubAt = item.UpdatedParsed
		}
		entryTitle := stripTags(item.Title)
		entrySummary := ""
		if item.Content != "" {
			entrySummary = sanitizeHTML(html.UnescapeString(item.Content))
		} else if item.Description != "" {
			entrySummary = sanitizeHTML(html.UnescapeString(item.Description))
		}
		h.DB.Exec(`
			INSERT OR IGNORE INTO entries (feed_id, guid, title, url, summary, published_at)
			VALUES (?, ?, ?, ?, ?, ?)
		`, feedID, guid, entryTitle, item.Link, entrySummary, db.NullTime(pubAt))
	}
}

func totalEntries(feeds []models.Feed) int {
	t := 0
	for _, f := range feeds {
		t += f.EntryCount
	}
	return t
}
