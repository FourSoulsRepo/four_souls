package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// 9 Volt (Passive Treasure Card)
//
//	At the end of your turn, you may recharge up to 2 items you control.
var card9Volt = engine.CardDef{
	Ref:    "9_volt",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "At the end of your turn, you may recharge up to 2 items you control.",
		Trigger: engine.AtEndOfYourTurn(),
		Effects: []engine.Effect{rechargeYourItem, rechargeYourItem},
	}},
}
