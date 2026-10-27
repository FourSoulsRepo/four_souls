package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// The card's own lines; their words are in en.json (ADR 021).
var (
	pillsRPutNext = engine.NewKey("card.pills_r.put_next")
)

// Pills! (Pill/Rune Card)
//
//	Roll-
//	1-2: Loot 3. Then put 3 cards from your hand on top of the loot deck in any order.
//	3-4: Loot 3. Then put 1 card from your hand on top of the loot deck.
//	5-6: Loot 3.
var pillsR = engine.CardDef{
	Ref:    "pills_r",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind: engine.LootAbility,
		Text: "Roll- 1-2: Loot 3. Then put 3 cards from your hand on top of the loot deck in any order. 3-4: Loot 3. Then put 1 card from your hand on top of the loot deck. 5-6: Loot 3.",
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
			Results(1, 2, engine.Loot(3), handToLootTop(3)).
			Results(3, 4, engine.Loot(3), handToLootTop(1)).
			Results(5, 6, engine.Loot(3)))},
	}},
}

// handToLootTop: the controller puts n cards from their hand on top of
// the loot deck, one at a time; the last one picked ends on top.
func handToLootTop(n int) engine.Effect {
	qs := make([]engine.Question, n)
	for i := range qs {
		text := cardkit.AskPutWhichCardOnTop
		if n > 1 {
			text = pillsRPutNext
		}
		qs[i] = engine.Question{Text: text, Options: func(c *engine.Ctx, a []int) []engine.Msg {
			return labelsOf(c.G, handLeft(c, a))
		}}
	}
	return engine.Ask(func(c *engine.Ctx, a []int) {
		left := append([]engine.ObjectID(nil), c.G.Players[c.Controller].Hand...)
		for _, i := range a {
			if i < 0 || i >= len(left) {
				return
			}
			c.G.HandToDeckTop(c.Controller, left[i])
			left = append(left[:i:i], left[i+1:]...)
		}
	}, qs...)
}

// handLeft is the controller's hand without the cards picked so far.
func handLeft(c *engine.Ctx, a []int) []engine.ObjectID {
	left := append([]engine.ObjectID(nil), c.G.Players[c.Controller].Hand...)
	for _, i := range a {
		if i < 0 || i >= len(left) {
			return nil
		}
		left = append(left[:i:i], left[i+1:]...)
	}
	return left
}
