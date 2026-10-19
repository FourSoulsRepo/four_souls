package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Consolation Prize (Passive Treasure Card)
//
//	At the start of your turn, loot x, where x is the number of souls the player with the most souls controls minus the number of souls you control.
var consolationPrize = engine.CardDef{
	Ref:    "consolation_prize",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "At the start of your turn, loot x, where x is the number of souls the player with the most souls controls minus the number of souls you control.",
		Trigger: engine.AtStartOfYourTurn(),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			if x := soulsBehind(c.G, c.Controller); x > 0 {
				c.Do(engine.Loot(x))
			}
		})},
	}},
}
