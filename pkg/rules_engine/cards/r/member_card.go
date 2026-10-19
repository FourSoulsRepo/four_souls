package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Member Card (Active Treasure Card)
//
//	{Tap Effect}Look at the top 3 cards of the treasure deck. Put them back in any order, then put a shop item into discard.
//	Shop items you purchase cost 4¢ less.
var memberCard = engine.CardDef{
	Ref:    "member_card",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{{
		Kind:    engine.Activated,
		Text:    "↷: Look at the top 3 cards of the treasure deck. Put them back in any order, then put a shop item into discard.",
		Costs:   []engine.Cost{engine.Tap()},
		Effects: []engine.Effect{putBackInOrder(engine.TreasureDeck, 3), discardAShopItem},
	}},
	Statics: []engine.Static{engine.YouHave(engine.StatShopPrice, -4)},
}
