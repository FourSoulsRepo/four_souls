package server

import (
	"context"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/b2"
)

// TestNoWaitingWithoutOptions: in a played game, no update ever asks a
// player who can only pass (N-05, N-09; 6.7 done when).
func TestNoWaitingWithoutOptions(t *testing.T) {
	h := startHub(t)
	players := startTable(t, h, 2)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	done := make(chan struct{})
	for i, p := range players {
		go func() {
			defer func() { done <- struct{}{} }()
			rng := rand.New(rand.NewPCG(uint64(i), 9)) //nolint:gosec // test moves
			for {
				select {
				case <-ctx.Done():
					return
				case msg := <-p.conn.out:
					e, err := protocol.Decode(msg)
					if err != nil || e.Type != protocol.TypeUpdate {
						continue
					}
					var u protocol.Update
					if e.Unpack(&u) != nil || u.View.Over || len(u.Allowed) == 0 {
						continue
					}
					w := u.View.Waiting
					if w.Kind == engine.PromptPriority && len(u.Allowed) == 1 && u.Allowed[0].Kind == engine.IntentPass {
						t.Errorf("step %d: player %d waits with only a pass", u.Step, w.Player)
						return
					}
					in := u.Allowed[rng.IntN(len(u.Allowed))]
					if in.Kind == engine.IntentDiscard {
						in.Objects = in.Objects[:w.Count]
					}
					out, err := protocol.Encode(protocol.TypeIntent, 1, protocol.Intent{Intent: in})
					if err != nil {
						t.Error(err)
						return
					}
					p.c.Receive(out)
				}
			}
		}()
	}
	<-done
	<-done
}

// TestSkipAll: the player passes on the stack they saw; a new item
// returns control; an empty stack ends it (N-06).
func TestSkipAll(t *testing.T) {
	g, err := engine.BuildTable(engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "isaac"}, {Character: "cain"}},
		Stack:   []engine.SituationStack{{DiceRoll: 3, Owner: 0}},
	}, b2.Set)
	if err != nil {
		t.Fatal(err)
	}
	r := &Room{game: g, seats: []roomSeat{{client: &Client{}}, {client: &Client{}}}} // both connected
	if _, ok := r.autoIntent(); ok {
		t.Fatal("a player who can act was skipped") // Isaac's character is charged
	}
	r.seats[0].skip = []int{g.Stack[0].ID}
	if in, ok := r.autoIntent(); !ok || in.Kind != engine.IntentPass || in.Player != 0 {
		t.Fatalf("skip all did not pass: %+v %v", in, ok)
	}
	r.seats[0].skip = []int{999} // the stack holds an item the player has not seen
	if _, ok := r.autoIntent(); ok || r.seats[0].skip != nil {
		t.Error("a new stack item did not return control")
	}
}

// TestWaitTurn: a waiting player is passed for until the player before
// them ends the turn; then they are asked again (N-16).
func TestWaitTurn(t *testing.T) {
	g, err := engine.BuildTable(engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "isaac"}, {Character: "cain"}, {Character: "maggy"}},
		Stack:   []engine.SituationStack{{DiceRoll: 3, Owner: 0}},
	}, b2.Set)
	if err != nil {
		t.Fatal(err)
	}
	r := &Room{game: g, seats: []roomSeat{{client: &Client{}}, {client: &Client{}}, {client: &Client{}}}}
	if _, err := g.Apply(engine.Intent{Player: 0, Kind: engine.IntentPass}); err != nil {
		t.Fatal(err)
	}
	if w := g.Prompt(); w.Player != 1 {
		t.Fatalf("player %d has priority, want 1", w.Player)
	}
	if _, ok := r.autoIntent(); ok {
		t.Fatal("a player who can act was skipped") // Cain's character is charged
	}
	r.seats[0].wait, r.seats[1].wait, r.seats[2].wait = true, true, true
	if in, ok := r.autoIntent(); !ok || in.Kind != engine.IntentPass || in.Player != 1 {
		t.Fatalf("wait for my turn did not pass: %+v %v", in, ok)
	}
	if r.seats[0].wait {
		t.Error("the active player still waits for their turn")
	}
	g.Turn.EndDeclared = true // player 0 ends the turn: player 1 is next
	if _, ok := r.autoIntent(); ok || r.seats[1].wait {
		t.Error("the end of the turn before did not return control")
	}
	if !r.seats[2].wait {
		t.Error("player 2 stopped waiting at the end of player 0's turn")
	}
}

// TestResponseTimer: with nobody answering, the timer moves the game on
// (N-07).
func TestResponseTimer(t *testing.T) {
	h := startHubWith(t, time.Millisecond)
	players := setupTable(t, h, 2, protocol.Options{ResponseTimer: 60})
	var u protocol.Update
	players[0].next(protocol.TypeUpdate, &u)
	if u.Deadline == 0 {
		t.Error("no deadline shown")
	}
	first := u.Step
	deadline := time.Now().Add(5 * time.Second)
	for u.Step < first+5 && time.Now().Before(deadline) {
		players[0].next(protocol.TypeUpdate, &u)
	}
	if u.Step < first+5 {
		t.Errorf("the game did not move on: step %d", u.Step)
	}
}

// atThrow is a game where Isaac (seat 0) attacks Pin and must throw
// the attack roll (15.59).
func atThrow(t *testing.T) *engine.Game {
	t.Helper()
	g, err := engine.BuildTable(engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac"}, {Character: "cain"}},
		Monsters: []engine.CardRef{"pin"},
	}, b2.Set)
	if err != nil {
		t.Fatal(err)
	}
	for _, in := range []engine.Intent{
		{Player: 0, Kind: engine.IntentAttack},
		{Player: 0, Kind: engine.IntentPass},
		{Player: 1, Kind: engine.IntentPass},
		{Player: 0, Kind: engine.IntentChoose}, // Pin
	} {
		if _, err := g.Apply(in); err != nil {
			t.Fatalf("%+v: %v", in, err)
		}
	}
	if w := g.Prompt(); w.Kind != engine.PromptThrow || w.Player != 0 {
		t.Fatalf("prompt %+v, want Isaac's throw", w)
	}
	return g
}

// thrown: the throw prompt is gone and the attack roll has its value.
func thrown(t *testing.T, g *engine.Game) {
	t.Helper()
	if w := g.Prompt(); w.Kind == engine.PromptThrow {
		t.Fatalf("still waiting for the throw: %+v", w)
	}
	if len(g.Stack) == 0 || g.Stack[len(g.Stack)-1].Roll == 0 {
		t.Fatalf("stack %+v: the roll has no value", g.Stack)
	}
}

// TestKickedSeatThrows: a kicked seat throws at once (15.60).
func TestKickedSeatThrows(t *testing.T) {
	g := atThrow(t)
	r := &Room{game: g, seats: []roomSeat{{kicked: true}, {client: &Client{}}}}
	in, ok := r.autoIntent()
	if !ok || in.Kind != engine.IntentThrow || in.Player != 0 {
		t.Fatalf("auto intent %+v %v, want Isaac's throw", in, ok)
	}
	r.after(nil, false)
	thrown(t, g)
}

// TestTimerThrows: the response timer throws for an idle player (15.60).
func TestTimerThrows(t *testing.T) {
	g := atThrow(t)
	r := &Room{game: g, seats: []roomSeat{{client: &Client{}}, {client: &Client{}}}}
	r.timeUp(r.timerGen)
	thrown(t, g)
}

// TestSkipAndWaitNeverThrow: "skip all" and "wait for my turn" pass,
// but never throw for the player (15.60).
func TestSkipAndWaitNeverThrow(t *testing.T) {
	g := atThrow(t)
	r := &Room{game: g, seats: []roomSeat{{client: &Client{}, wait: true}, {client: &Client{}, wait: true}}}
	r.seats[0].skip = []int{g.Stack[len(g.Stack)-1].ID}
	if in, ok := r.autoIntent(); ok {
		t.Fatalf("the server answered %+v for the roller", in)
	}
	if g.Prompt().Kind != engine.PromptThrow {
		t.Error("the throw prompt is gone")
	}
}
