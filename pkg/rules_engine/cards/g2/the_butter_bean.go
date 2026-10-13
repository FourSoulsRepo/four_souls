package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Butter Bean (Active Treasure Card)
//
//	{Tap Effect}Roll-
//	6: Recharge this.
var theButterBean = engine.CardDef{
	Ref:    "the_butter_bean",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{{
		Kind:    engine.Activated,
		Text:    "↷: Roll- 6: Recharge this.",
		Costs:   []engine.Cost{engine.Tap()},
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.Results(6, 6, engine.RechargeSelf()))},
	}},
}
