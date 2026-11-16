package rulesengine

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
)

// A save made while a nested ask waits keeps the last clean state and
// the steps since it; a load replays those steps (19.177, ADR 020 item
// 13). A nested ask is one an effect makes while it runs: inside an
// Ask's answer (Gimpy's "Loot 1, then discard") or inside a wrapper
// (Plan C's DestroyThis). Its question lives in a Go closure, which no
// save can hold; the replay runs the card's own code again, so the
// closure is rebuilt. A clean state has no nested ask waiting. An ask
// that is the ability's own effect is found again by its index
// (Asking.Effect), so it needs no replay.

// ReplayStep is one step a save replays after its clean state.
type ReplayStep struct {
	Intent Intent `json:"intent"`
	// Apply: the step came from Apply, not Step, so it did not extend
	// the checksum chain (tests only; a server always steps).
	Apply bool `json:"apply,omitempty"`
}

// saveForm is the saved form: the clean state, and the steps to
// replay after it when a nested ask waits. Replay comes last, so a save
// without it is the state's JSON alone.
type saveForm struct {
	*Game
	Replay []ReplayStep `json:"replay,omitempty"`
}

// tracking is what a game keeps in memory to save at a nested ask: a
// clean saved state (origin), the steps applied since it, and how many
// of them led to the last clean state.
type tracking struct {
	origin  []byte
	since   []ReplayStep
	cleanAt int
}

// errNoCleanState: the game waits on a nested ask and has no clean
// state to save from.
var errNoCleanState = errors.New("save: a nested question waits and the game has no clean state to save from")

// clean reports whether no nested ask waits.
func (g *Game) clean() bool {
	for _, e := range g.asks {
		if e.nested {
			return false
		}
	}
	return true
}

// beforeApply keeps the state as the origin when the game has none yet:
// before the first applied intent, so changes made by direct calls
// before it (a test table) are kept too.
func (g *Game) beforeApply() {
	if g.saved.origin != nil || !g.clean() {
		return
	}
	data, err := json.Marshal(g)
	if err != nil {
		return // Save then reports there is no clean state
	}
	g.saved.origin = data
}

// afterApply logs an applied intent and notes a clean state.
func (g *Game) afterApply(in Intent) {
	if g.saved.origin == nil {
		return
	}
	g.saved.since = append(g.saved.since, ReplayStep{Intent: in, Apply: true})
	if g.clean() {
		g.saved.cleanAt = len(g.saved.since)
	}
}

// Save encodes the whole game state. While a nested ask waits, it
// encodes the last clean state and the steps since it instead; Load
// replays them (19.177).
func (g *Game) Save() ([]byte, error) {
	if g.clean() {
		data, err := json.Marshal(g)
		if err != nil {
			return nil, fmt.Errorf("save: %w", err)
		}
		g.saved = tracking{origin: slices.Clone(data)} // the caller owns data
		return data, nil
	}
	if err := g.moveOrigin(); err != nil {
		return nil, err
	}
	data, err := json.Marshal(saveForm{Game: g.originGame(), Replay: g.saved.since})
	if err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}
	return data, nil
}

// originGame decodes the origin for encoding; its bytes came from
// json.Marshal, so they decode.
func (g *Game) originGame() *Game {
	var o Game
	if err := json.Unmarshal(g.saved.origin, &o); err != nil {
		panic("rulesengine: the clean state does not decode: " + err.Error())
	}
	return &o
}

// moveOrigin makes the last clean state the origin: it replays the
// steps up to it from the old origin. It replays the rest too and
// checks the result is this game: a change made outside Apply since the
// origin (a test's direct call) cannot be saved this way.
func (g *Game) moveOrigin() error {
	if g.saved.origin == nil {
		return errNoCleanState
	}
	if g.saved.cleanAt == 0 {
		return nil
	}
	r, err := g.loadState(g.saved.origin)
	if err != nil {
		return err
	}
	since := g.saved.since
	if err = r.replaySteps(since[:g.saved.cleanAt]); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	origin, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("save: %w", err)
	}
	if err = r.replaySteps(since[g.saved.cleanAt:]); err != nil {
		return fmt.Errorf("save: %w", err)
	}
	a, errA := r.Checksum()
	b, errB := g.Checksum()
	if errA != nil || errB != nil {
		return fmt.Errorf("save: %w", errors.Join(errA, errB))
	}
	if a != b {
		return errors.New("save: the game changed outside Apply since its last clean state")
	}
	g.saved = tracking{origin: origin, since: slices.Clone(since[g.saved.cleanAt:])}
	return nil
}

// replaySteps applies saved steps as they were applied.
func (g *Game) replaySteps(steps []ReplayStep) error {
	for i, st := range steps {
		var err error
		if st.Apply {
			_, err = g.Apply(st.Intent)
		} else {
			_, err = g.Step(st.Intent)
		}
		if err != nil {
			return fmt.Errorf("replay step %d of the save: %w", i, err)
		}
	}
	return nil
}

// loadState decodes a clean state with this game's cards.
func (g *Game) loadState(data []byte) (*Game, error) {
	var r Game
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("save: %w", err)
	}
	r.cards = g.cards
	r.saved = tracking{origin: data}
	return &r, nil
}

// Load decodes a saved game. The same card sets must be given; their
// definitions are not part of the saved state. A save made while a
// nested ask waits replays its steps after the clean state (19.177).
func Load(data []byte, sets ...CardSet) (*Game, error) {
	var g Game
	s := saveForm{Game: &g}
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}
	names := make([]string, len(sets))
	for i, set := range sets {
		names[i] = set.Name
	}
	if !slices.Equal(names, g.Sets) {
		return nil, fmt.Errorf("load: game uses sets %v, got %v", g.Sets, names)
	}
	idx, err := newCardIndex(sets...)
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}
	g.cards = idx
	if len(s.Replay) == 0 {
		g.saved = tracking{origin: slices.Clone(data)}
		return &g, nil
	}
	origin, err := json.Marshal(&g)
	if err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}
	g.saved = tracking{origin: origin}
	if err := g.replaySteps(s.Replay); err != nil {
		return nil, fmt.Errorf("load: %w", err)
	}
	return &g, nil
}

// NestedAsk reports whether a nested ask waits: a save now keeps the
// last clean state and replays the steps since it (19.177).
func (g *Game) NestedAsk() bool { return !g.clean() }
