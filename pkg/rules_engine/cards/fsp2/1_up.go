package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// 1-Up (Passive Treasure Card)
//
//	When you would die on your turn, destroy this. If you do, prevent death, heal to full {HP}, and cancel your attack. You may attack an additional time this turn.
var card1Up = engine.CardDef{
	Ref:    "1_up",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "When you would die on your turn, destroy this. If you do, prevent death, heal to full HP, and cancel your attack. You may attack an additional time this turn.",
		Trigger: engine.Trigger{On: engine.EvDeathPending, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
			p := g.Object(self).Controller
			return e.Player == p && g.Turn.Active == p
		}},
		Effects: []engine.Effect{engine.DestroyThis(engine.PreventYourDeath(), engine.EffectFunc(func(c *engine.Ctx) {
			c.G.HealPlayer(c.Controller, c.G.Players[c.Controller].Damage)
		}), engine.CancelAttack(), engine.ActivePlayerAttacksAgain())},
	}},
}
