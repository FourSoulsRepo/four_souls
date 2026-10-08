package legal

import (
	"os"
	"strings"
	"testing"
)

// The screenshot fixture copies the notice; keep both in sync.
func TestScreenshotFixtureMatchesNotice(t *testing.T) {
	data, err := os.ReadFile("../../frontend/src/screenshotBridge.ts")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	fixture := string(data)
	for _, line := range Notice {
		for _, p := range line.Parts {
			if !strings.Contains(fixture, "'"+p.Text+"'") {
				t.Errorf("fixture misses text %q", p.Text)
			}
			if p.URL != "" && !strings.Contains(fixture, p.URL) {
				t.Errorf("fixture misses url %q", p.URL)
			}
		}
	}
}
