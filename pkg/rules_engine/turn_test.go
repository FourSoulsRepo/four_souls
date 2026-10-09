package rulesengine

import (
	"errors"
	"reflect"
	"testing"
)

func newTestGame(t *testing.T, players int) *Game {
	t.Helper()
	g, _, err := NewGame(Setup{Seed: 1, Players: players, Sets: []CardSet{testSet}, BonusSouls: true})
	if err != nil {
		t.Fatal(err)
	}
	return g
}

// passAll makes every player pass until the prompt is no longer priority
// for someone, or until limit passes.
func passUntilStep(t *testing.T, g *Game, step Step) {
	t.Helper()
	for range 50 {
		if g.Turn.Step == step && g.Turn.Entered {
			return
		}
		p := g.Prompt()
		if p.Kind != PromptPriority {
			t.Fatalf("waiting for %v, want to pass to step %d", p, step)
		}
		if _, err := g.Apply(Intent{Player: p.Player, Kind: IntentPass}); err != nil {
			t.Fatal(err)
		}
	}
	t.Fatalf("never reached step %d", step)
}

func TestSetup(t *testing.T) {
	g := newTestGame(t, 3)
	for _, pl := range g.Players {
		if len(pl.Hand) != 3 || pl.Cents != 3 {
			t.Errorf("player %d: hand %d, cents %d (R-SETUP-10)", pl.ID, len(pl.Hand), pl.Cents)
		}
		if g.Object(pl.Character).Role != RoleCharacter {
			t.Errorf("player %d has no character (R-SETUP-07)", pl.ID)
		}
	}
	if len(g.Shop) != 2 || len(g.Monsters) != 2 {
		t.Fatalf("shop %d, monsters %d (R-SETUP-03, R-SETUP-04)", len(g.Shop), len(g.Monsters))
	}
	for i, s := range g.Monsters {
		top, _ := s.TopOf()
		if g.kindOf(top) == EventCard {
			t.Errorf("monster slot %d holds an event after setup (R-SETUP-05)", i)
		}
	}
	if len(g.BonusSouls) != 3 {
		t.Errorf("bonus souls = %d, want 3 (R-SETUP-06)", len(g.BonusSouls))
	}
	for _, id := range g.Decks[TreasureDeck] {
		if c := g.Object(id).Card; c == "item_a" || c == "item_b" {
			t.Error("starting items must not be in the treasure deck")
		}
	}
}

func TestCharactersStartDeactivatedItemsCharged(t *testing.T) {
	// Setup state before the first recharge: build without running.
	g := newTestGame(t, 2)
	active := g.Turn.Active
	other := g.next(active)
	// The first turn already recharged the active player (R-TURN-02);
	// the other player still shows the setup state (R-SETUP-08).
	if g.Object(g.Players[other].Character).Charged {
		t.Error("characters start deactivated (R-SETUP-08)")
	}
	for _, id := range g.Players[other].InPlay {
		if !g.Object(id).Charged {
			t.Error("starting items start charged (R-SETUP-08)")
		}
	}
	if !g.Object(g.Players[active].Character).Charged {
		t.Error("the active player recharged in the recharge step (R-TURN-02)")
	}
}

func TestFirstPromptIsStartOfTurnPriority(t *testing.T) {
	g := newTestGame(t, 3)
	p := g.Prompt()
	if p.Kind != PromptPriority || p.Player != g.Turn.Active || g.Turn.Step != StepStartTriggers {
		t.Fatalf("prompt %+v at step %d; want priority for the active player at start (R-TURN-03)", p, g.Turn.Step)
	}
}

func TestTurnCycle(t *testing.T) {
	g := newTestGame(t, 3)
	first := g.Turn.Active
	hand := len(g.Players[first].Hand)

	passUntilStep(t, g, StepAction)
	if got := len(g.Players[first].Hand); got != hand+1 {
		t.Errorf("hand %d -> %d; the loot step loots 1 (R-TURN-04)", hand, got)
	}
	if g.Turn.LootPlays != 1 || g.Turn.Attacks != 1 || g.Turn.Purchases != 1 {
		t.Errorf("turn allowances %+v (R-TURN-05, R-TURN-07)", g.Turn)
	}

	// The active player cannot just pass in an open action phase.
	_, err := g.Apply(Intent{Player: first, Kind: IntentPass})
	var re *RuleError
	if !errors.As(err, &re) || re.Rule != "R-TURN-09" {
		t.Errorf("pass in action phase: %v, want R-TURN-09", err)
	}
	// Someone else cannot end the turn.
	if _, err := g.Apply(Intent{Player: g.next(first), Kind: IntentEndTurn}); err == nil {
		t.Error("a non-active player ended the turn")
	}
	if _, err := g.Apply(Intent{Player: first, Kind: IntentEndTurn}); err != nil {
		t.Fatal(err)
	}
	// The others may respond to the end declaration, then the end phase starts.
	if p := g.Prompt(); p.Player != g.next(first) {
		t.Fatalf("after ending, priority goes to the next player; got %+v", p)
	}
	passUntilStep(t, g, StepEndTriggers)
	passUntilStep(t, g, StepStartTriggers)
	if g.Turn.Active != g.next(first) || g.Turn.Number != 2 {
		t.Errorf("turn %d active %d; want turn 2 for player %d (R-TURN-01)", g.Turn.Number, g.Turn.Active, g.next(first))
	}
}

func TestHandSizeDiscard(t *testing.T) {
	g := newTestGame(t, 2)
	p := g.Turn.Active
	g.loot(p, 9) // 3 + 9 = 12 cards; the loot step makes 13
	passUntilStep(t, g, StepAction)
	if _, err := g.Apply(Intent{Player: p, Kind: IntentEndTurn}); err != nil {
		t.Fatal(err)
	}
	passUntilStep(t, g, StepEndTriggers)
	for g.Prompt().Kind == PromptPriority {
		if _, err := g.Apply(Intent{Player: g.Prompt().Player, Kind: IntentPass}); err != nil {
			t.Fatal(err)
		}
	}
	pr := g.Prompt()
	if pr.Kind != PromptDiscard || pr.Player != p || pr.Count != 3 {
		t.Fatalf("prompt %+v; want discard 3 down to 10 (R-TURN-11)", pr)
	}
	hand := g.Players[p].Hand
	if _, err := g.Apply(Intent{Player: p, Kind: IntentDiscard, Objects: hand[:2]}); err == nil {
		t.Error("discarding too few cards was accepted")
	}
	if _, err := g.Apply(Intent{Player: p, Kind: IntentDiscard, Objects: []ObjectID{hand[0], hand[0], hand[1]}}); err == nil {
		t.Error("discarding the same card twice was accepted")
	}
	if _, err := g.Apply(Intent{Player: p, Kind: IntentDiscard, Objects: append([]ObjectID(nil), hand[:3]...)}); err != nil {
		t.Fatal(err)
	}
	if len(g.Players[p].Hand) != 10 || len(g.Discards[LootDeck]) != 3 {
		t.Errorf("hand %d, discard %d", len(g.Players[p].Hand), len(g.Discards[LootDeck]))
	}
	if g.Turn.Active == p {
		t.Error("the turn should have passed after the hand-size step")
	}
}

func TestWinWithFourSouls(t *testing.T) {
	g := newTestGame(t, 2)
	p := g.Turn.Active
	for range 4 {
		id := g.newObject("boss", Zone{Kind: ZoneInPlay}, p)
		g.Object(id).Role = RoleSoul
		g.Players[p].InPlay = append(g.Players[p].InPlay, id)
	}
	g.Waiting = Prompt{}
	g.run()
	if !g.Over || !reflect.DeepEqual(g.Winners, []PlayerID{p}) {
		t.Fatalf("over %v winners %v; 4 souls win (R-WIN-01)", g.Over, g.Winners)
	}
	if _, err := g.Apply(Intent{Player: p, Kind: IntentPass}); err == nil {
		t.Error("intents after the game ended were accepted")
	}
}

func TestSameSeedSameGame(t *testing.T) {
	play := func() uint64 {
		g := newTestGame(t, 4)
		for range 40 {
			p := g.Prompt()
			in := Intent{Player: p.Player, Kind: IntentPass}
			if g.inOpenActionPhase() && p.Player == g.Turn.Active {
				in.Kind = IntentEndTurn
			}
			if p.Kind == PromptDiscard {
				in = Intent{Player: p.Player, Kind: IntentDiscard, Objects: g.Players[p.Player].Hand[:p.Count]}
			}
			if _, err := g.Apply(in); err != nil {
				t.Fatal(err)
			}
		}
		sum, err := g.Checksum()
		if err != nil {
			t.Fatal(err)
		}
		return sum
	}
	if a, b := play(), play(); a != b {
		t.Errorf("same seed and intents gave different states (A-08)")
	}
}

func TestCloneKeepsCardDefinitions(t *testing.T) {
	g := newTestGame(t, 2)
	c, err := g.Clone()
	if err != nil {
		t.Fatal(err)
	}
	top, _ := c.Monsters[0].TopOf()
	if c.kindOf(top) != MonsterCard {
		t.Error("a clone must still know card definitions")
	}
}
