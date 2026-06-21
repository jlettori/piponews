package handlers

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/jlettori/piponews/internal/db"
	"github.com/jlettori/piponews/internal/i18n"
	"github.com/jlettori/piponews/internal/templates"
	texttemplate "text/template"
)

func parseSignals(r *http.Request, target any) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if len(body) == 0 {
		return nil
	}
	return json.Unmarshal(body, target)
}

func parseInt64(s string) int64 {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return n
}

func safeHTML(s string) template.HTML {
	return template.HTML(s)
}

var (
	tagRegex         = regexp.MustCompile(`<[^>]*>`)
	dangerousTagsRE  = regexp.MustCompile(`(?i)<(script|iframe|object|embed|base|form|input|button|select|textarea|style|link|meta)[^>]*>`)
	dangerousAttrsRE = regexp.MustCompile(`(?i)\s+on\w+\s*=\s*(?:"[^"]*"|'[^']*'|[^\s>]*)`)
	javascriptURLsRE = regexp.MustCompile(`(?i)\s+(href|src|action|formaction)\s*=\s*["']?\s*javascript:`)
)

func sanitizeHTML(s string) string {
	s = dangerousTagsRE.ReplaceAllString(s, "")
	s = dangerousAttrsRE.ReplaceAllString(s, "")
	s = javascriptURLsRE.ReplaceAllString(s, "")
	return s
}

func validateFeedURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("only http and https URLs are allowed")
	}
	if u.Host == "" {
		return fmt.Errorf("URL must have a host")
	}
	return nil
}

func httpClientWithTimeout(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyFromEnvironment,
		},
	}
}

var stripTags = func(s string) string {
	return tagRegex.ReplaceAllString(html.UnescapeString(s), "")
}

var BuildVersion = "dev"

func version() string {
	return BuildVersion
}

func initial(s string) string {
	if s == "" {
		return "?"
	}
	r, _ := utf8.DecodeRuneInString(s)
	return string(r)
}

func dict(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("dict requires an even number of arguments, got %d", len(values))
	}
	m := make(map[string]any, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict keys must be strings, got %T", values[i])
		}
		m[key] = values[i+1]
	}
	return m, nil
}

func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"dict":       dict,
		"formatTime": formatTime,
		"initial":    initial,
		"safeHTML":   safeHTML,
		"T":          i18n.T,
		"N":          i18n.N,
		"version":    version,
	}
}

func textTemplateFuncs() texttemplate.FuncMap {
	return texttemplate.FuncMap{
		"formatTime": formatTime,
		"stripTags":  stripTags,
		"T":          i18n.T,
		"N":          i18n.N,
	}
}

func parseTemplates(files ...string) *template.Template {
	allFiles := []string{"base.html", "entry_card.html"}
	allFiles = append(allFiles, files...)
	return template.Must(template.New("base").Funcs(templateFuncs()).ParseFS(templates.FS, allFiles...))
}

func joinInts(ints []int) string {
	parts := make([]string, len(ints))
	for i, v := range ints {
		parts[i] = strconv.Itoa(v)
	}
	return strings.Join(parts, ", ")
}

// ErrResponse writes an HTML error fragment for DataStar
func ErrResponse(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, `<div id="entries-list"><div class="empty-state"><p>%s</p></div></div>`, html.EscapeString(msg))
}

func detectLocale(r *http.Request) i18n.Locale {
	return i18n.DetectLocale(r)
}

func selectedCount(database *db.DB, userID int64) int {
	var count int
	database.Get(&count, `
		SELECT COUNT(*)
		FROM entry_selections us
		JOIN entries e ON e.id = us.entry_id
		JOIN feeds f ON f.id = e.feed_id
		WHERE f.user_id = ? AND us.user_id = ?
	`, userID, userID)
	return count
}
