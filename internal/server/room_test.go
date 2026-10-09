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
	"github.com/FourSoulsRepo/rules_engine/cards"
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

// next waits for the next message.
func (m *memConn) next(t *testing.T) protocol.Envelope {
	t.Helper()
	select {
	case msg := <-m.out:
		e, err := protocol.Decode(msg)
		if err != nil {
			t.Fatal(err)
		}
		return e
	case <-time.After(5 * time.Second):
		t.Fatal("no message")
		return protocol.Envelope{}
	}
}

func startRoom(t *testing.T, players int) *Room {
	t.Helper()
	r, err := NewRoom(engine.Setup{Seed: 7, Players: players, Sets: cards.Sets()})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	go r.Run(ctx)
	t.Cleanup(cancel)
	return r
}

func hello(t *testing.T, c *Client, name, token string) {
	t.Helper()
	msg, err := protocol.Encode(protocol.TypeHello, 1, protocol.Hello{Protocol: protocol.Version, Name: name, Role: protocol.RolePlayer, Sets: []string{"b2"}, Token: token})
	if err != nil {
		t.Fatal(err)
	}
	c.Receive(msg)
}

func TestHelloSeatsAndNames(t *testing.T) {
	r := startRoom(t, 2)
	a, b, x := newMemConn(), newMemConn(), newMemConn()
	ca, cb, cx := r.Connect(a), r.Connect(b), r.Connect(x)
	hello(t, ca, "Ann", "")
	var w protocol.Welcome
	if e := a.next(t); e.Type != protocol.TypeWelcome || e.Unpack(&w) != nil || w.Seat != 0 || w.Token == "" {
		t.Fatalf("welcome %+v", w)
	}
	hello(t, cb, " Ann ", "")
	var wb protocol.Welcome
	if e := b.next(t); e.Unpack(&wb) != nil || wb.Seat != 1 {
		t.Fatalf("second welcome %+v", wb)
	}
	var u protocol.Update
	if e := b.next(t); e.Type != protocol.TypeUpdate || e.Unpack(&u) != nil || u.Seats[1].Name != "Ann (2)" {
		t.Fatalf("update seats %+v", u.Seats)
	}
	hello(t, cx, "Cy", "")
	var er protocol.Error
	if e := x.next(t); e.Type != protocol.TypeError || e.Unpack(&er) != nil || er.Code != protocol.ErrRoomFull {
		t.Fatalf("third player: %+v", er)
	}

	// Ann's connection drops; she comes back with her token.
	ca.Leave()
	a2 := newMemConn()
	c2 := r.Connect(a2)
	hello(t, c2, "Ann", w.Token)
	var back protocol.Welcome
	if e := a2.next(t); e.Unpack(&back) != nil || back.Seat != 0 || back.Token != w.Token {
		t.Fatalf("reconnect %+v", back)
	}
}

func TestBadMessages(t *testing.T) {
	r := startRoom(t, 2)
	m := newMemConn()
	c := r.Connect(m)
	c.Receive([]byte(`{"type":"intent","id":4,"data":{}}`))
	var er protocol.Error
	if e := m.next(t); e.Unpack(&er) != nil || er.Code != protocol.ErrNoHello || e.ID != 4 {
		t.Fatalf("%+v", er)
	}
	c.Receive([]byte(`nonsense`))
	if e := m.next(t); e.Unpack(&er) != nil || er.Code != protocol.ErrBadMessage {
		t.Fatalf("%+v", er)
	}
	old, err := protocol.Encode(protocol.TypeHello, 1, protocol.Hello{Protocol: protocol.Version - 1, Name: "Old"})
	if err != nil {
		t.Fatal(err)
	}
	c.Receive(old)
	if e := m.next(t); e.Unpack(&er) != nil || er.Code != protocol.ErrClientOutdated {
		t.Fatalf("%+v", er)
	}
}

// TestFourClientsRace: four clients play random allowed intents at the
// same time; run with -race (6.2).
func TestFourClientsRace(t *testing.T) {
	r := startRoom(t, 4)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var wg sync.WaitGroup
	var maxStep atomic.Int64
	for i := range 4 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m := newMemConn()
			c := r.Connect(m)
			hello(t, c, "P"+strconv.Itoa(i+1), "")
			rng := rand.New(rand.NewPCG(uint64(i), 1)) //nolint:gosec // test moves
			id := 10
			for {
				select {
				case <-ctx.Done():
					return
				case msg := <-m.out:
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
					id++
					out, err := protocol.Encode(protocol.TypeIntent, id, protocol.Intent{Intent: in})
					if err != nil {
						t.Error(err)
						return
					}
					c.Receive(out)
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
	r := startRoom(t, 2)
	srv := httptest.NewServer(Handler(r))
	defer srv.Close()
	url := "ws" + strings.TrimPrefix(srv.URL, "http")
	dial := func() *websocket.Conn {
		ws, resp, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			t.Fatal(err)
		}
		_ = resp.Body.Close() //nolint:errcheck // test
		return ws
	}
	a, b := dial(), dial()
	defer a.Close() //nolint:errcheck // test
	defer b.Close() //nolint:errcheck // test
	for i, ws := range []*websocket.Conn{a, b} {
		msg, err := protocol.Encode(protocol.TypeHello, 1, protocol.Hello{Protocol: protocol.Version, Name: "P" + strconv.Itoa(i)})
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
		if e, err := protocol.Decode(got); err != nil || e.Type != protocol.TypeWelcome {
			t.Fatalf("got %s", got)
		}
	}
	// Too large a message closes the connection.
	if err := a.WriteMessage(websocket.TextMessage, make([]byte, protocol.MaxClientMessage+1)); err != nil {
		t.Fatal(err)
	}
	if err := a.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		t.Fatal(err)
	}
	for {
		if _, _, err := a.ReadMessage(); err != nil {
			break // closed
		}
	}
}
