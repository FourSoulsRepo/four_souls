package rulesengine

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
