package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Devil Deal (Good Event Card)
//
//	Choose one- Put this into discard. Loot 2. Take 1 damage. Take 2 damage. Search the treasure deck for a guppy item, gain it, then shuffle the treasure deck.
var devilDeal = engine.CardDef{
	Ref:    "devil_deal",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Choose one- Put this into discard. Loot 2, take 1 damage. Take 2 damage, search the treasure deck for a guppy item, gain it, then shuffle the treasure deck.",
			Trigger: engine.WhenThisEntersPlay(),
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				switch a[0] {
				case 1:
					c.Do(engine.Loot(2), engine.DealDamage(1, engine.You))
				case 2:
					c.Do(engine.DealDamage(2, engine.You), searchGuppy)
				}
			}, engine.Question{Text: "Choose one", Options: func(*engine.Ctx, []int) []string {
				return []string{"put this into discard", "loot 2, take 1 damage", "take 2 damage, search for a guppy item"}
			}})},
		},
	},
}
