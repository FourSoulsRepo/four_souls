package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestCurvedHorn(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", "curved_horn"), seat("cain")},
		Monsters: []engine.CardRef{"conjoined_fatty"}, // 4 HP
	}, Set)
	tb.G.Turn.Attacks = 1
	fatty, _ := tb.G.Monsters[0].TopOf()
	tb.G.ForceRolls(6, 6, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "conjoined_fatty")
	// 2 damage on the first roll, then 1 and 1: dead after the third roll.
	if tb.G.Object(fatty).Zone.Kind == engine.ZoneInPlay {
		t.Error("the fatty survived 2 + 1 + 1 damage")
	}
	if tb.G.Turn.AttackRolls != 3 {
		t.Errorf("%d attack rolls, want 3", tb.G.Turn.AttackRolls)
	}
}

// The view shows Curved Horn's +1 ATK from before the attack until the
// first attack roll's damage is dealt (19.67).
func TestCurvedHornShownForTheFirstRoll(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", "curved_horn"), seat("cain")},
		Monsters: items("big_spider"), // 3 HP, DC 4
	}, Set)
	check := func(when string, atk int, horn bool) {
		t.Helper()
		pv := tb.G.View(engine.Viewer{Kind: engine.ViewSpectator}).Players[0]
		shown := false
		for _, s := range pv.Sources {
			shown = shown || (s.Card == "curved_horn" && s.Stat == "atk" && s.When.String() == "ask.when.first_attack_roll")
		}
		if pv.ATK != atk || shown != horn {
			t.Errorf("%s: ATK %d, horn shown %v; want %d, %v", when, pv.ATK, shown, atk, horn)
		}
	}
	check("before the attack", 2, true)
	tb.G.ForceRolls(6, 6)
	attackRoll(tb, "big_spider")
	check("first roll on the stack", 2, true)
	tb.Pass(0)
	tb.Pass(1)
	check("its damage on the stack", 2, true)
	tb.Pass(0)
	tb.Pass(1)
	check("its damage dealt", 1, false)
}

// A missed first attack roll takes Curved Horn's +1 ATK off the view
// (19.67).
func TestCurvedHornGoneAfterAMiss(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", "curved_horn"), seat("cain")},
		Monsters: items("big_spider"),
	}, Set)
	tb.G.ForceRolls(1)
	attackRoll(tb, "big_spider")
	tb.Pass(0)
	tb.Pass(1) // a miss: the spider's damage on the stack
	if atk := tb.G.View(engine.Viewer{Kind: engine.ViewSpectator}).Players[0].ATK; atk != 1 {
		t.Errorf("after a miss: ATK %d, want 1", atk)
	}
}
