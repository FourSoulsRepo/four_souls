package r

import (
	"strconv"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// Jar Of Flies (Active Treasure Card)
//
//	Each time a monster dies, put a counter on this.
//	{Tap Effect}Remove 1 or more counters from this:
//	Loot 1. If 3+ counters were removed, deal 1 damage to a monster instead.
var jarOfFlies = engine.CardDef{
	Ref:    "jar_of_flies",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a monster dies, put a counter on this.",
			Trigger: engine.WhenAMonsterDies(),
			Effects: []engine.Effect{engine.AddCounters(1)},
		},
		{
			Kind:  engine.Activated,
			Text:  "↷ Remove 1 or more counters from this: Loot 1. If 3+ counters were removed, deal 1 damage to a monster instead.",
			Costs: []engine.Cost{engine.Tap(), engine.RemoveCounters(1)},
			// One counter is the cost; how many more come off is asked first.
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				more := max(a[0], 0)
				c.G.AddCountersTo(c.Source, -more)
				if 1+more >= 3 {
					c.Do(damageAMonster(1))
					return
				}
				c.Do(engine.Loot(1))
			}, engine.Question{Text: "Remove how many counters in all?", Options: func(c *engine.Ctx, _ []int) []string {
				n := c.G.Object(c.Source).CountersOf("")
				if n == 0 {
					return nil
				}
				out := make([]string, n+1)
				for i := range out {
					out[i] = strconv.Itoa(i + 1)
				}
				return out
			}})},
		},
	},
}
