package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestTheRainmaker(t *testing.T) {
	tb := fight(t, "the_rainmaker")
	hurt(tb, 0, 1)
	tb.G.PutOnDeck(engine.MonsterDeck, "curse_of_the_hunted")
	tb.G.ForceRolls(1, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "the_rainmaker", "curse_of_the_hunted", "ask.slot_card 2 fly", foe)
	if len(tb.G.Players[1].InPlay) != 1 {
		t.Error("the curse came into slot 2 and went to Cain")
	}
	if len(tb.G.Players[0].Hand) != 1+3 {
		t.Errorf("the roll of 1 looted 1, then 3 rewards; hand %d", len(tb.G.Players[0].Hand))
	}
}
