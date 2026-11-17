package record

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func header() Header {
	return Header{
		App: "0.1.0", Engine: "0.1.0", Protocol: 1, Game: "g1", Started: time.Date(2026, 10, 10, 12, 0, 0, 0, time.UTC),
		Setup: Setup{Seed: 42, Players: 2, Sets: []string{"b2"}, BonusSouls: true, Characters: []string{"isaac", "cain"}},
		Seats: []Seat{{Seat: 0, Name: "Ann"}, {Seat: 1, Name: "Bo"}},
		Cards: []Card{{Ref: "the_d6", Name: "The D6", Effects: []string{"Reroll a dice roll."}}},
	}
}

func TestRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "g1"+Ext)
	w, err := Create(path, header())
	if err != nil {
		t.Fatal(err)
	}
	for i := range 3 {
		if serr := w.Step(map[string]int{"kind": i}, []map[string]string{{"kind": "passed"}}, json.RawMessage(fmt.Sprintf(`{"n":%d}`, i)), uint64(100+i)); serr != nil {
			t.Fatal(serr)
		}
	}
	if cerr := w.Close(End{Ended: time.Date(2026, 10, 10, 13, 0, 0, 0, time.UTC), Finished: true, Winners: []int{1}}); cerr != nil {
		t.Fatal(cerr)
	}
	rec, err := ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Header.Format != Format || rec.Header.Setup.Seed != 42 || rec.Header.Seats[1].Name != "Bo" || rec.Header.Cards[0].Ref != "the_d6" {
		t.Errorf("header %+v", rec.Header)
	}
	if len(rec.Steps) != 3 || rec.Steps[2].N != 3 || rec.Steps[2].Checksum != 102 || string(rec.Steps[1].Intent) != `{"kind":1}` || string(rec.Steps[1].View) != `{"n":1}` {
		t.Errorf("steps %+v", rec.Steps)
	}
	if rec.End == nil || !rec.End.Finished || rec.End.Winners[0] != 1 {
		t.Errorf("end %+v", rec.End)
	}
	if _, err := Create(path, header()); err == nil {
		t.Error("an existing record was overwritten")
	}
}

// A server that crashes leaves a record without an end; it is still
// readable up to the last step (RP-08).
func TestUnfinished(t *testing.T) {
	path := filepath.Join(t.TempDir(), "g2"+Ext)
	w, err := Create(path, header())
	if err != nil {
		t.Fatal(err)
	}
	for i := range 5 {
		if serr := w.Step(i, nil, nil, 1); serr != nil {
			t.Fatal(serr)
		}
	}
	// No Close: the file ends after the last flushed step.
	if ferr := w.Flush(); ferr != nil {
		t.Fatal(ferr)
	}
	data, err := os.ReadFile(path) //nolint:gosec // test file
	if err != nil {
		t.Fatal(err)
	}
	if serr := w.Stop(); serr != nil {
		t.Fatal(serr)
	}
	rec, err := Read(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if len(rec.Steps) != 5 || rec.End != nil {
		t.Errorf("%d steps, end %v", len(rec.Steps), rec.End)
	}
}

func TestNotARecord(t *testing.T) {
	for name, data := range map[string][]byte{"text": []byte("hello"), "empty": nil} {
		if _, err := Read(bytes.NewReader(data)); !errors.Is(err, ErrFormat) {
			t.Errorf("%s: %v", name, err)
		}
	}
}

// writeFormat writes a gzip record of one header line of format n and
// one step line without a view.
func writeFormat(t *testing.T, n int) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	h := header()
	h.Kind, h.Format = "header", n
	enc := json.NewEncoder(gz)
	if err := enc.Encode(h); err != nil {
		t.Fatal(err)
	}
	if err := enc.Encode(Step{Kind: "step", N: 1, Intent: json.RawMessage(`{}`), Events: json.RawMessage(`[]`), Checksum: 1}); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// TestFormats: format 1 still reads, with no views; a newer format is
// refused as ErrNewer, which is an ErrFormat too (ADR 031, 1.6).
func TestFormats(t *testing.T) {
	rec, err := Read(bytes.NewReader(writeFormat(t, 1)))
	if err != nil || rec.Header.Format != 1 || len(rec.Steps) != 1 || rec.Steps[0].View != nil {
		t.Errorf("format 1: %+v, %v", rec, err)
	}
	if _, err := Read(bytes.NewReader(writeFormat(t, Format+1))); !errors.Is(err, ErrNewer) || !errors.Is(err, ErrFormat) {
		t.Errorf("a newer format: %v", err)
	}
	if _, err := Read(bytes.NewReader(writeFormat(t, 0))); !errors.Is(err, ErrFormat) || errors.Is(err, ErrNewer) {
		t.Errorf("format 0: %v", err)
	}
}

func TestPrune(t *testing.T) {
	dir := t.TempDir()
	now := time.Date(2026, 10, 10, 0, 0, 0, 0, time.UTC)
	for name, age := range map[string]int{"old" + Ext: 40, "new" + Ext: 5, "notes.txt": 40} {
		p := filepath.Join(dir, name)
		if err := os.WriteFile(p, nil, 0o600); err != nil {
			t.Fatal(err)
		}
		when := now.AddDate(0, 0, -age)
		if err := os.Chtimes(p, when, when); err != nil {
			t.Fatal(err)
		}
	}
	if n, err := Prune(dir, 30, now); err != nil || n != 1 {
		t.Fatalf("pruned %d, %v", n, err)
	}
	if _, err := os.Stat(filepath.Join(dir, "new"+Ext)); err != nil {
		t.Error("a new record was deleted")
	}
	if _, err := os.Stat(filepath.Join(dir, "notes.txt")); err != nil {
		t.Error("another file was deleted")
	}
	if n, err := Prune(dir, 0, now); err != nil || n != 0 {
		t.Error("retention 0 deleted records")
	}
}

// TestAppend: a saved game continues in the same file (N-11).
func TestAppend(t *testing.T) {
	path := filepath.Join(t.TempDir(), "g3"+Ext)
	w, err := Create(path, header())
	if err != nil {
		t.Fatal(err)
	}
	for i := range 2 {
		if serr := w.Step(i, nil, nil, 1); serr != nil {
			t.Fatal(serr)
		}
	}
	if cerr := w.Close(End{Finished: false}); cerr != nil {
		t.Fatal(cerr)
	}
	w, err = Append(path, 2, time.Date(2026, 10, 11, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if serr := w.Step(2, nil, nil, 7); serr != nil {
		t.Fatal(serr)
	}
	if cerr := w.Close(End{Finished: true, Winners: []int{0}}); cerr != nil {
		t.Fatal(cerr)
	}
	rec, err := ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(rec.Steps) != 3 || rec.Steps[2].N != 3 || rec.Resumes != 1 || rec.End == nil || !rec.End.Finished {
		t.Errorf("steps %d (last %d), resumes %d, end %+v", len(rec.Steps), rec.Steps[len(rec.Steps)-1].N, rec.Resumes, rec.End)
	}
}

// TestStopThenAppend: a server restart stops the record without an end
// line; the next server goes on with it (N-20).
func TestStopThenAppend(t *testing.T) {
	path := filepath.Join(t.TempDir(), "g4"+Ext)
	w, err := Create(path, header())
	if err != nil {
		t.Fatal(err)
	}
	if serr := w.Step(0, nil, nil, 1); serr != nil {
		t.Fatal(serr)
	}
	if serr := w.Stop(); serr != nil {
		t.Fatal(serr)
	}
	rec, err := ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(rec.Steps) != 1 || rec.End != nil {
		t.Fatalf("stopped: steps %d, end %+v; want 1 and none", len(rec.Steps), rec.End)
	}
	w, err = Append(path, 1, time.Date(2026, 10, 21, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatal(err)
	}
	if serr := w.Step(1, nil, nil, 2); serr != nil {
		t.Fatal(serr)
	}
	if cerr := w.Close(End{Finished: true}); cerr != nil {
		t.Fatal(cerr)
	}
	rec, err = ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(rec.Steps) != 2 || rec.Steps[1].N != 2 || rec.Resumes != 1 || rec.End == nil {
		t.Errorf("steps %d, resumes %d, end %+v", len(rec.Steps), rec.Resumes, rec.End)
	}
}

// TestServerName: the header names a named server; an app's own host
// has no name, and the header leaves the field out (N-18).
func TestServerName(t *testing.T) {
	for _, name := range []string{"Basement", ""} {
		path := filepath.Join(t.TempDir(), "g1"+Ext)
		h := header()
		h.Server = name
		w, err := Create(path, h)
		if err != nil {
			t.Fatal(err)
		}
		if cerr := w.Close(End{}); cerr != nil {
			t.Fatal(cerr)
		}
		rec, err := ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if rec.Header.Server != name {
			t.Errorf("server %q, want %q", rec.Header.Server, name)
		}
		line, err := json.Marshal(rec.Header)
		if err != nil {
			t.Fatal(err)
		}
		if has := bytes.Contains(line, []byte(`"server"`)); has != (name != "") {
			t.Errorf("header %s: server field %v", line, has)
		}
	}
}

// TestRollbackLines: a rollback line takes a step number; the undone
// steps stay in the record, and Live leaves them out (ADR 020, item 9).
func TestRollbackLines(t *testing.T) {
	path := filepath.Join(t.TempDir(), "g5"+Ext)
	h := header()
	h.Rollback = "god"
	w, err := Create(path, h)
	if err != nil {
		t.Fatal(err)
	}
	step := func(n uint64) {
		t.Helper()
		if serr := w.Step(map[string]uint64{"kind": n}, nil, nil, n); serr != nil {
			t.Fatal(serr)
		}
	}
	back := func(to, seat, undid int) {
		t.Helper()
		if rerr := w.Rollback(to, seat, undid); rerr != nil {
			t.Fatal(rerr)
		}
	}
	// Steps 1-5; 6 undoes 4 and 5; 7, 8; 9 undoes 3, 7 and 8; 10.
	for n := range uint64(5) {
		step(n + 1)
	}
	back(4, 1, 2)
	step(7)
	step(8)
	back(3, 0, 3)
	step(10)
	if cerr := w.Close(End{Ended: time.Now().UTC()}); cerr != nil {
		t.Fatal(cerr)
	}
	rec, err := ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Header.Rollback != "god" {
		t.Errorf("header rollback %q", rec.Header.Rollback)
	}
	if len(rec.Steps) != 8 || rec.Steps[5].N != 7 || rec.Steps[7].N != 10 {
		t.Errorf("steps %+v", rec.Steps)
	}
	want := []Rollback{{Kind: "rollback", N: 6, To: 4, Seat: 1, Undid: 2}, {Kind: "rollback", N: 9, To: 3, Seat: 0, Undid: 3}}
	if len(rec.Rollbacks) != 2 || rec.Rollbacks[0] != want[0] || rec.Rollbacks[1] != want[1] {
		t.Errorf("rollbacks %+v", rec.Rollbacks)
	}
	var live []int
	for _, s := range rec.Live() {
		live = append(live, s.N)
	}
	if fmt.Sprint(live) != "[1 2 10]" {
		t.Errorf("live steps %v, want [1 2 10]", live)
	}
	// A rollback as the last line undoes too.
	rec.Rollbacks = append(rec.Rollbacks, Rollback{N: 11, To: 2})
	if got := rec.Live(); len(got) != 1 || got[0].N != 1 {
		t.Errorf("after a last rollback: %+v", got)
	}
}

// TestChallengeSetup: a challenge record keeps its challenge and
// difficulty, and a lost one ends finished with no winners; a regular
// record leaves both fields out (ADR 028, Decision 10.5).
func TestChallengeSetup(t *testing.T) {
	for _, c := range []struct{ challenge, difficulty string }{{"krampus", "hard"}, {"", ""}} {
		path := filepath.Join(t.TempDir(), "g1"+Ext)
		h := header()
		h.Setup.Challenge, h.Setup.Difficulty, h.Setup.BonusSouls = c.challenge, c.difficulty, false
		w, err := Create(path, h)
		if err != nil {
			t.Fatal(err)
		}
		if cerr := w.Close(End{Finished: true}); cerr != nil {
			t.Fatal(cerr)
		}
		rec, err := ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if s := rec.Header.Setup; s.Challenge != c.challenge || s.Difficulty != c.difficulty {
			t.Errorf("setup %+v, want %+v", s, c)
		}
		if rec.End == nil || !rec.End.Finished || len(rec.End.Winners) != 0 {
			t.Errorf("end %+v", rec.End)
		}
		line, err := json.Marshal(rec.Header.Setup)
		if err != nil {
			t.Fatal(err)
		}
		for _, field := range []string{`"challenge"`, `"difficulty"`} {
			if has := bytes.Contains(line, []byte(field)); has != (c.challenge != "") {
				t.Errorf("setup %s: %s field %v", line, field, has)
			}
		}
	}
}

// TestReadBrief: a brief read has the header, the end and every step's
// N, the rollbacks too, but no step's intent, events or view.
func TestReadBrief(t *testing.T) {
	path := filepath.Join(t.TempDir(), "g7"+Ext)
	w, err := Create(path, header())
	if err != nil {
		t.Fatal(err)
	}
	for i := range 4 {
		if serr := w.Step(map[string]int{"kind": i}, []string{"passed"}, json.RawMessage(`{"turn":{"number":1}}`), 1); serr != nil {
			t.Fatal(serr)
		}
	}
	if rerr := w.Rollback(3, 0, 2); rerr != nil {
		t.Fatal(rerr)
	}
	if cerr := w.Close(End{Ended: time.Now().UTC(), Finished: true, Winners: []int{0}}); cerr != nil {
		t.Fatal(cerr)
	}
	f, err := os.Open(path) //nolint:gosec // a test file
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close() //nolint:errcheck // reading only
	rec, err := ReadBrief(f)
	if err != nil {
		t.Fatal(err)
	}
	if rec.Header.Seats[1].Name != "Bo" || rec.End == nil || !rec.End.Finished || len(rec.Rollbacks) != 1 {
		t.Errorf("brief %+v", rec)
	}
	if len(rec.Steps) != 4 || rec.Steps[3].N != 4 || rec.Steps[0].View != nil || rec.Steps[0].Intent != nil {
		t.Errorf("steps %+v", rec.Steps)
	}
	if live := rec.Live(); len(live) != 2 {
		t.Errorf("live %+v", live)
	}
	if _, err := ReadBrief(bytes.NewReader(stepFirst(t))); !errors.Is(err, ErrFormat) {
		t.Errorf("a step before the header: %v", err)
	}
}

// stepFirst is a gzip file whose first line is a step.
func stepFirst(t *testing.T) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write([]byte(`{"kind":"step","n":1,"intent":{}}` + "\n")); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}
