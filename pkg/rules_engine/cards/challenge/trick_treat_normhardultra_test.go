package challenge_test

import (
	"slices"
	"strings"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/challenge"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// R-TRICK-05, R-TRICK-11, Q-43, Q-44: as each turn starts, the active
// player may flip the rule card; only the side up works; a no keeps it.
func TestTrickFlip(t *testing.T) {
	tb := trick(t, engine.DifficultyNormal, engine.SituationPlayer{}, engine.SituationPlayer{})
	g := tb.G
	if _, err := g.Apply(engine.Intent{Player: 0, Kind: engine.IntentEndTurn}); err != nil {
		t.Fatal(err)
	}
	for range 100 {
		tb.ThrowPending()
		if w := g.Prompt(); w.Kind == engine.PromptPriority {
			tb.Pass(int(w.Player))
			continue
		}
		break
	}
	if w := g.Prompt(); w.Player != 1 || g.Turn.Active != 1 || w.Text.Key.String() != "card.trick_treat.flip_to_treat" {
		t.Fatalf("flip asked of %d on %d's turn: %+v", w.Player, g.Turn.Active, w.Text)
	}
	settle(tb, "ask.yes")
	for range 100 {
		if g.Turn.Step == engine.StepAction && g.Prompt().Kind == engine.PromptPriority && len(g.Stack) == 0 {
			break
		}
		tb.ThrowPending()
		tb.Pass(int(g.Prompt().Player))
	}
	if g.CardOf(g.Challenge.Card) != treatCard || !g.Object(g.Challenge.Card).Flipped {
		t.Fatal("the rule card did not flip")
	}
	// Treat now: a kill removes no trick counter; the roll adds 3¢.
	ev := attack(tb, "fly", 1, 6, 1)
	if tricks(tb) != 12 || g.Players[1].Cents != 1+3 {
		t.Errorf("Treat kill: tricks %d, Cain %d¢; want 12 and 4¢", tricks(tb), g.Players[1].Cents)
	}
	if rollsOf(tb, ev, treatCard) != 1 || count(ev, engine.EvAbilitiesGained) != 0 {
		t.Errorf("Treat rolls %d", rollsOf(tb, ev, treatCard))
	}
	tb.EndTurn("ask.no")
	if g.CardOf(g.Challenge.Card) != treatCard {
		t.Error("a no flipped it")
	}
	tb.EndTurn("ask.yes")
	if g.CardOf(g.Challenge.Card) != trickCard {
		t.Error("a yes on Treat did not flip it back")
	}
}

// R-TRICK-12, Q-54: with Trick up, every monster death removes a trick
// counter from the final boss, in an attack or not.
func TestNobodysFool(t *testing.T) {
	tb := trick(t, engine.DifficultyNormal, engine.SituationPlayer{Hand: []engine.CardRef{"bomb"}}, engine.SituationPlayer{})
	attack(tb, "fly", 6, 6, 1)
	if tricks(tb) != 11 {
		t.Errorf("an attack kill: tricks %d, want 11", tricks(tb))
	}
	tb.G.ForceRolls(1)
	tb.ThrowPending()
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentPlayLoot, Objects: []engine.ObjectID{tb.Find(0, "bomb")}}); err != nil {
		t.Fatal(err)
	}
	settle(tb, "leech")
	if tricks(tb) != 10 {
		t.Errorf("a Bomb! kill: tricks %d, want 10", tricks(tb))
	}
	check(t, tb.G)
}

// R-TRICK-13, Q-46, Q-53: Ghostly Trickery rolls once on Normal, twice
// on Hard, three times on Ultra Hard, before the first attack roll; the
// active player rolls.
func TestGhostlyTrickeryRolls(t *testing.T) {
	for d, n := range map[string]int{engine.DifficultyNormal: 1, engine.DifficultyHard: 2, engine.DifficultyUltra: 3} {
		tb := trick(t, d, engine.SituationPlayer{}, engine.SituationPlayer{})
		rolls := append(slices.Repeat([]int{6}, n), 6, 1)
		ev := attack(tb, "fly", rolls...)
		first := enginetest.FirstAttackRoll(ev)
		if got := rollsOf(tb, ev, trickCard); got != n {
			t.Errorf("%s: %d rolls, want %d", d, got, n)
		}
		for i, e := range ev {
			if e.Kind == engine.EvStackResolved && e.Card == trickCard && strings.Contains(e.Label.Key.String(), "roll") && i > first {
				t.Errorf("%s: a Trick roll after the first attack roll", d)
			}
		}
		if p := tb.G.Players[0]; p.Cents != 1 {
			t.Errorf("%s: %d¢, want fly's 1¢ as 1¢", d, p.Cents)
		}
	}
}

// R-TRICK-13, Q-46: 1-2 makes each damage the active player takes 1
// more till end of turn; two such results add up (Hard). Daddy Haunt
// adds its own 1 (one rule, StatDamageTaken). The boost lasts the turn
// (TestIncreaseDamageTaken in cards/trick_blocks_test.go).
func TestGhostlyTrickeryDamage(t *testing.T) {
	tb := trick(t, engine.DifficultyHard, engine.SituationPlayer{}, engine.SituationPlayer{})
	hp(tb, 0, 5)
	ev := attack(tb, "fly", 1, 2, 1, 6, 1)
	if got := damageTo(ev, 0); len(got) != 1 || got[0] != 3 {
		t.Errorf("Hard, 1 and 2: damage %v, want [3]", got)
	}
	tb = trick(t, engine.DifficultyNormal, engine.SituationPlayer{Items: []engine.CardRef{"daddy_haunt"}}, engine.SituationPlayer{})
	hp(tb, 0, 5)
	ev = attack(tb, "fly", 1, 1, 6, 1)
	if got := damageTo(ev, 0); len(got) != 1 || got[0] != 3 {
		t.Errorf("Daddy Haunt and 1: damage %v, want [3]", got)
	}
}

// R-TRICK-13, Q-47: 3-4 reveals monster cards until a non-event one;
// the attacked monster gains its abilities, "when this dies" too:
// Boom Fly's damage to each player then brings fear (R-TRICK-07).
func TestGhostlyTrickeryGain(t *testing.T) {
	tb := trick(t, engine.DifficultyNormal, engine.SituationPlayer{}, engine.SituationPlayer{})
	tb.G.PutOnDeck(engine.MonsterDeck, "boom_fly")
	ev := attack(tb, "fly", 3, 6, 1)
	gained := slices.IndexFunc(ev, func(e engine.Event) bool { return e.Kind == engine.EvAbilitiesGained })
	if gained < 0 || ev[gained].Text != "boom_fly" || ev[gained].Card != "fly" {
		t.Fatalf("no abilities_gained of boom_fly: %+v", ev)
	}
	if got := damageTo(ev, 0); len(got) != 1 || got[0] != 1 {
		t.Errorf("Isaac took %v, want Boom Fly's 1", got)
	}
	if fear(tb, 0) != 1 || fear(tb, 1) != 1 || tb.G.Players[0].Cents != 1 {
		t.Errorf("fear %d and %d, %d¢; want 1, 1 and fly's 1¢", fear(tb, 0), fear(tb, 1), tb.G.Players[0].Cents)
	}
	check(t, tb.G)
}

// R-TRICK-13, Q-46, Q-48: 5-6 makes the rewards 1¢, twice still 1¢;
// Dinga killed on a 6 doubles the 1¢; Blessing of the Sack too.
func TestGhostlyTrickeryCents(t *testing.T) {
	tb := trick(t, engine.DifficultyHard, engine.SituationPlayer{}, engine.SituationPlayer{}, "conjoined_fatty", "dinga")
	boost(tb, 0, 3)
	attack(tb, "conjoined_fatty", 5, 6, 6, 1)
	if p := tb.G.Players[0]; p.Cents != 1 || len(p.Hand) != 0 {
		t.Errorf("Conjoined Fatty: %d¢, %d loot; want 1¢ only", p.Cents, len(p.Hand))
	}
	tb.G.Turn.Attacks++
	attack(tb, "dinga", 6, 6, 6, 1)
	if c := tb.G.Players[0].Cents; c != 1+2 {
		t.Errorf("Dinga on a 6: %d¢, want 1 + 1¢ doubled", c)
	}
	tb = trickOn(t, engine.DifficultyNormal, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{}, {}},
		Monsters: []engine.CardRef{"leech", "fly"}, Rooms: []engine.CardRef{"blessing_of_the_sack"},
	})
	attack(tb, "leech", 5, 6, 1)
	if p := tb.G.Players[0]; p.Cents != 2 || len(p.Hand) != 0 {
		t.Errorf("Blessing of the Sack: %d¢, %d loot; want 2¢", p.Cents, len(p.Hand))
	}
}

// R-TRICK-14, Q-51, Q-54: with Treat up, each time the active player
// gains a monster's rewards, they remove a fear counter from a player;
// a monster without rewards gives none.
func TestYouDontScareMe(t *testing.T) {
	tb := trick(t, engine.DifficultyNormal, engine.SituationPlayer{Hand: []engine.CardRef{"bomb"}}, engine.SituationPlayer{}, "fly", "bare_mob")
	flipToTreat(tb)
	tb.G.AddPlayerCounters(1, challenge.FearCounter, 2)
	tb.G.ForceRolls(1)
	tb.ThrowPending()
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentPlayLoot, Objects: []engine.ObjectID{tb.Find(0, "bomb")}}); err != nil {
		t.Fatal(err)
	}
	settle(tb, "bare_mob")
	if asking(tb) {
		t.Fatalf("a question for a monster without rewards: %v", tb.Options())
	}
	if fear(tb, 1) != 2 {
		t.Errorf("no rewards: Cain's fear %d, want 2", fear(tb, 1))
	}
	attackAsk(tb, "fly", 1, 6, 1)
	if opts := tb.Options(); !slices.Equal(opts, []string{player(tb, 1)}) {
		t.Fatalf("options %v, want Cain only", opts)
	}
	settle(tb, player(tb, 1))
	if fear(tb, 1) != 1 || tb.G.Players[0].Cents != 1+3 {
		t.Errorf("fear %d, %d¢; want 1 and 4¢", fear(tb, 1), tb.G.Players[0].Cents)
	}
}

// R-TRICK-15, Q-46.7, Q-48: Treat rolls once at every difficulty; the
// result adds to the attacked monster's rewards.
func TestHallowsEveBounty(t *testing.T) {
	for roll, want := range map[int][3]int{1: {1 + 3, 1, 0}, 3: {1, 2, 0}, 5: {1, 1, 1}} {
		tb := trick(t, engine.DifficultyUltra, engine.SituationPlayer{}, engine.SituationPlayer{}, "leech", "fly")
		flipToTreat(tb)
		tb.G.PutOnDeck(engine.TreasureDeck, "breakfast")
		ev := attack(tb, "fly", roll, 6, 1)
		p := tb.G.Players[0]
		hand := len(p.Hand)
		if got := [3]int{p.Cents, hand + 1, len(p.InPlay)}; got != want {
			t.Errorf("roll %d: ¢, loot+1, items %v; want %v", roll, got, want)
		}
		if n := rollsOf(tb, ev, treatCard); n != 1 {
			t.Errorf("roll %d: %d Treat rolls on Ultra Hard, want 1", roll, n)
		}
	}
}

// R-TRICK-01: the boss dies; the challenge is won, both players win.
// Both sides of the rule card hold the goal.
func TestTrickGoal(t *testing.T) {
	for _, treat := range []bool{false, true} {
		tb := trick(t, engine.DifficultyNormal, engine.SituationPlayer{}, engine.SituationPlayer{})
		if treat {
			flipToTreat(tb)
			tb.G.PutOnDeck(engine.TreasureDeck, "breakfast")
		}
		setTricks(tb, 0)
		boost(tb, 0, 3)
		attack(tb, hauntNorm, 6, 6)
		if !tb.G.Over || tb.G.Challenge.Lost || len(tb.G.Winners) != 2 {
			t.Errorf("treat %v: over %v, lost %v, winners %v", treat, tb.G.Over, tb.G.Challenge.Lost, tb.G.Winners)
		}
	}
}
