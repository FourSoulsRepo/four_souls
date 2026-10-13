package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Trap Door! (Good Event Card)
//
//	Roll-
//	1-3: Look at the top 3 cards of the monster deck and put them back in any order. You may attack an additional time this turn.
//	4-6: Take 1 damage.
var trapDoor = engine.CardDef{
	Ref:    "trap_door",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Roll- 1-3: Look at the top 3 cards of the monster deck and put them back in any order. You may attack an additional time this turn. 4-6: Take 1 damage.",
		Trigger: engine.WhenThisEntersPlay(),
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
			Results(1, 3, putBackInOrder(engine.MonsterDeck, 3), engine.ActivePlayerAttacksAgain()).
			Results(4, 6, engine.DealDamage(1, engine.You)))},
	}},
}
