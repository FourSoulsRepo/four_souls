package record

import (
	"bytes"
	"errors"
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
		if serr := w.Step(map[string]int{"kind": i}, []map[string]string{{"kind": "passed"}}, uint64(100+i)); serr != nil {
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
	if len(rec.Steps) != 3 || rec.Steps[2].N != 3 || rec.Steps[2].Checksum != 102 || string(rec.Steps[1].Intent) != `{"kind":1}` {
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
		if serr := w.Step(i, nil, 1); serr != nil {
			t.Fatal(serr)
		}
	}
	// No Close: the file ends after the last flushed step.
	data, err := os.ReadFile(path) //nolint:gosec // test file
	if err != nil {
		t.Fatal(err)
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
		if serr := w.Step(i, nil, 1); serr != nil {
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
	if serr := w.Step(2, nil, 7); serr != nil {
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
