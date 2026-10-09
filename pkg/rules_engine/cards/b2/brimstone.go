package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Brimstone (Passive Treasure Card)
//
//	{ATK}
//	Each time you deal combat damage to a monster, deal 1 damage to another player.
var brimstone = engine.CardDef{
	Ref:     "brimstone",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Statics: []engine.Static{engine.YouHave(engine.StatPlayerATK, 1)},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time you deal combat damage to a monster, deal 1 damage to another player.",
			Trigger: whenYouDealCombatDamage(),
			Effects: []engine.Effect{damageAPlayer(1, true)},
		},
	},
}
