package rulesengine

import "slices"

// NumberKind is a kind of number Magic Marker can change (R-ABIL-39).
type NumberKind int

// The numbers the engine knows in effect text.
const (
	NumLoot   NumberKind = iota // "loot N"
	NumCents                    // "gain N¢"
	NumDamage                   // "deal N damage"
	NumRoll                     // a roll result in a roll table: "6: …"
)

// Number is one number in an object's effect text.
type Number struct {
	Kind  NumberKind `json:"kind"`
	Value int        `json:"value"`
}

// numberKeys name a number as the card says it: "loot 2".
var numberKeys = [...]Key{
	NumLoot:   NewKey("ask.number.loot"),
	NumCents:  NewKey("ask.number.cents"),
	NumDamage: NewKey("ask.number.damage"),
	NumRoll:   NewKey("ask.number.roll"),
}

// Label names the number as the card says it, as an option.
func (n Number) Label() Msg {
	k := numberKeys[NumRoll]
	if n.Kind >= 0 && int(n.Kind) < len(numberKeys) {
		k = numberKeys[n.Kind]
	}
	return k.Msg(N("n", n.Value))
}

// Marker changes a number of an object's effect text till end of turn
// (Magic Marker).
type Marker struct {
	Object ObjectID   `json:"object"`
	Kind   NumberKind `json:"kind"`
	From   int        `json:"from"`
	To     int        `json:"to"`
}

// NumbersOf lists the numbers in the effect text of an object's cards,
// once each, in the order they are printed.
func (g *Game) NumbersOf(id ObjectID) []Number {
	var out []Number
	add := func(n Number) {
		if n.Value > 0 && !slices.Contains(out, n) {
			out = append(out, n)
		}
	}
	var scan func(effects []Effect)
	scanTable := func(t RollTable) {
		for r, effects := range t {
			if len(effects) > 0 {
				add(Number{NumRoll, r + 1})
			}
			scan(effects)
		}
	}
	scan = func(effects []Effect) {
		for _, e := range effects {
			switch x := e.(type) {
			case lootEffect:
				add(Number{NumLoot, x.n})
			case gainCentsEffect:
				add(Number{NumCents, x.n})
			case damageEffect:
				add(Number{NumDamage, x.n})
			case rollEffect:
				scanTable(x.table)
			case rollTimesEffect:
				scanTable(x.t)
			case destroySelfThen:
				scan(x.then)
			case eachPlayerEffect:
				scan(x.effects)
			case askEffect:
				if x.roll != nil {
					scanTable(*x.roll)
				}
			}
		}
	}
	for _, card := range g.facesOf(id) {
		d, ok := g.cards.find(card)
		if !ok {
			continue
		}
		for _, ab := range d.Abilities {
			scan(ab.Effects)
			for _, m := range ab.Modes {
				scan(m.Effects)
			}
		}
	}
	return out
}

// MarkNumber changes a number of an object's effect text by delta till
// end of turn; it stays between 1 and 6.
func (g *Game) MarkNumber(id ObjectID, n Number, delta int) {
	to := min(max(g.number(id, n.Kind, n.Value)+delta, 1), 6)
	for i, m := range g.Markers {
		if m.Object == id && m.Kind == n.Kind && m.From == n.Value {
			g.Markers[i].To = to
			return
		}
	}
	g.Markers = append(g.Markers, Marker{Object: id, Kind: n.Kind, From: n.Value, To: to})
}

// number is n as changed by markers on src.
func (g *Game) number(src ObjectID, k NumberKind, n int) int {
	for _, m := range g.Markers {
		if m.Object == src && m.Kind == k && m.From == n {
			return m.To
		}
	}
	return n
}

// markedTable applies roll markers on src to a roll table: the printed
// result From now reads To. A range ("3-4") grows or shrinks at that
// end; a single result moves.
func (g *Game) markedTable(src ObjectID, t RollTable) RollTable {
	for _, m := range g.Markers {
		if m.Object != src || m.Kind != NumRoll || m.From == m.To {
			continue
		}
		from, to := m.From-1, m.To-1
		entry := t[from]
		if len(entry) == 0 {
			continue
		}
		same := func(i int) bool { return i >= 0 && i < 6 && len(t[i]) > 0 && &t[i][0] == &entry[0] }
		inside := same(from-1) && same(from+1)
		if inside {
			continue // a number inside a range is not printed
		}
		out := t // a copy: RollTable is an array
		lo, hi := min(from, to), max(from, to)
		grow := (to > from && same(from-1)) || (to < from && same(from+1))
		shrink := (to < from && same(from-1)) || (to > from && same(from+1))
		switch {
		case grow:
			for i := lo; i <= hi; i++ {
				out[i] = entry
			}
		case shrink:
			out[from] = nil
		default:
			out[from], out[to] = nil, entry
		}
		t = out
	}
	return t
}
