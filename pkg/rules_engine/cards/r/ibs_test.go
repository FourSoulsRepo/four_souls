package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestIBS(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "the_soiled", Items: []engine.CardRef{"ibs"}, Hand: []engine.CardRef{"bomb", "a_nickel", "dice_shard"}}, seat("cain"),
	}})
	tb.Activate(0, "ibs", 0, "Gain 3¢.", "bomb")
	if tb.G.Players[0].Cents != 3 {
		t.Error("discard a card: gain 3¢")
	}
}

// A deck attack is an attack on the revealed monster (Q-24, 19.62).
func TestIBSCancelsDeckAttack(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "the_soiled", Items: []engine.CardRef{"ibs"}, Hand: []engine.CardRef{"bomb"}}, seat("cain")},
		Monsters: []engine.CardRef{"leech", "fly"},
	})
	tb.G.PutOnDeck(engine.MonsterDeck, "pooter")
	tb.G.ForceRolls(1)
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentAttack}); err != nil {
		t.Fatal(err)
	}
	tb.Pass(0)
	tb.Pass(1)
	tb.Choose("monster deck")
	tb.Choose("slot 1: leech")
	tb.ThrowPending() // the attack roll waits on the stack
	tb.Activate(0, "ibs", 0, "Cancel your attack on a monster.", "bomb")
	if tb.G.Attack.On {
		t.Error("IBS should cancel the attack on the revealed pooter")
	}
	if tb.G.Players[0].Damage != 0 {
		t.Error("the cancelled roll deals no damage")
	}
}
