package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Berserk (Eternal Treasure Card)
//
//	You may attack an additional time on your turn.
//	Each time you kill a monster, gain +1{HP} till end of turn.
//	-Eternal- This can't be destroyed or put into discard.
var berserk = engine.CardDef{
	Ref:     "berserk",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Statics: []engine.Static{engine.YouHave(engine.StatAttacks, 1)},
	Abilities: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "Each time you kill a monster, gain +1 HP till end of turn.",
		Trigger: engine.Trigger{On: engine.EvDied, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
			return youKilledMonster(g, g.Object(self).Controller, e)
		}},
		Effects: []engine.Effect{engine.GainHPThisTurn(1, engine.You)},
	}},
}
