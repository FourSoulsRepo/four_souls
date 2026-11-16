// Package replay reads a match record for the replay player (ADR 031).
// It serves the record's live steps with the judge view after each, as
// the server wrote them. It never runs the rules engine: a later rule
// fix never stops a record from playing.
package replay

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/FourSoulsRepo/record"
)

// viewsFrom is the first record format with a view in every step.
const viewsFrom = 2

// ErrBeforeReplays means a format 1 record: it has no views, so it
// cannot play (ADR 031, 1.5).
var ErrBeforeReplays = errors.New("replay: the record was made before replays")

// ErrNewer means a record of a newer format than this app's: update the
// app to watch it (ADR 031, 1.6).
var ErrNewer = errors.New("replay: the record is newer than the app; update the app")

// ErrNoView means a live step of a format 2 record has no view.
var ErrNoView = errors.New("replay: a step of the record has no view")

// Step is one live step: the intent, its events and the judge view
// after it, as raw JSON.
type Step struct {
	N      int             `json:"n"`
	Intent json.RawMessage `json:"intent"`
	Events json.RawMessage `json:"events"`
	View   json.RawMessage `json:"view"`
}

// Replay is a record ready to play: its header (seats, setup, card
// text), its live steps and its end, nil when the game never ended.
type Replay struct {
	Header record.Header
	Steps  []Step
	End    *record.End
}

// Open reads a record file.
func Open(path string) (*Replay, error) {
	f, err := os.Open(path) //nolint:gosec // the player picks the file
	if err != nil {
		return nil, fmt.Errorf("replay: %w", err)
	}
	defer f.Close() //nolint:errcheck // reading only
	return Read(f)
}

// Read reads a record. Steps a rollback undid are left out, as
// record.Live does (ADR 031, 2.3).
func Read(r io.Reader) (*Replay, error) {
	rec, err := record.Read(r)
	switch {
	case errors.Is(err, record.ErrNewer):
		return nil, ErrNewer
	case err != nil:
		return nil, err //nolint:wrapcheck // record's errors say what is wrong
	case rec.Header.Format < viewsFrom:
		return nil, ErrBeforeReplays
	}
	live := rec.Live()
	out := &Replay{Header: rec.Header, End: rec.End, Steps: make([]Step, len(live))}
	for i, s := range live {
		if len(s.View) == 0 {
			return nil, fmt.Errorf("%w: step %d", ErrNoView, s.N)
		}
		out.Steps[i] = Step{N: s.N, Intent: s.Intent, Events: s.Events, View: s.View}
	}
	return out, nil
}

// Turn is where a turn starts: the index of its first live step in
// Steps, the turn's number and its active player, read from the views.
type Turn struct {
	Number int `json:"number"`
	Active int `json:"active"`
	Step   int `json:"step"`
}

// Turns lists the turns in order: a new one starts at each step whose
// view has another turn number than the step before. A view it cannot
// read starts no turn.
func (r *Replay) Turns() []Turn {
	var out []Turn
	last := -1
	for i, s := range r.Steps {
		var v struct {
			Turn struct {
				Number int `json:"number"`
				Active int `json:"active"`
			} `json:"turn"`
		}
		if json.Unmarshal(s.View, &v) != nil || v.Turn.Number == last {
			continue
		}
		last = v.Turn.Number
		out = append(out, Turn{Number: v.Turn.Number, Active: v.Turn.Active, Step: i})
	}
	return out
}
