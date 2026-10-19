package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Blessing Of The Inner Eye (Room Card)
//
//	Each time a player declares an attack, before choosing what to attack, they may look at the top 3 cards of the monster deck and put them back in any order.
var blessingOfTheInnerEye = engine.CardDef{
	Ref:    "blessing_of_the_inner_eye",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time a player declares an attack, before choosing what to attack, they may look at the top 3 cards of the monster deck and put them back in any order.",
		Trigger: engine.Trigger{On: engine.EvAttackDeclared, Match: anyPlayer},
		Effects: []engine.Effect{cardkit.PutBackInOrder(engine.MonsterDeck, 3)},
	}},
}
