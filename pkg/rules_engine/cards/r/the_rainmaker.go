package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// The Rainmaker (Boss Card)
//
//	Each time a player rolls a ❶, they loot 1.
//	When this dies, search the monster deck for a curse card and put it in a monster slot not being attacked, then shuffle the monster deck.
var theRainmaker = engine.CardDef{
	Ref:     "the_rainmaker",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 3}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player rolls a ❶, they loot 1.",
			Trigger: engine.OnRollOf(1),
			Effects: []engine.Effect{cardkit.ForRoller(engine.Loot(1))},
		},
		{
			Kind:    engine.Triggered,
			Text:    "When this dies, search the monster deck for a curse card and put it in a monster slot not being attacked, then shuffle the monster deck.",
			Trigger: engine.WhenThisDies(),
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				curses, slots := cursesInDeck(c.G), freeSlots(c.G)
				if a[0] >= 0 && a[0] < len(curses) && a[1] >= 0 && a[1] < len(slots) {
					c.G.PlaceFromDeck(curses[a[0]], slots[a[1]])
				}
				c.G.ShuffleDeck(engine.MonsterDeck)
			}, engine.Question{Text: "Which curse?", Options: func(c *engine.Ctx, _ []int) []string {
				return labelsOf(c.G, cursesInDeck(c.G))
			}}, slotQuestionAfterPick("Put it in which monster slot?"))},
		},
	},
}
