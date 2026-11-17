package replay

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/FourSoulsRepo/record"
)

func header() record.Header {
	return record.Header{
		App: "1.0.0", Engine: "1.0.0", Protocol: 1, Game: "g1", Started: time.Date(2026, 11, 16, 12, 0, 0, 0, time.UTC),
		Setup: record.Setup{Seed: 7, Players: 2, Sets: []string{"b2"}},
		Seats: []record.Seat{{Seat: 0, Name: "Ann"}, {Seat: 1, Name: "Bo"}},
	}
}

// view is a step's stand-in view.
func view(n int) json.RawMessage {
	return json.RawMessage(fmt.Sprintf(`{"turn":{"number":%d}}`, n))
}

// TestLiveSteps: the steps come with their views; a rollback's undone
// steps are left out (ADR 031, 2.3).
func TestLiveSteps(t *testing.T) {
	path := filepath.Join(t.TempDir(), "g1"+record.Ext)
	w, err := record.Create(path, header())
	if err != nil {
		t.Fatal(err)
	}
	for i := 1; i <= 3; i++ { // steps 1, 2, 3
		if serr := w.Step(map[string]int{"i": i}, []string{"passed"}, view(i), 1); serr != nil {
			t.Fatal(serr)
		}
	}
	if rerr := w.Rollback(2, 0, 2); rerr != nil { // step 4 undoes steps 2 and 3
		t.Fatal(rerr)
	}
	if serr := w.Step(map[string]int{"i": 5}, nil, view(5), 1); serr != nil { // step 5
		t.Fatal(serr)
	}
	if cerr := w.Close(record.End{Ended: time.Date(2026, 11, 16, 13, 0, 0, 0, time.UTC), Finished: true, Winners: []int{0}}); cerr != nil {
		t.Fatal(cerr)
	}
	r, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Steps) != 2 || r.Steps[0].N != 1 || r.Steps[1].N != 5 {
		t.Fatalf("steps %+v", r.Steps)
	}
	if string(r.Steps[1].View) != string(view(5)) || string(r.Steps[0].Intent) != `{"i":1}` || string(r.Steps[0].Events) != `["passed"]` {
		t.Errorf("step %+v", r.Steps[0])
	}
	if r.Header.Seats[1].Name != "Bo" || r.End == nil || !r.End.Finished {
		t.Errorf("header %+v, end %+v", r.Header, r.End)
	}
}

// raw writes a record of format n by hand, with one step line.
func raw(t *testing.T, n int, step record.Step) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	h := header()
	h.Kind, h.Format = "header", n
	enc := json.NewEncoder(gz)
	if err := enc.Encode(h); err != nil {
		t.Fatal(err)
	}
	step.Kind, step.N = "step", 1
	if err := enc.Encode(step); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// TestRefused: format 1 was made before replays; a newer format needs
// a newer app; a step without a view cannot play (ADR 031, 1.5, 1.6).
func TestRefused(t *testing.T) {
	step := record.Step{Intent: json.RawMessage(`{}`), Events: json.RawMessage(`[]`)}
	withView := step
	withView.View = view(1)
	for name, c := range map[string]struct {
		data []byte
		want error
	}{
		"format 1":     {raw(t, 1, step), ErrBeforeReplays},
		"newer format": {raw(t, record.Format+1, withView), ErrNewer},
		"no view":      {raw(t, record.Format, step), ErrNoView},
		"not a record": {[]byte("hello"), record.ErrFormat},
	} {
		if _, err := Read(bytes.NewReader(c.data)); !errors.Is(err, c.want) {
			t.Errorf("%s: %v, want %v", name, err, c.want)
		}
	}
	if _, err := Read(bytes.NewReader(raw(t, record.Format, withView))); err != nil {
		t.Errorf("format %d: %v", record.Format, err)
	}
}

// TestTurns: a turn starts where the view's turn number changes.
func TestTurns(t *testing.T) {
	r := &Replay{}
	for i, v := range []string{
		`{"turn":{"number":0,"active":0}}`,
		`{"turn":{"number":1,"active":1}}`,
		`{"turn":{"number":1,"active":1}}`,
		`not a view`,
		`{"turn":{"number":2,"active":0}}`,
	} {
		r.Steps = append(r.Steps, Step{N: i + 1, View: json.RawMessage(v)})
	}
	got := fmt.Sprint(r.Turns())
	if want := "[{0 0 0} {1 1 1} {2 0 4}]"; got != want {
		t.Errorf("turns %s, want %s", got, want)
	}
}
