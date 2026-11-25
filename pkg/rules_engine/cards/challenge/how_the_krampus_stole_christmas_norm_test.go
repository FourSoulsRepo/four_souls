package challenge_test

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/challenge"
)

// The forgo itself, yes and no, is in
// testdata/situations/how_the_krampus_stole_christmas_norm (24.27).
// What stays here: the question's text and a save in the middle of it.

// R-KRAMPUS-07, Q-32: the forgo question asks the player who would gain
// the ¢; a save in the middle of the question goes on after a load.
func TestForgoCentsSaveMidQuestion(t *testing.T) {
	tb := krampus(t, engine.DifficultyNormal, engine.SituationPlayer{}, engine.SituationPlayer{})
	cents := tb.G.Players[0].Cents
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentAttack}); err != nil {
		t.Fatal(err)
	}
	tb.G.PutOnDeck(engine.MonsterDeck, "fly")
	tb.G.ForceRolls(6)
	if _, err := tb.G.Settle("fly"); err == nil {
		t.Fatal("no forgo question")
	}
	if w := tb.G.Prompt(); w.Player != 0 || w.Text.Key.String() != "card.how_the_krampus_stole_christmas.forgo_cents" || !slices.Equal(tb.Options(), []string{"ask.yes", "ask.no"}) {
		t.Fatalf("question %+v, options %v", w, tb.Options())
	}
	data, err := tb.G.Save()
	if err != nil {
		t.Fatal(err)
	}
	if tb.G, err = engine.Load(data, sets()...); err != nil {
		t.Fatal(err)
	}
	ev := tb.Settle("ask.yes")
	if tb.G.Players[0].Cents != cents || tb.G.PresentCounters(engine.PresentCents) != 9 || count(ev, engine.EvForgone) != 1 {
		t.Errorf("forgone: cents %d, counters %v", tb.G.Players[0].Cents, counters(tb.G))
	}
	check(t, tb.G)
}

// R-KRAMPUS-08 (Normal): the minions come one at a time from the top of
// the monster deck, so the top card takes the first minion slot. The
// view shows the destroyed present face up. The present, its discard
// and the minions are in
// how_the_krampus_stole_christmas_norm/a-present-at-0-is-destroyed-3-minions-come.json.
func TestWhatsInTheBoxOrderAndView(t *testing.T) {
	tb := krampus(t, engine.DifficultyNormal, engine.SituationPlayer{}, engine.SituationPlayer{})
	setCounters(tb, engine.PresentCents, 2)
	card := tb.G.Object(tb.G.Challenge.Presents[0].Object).Card
	tb.G.PutOnDeck(engine.MonsterDeck, "pooter", "dip", "squirt")
	tb.G.ForceRolls(1)
	tb.EndTurn()
	if got := minionCards(tb.G); !slices.Equal(got, []engine.CardRef{"squirt", "dip", "pooter"}) {
		t.Errorf("minions %v", got)
	}
	v := tb.G.View(engine.Viewer{Kind: engine.ViewSpectator}).Challenge.Presents[0]
	if !v.Destroyed || v.Card != card || v.Counters != 0 {
		t.Errorf("view %+v", v)
	}
	check(t, tb.G)
}

// R-KRAMPUS-01, R-KRAMPUS-09: the fourth destroyed present loses the
// challenge with the presents as its reason. The loss itself is in
// how_the_krampus_stole_christmas_norm/the-fourth-destroyed-present-loses-the-challenge.json.
func TestAllPresentsLostReason(t *testing.T) {
	tb := krampus(t, engine.DifficultyNormal, engine.SituationPlayer{}, engine.SituationPlayer{})
	for _, kind := range []string{engine.PresentCents, engine.PresentLoot, engine.PresentTreasure, engine.PresentSoul} {
		tb.G.PutOnDeck(engine.MonsterDeck, "pooter", "dip", "squirt")
		tb.G.AddPresentCounters(kind, -20, engine.NoPlayer)
		tb.G.Run()
		tb.Settle()
	}
	if c := tb.G.Challenge; !tb.G.Over || !c.Lost || c.Reason.Key != challenge.ReasonPresents {
		t.Errorf("over %v, lost %v (%v)", tb.G.Over, c.Lost, c.Reason)
	}
	check(t, tb.G)
}

// R-KRAMPUS-10, Q-40 (Normal): after Happy Holidays! the view shows each
// present face up with the player who gained it. Who gains each, and
// that it happens once, is in
// how_the_krampus_stole_christmas_norm/happy-holidays-at-5-souls-gives-each-present.json.
func TestHappyHolidaysView(t *testing.T) {
	tb := krampus(t, engine.DifficultyNormal, engine.SituationPlayer{Souls: souls(3)}, engine.SituationPlayer{Souls: souls(1)})
	tb.G.DestroyPresent(engine.PresentLoot)
	isaac, cain := player(tb, 0), player(tb, 1)
	boost(tb, 0, 5)
	tb.G.ForceRolls(6)
	tb.RevealFromDeck("mixed_bag", "ask.no", "ask.no", cain, isaac, cain)
	for i, want := range []engine.PlayerID{1, engine.NoPlayer, 0, 1} {
		if want == engine.NoPlayer {
			continue
		}
		o := tb.G.Object(tb.G.Challenge.Presents[i].Object)
		v := tb.G.View(engine.Viewer{Kind: engine.ViewSpectator}).Challenge.Presents[i]
		if !v.Flipped || v.Player != want || v.Card != o.Card {
			t.Errorf("present %d view %+v", i, v)
		}
	}
	check(t, tb.G)
}
