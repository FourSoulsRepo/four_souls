package records

import (
	"cmp"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/FourSoulsRepo/four_souls/internal/msgerr"
	"github.com/FourSoulsRepo/four_souls/internal/replay"
	"github.com/FourSoulsRepo/record"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// Record errors the player reads in their own language (ADR 021);
// their words are in frontend/src/locales/en.json.
var (
	errBeforeReplays = engine.NewKey("error.record_before_replays")
	errNewer         = engine.NewKey("error.record_newer")
	errUnreadable    = engine.NewKey("error.record_unreadable")
	errNotRecord     = engine.NewKey("error.record_not_a_record")
)

// ErrNotRecord means a file without the record extension: Delete and
// Export refuse it.
var ErrNotRecord = errors.New("records: not a record file")

// Why a listed record cannot be watched; "" when it can.
const (
	whyBeforeReplays = "before_replays" // format 1 (ADR 031, 1.5)
	whyNewer         = "newer"          // a newer app made it (ADR 031, 1.6)
	whyUnreadable    = "unreadable"     // not a record, or broken
)

// Entry is one record of the list.
type Entry struct {
	Path      string   `json:"path"`
	Name      string   `json:"name"`    // the file name
	Started   int64    `json:"started"` // Unix ms; 0 when the header is unreadable
	Players   []string `json:"players"` // in seat order
	Finished  bool     `json:"finished"`
	Winners   []string `json:"winners"`
	Steps     int      `json:"steps"` // live steps
	Challenge string   `json:"challenge"`
	// Why is "" when the record can be watched, else before_replays,
	// newer or unreadable.
	Why string `json:"why"`
}

// Seat is one player of a record.
type Seat struct {
	Seat int    `json:"seat"`
	Name string `json:"name"`
}

// Card is the text of one card a record holds (RP-10).
type Card struct {
	Ref     string   `json:"ref"`
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Effects []string `json:"effects"`
}

// Turn is where a turn starts: the index of its first live step.
type Turn struct {
	Number int `json:"number"`
	Active int `json:"active"` // the active player's seat
	Step   int `json:"step"`
}

// End is how a record ends.
type End struct {
	Ended    int64 `json:"ended"` // Unix ms
	Finished bool  `json:"finished"`
	Winners  []int `json:"winners"` // seats
}

// Summary is an opened record: what the player shows before its steps.
type Summary struct {
	Path       string   `json:"path"`
	Game       string   `json:"game"`
	Server     string   `json:"server"`
	Started    int64    `json:"started"` // Unix ms
	Seats      []Seat   `json:"seats"`
	Sets       []string `json:"sets"`
	Challenge  string   `json:"challenge"`
	Difficulty string   `json:"difficulty"`
	Cards      []Card   `json:"cards"`
	Steps      int      `json:"steps"` // live steps
	Turns      []Turn   `json:"turns"`
	End        *End     `json:"end"` // nil when the game never ended
}

// List gives the records in the records folder, newest first: the
// hosted games' records and the saved ones (M-01). A file it cannot
// read is listed with why; a missing folder lists nothing.
func (r *Records) List() ([]Entry, error) {
	files, err := os.ReadDir(r.dir)
	if errors.Is(err, os.ErrNotExist) {
		return []Entry{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("records: %w", err)
	}
	out := []Entry{}
	for _, f := range files {
		if f.IsDir() || filepath.Ext(f.Name()) != record.Ext {
			continue
		}
		out = append(out, entry(filepath.Join(r.dir, f.Name())))
	}
	slices.SortStableFunc(out, func(a, b Entry) int {
		return cmp.Or(cmp.Compare(b.Started, a.Started), cmp.Compare(a.Name, b.Name))
	})
	return out, nil
}

// entry reads one record's header and end, without its steps'
// contents (record.ReadBrief), so a long list stays fast.
func entry(path string) Entry {
	e := Entry{Path: path, Name: filepath.Base(path), Players: []string{}, Winners: []string{}}
	f, err := os.Open(path) //nolint:gosec // a file of the records folder
	if err != nil {
		e.Why = whyUnreadable
		return e
	}
	defer f.Close() //nolint:errcheck // reading only
	rec, err := record.ReadBrief(f)
	switch {
	case errors.Is(err, record.ErrNewer):
		e.Why = whyNewer
		return e
	case err != nil || rec == nil:
		e.Why = whyUnreadable
		return e
	case rec.Header.Format < record.Format:
		e.Why = whyBeforeReplays
	}
	h := rec.Header
	e.Started, e.Challenge, e.Steps = h.Started.UnixMilli(), h.Setup.Challenge, len(rec.Live())
	names := map[int]string{}
	for _, s := range h.Seats {
		e.Players = append(e.Players, s.Name)
		names[s.Seat] = s.Name
	}
	if rec.End != nil && rec.End.Finished {
		e.Finished = true
		for _, w := range rec.End.Winners {
			e.Winners = append(e.Winners, names[w])
		}
	}
	return e
}

// Open opens a record for watching and keeps it for the player, which
// asks for its steps by number (ADR 031, 3.3). It never runs the rules
// engine.
func (r *Records) Open(path string) (Summary, error) {
	rp, err := replay.Open(path)
	if err != nil {
		return Summary{}, playerError(err)
	}
	r.mu.Lock()
	r.open = rp
	r.mu.Unlock()
	return summary(path, rp), nil
}

// playerError gives a record error the player's words; any other
// error (the disk) stays as it is.
func playerError(err error) error {
	var key engine.Key
	switch {
	case errors.Is(err, replay.ErrBeforeReplays):
		key = errBeforeReplays
	case errors.Is(err, replay.ErrNewer):
		key = errNewer
	case errors.Is(err, record.ErrFormat), errors.Is(err, replay.ErrNoView):
		key = errUnreadable
	case errors.Is(err, ErrNotRecord):
		key = errNotRecord
	default:
		return err
	}
	return msgerr.New(key.Msg(), err)
}

// summary is what the player needs of an opened record before its
// steps.
func summary(path string, rp *replay.Replay) Summary {
	h := rp.Header
	s := Summary{
		Path: path, Game: h.Game, Server: h.Server, Started: h.Started.UnixMilli(),
		Seats: []Seat{}, Sets: append([]string{}, h.Setup.Sets...), Cards: []Card{}, Turns: []Turn{},
		Challenge: h.Setup.Challenge, Difficulty: h.Setup.Difficulty, Steps: len(rp.Steps),
	}
	for _, seat := range h.Seats {
		s.Seats = append(s.Seats, Seat{Seat: seat.Seat, Name: seat.Name})
	}
	for _, c := range h.Cards {
		s.Cards = append(s.Cards, Card{Ref: c.Ref, Name: c.Name, Type: c.Type, Effects: append([]string{}, c.Effects...)})
	}
	for _, t := range rp.Turns() {
		s.Turns = append(s.Turns, Turn{Number: t.Number, Active: t.Active, Step: t.Step})
	}
	if e := rp.End; e != nil {
		s.End = &End{Ended: e.Ended.UnixMilli(), Finished: e.Finished, Winners: append([]int{}, e.Winners...)}
	}
	return s
}

// Delete deletes a record file; the screen asks the player first. It
// deletes only files with the record extension.
func (r *Records) Delete(path string) error {
	if filepath.Ext(path) != record.Ext {
		return playerError(ErrNotRecord)
	}
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("records: %w", err)
	}
	return nil
}

// Export asks where to save a copy of a record and copies it there.
// The result is the copy, or "" when the player cancelled.
func (r *Records) Export(path string) (string, error) {
	if filepath.Ext(path) != record.Ext {
		return "", playerError(ErrNotRecord)
	}
	to, err := runtime.SaveFileDialog(r.api.Context(), runtime.SaveDialogOptions{
		Title:           "Export the match record",
		DefaultFilename: filepath.Base(path),
		Filters:         []runtime.FileFilter{{DisplayName: "Four Souls records", Pattern: "*" + record.Ext}},
	})
	if err != nil || to == "" {
		return "", err //nolint:wrapcheck // shown to the player as is
	}
	return to, copyFile(path, to)
}

// copyFile copies a record; a copy onto itself does nothing.
func copyFile(from, to string) error {
	if same(from, to) {
		return nil
	}
	data, err := os.ReadFile(from) //nolint:gosec // a record the player picked
	if err != nil {
		return fmt.Errorf("records: %w", err)
	}
	if err := os.WriteFile(to, data, 0o600); err != nil { //nolint:gosec // where the player chose to save
		return fmt.Errorf("records: %w", err)
	}
	return nil
}

// same tells whether two paths name one file.
func same(a, b string) bool {
	ia, err := os.Stat(a)
	if err != nil {
		return false
	}
	ib, err := os.Stat(b)
	return err == nil && os.SameFile(ia, ib)
}

// PickFile asks for a record anywhere on the computer (M-01); the
// screen then opens it. "" when the player cancelled.
func (r *Records) PickFile() (string, error) {
	path, err := runtime.OpenFileDialog(r.api.Context(), runtime.OpenDialogOptions{
		Title:   "Open a match record",
		Filters: []runtime.FileFilter{{DisplayName: "Four Souls records", Pattern: "*" + record.Ext}},
	})
	return path, err //nolint:wrapcheck // shown to the player as is
}
