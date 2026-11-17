package server

import (
	"reflect"
	"testing"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	"github.com/FourSoulsRepo/four_souls/internal/world"
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards"
)

// TestBotsMemoryLikeArena: a hard seat's memory in the room is the
// memory the arena builds on the same game (ADR 030, Decisions 19.13
// and 19.17). The room plays a game of a player, two hard bots and a
// standard one on every set; after each request a reference game
// replays the room's new steps one by one, from the start, and feeds
// arena-style memories (NewMemory, Observe after every step with the
// seat's ViewFor and FilterEvents). Room and reference memories must be
// equal, unexported state and all, its count of observed steps among
// it: any step the room's memory skipped shows, the server's own passes
// and forced answers among them.
// Standard keeps no memory (Decision 19.1).
func TestBotsMemoryLikeArena(t *testing.T) {
	seeds := []uint64{3, 11, 29}
	if raceEnabled || testing.Short() {
		seeds = seeds[:1]
	}
	known := 0
	for _, seed := range seeds {
		known += playMemoryGame(t, seed, 400)
	}
	// Bots think on a time budget, so a seed's game differs from run to
	// run; some games remember nothing. The seeds together must.
	if known == 0 {
		t.Error("the memories never held a card; the check proves nothing")
	}
}

// playMemoryGame plays up to steps steps of one seeded room game and
// checks its memories against the reference after every request. It
// returns the remembered cards summed over the checks.
func playMemoryGame(t *testing.T, seed uint64, steps int) int {
	t.Helper()
	seats := memorySeats
	r := memoryRoom(t, seed, seats)
	for i, m := range r.memory {
		if (m != nil) != (seats[i] == protocol.BotHard) {
			t.Fatalf("seed %d: seat %d (%q) has memory %v", seed, i, seats[i], m != nil)
		}
	}
	ref, err := engine.Load(r.base, r.setup.Sets...)
	if err != nil {
		t.Fatal(err)
	}
	mems := make([]*world.Memory, len(seats))
	for i, level := range seats {
		if level == protocol.BotHard {
			mems[i] = world.NewMemory(engine.PlayerID(i))
		}
	}
	feedArena(ref, mems, nil) // the first Observe takes the view as it is
	done, known := 0, 0
	check := func() {
		t.Helper()
		for ; done < len(r.live); done++ {
			st, err := ref.Step(r.live[done].rec.Intent)
			if err != nil {
				t.Fatalf("seed %d: the reference refused step %d: %v", seed, done, err)
			}
			if st.Checksum != r.live[done].rec.Checksum {
				t.Fatalf("seed %d: step %d differs from the room's", seed, done)
			}
			feedArena(ref, mems, st.Events)
		}
		for i, m := range mems {
			if m == nil {
				continue
			}
			if !reflect.DeepEqual(m, r.memory[i]) {
				t.Fatalf("seed %d, step %d: seat %d's memory in the room differs from the arena's:\nroom  %+v\narena %+v", seed, done, i, r.memory[i].Known(), m.Known())
			}
			k := m.Known()
			for _, h := range k.Hands {
				known += len(h)
			}
			for d := range k.Tops {
				known += len(k.Tops[d]) + len(k.Bottoms[d])
			}
		}
	}
	check()
	deadline := time.Now().Add(time.Minute)
	for wake(t, r); !r.game.Over && done < steps; wake(t, r) {
		if time.Now().After(deadline) {
			t.Fatalf("seed %d: only %d steps in time", seed, done)
		}
		if !r.playerAnswers() {
			select {
			case req := <-r.in:
				r.handle(req)
			case <-time.After(5 * time.Second):
				t.Fatalf("seed %d: no bot move came", seed)
			}
		}
		check()
	}
	t.Logf("seed %d: %d steps, %d remembered cards summed over the checks", seed, done, known)
	return known
}

// memorySeats are the seats of the memory tests: a player, two hard
// bots and a standard one.
var memorySeats = []string{"", protocol.BotHard, protocol.BotStandard, protocol.BotHard}

// memoryRoom starts a seeded game of the given seats on every set, with
// rollback on, so the room keeps its base and live steps. Nothing runs
// the room: the test hands it its requests.
func memoryRoom(t *testing.T, seed uint64, seats []string) *Room {
	t.Helper()
	var codes []string
	for _, s := range cards.Sets() {
		codes = append(codes, s.Code)
	}
	opts := protocol.Options{Rollback: protocol.RollbackGod}
	if err := checkOptions(&opts, codes, len(seats)); err != nil {
		t.Fatal(err)
	}
	rs := make([]roomSeat, len(seats))
	for i, level := range seats {
		if level != "" {
			rs[i] = roomSeat{name: botName(i), bot: level, botN: i}
			continue
		}
		rs[i] = roomSeat{name: "P0", token: "a", client: &Client{conn: newMemConn()}}
	}
	r, err := NewRoom(engineSetup(seed, len(seats), codes, opts), opts, rs)
	if err != nil {
		t.Fatal(err)
	}
	r.second, r.botThink = time.Millisecond, 2*time.Millisecond
	r.recInfo.game = "g1"
	t.Cleanup(r.stopThinkers)
	r.play(characters(setsOf(codes))[:len(seats)])
	if r.base == nil {
		t.Fatal("the room keeps no base to replay")
	}
	return r
}

// TestBotsMemoryReset: a rollback starts each hard seat's memory over
// from the restored game, as after a load or a restart (ADR 030,
// Decision 19.12): its undone steps stay in no memory.
func TestBotsMemoryReset(t *testing.T) {
	r := memoryRoom(t, 3, memorySeats)
	deadline := time.Now().Add(time.Minute)
	for wake(t, r); len(r.live) < 60; wake(t, r) {
		if time.Now().After(deadline) || r.game.Over {
			t.Fatalf("only %d steps", len(r.live))
		}
		if !r.playerAnswers() {
			r.handle(<-r.in)
		}
	}
	if err := r.rollBackTo(0, 10); err != nil {
		t.Fatal(err)
	}
	r.stopThinkers()
	// The reference: the restored game, new memories, then the passes
	// the server made after the rollback.
	ref, err := engine.Load(r.base, r.setup.Sets...)
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range r.live[:10] {
		if _, err := ref.Step(l.rec.Intent); err != nil {
			t.Fatal(err)
		}
	}
	mems := make([]*world.Memory, len(memorySeats))
	for i, level := range memorySeats {
		if level == protocol.BotHard {
			mems[i] = world.NewMemory(engine.PlayerID(i))
		}
	}
	feedArena(ref, mems, nil)
	for _, l := range r.live[10:] {
		st, err := ref.Step(l.rec.Intent)
		if err != nil {
			t.Fatal(err)
		}
		feedArena(ref, mems, st.Events)
	}
	for i, m := range mems {
		if m != nil && !reflect.DeepEqual(m, r.memory[i]) {
			t.Errorf("seat %d: after the rollback the memory did not start over: %+v", i, r.memory[i].Known())
		}
	}
}

// feedArena feeds the memories one step the way the arena does
// (internal/arena, observe).
func feedArena(g *engine.Game, mems []*world.Memory, events []engine.Event) {
	base := g.ViewBase()
	for i, m := range mems {
		if m != nil {
			v := engine.Viewer{Kind: engine.ViewPlayer, Player: engine.PlayerID(i)}
			m.Observe(g.ViewFor(base, v), engine.FilterEvents(events, v))
		}
	}
}
