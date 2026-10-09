package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestPsyHorf(t *testing.T) {
	tb := slayTable(t, "psy_horf", engine.SituationPlayer{Character: "isaac", Items: items("jawbone"), Deactivated: items("jawbone")}, seat("cain"))
	killWithSixes(t, tb)
	if !tb.G.Object(tb.Find(0, "jawbone")).Charged {
		t.Error("the jawbone is not recharged")
	}
}
