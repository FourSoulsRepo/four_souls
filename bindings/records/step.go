package records

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ErrNotOpen means Step was called before Open opened a record.
var ErrNotOpen = errors.New("records: no record is open")

// ErrNoStep means a step number outside the open record's live steps.
var ErrNoStep = errors.New("records: no such step")

// Step is one live step of the open record, as the player shows it
// (ADR 031, 3.3): the events it played and the judge view after it.
type Step struct {
	Index  int             `json:"index"` // in the live steps, from 0
	N      int             `json:"n"`     // the record's own step number
	Events json.RawMessage `json:"events" ts_type:"unknown"`
	View   json.RawMessage `json:"view" ts_type:"unknown"`
}

// Step gives live step i of the record Open opened, from 0. It reads
// the kept record only: no file, no rules engine.
func (r *Records) Step(i int) (Step, error) {
	r.mu.Lock()
	rp := r.open
	r.mu.Unlock()
	if rp == nil {
		return Step{}, ErrNotOpen
	}
	if i < 0 || i >= len(rp.Steps) {
		return Step{}, fmt.Errorf("%w: %d of %d", ErrNoStep, i, len(rp.Steps))
	}
	s := rp.Steps[i]
	events := s.Events
	if len(events) == 0 {
		events = json.RawMessage("null")
	}
	return Step{Index: i, N: s.N, Events: events, View: s.View}, nil
}
