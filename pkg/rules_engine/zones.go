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

// turnOrderFrom lists every player in turn order, starting with p.
func (g *Game) turnOrderFrom(p PlayerID) []PlayerID {
	out := make([]PlayerID, 0, len(g.Players))
	for range g.Players {
		out = append(out, p)
		p = g.next(p)
	}
	return out
}

// DeckToBottom moves a card that is in a deck to its bottom; it stays the
// same object.
func (g *Game) DeckToBottom(d DeckKind, id ObjectID) {
	g.Decks[d] = append([]ObjectID{id}, remove(g.Decks[d], id)...)
}

// DestroyObject destroys an item, curse or soul that p controls: it goes
// to the discard of its deck, or outside the game if it has none
// (R-MECH-23). Eternal objects stay (R-ABIL-18). It reports whether the
// object was destroyed.
func (g *Game) DestroyObject(p PlayerID, id ObjectID) bool {
	o := g.Object(id)
	if o.Zone.Kind != ZoneInPlay || o.Controller != p || g.Eternal(id) {
		return false
	}
	if g.def(id).SoulWhenDestroyed && o.Role == RoleItem {
		o.Role, o.Charged, o.Counters = RoleSoul, false, nil // The Chest
		g.emit(Event{Kind: EvGainedSoul, Player: p, Object: id, Card: o.Card})
		return false
	}
	g.Players[p].InPlay = remove(g.Players[p].InPlay, id)
	var nid ObjectID
	if d, ok := g.kindOf(id).Deck(); ok {
		nid = g.discard(id, d)
	} else {
		nid = g.move(id, Zone{Kind: ZoneOutside}, NoPlayer)
	}
	g.emit(Event{Kind: EvDestroyed, Player: p, Object: nid, Card: g.Object(nid).Card, Prev: id})
	return true
}

// GainControl moves an item in play to p: "steal an item". It stays the
// same object; only its controller changes (R-MECH-38). A shop item
// leaves its slot, which is refilled.
func (g *Game) GainControl(p PlayerID, id ObjectID) {
	o := g.Object(id)
	if o.Zone.Slot == ShopSlot && o.Controller == NoPlayer {
		g.removeFromSlot(id)
		nid := g.move(id, Zone{Kind: ZoneInPlay}, p)
		g.enterAsItem(p, nid)
		g.emit(Event{Kind: EvGainedTreasure, Player: p, Object: nid, Card: o.Card, Prev: id})
		g.enqueue(Action{Kind: ActRefillSlots, Player: NoPlayer})
		return
	}
	from := o.Controller
	g.Players[from].InPlay = remove(g.Players[from].InPlay, id)
	o.Controller = p
	g.Players[p].InPlay = append(g.Players[p].InPlay, id)
	g.emit(Event{Kind: EvGainedTreasure, Player: p, Object: id, Card: o.Card, Text: "stolen"})
}

// ShopItems lists the items on top of the shop slots.
func (g *Game) ShopItems() []ObjectID {
	var out []ObjectID
	for _, s := range g.Shop {
		if top, ok := s.TopOf(); ok {
			out = append(out, top)
		}
	}
	return out
}

// DiscardMonster puts a monster in a slot into the monster discard; the
// slot is refilled later (RefillSlots).
func (g *Game) DiscardMonster(id ObjectID) {
	g.removeFromSlot(id)
	nid := g.discard(id, MonsterDeck)
	g.emit(Event{Kind: EvDiscarded, Player: NoPlayer, Object: nid, Card: g.Object(nid).Card, Prev: id})
}

// RefillSlots queues the refill of empty slots (R-SHOP-06).
func (g *Game) RefillSlots() { g.enqueue(Action{Kind: ActRefillSlots, Player: NoPlayer}) }

// enterAsItem makes a new object in play an item of p: charged, unless
// the card says it enters deactivated, with its starting counters.
func (g *Game) enterAsItem(p PlayerID, id ObjectID) {
	o, d := g.Object(id), g.def(id)
	o.Role, o.Charged = RoleItem, !d.EntersDeactivated
	if d.EntersWithCounters > 0 {
		o.addCounters("", d.EntersWithCounters)
	}
	g.Players[p].InPlay = append(g.Players[p].InPlay, id)
}

// DealDamageTo puts n damage aimed at t on the stack (R-MECH-15), for
// effects that pick the target on resolution.
func (g *Game) DealDamageTo(t Target, n int, controller PlayerID, source ObjectID) {
	if n <= 0 {
		return
	}
	g.push(StackItem{Kind: StackDamage, Controller: controller, Source: source, Amount: n, Label: "damage", Target: t})
}

// PayHP makes p pay n HP (R-MECH-44): it is lost, not damage. It
// reports whether p could pay.
func (g *Game) PayHP(p PlayerID, n int) bool {
	if g.PlayerHP(p) < n || g.Players[p].Dead {
		return false
	}
	g.Players[p].Damage += n
	g.emit(Event{Kind: EvPaid, Player: p, Amount: n, Text: "HP"})
	if g.PlayerHP(p) == 0 {
		g.push(StackItem{Kind: StackDeath, Controller: NoPlayer, Label: "death", Target: Target{Player: p, IsPlayer: true}})
	}
	return true
}

// RerollItem destroys an item; if it is destroyed, its controller gains
// 1 treasure. A shop item is replaced by the top treasure card instead
// (R-MECH-48).
func (g *Game) RerollItem(id ObjectID) {
	o := g.Object(id)
	if o.Zone.Kind != ZoneInPlay || o.Role != RoleItem {
		return
	}
	if o.Controller == NoPlayer && o.Zone.Slot == ShopSlot {
		g.removeFromSlot(id)
		nid := g.discard(id, TreasureDeck)
		g.emit(Event{Kind: EvDestroyed, Player: NoPlayer, Object: nid, Card: g.Object(nid).Card, Prev: id})
		g.RefillSlots()
		return
	}
	if p := o.Controller; g.DestroyObject(p, id) {
		g.enqueue(Action{Kind: ActGainTreasure, Player: p, Amount: 1})
	}
}

// Kind returns the card kind of an object.
func (g *Game) Kind(id ObjectID) CardKind { return g.kindOf(id) }

// DiscardShopItem puts a shop item into the treasure discard; the slot is
// refilled later (RefillSlots).
func (g *Game) DiscardShopItem(id ObjectID) {
	g.removeFromSlot(id)
	nid := g.discard(id, TreasureDeck)
	g.emit(Event{Kind: EvDiscarded, Player: NoPlayer, Object: nid, Card: g.Object(nid).Card, Prev: id})
	g.RefillSlots()
}

// SlotToDeckBottom puts the top card of a shop or monster slot on the
// bottom of its deck; the slot is refilled later (RefillSlots).
func (g *Game) SlotToDeckBottom(id ObjectID) {
	d, ok := g.kindOf(id).Deck()
	if !ok {
		return
	}
	g.removeFromSlot(id)
	nid := g.move(id, DeckZone(d), NoPlayer)
	g.Decks[d] = append([]ObjectID{nid}, g.Decks[d]...)
	g.emit(Event{Kind: EvMovedToDeck, Player: NoPlayer, Object: nid, Card: g.Object(nid).Card, Text: d.String() + " deck bottom", Prev: id})
}

// DiscardToDeckTop puts a card from a discard pile on top of its deck.
func (g *Game) DiscardToDeckTop(d DeckKind, id ObjectID) {
	g.Discards[d] = remove(g.Discards[d], id)
	nid := g.move(id, DeckZone(d), NoPlayer)
	g.Decks[d] = append(g.Decks[d], nid)
	g.emit(Event{Kind: EvMovedToDeck, Player: NoPlayer, Object: nid, Card: g.Object(nid).Card, Text: d.String() + " deck"})
}

// RevealTop shows the top card of a deck to everyone.
func (g *Game) RevealTop(d DeckKind) {
	if top := g.DeckTop(d, 1); len(top) > 0 {
		g.emit(Event{Kind: EvCardRevealed, Player: NoPlayer, Object: top[0], Card: g.Object(top[0]).Card, Text: d.String() + " deck"})
	}
}

// lootFromDiscard puts the top n cards of the loot discard into p's
// hand (Compost).
func (g *Game) lootFromDiscard(p PlayerID, n int) {
	for range n {
		pile := g.Discards[LootDeck]
		if len(pile) == 0 {
			return
		}
		id := pile[len(pile)-1]
		g.Discards[LootDeck] = pile[:len(pile)-1]
		nid := g.move(id, Zone{Kind: ZoneHand}, p)
		g.Players[p].Hand = append(g.Players[p].Hand, nid)
		g.emit(Event{Kind: EvLooted, Player: p, Object: nid, Card: g.Object(nid).Card, Text: "from the discard"})
	}
}

// CoverMonsterSlot puts the top card of the monster deck on top of
// monster slot i, covering the card there (R-ZONE-10).
func (g *Game) CoverMonsterSlot(i int) {
	id, ok := g.drawTop(MonsterDeck)
	if !ok {
		return
	}
	if top, had := g.Monsters[i].TopOf(); had {
		g.Object(top).Zone = Zone{Kind: ZoneCovered, Slot: MonsterSlot, Index: i}
	}
	g.putInSlot(id, MonsterSlot, i)
}

// PutOnDeck puts new copies of cards on top of a deck, the last on top.
// For tests and the sandbox.
func (g *Game) PutOnDeck(d DeckKind, cards ...CardRef) {
	for _, c := range cards {
		g.Decks[d] = append(g.Decks[d], g.newObject(c, DeckZone(d), NoPlayer))
	}
}

// KillObject kills a monster or other object with HP: HP to 0, its
// death on the stack (R-MECH-23).
func (g *Game) KillObject(id ObjectID) {
	if g.Object(id).Zone.Kind != ZoneInPlay || g.Eternal(id) {
		return
	}
	g.Object(id).Damage = g.def(id).HP + g.bonus(StatMonsterHP, NoPlayer, id)
	g.push(StackItem{Kind: StackDeath, Controller: NoPlayer, Label: "death", Target: Target{Object: id}})
}

// LoseCentsNow queues "p loses n¢" (R-MECH-42).
func (g *Game) LoseCentsNow(p PlayerID, n int) {
	g.enqueue(Action{Kind: ActLoseCents, Player: p, Amount: n})
}
