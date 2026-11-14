//go:build embed && cardimages

package version

import "testing"

func TestEditionFromTags(t *testing.T) {
	if Edition != EditionFull {
		t.Errorf("embed cardimages build: Edition = %q, want %q", Edition, EditionFull)
	}
}
