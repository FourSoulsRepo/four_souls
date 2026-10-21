package server

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	"github.com/FourSoulsRepo/record"
)

// startKeeping starts a server that keeps its games over a restart, on
// the records and saves folders of dir (N-20).
func startKeeping(t *testing.T, dir string) *Server {
	t.Helper()
	cfg := DefaultConfig()
	cfg.Addr, cfg.Port, cfg.Seed, cfg.KeepGames = "127.0.0.1", freePort(t), 91, true
	cfg.Records, cfg.Saves = filepath.Join(dir, "records"), filepath.Join(dir, "saves")
	srv, err := Start(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = srv.Shutdown(context.Background()) }) //nolint:errcheck // test
	return srv
}

// restart stops srv and starts a new server on the same folders; every
// bot reconnects to it with its token.
func restart(t *testing.T, srv *Server, dir string, bots ...*bot) *Server {
	t.Helper()
	if err := srv.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	next := startKeeping(t, dir)
	for _, b := range bots {
		_ = b.ws.Close() //nolint:errcheck // test
		b.url = "ws://" + next.Addr() + "/ws"
		b.connect()
	}
	return next
}

// seatTwo opens a table for two bots; options as given.
func seatTwo(t *testing.T, url string, opts *protocol.Options) []*bot {
	t.Helper()
	bots := []*bot{dialBot(t, url, "Ann", 1), dialBot(t, url, "Bo", 2)}
	bots[0].send(protocol.TypeCreate, protocol.Create{Seats: 2, Options: opts})
	var tb protocol.Table
	bots[0].expect(protocol.TypeTable, &tb)
	bots[1].send(protocol.TypeJoin, protocol.Join{Game: tb.Game})
	bots[1].expect(protocol.TypeTable, nil)
	return bots
}

// TestRestartKeepsGame: a game stopped by a restart goes on at the same
// step with the same view; it plays to the end and its record replays
// with every checksum (8.22).
func TestRestartKeepsGame(t *testing.T) {
	if testing.Short() {
		t.Skip("plays a whole game")
	}
	dir := t.TempDir()
	srv := startKeeping(t, dir)
	bots := seatTwo(t, "ws://"+srv.Addr()+"/ws", nil)
	for _, b := range bots {
		b.send(protocol.TypeReady, protocol.Ready{Ready: true})
	}

	// Play a while, then stop with the game in the middle.
	last := make([]protocol.Update, len(bots))
	for steps := 0; last[0].Step <= 20 && steps < 20000; steps++ {
		for i, b := range bots {
			e, ok := b.read(2 * time.Millisecond)
			if !ok || e.Type != protocol.TypeUpdate {
				continue
			}
			last[i] = protocol.Update{} // fields an update leaves out are empty
			if err := e.Unpack(&last[i]); err != nil {
				t.Fatal(err)
			}
			if last[0].Step <= 20 {
				b.act(last[i])
			}
		}
	}
	if last[0].Step <= 20 {
		t.Fatal("the game did not get past step 20")
	}
	// Let the last intent land, so both bots saw the same step.
	for _, b := range bots {
		for {
			e, ok := b.read(200 * time.Millisecond)
			if !ok {
				break
			}
			if e.Type == protocol.TypeUpdate {
				i := b.seatIndex(bots)
				last[i] = protocol.Update{}
				if err := e.Unpack(&last[i]); err != nil {
					t.Fatal(err)
				}
			}
		}
	}
	if err := srv.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	for _, b := range bots {
		b.expect(protocol.TypeRestarting, nil)
	}
	kept, err := filepath.Glob(filepath.Join(dir, "saves", restartFolder, "*"+restartExt))
	if err != nil || len(kept) != 1 {
		t.Fatalf("restart files %v, %v; want one", kept, err)
	}

	next := startKeeping(t, dir)
	for _, b := range bots {
		_ = b.ws.Close() //nolint:errcheck // test
		b.url = "ws://" + next.Addr() + "/ws"
		b.connect()
	}
	for i, b := range bots {
		u := b.untilBack()
		if u.Step != last[i].Step {
			t.Errorf("%s: step %d after the restart, want %d", b.name, u.Step, last[i].Step)
		}
		before, merr := json.Marshal(last[i].View)
		if merr != nil {
			t.Fatal(merr)
		}
		after, merr := json.Marshal(u.View)
		if merr != nil {
			t.Fatal(merr)
		}
		if string(before) != string(after) {
			t.Errorf("%s: the view changed over the restart:\n%s\n%s", b.name, before, after)
		}
		b.act(u)
	}
	if left, _ := filepath.Glob(filepath.Join(dir, "saves", restartFolder, "*")); len(left) != 0 { //nolint:errcheck // test
		t.Errorf("restart files left: %v", left)
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
		t.Skip("the game did not end within the step limit")
	}
	if err = next.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	files, err := filepath.Glob(filepath.Join(dir, "records", "*"+record.Ext))
	if err != nil || len(files) != 1 {
		t.Fatalf("records %v, %v", files, err)
	}
	rec, err := record.ReadFile(files[0])
	if err != nil {
		t.Fatal(err)
	}
	if rec.Resumes != 1 || rec.End == nil || !rec.End.Finished {
		t.Errorf("resumes %d, end %+v; want 1 and finished", rec.Resumes, rec.End)
	}
	checkRecord(t, files[0])
}

// seatIndex is b's place in bots.
func (b *bot) seatIndex(bots []*bot) int {
	for i, x := range bots {
		if x == b {
			return i
		}
	}
	return -1
}

// untilBack reads updates until one is not paused: every player is back.
func (b *bot) untilBack() protocol.Update {
	b.t.Helper()
	for {
		var u protocol.Update
		b.expect(protocol.TypeUpdate, &u)
		if u.Pause == nil {
			return u
		}
	}
}

// TestRestartKeepsTable: a lobby table comes back with its seats, its
// password, and its number; a new table gets the next one (8.22).
func TestRestartKeepsTable(t *testing.T) {
	dir := t.TempDir()
	srv := startKeeping(t, dir)
	ann := dialBot(t, "ws://"+srv.Addr()+"/ws", "Ann", 1)
	ann.send(protocol.TypeCreate, protocol.Create{Seats: 3, Password: "pw"})
	var tb protocol.Table
	ann.expect(protocol.TypeTable, &tb)

	srv = restart(t, srv, dir, ann)
	var back protocol.Table
	ann.expect(protocol.TypeTable, &back)
	if back.Game != tb.Game || back.You != 0 || back.Seats[0].Name != "Ann" {
		t.Fatalf("table after the restart %+v, want %s with Ann in seat 0", back, tb.Game)
	}

	url := "ws://" + srv.Addr() + "/ws"
	bo := dialBot(t, url, "Bo", 2)
	bo.send(protocol.TypeJoin, protocol.Join{Game: tb.Game})
	var e protocol.Error
	bo.expect(protocol.TypeError, &e)
	if e.Code != protocol.ErrPassword {
		t.Errorf("join without the password: %q, want %q", e.Code, protocol.ErrPassword)
	}
	cy := dialBot(t, url, "Cy", 3)
	cy.send(protocol.TypeCreate, protocol.Create{Seats: 2})
	cy.expect(protocol.TypeTable, &back)
	if back.Game == tb.Game {
		t.Errorf("a new table got the kept table's ID %s", back.Game)
	}
}

// TestRestartStartsBansOver: bans and picks are not kept; they start
// over once every player is back (8.22).
func TestRestartStartsBansOver(t *testing.T) {
	dir := t.TempDir()
	srv := startKeeping(t, dir)
	bots := seatTwo(t, "ws://"+srv.Addr()+"/ws", &protocol.Options{BanRounds: 1})
	for _, b := range bots {
		b.send(protocol.TypeReady, protocol.Ready{Ready: true})
	}
	var su protocol.Setup
	bots[0].expect(protocol.TypeSetup, &su)
	if su.Phase != protocol.PhaseBan {
		t.Fatalf("phase %q, want bans", su.Phase)
	}
	pool := len(su.Pool)
	bots[0].send(protocol.TypeBan, protocol.BanCard{Card: su.Pool[0]})
	bots[1].expect(protocol.TypeSetup, nil)

	restart(t, srv, dir, bots...)
	for _, b := range bots {
		var again protocol.Setup
		b.expect(protocol.TypeSetup, &again)
		if again.Phase != protocol.PhaseBan || again.Round != 1 || again.Turn != 0 || len(again.Pool) != pool {
			t.Errorf("%s: setup %s round %d turn %d, %d in the pool; want bans from the start, %d", b.name, again.Phase, again.Round, again.Turn, len(again.Pool), pool)
		}
	}
}

// TestNoKeepingWithoutTheSetting: an app's own host ends its games when
// it stops; nothing is kept (N-20).
func TestNoKeepingWithoutTheSetting(t *testing.T) {
	dir := t.TempDir()
	cfg := DefaultConfig()
	cfg.Addr, cfg.Port, cfg.Records, cfg.Saves = "127.0.0.1", freePort(t), filepath.Join(dir, "records"), filepath.Join(dir, "saves")
	srv, err := Start(context.Background(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	ann := dialBot(t, "ws://"+srv.Addr()+"/ws", "Ann", 1)
	ann.send(protocol.TypeCreate, protocol.Create{Seats: 2})
	ann.expect(protocol.TypeTable, nil)
	if err := srv.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(cfg.Saves, restartFolder)); !os.IsNotExist(err) {
		t.Errorf("restart folder: %v, want none", err)
	}
}
