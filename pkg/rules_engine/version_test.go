package rulesengine

import (
	"regexp"
	"testing"
)

// TestVersionIsSemantic: the engine version names its tag
// rules_engine/vX.Y.Z (8.2).
func TestVersionIsSemantic(t *testing.T) {
	if !regexp.MustCompile(`^\d+\.\d+\.\d+$`).MatchString(Version) {
		t.Fatalf("Version %q is not X.Y.Z", Version)
	}
}
