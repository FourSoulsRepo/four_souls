package about

import (
	"testing"

	"github.com/FourSoulsRepo/four_souls/internal/legal"
	"github.com/FourSoulsRepo/four_souls/internal/version"
)

// The binding hands over the build's versions and the whole notice.
func TestAbout(t *testing.T) {
	a := New()
	if got := a.Versions(); got != version.Get() {
		t.Errorf("versions %+v, want %+v", got, version.Get())
	}
	if got := a.Notice(); len(got) == 0 || len(got) != len(legal.Notice) {
		t.Errorf("notice has %d lines, want %d", len(got), len(legal.Notice))
	}
}
