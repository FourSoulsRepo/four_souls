package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Battle Royale!!! (Room Card)
//
//	The active player may attack other players. Attacked players have 3+{DC}.
//	When a player dies, if that player was attacked this turn, that player gives the active player the item they would destroy for the death penalty.
var battleRoyale = engine.CardDef{
	Ref:       "battle_royale",
	Kind:      engine.RoomCard,
	Copies:    1,
	Statics:   []engine.Static{playersCanBeAttacked(3)},
	Abilities: []engine.Ability{attackedPlayersGivePenaltyItem},
}

// attackedPlayersGivePenaltyItem is Battle Royale's "When a player dies,
// if that player was attacked this turn, that player gives the active
// player the item they would destroy for the death penalty": it is set
// as the attack on them is declared, for the rest of the turn.
var attackedPlayersGivePenaltyItem = engine.Ability{
	Kind: engine.Triggered,
	Text: "When a player dies, if that player was attacked this turn, that player gives the active player the item they would destroy for the death penalty.",
	Trigger: engine.Trigger{On: engine.EvAttackTarget, Match: func(g *engine.Game, _ engine.ObjectID, _ engine.Event) bool {
		_, ok := g.AttackedPlayer()
		return ok
	}},
	Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
		if p, ok := c.G.AttackedPlayer(); ok {
			c.G.GivePenaltyItem(p, c.G.Turn.Active)
		}
	})},
}
