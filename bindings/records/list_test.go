package records

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/msgerr"
	"github.com/FourSoulsRepo/four_souls/internal/replay"
	"github.com/FourSoulsRepo/record"
)

func header(game string, started time.Time) record.Header {
	return record.Header{
		App: "1.0.0", Engine: "1.0.0", Protocol: 1, Game: game, Started: started,
		Setup: record.Setup{Seed: 7, Players: 2, Sets: []string{"b2"}},
		Seats: []record.Seat{{Seat: 0, Name: "Ann"}, {Seat: 1, Name: "Bo"}},
		Cards: []record.Card{{Ref: "the_d6", Name: "The D6", Effects: []string{"Reroll a dice roll."}}, {Ref: "blank", Name: "Blank"}},
	}
}

// writeRecord writes a format 2 record of turns, two steps each; end
// nil leaves the game unfinished.
func writeRecord(t *testing.T, path string, h record.Header, turns int, end *record.End) {
	t.Helper()
	w, err := record.Create(path, h)
	if err != nil {
		t.Fatal(err)
	}
	for i := range 2 * turns {
		view := json.RawMessage(fmt.Sprintf(`{"turn":{"number":%d,"active":%d}}`, i/2+1, i/2%2))
		if serr := w.Step(map[string]int{"i": i}, []string{}, view, 1); serr != nil {
			t.Fatal(serr)
		}
	}
	if end == nil {
		err = w.Stop()
	} else {
		err = w.Close(*end)
	}
	if err != nil {
		t.Fatal(err)
	}
}

// writeFormat writes a record of format n by hand, with one step.
func writeFormat(t *testing.T, path string, n int) {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	h := header("old", time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	h.Kind, h.Format = "header", n
	enc := json.NewEncoder(gz)
	if err := enc.Encode(h); err != nil {
		t.Fatal(err)
	}
	if err := enc.Encode(record.Step{Kind: "step", N: 1, Intent: json.RawMessage(`{}`), Events: json.RawMessage(`[]`)}); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, buf.Bytes(), 0o600); err != nil {
		t.Fatal(err)
	}
}

// TestList: newest first, with players, winners and live steps; a
// format 1, a newer and a broken record say why they can't play.
func TestList(t *testing.T) {
	dir := t.TempDir()
	day := time.Date(2026, 11, 16, 12, 0, 0, 0, time.UTC)
	writeRecord(t, filepath.Join(dir, "won"+record.Ext), header("g1", day), 2, &record.End{Ended: day.Add(time.Hour), Finished: true, Winners: []int{1}})
	writeRecord(t, filepath.Join(dir, "going"+record.Ext), header("g2", day.Add(time.Hour)), 1, nil)
	writeFormat(t, filepath.Join(dir, "old"+record.Ext), 1)
	writeFormat(t, filepath.Join(dir, "new"+record.Ext), record.Format+1)
	for name, data := range map[string]string{"broken" + record.Ext: "hello", "notes.txt": "skip me"} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(data), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(dir, "dir"+record.Ext), 0o750); err != nil {
		t.Fatal(err)
	}
	list, err := New(background{}, dir).List()
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, e := range list {
		var names []string
		for _, p := range e.Players {
			names = append(names, p.Name)
		}
		got = append(got, fmt.Sprintf("%s %s %v %v %d %v", e.Name, e.Why, names, e.Finished, e.Steps, e.Winners))
	}
	want := []string{
		"going.fsrec  [Ann Bo] false 2 []",
		"won.fsrec  [Ann Bo] true 4 [1]",
		"old.fsrec before_replays [Ann Bo] false 1 []",
		"broken.fsrec unreadable [] false 0 []",
		"new.fsrec newer [] false 0 []",
	}
	if strings.Join(got, "\n") != strings.Join(want, "\n") {
		t.Errorf("list:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
	if list[1].Path != filepath.Join(dir, "won"+record.Ext) || list[1].Started != day.UnixMilli() {
		t.Errorf("entry %+v", list[1])
	}
	empty, err := New(background{}, filepath.Join(dir, "none")).List()
	if err != nil || empty == nil || len(empty) != 0 {
		t.Errorf("a missing folder: %v, %v", empty, err)
	}
}

// TestBotSeats: a bot seat keeps its level and number in the list and
// in an opened record, so the app names it "Bot N (level)" with its
// robot (19.216); a player's seat has neither.
func TestBotSeats(t *testing.T) {
	dir := t.TempDir()
	day := time.Date(2026, 11, 17, 12, 0, 0, 0, time.UTC)
	h := header("g1", day)
	h.Seats[1] = record.Seat{Seat: 1, Name: "Bot 1", Bot: "hard", BotN: 1}
	path := filepath.Join(dir, "bots"+record.Ext)
	writeRecord(t, path, h, 1, &record.End{Ended: day.Add(time.Hour), Finished: true, Winners: []int{1}})
	r := New(background{}, dir)
	list, err := r.List()
	if err != nil || len(list) != 1 {
		t.Fatalf("list %v, %v", list, err)
	}
	want := []Seat{{Seat: 0, Name: "Ann"}, {Seat: 1, Name: "Bot 1", Bot: "hard", BotN: 1}}
	if e := list[0]; fmt.Sprint(e.Players) != fmt.Sprint(want) || fmt.Sprint(e.Winners) != "[1]" {
		t.Errorf("entry players %+v, winners %v", e.Players, e.Winners)
	}
	s, err := r.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprint(s.Seats) != fmt.Sprint(want) {
		t.Errorf("summary seats %+v, want %+v", s.Seats, want)
	}
}

// TestOpen: a summary of seats, cards, steps, turns and the end; the
// binding keeps the record. Old, newer and broken records give the
// player's message.
func TestOpen(t *testing.T) {
	dir := t.TempDir()
	day := time.Date(2026, 11, 16, 12, 0, 0, 0, time.UTC)
	path := filepath.Join(dir, "won"+record.Ext)
	writeRecord(t, path, header("g1", day), 3, &record.End{Ended: day.Add(time.Hour), Finished: true, Winners: []int{1}})
	r := New(background{}, dir)
	s, err := r.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if s.Game != "g1" || s.Steps != 6 || len(s.Seats) != 2 || s.Seats[1].Name != "Bo" || s.Cards[0].Name != "The D6" || s.Sets[0] != "b2" {
		t.Errorf("summary %+v", s)
	}
	// A card with no effects still sends a list: the client reads it as is.
	if c := s.Cards[1]; c.Effects == nil {
		t.Errorf("card %+v: effects are nil", c)
	}
	if fmt.Sprint(s.Turns) != "[{1 0 0} {2 1 2} {3 0 4}]" {
		t.Errorf("turns %v", s.Turns)
	}
	if s.End == nil || !s.End.Finished || s.End.Winners[0] != 1 || s.End.Ended != day.Add(time.Hour).UnixMilli() {
		t.Errorf("end %+v", s.End)
	}
	if r.open == nil || len(r.open.Steps) != 6 {
		t.Error("the binding did not keep the record")
	}

	writeFormat(t, filepath.Join(dir, "old"+record.Ext), 1)
	writeFormat(t, filepath.Join(dir, "new"+record.Ext), record.Format+1)
	writeFormat(t, filepath.Join(dir, "noview"+record.Ext), record.Format)
	if err := os.WriteFile(filepath.Join(dir, "broken"+record.Ext), []byte("hello"), 0o600); err != nil {
		t.Fatal(err)
	}
	for name, c := range map[string]struct {
		key   string
		cause error
	}{
		"old":    {"error.record_before_replays", replay.ErrBeforeReplays},
		"new":    {"error.record_newer", replay.ErrNewer},
		"noview": {"error.record_unreadable", replay.ErrNoView},
		"broken": {"error.record_unreadable", record.ErrFormat},
	} {
		_, err := r.Open(filepath.Join(dir, name+record.Ext))
		var m msgerr.Error
		if !errors.As(err, &m) || m.Msg.Key.String() != c.key || !errors.Is(err, c.cause) {
			t.Errorf("%s: %v, want %s", name, err, c.key)
		}
	}
	if _, err := r.Open(filepath.Join(dir, "gone"+record.Ext)); err == nil || errors.As(err, new(msgerr.Error)) {
		t.Errorf("a missing file: %v", err)
	}
}

// TestDelete deletes records only.
func TestDelete(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "g1"+record.Ext)
	writeRecord(t, path, header("g1", time.Now()), 1, nil)
	other := filepath.Join(dir, "notes.txt")
	if err := os.WriteFile(other, []byte("keep"), 0o600); err != nil {
		t.Fatal(err)
	}
	r := New(background{}, dir)
	if err := r.Delete(other); !errors.Is(err, ErrNotRecord) {
		t.Errorf("a text file: %v", err)
	}
	if _, err := os.Stat(other); err != nil {
		t.Errorf("the text file is gone: %v", err)
	}
	if err := r.Delete(path); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("the record is still there: %v", err)
	}
	if err := r.Delete(path); err == nil {
		t.Error("deleting a missing record gave no error")
	}
}

// TestExport copies a record; only records, and never onto itself.
func TestExport(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "g1"+record.Ext)
	writeRecord(t, path, header("g1", time.Now()), 1, nil)
	to := filepath.Join(dir, "copy"+record.Ext)
	if err := copyFile(path, to); err != nil {
		t.Fatal(err)
	}
	a, b := readFile(t, path), readFile(t, to)
	if len(a) == 0 || !bytes.Equal(a, b) {
		t.Errorf("the copy differs: %d and %d bytes", len(a), len(b))
	}
	if err := copyFile(path, path); err != nil {
		t.Errorf("a copy onto itself: %v", err)
	}
	if !bytes.Equal(a, readFile(t, path)) {
		t.Error("a copy onto itself changed the record")
	}
	if _, err := New(background{}, dir).Export(filepath.Join(dir, "notes.txt")); !errors.Is(err, ErrNotRecord) {
		t.Errorf("a text file: %v", err)
	}
}

// readFile reads a test file.
func readFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path) //nolint:gosec // a test file
	if err != nil {
		t.Fatal(err)
	}
	return data
}
