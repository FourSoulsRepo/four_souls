package bot

import (
	_ "embed"
	"encoding/json"
	"fmt"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// valuesJSON is the measured value of each loot card and item (ADR 033,
// Decision 20): written by make bot-values (cmd/bot_values), never by
// hand. It is data: the bot's code names no card.
//
//go:embed values.json
var valuesJSON []byte

// tauFloor2 is the least spread τ² the cards of a kind are taken to
// have (Decision 20.8.4): when their measured spread is no more than
// their noise, every card is pulled nearly to the kind's mean, never
// divided by 0.
const tauFloor2 = 0.01

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

// How a value was measured (Decision 20.4): a loot card's play, an
// item's playouts, an active item's uses.
const (
	MeasurePlay    = "play"
	MeasurePlayout = "playout"
	MeasureUse     = "use"
)

// ValueEntry is one measured card of the table file: Key is the card's
// ref; its value in settledScore units, the positions it was measured
// on and the standard error of the value; Measure is how. Kind is
// KindLoot, KindActive or KindPassive. The file's keys are the field
// names.
type ValueEntry struct {
	Key     string
	Kind    string
	Measure string
	Value   float64
	N       int
	SE      float64
}

// ValueFile is the table file: how it was measured, and its cards.
type ValueFile struct {
	Note string
	// Seed, Games, Samples and ItemSamples are the run's flags: the
	// seed, the arena games of each size the positions come from, the
	// samples a loot card and an item at most.
	Seed        uint64
	Games       int
	Samples     int
	ItemSamples int
	Cards       []ValueEntry
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
// to its kind's mean, by kind; measured items scaled so their mean is
// wItem; a kind with no measured card at its flat value, unscaled.
type values struct {
	cards [kindCount]map[string]float64
	means [kindCount]float64
	// tau2 is each kind's τ²: the spread of its cards' true values.
	tau2 [kindCount]float64
	// scale is what item values are multiplied by: wItem over the mean
	// of every measured item (Decision 20.8.5).
	scale float64
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

// newValues builds the table (Decision 20.8). Each card's value is
// weighed against its kind's mean by precision: v' = (v/se² + m/τ²) /
// (1/se² + 1/τ²), τ² the spread of the kind's values less their mean
// se², at least tauFloor2. A card with no sample gets the mean; one
// with no error keeps its value. A kind with no measured card has the
// flat value. Items are then scaled by wItem over the mean of every
// measured item: they were measured over one round, but stay for the
// game.
func newValues(f ValueFile) *values {
	v := &values{scale: 1}
	var vals, ses [kindCount][]float64
	var measured [kindCount]bool // the kind has a card with a sample
	items, itemSum := 0, 0.0
	for _, e := range f.Cards {
		k, _ := parseKind(e.Kind)
		if e.N == 0 {
			continue
		}
		measured[k] = true
		vals[k] = append(vals[k], e.Value)
		ses[k] = append(ses[k], e.SE*e.SE)
		if k != kindLoot {
			items++
			itemSum += e.Value
		}
	}
	for k := range kindCount {
		v.cards[k] = map[string]float64{}
		v.means[k] = wItem
		if k == kindLoot {
			v.means[k] = wLoot
		}
		v.tau2[k] = tauFloor2
		if n := len(vals[k]); n > 0 {
			m := mean(vals[k])
			v.means[k] = m
			if n > 1 {
				ss := 0.0
				for _, x := range vals[k] {
					ss += (x - m) * (x - m)
				}
				v.tau2[k] = max(ss/float64(n-1)-mean(ses[k]), tauFloor2)
			}
		}
	}
	if items > 0 && itemSum > 0 {
		v.scale = wItem * float64(items) / itemSum
	}
	for _, e := range f.Cards {
		k, _ := parseKind(e.Kind)
		m, x := v.means[k], e.Value
		switch {
		case e.N == 0:
			x = m
		case e.SE > 0:
			p := 1 / (e.SE * e.SE)
			x = (x*p + m/v.tau2[k]) / (p + 1/v.tau2[k])
		}
		if k != kindLoot && measured[k] {
			x *= v.scale
		}
		v.cards[k][e.Key] = x
	}
	for k := kindActive; k < kindCount; k++ {
		if measured[k] { // a kind with no measure keeps the flat value (Decision 20.8)
			v.means[k] *= v.scale
		}
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
