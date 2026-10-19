package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestEmergencyContact(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), {Character: "cain", Hand: []engine.CardRef{"emergency_contact"}}},
		Monsters: []engine.CardRef{"leech", "fly"},
	})
	tb.G.ForceRolls(1)
	tb.ThrowPending()
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentAttack}); err != nil {
		t.Fatal(err)
	}
	tb.Pass(0)
	tb.Pass(1)
	tb.Choose("leech")
	tb.ThrowPending() // the attack roll waits on the stack
	tb.Pass(0)
	tb.G.Players[1].ExtraLootPlays = 1 // as from his ↷
	tb.Play(1, "emergency_contact")
	if tb.G.Attack.On {
		t.Error("the attack should be cancelled")
	}
	if top(tb, 0) == "leech" {
		t.Error("the leech should leave its slot")
	}
	bottom := tb.G.Decks[engine.MonsterDeck][0]
	if tb.G.Object(bottom).Card != "leech" {
		t.Error("the leech goes to the bottom of the monster deck")
	}
	if tb.G.Players[0].Damage != 0 {
		t.Error("the cancelled roll deals no damage")
	}
}

func TestEmergencyContactNotADeckAttack(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), {Character: "cain", Hand: []engine.CardRef{"emergency_contact"}}},
		Monsters: []engine.CardRef{"leech", "fly"},
	})
	tb.G.PutOnDeck(engine.MonsterDeck, "pooter")
	tb.G.ForceRolls(1, 6, 6) // a miss, then two hits kill the pooter
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentAttack}); err != nil {
		t.Fatal(err)
	}
	tb.Pass(0)
	tb.Pass(1)
	tb.Choose("monster deck")
	tb.Choose("slot 1: leech")
	tb.ThrowPending()
	tb.Pass(0)
	tb.G.Players[1].ExtraLootPlays = 1
	tb.Play(1, "emergency_contact")
	if !inDiscard(tb.G, engine.MonsterDeck, "pooter") {
		t.Error("an attack on the deck is not an attack on a monster (R-ATK-05): it goes on and kills the pooter")
	}
}
