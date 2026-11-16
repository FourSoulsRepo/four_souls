package rulesengine

import (
	"fmt"
	"reflect"
	"slices"
)

// Known is what a seat remembers of cards it cannot see now (ADR 030,
// Decision 2.7 – 2.9 and 19). The zero Known remembers nothing.
type Known struct {
	// Hands are the cards remembered in each player's hand, by seat. The
	// seat's own hand is never hidden from it, so its entry does nothing.
	Hands [][]CardRef
	// Tops are the cards remembered on each deck, top first; Bottoms
	// from the bottom up. "" is a place the seat does not remember.
	Tops    [deckCount][]CardRef
	Bottoms [deckCount][]CardRef
}

// Imagine returns a copy of the game as seat p may imagine it (ADR 030,
// Decision 2): every card p cannot see now is pooled with the hidden
// cards of its kind, shuffled and dealt back, so each zone keeps its
// size and p's view of the copy is p's view of the game. Cards in
// known are dealt where p remembers them, when such a card is among
// the hidden ones; else it is skipped. The copy's RNG starts from seed:
// its dice and draws are not the game's. The game is not changed.
//
// What is hidden is listed in hiddenFrom; a new place that hides a card
// from a seat must be added there (docs/wiki/Rules-engine.md, "Imagined
// game").
func (g *Game) Imagine(p PlayerID, seed uint64, known Known) (*Game, error) {
	if p < 0 || int(p) >= len(g.Players) {
		return nil, fmt.Errorf("imagine: no player %d", p)
	}
	c, err := g.Clone()
	if err != nil {
		return nil, err
	}
	rng := NewRNG(seed)
	slots := g.hiddenFrom(p)
	dealt := g.dealHidden(slots, known, &rng)
	// relabel says which card stands in for a hidden card in the copy, so
	// the copy's past and its other players' prompts name no real one.
	relabel := map[CardRef]CardRef{}
	for i, s := range slots {
		was := g.Objects[s.id].Card
		c.Objects[s.id].Card = dealt[i]
		if _, ok := relabel[was]; !ok {
			relabel[was] = dealt[i]
		}
	}
	c.relabelHidden(p, relabel)
	c.RNG = rng
	return c, nil
}

// hiddenSlot is one object whose card a seat cannot see.
type hiddenSlot struct {
	id ObjectID
	// pool is the card's deck, or deckCount plus its kind for a card of
	// no deck: cards are redealt only within their pool.
	pool int
	// hand is the player whose hand holds it, or NoPlayer; deck and pos
	// (from the bottom) place a deck card, pos is -1 for others.
	hand PlayerID
	deck DeckKind
	pos  int
}

// hiddenFrom lists the objects whose card seat p cannot see now, in a
// fixed order: deck cards bottom up, deck by deck; the other players'
// hands in seat order; face-down presents. A card p looks at now stays
// out of the list (seenBy).
func (g *Game) hiddenFrom(p PlayerID) []hiddenSlot {
	seen := g.seenBy(p)
	var out []hiddenSlot
	for d := range deckCount {
		for i, id := range g.Decks[d] {
			if !seen[id] {
				out = append(out, hiddenSlot{id: id, pool: int(d), hand: NoPlayer, deck: d, pos: i})
			}
		}
	}
	for _, pl := range g.Players {
		if g.seesHand(p, pl.ID) {
			continue
		}
		for _, id := range pl.Hand {
			if !seen[id] {
				out = append(out, hiddenSlot{id: id, pool: g.poolOf(g.Objects[id].Card), hand: pl.ID, pos: -1})
			}
		}
	}
	if ch := g.Challenge; ch != nil {
		for _, pr := range ch.Presents {
			// A live present lies face down outside the game (R-KRAMPUS-05).
			if !pr.Destroyed && !pr.Flipped && g.Objects[pr.Object].Zone.Kind == ZoneOutside && !seen[pr.Object] {
				out = append(out, hiddenSlot{id: pr.Object, pool: g.poolOf(g.Objects[pr.Object].Card), hand: NoPlayer, pos: -1})
			}
		}
	}
	return out
}

// seesHand: p's view shows q's hand (ViewFor): p's own, a revealed
// hand, or any hand of a co-op challenge (R-COOP-07).
func (g *Game) seesHand(p, q PlayerID) bool {
	return p == q || g.bonus(StatRevealHand, q, 0) > 0 || g.Challenge != nil
}

// seenBy lists the deck and hand cards p looks at now, as ViewFor shows
// them: the treasure peek, and the cards of p's open question: its
// options, subject, deck look and order.
func (g *Game) seenBy(p PlayerID) map[ObjectID]bool {
	seen := map[ObjectID]bool{}
	if g.peeksTreasure(p) {
		for _, id := range g.DeckTop(TreasureDeck, 1) {
			seen[id] = true
		}
	}
	w, ch := g.Waiting, g.Choice
	if i := g.startIndex(p); i >= 0 {
		w, ch = g.startPrompt(i), &g.StartChoices[i]
	}
	if ch == nil || w.Player != p || (w.Kind != PromptChoose && w.Kind != PromptOrder) {
		return seen
	}
	for _, id := range ch.Objects {
		seen[id] = true
	}
	for _, r := range ch.Replacements {
		seen[r.Object] = true
	}
	for _, t := range ch.Targets {
		seen[t.Object] = true
	}
	if ch.Subject != 0 {
		seen[ch.Subject] = true
	}
	// A deck look without its cards in Objects is a look at the top
	// (Question.Looks, 20.18).
	if l := ch.Look; l != nil && len(ch.Objects) == 0 {
		for _, id := range g.DeckTop(l.Deck, l.Count) {
			seen[id] = true
		}
	}
	return seen
}

// poolOf is the pool of a card: its deck, or deckCount plus its kind.
// A card with no definition gets a pool of its own (-1).
func (g *Game) poolOf(ref CardRef) int {
	d, ok := g.cards.find(ref)
	if !ok {
		return -1
	}
	if deck, ok := d.Kind.Deck(); ok {
		return int(deck)
	}
	return int(deckCount) + int(d.Kind)
}

// dealHidden returns the card each slot gets: first the known cards,
// each where p remembers it while its pool still holds it, then the
// rest of each pool, shuffled, in slot order.
func (g *Game) dealHidden(slots []hiddenSlot, known Known, rng *RNG) []CardRef {
	pools := map[int][]CardRef{}
	var keys []int
	for _, s := range slots {
		if _, ok := pools[s.pool]; !ok {
			keys = append(keys, s.pool)
		}
		pools[s.pool] = append(pools[s.pool], g.Objects[s.id].Card)
	}
	dealt := make([]CardRef, len(slots))
	// place deals ref to slot i, if the slot is free and ref is still in
	// the slot's pool.
	place := func(i int, ref CardRef) bool {
		if i < 0 || ref == "" || dealt[i] != "" || g.poolOf(ref) != slots[i].pool {
			return false
		}
		pool := pools[slots[i].pool]
		j := slices.Index(pool, ref)
		if j < 0 {
			return false
		}
		pools[slots[i].pool] = slices.Delete(pool, j, j+1)
		dealt[i] = ref
		return true
	}
	deckSlot := func(d DeckKind, pos int) int {
		return slices.IndexFunc(slots, func(s hiddenSlot) bool { return s.hand == NoPlayer && s.pos == pos && s.deck == d })
	}
	for q, hand := range known.Hands {
		for _, ref := range hand {
			for i, s := range slots {
				if s.hand == PlayerID(q) && place(i, ref) {
					break
				}
			}
		}
	}
	for d := range deckCount {
		n := len(g.Decks[d])
		for k, ref := range known.Tops[d] {
			place(deckSlot(d, n-1-k), ref)
		}
		for k, ref := range known.Bottoms[d] {
			place(deckSlot(d, k), ref)
		}
	}
	slices.Sort(keys)
	for _, k := range keys {
		pool := pools[k]
		rng.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
		for i, s := range slots {
			if s.pool == k && dealt[i] == "" {
				dealt[i], pool = pool[0], pool[1:]
			}
		}
	}
	return dealt
}

// relabelHidden renames hidden cards where the copy still names them
// after the deal: the past objects of the copy (a card drawn unseen
// left its old object behind) and the options of another player's
// prompt, which p's view does not show. A past object the state still
// points at keeps its card: the stack or a boost may show it (a
// trigger of a monster that left play).
func (g *Game) relabelHidden(p PlayerID, relabel map[CardRef]CardRef) {
	used := g.objectRefs()
	for i := range g.Objects {
		o := &g.Objects[i]
		if o.Zone.Kind != ZoneGone || used[o.ID] {
			continue
		}
		if r, ok := relabel[o.Card]; ok {
			o.Card = r
		}
	}
	if g.Waiting.Kind == PromptChoose && g.Waiting.Player != p {
		for i := range g.Waiting.Options {
			args := g.Waiting.Options[i].Args
			for j, a := range args {
				if v, ok := a.Value.(CardValue); ok {
					if r, ok := relabel[v.Card]; ok {
						args[j].Value = CardValue{Card: r}
					}
				}
			}
		}
	}
}

// objectRefs lists every object the state points at outside the object
// list: zones, the stack, triggers, boosts, choices and the rest. It
// walks the state's types, so a new field is covered without a line here.
func (g *Game) objectRefs() []bool {
	out := make([]bool, len(g.Objects))
	v := reflect.ValueOf(g).Elem()
	t := v.Type()
	for i := range t.NumField() {
		switch t.Field(i).Name {
		case "Objects", "cards", "events", "asks":
			// The objects themselves; card definitions; the last step's
			// events, which a copy drops; ask effects, functions only.
			continue
		}
		addObjectRefs(v.Field(i), out)
	}
	return out
}

var objectIDType = reflect.TypeFor[ObjectID]()

// addObjectRefs marks every ObjectID inside v in out.
func addObjectRefs(v reflect.Value, out []bool) {
	if v.Type() == objectIDType {
		if id := v.Int(); id >= 0 && id < int64(len(out)) {
			out[id] = true
		}
		return
	}
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface:
		if !v.IsNil() {
			addObjectRefs(v.Elem(), out)
		}
	case reflect.Struct:
		for i := range v.NumField() {
			addObjectRefs(v.Field(i), out)
		}
	case reflect.Slice, reflect.Array:
		if v.Type().Elem().Kind() == reflect.String {
			return
		}
		for i := range v.Len() {
			addObjectRefs(v.Index(i), out)
		}
	case reflect.Map:
		it := v.MapRange()
		for it.Next() {
			addObjectRefs(it.Key(), out)
			addObjectRefs(it.Value(), out)
		}
	default:
	}
}
