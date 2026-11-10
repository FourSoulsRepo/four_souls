package challenge_test

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/challenge"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// kill attacks a monster in a slot with a hit of 6 and Isaac's ATK
// boosted, then answers the prompts. A fly refills the slot.
func kill(tb *enginetest.Table, monster engine.CardRef, choices ...string) []engine.Event {
	tb.T.Helper()
	boost(tb, 0, 5)
	tb.G.PutOnDeck(engine.MonsterDeck, "fly")
	tb.G.ForceRolls(6)
	return tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, append([]string{string(monster)}, choices...)...)
}

// R-KRAMPUS-07, Q-32: the player who would gain a monster's ¢ may forgo
// them for as many counters on the ¢ present; a no gains them.
func TestForgoCents(t *testing.T) {
	tb := krampus(t, engine.DifficultyNormal, engine.SituationPlayer{}, engine.SituationPlayer{})
	cents := tb.G.Players[0].Cents
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentAttack}); err != nil {
		t.Fatal(err)
	}
	boost(tb, 0, 5)
	tb.G.PutOnDeck(engine.MonsterDeck, "fly")
	tb.G.ForceRolls(6)
	if _, err := tb.G.Settle("fly"); err == nil {
		t.Fatal("no forgo question")
	}
	if w := tb.G.Prompt(); w.Player != 0 || w.Text.Key.String() != "card.how_the_krampus_stole_christmas.forgo_cents" || !slices.Equal(tb.Options(), []string{"ask.yes", "ask.no"}) {
		t.Fatalf("question %+v, options %v", w, tb.Options())
	}
	// A save in the middle of the question goes on after a load.
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

	tb = krampus(t, engine.DifficultyNormal, engine.SituationPlayer{}, engine.SituationPlayer{})
	kill(tb, "fly", "ask.no")
	if tb.G.Players[0].Cents != cents+1 || tb.G.PresentCounters(engine.PresentCents) != 8 {
		t.Errorf("kept: cents %d, counters %v", tb.G.Players[0].Cents, counters(tb.G))
	}
	check(t, tb.G)
}

// Q-32: each part of a mixed reward is offered on its own, then the
// soul; a forgone soul card goes to discard. Q-33: a part whose present
// is destroyed is not offered; it is gained.
func TestForgoMixed(t *testing.T) {
	tb := krampus(t, engine.DifficultyNormal, engine.SituationPlayer{}, engine.SituationPlayer{})
	tb.G.DestroyPresent(engine.PresentLoot)
	cents, hand := tb.G.Players[0].Cents, len(tb.G.Players[0].Hand)
	tb.G.ForceRolls(6)
	boost(tb, 0, 5)
	tb.RevealFromDeck("mixed_bag", "ask.yes", "ask.yes")
	if got := counters(tb.G); got != [4]int{10, 0, 3, 4} {
		t.Errorf("counters %v: 2¢ and the soul forgone", got)
	}
	if tb.G.Players[0].Cents != cents || len(tb.G.Players[0].Hand) != hand+1 || tb.G.SoulValue(0) != 0 {
		t.Errorf("cents %d, hand %d, souls %d: the loot is gained", tb.G.Players[0].Cents, len(tb.G.Players[0].Hand), tb.G.SoulValue(0))
	}
	d := tb.G.Discards[engine.MonsterDeck]
	if len(d) == 0 || tb.G.Object(d[len(d)-1]).Card != "mixed_bag" {
		t.Error("the forgone soul card goes to discard")
	}
	check(t, tb.G)
}

// Q-32: a rolled reward is offered after its roll, with the result.
func TestForgoRolledReward(t *testing.T) {
	tb := krampus(t, engine.DifficultyNormal, engine.SituationPlayer{}, engine.SituationPlayer{})
	boost(tb, 0, 5)
	tb.G.ForceRolls(6, 4) // the hit, then the reward roll
	tb.RevealFromDeck("cursed_keeper_head", "ask.yes")
	if got := tb.G.PresentCounters(engine.PresentCents); got != 12 {
		t.Errorf("¢ present %d, want 8 + the roll of 4", got)
	}
}

// R-KRAMPUS-07, Q-32: a soul forgone gives its soul value (Mom!, 2);
// the treasure reward is offered too.
func TestForgoSoulValue(t *testing.T) {
	tb := krampus(t, engine.DifficultyNormal, engine.SituationPlayer{Items: []engine.CardRef{"chaos_card"}}, engine.SituationPlayer{})
	tb.G.PutOnDeck(engine.MonsterDeck, "mom")
	tb.G.MillTop(engine.MonsterDeck)
	tb.G.PlaceFromDiscard(tb.G.Discards[engine.MonsterDeck][len(tb.G.Discards[engine.MonsterDeck])-1], 0)
	tb.G.Run()
	tb.Settle()
	tb.G.PutOnDeck(engine.MonsterDeck, "fly") // Mom's new slot
	tb.Activate(0, "chaos_card", 0, "Kill a player or monster.", "mom", "ask.no", "ask.yes")
	if got := counters(tb.G); got != [4]int{8, 2, 3, 5} {
		t.Errorf("counters %v: Mom's soul value 2 on the soul present", got)
	}
	check(t, tb.G)
}

// Q-32: doubled rewards (Blessing of the Sack) are forgone doubled.
// Q-32: other gains (A Penny!) are no monster rewards: no question.
func TestForgoDoubledAndOther(t *testing.T) {
	tb := krampusOn(t, engine.DifficultyNormal, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Hand: []engine.CardRef{"a_penny"}}, {}},
		Monsters: []engine.CardRef{"leech", "fly"},
		Rooms:    []engine.CardRef{"blessing_of_the_sack"},
	})
	cents := tb.G.Players[0].Cents
	tb.Play(0, "a_penny")
	if tb.G.Players[0].Cents != cents+1 {
		t.Fatalf("A Penny! gives 1¢ with no question: %d", tb.G.Players[0].Cents)
	}
	kill(tb, "fly", "ask.yes")
	if got := tb.G.PresentCounters(engine.PresentCents); got != 10 {
		t.Errorf("¢ present %d, want 8 + 2 doubled", got)
	}
}

// Q-33: after Happy Holidays! there are no presents: no question.
func TestNoForgoAfterHolidays(t *testing.T) {
	tb := krampus(t, engine.DifficultyNormal, engine.SituationPlayer{}, engine.SituationPlayer{})
	tb.G.FlipPresents()
	cents := tb.G.Players[0].Cents
	kill(tb, "fly")
	if tb.G.Players[0].Cents != cents+1 {
		t.Errorf("cents %d", tb.G.Players[0].Cents)
	}
}

// R-KRAMPUS-08 (Normal): a present at 0 is destroyed into the treasure
// discard, face up; then the top 3 monster cards become minions, one at
// a time. R-KRAMPUS-09: counters stop at 0; a destroyed present takes
// none.
func TestWhatsInTheBox(t *testing.T) {
	tb := krampus(t, engine.DifficultyNormal, engine.SituationPlayer{}, engine.SituationPlayer{})
	setCounters(tb, engine.PresentCents, 2)
	card := tb.G.Object(tb.G.Challenge.Presents[0].Object).Card
	tb.G.PutOnDeck(engine.MonsterDeck, "pooter", "dip", "squirt")
	tb.G.ForceRolls(1, 1, 1)
	tb.EndTurn()
	p := tb.G.Challenge.Presents[0]
	d := tb.G.Discards[engine.TreasureDeck]
	if !p.Destroyed || len(d) == 0 || d[len(d)-1] != p.Object || tb.G.Object(p.Object).Card != card {
		t.Errorf("the ¢ present is not in the treasure discard: %+v", p)
	}
	if got := minionCards(tb.G); !slices.Equal(got, []engine.CardRef{"squirt", "dip", "pooter"}) {
		t.Errorf("minions %v", got)
	}
	if n := tb.G.AddPresentCounters(engine.PresentCents, -4, engine.NoPlayer); n != 0 || tb.G.PresentCounters(engine.PresentCents) != 0 {
		t.Error("a destroyed present took counters")
	}
	v := tb.G.View(engine.Viewer{Kind: engine.ViewSpectator}).Challenge.Presents[0]
	if !v.Destroyed || v.Card != card || v.Counters != 0 {
		t.Errorf("view %+v", v)
	}
	check(t, tb.G)
}

// R-KRAMPUS-01, R-KRAMPUS-09: the fourth destroyed present loses the
// challenge at once, with its reason; no minions come for it.
func TestAllPresentsLost(t *testing.T) {
	tb := krampus(t, engine.DifficultyNormal, engine.SituationPlayer{}, engine.SituationPlayer{})
	for i, kind := range []string{engine.PresentCents, engine.PresentLoot, engine.PresentTreasure, engine.PresentSoul} {
		tb.G.PutOnDeck(engine.MonsterDeck, "pooter", "dip", "squirt")
		tb.G.AddPresentCounters(kind, -20, engine.NoPlayer)
		tb.G.Run()
		ev := tb.Settle()
		if i < 3 && (tb.G.Over || len(tb.G.Minions) != 3*(i+1)) {
			t.Fatalf("present %d: over %v, %d minions", i+1, tb.G.Over, len(tb.G.Minions))
		}
		if i == 3 && count(ev, engine.EvChallengeLost) != 1 {
			t.Error("no challenge_lost")
		}
	}
	c := tb.G.Challenge
	if !tb.G.Over || !c.Lost || c.Reason.Key != challenge.ReasonPresents || len(tb.G.Minions) != 9 || len(tb.G.Winners) != 0 {
		t.Errorf("over %v, lost %v (%v), %d minions", tb.G.Over, c.Lost, c.Reason, len(tb.G.Minions))
	}
	check(t, tb.G)
}

// R-KRAMPUS-10, Q-30, Q-40 (Normal): the first time both players hold 5
// souls together, every present still there flips, loses its counters
// and becomes an item; the active player picks who gains each, in
// order ¢, loot, treasure, soul. A destroyed present is skipped. It
// happens once.
func TestHappyHolidays(t *testing.T) {
	tb := krampus(t, engine.DifficultyNormal, engine.SituationPlayer{Souls: souls(3)}, engine.SituationPlayer{Souls: souls(1)})
	tb.G.DestroyPresent(engine.PresentLoot)
	isaac, cain := player(tb, 0), player(tb, 1)
	boost(tb, 0, 5)
	tb.G.ForceRolls(6)
	ev := tb.RevealFromDeck("mixed_bag", "ask.no", "ask.no", cain, isaac, cain)
	if count(ev, engine.EvPresentsFlipped) != 1 || count(ev, engine.EvPresentGiven) != 3 || !tb.G.Challenge.Holidays {
		t.Fatalf("flipped %d, given %d", count(ev, engine.EvPresentsFlipped), count(ev, engine.EvPresentGiven))
	}
	for i, want := range []engine.PlayerID{1, engine.NoPlayer, 0, 1} {
		p := tb.G.Challenge.Presents[i]
		if want == engine.NoPlayer {
			if !p.Destroyed || p.Flipped {
				t.Errorf("%s present %+v", p.Kind, p)
			}
			continue
		}
		o := tb.G.Object(p.Object)
		if !p.Flipped || o.Controller != want || o.Role != engine.RoleItem || !o.Charged || len(o.Counters) != 0 ||
			!slices.Contains(tb.G.Players[want].InPlay, p.Object) {
			t.Errorf("%s present: %+v", p.Kind, o)
		}
		v := tb.G.View(engine.Viewer{Kind: engine.ViewSpectator}).Challenge.Presents[i]
		if !v.Flipped || v.Player != want || v.Card != o.Card {
			t.Errorf("%s present view %+v", p.Kind, v)
		}
	}
	tb.G.Turn.Attacks = 1
	tb.G.ForceRolls(6)
	if ev := tb.RevealFromDeck("mixed_bag"); count(ev, engine.EvPresentsFlipped) != 0 || tb.G.SoulsInPlay() != 6 {
		t.Error("Happy Holidays! twice, or a forgo question with no presents")
	}
	check(t, tb.G)
}

// R-KRAMPUS-01: Pilfering Krampus dies: both players win at once.
func TestKrampusGoal(t *testing.T) {
	tb := krampus(t, engine.DifficultyNormal, engine.SituationPlayer{Souls: souls(5)}, engine.SituationPlayer{})
	boost(tb, 0, 4)
	ev := tb.Attack(normBoss, 6)
	if !tb.G.Over || tb.G.Challenge.Lost || len(tb.G.Winners) != 2 || count(ev, engine.EvChallengeWon) != 1 {
		t.Errorf("over %v, lost %v, winners %v", tb.G.Over, tb.G.Challenge.Lost, tb.G.Winners)
	}
}

// R-ABIL-33, R-KRAMPUS-07: with another replacement (Lucky Toe), Isaac
// orders them; the forgo is still asked, for the loot as it is then.
func TestForgoWithLuckyToe(t *testing.T) {
	tb := krampus(t, engine.DifficultyNormal, engine.SituationPlayer{Items: []engine.CardRef{"lucky_toe"}}, engine.SituationPlayer{})
	kill(tb, "leech", "lucky_toe", "ask.yes")
	if got := tb.G.PresentCounters(engine.PresentLoot); got != 4 {
		t.Errorf("loot present %d, want 2 + Lucky Toe's 2", got)
	}
	tb = krampus(t, engine.DifficultyNormal, engine.SituationPlayer{Items: []engine.CardRef{"lucky_toe"}}, engine.SituationPlayer{})
	hand := len(tb.G.Players[0].Hand)
	kill(tb, "leech", "how_the_krampus_stole_christmas_norm", "ask.no")
	if got := len(tb.G.Players[0].Hand); got != hand+2 || tb.G.PresentCounters(engine.PresentLoot) != 2 {
		t.Errorf("declined: hand %d, loot present %d", got, tb.G.PresentCounters(engine.PresentLoot))
	}
	check(t, tb.G)
}
