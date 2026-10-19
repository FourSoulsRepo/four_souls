package r

import (
	"strconv"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// Red Key (Active Treasure Card)
//
//	{Tap Effect}Look at the top 3 cards of the room or monster deck. You may put one of those in a slot and the rest back. This can't be activated during an attack.
var redKey = engine.CardDef{
	Ref:    "red_key",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{{
		Kind: engine.Activated,
		Text: "↷: Look at the top 3 cards of the room or monster deck. You may put one of those in a slot and the rest back. This can't be activated during an attack.",
		Costs: []engine.Cost{engine.Tap(), engine.If("not during an attack", func(g *engine.Game, _ engine.PlayerID, _ engine.ObjectID) bool {
			return !g.Attack.On
		})},
		Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
			ds := monsterOrRoomDecks(c.G)
			if a[0] < 0 || a[0] >= len(ds) {
				return
			}
			d := ds[a[0]]
			top := c.G.DeckTop(d, 3)
			if a[1] < 0 || a[1] >= len(top) {
				return // "none": they all stay
			}
			slots := keySlots(c.G, d)
			if a[2] < 0 || a[2] >= len(slots) {
				return
			}
			if d == engine.RoomDeck {
				c.G.PlaceRoomFromDeck(top[a[1]], slots[a[2]])
			} else {
				c.G.PlaceFromDeck(top[a[1]], slots[a[2]])
			}
		}, engine.Question{Text: "Look at which deck?", Options: func(c *engine.Ctx, _ []int) []string {
			return deckLabels(monsterOrRoomDecks(c.G))
		}}, engine.Question{
			Text: "Put which card in a slot?",
			Options: func(c *engine.Ctx, a []int) []string {
				if ds := monsterOrRoomDecks(c.G); a[0] >= 0 && a[0] < len(ds) {
					return append(labelsOf(c.G, c.G.DeckTop(ds[a[0]], 3)), "none")
				}
				return nil
			},
		}, engine.Question{
			Text: "Put it in which slot?",
			Options: func(c *engine.Ctx, a []int) []string {
				ds := monsterOrRoomDecks(c.G)
				if a[0] < 0 || a[0] >= len(ds) || a[1] < 0 || a[1] >= len(c.G.DeckTop(ds[a[0]], 3)) {
					return nil
				}
				var out []string
				for _, i := range keySlots(c.G, ds[a[0]]) {
					out = append(out, "slot "+strconv.Itoa(i+1))
				}
				return out
			},
		})},
	}},
}

// keySlots are the slots a card from deck d may go to: monster slots not
// being attacked, or the room slots.
func keySlots(g *engine.Game, d engine.DeckKind) []int {
	if d == engine.MonsterDeck {
		return freeSlots(g)
	}
	out := make([]int, len(g.Rooms))
	for i := range out {
		out[i] = i
	}
	return out
}
