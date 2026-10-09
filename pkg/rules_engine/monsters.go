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
