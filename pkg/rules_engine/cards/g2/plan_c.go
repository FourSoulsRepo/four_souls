package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Plan C (One-Use Treasure Card)
//
//	{Tap Effect}Destroy this. If you do, kill a monster. At the end of the turn, kill all players.
var planC = engine.CardDef{
	Ref:    "plan_c",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{{
		Kind:    engine.Activated,
		Text:    "↷: Destroy this. If you do, kill a monster. At the end of the turn, kill all players.",
		Costs:   []engine.Cost{engine.Tap()},
		Effects: []engine.Effect{engine.DestroyThis(killAMonster, engine.Linger(0))},
	}},
	Lingering: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "At the end of the turn, kill all players.",
		Trigger: engine.Trigger{On: engine.EvEndOfTurn, Match: func(*engine.Game, engine.ObjectID, engine.Event) bool { return true }},
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			p := c.G.Turn.Active // in turn order from the active player
			for range c.G.Players {
				them := *c
				them.Targets = []engine.Chosen{{Kind: engine.TargetPlayer, Player: p}}
				them.Do(engine.Kill(0))
				p = c.G.PlayerAfter(p)
			}
		})},
	}},
}
