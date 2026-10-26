package server

import (
	"context"
	"math/rand/v2"
	"net"
	"net/http/httptest"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// memConn is an in-memory Conn for tests.
type memConn struct {
	out    chan []byte
	once   sync.Once
	closed chan struct{}
}

func newMemConn() *memConn {
	return &memConn{out: make(chan []byte, 1024), closed: make(chan struct{})}
}

func (m *memConn) Send(msg []byte) bool {
	select {
	case <-m.closed:
		return false
	case m.out <- msg:
		return true
	default:
		return false
	}
}

func (m *memConn) Close() { m.once.Do(func() { close(m.closed) }) }

// player is a test client.
type player struct {
	t    *testing.T
	conn *memConn
	c    *Client
	id   int
}

// next waits for the next message of the given type, into v.
func (p *player) next(typ string, v any) protocol.Envelope {
	p.t.Helper()
	for {
		select {
		case msg := <-p.conn.out:
			e, err := protocol.Decode(msg)
			if err != nil {
				p.t.Fatal(err)
			}
			if e.Type != typ {
				continue
			}
			if v != nil {
				if err := e.Unpack(v); err != nil {
					p.t.Fatal(err)
				}
			}
			return e
		case <-time.After(5 * time.Second):
			p.t.Fatalf("no %s message", typ)
			return protocol.Envelope{}
		}
	}
}

func (p *player) say(typ string, data any) {
	p.t.Helper()
	p.id++
	msg, err := protocol.Encode(typ, p.id, data)
	if err != nil {
		p.t.Fatal(err)
	}
	p.c.Receive(msg)
}

// failure waits for an error and returns its code.
func (p *player) failure() string {
	p.t.Helper()
	var er protocol.Error
	p.next(protocol.TypeError, &er)
	return er.Code
}

func startHub(t *testing.T) *Hub {
	t.Helper()
	return startHubWith(t, time.Second)
}

// startHubWith runs a hub whose timers count second as one second.
func startHubWith(t *testing.T, second time.Duration) *Hub {
	t.Helper()
	h := NewHub()
	h.second = second
	ctx, cancel := context.WithCancel(context.Background())
	go h.Run(ctx)
	t.Cleanup(cancel)
	return h
}

// arrive connects a player and says hello with the Base Game.
func arrive(t *testing.T, h *Hub, name, token string) (*player, protocol.Welcome) {
	t.Helper()
	p := &player{t: t, conn: newMemConn()}
	p.c = h.Connect(p.conn)
	p.say(protocol.TypeHello, protocol.Hello{Protocol: protocol.Version, Name: name, Role: protocol.RolePlayer, Sets: []string{"b2"}, Token: token})
	var w protocol.Welcome
	p.next(protocol.TypeWelcome, &w)
	return p, w
}

// startTable starts a game of n players and returns them, seated.
func startTable(t *testing.T, h *Hub, n int) []*player {
	t.Helper()
	return startTableWith(t, h, protocol.Create{Seats: n})
}

// startTableWith starts a game created by m and returns its players.
func startTableWith(t *testing.T, h *Hub, m protocol.Create) []*player {
	t.Helper()
	n := m.Seats
	host, _ := arrive(t, h, "P0", "")
	host.say(protocol.TypeCreate, m)
	var tb protocol.Table
	host.next(protocol.TypeTable, &tb)
	players := []*player{host}
	for i := 1; i < n; i++ {
		p, _ := arrive(t, h, "P"+strconv.Itoa(i), "")
		p.say(protocol.TypeJoin, protocol.Join{Game: tb.Game})
		p.next(protocol.TypeTable, nil)
		players = append(players, p)
	}
	for _, p := range players {
		p.say(protocol.TypeReady, protocol.Ready{Ready: true})
	}
	return players
}

func TestLobby(t *testing.T) {
	h := startHub(t)
	ann, w := arrive(t, h, "Ann", "")
	if w.Seat != -1 || w.Token == "" {
		t.Fatalf("welcome %+v", w)
	}
	var games protocol.Games
	ann.next(protocol.TypeGames, &games)
	ann.say(protocol.TypeCreate, protocol.Create{Seats: 2, Sets: []string{"b2"}})
	var tb protocol.Table
	ann.next(protocol.TypeTable, &tb)
	if tb.You != 0 || tb.Seats[0].Name != "Ann" || len(tb.Seats) != 2 {
		t.Fatalf("table %+v", tb)
	}

	bob, _ := arrive(t, h, "Ann", "") // the same nickname
	bob.next(protocol.TypeGames, &games)
	if len(games.Games) != 1 || games.Games[0].Taken != 1 || games.Games[0].Host != "Ann" {
		t.Fatalf("games %+v", games)
	}
	bob.say(protocol.TypeJoin, protocol.Join{Game: tb.Game})
	bob.next(protocol.TypeTable, &tb)
	if tb.You != 1 || tb.Seats[1].Name != "Ann (2)" {
		t.Fatalf("second seat %+v", tb)
	}

	ann.say(protocol.TypeReady, protocol.Ready{Ready: true})
	for !tb.Seats[0].Ready {
		ann.next(protocol.TypeTable, &tb)
	}
	if tb.Seats[1].Ready {
		t.Fatalf("ready %+v", tb.Seats)
	}
	bob.say(protocol.TypeReady, protocol.Ready{Ready: true})
	var u protocol.Update
	bob.next(protocol.TypeUpdate, &u)
	if len(u.Seats) != 2 || u.Seats[0].Name != "Ann" || len(u.View.Players) != 2 {
		t.Fatalf("game start %+v", u.Seats)
	}
	ann.next(protocol.TypeUpdate, &u)

	late, _ := arrive(t, h, "Cy", "")
	late.say(protocol.TypeJoin, protocol.Join{Game: tb.Game})
	if code := late.failure(); code != protocol.ErrStarted {
		t.Errorf("joining a started game: %s", code)
	}
}

func TestLobbyRefusals(t *testing.T) {
	h := startHub(t)
	ann, _ := arrive(t, h, "Ann", "")
	ann.say(protocol.TypeCreate, protocol.Create{Seats: 5})
	if code := ann.failure(); code != protocol.ErrBadSetup {
		t.Errorf("5 seats: %s", code)
	}
	ann.say(protocol.TypeCreate, protocol.Create{Seats: 2, Sets: []string{"zz"}})
	if code := ann.failure(); code != protocol.ErrBadSetup {
		t.Errorf("unknown set: %s", code)
	}
	ann.say(protocol.TypeCreate, protocol.Create{Seats: 2})
	var tb protocol.Table
	ann.next(protocol.TypeTable, &tb)

	poor := &player{t: t, conn: newMemConn()}
	poor.c = h.Connect(poor.conn)
	poor.say(protocol.TypeHello, protocol.Hello{Protocol: protocol.Version, Name: "Poor"})
	poor.next(protocol.TypeWelcome, nil)
	poor.say(protocol.TypeJoin, protocol.Join{Game: tb.Game})
	if code := poor.failure(); code != protocol.ErrMissingSets {
		t.Errorf("no card sets: %s", code)
	}
	poor.say(protocol.TypeJoin, protocol.Join{Game: "nope"})
	if code := poor.failure(); code != protocol.ErrNoGame {
		t.Errorf("no game: %s", code)
	}

	bob, _ := arrive(t, h, "Bob", "")
	bob.say(protocol.TypeJoin, protocol.Join{Game: tb.Game})
	bob.next(protocol.TypeTable, nil)
	cy, _ := arrive(t, h, "Cy", "")
	cy.say(protocol.TypeJoin, protocol.Join{Game: tb.Game})
	if code := cy.failure(); code != protocol.ErrRoomFull {
		t.Errorf("full table: %s", code)
	}
	ann.say(protocol.TypeIntent, protocol.Intent{})
	if code := ann.failure(); code != protocol.ErrBadMessage {
		t.Errorf("intent before the start: %s", code)
	}
}

func TestBadMessages(t *testing.T) {
	h := startHub(t)
	p := &player{t: t, conn: newMemConn()}
	p.c = h.Connect(p.conn)
	p.c.Receive([]byte(`{"type":"list","id":4}`))
	if e := p.next(protocol.TypeError, nil); e.ID != 4 {
		t.Errorf("reply id %d", e.ID)
	}
	p.c.Receive([]byte(`nonsense`))
	if code := p.failure(); code != protocol.ErrBadMessage {
		t.Errorf("%s", code)
	}
	p.say(protocol.TypeHello, protocol.Hello{Protocol: protocol.Version, Name: "  "})
	if code := p.failure(); code != protocol.ErrBadName {
		t.Errorf("%s", code)
	}
	p.say(protocol.TypeHello, protocol.Hello{Protocol: protocol.Version - 1, Name: "Old"})
	if code := p.failure(); code != protocol.ErrClientOutdated {
		t.Errorf("%s", code)
	}
}

// TestReconnect: a player comes back with their token, at the table and
// during the game (N-08).
func TestReconnect(t *testing.T) {
	h := startHub(t)
	ann, w := arrive(t, h, "Ann", "")
	ann.say(protocol.TypeCreate, protocol.Create{Seats: 2})
	var tb protocol.Table
	ann.next(protocol.TypeTable, &tb)
	bob, _ := arrive(t, h, "Bob", "")
	bob.say(protocol.TypeJoin, protocol.Join{Game: tb.Game})
	bob.next(protocol.TypeTable, nil)

	// Ann's connection drops at the table: her seat is freed.
	ann.c.Leave()
	for tb.Seats[0].Name != "" {
		bob.next(protocol.TypeTable, &tb)
	}
	ann2, _ := arrive(t, h, "Ann", w.Token)
	ann2.say(protocol.TypeJoin, protocol.Join{Game: tb.Game})
	ann2.next(protocol.TypeTable, nil)
	ann2.say(protocol.TypeReady, protocol.Ready{Ready: true})
	bob.say(protocol.TypeReady, protocol.Ready{Ready: true})
	ann2.next(protocol.TypeUpdate, nil)

	// During the game the seat waits for its token.
	ann2.c.Leave()
	var u protocol.Update
	for u.Seats == nil || u.Seats[0].Connected {
		bob.next(protocol.TypeUpdate, &u)
	}
	ann3, back := arrive(t, h, "Ann", w.Token)
	if back.Token != w.Token || back.Seat != 0 {
		t.Fatalf("token %q, seat %d; want the old token and seat 0", back.Token, back.Seat)
	}
	ann3.next(protocol.TypeUpdate, &u)
	if !u.Seats[0].Connected || u.View.Players[0].HandSize == 0 {
		t.Errorf("back in the game: %+v", u.Seats)
	}
}

// TestFinishedGameLetsGo: a token from a finished game gets the lobby,
// keeps the token, and takes its player to a new game instead (19.32).
func TestFinishedGameLetsGo(t *testing.T) {
	h := startHub(t)
	const token = "0123456789abcdef0123456789abcdef"
	old := &Room{tokens: []string{token, "other"}}
	old.over.Store(true)
	done := make(chan struct{})
	h.send(hubRequest{run: func(h *Hub) {
		h.games = append(h.games, started{info: protocol.GameInfo{ID: "g0", Started: true}, room: old})
		close(done)
	}})
	<-done

	ann, w := arrive(t, h, "Ann", token)
	if w.Token != token || w.Seat != -1 {
		t.Fatalf("token %q, seat %d; want the old token and no seat", w.Token, w.Seat)
	}
	ann.next(protocol.TypeGames, nil)
	if ann.c.room.Load() != nil {
		t.Fatal("the hello took Ann back to the finished game")
	}

	// A new game with the same token is the one she comes back to.
	ann.say(protocol.TypeCreate, protocol.Create{Seats: 2})
	var tb protocol.Table
	ann.next(protocol.TypeTable, &tb)
	bob, _ := arrive(t, h, "Bob", "")
	bob.say(protocol.TypeJoin, protocol.Join{Game: tb.Game})
	bob.next(protocol.TypeTable, nil)
	ann.say(protocol.TypeReady, protocol.Ready{Ready: true})
	bob.say(protocol.TypeReady, protocol.Ready{Ready: true})
	ann.next(protocol.TypeUpdate, nil)
	ann.c.Leave()
	var u protocol.Update
	for u.Seats == nil || u.Seats[0].Connected {
		bob.next(protocol.TypeUpdate, &u)
	}
	ann2, back := arrive(t, h, "Ann", token)
	if back.Seat != 0 {
		t.Fatalf("seat %d; want 0 in the new game", back.Seat)
	}
	ann2.next(protocol.TypeUpdate, &u)
	if !u.Seats[0].Connected {
		t.Errorf("not back in the new game: %+v", u.Seats)
	}
	if r := ann2.c.room.Load(); r == nil || r == old {
		t.Errorf("room %p; want the new game, not %p", r, old)
	}
}

// arriveWithDice connects a player whose hello carries a dice look.
func arriveWithDice(t *testing.T, h *Hub, name, token, dice string) *player {
	t.Helper()
	p := &player{t: t, conn: newMemConn()}
	p.c = h.Connect(p.conn)
	p.say(protocol.TypeHello, protocol.Hello{Protocol: protocol.Version, Name: name, Role: protocol.RolePlayer, Sets: []string{"b2"}, Token: token, Dice: dice})
	p.next(protocol.TypeWelcome, nil)
	return p
}

// TestSeatsShowDice: each seat's dice look reaches the other players,
// at the table and in the game; a bad ID is the default (V-19).
func TestSeatsShowDice(t *testing.T) {
	h := startHub(t)
	ann := arriveWithDice(t, h, "Ann", "", "tear")
	ann.say(protocol.TypeCreate, protocol.Create{Seats: 2})
	var tb protocol.Table
	ann.next(protocol.TypeTable, &tb)
	bob := arriveWithDice(t, h, "Bob", "", "Bad look!")
	bob.say(protocol.TypeJoin, protocol.Join{Game: tb.Game})
	bob.next(protocol.TypeTable, &tb)
	if tb.Seats[0].Dice != "tear" || tb.Seats[1].Dice != "" {
		t.Errorf("table seats %+v; want tear and the default", tb.Seats)
	}
	ann.say(protocol.TypeReady, protocol.Ready{Ready: true})
	bob.say(protocol.TypeReady, protocol.Ready{Ready: true})
	var u protocol.Update
	bob.next(protocol.TypeUpdate, &u)
	if u.Seats[0].Dice != "tear" || u.Seats[1].Dice != "" {
		t.Errorf("game seats %+v; want tear and the default", u.Seats)
	}

	// Away, the seat keeps its look; back, it takes the new hello's.
	token := ann.c.token
	ann.c.Leave()
	u = bob.untilUpdate(func(u protocol.Update) bool { return u.Pause != nil })
	if u.Seats[0].Dice != "tear" {
		t.Errorf("away seat %+v; want it to keep tear", u.Seats[0])
	}
	arriveWithDice(t, h, "Ann", token, "gold")
	u = bob.untilUpdate(func(u protocol.Update) bool { return u.Pause == nil })
	if u.Seats[0].Dice != "gold" {
		t.Errorf("returned seat %+v; want gold", u.Seats[0])
	}
}

// TestFourClientsRace: four clients play random allowed intents at the
// same time; run with -race (6.2).
func TestFourClientsRace(t *testing.T) {
	h := startHub(t)
	players := startTable(t, h, 4)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var wg sync.WaitGroup
	var maxStep atomic.Int64
	for i, p := range players {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rng := rand.New(rand.NewPCG(uint64(i), 1)) //nolint:gosec // test moves
			for {
				select {
				case <-ctx.Done():
					return
				case msg := <-p.conn.out:
					e, err := protocol.Decode(msg)
					if err != nil {
						t.Error(err)
						return
					}
					var u protocol.Update
					if e.Type != protocol.TypeUpdate || e.Unpack(&u) != nil {
						continue
					}
					if int64(u.Step) > maxStep.Load() {
						maxStep.Store(int64(u.Step))
					}
					if len(u.Allowed) == 0 || u.View.Over {
						continue
					}
					in := u.Allowed[rng.IntN(len(u.Allowed))]
					if in.Kind == engine.IntentDiscard {
						in.Objects = in.Objects[:u.View.Waiting.Count]
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
	wg.Wait()
	if n := maxStep.Load(); n < 50 {
		t.Errorf("only %d steps in 3 seconds", n)
	}
}

func TestWebSocket(t *testing.T) {
	h := startHub(t)
	srv := httptest.NewServer(Handler(h))
	defer srv.Close()
	url := "ws" + strings.TrimPrefix(srv.URL, "http")
	ws, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close() //nolint:errcheck // test
	defer ws.Close()      //nolint:errcheck // test
	msg, err := protocol.Encode(protocol.TypeHello, 1, protocol.Hello{Protocol: protocol.Version, Name: "Ann"})
	if err != nil {
		t.Fatal(err)
	}
	if werr := ws.WriteMessage(websocket.TextMessage, msg); werr != nil {
		t.Fatal(werr)
	}
	_, got, err := ws.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	if e, derr := protocol.Decode(got); derr != nil || e.Type != protocol.TypeWelcome {
		t.Fatalf("got %s", got)
	}
	// Too large a message closes the connection.
	if werr := ws.WriteMessage(websocket.TextMessage, make([]byte, protocol.MaxClientMessage+1)); werr != nil {
		t.Fatal(werr)
	}
	if derr := ws.SetReadDeadline(time.Now().Add(5 * time.Second)); derr != nil {
		t.Fatal(derr)
	}
	for {
		if _, _, rerr := ws.ReadMessage(); rerr != nil {
			if strings.Contains(rerr.Error(), "timeout") {
				t.Error("the connection is still open")
			}
			break
		}
	}
}

// A game password keeps others out (N-14); it is never sent back.
func TestGamePassword(t *testing.T) {
	const secret = "s3cret-pw"
	h := startHub(t)
	ann, _ := arrive(t, h, "Ann", "")
	ann.say(protocol.TypeCreate, protocol.Create{Seats: 2, Password: strings.Repeat("x", protocol.MaxPassword+1)})
	if code := ann.failure(); code != protocol.ErrBadSetup {
		t.Errorf("too long a password: %s", code)
	}
	ann.say(protocol.TypeCreate, protocol.Create{Seats: 2, Password: secret})
	var tb protocol.Table
	if e := ann.next(protocol.TypeTable, &tb); strings.Contains(string(e.Data), secret) {
		t.Errorf("the table carries the password: %s", e.Data)
	}

	bo, _ := arrive(t, h, "Bo", "")
	var games protocol.Games
	if e := bo.next(protocol.TypeGames, &games); strings.Contains(string(e.Data), secret) {
		t.Errorf("the games list carries the password: %s", e.Data)
	}
	if len(games.Games) != 1 || !games.Games[0].Locked {
		t.Fatalf("games %+v, want one locked game", games.Games)
	}
	for _, c := range []struct{ password, reason string }{
		{"", "this game needs a password"},
		{"guess", "wrong password"},
	} {
		bo.say(protocol.TypeJoin, protocol.Join{Game: tb.Game, Password: c.password})
		var er protocol.Error
		bo.next(protocol.TypeError, &er)
		if er.Code != protocol.ErrPassword || er.Message != c.reason {
			t.Errorf("password %q: %+v, want %q", c.password, er, c.reason)
		}
	}
	bo.say(protocol.TypeJoin, protocol.Join{Game: tb.Game, Password: secret})
	if e := bo.next(protocol.TypeTable, &tb); tb.You != 1 || strings.Contains(string(e.Data), secret) {
		t.Errorf("the right password: table %+v", tb)
	}
}

func TestWithBaseGame(t *testing.T) {
	for _, c := range []struct{ in, want []string }{
		{nil, []string{"b2"}},
		{[]string{"g2"}, []string{"b2", "g2"}},
		{[]string{"fsp2", "b2", "fsp2"}, []string{"b2", "fsp2"}},
	} {
		if got := withBaseGame(c.in); !slices.Equal(got, c.want) {
			t.Errorf("withBaseGame(%v) = %v, want %v", c.in, got, c.want)
		}
	}
}

// localMemConn is a memConn whose client runs on the server's computer.
type localMemConn struct{ *memConn }

func (localMemConn) Local() bool { return true }

// arriveLocal is arrive from the server's own computer.
func arriveLocal(t *testing.T, h *Hub, name string) (*player, protocol.Welcome) {
	t.Helper()
	p := &player{t: t, conn: newMemConn()}
	p.c = h.Connect(localMemConn{p.conn})
	p.say(protocol.TypeHello, protocol.Hello{Protocol: protocol.Version, Name: name, Role: protocol.RolePlayer, Sets: []string{"b2"}})
	var w protocol.Welcome
	p.next(protocol.TypeWelcome, &w)
	return p, w
}

func TestDedicatedLetsAnyoneCreate(t *testing.T) {
	h := startHub(t)
	ann, w := arrive(t, h, "Ann", "")
	if !w.Create {
		t.Fatal("a dedicated server does not offer create")
	}
	ann.say(protocol.TypeCreate, protocol.Create{Seats: 2})
	ann.next(protocol.TypeTable, nil)
}

func TestHostOnlyCreates(t *testing.T) {
	dir := t.TempDir()
	writeSave(t, dir, "g1", engine.Version)
	h := startHubWithSaves(t, dir)
	h.hostOnly = true // before anyone connects; the hub reads it later

	bob, w := arrive(t, h, "Bob", "")
	if w.Create {
		t.Error("a friend is offered create on a hosted server")
	}
	bob.say(protocol.TypeCreate, protocol.Create{Seats: 2})
	if code := bob.failure(); code != protocol.ErrNotHost {
		t.Errorf("a friend creates: %s", code)
	}
	bob.say(protocol.TypeLoad, protocol.Load{Save: "g1"})
	if code := bob.failure(); code != protocol.ErrNotHost {
		t.Errorf("a friend loads: %s", code)
	}

	ann, w := arriveLocal(t, h, "Ann")
	if !w.Create {
		t.Error("the host is not offered create")
	}
	ann.say(protocol.TypeLoad, protocol.Load{Save: "g1"})
	ann.next(protocol.TypeTable, nil)
	// The host opened the saved game; the friend sits down in it.
	bob.say(protocol.TypeLoad, protocol.Load{Save: "g1"})
	var tb protocol.Table
	bob.next(protocol.TypeTable, &tb)
	if tb.You != 1 {
		t.Errorf("friend's seat %d", tb.You)
	}
}

func TestLoopback(t *testing.T) {
	for _, c := range []struct {
		addr net.Addr
		want bool
	}{
		{&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)}, true},
		{&net.TCPAddr{IP: net.IPv6loopback}, true},
		{&net.TCPAddr{IP: net.IPv4(192, 168, 1, 5)}, false},
		{&net.UnixAddr{Name: "x"}, false},
	} {
		if got := loopback(c.addr); got != c.want {
			t.Errorf("%v: %v", c.addr, got)
		}
	}
}

// TestServerName: the welcome carries the server's name; an app's own
// host has none (N-18).
func TestServerName(t *testing.T) {
	if _, w := arrive(t, startHub(t), "Ann", ""); w.Name != "" {
		t.Errorf("default hub sends name %q", w.Name)
	}
	h := NewHub()
	h.name = "Basement"
	ctx, cancel := context.WithCancel(context.Background())
	go h.Run(ctx)
	t.Cleanup(cancel)
	if _, w := arrive(t, h, "Ann", ""); w.Name != "Basement" {
		t.Errorf("welcome name %q", w.Name)
	}
}
