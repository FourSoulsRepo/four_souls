package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Floor Spikes (Room Card)
//
//	Each time a player rolls a ❶, they take 1 damage.
var floorSpikes = engine.CardDef{
	Ref:    "floor_spikes",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "Each time a player rolls a 1, they take 1 damage.",
		Trigger: engine.Trigger{On: engine.EvRollResolved, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
			return e.Amount == 1 && anyPlayer(g, self, e)
		}},
		Effects: []engine.Effect{theyTake(1)},
	}},
}
