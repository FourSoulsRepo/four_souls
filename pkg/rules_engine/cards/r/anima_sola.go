package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Anima Sola/The Revenant (Eternal Treasure Card)
//
//	{Tap Effect}Flip this item.
//	When this is flipped to this side, loot 1.
//	The Revenant: This item can be attacked.
//	When this reaches 0{HP}, the player attacking this gains its reward, then you flip it. That player may attack an additional time this turn.
//	-Eternal- This can't be destroyed or put into discard.
var animaSola = engine.CardDef{
	Ref:     "anima_sola",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
	Back:    "the_revenant",
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Flip this item.",
			Costs:   []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{flipThis},
		},
		{
			Kind:    engine.Triggered,
			Text:    "When this is flipped to this side, loot 1.",
			Trigger: engine.WhenFlippedToThisSide(),
			Effects: []engine.Effect{engine.Loot(1)},
		},
	},
}

// theRevenant is Anima Sola's back face. It can be attacked; it is
// eternal, so at 0 HP it does not die: its trigger flips it back.
var theRevenant = engine.CardDef{
	Ref:        "the_revenant",
	Kind:       engine.TreasureCard,
	Eternal:    true,
	Outside:    true,
	HP:         2,
	DC:         3,
	ATK:        2,
	Attackable: true,
	Rewards:    []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
	Abilities: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "When this reaches 0 HP, the player attacking this gains its reward, then you flip it. That player may attack an additional time this turn.",
		Trigger: engine.Trigger{On: engine.EvDamaged, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
			return e.Object == self && g.HP(self) == 0
		}},
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			attacker := c.G.Turn.Active
			them := *c
			them.Controller = attacker
			them.Do(engine.Loot(2))
			c.G.Flip(c.Source)
			c.G.Turn.Attacks++
		})},
	}},
}

// flipThis flips the ability's own object.
var flipThis = engine.EffectFunc(func(c *engine.Ctx) { c.G.Flip(c.Source) })
