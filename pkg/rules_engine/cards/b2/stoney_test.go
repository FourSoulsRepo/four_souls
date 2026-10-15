package b2

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestStoney(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac"), seat("cain")}, Monsters: items("fly", "stoney")}, Set)
	stoney, _ := tb.G.Monsters[1].TopOf()
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentAttack}); err != nil {
		t.Fatal(err)
	}
	tb.Pass(0)
	tb.Pass(1)
	if opts := tb.Options(); slices.Contains(opts, "stoney") {
		t.Errorf("Stoney can be attacked: %v", opts)
	}
	tb.G.ForceRolls(6, 6)
	tb.Settle("fly")
	if tb.G.Object(stoney).Zone.Kind == engine.ZoneInPlay {
		t.Error("Stoney did not die with the fly")
	}
}

// A covered Stoney is not in play: the monster on top of it dying does
// not kill it, and it comes back into play.
func TestStoneyCovered(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac"), seat("cain")}, Monsters: items("stoney", "fly")}, Set)
	tb.G.PutOnDeck(engine.MonsterDeck, "pooter")
	tb.G.CoverMonsterSlot(0)
	pooter, _ := tb.G.Monsters[0].TopOf()
	tb.G.KillObject(pooter)
	tb.Settle()
	stoney, ok := tb.G.Monsters[0].TopOf()
	if !ok || tb.G.Object(stoney).Card != "stoney" || tb.G.Object(stoney).Zone.Kind != engine.ZoneInPlay {
		t.Error("Stoney died with the monster that covered it")
	}
}
