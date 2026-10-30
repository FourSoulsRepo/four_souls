package record

import (
	"path/filepath"
	"testing"
	"time"
)

// writeSteps queues n steps on w.
func writeSteps(t *testing.T, w *Writer, n int) {
	t.Helper()
	for range n {
		if err := w.Step(map[string]int{"k": 1}, nil, 7); err != nil {
			t.Fatal(err)
		}
	}
}

// TestStopWritesEveryLine: a graceful stop writes every batched line (18.17).
func TestStopWritesEveryLine(t *testing.T) {
	path := filepath.Join(t.TempDir(), "g.fsrec")
	w, err := Create(path, header())
	if err != nil {
		t.Fatal(err)
	}
	writeSteps(t, w, 500)
	if err = w.Stop(); err != nil {
		t.Fatal(err)
	}
	rec, err := ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(rec.Steps) != 500 {
		t.Errorf("read %d steps after Stop, want 500", len(rec.Steps))
	}
}

// TestFlushMakesLinesReadable: Flush puts the queued lines in the file
// while the record is still open.
func TestFlushMakesLinesReadable(t *testing.T) {
	path := filepath.Join(t.TempDir(), "g.fsrec")
	w, err := Create(path, header())
	if err != nil {
		t.Fatal(err)
	}
	defer w.Stop() //nolint:errcheck // the test ends here
	writeSteps(t, w, 3)
	if err = w.Flush(); err != nil {
		t.Fatal(err)
	}
	rec, err := ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(rec.Steps) != 3 {
		t.Errorf("read %d steps after Flush, want 3", len(rec.Steps))
	}
}

// TestTickFlushes: lines wait at most flushEvery before they are in the
// file, without a Flush or a stop.
func TestTickFlushes(t *testing.T) {
	path := filepath.Join(t.TempDir(), "g.fsrec")
	w, err := Create(path, header())
	if err != nil {
		t.Fatal(err)
	}
	defer w.Stop() //nolint:errcheck // the test ends here
	writeSteps(t, w, 3)
	time.Sleep(flushEvery + 500*time.Millisecond)
	rec, err := ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(rec.Steps) != 3 {
		t.Errorf("read %d steps after a tick, want 3", len(rec.Steps))
	}
}
