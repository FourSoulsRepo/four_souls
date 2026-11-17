package world

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"sync/atomic"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards"
	"github.com/FourSoulsRepo/rules_engine/cards/b2"
	"github.com/FourSoulsRepo/rules_engine/cards/challenge"
)

// memoryGames and memorySteps size test 6 of 23.15: games, and steps
// each game plays at most.
const (
	memoryGames = 100
	memorySteps = 500
)

// memorySetups are games of every set: b2 alone and with each other
// set, all sets at 2 to 4 players with rooms and deck ratios, and each
// challenge at each difficulty.
func memorySetups() []engine.Setup {
	var out []engine.Setup
	all := cards.Sets()
	for i, s := range all {
		sets := []engine.CardSet{all[0]}
		if i > 0 {
			sets = append(sets, s)
		}
		out = append(out, engine.Setup{Players: 2 + i%3, Sets: sets, BonusSouls: true, Rooms: true})
	}
	for n := 2; n <= 4; n++ {
		out = append(out, engine.Setup{Players: n, Sets: all, BonusSouls: true, Rooms: n != 3, DeckRatios: n == 3})
	}
	for _, c := range cards.Challenges() {
		for _, d := range []string{engine.DifficultyNormal, engine.DifficultyHard, engine.DifficultyUltra} {
			out = append(out, engine.Setup{Players: 2, Sets: append(cards.Sets(), challenge.Set), Challenge: c, Difficulty: d})
		}
	}
	return out
}

// pickIntent is a random allowed intent, two times in three an action
// over a pass or an end of turn, a discard cut to its count.
func pickIntent(g *engine.Game, rng *engine.RNG) (engine.Intent, bool) {
	w := g.Prompt()
	allowed := g.Allowed(w.Player)
	if len(allowed) == 0 {
		return engine.Intent{}, false
	}
	var acts []engine.Intent
	for _, in := range allowed {
		if in.Kind != engine.IntentPass && in.Kind != engine.IntentEndTurn {
			acts = append(acts, in)
		}
	}
	in := allowed[rng.Intn(len(allowed))]
	if len(acts) > 0 && rng.Intn(3) > 0 {
		in = acts[rng.Intn(len(acts))]
	}
	if in.Kind == engine.IntentDiscard {
		in.Objects = in.Objects[:g.PromptFor(w.Player).Count]
	}
	return in, true
}

// feed gives memory m the step as its seat's client gets it.
func feed(m *Memory, g *engine.Game, events []engine.Event) {
	v := engine.Viewer{Kind: engine.ViewPlayer, Player: m.seat}
	m.Observe(g.View(v), engine.FilterEvents(events, v))
}

// wrongMemory checks every remembered card against the real game: a
// hand card must be in that hand (as many copies as remembered), a deck
// card at its place. It returns the cards checked.
func wrongMemory(g *engine.Game, k engine.Known) (int, error) {
	return wrongMemorySeen(g, k, nil)
}

// wrongMemorySeen is wrongMemory that counts no card of a hand the seat
// sees now (seen, by seat): only cards out of sight are memory.
func wrongMemorySeen(g *engine.Game, k engine.Known, seen []bool) (int, error) {
	checked := 0
	for q, refs := range k.Hands {
		if q < len(seen) && seen[q] {
			checked -= len(refs)
		}
		var hand []engine.CardRef
		for _, id := range g.Players[q].Hand {
			hand = append(hand, g.Object(id).Card)
		}
		for _, ref := range refs {
			i := slices.Index(hand, ref)
			if i < 0 {
				return checked, fmt.Errorf("hand %d: %s is not there (hand %v, remembered %v)", q, ref, hand, refs)
			}
			hand = slices.Delete(hand, i, i+1)
			checked++
		}
	}
	for d := range k.Tops {
		deck := g.Decks[d]
		at := func(pos int, ref engine.CardRef, where string) error {
			if ref == "" {
				return nil
			}
			checked++
			if pos < 0 || pos >= len(deck) {
				return fmt.Errorf("%s deck: %s %s is outside a deck of %d", engine.DeckKind(d), where, ref, len(deck))
			}
			if got := g.Object(deck[pos]).Card; got != ref {
				return fmt.Errorf("%s deck: %s %s, the card there is %s", engine.DeckKind(d), where, ref, got)
			}
			return nil
		}
		for i, ref := range k.Tops[d] {
			if err := at(len(deck)-1-i, ref, fmt.Sprintf("top %d", i)); err != nil {
				return checked, err
			}
		}
		for i, ref := range k.Bottoms[d] {
			if err := at(i, ref, fmt.Sprintf("bottom %d", i)); err != nil {
				return checked, err
			}
		}
	}
	return checked, nil
}

// TestMemoryIsTrue is test 6 of 23.15: in simulated games of every set,
// every card each seat's memory holds is truly where it says, at every
// step (Decision 19.9).
func TestMemoryIsTrue(t *testing.T) {
	memoryIsTrue(t, memoryGames, memorySteps)
}

// stressGames and stressSteps size the long memory run (23.15): eight
// times the games of TestMemoryIsTrue, four times as long, so late-game
// moves and rare cards come up.
const (
	stressGames = 800
	stressSteps = 2000
)

// TestMemoryStress is TestMemoryIsTrue at 800 games of 2000 steps, a
// few minutes on all cores. It runs only with FOUR_SOULS_MEMORY_STRESS
// set and never with -short: make sims runs it after the simulations
// (docs/wiki/Fuzzing.md). Plain go test checks the same with
// TestMemoryIsTrue.
func TestMemoryStress(t *testing.T) {
	if testing.Short() || os.Getenv("FOUR_SOULS_MEMORY_STRESS") == "" {
		t.Skip("long run: FOUR_SOULS_MEMORY_STRESS=1, or make sims")
	}
	memoryIsTrue(t, stressGames, stressSteps)
}

// memoryIsTrue plays games simulated games of at most maxSteps steps
// each and checks every seat's memory against the real game at every
// step.
func memoryIsTrue(t *testing.T, games uint64, maxSteps int) {
	setups := memorySetups()
	var steps, checked, remembering atomic.Int64
	t.Run("games", func(t *testing.T) {
		for k := range games {
			s := setups[k%uint64(len(setups))]
			s.Seed = 100 + k
			t.Run(fmt.Sprint(k), func(t *testing.T) {
				t.Parallel()
				g, setup, err := engine.NewGame(s)
				if err != nil {
					t.Fatal(err)
				}
				mems := make([]*Memory, s.Players)
				for p := range mems {
					mems[p] = NewMemory(engine.PlayerID(p))
					feed(mems[p], g, setup)
				}
				rng := engine.NewRNG(s.Seed)
				for step := 0; step < maxSteps && !g.Over; step++ {
					in, ok := pickIntent(g, &rng)
					if !ok {
						t.Fatalf("seed %d step %d: nothing allowed", s.Seed, step)
					}
					events, err := g.Apply(in)
					if err != nil {
						t.Fatalf("seed %d step %d: %+v refused: %v", s.Seed, step, in, err)
					}
					steps.Add(1)
					for p, m := range mems {
						v := engine.Viewer{Kind: engine.ViewPlayer, Player: engine.PlayerID(p)}
						view := g.View(v)
						m.Observe(view, engine.FilterEvents(events, v))
						seen := make([]bool, len(view.Players))
						for q, pl := range view.Players {
							seen[q] = pl.Hand != nil
						}
						n, err := wrongMemorySeen(g, m.Known(), seen)
						if err != nil {
							t.Fatalf("seed %d step %d seat %d: wrong memory: %v", s.Seed, step, p, err)
						}
						checked.Add(int64(n))
						if n > 0 {
							remembering.Add(1)
						}
					}
				}
			})
		}
	})
	t.Logf("%d games, %d steps, %d seat-steps remembering a card out of sight, %d such cards checked, 0 wrong",
		games, steps.Load(), remembering.Load(), checked.Load())
	if checked.Load() == 0 {
		t.Error("no seat ever remembered a card: the test checks nothing")
	}
}

// addStrings adds every string in v's JSON to seen.
func addStrings(t *testing.T, seen map[string]bool, v any) {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var walk func(x any)
	walk = func(x any) {
		switch x := x.(type) {
		case string:
			seen[x] = true
		case []any:
			for _, y := range x {
				walk(y)
			}
		case map[string]any:
			for _, y := range x {
				walk(y)
			}
		}
	}
	var x any
	if err := json.Unmarshal(b, &x); err != nil {
		t.Fatal(err)
	}
	walk(x)
}

// TestMemoryNamesOnlyShownCards is test 7 of 23.15: a memory built from
// seat A's feed never names a card that seat A's views and events did
// not show (Decision 19.4).
func TestMemoryNamesOnlyShownCards(t *testing.T) {
	setups := memorySetups()
	for k := range uint64(12) {
		s := setups[k%uint64(len(setups))]
		s.Seed = 300 + k
		g, _, err := engine.NewGame(s)
		if err != nil {
			t.Fatal(err)
		}
		m := NewMemory(0)
		v := engine.Viewer{Kind: engine.ViewPlayer, Player: 0}
		shown := map[string]bool{}
		observe := func(events []engine.Event) {
			view, filtered := g.View(v), engine.FilterEvents(events, v)
			addStrings(t, shown, view)
			addStrings(t, shown, filtered)
			m.Observe(view, filtered)
		}
		observe(nil)
		rng := engine.NewRNG(s.Seed)
		for step := 0; step < memorySteps && !g.Over; step++ {
			in, ok := pickIntent(g, &rng)
			if !ok {
				t.Fatalf("seed %d step %d: nothing allowed", s.Seed, step)
			}
			events, err := g.Apply(in)
			if err != nil {
				t.Fatal(err)
			}
			observe(events)
			known := m.Known()
			var named []engine.CardRef
			for _, h := range known.Hands {
				named = append(named, h...)
			}
			for d := range known.Tops {
				named = append(named, known.Tops[d]...)
				named = append(named, known.Bottoms[d]...)
			}
			for _, ref := range named {
				if ref != "" && !shown[string(ref)] {
					t.Fatalf("seed %d step %d: the memory names %s, which seat 0 was never shown", s.Seed, step, ref)
				}
			}
		}
	}
}

// table is a b2 game of the given seats, seat 0 active, with a memory
// for seat 0 that saw the start.
func table(t *testing.T, seats ...engine.SituationPlayer) (*engine.Game, *Memory) {
	t.Helper()
	g, err := engine.BuildTable(engine.SituationSetup{Players: seats}, b2.Set)
	if err != nil {
		t.Fatal(err)
	}
	m := NewMemory(0)
	feed(m, g, nil)
	return g, m
}

// top is the card on top of deck d.
func top(g *engine.Game, d engine.DeckKind) engine.CardRef {
	return g.Object(g.DeckTop(d, 1)[0]).Card
}

// toTurnOf ends turns until it is seat p's action phase, feeding m each
// step: p's turn starts with its loot draw.
func toTurnOf(t *testing.T, g *engine.Game, m *Memory, p engine.PlayerID) {
	t.Helper()
	for range 200 {
		if g.Turn.Active == p && g.Turn.Step == engine.StepAction && g.Prompt().Kind == engine.PromptPriority && len(g.Stack) == 0 {
			return
		}
		w := g.Prompt()
		in := engine.Intent{Player: w.Player, Kind: engine.IntentPass}
		switch {
		case w.Kind == engine.PromptThrow:
			in.Kind = engine.IntentThrow
		case w.Kind == engine.PromptPriority && w.Player == g.Turn.Active && g.Turn.Step == engine.StepAction && len(g.Stack) == 0:
			in.Kind = engine.IntentEndTurn
		case w.Kind != engine.PromptPriority:
			in = g.Allowed(w.Player)[0]
			if in.Kind == engine.IntentDiscard {
				in.Objects = in.Objects[:g.PromptFor(w.Player).Count]
			}
		}
		events, err := g.Apply(in)
		if err != nil {
			t.Fatalf("%+v refused: %v", in, err)
		}
		feed(m, g, events)
		if _, err := wrongMemory(g, m.Known()); err != nil {
			t.Fatal(err)
		}
	}
	t.Fatalf("no action phase of seat %d", p)
}

// TestMemoryLookedAtTop is test 9 of 23.15: the seat remembers a deck
// top it looked at, and where it goes when another seat draws it; a
// world deals it there.
func TestMemoryLookedAtTop(t *testing.T) {
	g, m := table(t, engine.SituationPlayer{Character: "isaac"}, engine.SituationPlayer{Character: "cain"})
	card := top(g, engine.LootDeck)
	g.LookAt(0, g.DeckTop(engine.LootDeck, 1)...)
	feed(m, g, g.Run())
	if got := m.Known().Tops[engine.LootDeck]; !slices.Equal(got, []engine.CardRef{card}) {
		t.Fatalf("loot deck tops %v, want [%s]", got, card)
	}
	worlds, err := For(g, 0, m.Known(), 5)
	if err != nil {
		t.Fatal(err)
	}
	w := world(t, worlds, 0)
	if got := top(w, engine.LootDeck); got != card {
		t.Errorf("the world's loot deck top is %s, want the remembered %s", got, card)
	}
	// Seat 1's turn starts with a loot draw: the card goes to its hand.
	toTurnOf(t, g, m, 1)
	k := m.Known()
	if !slices.Contains(k.Hands[1], card) {
		t.Errorf("hand 1 remembers %v, want %s drawn from the top", k.Hands[1], card)
	}
	if len(k.Tops[engine.LootDeck]) != 0 {
		t.Errorf("loot deck tops %v after the draw, want none", k.Tops[engine.LootDeck])
	}
}

// TestMemoryShownHandCard is test 9 of 23.15: the seat remembers a card
// it was shown from another hand, until that card leaves it; a card
// given unseen makes it forget that hand.
func TestMemoryShownHandCard(t *testing.T) {
	g, m := table(t,
		engine.SituationPlayer{Character: "isaac"},
		engine.SituationPlayer{Character: "cain", Hand: []engine.CardRef{"a_penny", "a_nickel"}},
		engine.SituationPlayer{Character: "maggy"})
	nickel := g.Players[1].Hand[1]
	g.LookAt(0, nickel)
	feed(m, g, g.Run())
	if got := m.Known().Hands[1]; !slices.Equal(got, []engine.CardRef{"a_nickel"}) {
		t.Fatalf("hand 1 remembers %v, want [a_nickel]", got)
	}
	// Looked at again: still one.
	g.LookAt(0, nickel)
	feed(m, g, g.Run())
	if got := m.Known().Hands[1]; len(got) != 1 {
		t.Fatalf("hand 1 remembers %v after a second look, want one card", got)
	}
	// The other card goes to seat 2 unseen: seat 0 cannot tell which.
	g.GiveHandCard(1, 2, g.Players[1].Hand[0])
	feed(m, g, g.Run())
	if got := m.Known().Hands[1]; len(got) != 0 {
		t.Errorf("hand 1 remembers %v after a card left unseen, want none", got)
	}
	// Shown again, then discarded in sight: gone.
	g.LookAt(0, g.Players[1].Hand[0])
	feed(m, g, g.Run())
	if got := m.Known().Hands[1]; !slices.Equal(got, []engine.CardRef{"a_nickel"}) {
		t.Fatalf("hand 1 remembers %v, want [a_nickel]", got)
	}
	g.DiscardFromHand(1, g.Players[1].Hand[0])
	feed(m, g, g.Run())
	if got := m.Known().Hands[1]; len(got) != 0 {
		t.Errorf("hand 1 remembers %v after the discard, want none", got)
	}
}

// TestMemoryDiscardTakenBack is test 9 of 23.15: a card taken from a
// discard pile into a hand is remembered there.
func TestMemoryDiscardTakenBack(t *testing.T) {
	g, m := table(t,
		engine.SituationPlayer{Character: "isaac"},
		engine.SituationPlayer{Character: "cain", Hand: []engine.CardRef{"a_dime"}})
	g.DiscardFromHand(1, g.Players[1].Hand[0])
	feed(m, g, g.Run())
	pile := g.Discards[engine.LootDeck]
	g.DiscardToHand(engine.LootDeck, pile[len(pile)-1], 1)
	feed(m, g, g.Run())
	if got := m.Known().Hands[1]; !slices.Equal(got, []engine.CardRef{"a_dime"}) {
		t.Errorf("hand 1 remembers %v, want [a_dime]", got)
	}
	if _, err := wrongMemory(g, m.Known()); err != nil {
		t.Error(err)
	}
}

// TestMemoryForgetsShuffledDeck is test 8 of 23.15: a shuffle, a reorder
// of the top cards or a card moved within the deck empties that deck's
// memory; the other decks keep theirs.
func TestMemoryForgetsShuffledDeck(t *testing.T) {
	for name, change := range map[string]func(g *engine.Game){
		"shuffle": func(g *engine.Game) { g.ShuffleDeck(engine.LootDeck) },
		"reorder": func(g *engine.Game) {
			top := g.DeckTop(engine.LootDeck, 2)
			g.SetDeckTop(engine.LootDeck, []engine.ObjectID{top[1], top[0]})
		},
		"to the bottom": func(g *engine.Game) { g.DeckToBottom(engine.LootDeck, g.DeckTop(engine.LootDeck, 1)[0]) },
	} {
		t.Run(name, func(t *testing.T) {
			g, m := table(t, engine.SituationPlayer{Character: "isaac"}, engine.SituationPlayer{Character: "cain"})
			g.LookAt(0, g.DeckTop(engine.LootDeck, 1)...)
			g.LookAt(0, g.DeckTop(engine.TreasureDeck, 1)...)
			feed(m, g, g.Run())
			if len(m.Known().Tops[engine.LootDeck]) != 1 {
				t.Fatal("the look was not remembered")
			}
			change(g)
			feed(m, g, g.Run())
			k := m.Known()
			if len(k.Tops[engine.LootDeck]) != 0 || len(k.Bottoms[engine.LootDeck]) != 0 {
				t.Errorf("loot deck remembers %v and %v, want nothing", k.Tops[engine.LootDeck], k.Bottoms[engine.LootDeck])
			}
			if len(k.Tops[engine.TreasureDeck]) != 1 {
				t.Errorf("treasure deck remembers %v, want its top", k.Tops[engine.TreasureDeck])
			}
		})
	}
}

// TestMemoryReset: Reset forgets everything; Known of a memory that saw
// nothing is empty.
func TestMemoryReset(t *testing.T) {
	g, m := table(t, engine.SituationPlayer{Character: "isaac"}, engine.SituationPlayer{Character: "cain"})
	g.LookAt(0, g.DeckTop(engine.LootDeck, 1)...)
	feed(m, g, g.Run())
	m.Reset()
	k := m.Known()
	if k.Hands != nil || len(k.Tops[engine.LootDeck]) != 0 {
		t.Errorf("after Reset: %+v", k)
	}
}

// TestShuffledDeckImaginedAgain: main's shuffled event (19.178) makes the
// memory drop the deck's known top, so For no longer deals the
// remembered card there: across worlds, the imagined top varies again,
// while before the shuffle every world put the remembered card on top.
func TestShuffledDeckImaginedAgain(t *testing.T) {
	g, m := table(t, engine.SituationPlayer{Character: "isaac"}, engine.SituationPlayer{Character: "cain"})
	g.LookAt(0, g.DeckTop(engine.LootDeck, 1)...)
	feed(m, g, g.Run())
	tops := func() map[engine.CardRef]bool {
		t.Helper()
		worlds, err := For(g, 0, m.Known(), 7)
		if err != nil {
			t.Fatal(err)
		}
		seen := map[engine.CardRef]bool{}
		for i := range 20 {
			w, err := worlds(i)
			if err != nil {
				t.Fatal(err)
			}
			seen[top(w, engine.LootDeck)] = true
		}
		return seen
	}
	if seen := tops(); len(seen) != 1 || !seen[top(g, engine.LootDeck)] {
		t.Fatalf("before the shuffle, worlds put %v on top; want only the remembered card", seen)
	}
	g.ShuffleDeck(engine.LootDeck)
	events := g.Run()
	if !slices.ContainsFunc(events, func(e engine.Event) bool { return e.Kind == engine.EvShuffled }) {
		t.Fatal("ShuffleDeck sent no shuffled event")
	}
	feed(m, g, events)
	if k := m.Known(); len(k.Tops[engine.LootDeck]) != 0 {
		t.Fatalf("after shuffled, the memory keeps the loot top %v", k.Tops[engine.LootDeck])
	}
	if seen := tops(); len(seen) < 2 {
		t.Errorf("after the shuffle, every world put %v on top; want the deck imagined shuffled", seen)
	}
}
