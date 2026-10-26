// Package record reads and writes match records (ADR 007): a .fsrec
// file is gzip-compressed JSON lines, a header, one line per applied
// step or rollback and, for a game that ended, an end line. Intents
// and events are kept as raw JSON, so the package needs no engine and
// old records stay readable.
package record

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"slices"
	"time"
)

// Format is the record format version.
const Format = 1

// Ext is the file extension of records (RP-09).
const Ext = ".fsrec"

// Header opens a record.
type Header struct {
	Kind     string    `json:"kind"` // "header"
	Format   int       `json:"format"`
	App      string    `json:"app"`
	Engine   string    `json:"engine"`
	Protocol int       `json:"protocol"`
	Game     string    `json:"game"`
	Server   string    `json:"server,omitempty"` // the server's name; none on an app's own host (N-18)
	Started  time.Time `json:"started"`
	Setup    Setup     `json:"setup"`
	Seats    []Seat    `json:"seats"`
	// Cards holds the text of the cards, so outdated clients can show
	// them (RP-10).
	Cards []Card `json:"cards,omitempty"`
	// Rollback is the game's rollback mode: "", "soft" or "god" (ADR 020).
	Rollback string `json:"rollback,omitempty"`
}

// Setup is what is needed to rebuild the game: the engine's setup
// without card definitions.
type Setup struct {
	Seed       uint64   `json:"seed"`
	Players    int      `json:"players"`
	Sets       []string `json:"sets"`
	BonusSouls bool     `json:"bonus_souls"`
	Characters []string `json:"characters,omitempty"`
	Banned     []string `json:"banned,omitempty"`      // left out of the game (GS-13)
	DeckRatios bool     `json:"deck_ratios,omitempty"` // decks built from the deck ratios (GS-07)
}

// Seat is one player of the match.
type Seat struct {
	Seat int    `json:"seat"`
	Name string `json:"name"`
}

// Card is the text of one card.
type Card struct {
	Ref     string   `json:"ref"`
	Name    string   `json:"name"`
	Type    string   `json:"type,omitempty"`
	Effects []string `json:"effects,omitempty"`
}

// Step is one applied step: the intent, every event (hands too, RP-04)
// and the state checksum after it (RP-12).
type Step struct {
	Kind     string          `json:"kind"` // "step"
	N        int             `json:"n"`
	Intent   json.RawMessage `json:"intent"`
	Events   json.RawMessage `json:"events"`
	Checksum uint64          `json:"checksum"`
}

// Rollback undid the live steps from step To on (ADR 020, item 9).
// Undone steps stay in the record. N is the rollback's own number:
// step numbers never repeat, so the next step is N+1.
type Rollback struct {
	Kind  string `json:"kind"` // "rollback"
	N     int    `json:"n"`
	To    int    `json:"to"`    // the first undone step
	Seat  int    `json:"seat"`  // who rolled back
	Undid int    `json:"undid"` // how many live steps it undid
}

// End closes a record.
type End struct {
	Kind     string    `json:"kind"` // "end"
	Ended    time.Time `json:"ended"`
	Finished bool      `json:"finished"` // false: stopped before a winner (RP-08)
	Winners  []int     `json:"winners,omitempty"`
}

// Writer appends to a record file. Each step is flushed, so a record of
// a game that never ended is still readable (RP-08).
type Writer struct {
	f     *os.File
	gz    *gzip.Writer
	enc   *json.Encoder
	steps int
}

// Create starts a record file with its header.
func Create(path string, h Header) (*Writer, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return nil, fmt.Errorf("record: %w", err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600) //nolint:gosec // the server names the file
	if err != nil {
		return nil, fmt.Errorf("record: %w", err)
	}
	w := &Writer{f: f, gz: gzip.NewWriter(f)}
	w.enc = json.NewEncoder(w.gz)
	h.Kind, h.Format = "header", Format
	if err := w.write(h); err != nil {
		_ = f.Close() //nolint:errcheck // the write error matters
		return nil, err
	}
	return w, nil
}

// Step appends one applied step.
func (w *Writer) Step(intent, events any, checksum uint64) error {
	in, err := json.Marshal(intent)
	if err != nil {
		return fmt.Errorf("record: %w", err)
	}
	ev, err := json.Marshal(events)
	if err != nil {
		return fmt.Errorf("record: %w", err)
	}
	w.steps++
	return w.write(Step{Kind: "step", N: w.steps, Intent: in, Events: ev, Checksum: checksum})
}

// Rollback appends a rollback to before step to; it takes the next
// step number.
func (w *Writer) Rollback(to, seat, undid int) error {
	w.steps++
	return w.write(Rollback{Kind: "rollback", N: w.steps, To: to, Seat: seat, Undid: undid})
}

// Close writes the end line and closes the file.
func (w *Writer) Close(end End) error {
	end.Kind = "end"
	err := w.write(end)
	if cerr := w.gz.Close(); err == nil && cerr != nil {
		err = fmt.Errorf("record: %w", cerr)
	}
	if cerr := w.f.Close(); err == nil && cerr != nil {
		err = fmt.Errorf("record: %w", cerr)
	}
	return err
}

// Stop closes the file without an end line: the game goes on after a
// server restart, and Append continues the record (N-20).
func (w *Writer) Stop() error {
	err := w.gz.Close()
	if err != nil {
		err = fmt.Errorf("record: %w", err)
	}
	if cerr := w.f.Close(); err == nil && cerr != nil {
		err = fmt.Errorf("record: %w", cerr)
	}
	return err
}

func (w *Writer) write(v any) error {
	if err := w.enc.Encode(v); err != nil {
		return fmt.Errorf("record: %w", err)
	}
	if err := w.gz.Flush(); err != nil {
		return fmt.Errorf("record: %w", err)
	}
	return nil
}

// Record is a read record. End is nil when the game never ended or the
// file was cut off.
type Record struct {
	Header Header
	Steps  []Step
	// Rollbacks are the rollback lines, in order; Live applies them.
	Rollbacks []Rollback
	End       *End
	// Resumes counts how often the game was saved and loaded again.
	Resumes int
}

// ErrFormat means a file is not a record this package can read.
var ErrFormat = errors.New("record: not a record of a known format")

// Read reads a whole record. A file cut off in the middle (a crash) is
// read up to its last complete line.
func Read(r io.Reader) (*Record, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFormat, err)
	}
	defer gz.Close() //nolint:errcheck // reading only
	sc := bufio.NewScanner(gz)
	sc.Buffer(make([]byte, 64<<10), 16<<20)
	var rec Record
	first := true
	for sc.Scan() {
		var kind struct {
			Kind string `json:"kind"`
		}
		line := sc.Bytes()
		if err := json.Unmarshal(line, &kind); err != nil {
			return &rec, fmt.Errorf("%w: %w", ErrFormat, err)
		}
		switch {
		case first && kind.Kind != "header":
			return nil, ErrFormat
		case kind.Kind == "header":
			if err := json.Unmarshal(line, &rec.Header); err != nil || rec.Header.Format != Format {
				return nil, ErrFormat
			}
		case kind.Kind == "step":
			var s Step
			if err := json.Unmarshal(line, &s); err != nil {
				return &rec, fmt.Errorf("%w: %w", ErrFormat, err)
			}
			rec.Steps = append(rec.Steps, s)
		case kind.Kind == "rollback":
			var b Rollback
			if err := json.Unmarshal(line, &b); err != nil {
				return &rec, fmt.Errorf("%w: %w", ErrFormat, err)
			}
			rec.Rollbacks = append(rec.Rollbacks, b)
		case kind.Kind == "end":
			var e End
			if err := json.Unmarshal(line, &e); err != nil {
				return &rec, fmt.Errorf("%w: %w", ErrFormat, err)
			}
			rec.End = &e
		case kind.Kind == "resume":
			rec.End = nil // the game went on after a save (N-11)
			rec.Resumes++
		}
		first = false
	}
	if err := sc.Err(); err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return &rec, fmt.Errorf("record: %w", err)
	}
	if first {
		return nil, ErrFormat
	}
	return &rec, nil
}

// Live is the steps a replay applies: every step, less those a
// rollback undid (ADR 020, item 9). A rollback never reaches behind a
// resume, so a replay from the start with these steps rebuilds the game.
func (r *Record) Live() []Step {
	var out []Step
	b := r.Rollbacks
	undo := func(before int) {
		for len(b) > 0 && b[0].N < before {
			to := b[0].To
			out = slices.DeleteFunc(out, func(s Step) bool { return s.N >= to })
			b = b[1:]
		}
	}
	for _, s := range r.Steps {
		undo(s.N)
		out = append(out, s)
	}
	undo(math.MaxInt)
	return out
}

// ReadFile reads a record file.
func ReadFile(path string) (*Record, error) {
	f, err := os.Open(path) //nolint:gosec // the caller names the file
	if err != nil {
		return nil, fmt.Errorf("record: %w", err)
	}
	defer f.Close() //nolint:errcheck // reading only
	return Read(f)
}

// Prune deletes records in dir older than days (RP-07); 0 keeps them
// all. It returns how many it deleted.
func Prune(dir string, days int, now time.Time) (int, error) {
	if days <= 0 {
		return 0, nil
	}
	entries, err := os.ReadDir(dir)
	if errors.Is(err, os.ErrNotExist) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("record: %w", err)
	}
	cutoff := now.AddDate(0, 0, -days)
	n := 0
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != Ext {
			continue
		}
		info, err := e.Info()
		if err != nil || !info.ModTime().Before(cutoff) {
			continue
		}
		if err := os.Remove(filepath.Join(dir, e.Name())); err == nil {
			n++
		}
	}
	return n, nil
}

// Resume marks where a saved game was loaded again (N-11), or where a
// game went on after a server restart (N-20).
type Resume struct {
	Kind    string    `json:"kind"` // "resume"
	Resumed time.Time `json:"resumed"`
}

// Append continues a record after a saved game was loaded: a new gzip
// part at the end of the same file (gzip readers read every part), a
// resume line, then the next steps from step steps+1.
func Append(path string, steps int, at time.Time) (*Writer, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o600) //nolint:gosec // the server names the file
	if err != nil {
		return nil, fmt.Errorf("record: %w", err)
	}
	w := &Writer{f: f, gz: gzip.NewWriter(f), steps: steps}
	w.enc = json.NewEncoder(w.gz)
	if err := w.write(Resume{Kind: "resume", Resumed: at}); err != nil {
		_ = f.Close() //nolint:errcheck // the write error matters
		return nil, err
	}
	return w, nil
}
