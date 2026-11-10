package rulesengine

// Final boss and minion slots: a challenge's own slot rows (ADR 028,
// Decision 7; R-CHAL-02 to R-CHAL-12). "Monsters in play" means the
// tops of all three monster rows (MonstersInPlay); "monster slot"
// effects (cover, refill, expand, put into a monster slot, deck attacks,
// Indomitable) use g.Monsters only.

// bossArmor is the most damage a final boss takes from one non-combat
// instance, and what an outright kill becomes (R-CHAL-05).
const bossArmor = 3

// MonsterRows lists the slot rows that hold monsters: the monster
// slots, then the final boss slots, then the minion slots (ADR 028).
func (g *Game) MonsterRows() [][]Slot {
	return [][]Slot{g.Monsters, g.Bosses, g.Minions}
}

// MonstersInPlay lists the monsters in play: the monster on top of each
// monster, final boss and minion slot, in that order (R-CHAL-03,
// R-CHAL-07). Events in a slot are not monsters.
func (g *Game) MonstersInPlay() []ObjectID {
	var out []ObjectID
	for _, rows := range g.MonsterRows() {
		for _, s := range rows {
			if top, ok := s.TopOf(); ok && g.Object(top).Role == RoleMonster {
				out = append(out, top)
			}
		}
	}
	return out
}

// IsFinalBoss reports whether id is a monster in play in a final boss
// slot (R-CHAL-03).
func (g *Game) IsFinalBoss(id ObjectID) bool {
	o := g.Object(id)
	return o.Zone.Kind == ZoneInPlay && o.Zone.Slot == FinalBossSlot && o.Role == RoleMonster
}

// IsMinion reports whether id is a minion: a monster in play in a
// minion slot, and nothing else (R-CHAL-08).
func (g *Game) IsMinion(id ObjectID) bool {
	o := g.Object(id)
	return o.Zone.Kind == ZoneInPlay && o.Zone.Slot == MinionSlot && o.Role == RoleMonster
}

// Held reports whether id can't leave its slot unless it is killed: a
// final boss or a minion (R-CHAL-04, R-CHAL-10, Q-37). Effects that
// would move it, put it into discard or on a deck do nothing to it.
func (g *Game) Held(id ObjectID) bool { return g.IsFinalBoss(id) || g.IsMinion(id) }

// armored caps damage to a final boss: more than 3 non-combat damage in
// one instance becomes 3; combat damage is never capped (R-CHAL-05).
func (g *Game) armored(id ObjectID, n int, combat bool) int {
	if combat || n <= bossArmor || !g.IsFinalBoss(id) {
		return n
	}
	return bossArmor
}

// PutInFinalBossSlot puts a new copy of card into a new final boss slot
// (R-CHAL-02), for a challenge's setup. It returns the new object, or 0
// if card is no monster of the game's sets.
func (g *Game) PutInFinalBossSlot(card CardRef) ObjectID {
	if d, ok := g.cards.find(card); !ok || d.Kind != MonsterCard {
		return 0
	}
	g.Bosses = append(g.Bosses, Slot{})
	id := g.newObject(card, Zone{Kind: ZoneOutside}, NoPlayer)
	return g.putInSlot(id, FinalBossSlot, len(g.Bosses)-1)
}

// PutInMinionSlot takes a card out of the monster deck and puts it into
// a new minion slot of its own (R-CHAL-09). An event there does its
// abilities, then goes to discard (R-CHAL-12). It returns the new
// object, or 0 if the card is not in the monster deck.
func (g *Game) PutInMinionSlot(id ObjectID) ObjectID {
	if !contains(g.Decks[MonsterDeck], id) {
		return 0
	}
	g.Decks[MonsterDeck] = remove(g.Decks[MonsterDeck], id)
	return g.toMinionSlot(id)
}

// MinionFromDeck puts the top card of the monster deck into a new
// minion slot (R-CHAL-09, R-KRAMPUS-08). It returns the new object, or
// 0 if the deck and its discard are empty.
func (g *Game) MinionFromDeck() ObjectID {
	id, ok := g.drawTop(MonsterDeck)
	if !ok {
		return 0
	}
	return g.toMinionSlot(id)
}

// toMinionSlot puts a card already taken from its zone list into a new
// minion slot; an event's abilities trigger (R-CHAL-12).
func (g *Game) toMinionSlot(id ObjectID) ObjectID {
	g.Minions = append(g.Minions, Slot{})
	nid := g.putInSlot(id, MinionSlot, len(g.Minions)-1)
	g.enterMonsterSlot(nid)
	return nid
}

// dropMinionSlot removes the empty minion slot i; the later minions'
// slot numbers shift down by 1 (R-CHAL-11, ADR 028 Decision 7.7).
func (g *Game) dropMinionSlot(i int) {
	g.Minions = append(g.Minions[:i:i], g.Minions[i+1:]...)
	for j := i; j < len(g.Minions); j++ {
		for _, id := range g.Minions[j].Cards {
			g.Object(id).Zone.Index = j
		}
	}
}

// eventSlot: z is a slot where an event resolves and then goes to
// discard: a monster slot or a minion slot (R-CARD-15, R-CHAL-12).
func eventSlot(z Zone) bool { return z.Slot == MonsterSlot || z.Slot == MinionSlot }
