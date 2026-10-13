package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Euthanasia (Passive Treasure Card)
//
//	Each time you roll the same result twice in a row on an attack roll on the same turn, kill the monster you're attacking.
var euthanasia = engine.CardDef{
	Ref:    "euthanasia",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "Each time you roll the same result twice in a row on an attack roll on the same turn, kill the monster you're attacking.",
		Trigger: engine.Trigger{On: engine.EvAttackRollRepeated, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
			return e.Player == g.Object(self).Controller // R-DICE-09, Q-12
		}},
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			if c.G.Attack.Started && c.G.Kind(c.G.Attack.Target) == engine.MonsterCard {
				c.G.KillObject(c.G.Attack.Target)
			}
		})},
	}},
}
