package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Succubus (Paid Treasure Card)
//
//	{Paid Effect}Pay 3¢:
//	You may play up to 2 additional loot cards this turn.
var succubus = engine.CardDef{
	Ref:    "succubus",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Activated,
		Text:    "Pay 3¢: You may play up to 2 additional loot cards this turn.",
		Costs:   []engine.Cost{engine.PayCents(3)},
		Effects: []engine.Effect{engine.AddLootPlays(2)},
	}},
}
