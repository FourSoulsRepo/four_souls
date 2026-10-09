package rulesengine

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
)

// StepRecord is one applied intent with what happened and the state
// checksum after it (RP-03, RP-12).
type StepRecord struct {
	Intent   Intent  `json:"intent"`
	Events   []Event `json:"events"`
	Checksum uint64  `json:"checksum"`
}

// Step applies an intent like Apply and also returns the record.
func (g *Game) Step(in Intent) (StepRecord, error) {
	ev, err := g.Apply(in)
	if err != nil {
		return StepRecord{}, err
	}
	sum, err := g.Checksum()
	if err != nil {
		return StepRecord{}, err
	}
	return StepRecord{Intent: in, Events: ev, Checksum: sum}, nil
}

// ErrMismatch means a replayed step ended in a different state.
var ErrMismatch = errors.New("replay: state differs from the record")

// Replay rebuilds a game from its setup and recorded steps. It stops at
// the first step whose checksum differs and returns that step's index
// with ErrMismatch (RP-02, RP-12).
func Replay(s Setup, steps []StepRecord) (*Game, error) {
	g, _, err := NewGame(s)
	if err != nil {
		return nil, err
	}
	for i, st := range steps {
		got, err := g.Step(st.Intent)
		if err != nil {
			return g, fmt.Errorf("replay step %d: %w", i, err)
		}
		if got.Checksum != st.Checksum {
			return g, fmt.Errorf("step %d: %w", i, ErrMismatch)
		}
	}
	return g, nil
}

// Save encodes the whole game state.
func (g *Game) Save() ([]byte, error) {
	data, err := json.Marshal(g)
	if err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}
	return data, nil
}

// Load decodes a saved game. The same card sets must be given; their
// definitions are not part of the saved state.
func Load(data []byte, sets ...CardSet) (*Game, error) {
	var g Game
	if err := json.Unmarshal(data, &g); err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}
	names := make([]string, len(sets))
	for i, s := range sets {
		names[i] = s.Name
	}
	if !slices.Equal(names, g.Sets) {
		return nil, fmt.Errorf("load: game uses sets %v, got %v", g.Sets, names)
	}
	idx, err := newCardIndex(sets...)
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}
	g.cards = idx
	return &g, nil
}
