package rulesengine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"slices"
)

// SituationCounter is a number of counters on one card (23.3). Card is
// the card's ref, or "character" for a player's character; Name names
// the counters ("gold", a level), empty for plain counters.
type SituationCounter struct {
	Card  CardRef `json:"card"`
	Name  string  `json:"name,omitempty"`
	Count int     `json:"count"`
}

// UnmarshalJSON rejects unknown fields, so a misspelt field is an error,
// not a counter of 0.
func (c *SituationCounter) UnmarshalJSON(data []byte) error {
	type plain SituationCounter
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var p plain
	if err := dec.Decode(&p); err != nil {
		return fmt.Errorf("counter (fields card, name, count): %w", err)
	}
	*c = SituationCounter(p)
	return nil
}

// setUpPlayer applies a seat's counters, eternal cards and extra loot
// plays once its cards are built; seat is its index, as in errors.
func (g *Game) setUpPlayer(pl *Player, seat int, sp SituationPlayer) error {
	if sp.ExtraLootPlays < 0 {
		return fmt.Errorf("situation: player %d extra_loot_plays %d is below 0", seat, sp.ExtraLootPlays)
	}
	pl.ExtraLootPlays = sp.ExtraLootPlays
	// refs are the seat's cards in play by ref; "character" first.
	refs := []CardRef{"character"}
	ids := []ObjectID{pl.Character}
	for _, id := range pl.InPlay {
		refs, ids = append(refs, g.Object(id).Card), append(ids, id)
	}
	find := func(field string, ref CardRef, from int) (ObjectID, error) {
		if i := slices.Index(refs[from:], ref); i >= 0 {
			return ids[from+i], nil
		}
		return 0, fmt.Errorf("situation: player %d %s: %q is not in play, options %v", seat, field, ref, refs[from:])
	}
	for _, ref := range sp.Eternal {
		id, err := find("eternal", ref, 1) // an item or soul, not the character
		if err != nil {
			return err
		}
		g.Object(id).Eternal = true
	}
	for _, c := range sp.Counters {
		id, err := find("counters", c.Card, 0)
		if err != nil {
			return err
		}
		if err := putCounters(g.Object(id), c); err != nil {
			return fmt.Errorf("situation: player %d counters: %w", seat, err)
		}
	}
	return nil
}

func putCounters(o *Object, c SituationCounter) error {
	if c.Count < 1 {
		return fmt.Errorf("%q count %d is below 1", c.Card, c.Count)
	}
	o.addCounters(c.Name, c.Count)
	return nil
}

// setUpDecks puts the setup's cards on top of their decks, top first.
func (g *Game) setUpDecks(decks map[string][]CardRef) error {
	for name := range decks {
		if !slices.Contains(deckNames[:], name) {
			return fmt.Errorf("situation: decks: %q is no deck, options %v", name, deckNames)
		}
	}
	for d := range deckCount {
		refs := decks[d.String()]
		for _, ref := range refs {
			def, ok := g.cards.find(ref)
			if !ok {
				return fmt.Errorf("situation: decks %s: unknown card %q", d, ref)
			}
			if k, ok := def.Kind.Deck(); !ok || k != d {
				return fmt.Errorf("situation: decks %s: %q is no card of the %s deck", d, ref, d)
			}
		}
		for i := len(refs) - 1; i >= 0; i-- {
			g.PutOnDeck(d, refs[i])
		}
	}
	return nil
}

// setUpSlots puts damage and counters on cards in slots and on bonus
// souls. A ref names the first such card, in slot order.
func (g *Game) setUpSlots(s SituationSetup) error {
	var refs []CardRef
	var ids []ObjectID
	for _, rows := range append(g.MonsterRows(), g.Rooms) {
		for _, slot := range rows {
			if id, ok := slot.TopOf(); ok {
				refs, ids = append(refs, g.Object(id).Card), append(ids, id)
			}
		}
	}
	slotted := len(refs) // bonus souls take counters, not damage
	for _, id := range g.BonusSouls {
		refs, ids = append(refs, g.Object(id).Card), append(ids, id)
	}
	damaged := make([]CardRef, 0, len(s.MonsterDamage))
	for ref := range s.MonsterDamage {
		damaged = append(damaged, ref)
	}
	slices.Sort(damaged)
	for _, ref := range damaged {
		i := slices.Index(refs[:slotted], ref)
		if i < 0 {
			return fmt.Errorf("situation: monster_damage: %q is in no slot, options %v", ref, refs[:slotted])
		}
		n, hp := s.MonsterDamage[ref], g.objectMaxHP(ids[i])
		if n < 0 || n >= hp {
			return fmt.Errorf("situation: monster_damage: %q has %d HP; damage %d is not 0 to %d", ref, hp, n, hp-1)
		}
		g.Object(ids[i]).Damage = n
	}
	for _, c := range s.Counters {
		i := slices.Index(refs, c.Card)
		if i < 0 {
			return fmt.Errorf("situation: counters: %q is in no slot and no bonus soul, options %v", c.Card, refs)
		}
		if err := putCounters(g.Object(ids[i]), c); err != nil {
			return fmt.Errorf("situation: counters: %w", err)
		}
	}
	return nil
}

// setUpTurn applies the turn's overrides. loot_plays is unchecked, as
// before 23.3, so old payloads keep their answers.
func (s SituationSetup) setUpTurn(t *Turn) error {
	if s.LootPlays != nil {
		t.LootPlays = *s.LootPlays
	}
	for _, o := range []struct {
		name string
		n    *int
		to   *int
	}{{"attacks", s.Attacks, &t.Attacks}, {"purchases", s.Purchases, &t.Purchases}} {
		if o.n == nil {
			continue
		}
		if *o.n < 0 {
			return fmt.Errorf("situation: %s %d is below 0", o.name, *o.n)
		}
		*o.to = *o.n
	}
	return nil
}

// setUpRolls forces the next dice results.
func (g *Game) setUpRolls(rolls []int) error {
	for i, r := range rolls {
		if r < 1 || r > 6 {
			return fmt.Errorf("situation: rolls: roll %d is %d, not 1 to 6", i+1, r)
		}
	}
	g.ForceRolls(rolls...)
	return nil
}

// openSetupWindow opens the first window: for holder, the top item's owner
// or the active player, unless the setup names the seat (23.4). With an
// empty stack, the players from the active player up to that seat have
// passed (R-PRIO-02), so when it and the rest pass the window closes as
// in a game (R-PRIO-05).
func (g *Game) openSetupWindow(seat *int, holder PlayerID) error {
	if seat == nil {
		g.openWindow(holder)
		return nil
	}
	n := len(g.Players)
	if *seat < 0 || *seat >= n {
		return fmt.Errorf("situation: priority %d is not a seat, options 0 to %d", *seat, n-1)
	}
	g.openWindow(PlayerID(*seat))
	if len(g.Stack) == 0 {
		g.Priority.Passes = (*seat - int(g.Turn.Active) + n) % n
	}
	return nil
}
