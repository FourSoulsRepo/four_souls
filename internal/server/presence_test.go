package server

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	"github.com/FourSoulsRepo/record"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// presences collects the presence messages that come within d.
func (p *player) presences(d time.Duration) []protocol.Presence {
	p.t.Helper()
	var out []protocol.Presence
	end := time.After(d)
	for {
		select {
		case msg := <-p.conn.out:
			e, err := protocol.Decode(msg)
			if err != nil {
				p.t.Fatal(err)
			}
			if e.Type != protocol.TypePresence {
				continue
			}
			if bytes.Contains(e.Data, []byte(`"object"`)) && bytes.Contains(e.Data, []byte(`"hand"`)) {
				p.t.Errorf("a presence names a hand slot and an object: %s", e.Data)
			}
			var pr protocol.Presence
			if err := e.Unpack(&pr); err != nil {
				p.t.Fatal(err)
			}
			out = append(out, pr)
		case <-end:
			return out
		}
	}
}

// playing starts a game of n players and waits for everyone's first
// update.
func playing(t *testing.T, h *Hub, n int) ([]*player, []protocol.Update) {
	t.Helper()
	players := startTable(t, h, n)
	ups := make([]protocol.Update, n)
	for i, p := range players {
		p.next(protocol.TypeUpdate, &ups[i])
	}
	return players, ups
}

// TestPresenceReachesOthers: an aim reaches the other players with the
// sender's seat, never the sender (15.35).
func TestPresenceReachesOthers(t *testing.T) {
	h := startHub(t)
	if _, w := arrive(t, h, "Ann", ""); !w.Presence {
		t.Error("welcome says presence is off")
	}
	players, _ := playing(t, h, 3)
	target := engine.PlayerID(2)
	deck := engine.MonsterDeck
	players[1].say(protocol.TypePresence, protocol.Presence{Seat: 0, At: &protocol.Spot{Player: &target}, Point: &protocol.Point{X: 2, Y: 0.5}})
	players[1].say(protocol.TypePresence, protocol.Presence{At: &protocol.Spot{Deck: &deck}})
	for _, i := range []int{0, 2} {
		got := players[i].presences(200 * time.Millisecond)
		if len(got) != 2 {
			t.Fatalf("player %d got %d presences", i, len(got))
		}
		first := got[0]
		if first.Seat != 1 || first.At == nil || first.At.Player == nil || *first.At.Player != target || first.Point == nil || first.Point.X != 1 {
			t.Errorf("player %d got %+v", i, first)
		}
		if got[1].At == nil || got[1].At.Deck == nil || *got[1].At.Deck != deck {
			t.Errorf("player %d got %+v", i, got[1])
		}
	}
	if got := players[1].presences(100 * time.Millisecond); len(got) != 0 {
		t.Errorf("the sender got its own presence: %+v", got)
	}
}

// TestPresenceHidesHandCards: a hidden card aims as its hand slot; its
// object never reaches the others (A-16).
func TestPresenceHidesHandCards(t *testing.T) {
	players, ups := playing(t, startHub(t), 2)
	hand := ups[0].View.Players[0].Hand
	if len(hand) == 0 {
		t.Fatal("no hand card")
	}
	id := hand[0].ID
	players[0].say(protocol.TypePresence, protocol.Presence{From: &protocol.Spot{Object: id}})
	players[0].say(protocol.TypePresence, protocol.Presence{From: &protocol.Spot{Hand: 1, Object: id}})
	players[0].say(protocol.TypePresence, protocol.Presence{From: &protocol.Spot{Hand: 1}})
	players[0].say(protocol.TypePresence, protocol.Presence{From: &protocol.Spot{Hand: 99}})
	got := players[1].presences(200 * time.Millisecond)
	if len(got) != 4 {
		t.Fatalf("got %d presences", len(got))
	}
	for i, want := range []int{0, 0, 1, 0} {
		from := got[i].From
		switch {
		case want == 0 && from != nil:
			t.Errorf("presence %d kept %+v", i, *from)
		case want != 0 && (from == nil || from.Hand != want || from.Object != 0):
			t.Errorf("presence %d: %+v, want hand slot %d", i, from, want)
		}
	}
	character := ups[0].View.Players[0].Character.ID
	players[0].say(protocol.TypePresence, protocol.Presence{From: &protocol.Spot{Object: character}})
	if got := players[1].presences(200 * time.Millisecond); len(got) != 1 || got[0].From == nil || got[0].From.Object != character {
		t.Errorf("a public card does not aim: %+v", got)
	}
}

// TestPresencePace: a player's flood is cut to MaxPresence a second.
func TestPresencePace(t *testing.T) {
	players, _ := playing(t, startHub(t), 2)
	for range 100 {
		players[0].say(protocol.TypePresence, protocol.Presence{Point: &protocol.Point{X: 0.5, Y: 0.5}})
	}
	got := players[1].presences(300 * time.Millisecond)
	if len(got) < protocol.MaxPresence || len(got) > protocol.MaxPresence+10 {
		t.Errorf("%d of 100 presences went through", len(got))
	}
}

// TestPresenceGoesOnDrop: a player who leaves takes their arrow along.
func TestPresenceGoesOnDrop(t *testing.T) {
	players, _ := playing(t, startHub(t), 2)
	players[0].c.Leave()
	got := players[1].presences(200 * time.Millisecond)
	if len(got) != 1 || got[0].Seat != 0 || got[0].From != nil || got[0].At != nil || got[0].Point != nil {
		t.Errorf("got %+v, want an empty presence of seat 0", got)
	}
}

// TestPresenceNotRecorded: the record has no presence messages (15.35).
func TestPresenceNotRecorded(t *testing.T) {
	dir := t.TempDir()
	h := NewHub()
	h.records = dir
	ctx, cancel := context.WithCancel(context.Background())
	stopped := make(chan struct{})
	go func() { h.Run(ctx); close(stopped) }()
	players, ups := playing(t, h, 2)
	players[0].say(protocol.TypePresence, protocol.Presence{Point: &protocol.Point{X: 0.123456, Y: 0.654321}})
	if got := players[1].presences(200 * time.Millisecond); len(got) != 1 {
		t.Fatalf("got %d presences", len(got))
	}
	for i, p := range players {
		if len(ups[i].Allowed) > 0 {
			p.say(protocol.TypeIntent, protocol.Intent{Intent: ups[i].Allowed[0]})
		}
	}
	time.Sleep(50 * time.Millisecond)
	cancel()
	<-stopped

	files, err := filepath.Glob(filepath.Join(dir, "*"+record.Ext))
	if err != nil || len(files) != 1 {
		t.Fatalf("records %v, %v", files, err)
	}
	f, err := os.Open(files[0])
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close() //nolint:errcheck // read only
	zr, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	text, err := io.ReadAll(zr)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range []string{protocol.TypePresence, "0.123456", "0.654321"} {
		if bytes.Contains(text, []byte(s)) {
			t.Errorf("the record holds %q", s)
		}
	}
	rec, err := record.ReadFile(files[0])
	if err != nil {
		t.Fatal(err)
	}
	if len(rec.Steps) == 0 {
		t.Error("the record keeps no step; the check above proves nothing")
	}
}

// TestPresenceOff: a server with presence off says so and relays none.
func TestPresenceOff(t *testing.T) {
	h := NewHub()
	h.presence = false
	ctx, cancel := context.WithCancel(context.Background())
	go h.Run(ctx)
	t.Cleanup(cancel)
	if _, w := arrive(t, h, "Ann", ""); w.Presence {
		t.Error("welcome says presence is on")
	}
	players, _ := playing(t, h, 2)
	players[0].say(protocol.TypePresence, protocol.Presence{Point: &protocol.Point{X: 0.5, Y: 0.5}})
	players[0].c.Leave()
	if got := players[1].presences(200 * time.Millisecond); len(got) != 0 {
		t.Errorf("got %+v with presence off", got)
	}
}
