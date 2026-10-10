package server

import (
	"context"
	"math/rand/v2"
	"net/http/httptest"
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
	host, _ := arrive(t, h, "P0", "")
	host.say(protocol.TypeCreate, protocol.Create{Seats: n})
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
	if back.Token != w.Token {
		t.Fatalf("token %q", back.Token)
	}
	ann3.next(protocol.TypeUpdate, &u)
	if !u.Seats[0].Connected || u.View.Players[0].HandSize == 0 {
		t.Errorf("back in the game: %+v", u.Seats)
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
