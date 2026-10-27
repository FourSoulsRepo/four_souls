package server

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
)

// fakeClock is a clock the test moves (N-20).
type fakeClock struct {
	mu sync.Mutex
	t  time.Time
}

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.t
}

func (c *fakeClock) Add(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.t = c.t.Add(d)
}

// startLimitHub runs a hub with the default limits on a fake clock that
// is checked every few milliseconds.
func startLimitHub(t *testing.T) (*Hub, *fakeClock) {
	t.Helper()
	clock := &fakeClock{t: time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)}
	h := NewHub()
	h.now, h.limitTick = clock.Now, 5*time.Millisecond
	h.emptyLimit, h.gameLimit = 30*time.Minute, 24*time.Hour
	ctx, cancel := context.WithCancel(context.Background())
	go h.Run(ctx)
	t.Cleanup(cancel)
	return h, clock
}

// seatedTwo starts a game of two and returns the players and tokens.
func seatedTwo(t *testing.T, h *Hub) ([]*player, []string) {
	t.Helper()
	ann, wa := arrive(t, h, "Ann", "")
	ann.say(protocol.TypeCreate, protocol.Create{Seats: 2})
	var tb protocol.Table
	ann.next(protocol.TypeTable, &tb)
	bo, wb := arrive(t, h, "Bo", "")
	bo.say(protocol.TypeJoin, protocol.Join{Game: tb.Game})
	bo.next(protocol.TypeTable, nil)
	players := []*player{ann, bo}
	for _, p := range players {
		p.say(protocol.TypeReady, protocol.Ready{Ready: true})
	}
	for _, p := range players {
		p.next(protocol.TypeUpdate, nil)
	}
	return players, []string{wa.Token, wb.Token}
}

// gamesLeft lists the games a lobby client sees.
func gamesLeft(t *testing.T, h *Hub) int {
	t.Helper()
	p, _ := arrive(t, h, "Watcher", "")
	p.say(protocol.TypeList, struct{}{})
	var g protocol.Games
	p.next(protocol.TypeGames, &g)
	p.c.Leave()
	return len(g.Games)
}

// untilNoGames waits until the lobby lists no game.
func untilNoGames(t *testing.T, h *Hub) {
	t.Helper()
	for range 200 {
		if gamesLeft(t, h) == 0 {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatal("the game is still listed")
}

// TestEmptyGameEnds: a game nobody is connected to ends after 30
// minutes; a player who comes back hears why, once (8.25).
func TestEmptyGameEnds(t *testing.T) {
	h, clock := startLimitHub(t)
	players, tokens := seatedTwo(t, h)
	for _, p := range players {
		p.c.Leave()
	}
	time.Sleep(50 * time.Millisecond) // a check sees the game empty
	clock.Add(29 * time.Minute)
	time.Sleep(50 * time.Millisecond)
	if gamesLeft(t, h) != 1 {
		t.Fatal("the game ended before 30 minutes")
	}
	clock.Add(2 * time.Minute)
	untilNoGames(t, h)

	back, w := arrive(t, h, "Ann", tokens[0])
	if w.Seat != -1 {
		t.Errorf("seat %d, want the lobby", w.Seat)
	}
	var n protocol.Notice
	back.next(protocol.TypeNotice, &n)
	if n.Kind != protocol.NoticeIdle || n.Reason != protocol.IdleEmpty {
		t.Errorf("notice %+v, want idle, empty", n)
	}
}

// TestLongGameEnds: a game ends 24 hours after it started, even with
// its players connected; they hear why (8.25).
func TestLongGameEnds(t *testing.T) {
	h, clock := startLimitHub(t)
	players, _ := seatedTwo(t, h)
	clock.Add(23 * time.Hour)
	time.Sleep(50 * time.Millisecond)
	if gamesLeft(t, h) != 1 {
		t.Fatal("the game ended before 24 hours")
	}
	clock.Add(2 * time.Hour)
	for _, p := range players {
		var n protocol.Notice
		p.next(protocol.TypeNotice, &n)
		if n.Kind != protocol.NoticeIdle || n.Reason != protocol.IdleLong {
			t.Errorf("notice %+v, want idle, long", n)
		}
		p.next(protocol.TypeGames, nil) // back in the lobby
	}
	untilNoGames(t, h)
}

// TestNoLimits: 0 turns both limits off (8.25).
func TestNoLimits(t *testing.T) {
	h, clock := startLimitHub(t)
	h.send(hubRequest{run: func(h *Hub) { h.emptyLimit, h.gameLimit = 0, 0 }})
	players, _ := seatedTwo(t, h)
	players[0].c.Leave()
	players[1].c.Leave()
	time.Sleep(50 * time.Millisecond)
	clock.Add(48 * time.Hour)
	time.Sleep(50 * time.Millisecond)
	if gamesLeft(t, h) != 1 {
		t.Error("a game ended with the limits off")
	}
}

// TestEmptyKeptTableCloses: a table kept over a restart that nobody
// comes back to closes after 30 minutes, downtime included (8.25).
func TestEmptyKeptTableCloses(t *testing.T) {
	h, clock := startLimitHub(t)
	h.send(hubRequest{run: func(h *Hub) {
		f := &restartFile{saveFile: saveFile{Game: "g7", Sets: []string{"b2"}, Saved: clock.Now().Add(-20 * time.Minute)}, Kind: keptTable, N: 7, Host: "Ann"}
		f.Seats = []savedSeat{{Name: "Ann", Token: newToken()}, {Name: "Bo", Token: newToken()}}
		h.restoreTable(f)
	}})
	time.Sleep(50 * time.Millisecond)
	if gamesLeft(t, h) != 1 {
		t.Fatal("the kept table is not listed")
	}
	clock.Add(11 * time.Minute)
	untilNoGames(t, h)
}

// TestHostEndsGame: only the host ends a game; the players hear why,
// one who is away when they come back (19.84).
func TestHostEndsGame(t *testing.T) {
	h, _ := startLimitHub(t)
	players, tokens := seatedTwo(t, h)
	ann, bo := players[0], players[1]
	bo.say(protocol.TypeEnd, struct{}{})
	if code := bo.failure(); code != protocol.ErrNotHost {
		t.Errorf("another seat ends it: %s", code)
	}
	bo.c.Leave()
	ann.say(protocol.TypeEnd, struct{}{})
	var n protocol.Notice
	ann.next(protocol.TypeNotice, &n)
	if n.Kind != protocol.NoticeEnded {
		t.Errorf("notice %+v, want ended", n)
	}
	ann.next(protocol.TypeGames, nil) // back in the lobby
	untilNoGames(t, h)

	back, w := arrive(t, h, "Bo", tokens[1])
	if w.Seat != -1 {
		t.Errorf("seat %d, want the lobby", w.Seat)
	}
	back.next(protocol.TypeNotice, &n)
	if n.Kind != protocol.NoticeEnded {
		t.Errorf("notice %+v, want ended", n)
	}
}
