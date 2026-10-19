package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Flip (Eternal Treasure Card)
//
//	Each time you die, after paying penalties, flip your character if able. Then recharge it. Discard your hand and loot 3.
//	-Eternal- This can't be destroyed or put into discard.
var flip = engine.CardDef{
	Ref:     "flip",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time you die, after paying penalties, flip your character if able. Then recharge it. Discard your hand and loot 3.",
		Trigger: engine.AfterYourDeathPenalty(),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			ch := c.G.Players[c.Controller].Character
			c.G.Flip(ch)
			c.G.Recharge(ch)
		}), discardHandLoot(3)},
	}},
}
