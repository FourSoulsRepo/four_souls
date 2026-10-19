package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Static Shock (Room Card)
//
//	Each time a player activates an item, they take 1 damage.
var staticShock = engine.CardDef{
	Ref:    "static_shock",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "Each time a player activates an item, they take 1 damage.",
		Trigger: engine.Trigger{On: engine.EvActivated, Match: func(g *engine.Game, _ engine.ObjectID, e engine.Event) bool {
			return e.Player != engine.NoPlayer && e.Object != 0 && g.Object(e.Object).Role == engine.RoleItem
		}},
		Effects: []engine.Effect{theyTake(1)},
	}},
}
