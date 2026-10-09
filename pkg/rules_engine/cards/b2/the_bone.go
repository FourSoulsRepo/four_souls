package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Bone (Eternal Treasure Card)
//
//	{Tap Effect}Put a counter on this.
//	{Paid Effect}Remove 1 counter from this: Add +1 to a dice roll.
//	{Paid Effect}Remove 2 counters from this: Deal 1 damage to a monster or player.
//	{Paid Effect}Remove 5 counters from this: This becomes a soul and loses all abilities.
//	-Eternal- This can't be destroyed or put into discard.
var theBone = engine.CardDef{
	Ref:     "the_bone",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Soul:    1,
	Eternal: true,
	Outside: true,
	Tap:     true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Put a counter on this.",
			Costs:   []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{engine.AddCounters(1)},
		},
		{
			Kind:    engine.Activated,
			Text:    "Remove 1 counter from this: Add +1 to a dice roll.",
			Costs:   []engine.Cost{engine.RemoveCounters(1)},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetDiceRoll)},
			Effects: []engine.Effect{engine.ModifyRoll(1, 0)},
		},
		{
			Kind:    engine.Activated,
			Text:    "Remove 2 counters from this: Deal 1 damage to a monster or player.",
			Costs:   []engine.Cost{engine.RemoveCounters(2)},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetMonsterOrPlayer)},
			Effects: []engine.Effect{engine.DealDamage(1, 0)},
		},
		{
			Kind:    engine.Activated,
			Text:    "Remove 5 counters from this: This becomes a soul and loses all abilities.",
			Costs:   []engine.Cost{engine.RemoveCounters(5)},
			Effects: []engine.Effect{engine.BecomeSoul()},
		},
	},
}
