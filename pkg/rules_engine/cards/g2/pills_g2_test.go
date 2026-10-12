package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestPillsG2Recharges(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"the_d6"}, Deactivated: []engine.CardRef{"the_d6"}, Hand: []engine.CardRef{"pills_g2"}},
		seat("cain"),
	}})
	tb.G.ForceRolls(2)
	tb.Play(0, "pills_g2")
	if !tb.G.Object(tb.Find(0, "the_d6")).Charged {
		t.Error("1-2: the D6 was not recharged")
	}
}

func TestPillsG2ChangesDC(t *testing.T) {
	for _, c := range []struct{ roll, want int }{{3, 3}, {5, 6}} {
		tb := lootTable(t, "pills_g2")
		tb.G.ForceRolls(c.roll)
		tb.Play(0, "pills_g2")
		declareAttack(tb, "leech")
		leech, _ := tb.G.Monsters[1].TopOf()
		if dc := tb.G.Evasion(leech); dc != c.want {
			t.Errorf("roll %d: the attacked leech has DC %d, want %d", c.roll, dc, c.want)
		}
		fly, _ := tb.G.Monsters[0].TopOf()
		if dc := tb.G.Evasion(fly); dc != 2 {
			t.Errorf("roll %d: the fly is not attacked; DC %d, want 2", c.roll, dc)
		}
	}
}
