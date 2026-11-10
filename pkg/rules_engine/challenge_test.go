package rulesengine_test

import (
	"bytes"
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// trialSet is a tiny set for a test challenge (ADR 028): a rule card
// whose goal is trivial, a boss, and just enough cards to play turns.
var trialSet = engine.CardSet{Code: "trial", Name: "trial", Cards: []engine.CardDef{
	{Ref: "trial_hero", Kind: engine.CharacterCard, HP: 2, ATK: 1, Copies: 2},
	// trial_saint's soul value alone reaches the soul target (R-CARD-35).
	{Ref: "trial_saint", Kind: engine.CharacterCard, HP: 2, ATK: 1, Soul: 4},
	{Ref: "trial_item", Kind: engine.TreasureCard, Copies: 6},
	// trial_crowd has the 3+ player symbol (R-CARD-32).
	{Ref: "trial_crowd", Kind: engine.TreasureCard, Copies: 6, ThreePlus: true},
	{Ref: "trial_penny", Kind: engine.LootCard, Copies: 60},
	{Ref: "trial_fly", Kind: engine.MonsterCard, HP: 1, DC: 6, ATK: 0, Copies: 4},
	{Ref: "trial_room", Kind: engine.RoomCard, Copies: 8},
	{Ref: "trial_boss", Kind: engine.MonsterCard, HP: 3, DC: 4, ATK: 1, Outside: true},
	// The goal: hold 10¢ as a turn ends.
	{Ref: "trial_rules", Kind: engine.ChallengeCard, Abilities: []engine.Ability{{
		Kind: engine.Triggered, Text: "At the end of each turn, if the active player has 10¢, the players win.",
		Trigger: engine.Trigger{On: engine.EvEndOfTurn, Match: func(g *engine.Game, _ engine.ObjectID, e engine.Event) bool {
			return g.Players[e.Player].Cents >= 10
		}},
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.WinChallenge() })},
	}}},
}}

// trial is the test challenge; its own setup gives seat 0 2¢.
var trial = engine.ChallengeDef{
	ID: "trial", Set: "trial",
	Levels: []engine.ChallengeLevel{
		{Difficulty: engine.DifficultyNormal, Boss: "trial_boss", Card: "trial_rules"},
		{Difficulty: engine.DifficultyHard, Boss: "trial_boss", Card: "trial_rules"},
	},
	Setup: func(g *engine.Game, _ string) error {
		g.Players[0].Cents += 2
		return nil
	},
}

func trialSetup(seed uint64) engine.Setup {
	return engine.Setup{
		Seed: seed, Players: 2, Sets: []engine.CardSet{trialSet},
		Challenge: &trial, Difficulty: engine.DifficultyNormal,
		Characters: []engine.CardRef{"trial_hero", "trial_hero"},
	}
}

// newTrial starts a test challenge and runs it to turn 1's action phase.
func newTrial(t *testing.T, s engine.Setup) *engine.Game {
	t.Helper()
	g, _, err := engine.NewGame(s)
	if err != nil {
		t.Fatal(err)
	}
	enginetest.ThrowFirstRoll(t, g)
	untilAction(t, g)
	return g
}

// untilAction answers prompts until the active player may act in an
// open action phase, or the game is over. It returns the events.
func untilAction(t *testing.T, g *engine.Game) []engine.Event {
	t.Helper()
	var out []engine.Event
	for range 1000 {
		w := g.Prompt()
		if g.Over || (g.Turn.Step == engine.StepAction && w.Kind == engine.PromptPriority && w.Player == g.Turn.Active && len(g.Stack) == 0 && !g.Turn.EndDeclared) {
			return out
		}
		in := engine.Intent{Player: w.Player, Kind: engine.IntentPass}
		switch w.Kind { //nolint:exhaustive // the trial set asks nothing else
		case engine.PromptDiscard:
			in = engine.Intent{Player: w.Player, Kind: engine.IntentDiscard, Objects: g.Players[w.Player].Hand[:w.Count]}
		case engine.PromptThrow:
			in.Kind = engine.IntentThrow
		case engine.PromptPriority:
		default:
			t.Fatalf("unexpected prompt %+v", w)
		}
		ev, err := g.Apply(in)
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, ev...)
		if err := g.CheckInvariants(); err != nil {
			t.Fatal(err)
		}
	}
	t.Fatal("the game never got to an action phase")
	return nil
}

// endTurn ends the active player's turn and runs to the next action
// phase, or to the game's end.
func endTurn(t *testing.T, g *engine.Game) []engine.Event {
	t.Helper()
	ev, err := g.Apply(engine.Intent{Player: g.Turn.Active, Kind: engine.IntentEndTurn})
	if err != nil {
		t.Fatal(err)
	}
	return append(ev, untilAction(t, g)...)
}

func countKind(events []engine.Event, kind engine.EventKind) int {
	n := 0
	for _, e := range events {
		if e.Kind == kind {
			n++
		}
	}
	return n
}

// R-COOP-01, R-COOP-02, ADR 028 Decision 4: a challenge needs 2
// players, no bonus souls, its card set, and one of its difficulties.
func TestChallengeSetupChecks(t *testing.T) {
	for name, change := range map[string]func(*engine.Setup){
		"3 players (R-COOP-01)":   func(s *engine.Setup) { s.Players, s.Characters = 3, nil },
		"bonus souls (R-COOP-02)": func(s *engine.Setup) { s.BonusSouls = true },
		"no challenge set": func(s *engine.Setup) {
			s.Sets = []engine.CardSet{{Code: "other", Name: "other", Cards: trialSet.Cards}}
		},
		"unknown difficulty": func(s *engine.Setup) { s.Difficulty = engine.DifficultyUltra },
		"no difficulty":      func(s *engine.Setup) { s.Difficulty = "" },
		"rule card is no challenge": func(s *engine.Setup) {
			c := trial
			c.Levels = []engine.ChallengeLevel{{Difficulty: "normal", Card: "trial_item"}}
			s.Challenge = &c
		},
	} {
		s := trialSetup(1)
		change(&s)
		if _, _, err := engine.NewGame(s); err == nil {
			t.Errorf("%s: the challenge started", name)
		}
	}
}

// R-COOP-02, R-CHAL-01, ADR 028 Decision 5: the co-op setup, then the
// timer at 8, the rule card outside the game, and the challenge's own
// setup.
func TestChallengeSetup(t *testing.T) {
	g, _, err := engine.NewGame(trialSetup(1))
	if err != nil {
		t.Fatal(err)
	}
	c := g.Challenge
	if c == nil || c.ID != "trial" || c.Difficulty != engine.DifficultyNormal || c.Timer != 8 {
		t.Fatalf("challenge %+v", c)
	}
	if o := g.Object(c.Card); o.Card != "trial_rules" || o.Zone.Kind != engine.ZoneOutside {
		t.Errorf("rule card %+v", o)
	}
	if g.Players[0].Cents != 5 || g.Players[1].Cents != 3 {
		t.Errorf("cents %d, %d: the challenge's own setup did not run", g.Players[0].Cents, g.Players[1].Cents)
	}
	if len(g.BonusSouls) != 0 || len(g.Players[0].Hand) != 3 || len(g.Shop) != 2 || len(g.Monsters) != 2 {
		t.Errorf("not the co-op setup: bonus souls %d, hand %d, shop %d, monsters %d", len(g.BonusSouls), len(g.Players[0].Hand), len(g.Shop), len(g.Monsters))
	}
	if err := g.CheckInvariants(); err != nil {
		t.Fatal(err)
	}
}

// R-COOP-03: cards with the 3+ player symbol are out of a challenge.
func TestChallengeLeavesThreePlusOut(t *testing.T) {
	for seed := range uint64(5) {
		g, _, err := engine.NewGame(trialSetup(seed))
		if err != nil {
			t.Fatal(err)
		}
		for _, o := range g.Objects {
			if o.Card == "trial_crowd" {
				t.Fatalf("seed %d: trial_crowd is in the game, zone %+v", seed, o.Zone)
			}
		}
	}
}

// R-COOP-08, R-COOP-09, R-CHAL-13, R-CHAL-15: the timer ticks once
// per round, after both players' turns; at 0 both players lose at
// once, before another turn starts.
func TestChallengeTimerLoses(t *testing.T) {
	g := newTrial(t, trialSetup(3))
	starter := g.Turn.Active
	for round := 1; round <= 8; round++ {
		if ev := endTurn(t, g); countKind(ev, engine.EvChallengeTimer) != 0 {
			t.Fatalf("round %d: the timer ticked after the first turn", round)
		}
		ev := endTurn(t, g)
		if countKind(ev, engine.EvChallengeTimer) != 1 || g.Challenge.Timer != 8-round {
			t.Fatalf("round %d: timer %d, ticks %d", round, g.Challenge.Timer, countKind(ev, engine.EvChallengeTimer))
		}
		if round < 8 && (g.Over || g.Turn.Active != starter) {
			t.Fatalf("round %d: over %v, active %d", round, g.Over, g.Turn.Active)
		}
		if round == 8 {
			i := slices.IndexFunc(ev, func(e engine.Event) bool { return e.Kind == engine.EvChallengeLost })
			if i < 0 || ev[i].Label.Key != engine.ReasonTimer || countKind(ev, engine.EvTurnStarted) != 0 {
				t.Fatalf("no loss at once: %+v", ev)
			}
		}
	}
	c := g.Challenge
	if !g.Over || len(g.Winners) != 0 || !c.Lost || c.Reason.Key != engine.ReasonTimer || g.Prompt().Kind != engine.PromptGameOver {
		t.Errorf("over %v, winners %v, challenge %+v, prompt %+v", g.Over, g.Winners, c, g.Prompt())
	}
	if err := g.CheckInvariants(); err != nil {
		t.Fatal(err)
	}
}

// R-COOP-08, R-TURN-16, Q-34: an extra turn adds no tick; a skipped
// turn counts as taken.
func TestChallengeTimerExtraAndSkippedTurns(t *testing.T) {
	g := newTrial(t, trialSetup(3))
	starter := g.Turn.Active
	g.ExtraTurn = true // the starting player takes another turn
	endTurn(t, g)
	endTurn(t, g)
	if g.Challenge.Timer != 8 || g.Turn.Active == starter {
		t.Fatalf("after the starter's extra turn: timer %d, active %d", g.Challenge.Timer, g.Turn.Active)
	}
	g.ExtraTurn = true // the partner takes another turn: the round ends after it
	if endTurn(t, g); g.Challenge.Timer != 8 {
		t.Fatalf("the timer ticked before the partner's extra turn: %d", g.Challenge.Timer)
	}
	if endTurn(t, g); g.Challenge.Timer != 7 || g.Turn.Active != starter {
		t.Fatalf("after the round: timer %d, active %d", g.Challenge.Timer, g.Turn.Active)
	}
	partner := g.Turn.Active ^ 1
	g.Players[partner].SkipTurns = 1
	if endTurn(t, g); g.Challenge.Timer != 6 || g.Turn.Active != starter {
		t.Fatalf("after a skipped turn: timer %d, active %d", g.Challenge.Timer, g.Turn.Active)
	}
	g.Players[starter].SkipTurns = 1
	if endTurn(t, g); g.Challenge.Timer != 6 || g.Turn.Active != partner {
		t.Fatalf("to the partner: timer %d, active %d", g.Challenge.Timer, g.Turn.Active)
	}
	// The round ends; the starter's next turn is skipped, but counts.
	if ev := endTurn(t, g); g.Challenge.Timer != 5 || g.Turn.Active != partner || countKind(ev, engine.EvChallengeTimer) != 1 {
		t.Fatalf("the starter's turn skipped: timer %d, active %d", g.Challenge.Timer, g.Turn.Active)
	}
	if endTurn(t, g); g.Challenge.Timer != 4 || g.Turn.Active != starter {
		t.Fatalf("after the round with the skipped turn: timer %d, active %d", g.Challenge.Timer, g.Turn.Active)
	}
}

// R-COOP-10, R-CHAL-15, ADR 028 Decision 6: the rule card's abilities
// work; its goal met, both players win at once.
func TestChallengeWin(t *testing.T) {
	g := newTrial(t, trialSetup(3))
	endTurn(t, g)
	g.Players[g.Turn.Active].Cents = 10
	ev := endTurn(t, g)
	if !g.Over || len(g.Winners) != 2 || g.Challenge.Lost || g.Prompt().Kind != engine.PromptGameOver {
		t.Fatalf("over %v, winners %v, challenge %+v", g.Over, g.Winners, g.Challenge)
	}
	if countKind(ev, engine.EvChallengeWon) != 1 || countKind(ev, engine.EvGameWon) != 0 {
		t.Errorf("events %+v", ev)
	}
	if err := g.CheckInvariants(); err != nil {
		t.Fatal(err)
	}
}

// R-CHAL-14: the soul target and "that player wins" effects do not end
// a challenge.
func TestChallengeIgnoresSoulTargetAndWinEffects(t *testing.T) {
	s := trialSetup(3)
	s.Characters = []engine.CardRef{"trial_saint", "trial_hero"}
	regular := s
	regular.Challenge, regular.Difficulty = nil, ""
	if g, _, err := engine.NewGame(regular); err != nil || !g.Over {
		t.Fatalf("a regular game with 4 souls goes on: %v", err)
	}
	g := newTrial(t, s)
	if g.Over || g.SoulValue(0) < g.WinSouls {
		t.Fatalf("over %v, souls %d", g.Over, g.SoulValue(0))
	}
	g.WinGame(0)
	if g.Over {
		t.Fatal("WinGame ended a challenge")
	}
}

// R-COOP-11, Q-38, Q-39: with rooms, a room whose top card was in its
// slot when the active player's previous turn started goes to discard
// after the room step, and the slot refills. Nothing goes on the first
// round.
func TestChallengeRoomRule(t *testing.T) {
	s := trialSetup(3)
	s.Rooms = true
	g := newTrial(t, s)
	first, _ := g.Rooms[0].TopOf()
	for turn := 1; turn <= 2; turn++ {
		if endTurn(t, g); g.Rooms[0].Cards[len(g.Rooms[0].Cards)-1] != first {
			t.Fatalf("turn %d: the room went in the first round", turn)
		}
	}
	ev := endTurn(t, g) // the starter's second turn
	second, _ := g.Rooms[0].TopOf()
	if second == first || g.Object(first).Zone.Kind != engine.ZoneGone || len(g.Discards[engine.RoomDeck]) != 1 {
		t.Fatalf("turn 3: the room stayed")
	}
	if !slices.ContainsFunc(ev, func(e engine.Event) bool { return e.Kind == engine.EvDiscarded && e.Prev == first }) {
		t.Error("no discarded event for the room")
	}
	for turn := 4; turn <= 5; turn++ {
		if endTurn(t, g); g.Rooms[0].Cards[len(g.Rooms[0].Cards)-1] != second {
			t.Fatalf("turn %d: the new room went before a whole round", turn)
		}
	}
	if endTurn(t, g); g.Rooms[0].Cards[len(g.Rooms[0].Cards)-1] == second {
		t.Fatal("turn 6: the room stayed a whole round")
	}
}

// R-COOP-07, Q-42, ADR 028 Decision 9: the view shows the challenge to
// everyone, and each player the partner's hand; spectators see no hand.
func TestChallengeView(t *testing.T) {
	g := newTrial(t, trialSetup(3))
	v := g.View(engine.Viewer{Kind: engine.ViewPlayer, Player: 0})
	c := v.Challenge
	if c == nil || c.ID != "trial" || c.Difficulty != engine.DifficultyNormal || c.Timer != 8 || c.Card.Card != "trial_rules" || c.Lost {
		t.Fatalf("challenge view %+v", c)
	}
	if len(v.Players[1].Hand) != len(g.Players[1].Hand) || len(v.Players[1].Hand) == 0 {
		t.Errorf("player 0 sees %d of the partner's %d cards", len(v.Players[1].Hand), len(g.Players[1].Hand))
	}
	if s := g.View(engine.Viewer{Kind: engine.ViewSpectator}); len(s.Players[0].Hand)+len(s.Players[1].Hand) != 0 {
		t.Error("a spectator sees a hand")
	}
	r, _, err := engine.NewGame(engine.Setup{Seed: 3, Players: 2, Sets: []engine.CardSet{trialSet}, Characters: s2("trial_hero")})
	if err != nil {
		t.Fatal(err)
	}
	if rv := r.View(engine.Viewer{Kind: engine.ViewPlayer, Player: 0}); rv.Challenge != nil || len(rv.Players[1].Hand) != 0 {
		t.Error("a regular game shows a challenge or the other hand")
	}
	for range 8 {
		endTurn(t, g)
	}
	if cv := g.View(engine.Viewer{Kind: engine.ViewSpectator}).Challenge; cv.Timer != 4 || cv.Lost {
		t.Errorf("after 4 rounds: %+v", cv)
	}
}

func s2(c engine.CardRef) []engine.CardRef { return []engine.CardRef{c, c} }

// ADR 019, ADR 028 Decision 5: a challenge saves, resumes to the same
// bytes, and plays on.
func TestChallengeSaveResume(t *testing.T) {
	g := newTrial(t, trialSetup(3))
	endTurn(t, g)
	endTurn(t, g)
	data, err := g.Save()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(data, []byte(`"challenge":{"id":"trial"`)) {
		t.Fatalf("no challenge in the save: %s", data)
	}
	r, err := engine.Resume(data, engine.StateMeaning, trialSet)
	if err != nil {
		t.Fatal(err)
	}
	if r.Challenge.Timer != 7 || r.Challenge.Card != g.Challenge.Card {
		t.Fatalf("resumed challenge %+v", r.Challenge)
	}
	c, err := r.Clone()
	if err != nil || c.Challenge == nil || c.Challenge == r.Challenge {
		t.Fatalf("clone: %v", err)
	}
	r.Players[r.Turn.Active].Cents = 10
	if endTurn(t, r); !r.Over || len(r.Winners) != 2 {
		t.Fatal("the resumed challenge's rule card does not work")
	}
}

// ADR 028: Replay rebuilds a challenge from its setup and steps.
func TestChallengeReplays(t *testing.T) {
	s := trialSetup(5)
	g, _, err := engine.NewGame(s)
	if err != nil {
		t.Fatal(err)
	}
	var steps []engine.StepRecord
	for i := 0; i < 200 && !g.Over; i++ {
		p := g.Prompt()
		allowed := g.Allowed(p.Player)
		in := allowed[(i*7)%len(allowed)]
		if in.Kind == engine.IntentDiscard {
			in.Objects = in.Objects[:p.Count]
		}
		st, serr := g.Step(in)
		if serr != nil {
			t.Fatal(serr)
		}
		steps = append(steps, st)
	}
	r, err := engine.Replay(s, steps)
	if err != nil {
		t.Fatal(err)
	}
	if r.Challenge == nil || r.Challenge.Timer != g.Challenge.Timer {
		t.Fatalf("replayed challenge %+v, want %+v", r.Challenge, g.Challenge)
	}
}
