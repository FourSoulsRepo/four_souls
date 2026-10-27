package rulesengine

import (
	"slices"
	"strconv"
)

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
// The receiver learns the card; the giver knew it already, so both see
// it in the event and nobody else does (V-14).
func (g *Game) GiveHandCard(from, to PlayerID, id ObjectID) ObjectID {
	g.Players[from].Hand = remove(g.Players[from].Hand, id)
	nid := g.move(id, Zone{Kind: ZoneHand}, to)
	g.Players[to].Hand = append(g.Players[to].Hand, nid)
	g.emit(Event{
		Kind: EvGaveCard, Player: to, Object: nid, Card: g.Object(nid).Card,
		Text: strconv.Itoa(int(from)), Private: true, SeenBy: []PlayerID{from},
	})
	return nid
}

// LookAt shows cards to p only. A hand card's event names the hand's
// owner in By, so the look can say whose hand it is (19.71).
func (g *Game) LookAt(p PlayerID, ids ...ObjectID) {
	for _, id := range ids {
		o := g.Object(id)
		by := 0
		if o.Zone.Kind == ZoneHand {
			by = int(o.Controller) + 1
		}
		g.emit(Event{Kind: EvLookedAt, Player: p, Object: id, Card: o.Card, Private: true, By: by})
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

// KeepOrBottom is p's answer after looking at the top card id of deck
// d: it stays on top, or goes to the bottom. Everyone learns which;
// only p sees the card (19.90).
func (g *Game) KeepOrBottom(p PlayerID, d DeckKind, id ObjectID, bottom bool) {
	where := " deck top"
	if bottom {
		g.DeckToBottom(d, id)
		where = " deck bottom"
	}
	g.emit(Event{Kind: EvTopChoice, Player: p, Object: id, Card: g.Object(id).Card, Text: d.String() + where, Private: true})
}

// ToBottomRandom puts cards on the bottom of a deck in random order; no
// player orders them (Q-22). A card in a deck stays the same object; a
// card in a slot leaves it.
func (g *Game) ToBottomRandom(d DeckKind, ids []ObjectID) {
	ids = slices.Clone(ids)
	g.RNG.Shuffle(len(ids), func(i, j int) { ids[i], ids[j] = ids[j], ids[i] })
	for _, id := range ids {
		if g.Object(id).Zone.Kind == ZoneInPlay {
			g.SlotToDeckBottom(id)
		} else {
			g.DeckToBottom(d, id)
		}
	}
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
	if o.Zone.Slot == ShopSlot && o.Controller == NoPlayer {
		// A shop item: its abilities don't work in the shop (R-CARD-05),
		// so it is simply discarded and its slot refilled (R-SHOP-06).
		g.removeFromSlot(id)
		nid := g.discard(id, TreasureDeck)
		g.emit(Event{Kind: EvDestroyed, Player: NoPlayer, Object: nid, Card: g.Object(nid).Card, Prev: id})
		g.enqueue(Action{Kind: ActRefillSlots, Player: NoPlayer})
		return true
	}
	if f := g.def(id).InsteadOfDestroy; f != nil && f(g, id) {
		return false // replaced (Bag Lunch)
	}
	if g.def(id).SoulWhenDestroyed && o.Role == RoleItem && g.CanGainSouls(p) {
		o.Role, o.Charged, o.Counters, o.NoAbilities = RoleSoul, false, nil, true // The Chest
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
	if o.Zone.Kind != ZoneInPlay || o.Controller == p {
		return // it left play in the meantime, or is p's already
	}
	if o.Zone.Slot == ShopSlot && o.Controller == NoPlayer {
		g.removeFromSlot(id)
		nid := g.move(id, Zone{Kind: ZoneInPlay}, p)
		g.enterAsItem(p, nid)
		g.emit(Event{Kind: EvGainedTreasure, Player: p, Object: nid, Card: o.Card, Prev: id, Text: "stolen"})
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

// DiscardRoom puts a room in a room slot into discard; the slot is
// refilled (R-TURN-12, R-SHOP-06).
func (g *Game) DiscardRoom(id ObjectID) {
	if o := g.Object(id); o.Zone.Kind != ZoneInPlay || o.Zone.Slot != RoomSlot {
		return
	}
	if stay := g.def(id).NoDiscard; stay != nil && stay(g, id) {
		return
	}
	g.removeFromSlot(id)
	nid := g.discard(id, RoomDeck)
	g.emit(Event{Kind: EvDiscarded, Player: NoPlayer, Object: nid, Card: g.Object(nid).Card, Prev: id})
	g.RefillSlots()
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
	g.push(StackItem{Kind: StackDamage, Controller: controller, Source: source, Amount: n, Label: eventDamage.Msg(), Target: t})
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
		g.push(StackItem{Kind: StackDeath, Controller: NoPlayer, Label: eventDeath.Msg(), Target: Target{Player: p, IsPlayer: true}})
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
// monster slot i, covering the card there (R-ZONE-10); an event's
// abilities trigger (R-CARD-15).
func (g *Game) CoverMonsterSlot(i int) {
	id, ok := g.drawTop(MonsterDeck)
	if !ok {
		return
	}
	g.enterMonsterSlot(g.coverMonsterSlot(id, i))
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
	g.push(StackItem{Kind: StackDeath, Controller: NoPlayer, Label: eventDeath.Msg(), Target: Target{Object: id}})
}

// LoseCentsNow queues "p loses n¢" (R-MECH-42).
func (g *Game) LoseCentsNow(p PlayerID, n int) {
	g.enqueue(Action{Kind: ActLoseCents, Player: p, Amount: n})
}

// PlayLootFree makes p play the loot card id from their hand without
// using a loot play (Distant Admiration).
func (g *Game) PlayLootFree(p PlayerID, id ObjectID) {
	g.enqueue(Action{Kind: ActPlayLootFree, Player: p, Object: id})
}

// DiscardToHand puts a card from a discard pile into p's hand ("put a
// card named Pills! from the loot discard into your hand").
func (g *Game) DiscardToHand(d DeckKind, id ObjectID, p PlayerID) {
	if !contains(g.Discards[d], id) {
		return
	}
	g.Discards[d] = remove(g.Discards[d], id)
	nid := g.move(id, Zone{Kind: ZoneHand}, p)
	g.Players[p].Hand = append(g.Players[p].Hand, nid)
	g.emit(Event{Kind: EvLooted, Player: p, Object: nid, Card: g.Object(nid).Card, Text: "from the discard"})
}

// PlaceFromDiscard puts a monster card from the monster discard on top
// of monster slot i, covering the card there ("put them in monster
// slots": Dad's Note).
func (g *Game) PlaceFromDiscard(id ObjectID, i int) {
	if !contains(g.Discards[MonsterDeck], id) || i < 0 || i >= len(g.Monsters) {
		return
	}
	g.Discards[MonsterDeck] = remove(g.Discards[MonsterDeck], id)
	g.enterMonsterSlot(g.coverMonsterSlot(id, i))
}

// AddCountersTo puts n generic counters on an object in play, or removes
// them if n is negative.
func (g *Game) AddCountersTo(id ObjectID, n int) { g.AddNamedCounters(id, "", n) }

// AddNamedCounters puts n counters named name on an object in play ("a
// gold counter", "levels up"), or removes them if n is negative.
func (g *Game) AddNamedCounters(id ObjectID, name string, n int) {
	o := g.Object(id)
	if (o.Zone.Kind != ZoneInPlay && !contains(g.BonusSouls, id)) || n == 0 {
		return // in play, or an active bonus soul (Soul of Lust)
	}
	before := o.CountersOf(name)
	o.addCounters(name, n)
	if d := o.CountersOf(name) - before; d != 0 {
		g.emit(Event{Kind: EvCounters, Player: o.Controller, Object: id, Card: o.Card, Amount: d, Text: name})
	}
}

// PutUnder puts a monster in a slot under holder: it leaves its slot,
// which is refilled, and is outside the game (R-ZONE-13: Friendly Ball).
func (g *Game) PutUnder(holder, id ObjectID) {
	o := g.Object(id)
	if o.Zone.Kind != ZoneInPlay || o.Zone.Slot != MonsterSlot || g.Object(holder).Zone.Kind != ZoneInPlay {
		return
	}
	g.removeFromSlot(id)
	nid := g.move(id, Zone{Kind: ZoneOutside}, NoPlayer)
	h := g.Object(holder)
	h.Under = append(h.Under, nid)
	g.emit(Event{Kind: EvMovedToDeck, Player: h.Controller, Object: nid, Card: o.Card, Prev: id, Text: "under " + string(h.Card)})
	g.RefillSlots()
}

// FromUnder puts a card from under holder on top of monster slot i,
// covering the card there; it enters play.
func (g *Game) FromUnder(holder, id ObjectID, i int) ObjectID {
	h := g.Object(holder)
	if !contains(h.Under, id) || i < 0 || i >= len(g.Monsters) {
		return 0
	}
	h.Under = remove(h.Under, id)
	nid := g.coverMonsterSlot(id, i)
	g.enterMonsterSlot(nid)
	return nid
}

// DeckToHand takes a card out of a deck into p's hand; only p sees it.
func (g *Game) DeckToHand(d DeckKind, id ObjectID, p PlayerID) {
	if !contains(g.Decks[d], id) {
		return
	}
	g.Decks[d] = remove(g.Decks[d], id)
	nid := g.move(id, Zone{Kind: ZoneHand}, p)
	g.Players[p].Hand = append(g.Players[p].Hand, nid)
	g.emit(Event{Kind: EvLooted, Player: p, Object: nid, Card: g.Object(nid).Card, Private: true})
}

// TurnOrderFrom lists every player in turn order, starting with p.
func (g *Game) TurnOrderFrom(p PlayerID) []PlayerID { return g.turnOrderFrom(p) }

// coverMonsterSlot puts a card on top of monster slot i, covering the
// card there (R-ZONE-10). An Indomitable card, or one that would cover
// an Indomitable monster, goes into a new slot instead (R-ABIL-23).
func (g *Game) coverMonsterSlot(id ObjectID, i int) ObjectID {
	top, had := g.Monsters[i].TopOf()
	if g.def(id).Indomitable || (had && g.Object(top).Zone.Kind == ZoneInPlay && g.def(top).Indomitable) {
		return g.putInSlot(id, MonsterSlot, g.newMonsterSlot())
	}
	if had {
		o := g.Object(top)
		o.Zone, o.Role = Zone{Kind: ZoneCovered, Slot: MonsterSlot, Index: i}, RoleNone
	}
	return g.putInSlot(id, MonsterSlot, i)
}

// newMonsterSlot expands monster slots by 1 and returns the new slot
// (R-MECH-31).
func (g *Game) newMonsterSlot() int {
	g.Monsters = append(g.Monsters, Slot{})
	g.emit(Event{Kind: EvExpanded, Player: NoPlayer, Amount: 1, Text: "monster"})
	return len(g.Monsters) - 1
}

// GiveCents makes from give up to n¢ to another player and returns how
// much was given; giving is not losing (R-MECH-35).
func (g *Game) GiveCents(from, to PlayerID, n int) int {
	n = min(n, g.Players[from].Cents)
	if n <= 0 || from == to {
		return 0
	}
	g.Players[from].Cents -= n
	g.Players[to].Cents += n
	g.emit(Event{Kind: EvGaveCents, Player: from, Amount: n, By: int(to) + 1})
	return n
}
