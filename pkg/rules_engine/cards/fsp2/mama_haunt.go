package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mama Haunt (Passive Treasure Card)
//
//	{Curse Effect}Your character doesn't recharge during your recharge step.
//	When you die, before paying penalties, give this to another player.
var mamaHaunt = engine.CardDef{
	Ref:     "mama_haunt",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Statics: []engine.Static{engine.YouHave(engine.StatNoCharRecharge, 1)},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "When you die, before paying penalties, give this to another player.",
		Trigger: engine.WhenYouDie(),
		Effects: []engine.Effect{giveThisAway},
	}},
}
