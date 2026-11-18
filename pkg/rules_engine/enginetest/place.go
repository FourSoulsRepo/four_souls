package enginetest

import (
	"fmt"
	"slices"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// PlaceFromDeck moves the card of object id from its deck to seat p: a
// loot card into p's hand, a treasure card or a trinket into p's play
// area as an item, as a draw and a gain would (ADR 033, Decision 20.5).
// It makes a new object for the card (R-ZONE-01) and returns it.
//
// It is quiet: no event and no trigger, so a card that does something as
// it enters play does not do it. It is for tools that measure a card on
// a copy of a game (cmd/bot_values); never for the server or a bot,
// which a depguard rule checks. The state is checked by CheckInvariants
// after the move.
func PlaceFromDeck(g *engine.Game, id engine.ObjectID, p engine.PlayerID, toPlay bool) (engine.ObjectID, error) {
	if p < 0 || int(p) >= len(g.Players) {
		return 0, fmt.Errorf("enginetest: no seat %d", p)
	}
	o := g.Object(id)
	if o.Zone.Kind != engine.ZoneDeck {
		return 0, fmt.Errorf("enginetest: %s is not in a deck", o.Card)
	}
	d := o.Zone.Deck
	at := slices.Index(g.Decks[d], id)
	if at < 0 {
		return 0, fmt.Errorf("enginetest: %s is not in the %s deck", o.Card, d)
	}
	def := g.Def(id)
	if toPlay && def.Kind != engine.TreasureCard && !def.Trinket {
		return 0, fmt.Errorf("enginetest: %s is no item", o.Card)
	}
	if !toPlay && def.Kind != engine.LootCard {
		return 0, fmt.Errorf("enginetest: %s is no loot card", o.Card)
	}
	g.Decks[d] = slices.Delete(slices.Clone(g.Decks[d]), at, at+1)
	card := o.Card
	o.Zone = engine.Zone{Kind: engine.ZoneGone}
	nid := engine.ObjectID(len(g.Objects))
	n := engine.Object{ID: nid, Card: card, Controller: p}
	pl := &g.Players[p]
	if toPlay {
		n.Zone = engine.Zone{Kind: engine.ZoneInPlay}
		n.Role, n.Charged = engine.RoleItem, !def.EntersDeactivated
		if def.EntersWithCounters > 0 {
			n.Counters = []engine.Counter{{Count: def.EntersWithCounters}}
		}
		pl.InPlay = append(slices.Clone(pl.InPlay), nid)
	} else {
		n.Zone = engine.Zone{Kind: engine.ZoneHand}
		pl.Hand = append(slices.Clone(pl.Hand), nid)
	}
	g.Objects = append(g.Objects, n)
	if err := g.CheckInvariants(); err != nil {
		return 0, fmt.Errorf("enginetest: after placing %s: %w", card, err)
	}
	return nid, nil
}
