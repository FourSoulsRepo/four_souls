package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Equality (Room Card)
//
//	If the player would gain any amount of ¢, instead each player gains that much ¢.
var equality = engine.CardDef{
	Ref:    "equality",
	Kind:   engine.RoomCard,
	Copies: 1,
	Replacements: []engine.Replacement{{
		Text: "If the player would gain any amount of ¢, instead each player gains that much ¢.",
		When: func(_ *engine.Game, _ engine.ObjectID, a engine.Action) bool {
			return a.Kind == engine.ActGainCents && a.Amount > 0 && a.Player != engine.NoPlayer
		},
		Do: func(g *engine.Game, _ engine.ObjectID, a engine.Action) []engine.Action {
			var out []engine.Action
			for _, p := range g.TurnOrderFrom(a.Player) {
				b := a
				b.Player = p
				out = append(out, b)
			}
			return out
		},
	}},
}
