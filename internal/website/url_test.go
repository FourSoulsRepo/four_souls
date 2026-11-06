package website

import (
	"net/url"
	"os"
	"strings"
	"testing"
)

// The site's address must match the path Vite builds it for, or every
// asset link on the deployed page breaks.
func TestURLMatchesViteBase(t *testing.T) {
	u, err := url.Parse(URL)
	if err != nil {
		t.Fatal(err)
	}
	if u.Scheme != "https" || u.Host == "" || !strings.HasSuffix(u.Path, "/") {
		t.Fatalf("URL = %q: want https, a host and a path ending in /", URL)
	}
	cfg, err := os.ReadFile("../../frontend/vite.config.ts")
	if err != nil {
		t.Fatal(err)
	}
	if want := "base: '" + u.Path + "'"; !strings.Contains(string(cfg), want) {
		t.Errorf("frontend/vite.config.ts has no %s", want)
	}
}
