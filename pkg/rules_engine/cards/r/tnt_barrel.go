package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Tnt Barrel (Good Event Card)
//
//	Choose a player or monster, then roll-
//	1-4: You and that player or monster each takes 2 damage.
//	5-6: That player or monster take 4 damage.
var tntBarrel = engine.CardDef{
	Ref:    "tnt_barrel",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Choose a player or monster, then roll- 1-4: You and that player or monster each takes 2 damage. 5-6: That player or monster take 4 damage.",
		Trigger: engine.WhenThisEntersPlay(),
		Effects: []engine.Effect{engine.AskWithRoll(func(c *engine.Ctx, a []int) {
			ts := monstersAndPlayers(c.G)
			if a[0] < 0 || a[0] >= len(ts) {
				return
			}
			them := *c
			them.Targets = []engine.Chosen{ts[a[0]]}
			them.Do(engine.Roll(tntTable))
		}, tntTable, engine.Question{Text: "Choose a player or monster.", Options: func(c *engine.Ctx, _ []int) []string {
			return chosenLabels(c.G, monstersAndPlayers(c.G))
		}})},
	}},
}

var tntTable = engine.RollTable{}.
	Results(1, 4, engine.DealDamage(2, engine.You), engine.DealDamage(2, 0)).
	Results(5, 6, engine.DealDamage(4, 0))
