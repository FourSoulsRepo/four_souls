package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Algiz (Pill/Rune Card)
//
//	Choose a player. Each item they control gains eternal till end of turn.
var algiz = engine.CardDef{
	Ref:    "algiz",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.LootAbility,
		Text:    "Choose a player. Each item they control gains eternal till end of turn.",
		Targets: []engine.TargetSpec{engine.Choose(engine.TargetPlayer)},
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			// Only the items they control now; later ones don't (S-DIS).
			for _, id := range itemsOf(c.G, c.Targets[0].Player) {
				c.G.Boosts = append(c.G.Boosts, engine.Boost{Source: c.Source, Stat: engine.StatEternal, Player: engine.NoPlayer, Object: id, Amount: 1})
			}
		})},
	}},
}
