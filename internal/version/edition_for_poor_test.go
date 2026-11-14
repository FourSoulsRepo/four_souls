//go:build !(embed && cardimages)

package version

import "testing"

func TestEditionFromTags(t *testing.T) {
	if Edition != EditionForPoor {
		t.Errorf("build without embed cardimages: Edition = %q, want %q", Edition, EditionForPoor)
	}
}
