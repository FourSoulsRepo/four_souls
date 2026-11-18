package b2

// Interaction tests: card combinations and rulings from step 3 (TS-04).
// Each test names its source: the re-checked old rulings (R-07) and the
// S-RU FAQ rulings in docs/rules/open-questions.md, and Yuggy's rulings.
// The simple ones are situation files (ADR 033): testdata/situations/.

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// R-07 #4, narrowed by Q-07: a dead player takes untargeted damage, but
// nothing is marked (R-MECH-16). The rest of this case is the file
// interactions/xiii-death-xvi-the-tower-damage-after-death.json.
func TestRulingDamageToDeadPlayer(t *testing.T) {
	tb := lootTable(t, "xiii_death", "xvi_the_tower")
	tb.G.Players[0].ExtraLootPlays = 1
	tb.Play(0, "xiii_death", foe)
	before := tb.G.Players[1].Damage
	tb.G.ForceRolls(1) // each player takes 1 damage
	tb.Play(0, "xvi_the_tower")
	if d := tb.G.Players[1].Damage; d != before {
		t.Errorf("damage %d, want %d", d, before)
	}
}

// S-RU FAQ #1, Q-09: Two of Clubs used again still only doubles; the
// wording matters ("loot double that number"), so uses don't stack.
func TestRulingTwoOfClubsTwice(t *testing.T) {
	for uses := 1; uses <= 3; uses++ {
		tb := itemTable(t, items("two_of_clubs", "bum_friend"), "a_dime")
		for i := range uses {
			if i > 0 {
				tb.G.Recharge(tb.Find(0, "two_of_clubs"))
			}
			tb.Activate(0, "two_of_clubs", 0, me)
		}
		tb.Activate(0, "bum_friend", 0, "a_dime") // loot 1, doubled once; one goes back
		if h := len(tb.G.Players[0].Hand); h != 2 {
			t.Errorf("%d uses: hand = %d, want 2", uses, h)
		}
	}
}

// Yuggy (X, 2024-08-09): a roll pushed above 6 counts as 6; minus 1 is 5.
func TestRulingRollAboveSix(t *testing.T) {
	tb := rollTable(t, 6, seat("judas", "spoon_bender", "book_of_belial"), seat("cain"))
	tb.Pass(1)
	tb.Start(0, "spoon_bender", 0)
	tb.Choose("ask.roll_of 6")
	tb.Pass(0)
	tb.Pass(1) // the +1 resolves: still 6
	tb.Activate(0, "book_of_belial", 0, "Subtract 1 from a roll.", "ask.roll_of 6")
	if r, ok := firstRoll(tb.G); ok && r != 5 {
		t.Errorf("roll = %d, want 5", r)
	}
}

func firstRoll(g *engine.Game) (int, bool) {
	for _, it := range g.Stack {
		if it.Kind == engine.StackRoll {
			return it.Roll, true
		}
	}
	return 0, false
}
