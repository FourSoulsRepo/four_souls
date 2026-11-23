package cards

import (
	"embed"
	"testing"

	"github.com/FourSoulsRepo/rules_engine/cards/challenge"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

//go:embed testdata/situations
var situations embed.FS

// TestSituations runs the cross-set situation files: cards of several
// sets together (ADR 034 item 5.3).
func TestSituations(t *testing.T) {
	enginetest.RunSituations(t, situations, append(Sets(), challenge.Set)...)
}
