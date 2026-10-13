package g2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTheBibleKillsMom(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"the_bible"}}, seat("cain")},
		Monsters: []engine.CardRef{"mom", "fatty"},
	})
	tb.Activate(0, "the_bible", 0, "Kill a monster named Mom!, Mom's Heart!, or It Lives!", "mom")
	if top, _ := tb.G.Monsters[0].TopOf(); tb.G.Object(top).Card == "mom" {
		t.Error("Mom is still alive")
	}
}

func TestTheBibleCancelsYourAttack(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"the_bible"}}, seat("cain")},
		Monsters: []engine.CardRef{"leech", "fatty"},
	})
	tb.G.ForceRolls(1)
	declareAttack(tb, "leech")
	tb.Activate(0, "the_bible", 0, "Cancel your attack.")
	if tb.G.Attack.On {
		t.Error("the attack was not cancelled")
	}
	if d := tb.G.Players[0].Damage; d != 0 {
		t.Errorf("damage = %d: the missed roll must not resolve", d)
	}
}

func TestTheBibleDestroysACurse(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"the_bible"}}, seat("cain")},
		Monsters: []engine.CardRef{"fatty", "monstro"},
	})
	tb.RevealFromDeck("curse_of_pain", me)
	tb.Activate(0, "the_bible", 0, "Destroy a curse.", "curse_of_pain")
	for _, id := range tb.G.Players[0].InPlay {
		if tb.G.Object(id).Card == "curse_of_pain" {
			t.Error("the curse is still there")
		}
	}
}
