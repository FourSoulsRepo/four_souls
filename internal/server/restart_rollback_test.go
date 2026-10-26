package server

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"sync"
	"testing"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
)

// keepingHub runs a hub that keeps its games over a restart in dir's
// restart folder, with records in dir's records folder (N-20); stop
// ends it as a server stop does.
func keepingHub(t *testing.T, dir string) (*Hub, func()) {
	t.Helper()
	h := NewHub()
	h.restart, h.records = filepath.Join(dir, restartFolder), filepath.Join(dir, "records")
	ctx, cancel := context.WithCancel(context.Background())
	stopped := make(chan struct{})
	go func() { h.Run(ctx); close(stopped) }()
	var once sync.Once
	stop := func() { once.Do(func() { cancel(); <-stopped }) }
	t.Cleanup(stop)
	return h, stop
}

// startKept starts a 2-player game with opts on a keeping hub.
func startKept(t *testing.T, dir string, opts protocol.Options) (*driver, func()) {
	t.Helper()
	h, stop := keepingHub(t, dir)
	players := startTableWith(t, h, protocol.Create{Seats: 2, Options: &opts})
	tb := &driver{t: t, players: players, last: make([]protocol.Update, 2)}
	tb.updates(-1) // the start
	return tb, stop
}

// keptFile reads the one restart file kept in dir.
func keptFile(t *testing.T, dir string) (string, restartFile) {
	t.Helper()
	files, err := filepath.Glob(filepath.Join(dir, restartFolder, "*"+restartExt))
	if err != nil || len(files) != 1 {
		t.Fatalf("restart files %v, %v; want one", files, err)
	}
	data, err := os.ReadFile(files[0])
	if err != nil {
		t.Fatal(err)
	}
	var f restartFile
	if err := json.Unmarshal(data, &f); err != nil {
		t.Fatal(err)
	}
	return files[0], f
}

// stopKept stops the hub; its game goes to a restart file. The tokens
// are the seats', to come back with.
func (tb *driver) stopKept(stop func()) []string {
	tb.t.Helper()
	stop()
	tokens := make([]string, len(tb.players))
	for i, p := range tb.players {
		tokens[i] = p.c.token // the hub has stopped
	}
	return tokens
}

// comeBack starts a new hub on dir; every player takes their seat back
// by token. It returns a driver whose last updates are the first ones
// with every player back.
func comeBack(t *testing.T, dir string, tokens []string) (*driver, func()) {
	t.Helper()
	h, stop := keepingHub(t, dir)
	tb := &driver{t: t, last: make([]protocol.Update, len(tokens))}
	for i, token := range tokens {
		p, _ := arrive(t, h, "P"+strconv.Itoa(i), token)
		tb.players = append(tb.players, p)
	}
	for i, p := range tb.players {
		for {
			var u protocol.Update
			p.next(protocol.TypeUpdate, &u)
			if u.Pause == nil {
				tb.last[i] = u
				break
			}
		}
	}
	tb.seen, tb.authors, tb.sent = []protocol.Update{tb.last[0]}, []int{-1}, []protocol.Intent{{}}
	return tb, stop
}

// TestRestartKeepsGodmodeHistory: with the same engine, the base and the
// live steps replay to the kept state, so a Godmode rollback reaches
// past the restart to the exact table from before (ADR 020, item 10).
// The record's live steps still replay from the start.
func TestRestartKeepsGodmodeHistory(t *testing.T) {
	dir := t.TempDir()
	tb, stop := startKept(t, dir, protocol.Options{Rollback: protocol.RollbackGod})
	for range 12 {
		tb.act()
	}
	before := tb.seen[7]
	k, view := before.Step+1, viewJSON(t, tb.last[0])
	tokens := tb.stopKept(stop)
	_, f := keptFile(t, dir)
	if len(f.Base) == 0 || len(f.Live) == 0 || f.Live[len(f.Live)-1].N != f.Steps {
		t.Fatalf("restart file: base %d bytes, %d live steps up to %d; want both, up to %d", len(f.Base), len(f.Live), f.Live[len(f.Live)-1].N, f.Steps)
	}

	back, stop := comeBack(t, dir, tokens)
	for i, u := range back.last {
		if u.RollbackFrom != 1 || u.RollbackCut != 0 {
			t.Errorf("player %d: from %d, cut %d; want the history from 1, no cut", i, u.RollbackFrom, u.RollbackCut)
		}
	}
	if viewJSON(t, back.last[0]) != view {
		t.Fatal("the table changed over the restart")
	}
	back.players[1].say(protocol.TypeRollback, protocol.RollbackTo{Step: k})
	got := back.updates(-1)
	if got[0].Rollback == nil || got[0].Rollback.To != k {
		t.Fatalf("rollback %+v, want to before %d", got[0].Rollback, k)
	}
	if viewJSON(t, got[0]) != viewJSON(t, before) {
		t.Error("the table rolled back over the restart differs from the one before the target")
	}
	stop()
	checkRollbackRecord(t, filepath.Join(dir, "records"), 1)
}

// TestRestartKeepsSoftHistory: Soft keeps its targets and barriers over
// the restart; a vote rolls back past it (ADR 020, item 10).
func TestRestartKeepsSoftHistory(t *testing.T) {
	dir := t.TempDir()
	tb, stop := startKept(t, dir, protocol.Options{Rollback: protocol.RollbackSoft})
	target := tb.toTarget()
	from, targets := tb.last[0].RollbackFrom, tb.last[0].RollbackSteps
	i := slices.IndexFunc(tb.seen, func(u protocol.Update) bool { return u.Step == target.Step-1 })
	if i < 0 {
		t.Fatalf("no update ends at step %d", target.Step-1)
	}
	before := tb.seen[i]
	tokens := tb.stopKept(stop)

	back, _ := comeBack(t, dir, tokens)
	u := back.last[0]
	if u.RollbackFrom != from || !slices.Equal(u.RollbackSteps, targets) || u.RollbackCut != 0 {
		t.Fatalf("after the restart: from %d, targets %v, cut %d; want %d, %v, none", u.RollbackFrom, u.RollbackSteps, u.RollbackCut, from, targets)
	}
	opener := 1 - target.Seat
	back.players[opener].say(protocol.TypeRollback, protocol.RollbackTo{Step: target.Step})
	back.rolledBack(target.Step, opener) // 2 players: the opener's yes decides
	if viewJSON(t, back.last[0]) != viewJSON(t, before) {
		t.Error("the table rolled back over the restart differs from the one before the target")
	}
}

// TestRestartCutsChangedHistory: a history that no longer replays to
// the kept state, or a base this engine cannot resume, ends at the
// restart: the kept state is the new base, the game goes on, and every
// update says where the history was cut (ADR 020, items 10.5 and 10.6).
func TestRestartCutsChangedHistory(t *testing.T) {
	edits := map[string]func(*restartFile){
		"no live steps": func(f *restartFile) { f.Live = nil }, // the base alone is not the kept state
		"a base that does not resume": func(f *restartFile) {
			f.Base = json.RawMessage(`{"sets":["nope"]}`)
		},
		"a step that is refused": func(f *restartFile) {
			f.Live[0].Intent = json.RawMessage(`{"player":1,"kind":99}`)
		},
	}
	for name, edit := range edits {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			tb, stop := startKept(t, dir, protocol.Options{Rollback: protocol.RollbackGod})
			for range 6 {
				tb.act()
			}
			tokens := tb.stopKept(stop)
			path, f := keptFile(t, dir)
			edit(&f)
			data, err := json.Marshal(&f) //nolint:gosec // G117: a test file; no password
			if err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, data, 0o600); err != nil {
				t.Fatal(err)
			}

			back, _ := comeBack(t, dir, tokens)
			for i, u := range back.last {
				if u.RollbackFrom != f.Step+1 || u.RollbackCut != f.Step || u.Step != f.Step {
					t.Errorf("player %d: step %d, from %d, cut %d; want the restart at %d as a barrier", i, u.Step, u.RollbackFrom, u.RollbackCut, f.Step)
				}
			}
			back.players[1].say(protocol.TypeRollback, protocol.RollbackTo{Step: f.Step})
			if code := back.players[1].failure(); code != protocol.ErrRollback {
				t.Errorf("a rollback behind the restart: %s", code)
			}
			// The game goes on, and a step after the restart rolls back.
			back.act()
			k := f.Step + 1
			back.players[1].say(protocol.TypeRollback, protocol.RollbackTo{Step: k})
			if got := back.updates(-1); got[0].Rollback == nil || got[0].RollbackCut != f.Step {
				t.Errorf("rollback after the restart %+v, cut %d", got[0].Rollback, got[0].RollbackCut)
			}
		})
	}
}

// TestRestartDisabledKeepsNoHistory: a game without a rollback mode
// writes no history into its restart file, and no cut.
func TestRestartDisabledKeepsNoHistory(t *testing.T) {
	dir := t.TempDir()
	tb, stop := startKept(t, dir, protocol.Options{})
	for range 4 {
		tb.act()
	}
	tokens := tb.stopKept(stop)
	_, f := keptFile(t, dir)
	if f.Base != nil || f.Live != nil {
		t.Errorf("a Disabled game kept a history: base %d bytes, %d live steps", len(f.Base), len(f.Live))
	}
	back, _ := comeBack(t, dir, tokens)
	if u := back.last[0]; u.RollbackMode != "" || u.RollbackCut != 0 {
		t.Errorf("mode %q, cut %d; want none", u.RollbackMode, u.RollbackCut)
	}
}
