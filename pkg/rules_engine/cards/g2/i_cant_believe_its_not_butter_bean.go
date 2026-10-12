package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// I Can’t Believe It’s Not Butter Bean (Active Treasure Card)
//
//	{Tap Effect}Cancel the ↷ or $ ability of an item or a loot card being played that targets an item or dice roll you control.
var iCantBelieveItsNotButterBean = engine.CardDef{
	Ref:    "i_cant_believe_its_not_butter_bean",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{{
		Kind:    engine.Activated,
		Text:    "↷: Cancel the ↷ or $ ability of an item or a loot card being played that targets an item or dice roll you control.",
		Costs:   []engine.Cost{engine.Tap()},
		Targets: []engine.TargetSpec{engine.ChooseWhere(engine.TargetStackAbility, targetsYours)},
		Effects: []engine.Effect{engine.CancelTarget(0)},
	}},
}
