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
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards"
)

// driver drives a game one intent at a time: each intent gets one
// update per player. It keeps every update player 0 got.
type driver struct {
	t       *testing.T
	players []*player
	last    []protocol.Update // each player's latest update
	seen    []protocol.Update // player 0's updates, in order
	authors []int             // who acted for each of seen; -1: nobody
	sent    []protocol.Intent // what they sent, by the same index
}

// updates reads one update for every player.
func (tb *driver) updates(author int) []protocol.Update {
	tb.t.Helper()
	for i, p := range tb.players {
		var u protocol.Update // fresh: decoding reuses slices
		p.next(protocol.TypeUpdate, &u)
		tb.last[i] = u
	}
	tb.seen = append(tb.seen, tb.last[0])
	tb.authors = append(tb.authors, author)
	tb.sent = append(tb.sent, protocol.Intent{})
	return tb.last
}

// act makes the asked player take an attack if they can, else their
// first allowed intent; it returns that player's seat.
func (tb *driver) act() int {
	tb.t.Helper()
	for i, u := range tb.last {
		if len(u.Allowed) == 0 || u.View.Over {
			continue
		}
		in := u.Allowed[0]
		for _, a := range u.Allowed {
			if a.Kind == engine.IntentAttack {
				in = a // attacks roll: the rollback must repeat it
			}
		}
		if in.Kind == engine.IntentDiscard {
			in.Objects = in.Objects[:u.View.Waiting.Count]
		}
		tb.players[i].say(protocol.TypeIntent, protocol.Intent{Intent: in})
		tb.updates(i)
		tb.sent[len(tb.sent)-1] = protocol.Intent{Intent: in}
		return i
	}
	tb.t.Fatal("nobody may act")
	return -1
}

// startGodmode starts a 2-player Godmode game with records in dir.
func startGodmode(t *testing.T, dir string) (*driver, context.CancelFunc, chan struct{}) {
	t.Helper()
	h := NewHub()
	h.records = dir
	ctx, cancel := context.WithCancel(context.Background())
	stopped := make(chan struct{})
	go func() { h.Run(ctx); close(stopped) }()
	players := startTableWith(t, h, protocol.Create{Seats: 2, Options: &protocol.Options{Rollback: protocol.RollbackGod}})
	tb := &driver{t: t, players: players, last: make([]protocol.Update, 2)}
	tb.updates(-1) // the start
	return tb, cancel, stopped
}

func viewJSON(t *testing.T, u protocol.Update) string {
	t.Helper()
	b, err := json.Marshal(struct {
		V engine.GameView `json:"v"`
		A []engine.Intent `json:"a"`
	}{u.View, u.Allowed})
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// TestGodmodeRollback: in Godmode the author rolls back their own step,
// then another player rolls back an earlier one. Each time the table is
// the one from before the target, and the step number goes on after
// the highest. The same intent again repeats the same events, dice
// included (ADR 020, items 6 to 9). The record keeps every step and
// both rollbacks, and replays without the undone ones.
func TestGodmodeRollback(t *testing.T) {
	dir := t.TempDir()
	tb, cancel, stopped := startGodmode(t, dir)
	if u := tb.last[0]; u.RollbackMode != protocol.RollbackGod || u.RollbackFrom != 1 || u.Rollback != nil {
		t.Fatalf("start: mode %q from %d rollback %+v", u.RollbackMode, u.RollbackFrom, u.Rollback)
	}
	for range 12 {
		tb.act()
	}
	// The target: the first step of the 8th update; its author rolls back.
	j := 8
	before, target, author := tb.seen[j-1], tb.seen[j], tb.authors[j]
	k, highest := before.Step+1, tb.last[0].Step
	tb.players[author].say(protocol.TypeRollback, protocol.RollbackTo{Step: k})
	got := tb.updates(-1)
	for i, u := range got {
		if u.Rollback == nil || *u.Rollback != (protocol.Rollback{Seat: author, To: k, Undid: highest - k + 1}) {
			t.Fatalf("player %d: rollback %+v, want seat %d to %d undid %d", i, u.Rollback, author, k, highest-k+1)
		}
		if u.Step != highest+1 || len(u.Events) != 0 {
			t.Errorf("player %d: step %d (want %d), %d events", i, u.Step, highest+1, len(u.Events))
		}
	}
	if viewJSON(t, got[0]) != viewJSON(t, before) {
		t.Error("the restored table differs from the one before the target")
		for i, u := range tb.seen {
			t.Logf("%d: step %d by %d: %s", i, u.Step, tb.authors[i], eventsJSON(t, u.Events))
		}
	}
	// The undone step is no target any more; nor the start.
	other := tb.players[1-author]
	for _, bad := range []int{k + 1, 0} {
		other.say(protocol.TypeRollback, protocol.RollbackTo{Step: bad})
		if code := other.failure(); code != protocol.ErrRollback {
			t.Errorf("step %d: %s, want %s", bad, code, protocol.ErrRollback)
		}
	}
	// The same intent again: the same events, the same roll; its step
	// goes on after the rollback's number.
	tb.players[author].say(protocol.TypeIntent, tb.sent[j])
	redo := tb.updates(author)[0]
	if a, b := eventsJSON(t, redo.Events), eventsJSON(t, target.Events); a != b {
		t.Errorf("the same intent after the rollback did otherwise:\n%s\nwas\n%s", a, b)
	}
	if redo.Step != highest+1+(target.Step-before.Step) {
		t.Errorf("the step after the rollback is %d, want %d", redo.Step, highest+1+(target.Step-before.Step))
	}
	for range 4 {
		tb.act()
	}
	// Another player rolls back further, past the first rollback: the
	// live steps from k2 to k-1, and those after the first rollback.
	before2 := tb.seen[3]
	k2, last := before2.Step+1, tb.last[0].Step
	second := 1 - tb.authors[4]
	tb.players[second].say(protocol.TypeRollback, protocol.RollbackTo{Step: k2})
	got = tb.updates(-1)
	want := protocol.Rollback{Seat: second, To: k2, Undid: k - k2 + last - (highest + 1)}
	if got[0].Rollback == nil || *got[0].Rollback != want {
		t.Fatalf("second rollback %+v, want %+v", got[0].Rollback, want)
	}
	if viewJSON(t, got[0]) != viewJSON(t, before2) {
		t.Error("the second restored table differs from the one before its target")
	}
	cancel()
	<-stopped
	checkRollbackRecord(t, dir, 2)
}

// TestRollbackRefused: no rollback in a game without the mode, nor
// for a seat that does not play (ADR 020, item 2); no unknown mode.
func TestRollbackRefused(t *testing.T) {
	h := startHub(t)
	players := startTable(t, h, 2)
	var u protocol.Update
	players[0].next(protocol.TypeUpdate, &u)
	if u.RollbackMode != "" || u.RollbackFrom != 0 {
		t.Errorf("a Disabled game offers rollback: %q from %d", u.RollbackMode, u.RollbackFrom)
	}
	players[1].say(protocol.TypeRollback, protocol.RollbackTo{Step: 1})
	if code := players[1].failure(); code != protocol.ErrRollback {
		t.Errorf("Disabled: %s", code)
	}
	host, _ := arrive(t, h, "Chaos", "")
	host.say(protocol.TypeCreate, protocol.Create{Seats: 2, Options: &protocol.Options{Rollback: "chaos"}})
	if code := host.failure(); code != protocol.ErrBadSetup {
		t.Errorf("an unknown mode: %s", code)
	}
}

func eventsJSON(t *testing.T, events []engine.Event) string {
	t.Helper()
	b, err := json.Marshal(events)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

// checkRollbackRecord: the record keeps its rollbacks, and its live
// steps replay with every checksum.
func checkRollbackRecord(t *testing.T, dir string, rollbacks int) {
	t.Helper()
	files, err := filepath.Glob(filepath.Join(dir, "*"+record.Ext))
	if err != nil || len(files) != 1 {
		t.Fatalf("records %v, %v", files, err)
	}
	rec, err := record.ReadFile(files[0])
	if err != nil {
		t.Fatal(err)
	}
	if rec.Header.Rollback != protocol.RollbackGod || len(rec.Rollbacks) != rollbacks {
		t.Fatalf("record: mode %q, %d rollbacks", rec.Header.Rollback, len(rec.Rollbacks))
	}
	s := rec.Header.Setup
	setup := engine.Setup{Seed: s.Seed, Players: s.Players, BonusSouls: s.BonusSouls, DeckRatios: s.DeckRatios}
	for _, code := range s.Sets {
		set, _ := cards.Find(code)
		setup.Sets = append(setup.Sets, set)
	}
	for _, c := range s.Characters {
		setup.Characters = append(setup.Characters, engine.CardRef(c))
	}
	for _, c := range s.Banned {
		setup.Banned = append(setup.Banned, engine.CardRef(c))
	}
	live := rec.Live()
	if len(live) >= len(rec.Steps) {
		t.Errorf("%d live of %d steps: nothing undone", len(live), len(rec.Steps))
	}
	steps := make([]engine.StepRecord, len(live))
	for i, st := range live {
		if err := json.Unmarshal(st.Intent, &steps[i].Intent); err != nil {
			t.Fatal(err)
		}
		steps[i].Checksum = st.Checksum
	}
	if _, err := engine.Replay(setup, steps); err != nil {
		t.Fatalf("the live steps do not replay: %v", err)
	}
}

// TestRestartOldFileIsRollbackBase: a Godmode game from a restart file
// written before 17.5, without the history, keeps its mode; the restart
// is its new base, so nothing before it can be rolled back, and every
// update says where the history was cut (ADR 020, item 10).
func TestRestartOldFileIsRollbackBase(t *testing.T) {
	dir := t.TempDir()
	bots, path := stoppedGameWith(t, dir, &protocol.Options{Rollback: protocol.RollbackGod})
	data, err := os.ReadFile(path) //nolint:gosec // test
	if err != nil {
		t.Fatal(err)
	}
	var f restartFile
	if err := json.Unmarshal(data, &f); err != nil {
		t.Fatal(err)
	}
	if f.Options.Rollback != protocol.RollbackGod || f.Step == 0 || len(f.Base) == 0 {
		t.Fatalf("restart file: mode %q, step %d, base %d bytes", f.Options.Rollback, f.Step, len(f.Base))
	}
	editRestart(t, path, func(data []byte) []byte {
		var m map[string]json.RawMessage
		if err := json.Unmarshal(data, &m); err != nil {
			t.Fatal(err)
		}
		delete(m, "base")
		delete(m, "base_step")
		delete(m, "live")
		out, err := json.Marshal(m)
		if err != nil {
			t.Fatal(err)
		}
		return out
	})
	reconnect(startKeeping(t, dir), bots...)
	for _, b := range bots {
		if u := b.untilBack(); u.RollbackMode != protocol.RollbackGod || u.RollbackFrom != f.Step+1 || u.RollbackCut != f.Step {
			t.Errorf("%s: mode %q from %d cut %d, want from %d cut %d", b.name, u.RollbackMode, u.RollbackFrom, u.RollbackCut, f.Step+1, f.Step)
		}
	}
	bots[1].send(protocol.TypeRollback, protocol.RollbackTo{Step: f.Step})
	for {
		e, ok := bots[1].read(5 * time.Second)
		if !ok {
			t.Fatal("no answer to a rollback behind the base")
		}
		if e.Type != protocol.TypeError {
			continue
		}
		var er protocol.Error
		if err := e.Unpack(&er); err != nil || er.Code != protocol.ErrRollback {
			t.Errorf("a rollback behind the base: %+v, %v", er, err)
		}
		return
	}
}
