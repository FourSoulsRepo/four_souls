package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// The Capricious (Character Card)
//
//	{Tap Effect}Choose one- Play an additional loot card this turn. Put a shop item into discard.
var theCapricious = engine.CardDef{
	Ref:          "the_capricious",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "glitch",
	Tap:          true,
	Abilities: []engine.Ability{{
		Kind:  engine.Activated,
		Text:  "↷: Choose one- Play an additional loot card this turn. Put a shop item into discard.",
		Costs: []engine.Cost{engine.Tap()},
		Modes: []engine.Mode{
			{Text: "Play an additional loot card this turn.", Effects: []engine.Effect{engine.AddLootPlays(1)}},
			{Text: "Put a shop item into discard.", Targets: []engine.TargetSpec{engine.ChooseWhere(engine.TargetItem, inShop)}, Effects: []engine.Effect{cardkit.DiscardTarget}},
		},
	}},
}
