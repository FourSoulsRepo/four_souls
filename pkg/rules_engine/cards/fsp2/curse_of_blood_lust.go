package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Curse Of Blood Lust (Curse Card)
//
//	{Curse Effect}You must attack on your turn if able.
//	-Curse- When this enters play, give this to a player. When they die, they put this into discard.
var curseOfBloodLust = engine.CardDef{
	Ref:    "curse_of_blood_lust",
	Kind:   engine.EventCard,
	Copies: 1,
	Curse:  true,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "You must attack on your turn if able.",
		Trigger: engine.AtStartOfYourTurn(),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.Turn.MustAttacks = max(c.G.Turn.MustAttacks, 1) })},
	}},
}
