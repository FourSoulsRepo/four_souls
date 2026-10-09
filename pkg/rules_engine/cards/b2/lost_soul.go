package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Lost Soul (Lost Soul Card)
//
//	When this enters play, it becomes a soul.
//	(It's no longer an item.)
//	-Trinket- This loot becomes an item under your control when it resolves.
var lostSoul = engine.CardDef{
	Ref:     "lost_soul",
	Kind:    engine.LootCard,
	Copies:  1,
	Soul:    1,
	Trinket: true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When this enters play, it becomes a soul.",
			Trigger: engine.WhenThisEntersPlay(),
			Effects: []engine.Effect{engine.BecomeSoul()},
		},
	},
}
