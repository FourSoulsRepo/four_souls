package r

import (
	"strings"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// Bomb Bum (Room Card)
//
//	{Paid Effect}Discard a loot card: Roll. if the discarded card had "bomb" in its name, roll 3x instead-
//	1-4: Deal 1 damage to a monster or player.
//	5-6: Take 1 damage.
var bombBum = engine.CardDef{
	Ref:    "bomb_bum",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Activated,
		Text:    "$ Discard a loot card: Roll. If the discarded card had \"bomb\" in its name, roll 3x instead- 1-4: Deal 1 damage to a monster or player. 5-6: Take 1 damage.",
		Costs:   []engine.Cost{engine.DiscardChosen(0)},
		Targets: []engine.TargetSpec{engine.Choose(engine.TargetYourHandCard)},
		Effects: []engine.Effect{rollTimesIfNamed("bomb", engine.RollTable{}.
			Results(1, 4, damageAnything(1)).
			Results(5, 6, youTake(1)))},
	}},
}

// rollTimesIfNamed rolls once, or 3 times if the discarded card (the
// first target) has word in its name.
func rollTimesIfNamed(word string, t engine.RollTable) engine.Effect {
	return engine.RollTimes(func(c *engine.Ctx) int {
		if strings.Contains(strings.ToLower(string(c.G.Object(c.Targets[0].Object).Card)), word) {
			return 3
		}
		return 1
	}, t)
}
