package rulesengine

import "strconv"

// DeckZone is the zone of a deck.
func DeckZone(d DeckKind) Zone { return Zone{Kind: ZoneDeck, Deck: d} }

// DiscardZone is the zone of a discard pile.
func DiscardZone(d DeckKind) Zone { return Zone{Kind: ZoneDiscard, Deck: d} }

// AddToDeck puts new cards on top of a deck, in order (the last one ends
// on top). It is for setup and tests.
func (g *Game) AddToDeck(d DeckKind, cards ...CardRef) {
	for _, c := range cards {
		g.Decks[d] = append(g.Decks[d], g.newObject(c, DeckZone(d), NoPlayer))
	}
}

// ShuffleDeck puts a deck in random order.
func (g *Game) ShuffleDeck(d DeckKind) {
	deck := g.Decks[d]
	g.RNG.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })
}

// drawTop removes the top card of a deck and returns its object.
// An empty deck is refilled from its shuffled discard first (R-ZONE-05).
// ok is false when both are empty.
func (g *Game) drawTop(d DeckKind) (id ObjectID, ok bool) {
	if len(g.Decks[d]) == 0 {
		g.refillDeck(d)
	}
	n := len(g.Decks[d])
	if n == 0 {
		return 0, false
	}
	id = g.Decks[d][n-1]
	g.Decks[d] = g.Decks[d][:n-1]
	return id, true
}

// refillDeck shuffles the discard pile into the empty deck (R-ZONE-05).
func (g *Game) refillDeck(d DeckKind) {
	for _, id := range g.Discards[d] {
		g.Decks[d] = append(g.Decks[d], g.move(id, DeckZone(d), NoPlayer))
	}
	g.Discards[d] = nil
	g.ShuffleDeck(d)
}

// discard puts an object on top of its deck's discard pile (R-ZONE-06)
// and returns the new object. Callers remove it from its old zone list.
func (g *Game) discard(id ObjectID, d DeckKind) ObjectID {
	nid := g.move(id, DiscardZone(d), NoPlayer)
	g.Discards[d] = append(g.Discards[d], nid)
	return nid
}

// TopOf returns the object on top of a slot: the one in play (R-ZONE-02).
func (s Slot) TopOf() (ObjectID, bool) {
	if len(s.Cards) == 0 {
		return 0, false
	}
	return s.Cards[len(s.Cards)-1], true
}

// deckNames are the names used in prompts and events.
var deckNames = [deckCount]string{"treasure", "loot", "monster", "room"}

func (d DeckKind) String() string { return deckNames[d] }

// DecksInPlay lists the decks this game uses (a deck or its discard has
// cards), in a fixed order: the options of "choose a deck".
func (g *Game) DecksInPlay() []DeckKind {
	var out []DeckKind
	for d := range deckCount {
		if len(g.Decks[d])+len(g.Discards[d]) > 0 {
			out = append(out, d)
		}
	}
	return out
}

// DeckTop returns up to n cards from the top of a deck, top first.
func (g *Game) DeckTop(d DeckKind, n int) []ObjectID {
	deck := g.Decks[d]
	n = min(n, len(deck))
	out := make([]ObjectID, 0, n)
	for i := len(deck) - 1; i >= len(deck)-n; i-- {
		out = append(out, deck[i])
	}
	return out
}

// SetDeckTop puts cards that are already on top of a deck back in a new
// order, order[0] on top. They stay the same objects: only the order
// changes.
func (g *Game) SetDeckTop(d DeckKind, order []ObjectID) {
	deck := g.Decks[d]
	rest := deck[:len(deck)-len(order)]
	for i := len(order) - 1; i >= 0; i-- {
		rest = append(rest, order[i])
	}
	g.Decks[d] = rest
}

// MillTop puts the top card of a deck into its discard.
func (g *Game) MillTop(d DeckKind) {
	id, ok := g.drawTop(d)
	if !ok {
		return
	}
	nid := g.discard(id, d)
	g.emit(Event{Kind: EvDiscarded, Player: NoPlayer, Object: nid, Card: g.Object(nid).Card, Text: d.String() + " deck"})
}

// DiscardTopToDeck puts the top card of a discard on top of its deck.
func (g *Game) DiscardTopToDeck(d DeckKind) {
	pile := g.Discards[d]
	if len(pile) == 0 {
		return
	}
	id := pile[len(pile)-1]
	g.Discards[d] = pile[:len(pile)-1]
	nid := g.move(id, DeckZone(d), NoPlayer)
	g.Decks[d] = append(g.Decks[d], nid)
	g.emit(Event{Kind: EvMovedToDeck, Player: NoPlayer, Object: nid, Card: g.Object(nid).Card, Text: d.String() + " deck"})
}

// HandToDeckTop puts a loot card from p's hand on top of the loot deck.
// Only p knows which card it was.
func (g *Game) HandToDeckTop(p PlayerID, id ObjectID) {
	g.Players[p].Hand = remove(g.Players[p].Hand, id)
	nid := g.move(id, DeckZone(LootDeck), NoPlayer)
	g.Decks[LootDeck] = append(g.Decks[LootDeck], nid)
	g.emit(Event{Kind: EvMovedToDeck, Player: p, Object: nid, Card: g.Object(nid).Card, Text: "loot deck", Private: true})
}

// GiveHandCard moves a loot card from one hand to another (R-MECH-37).
// The receiver learns the card; the giver knew it already.
func (g *Game) GiveHandCard(from, to PlayerID, id ObjectID) ObjectID {
	g.Players[from].Hand = remove(g.Players[from].Hand, id)
	nid := g.move(id, Zone{Kind: ZoneHand}, to)
	g.Players[to].Hand = append(g.Players[to].Hand, nid)
	g.emit(Event{Kind: EvGaveCard, Player: to, Object: nid, Card: g.Object(nid).Card, Text: strconv.Itoa(int(from)), Private: true})
	return nid
}

// LookAt shows cards to p only.
func (g *Game) LookAt(p PlayerID, ids ...ObjectID) {
	for _, id := range ids {
		g.emit(Event{Kind: EvLookedAt, Player: p, Object: id, Card: g.Object(id).Card, Private: true})
	}
}
