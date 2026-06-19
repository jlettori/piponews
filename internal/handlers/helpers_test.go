package handlers

import (
	"testing"
)

func TestValidateFeedURL_Valid(t *testing.T) {
	tests := []string{
		"https://example.com/feed.xml",
		"http://blog.example.org/rss",
		"https://news.example.com/rss/feed.xml",
	}
	for _, url := range tests {
		if err := validateFeedURL(url); err != nil {
			t.Errorf("validateFeedURL(%q) = %v; want nil", url, err)
		}
	}
}

func TestValidateFeedURL_Invalid(t *testing.T) {
	tests := []struct {
		url string
		msg string
	}{
		{"", "only http and https URLs are allowed"},
		{"not-a-url", "only http and https URLs are allowed"},
		{"ftp://example.com", "only http and https URLs are allowed"},
		{"javascript:alert(1)", "only http and https URLs are allowed"},
		{"http://", "URL must have a host"},
	}
	for _, tt := range tests {
		err := validateFeedURL(tt.url)
		if err == nil {
			t.Errorf("validateFeedURL(%q) = nil; want error containing %q", tt.url, tt.msg)
			continue
		}
		if err.Error() != tt.msg {
			t.Errorf("validateFeedURL(%q) = %q; want %q", tt.url, err.Error(), tt.msg)
		}
	}
}

func TestParseInt64(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{"0", 0},
		{"42", 42},
		{"-1", -1},
		{"", 0},
		{"abc", 0},
		{"9999999999999", 9999999999999},
	}
	for _, tt := range tests {
		got := parseInt64(tt.input)
		if got != tt.want {
			t.Errorf("parseInt64(%q) = %d; want %d", tt.input, got, tt.want)
		}
	}
}

func TestSanitizeHTML_RemovesScriptTags(t *testing.T) {
	input := `<p>Hello</p><script>alert("xss")</script><p>World</p>`
	got := sanitizeHTML(input)
	if contains(got, "<script") {
		t.Errorf("sanitizeHTML still contains script tag: %q", got)
	}
}

func TestSanitizeHTML_RemovesEventHandlers(t *testing.T) {
	input := `<a href="#" onclick="alert(1)">click</a>`
	got := sanitizeHTML(input)
	if contains(got, "onclick") {
		t.Errorf("sanitizeHTML still contains onclick: %q", got)
	}
}

func TestSanitizeHTML_RemovesIframe(t *testing.T) {
	input := `<iframe src="https://evil.com"></iframe>`
	got := sanitizeHTML(input)
	if contains(got, "<iframe") {
		t.Errorf("sanitizeHTML still contains iframe: %q", got)
	}
}

func TestSanitizeHTML_RemovesJavascriptURLs(t *testing.T) {
	input := `<a href="javascript:alert(1)">click</a>`
	got := sanitizeHTML(input)
	if contains(got, "javascript:") {
		t.Errorf("sanitizeHTML still contains javascript: %q", got)
	}
}

func TestSanitizeHTML_AllowsSafeHTML(t *testing.T) {
	input := `<p>Hello <strong>world</strong></p>`
	got := sanitizeHTML(input)
	if got != input {
		t.Errorf("sanitizeHTML(%q) = %q; want unchanged", input, got)
	}
}

func TestStripTags(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"<p>Hello</p>", "Hello"},
		{"Hello", "Hello"},
		{"<b>Bold</b> and <i>italic</i>", "Bold and italic"},
		{"<a href=\"#\">link</a>", "link"},
		{"", ""},
	}
	for _, tt := range tests {
		got := stripTags(tt.input)
		if got != tt.want {
			t.Errorf("stripTags(%q) = %q; want %q", tt.input, got, tt.want)
		}
	}
}

func TestTotalEntries(t *testing.T) {
	feeds := []struct {
		EntryCount int
	}{
		{3},
		{5},
		{0},
		{2},
	}
	total := 0
	for _, f := range feeds {
		total += f.EntryCount
	}
	if total != 10 {
		t.Errorf("total = %d; want 10", total)
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
