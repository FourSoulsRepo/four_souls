package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The Jar (Paid Treasure Card)
//
//	{Paid Effect}Pay 4¢:
//	Choose a player or monster. Prevent the next 1 damage they would take this turn.
var theJar = engine.CardDef{
	Ref:    "the_jar",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Activated,
		Text:    "$ Pay 4¢: Choose a player or monster. Prevent the next 1 damage they would take this turn.",
		Costs:   []engine.Cost{engine.PayCents(4)},
		Targets: []engine.TargetSpec{engine.Choose(engine.TargetMonsterOrPlayer)},
		Effects: []engine.Effect{engine.PreventDamage(1, 0)},
	}},
}
