package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Gimpy (Eternal Treasure Card)
//
//	Each time you take damage, choose one- Gain 2¢. Loot 1, then discard a loot card. Gain +2{ATK} till end of turn.
//	-Eternal- This can't be destroyed or put into discard.
var gimpy = engine.CardDef{
	Ref:     "gimpy",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time you take damage, choose one- Gain 2¢. Loot 1, then discard a loot card. Gain +2 ATK till end of turn.",
		Trigger: engine.WhenYouTakeDamage(),
		Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
			switch a[0] {
			case 0:
				c.Do(engine.GainCents(2))
			case 1:
				c.Do(engine.Loot(1), discardOne)
			case 2:
				c.G.Boosts = append(c.G.Boosts, engine.Boost{Stat: engine.StatPlayerATK, Player: c.Controller, Amount: 2})
			}
		}, engine.Question{Text: "Choose one", Options: func(*engine.Ctx, []int) []string {
			return []string{"gain 2¢", "loot 1, then discard a loot card", "+2 ATK till end of turn"}
		}})},
	}},
}
