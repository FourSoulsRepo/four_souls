package rulesengine

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"hash/fnv"
	"slices"
)

// StepRecord is one applied intent with what happened and the state
// checksum after it (RP-03, RP-12).
type StepRecord struct {
	Intent   Intent  `json:"intent"`
	Events   []Event `json:"events"`
	Checksum uint64  `json:"checksum"`
}

// anchorEvery is how many steps pass between full-state checks: the
// checksum of every such step also covers the whole state (RP-12).
const anchorEvery = 100

// Step applies an intent like Apply and also returns the record. Its
// checksum extends the history chain; on an anchor step it also covers
// the whole state, so a state that differs without showing in the
// history is caught there.
func (g *Game) Step(in Intent) (StepRecord, error) {
	ev, err := g.Apply(in)
	if err != nil {
		return StepRecord{}, err
	}
	chain, err := g.chainWith(in, ev)
	if err != nil {
		return StepRecord{}, err
	}
	g.Chain = chain
	g.Applied++
	sum := chain
	if g.Applied%anchorEvery == 0 {
		full, err := g.Checksum()
		if err != nil {
			return StepRecord{}, err
		}
		sum = mixSums(chain, full)
	}
	return StepRecord{Intent: in, Events: ev, Checksum: sum}, nil
}

// chainWith hashes the chain so far with one step: its intent, what
// happened and the random state after it. The events and the RNG show
// what the state did, so the chain catches any change to them.
func (g *Game) chainWith(in Intent, ev []Event) (uint64, error) {
	h := fnv.New64a()
	writeSum(h, g.Chain)
	for _, v := range []any{in, ev, g.RNG} {
		if err := json.NewEncoder(h).Encode(v); err != nil {
			return 0, fmt.Errorf("checksum: %w", err)
		}
	}
	return h.Sum64(), nil
}

// mixSums folds the whole-state sum into a chain value.
func mixSums(chain, full uint64) uint64 {
	h := fnv.New64a()
	writeSum(h, chain)
	writeSum(h, full)
	return h.Sum64()
}

// writeSum writes one sum as eight bytes.
func writeSum(h hash.Hash64, v uint64) {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], v)
	h.Write(b[:]) //nolint:gosec // a hash never fails to write
}

// ErrMismatch means a replayed step ended in a different state.
var ErrMismatch = errors.New("replay: state differs from the record")

// CanReplay says whether a record made by engine version v replays on
// this engine. Records made before the throw (ADR 014) do not: their
// steps have no throw intents. Nor do records made before messages
// (16.4): their states, so their checksums, held English text.
func CanReplay(v string) error {
	if older(v, throwVersion) {
		return fmt.Errorf("replay: the record was made by rules engine %s; engines from %s on need a throw for every roll, so it cannot replay", v, throwVersion)
	}
	if older(v, msgVersion) {
		return fmt.Errorf("replay: the record was made by rules engine %s; engines from %s on keep messages instead of text, so its checksums differ and it cannot replay", v, msgVersion)
	}
	if older(v, chainVersion) {
		return fmt.Errorf("replay: the record was made by rules engine %s; engines from %s on check a hash chain over the steps, so its checksums differ and it cannot replay", v, chainVersion)
	}
	if older(v, orderVersion) {
		return fmt.Errorf("replay: the record was made by rules engine %s; engines from %s on order some steps differently, so its checksums differ and it cannot replay", v, orderVersion)
	}
	return nil
}

// older reports whether version a (X.Y.Z) is older than b. A version
// that does not parse counts as older.
func older(a, b string) bool {
	var x, y [3]int
	if n, err := fmt.Sscanf(a, "%d.%d.%d", &x[0], &x[1], &x[2]); err != nil || n != 3 {
		return true
	}
	if _, err := fmt.Sscanf(b, "%d.%d.%d", &y[0], &y[1], &y[2]); err != nil {
		panic("rulesengine: bad version " + b)
	}
	return slices.Compare(x[:], y[:]) < 0
}

// Replay rebuilds a game from its setup and recorded steps. It stops at
// the first step whose checksum differs and returns that step's index
// with ErrMismatch (RP-02, RP-12).
func Replay(s Setup, steps []StepRecord) (*Game, error) {
	g, _, err := NewGame(s)
	if err != nil {
		return nil, err
	}
	return g, g.replay(steps)
}

// ReplayFrom rebuilds a game from a saved base state and the steps
// played since it, like Replay. A rollback to before step k replays the
// steps before k (ADR 020, item 7). The base keeps the RNG, so the
// rebuilt game rolls and shuffles as it did then (N-12, item 7). The
// same card sets as for Load must be given.
func ReplayFrom(state []byte, steps []StepRecord, sets ...CardSet) (*Game, error) {
	g, err := Load(state, sets...)
	if err != nil {
		return nil, err
	}
	return g, g.replay(steps)
}

// replay applies recorded steps and checks each checksum. It stops at
// the first step that is refused or ends in a different state.
func (g *Game) replay(steps []StepRecord) error {
	for i, st := range steps {
		got, err := g.Step(st.Intent)
		if err != nil {
			return fmt.Errorf("replay step %d: %w", i, err)
		}
		if got.Checksum != st.Checksum {
			return fmt.Errorf("step %d: %w", i, ErrMismatch)
		}
	}
	return nil
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

// ErrIncompatible means a game saved by another engine version does not
// read back exactly here; it cannot go on (N-20).
var ErrIncompatible = errors.New("load: the saved game does not fit this engine")

// Resume loads a game saved by any engine version, for a server restart
// (N-20). The game must write back to the same bytes: a renamed,
// removed, or retyped field fails, and so does a new field the old game
// lacks, unless it has omitempty (omitzero for a struct). meaning is
// the saving engine's StateMeaning.
func Resume(data []byte, meaning int, sets ...CardSet) (*Game, error) {
	if meaning != StateMeaning {
		return nil, fmt.Errorf("%w: state meaning %d, this engine %d", ErrIncompatible, meaning, StateMeaning)
	}
	g, err := Load(data, sets...)
	if err != nil {
		return nil, err
	}
	again, err := g.Save()
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(again, data) {
		return nil, fmt.Errorf("%w: it writes back differently", ErrIncompatible)
	}
	return g, nil
}
