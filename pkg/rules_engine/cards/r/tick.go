package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Tick (Trinket Card)
//
//	The first time you kill a monster on your turn, heal +1{HP} and you may attack an additional time this turn.
//	-Eternal- This can't be destroyed or put into discard.
//	-Trinket- This loot becomes an item under your control when it resolves.
var tick = engine.CardDef{
	Ref:     "tick",
	Kind:    engine.LootCard,
	Copies:  1,
	Trinket: true,
	Eternal: true,
	Abilities: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "The first time you kill a monster on your turn, heal +1 HP and you may attack an additional time this turn.",
		Trigger: engine.Trigger{On: engine.EvDied, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
			me := g.Object(self).Controller
			return me == g.Turn.Active && youKilledMonster(g, me, e) && !cardkit.Used(g, self)
		}},
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			cardkit.MarkUsed(c.G, c.Source)
			c.G.HealPlayer(c.Controller, 1)
		}), engine.AddAttacks(1, engine.You)},
	}},
}

// youKilledMonster: the event is the death of a monster p killed.
func youKilledMonster(g *engine.Game, p engine.PlayerID, e engine.Event) bool {
	return e.Kind == engine.EvDied && e.Player == engine.NoPlayer && e.Prev != 0 && g.Object(e.Prev).KilledBy == int(p)+1
}
