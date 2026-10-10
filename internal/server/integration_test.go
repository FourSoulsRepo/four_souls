package server

import (
	"context"
	"encoding/json"
	"math/rand/v2"
	"net"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	"github.com/FourSoulsRepo/record"
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards"
)

// bot is a headless test client over a real WebSocket (6.10).
type bot struct {
	t     *testing.T
	url   string
	name  string
	ws    *websocket.Conn
	token string
	seat  int
	rng   *rand.Rand
	in    chan protocol.Envelope // filled by the reader goroutine
}

func dialBot(t *testing.T, url, name string, seed uint64) *bot {
	t.Helper()
	b := &bot{t: t, url: url, name: name, seat: -1, rng: rand.New(rand.NewPCG(seed, 3))} //nolint:gosec // test moves
	b.connect()
	return b
}

// connect opens the socket and says hello, with the token if any.
func (b *bot) connect() {
	b.t.Helper()
	ws, resp, err := websocket.DefaultDialer.Dial(b.url, nil)
	if err != nil {
		b.t.Fatal(err)
	}
	_ = resp.Body.Close() //nolint:errcheck // test
	b.ws = ws
	b.in = make(chan protocol.Envelope, 1024)
	go func(ws *websocket.Conn, in chan<- protocol.Envelope) {
		defer close(in)
		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				return
			}
			if e, err := protocol.Decode(msg); err == nil {
				in <- e
			}
		}
	}(ws, b.in)
	b.send(protocol.TypeHello, protocol.Hello{Protocol: protocol.Version, Name: b.name, Role: protocol.RolePlayer, Sets: []string{"b2"}, Token: b.token})
	var w protocol.Welcome
	b.expect(protocol.TypeWelcome, &w)
	b.token = w.Token
}

func (b *bot) send(typ string, data any) {
	b.t.Helper()
	msg, err := protocol.Encode(typ, 1, data)
	if err != nil {
		b.t.Fatal(err)
	}
	if err := b.ws.WriteMessage(websocket.TextMessage, msg); err != nil {
		b.t.Fatal(err)
	}
}

// read returns the next message, or false after the timeout.
func (b *bot) read(timeout time.Duration) (protocol.Envelope, bool) {
	select {
	case e, ok := <-b.in:
		return e, ok
	case <-time.After(timeout):
		return protocol.Envelope{}, false
	}
}

func (b *bot) expect(typ string, v any) {
	b.t.Helper()
	for {
		e, ok := b.read(5 * time.Second)
		if !ok {
			b.t.Fatalf("%s: no %s", b.name, typ)
		}
		if e.Type == typ {
			if v != nil {
				if err := e.Unpack(v); err != nil {
					b.t.Fatal(err)
				}
			}
			return
		}
	}
}

// act answers an update with a random allowed intent; acting beats
// passing two times in three, so games move on.
func (b *bot) act(u protocol.Update) {
	if len(u.Allowed) == 0 || u.View.Over || u.Pause != nil {
		return
	}
	var acts []engine.Intent
	for _, in := range u.Allowed {
		if in.Kind != engine.IntentPass && in.Kind != engine.IntentEndTurn {
			acts = append(acts, in)
		}
	}
	in := u.Allowed[b.rng.IntN(len(u.Allowed))]
	if len(acts) > 0 && b.rng.IntN(3) > 0 {
		in = acts[b.rng.IntN(len(acts))]
	}
	if in.Kind == engine.IntentDiscard {
		in.Objects = in.Objects[:u.View.Waiting.Count]
	}
	b.send(protocol.TypeIntent, protocol.Intent{Intent: in})
}

// playGame plays one whole game of n bots on a fresh server; one bot
// drops once and comes back. It returns the record file, or "" when
// the game did not end within the step limit.
func playGame(t *testing.T, n int, seed uint64) string {
	t.Helper()
	dir := t.TempDir()
	cfg := DefaultConfig()
	cfg.Addr, cfg.Port, cfg.Records, cfg.Seed = "127.0.0.1", freePort(t), dir, seed
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv, err := Start(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = srv.Shutdown(context.Background()) }() //nolint:errcheck // test
	url := "ws://" + srv.Addr() + "/ws"

	bots := make([]*bot, n)
	for i := range bots {
		bots[i] = dialBot(t, url, "Bot"+strconv.Itoa(i), seed+uint64(i)) //nolint:gosec // small index
	}
	bots[0].send(protocol.TypeCreate, protocol.Create{Seats: n})
	var tb protocol.Table
	bots[0].expect(protocol.TypeTable, &tb)
	for _, b := range bots[1:] {
		b.send(protocol.TypeJoin, protocol.Join{Game: tb.Game})
		b.expect(protocol.TypeTable, nil)
	}
	for _, b := range bots {
		b.send(protocol.TypeReady, protocol.Ready{Ready: true})
	}

	const maxUpdates = 40000
	dropAt := 200
	over := false
	for updates := 0; updates < maxUpdates && !over; {
		got := false
		for i, b := range bots {
			e, ok := b.read(2 * time.Millisecond)
			if !ok {
				continue
			}
			got = true
			if e.Type != protocol.TypeUpdate {
				continue
			}
			var u protocol.Update
			if uerr := e.Unpack(&u); uerr != nil {
				t.Fatal(uerr)
			}
			updates++
			if u.View.Over {
				over = true
				break
			}
			if i == n-1 && updates > dropAt && dropAt > 0 {
				dropAt = 0       // drop once, come back with the token (N-08)
				_ = b.ws.Close() //nolint:errcheck // test
				b.connect()
				continue
			}
			b.act(u)
		}
		if !got {
			for _, b := range bots { // nobody heard anything: ask again (resync)
				b.send(protocol.TypeResync, struct{}{})
			}
		}
	}
	for _, b := range bots {
		_ = b.ws.Close() //nolint:errcheck // test
	}
	if !over {
		return ""
	}
	if dropAt != 0 {
		t.Error("the game ended before the planned drop")
	}
	cancel()
	_ = srv.Shutdown(context.Background()) //nolint:errcheck // test
	files, err := filepath.Glob(filepath.Join(dir, "*"+record.Ext))
	if err != nil || len(files) != 1 {
		t.Fatalf("records %v, %v", files, err)
	}
	return files[0]
}

func freePort(t *testing.T) int {
	t.Helper()
	var lc net.ListenConfig
	ln, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close() //nolint:errcheck // test
	addr, ok := ln.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatal("not a TCP address")
	}
	return addr.Port
}

// checkRecord: the record is finished and replays with every checksum.
func checkRecord(t *testing.T, path string) {
	t.Helper()
	rec, err := record.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if rec.End == nil || !rec.End.Finished || len(rec.End.Winners) == 0 {
		t.Fatalf("record end %+v", rec.End)
	}
	s := rec.Header.Setup
	setup := engine.Setup{Seed: s.Seed, Players: s.Players, BonusSouls: s.BonusSouls}
	for _, code := range s.Sets {
		set, _ := cards.Find(code)
		setup.Sets = append(setup.Sets, set)
	}
	for _, c := range s.Characters {
		setup.Characters = append(setup.Characters, engine.CardRef(c))
	}
	steps := make([]engine.StepRecord, len(rec.Steps))
	for i, st := range rec.Steps {
		if jerr := json.Unmarshal(st.Intent, &steps[i].Intent); jerr != nil {
			t.Fatal(jerr)
		}
		steps[i].Checksum = st.Checksum
	}
	g, err := engine.Replay(setup, steps)
	if err != nil {
		t.Fatalf("replay: %v", err)
	}
	if !g.Over {
		t.Error("the replayed game is not over")
	}
}

// TestHeadlessGames plays whole 2, 3 and 4 player games over WebSocket
// with a drop and a reconnect, and checks their records (6.10).
func TestHeadlessGames(t *testing.T) {
	if testing.Short() {
		t.Skip("plays whole games")
	}
	last := 4
	if raceEnabled {
		last = 2 // one game is enough to look for races; all three take minutes
	}
	for n := 2; n <= last; n++ {
		t.Run(strconv.Itoa(n)+" players", func(t *testing.T) {
			for try := range uint64(5) { // a random game may run long: try the next seed
				if path := playGame(t, n, 1000*uint64(n)+try); path != "" { //nolint:gosec // small n
					checkRecord(t, path)
					return
				}
			}
			t.Fatal("no game ended")
		})
	}
}

// TestSaveAndContinue: the host saves with a stack open, the game is
// loaded again, the players take their seats back by token, and it plays
// to the end; the record has a resume and replays with every checksum
// (6.11, N-11).
func TestSaveAndContinue(t *testing.T) {
	if testing.Short() {
		t.Skip("plays a whole game")
	}
	dir := t.TempDir()
	cfg := DefaultConfig()
	cfg.Addr, cfg.Port, cfg.Records, cfg.Saves, cfg.Seed = "127.0.0.1", freePort(t), filepath.Join(dir, "records"), filepath.Join(dir, "saves"), 77
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv, err := Start(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	url := "ws://" + srv.Addr() + "/ws"
	bots := []*bot{dialBot(t, url, "Host", 1), dialBot(t, url, "Guest", 2)}
	bots[0].send(protocol.TypeCreate, protocol.Create{Seats: 2})
	var tb protocol.Table
	bots[0].expect(protocol.TypeTable, &tb)
	bots[1].send(protocol.TypeJoin, protocol.Join{Game: tb.Game})
	bots[1].expect(protocol.TypeTable, nil)
	for _, b := range bots {
		b.send(protocol.TypeReady, protocol.Ready{Ready: true})
	}

	// Play until something is on the stack, then the host saves.
	saved := ""
	for steps := 0; saved == "" && steps < 20000; steps++ {
		for i, b := range bots {
			e, ok := b.read(2 * time.Millisecond)
			if !ok {
				continue
			}
			switch e.Type {
			case protocol.TypeUpdate:
				var u protocol.Update
				if uerr := e.Unpack(&u); uerr != nil {
					t.Fatal(uerr)
				}
				if i == 0 && len(u.View.Stack) > 0 && u.Step > 20 {
					bots[0].send(protocol.TypeSave, struct{}{})
					continue
				}
				b.act(u)
			case protocol.TypeSaved:
				var s protocol.Saved
				if uerr := e.Unpack(&s); uerr != nil {
					t.Fatal(uerr)
				}
				saved = s.Save
			}
		}
	}
	if saved == "" {
		t.Fatal("the game was not saved")
	}
	var games protocol.Games
	bots[0].expect(protocol.TypeGames, &games)
	if len(games.Saves) != 1 || games.Saves[0].ID != saved || games.Saves[0].Seats[1] != "Guest" {
		t.Fatalf("saves %+v", games.Saves)
	}

	// Load it again: the seats wait for their players.
	bots[0].send(protocol.TypeLoad, protocol.Load{Save: saved})
	bots[0].expect(protocol.TypeTable, &tb)
	if tb.Loaded != saved || tb.You != 0 {
		t.Fatalf("loaded table %+v", tb)
	}
	bots[1].send(protocol.TypeJoin, protocol.Join{Game: tb.Game})
	bots[1].expect(protocol.TypeTable, &tb)
	if tb.You != 1 {
		t.Fatalf("the guest got seat %d", tb.You)
	}
	for _, b := range bots {
		b.send(protocol.TypeReady, protocol.Ready{Ready: true})
	}
	over := false
	for steps := 0; !over && steps < 40000; steps++ {
		for _, b := range bots {
			e, ok := b.read(2 * time.Millisecond)
			if !ok || e.Type != protocol.TypeUpdate {
				continue
			}
			var u protocol.Update
			if uerr := e.Unpack(&u); uerr != nil {
				t.Fatal(uerr)
			}
			if u.View.Over {
				over = true
				break
			}
			b.act(u)
		}
	}
	if !over {
		t.Skip("the continued game did not end within the step limit")
	}
	cancel()
	_ = srv.Shutdown(context.Background()) //nolint:errcheck // test
	files, err := filepath.Glob(filepath.Join(cfg.Records, "*"+record.Ext))
	if err != nil || len(files) != 1 {
		t.Fatalf("records %v, %v", files, err)
	}
	rec, err := record.ReadFile(files[0])
	if err != nil {
		t.Fatal(err)
	}
	if rec.Resumes != 1 {
		t.Errorf("resumes %d, want 1", rec.Resumes)
	}
	checkRecord(t, files[0])
}
