package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestCancer(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Hand: []engine.CardRef{"cancer"}}, seat("cain")},
		Monsters: []engine.CardRef{"leech", "fly"},
	})
	tb.Play(0, "cancer")
	if !hasItem(tb.G, 0, "cancer") {
		t.Fatal("Cancer is a trinket: it stays in play")
	}
	// The first roll gets +1: a 3 hits the leech (DC 4).
	tb.Attack("leech", 3)
	leech, _ := tb.G.Monsters[0].TopOf()
	if tb.G.Object(leech).Card == "leech" {
		t.Error("the first roll did not get +1")
	}
}

func TestCancerOnlyFirstRoll(t *testing.T) {
	plays := 2
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{
			Character: "isaac", Items: []engine.CardRef{"cancer", "the_d6"}, Deactivated: []engine.CardRef{"the_d6"},
			Hand: []engine.CardRef{"pills_g2", "pills_g2"},
		}, seat("cain")},
		LootPlays: &plays,
	})
	tb.G.ForceRolls(2, 2)
	tb.Play(0, "pills_g2") // 2 + 1 = 3: monsters you attack have -1 DC
	if tb.G.Object(tb.Find(0, "the_d6")).Charged {
		t.Fatal("the first roll got no +1")
	}
	tb.Play(0, "pills_g2") // 2: recharge each item
	if !tb.G.Object(tb.Find(0, "the_d6")).Charged {
		t.Error("the second roll got +1 too")
	}
}
