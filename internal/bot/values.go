package bot

import (
	_ "embed"
	"encoding/json"
	"fmt"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// valuesJSON is the measured value of each loot card and item (ADR 030,
// Decision 20): written by make bot-values (cmd/bot_values), never by
// hand. It is data: the bot's code names no card.
//
//go:embed values.json
var valuesJSON []byte

// valueShrink is how many samples the kind's mean counts as, when a
// card's value is shrunk towards it (Decision 20.8): a card measured on
// n positions counts n / (n + valueShrink) of its own value. A card
// with no sample gets its kind's mean.
const valueShrink = 8.0

// valueKind is the kind a card is valued as (Decision 20.8.1).
type valueKind int

// The kinds: a loot card in hand, an item with a ↷ ability, any other
// item.
const (
	kindLoot valueKind = iota
	kindActive
	kindPassive
	kindCount
)

// The kinds' names in the table file.
const (
	KindLoot    = "loot"
	KindActive  = "active"
	KindPassive = "passive"
)

// kindNames name the kinds in the table file.
var kindNames = [kindCount]string{KindLoot, KindActive, KindPassive}

// ValueEntry is one measured card of the table file: Key is the card's
// ref; its value in settledScore units, the positions it was measured
// on and the standard error of the value. Kind is KindLoot, KindActive
// or KindPassive. The file's keys are the field names.
type ValueEntry struct {
	Key   string
	Kind  string
	Value float64
	N     int
	SE    float64
}

// ValueFile is the table file: how it was measured, and its cards.
type ValueFile struct {
	Note string
	// Seed, Games and Samples are the run's flags: the seed, the arena
	// games the positions come from, the samples a card at most.
	Seed    uint64
	Games   int
	Samples int
	Cards   []ValueEntry
}

// ParseValues reads a table file.
func ParseValues(data []byte) (ValueFile, error) {
	var f ValueFile
	if err := json.Unmarshal(data, &f); err != nil {
		return ValueFile{}, fmt.Errorf("bot: card values: %w", err)
	}
	for _, e := range f.Cards {
		if _, ok := parseKind(e.Kind); !ok {
			return ValueFile{}, fmt.Errorf("bot: card values: %s has kind %q", e.Key, e.Kind)
		}
	}
	return f, nil
}

// parseKind is the kind named s.
func parseKind(s string) (valueKind, bool) {
	for k, name := range kindNames {
		if name == s {
			return valueKind(k), true
		}
	}
	return 0, false
}

// values is the table as the score reads it: each card's value shrunk
// to its kind's mean, by kind.
type values struct {
	cards [kindCount]map[string]float64
	means [kindCount]float64
}

// cardTable is the embedded table.
var cardTable = mustValues(valuesJSON)

// mustValues builds the table from a table file; a broken file is a
// broken build (TestValuesFile).
func mustValues(data []byte) *values {
	f, err := ParseValues(data)
	if err != nil {
		panic(err)
	}
	return newValues(f)
}

// newValues shrinks each card's value to its kind's mean by its samples
// (Decision 20.8). A kind with no measured card has the flat value.
func newValues(f ValueFile) *values {
	v := &values{}
	var sum [kindCount]float64
	var n [kindCount]int
	for _, e := range f.Cards {
		k, _ := parseKind(e.Kind)
		if e.N > 0 {
			sum[k] += e.Value
			n[k]++
		}
	}
	for k := range kindCount {
		v.cards[k] = map[string]float64{}
		v.means[k] = wItem
		if k == kindLoot {
			v.means[k] = wLoot
		}
		if n[k] > 0 {
			v.means[k] = sum[k] / float64(n[k])
		}
	}
	for _, e := range f.Cards {
		k, _ := parseKind(e.Kind)
		w := float64(e.N)
		v.cards[k][e.Key] = (w*e.Value + valueShrink*v.means[k]) / (w + valueShrink)
	}
	return v
}

// mean is the mean value of kind k: a card of it the table lacks, and
// a card the seat does not see.
func (v *values) mean(k valueKind) float64 { return v.means[k] }

// value is the value of object id, valued as kind k: its card's shrunk
// value, or k's mean when the table lacks it.
func (v *values) value(g *engine.Game, id engine.ObjectID, k valueKind) float64 {
	if x, ok := v.cards[k][tableKey(g, id)]; ok {
		return x
	}
	return v.means[k]
}

// itemKind is how an item in play is valued: active with a ↷ ability.
func itemKind(g *engine.Game, id engine.ObjectID) valueKind {
	if g.HasTap(id) {
		return kindActive
	}
	return kindPassive
}

// tableKey is the card of object id, as the table names it. It is the
// one place the bot reads a card's identity, and only to look it up in
// the measured table (TestBotRules allows it here alone): the score
// reads it for a card the seat sees, in its hand or in play.
func tableKey(g *engine.Game, id engine.ObjectID) string {
	return string(g.Object(id).Card)
}
