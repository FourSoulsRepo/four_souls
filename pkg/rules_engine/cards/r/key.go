package r

import (
	"strconv"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// Key (Wildcard Card)
//
//	Choose one-
//	The active player may attack an additional time this turn.
//	Put the top card of the monster deck in a monster slot not being attacked.
var key = engine.CardDef{
	Ref:    "key",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind: engine.LootAbility,
		Text: "Choose one- The active player may attack an additional time this turn. Put the top card of the monster deck in a monster slot not being attacked.",
		Modes: []engine.Mode{
			{Text: "The active player may attack an additional time this turn.", Effects: []engine.Effect{engine.ActivePlayerAttacksAgain()}},
			{Text: "Put the top card of the monster deck in a monster slot not being attacked.", Effects: []engine.Effect{deckTopToSlot}},
		},
	}},
}

// deckTopToSlot puts the top monster card in a slot not being attacked,
// covering the card there.
var deckTopToSlot = engine.Ask(func(c *engine.Ctx, a []int) {
	if slots := freeSlots(c.G); a[0] >= 0 && a[0] < len(slots) {
		c.G.CoverMonsterSlot(slots[a[0]]) // an event resolves at once
	}
}, slotQuestion("Put it in which monster slot?"))

// slotQuestion asks for a monster slot not being attacked.
func slotQuestion(text string) engine.Question {
	return engine.Question{Text: text, Options: func(c *engine.Ctx, _ []int) []string {
		var out []string
		for _, i := range freeSlots(c.G) {
			label := "slot " + strconv.Itoa(i+1)
			if top, ok := c.G.Monsters[i].TopOf(); ok {
				label += ": " + string(c.G.Object(top).Card)
			}
			out = append(out, label)
		}
		return out
	}}
}
