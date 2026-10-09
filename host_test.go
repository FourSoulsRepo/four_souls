package main

import (
	"errors"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// TestHostAndRemoteClient: the app hosts; the host and a friend join
// over the LAN address and both play (6.4).
func TestHostAndRemoteClient(t *testing.T) {
	h := &Host{}
	info, err := h.Start(2)
	if err != nil {
		t.Skipf("cannot host here: %v", err)
	}
	defer func() {
		if err := h.Stop(); err != nil {
			t.Error(err)
		}
	}()
	if _, again := h.Start(2); !errors.Is(again, ErrHosting) {
		t.Error("hosted twice")
	}
	addr := "127.0.0.1"
	if len(info.Addresses) > 0 {
		addr = info.Addresses[0] // as a friend on the LAN would
	}
	url := "ws://" + net.JoinHostPort(addr, strconv.Itoa(info.Port)) + "/ws"
	var players [2]*websocket.Conn
	var seats [2]int
	for i := range players {
		ws, resp, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			t.Fatal(err)
		}
		_ = resp.Body.Close() //nolint:errcheck // test
		defer ws.Close()      //nolint:errcheck // test
		players[i] = ws
		msg, err := protocol.Encode(protocol.TypeHello, 1, protocol.Hello{Protocol: protocol.Version, Name: "P" + strconv.Itoa(i)})
		if err != nil {
			t.Fatal(err)
		}
		if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
			t.Fatal(err)
		}
		var w protocol.Welcome
		readType(t, ws, protocol.TypeWelcome, &w)
		seats[i] = w.Seat
	}
	// Play a few steps: whoever is asked answers with an allowed intent.
	steps := 0
	for steps < 6 {
		for i, ws := range players {
			var u protocol.Update
			readType(t, ws, protocol.TypeUpdate, &u)
			if len(u.Allowed) == 0 || u.View.Over || int(u.View.Waiting.Player) != seats[i] {
				continue
			}
			in := u.Allowed[0]
			if in.Kind == engine.IntentDiscard { // a discard needs Count cards
				in.Objects = in.Objects[:u.View.Waiting.Count]
			}
			msg, err := protocol.Encode(protocol.TypeIntent, 2, protocol.Intent{Intent: in})
			if err != nil {
				t.Fatal(err)
			}
			if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
				t.Fatal(err)
			}
			steps++
		}
	}
}

// readType reads messages until one of the given type, into v.
func readType(t *testing.T, ws *websocket.Conn, typ string, v any) {
	t.Helper()
	if err := ws.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		t.Fatal(err)
	}
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			t.Fatal(err)
		}
		e, err := protocol.Decode(msg)
		if err != nil {
			t.Fatal(err)
		}
		if e.Type == typ {
			if err := e.Unpack(v); err != nil {
				t.Fatal(err)
			}
			return
		}
	}
}
