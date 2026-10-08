package legal

import (
	"net/url"
	"strings"
	"testing"
)

func TestNoticeLinksAreHTTPS(t *testing.T) {
	for _, line := range Notice {
		for _, p := range line.Parts {
			if p.URL == "" {
				continue
			}
			u, err := url.Parse(p.URL)
			if err != nil || u.Scheme != "https" || u.Host == "" {
				t.Errorf("bad link %q", p.URL)
			}
		}
	}
}

func TestPlainTextHasEveryLink(t *testing.T) {
	text := PlainText()
	for _, line := range Notice {
		for _, p := range line.Parts {
			if p.URL != "" && !strings.Contains(text, p.URL) {
				t.Errorf("plain text misses %q", p.URL)
			}
		}
	}
}
