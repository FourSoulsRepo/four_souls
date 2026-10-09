package rulesengine

import "fmt"

// CheckInvariants reports the first broken invariant of the state, or nil.
// Fuzz tests call it after every step; it is cheap enough for debugging.
func (g *Game) CheckInvariants() error {
	where := make([]int, len(g.Objects)) // how many lists hold each object
	mark := func(list []ObjectID, kind ZoneKind, what string) error {
		for _, id := range list {
			if id < 0 || int(id) >= len(g.Objects) {
				return fmt.Errorf("%s holds unknown object %d", what, id)
			}
			where[id]++
			if k := g.Objects[id].Zone.Kind; k != kind {
				return fmt.Errorf("%s holds object %d (%s) whose zone is %d, want %d", what, id, g.Objects[id].Card, k, kind)
			}
		}
		return nil
	}
	for d := range deckCount {
		if err := mark(g.Decks[d], ZoneDeck, "deck"); err != nil {
			return err
		}
		if err := mark(g.Discards[d], ZoneDiscard, "discard"); err != nil {
			return err
		}
		// A card only ever sits in its own deck or discard (R-ZONE-04).
		for _, id := range append(append([]ObjectID(nil), g.Decks[d]...), g.Discards[d]...) {
			if own, ok := g.kindOf(id).Deck(); !ok || own != d {
				return fmt.Errorf("%s is in the %s deck or discard", g.Objects[id].Card, d)
			}
		}
	}
	for _, s := range g.Shop {
		if top, ok := s.TopOf(); ok && g.kindOf(top) != TreasureCard {
			return fmt.Errorf("%s is in a shop slot", g.Objects[top].Card)
		}
	}
	for _, pl := range g.Players {
		if err := mark(pl.Hand, ZoneHand, fmt.Sprintf("player %d hand", pl.ID)); err != nil {
			return err
		}
		if err := mark(append([]ObjectID{pl.Character}, pl.InPlay...), ZoneInPlay, fmt.Sprintf("player %d", pl.ID)); err != nil {
			return err
		}
		if pl.Cents < 0 {
			return fmt.Errorf("player %d has %d cents", pl.ID, pl.Cents)
		}
		if pl.Damage < 0 || g.PlayerHP(pl.ID) < 0 {
			return fmt.Errorf("player %d has bad damage %d", pl.ID, pl.Damage)
		}
		for _, id := range pl.InPlay {
			if g.Objects[id].Controller != pl.ID {
				return fmt.Errorf("player %d holds object %d controlled by %d", pl.ID, id, g.Objects[id].Controller)
			}
		}
	}
	for _, rows := range [][]Slot{g.Shop, g.Monsters, g.Rooms} {
		for _, s := range rows {
			for i, id := range s.Cards {
				want := ZoneCovered
				if i == len(s.Cards)-1 {
					want = ZoneInPlay
				}
				if err := mark([]ObjectID{id}, want, "slot"); err != nil {
					return err
				}
			}
		}
	}
	if err := mark(g.BonusSouls, ZoneOutside, "bonus souls"); err != nil {
		return err
	}
	for _, it := range g.Stack {
		if it.Kind == StackLoot {
			if err := mark([]ObjectID{it.Source}, ZoneStack, "stack"); err != nil {
				return err
			}
		}
	}
	for id, n := range where {
		o := g.Objects[id]
		switch {
		case n > 1:
			return fmt.Errorf("object %d (%s) is in %d places", id, o.Card, n)
		case n == 0 && o.Zone.Kind != ZoneGone && o.Zone.Kind != ZoneOutside:
			return fmt.Errorf("object %d (%s) in zone %d is in no list", id, o.Card, o.Zone.Kind)
		}
		if o.Zone.Kind == ZoneInPlay && o.Role == RoleMonster && g.HP(o.ID) < 0 {
			return fmt.Errorf("monster %d has negative HP", id)
		}
	}
	if !g.Over && g.Waiting.Kind == PromptNone {
		return fmt.Errorf("the game waits for nothing")
	}
	return nil
}
