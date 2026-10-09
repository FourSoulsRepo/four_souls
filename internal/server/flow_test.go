package server

import (
	"context"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards"
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
	}, cards.Sets()...)
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
