package rulesengine

import "strconv"

// Events, curses and other monster-deck mechanics (step 5.6).

// enterMonsterSlot runs what happens when a card enters a monster slot:
// an event's abilities trigger and it then goes to discard (R-CARD-15);
// a curse is given to a player (R-ABIL-20).
func (g *Game) enterMonsterSlot(id ObjectID) {
	if g.Object(id).Role != RoleEvent {
		return
	}
	before := len(g.PendingTriggers)
	g.emit(Event{Kind: EvEnteredPlay, Player: NoPlayer, Object: id, Card: g.Object(id).Card})
	if g.def(id).Curse {
		g.enqueue(Action{Kind: ActGiveCurse, Player: g.Turn.Active, Object: id})
		return
	}
	if len(g.PendingTriggers) == before {
		g.finishEvent(id)
	}
}

// finishEvent puts an event whose abilities are done into discard; its
// slot is refilled (R-CARD-15, R-SHOP-06).
func (g *Game) finishEvent(id ObjectID) {
	o := g.Object(id)
	if o.Zone.Kind != ZoneInPlay || o.Zone.Slot != MonsterSlot || o.Role != RoleEvent || g.def(id).Curse {
		return
	}
	g.removeFromSlot(id)
	g.discard(id, MonsterDeck)
	g.RefillSlots()
}

// askCurseTarget: the active player picks who gains the curse.
func (g *Game) askCurseTarget(curse ObjectID) {
	if g.Object(curse).Zone.Kind != ZoneInPlay {
		return
	}
	var labels []string
	var players []int
	for _, pl := range g.Players {
		players = append(players, int(pl.ID))
		labels = append(labels, "player "+strconv.Itoa(int(pl.ID)+1)+" ("+string(g.Object(pl.Character).Card)+")")
	}
	g.ask(Choice{Purpose: ChooseCursed, Player: g.Turn.Active, Rule: "R-ABIL-20", Objects: []ObjectID{curse}, Slots: players}, labels)
}

// giveCurse moves a curse from its slot to player p.
func (g *Game) giveCurse(curse ObjectID, p PlayerID) {
	g.removeFromSlot(curse)
	nid := g.move(curse, Zone{Kind: ZoneInPlay}, p)
	g.Object(nid).Role = RoleCurse
	g.Players[p].InPlay = append(g.Players[p].InPlay, nid)
	g.emit(Event{Kind: EvCursed, Player: p, Object: nid, Card: g.Object(nid).Card})
	g.RefillSlots()
}

// dropCurses puts a dead player's curses into discard (R-ABIL-20).
func (g *Game) dropCurses(p PlayerID) {
	for _, id := range append([]ObjectID(nil), g.Players[p].InPlay...) {
		if g.Object(id).Role == RoleCurse {
			g.Players[p].InPlay = remove(g.Players[p].InPlay, id)
			nid := g.discard(id, MonsterDeck)
			g.emit(Event{Kind: EvDiscarded, Player: p, Object: nid, Card: g.Object(nid).Card, Prev: id})
		}
	}
}

// ExpandShop adds n shop slots; ExpandMonsters adds n monster slots.
// New slots are filled the next time a player would get priority.
func (g *Game) ExpandShop(n int) {
	for range n {
		g.Shop = append(g.Shop, Slot{})
	}
	g.emit(Event{Kind: EvExpanded, Player: NoPlayer, Amount: n, Text: "shop"})
	g.RefillSlots()
}

// ExpandMonsters adds n monster slots (R-MECH-31).
func (g *Game) ExpandMonsters(n int) {
	for range n {
		g.Monsters = append(g.Monsters, Slot{})
	}
	g.emit(Event{Kind: EvExpanded, Player: NoPlayer, Amount: n, Text: "monster"})
	g.RefillSlots()
}

// ForceAttacks gives the active player n more attacks they must make;
// deck makes them attacks on the monster deck ("must attack the monster
// deck 2 times this turn").
func (g *Game) ForceAttacks(n int, deck bool) {
	g.Turn.Attacks += n
	g.Turn.MustAttacks += n
	if deck {
		g.Turn.MustAttackDeck += n
	}
}

// TakeFromDeck takes a card out of a deck and puts it into play as an
// item of p ("search the treasure deck for …, gain it").
func (g *Game) TakeFromDeck(d DeckKind, id ObjectID, p PlayerID) {
	g.Decks[d] = remove(g.Decks[d], id)
	nid := g.move(id, Zone{Kind: ZoneInPlay}, p)
	g.enterAsItem(p, nid)
	g.emit(Event{Kind: EvGainedTreasure, Player: p, Object: nid, Card: g.Object(nid).Card, Text: "searched"})
}

// HealPlayer removes up to n damage from a living player (R-MECH-40).
func (g *Game) HealPlayer(p PlayerID, n int) {
	pl := &g.Players[p]
	if pl.Dead || n <= 0 {
		return
	}
	h := min(n, pl.Damage)
	pl.Damage -= h
	g.emit(Event{Kind: EvHealed, Player: p, Amount: h})
}

// HealObject removes up to n damage from a monster or object.
func (g *Game) HealObject(id ObjectID, n int) {
	o := g.Object(id)
	h := min(n, o.Damage)
	o.Damage -= h
	g.emit(Event{Kind: EvHealed, Player: NoPlayer, Object: id, Card: o.Card, Amount: h})
}

// AddDeckAttack lets the active player attack the monster deck one more
// time this turn ("may attack the monster deck an additional time").
func (g *Game) AddDeckAttack() {
	g.Turn.Attacks++
	g.Turn.DeckAttacks++
}

// PutIntoDeck puts a card that is outside the game or in a slot into a
// deck, below the top n cards (0: on top).
func (g *Game) PutIntoDeck(d DeckKind, id ObjectID, n int) {
	if g.Object(id).Zone.Kind == ZoneInPlay {
		g.removeFromSlot(id) // a monster or shop card in its slot
	}
	nid := g.move(id, DeckZone(d), NoPlayer)
	deck := g.Decks[d]
	at := max(len(deck)-n, 0)
	g.Decks[d] = append(deck[:at:at], append([]ObjectID{nid}, deck[at:]...)...)
	g.emit(Event{Kind: EvMovedToDeck, Player: NoPlayer, Object: nid, Card: g.Object(nid).Card, Amount: n, Text: d.String() + " deck"})
}

// PlaceFromDeck takes a card out of the monster deck and puts it on
// top of monster slot i ("put it in a monster slot").
func (g *Game) PlaceFromDeck(id ObjectID, i int) {
	g.Decks[MonsterDeck] = remove(g.Decks[MonsterDeck], id)
	if top, had := g.Monsters[i].TopOf(); had {
		g.Object(top).Zone = Zone{Kind: ZoneCovered, Slot: MonsterSlot, Index: i}
	}
	nid := g.putInSlot(id, MonsterSlot, i)
	g.enterMonsterSlot(nid)
}

// checkBonusSouls gives each active bonus soul to the first player, in
// turn order from the active player, who meets its condition (R-ZONE-14).
func (g *Game) checkBonusSouls() bool {
	gained := false
	for _, id := range append([]ObjectID(nil), g.BonusSouls...) {
		cond := g.def(id).BonusSoul
		if cond == nil || g.Object(id).Zone.Kind != ZoneOutside {
			continue
		}
		for _, p := range g.turnOrderFrom(g.Turn.Active) {
			if cond(g, p) {
				g.BonusSouls = remove(g.BonusSouls, id)
				nid := g.move(id, Zone{Kind: ZoneInPlay}, p)
				g.Object(nid).Role = RoleSoul
				g.Players[p].InPlay = append(g.Players[p].InPlay, nid)
				g.emit(Event{Kind: EvGainedSoul, Player: p, Object: nid, Card: g.Object(nid).Card, Text: "bonus soul"})
				gained = true
				break
			}
		}
	}
	return gained
}

// SoulFromDiscard takes the latest card from a discard pile and gives it
// to p as a soul ("This becomes a soul. Gain it.").
func (g *Game) SoulFromDiscard(d DeckKind, card CardRef, p PlayerID) {
	pile := g.Discards[d]
	for i := len(pile) - 1; i >= 0; i-- {
		if g.Object(pile[i]).Card != card {
			continue
		}
		g.Discards[d] = remove(pile, pile[i])
		nid := g.move(pile[i], Zone{Kind: ZoneInPlay}, p)
		g.Object(nid).Role = RoleSoul
		g.Players[p].InPlay = append(g.Players[p].InPlay, nid)
		g.emit(Event{Kind: EvGainedSoul, Player: p, Object: nid, Card: card})
		return
	}
}
