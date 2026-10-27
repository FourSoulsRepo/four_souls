package g2

import "testing"

func TestFruitcakeRolls(t *testing.T) {
	tb := table(t, seat("isaac", "fruitcake"), seat("cain"))
	tb.G.ForceRolls(2)
	tb.Activate(0, "fruitcake", 0)
	if c := tb.G.Players[0].Cents; c != 1 {
		t.Errorf("2: cents = %d, want 1", c)
	}
}

func TestFruitcakeRechargesOnTwiceInARow(t *testing.T) {
	tb := table(t, seat("isaac", "fruitcake", "the_butter_bean"), seat("cain"))
	tb.G.ForceRolls(2, 2)
	tb.Activate(0, "fruitcake", 0)       // 2: gain 1¢; Fruitcake is deactivated
	tb.Activate(0, "the_butter_bean", 0) // 2 again: Fruitcake recharges
	if !tb.G.Object(tb.Find(0, "fruitcake")).Charged {
		t.Error("two 2s in a row should recharge Fruitcake")
	}
}

func TestFruitcakeOtherPlayersRollsBetween(t *testing.T) {
	// Q-10 (a): each player's own rolls count.
	tb := table(t, seat("isaac", "fruitcake", "the_butter_bean"), seat("cain", "crooked_penny"))
	tb.G.ForceRolls(4, 1, 4)
	tb.Activate(0, "fruitcake", 0, "ask.deck.loot") // Isaac 4
	tb.Start(0, "isaac", 0)
	tb.Pass(0)
	tb.Activate(1, "crooked_penny", 0, me) // Cain 1
	tb.Activate(0, "the_butter_bean", 0)   // Isaac 4 again
	if !tb.G.Object(tb.Find(0, "fruitcake")).Charged {
		t.Error("Isaac rolled 4 twice in a row; Cain's roll between does not count")
	}
}
