package r

import (
	"embed"
	"testing"

	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

//go:embed testdata/situations
var situations embed.FS

// TestSituations runs this set's situation files (ADR 033).
func TestSituations(t *testing.T) {
	enginetest.RunSituations(t, situations, sets...)
}
