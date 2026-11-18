package rulesengine_test

import (
	"embed"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

//go:embed testdata/situations
var situations embed.FS

// TestSituations runs the engine's rule tree, test sets only (LR-07,
// ADR 034): a folder per rule topic, each file must match its expect.
func TestSituations(t *testing.T) {
	enginetest.RunRuleSituations(t, situations, engine.TestSet, engine.AbilitySet)
}
