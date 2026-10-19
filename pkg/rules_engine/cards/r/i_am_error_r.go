package r

import engine "github.com/FourSoulsRepo/rules_engine"

// I Am Error (Room Card)
//
//	At the start of your turn roll-
//	1-2: Put a monster into discard and replace it with the top card of the monster deck
//	3-4: Reroll an item you control.
//	5: Discard your hand and loot 3.
//	6: Do all of the above.
var iAmErrorR = engine.CardDef{
	Ref:    "i_am_error_r",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "At the start of your turn roll- 1-2: Put a monster into discard and replace it with the top card of the monster deck. 3-4: Reroll an item you control. 5: Discard your hand and loot 3. 6: Do all of the above.",
		Trigger: atStartOfTheTurn(),
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
			Results(1, 2, replaceAMonster).
			Results(3, 4, rerollYourItem).
			Results(5, 5, discardHandLoot(3)).
			Results(6, 6, replaceAMonster, rerollYourItem, discardHandLoot(3)))},
	}},
}

// replaceAMonster: put a monster into discard; its slot is refilled.
var replaceAMonster = engine.Ask(func(c *engine.Ctx, a []int) {
	if ms := monstersInPlay(c.G); a[0] >= 0 && a[0] < len(ms) {
		c.G.DiscardMonster(ms[a[0]])
		c.G.RefillSlots()
	}
}, engine.Question{Text: "Put which monster into discard?", Options: func(c *engine.Ctx, _ []int) []string {
	return labelsOf(c.G, monstersInPlay(c.G))
}})

// rerollYourItem: reroll an item the controller controls.
var rerollYourItem = engine.Ask(func(c *engine.Ctx, a []int) {
	if ids := itemsOf(c.G, c.Controller); a[0] >= 0 && a[0] < len(ids) {
		c.G.RerollItem(ids[a[0]])
	}
}, engine.Question{Text: "Reroll which item?", Options: func(c *engine.Ctx, _ []int) []string {
	return labelsOf(c.G, itemsOf(c.G, c.Controller))
}})
