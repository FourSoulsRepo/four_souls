package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Curse Of Fatigue (Curse Card)
//
//	{Curse Effect}You can't recharge more than 1 item during your recharge step.
//	-Curse- When this enters play, give this to a player. When they die, they put this into discard.
var curseOfFatigue = engine.CardDef{
	Ref:     "curse_of_fatigue",
	Kind:    engine.EventCard,
	Copies:  1,
	Curse:   true,
	Statics: []engine.Static{engine.YouHave(engine.StatRechargeLimit, 1)},
}
