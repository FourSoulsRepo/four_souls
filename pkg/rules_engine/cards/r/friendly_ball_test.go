package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestFriendlyBall(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", "friendly_ball"), seat("cain")},
		Monsters: []engine.CardRef{"leech", "fly"},
	})
	if rule := tb.Refused(0, "friendly_ball", 1); rule != "R-ABIL-07" {
		t.Errorf("nothing under it: the second ability is refused by %s", rule)
	}
	tb.Activate(0, "friendly_ball", 0, "leech")
	ball := tb.G.Object(tb.Find(0, "friendly_ball"))
	if len(ball.Under) != 1 || top(tb, 0) == "leech" {
		t.Fatal("the leech goes under the ball and its slot refills")
	}
	ball.Charged = true
	if rule := tb.Refused(0, "friendly_ball", 0); rule != "R-ABIL-07" {
		t.Errorf("a card under it: the first ability is refused by %s", rule)
	}
	tb.Activate(0, "friendly_ball", 1, "ask.slot 2")
	if top(tb, 1) != "leech" {
		t.Fatalf("slot 2 = %s, want the leech back", top(tb, 1))
	}
	// "The active player must make an additional attack on it this
	// turn": one more attack, owed on the leech only (19.241).
	leech, _ := tb.G.Monsters[1].TopOf()
	if one := tb.G.Turn.MustAttackOneOf; len(one) != 1 || one[0] != leech || tb.G.Turn.MustAttacks != 0 || tb.G.Turn.Attacks != 2 {
		t.Errorf("must attack one of %v, must attacks %d, attacks %d; want one more attack, on the leech", one, tb.G.Turn.MustAttacks, tb.G.Turn.Attacks)
	}
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentEndTurn}); err == nil {
		t.Error("the turn can't end before the forced attack")
	}
}
