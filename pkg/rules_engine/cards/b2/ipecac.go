package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Ipecac (Passive Treasure Card)
//
//	{ATK}
//	Each time you roll an attack roll of 6, deal 1 damage to each other player.
var ipecac = engine.CardDef{
	Ref:     "ipecac",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Statics: []engine.Static{engine.YouHave(engine.StatPlayerATK, 1)},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time you roll an attack roll of 6, deal 1 damage to each other player.",
			Trigger: onYourAttackRollOf(6),
			Effects: []engine.Effect{engine.EachOtherPlayer(engine.DealDamage(1, engine.You))},
		},
	},
}
