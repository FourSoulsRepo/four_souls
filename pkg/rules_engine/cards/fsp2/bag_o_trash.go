package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Bag-O-Trash (Eternal Treasure Card)
//
//	{Paid Effect}Pay 3¢: Play an additional loot card this turn.
//	{Paid Effect}Pay 4¢: Loot 1.
//	{Paid Effect}Pay 6¢: Deal 1 damage to a monster or player.
//	-Eternal- This can't be destroyed or put into discard.
var bagOTrash = engine.CardDef{
	Ref:     "bag_o_trash",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Abilities: []engine.Ability{
		{Kind: engine.Activated, Text: "Pay 3¢: Play an additional loot card this turn.", Costs: []engine.Cost{engine.PayCents(3)}, Effects: []engine.Effect{engine.AddLootPlays(1)}},
		{Kind: engine.Activated, Text: "Pay 4¢: Loot 1.", Costs: []engine.Cost{engine.PayCents(4)}, Effects: []engine.Effect{engine.Loot(1)}},
		{
			Kind:    engine.Activated,
			Text:    "Pay 6¢: Deal 1 damage to a monster or player.",
			Costs:   []engine.Cost{engine.PayCents(6)},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetMonsterOrPlayer)},
			Effects: []engine.Effect{engine.DealDamage(1, 0)},
		},
	},
}
