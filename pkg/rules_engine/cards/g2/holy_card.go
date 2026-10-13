package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Holy Card (Wildcard Card)
//
//	Choose a player. The next time that player would die this turn, prevent death. If it's their turn, cancel everything that hasn't resolved and end it.
var holyCard = engine.CardDef{
	Ref:    "holy_card",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.LootAbility,
		Text:    "Choose a player. The next time that player would die this turn, prevent death. If it's their turn, cancel everything that hasn't resolved and end it.",
		Targets: []engine.TargetSpec{engine.Choose(engine.TargetPlayer)},
		Effects: []engine.Effect{engine.WardDeath(0)},
	}},
}
