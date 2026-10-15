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

// TestHostAndRemoteClient: the app hosts; the host joins over loopback,
// as the app does, and a friend over the LAN address; both play (6.4).
// Only the host may open the game (N-17).
func TestHostAndRemoteClient(t *testing.T) {
	h := &Host{}
	info, err := h.Start()
	if err != nil {
		t.Skipf("cannot host here: %v", err)
	}
	defer func() {
		if err := h.Stop(); err != nil {
			t.Error(err)
		}
	}()
	if _, again := h.Start(); !errors.Is(again, ErrHosting) {
		t.Error("hosted twice")
	}
	addrs := [2]string{"127.0.0.1", "127.0.0.1"}
	if len(info.Addresses) > 0 {
		addrs[1] = info.Addresses[0] // as a friend on the LAN would
	}
	var players [2]*websocket.Conn
	var seats [2]int
	for i := range players {
		url := "ws://" + net.JoinHostPort(addrs[i], strconv.Itoa(info.Port)) + "/ws"
		ws, resp, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			t.Fatal(err)
		}
		_ = resp.Body.Close() //nolint:errcheck // test
		defer ws.Close()      //nolint:errcheck // test
		players[i] = ws
		say(t, ws, protocol.TypeHello, protocol.Hello{Protocol: protocol.Version, Name: "P" + strconv.Itoa(i), Sets: []string{"b2"}})
		var w protocol.Welcome
		readType(t, ws, protocol.TypeWelcome, &w)
		if friend := addrs[i] != "127.0.0.1"; w.Create == friend {
			t.Errorf("player %d over %s: create %v", i, addrs[i], w.Create)
		}
	}
	if addrs[1] != "127.0.0.1" {
		say(t, players[1], protocol.TypeCreate, protocol.Create{Seats: 2})
		var er protocol.Error
		if readType(t, players[1], protocol.TypeError, &er); er.Code != protocol.ErrNotHost {
			t.Errorf("a friend creates: %+v", er)
		}
	}
	// The host creates a game in its own lobby; the friend joins it.
	say(t, players[0], protocol.TypeCreate, protocol.Create{Seats: 2})
	var tb protocol.Table
	readType(t, players[0], protocol.TypeTable, &tb)
	say(t, players[1], protocol.TypeJoin, protocol.Join{Game: tb.Game})
	for i, ws := range players {
		readType(t, ws, protocol.TypeTable, &tb)
		seats[i] = tb.You
		say(t, ws, protocol.TypeReady, protocol.Ready{Ready: true})
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
			say(t, ws, protocol.TypeIntent, protocol.Intent{Intent: in})
			steps++
		}
	}
}

// say sends one message.
func say(t *testing.T, ws *websocket.Conn, typ string, data any) {
	t.Helper()
	msg, err := protocol.Encode(typ, 1, data)
	if err != nil {
		t.Fatal(err)
	}
	if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
		t.Fatal(err)
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
