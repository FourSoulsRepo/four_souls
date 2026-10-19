package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Soul Of Lust (Bonus Soul Card)
//
//	Each time a player kills a monster, put a counter on this.
//	When the 6th counter is put on this, the active player gains this soul.
var soulOfLust = engine.CardDef{
	Ref:    "soul_of_lust",
	Kind:   engine.BonusSoulCard,
	Copies: 1,
	Soul:   1,
	Abilities: []engine.Ability{
		{
			Kind: engine.Triggered,
			Text: "Each time a player kills a monster, put a counter on this.",
			Trigger: engine.Trigger{On: engine.EvDied, Match: func(g *engine.Game, _ engine.ObjectID, e engine.Event) bool {
				return e.Player == engine.NoPlayer && g.Object(e.Object).Killer() != engine.NoPlayer
			}},
			Effects: []engine.Effect{counterOnBonusSoul},
		},
		sixthCounterGivesSoul,
	},
}
