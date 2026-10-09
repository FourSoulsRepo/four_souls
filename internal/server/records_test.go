package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	"github.com/FourSoulsRepo/record"
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards"
)

// TestRecordReplays: a game's record replays in the engine with every
// checksum matching; a stopped game is marked unfinished (6.9).
func TestRecordReplays(t *testing.T) {
	dir := t.TempDir()
	h := NewHub()
	h.records = dir
	h.cards = []record.Card{{Ref: "the_d6", Name: "The D6"}}
	ctx, cancel := context.WithCancel(context.Background())
	stopped := make(chan struct{})
	go func() { h.Run(ctx); close(stopped) }()
	players := startTable(t, h, 2)
	for range 30 { // whoever is asked answers with the first allowed intent
		for _, p := range players {
			up, ok := p.latest()
			if !ok || len(up.Allowed) == 0 || up.View.Over {
				continue
			}
			in := up.Allowed[0]
			if in.Kind == engine.IntentDiscard {
				in.Objects = in.Objects[:up.View.Waiting.Count]
			}
			p.say(protocol.TypeIntent, protocol.Intent{Intent: in})
		}
		time.Sleep(5 * time.Millisecond)
	}
	cancel()
	<-stopped

	files, err := filepath.Glob(filepath.Join(dir, "*"+record.Ext))
	if err != nil || len(files) != 1 {
		t.Fatalf("records %v, %v", files, err)
	}
	rec, err := record.ReadFile(files[0])
	if err != nil {
		t.Fatal(err)
	}
	if rec.End == nil || rec.End.Finished || len(rec.Steps) == 0 || rec.Header.Seats[0].Name != "P0" || len(rec.Header.Cards) != 1 {
		t.Fatalf("record: %d steps, end %+v, header %+v", len(rec.Steps), rec.End, rec.Header.Seats)
	}
	s := rec.Header.Setup
	setup := engine.Setup{Seed: s.Seed, Players: s.Players, BonusSouls: s.BonusSouls}
	for _, code := range s.Sets {
		set, ok := cards.Find(code)
		if !ok {
			t.Fatalf("unknown set %s", code)
		}
		setup.Sets = append(setup.Sets, set)
	}
	for _, c := range s.Characters {
		setup.Characters = append(setup.Characters, engine.CardRef(c))
	}
	var steps []engine.StepRecord
	for _, st := range rec.Steps {
		var in engine.Intent
		if err := json.Unmarshal(st.Intent, &in); err != nil {
			t.Fatal(err)
		}
		steps = append(steps, engine.StepRecord{Intent: in, Checksum: st.Checksum})
	}
	if _, err := engine.Replay(setup, steps); err != nil {
		t.Fatalf("the record does not replay: %v", err)
	}
}

// TestDownloadRecord: only the game's players, only after it is over
// (RP-05).
func TestDownloadRecord(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "g1"+record.Ext)
	if err := os.WriteFile(path, []byte("rec"), 0o600); err != nil {
		t.Fatal(err)
	}
	h := startHub(t)
	room := &Room{tokens: []string{"tok"}, path: path}
	done := make(chan struct{})
	h.send(hubRequest{run: func(h *Hub) {
		h.games = append(h.games, started{info: protocol.GameInfo{ID: "g1", Started: true}, room: room})
		close(done)
	}})
	<-done
	srv := httptest.NewServer(http.HandlerFunc(h.downloadRecord))
	defer srv.Close()
	get := func(q string) int {
		resp, err := http.Get(srv.URL + "/record?" + q) //nolint:noctx // test
		if err != nil {
			t.Fatal(err)
		}
		_ = resp.Body.Close() //nolint:errcheck // test
		return resp.StatusCode
	}
	if code := get("game=g1&token=tok"); code != http.StatusConflict {
		t.Errorf("before the end: %d", code)
	}
	room.over.Store(true)
	if code := get("game=g1&token=other"); code != http.StatusForbidden {
		t.Errorf("a stranger: %d", code)
	}
	if code := get("game=nope&token=tok"); code != http.StatusNotFound {
		t.Errorf("no game: %d", code)
	}
	if code := get("game=g1&token=tok"); code != http.StatusOK {
		t.Errorf("a player after the end: %d", code)
	}
}
