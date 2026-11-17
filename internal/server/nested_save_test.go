package server

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/b2"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// TestNestedAskSaveRestartRollback: a game saved while Mom's Shovel's
// question waits inside DestroyThis goes on in every path that loads a
// state (19.177): the host's save, the restart keep, and a Godmode
// rollback from that base. Each save keeps the last clean state and
// the steps since it; each load replays them.
func TestNestedAskSaveRestartRollback(t *testing.T) {
	dir := t.TempDir()
	saves := filepath.Join(dir, "saves")
	built := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"moms_shovel"}},
		{Character: "cain", Souls: []engine.CardRef{"monstro"}},
	}}, b2.Set)
	built.Start(0, "moms_shovel", 0)
	built.UntilChoose()
	if !built.G.NestedAsk() {
		t.Fatal("Mom's Shovel's question is not nested in DestroyThis")
	}
	state, err := built.G.Save()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(state), `"replay":`) {
		t.Fatal("the save at a nested ask has no steps to replay")
	}
	f := saveFile{
		Format: saveFormat, Engine: engine.Version, Game: "g1", Sets: []string{"b2"},
		Options: protocol.Options{Rollback: protocol.RollbackGod},
		Seats:   []savedSeat{{Name: "Ann", Token: newToken()}, {Name: "Bob", Token: newToken()}}, State: state,
	}
	data, err := json.Marshal(f)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(saves, 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(saves, "g1"+saveExt), data, 0o600); err != nil {
		t.Fatal(err)
	}

	// The host's save loads: the question is Mom's Shovel's.
	h, stop := keepingHubWithSaves(t, dir, saves)
	ann, _ := arrive(t, h, "Ann", "")
	ann.say(protocol.TypeLoad, protocol.Load{Save: "g1"})
	ann.next(protocol.TypeTable, nil)
	bob, _ := arrive(t, h, "Bob", "")
	bob.say(protocol.TypeJoin, protocol.Join{Game: "g1"})
	bob.next(protocol.TypeTable, nil)
	ann.say(protocol.TypeReady, protocol.Ready{Ready: true})
	bob.say(protocol.TypeReady, protocol.Ready{Ready: true})
	tb := &driver{t: t, players: []*player{ann, bob}, last: make([]protocol.Update, 2)}
	tb.updates(-1)
	asked := viewJSON(t, tb.last[0])
	if w := tb.last[0].View.Waiting; w.Kind != engine.PromptChoose || w.Player != 0 {
		t.Fatalf("loaded: prompt %+v, want Ann's question", w)
	}

	// The restart keeps it; the game comes back at the same question.
	tokens := tb.stopKept(stop)
	if _, kept := keptFile(t, dir); !strings.Contains(string(kept.State), `"replay":`) {
		t.Fatal("the kept state has no steps to replay")
	}
	back, stop := comeBack(t, dir, tokens)
	defer stop()
	if viewJSON(t, back.last[0]) != asked {
		t.Fatal("the question changed over the restart")
	}
	if u := back.last[0]; u.RollbackCut != 0 {
		t.Fatalf("the restart cut the history at %d", u.RollbackCut)
	}

	answer := engine.Intent{Player: 0, Kind: engine.IntentChoose, Choice: 0}
	k := back.do(answer)
	if got := back.last[0].View.Players[0].Souls; got != 1 {
		t.Fatalf("Ann has %d souls after the steal, want 1", got)
	}

	// A Godmode rollback replays from the base saved at the question.
	back.players[1].say(protocol.TypeRollback, protocol.RollbackTo{Step: k})
	got := back.updates(-1)
	if got[0].Rollback == nil || got[0].Rollback.To != k {
		t.Fatalf("rollback %+v, want to before %d", got[0].Rollback, k)
	}
	if viewJSON(t, got[0]) != asked {
		t.Fatal("the rolled back table is not the question")
	}
	back.do(answer)
	if got := back.last[0].View.Players[0].Souls; got != 1 {
		t.Fatalf("Ann has %d souls after the steal again, want 1", got)
	}
}

// keepingHubWithSaves is keepingHub with a saves folder too.
func keepingHubWithSaves(t *testing.T, dir, saves string) (*Hub, func()) {
	t.Helper()
	h := NewHub()
	h.restart, h.records, h.saves = filepath.Join(dir, restartFolder), filepath.Join(dir, "records"), saves
	ctx, cancel := context.WithCancel(context.Background())
	stopped := make(chan struct{})
	go func() { h.Run(ctx); close(stopped) }()
	var once sync.Once
	stop := func() { once.Do(func() { cancel(); <-stopped }) }
	t.Cleanup(stop)
	return h, stop
}
