package challenge_test

import (
	"embed"
	"testing"

	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

//go:embed testdata/situations
var situations embed.FS

// TestSituations runs this set's situation files (ADR 034). A file
// starts its challenge with setup.game (24.34); the test sets
// krampus_trial and trick_trial come too, so a file may name their
// cards in its sets.
func TestSituations(t *testing.T) {
	enginetest.RunSituations(t, situations, trickSets()...)
}
