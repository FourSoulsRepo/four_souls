package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Mega Satan! (Epic Boss Card)
//
//	Players can be attacked and have 4+{DC}.
//	At the end of each turn, if the active player is the only living player, they win.
//	-Indomitable- Indomitable monsters can't be covered. If they would be put in a monster slot, expand monster slots by 1 and put them in that new slot.
var megaSatan = engine.CardDef{
	Ref:         "mega_satan",
	Kind:        engine.MonsterCard,
	Copies:      1,
	HP:          6,
	DC:          6,
	ATK:         6,
	Soul:        2,
	ThreePlus:   true,
	Indomitable: true,
	Rewards:     []engine.Reward{{Kind: engine.RewardTreasure, Amount: 2}},
	Statics:     []engine.Static{playersCanBeAttacked(4)},
	Abilities: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "At the end of each turn, if the active player is the only living player, they win.",
		Trigger: engine.Trigger{On: engine.EvEndOfTurn, Match: func(g *engine.Game, _ engine.ObjectID, _ engine.Event) bool {
			return onlyLiving(g, g.Turn.Active)
		}},
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			if onlyLiving(c.G, c.G.Turn.Active) {
				c.G.WinGame(c.G.Turn.Active)
			}
		})},
	}},
}

// onlyLiving: p is alive and every other player is dead.
func onlyLiving(g *engine.Game, p engine.PlayerID) bool {
	for _, pl := range g.Players {
		if pl.Dead != (pl.ID != p) {
			return false
		}
	}
	return true
}
