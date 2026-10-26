package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Keeper’s Sack (Paid Treasure Card)
//
//	When you gain or purchase this, gain 14¢.
//	{Paid Effect}Pay 4¢: Give this to another player.
//	At the end of your turn, discard a loot card. If you can't, take 2 damage.
//	-Eternal- This can't be destroyed or put into discard.
var keepersSack = engine.CardDef{
	Ref:     "keepers_sack",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When you gain or purchase this, gain 14¢.",
			Trigger: cardkit.WhenYouGainThis(),
			Effects: []engine.Effect{engine.GainCents(14)},
		},
		{
			Kind:    engine.Activated,
			Text:    "$ Pay 4¢: Give this to another player.",
			Costs:   []engine.Cost{engine.PayCents(4)},
			Effects: []engine.Effect{cardkit.GiveThisAway},
		},
		{
			Kind:    engine.Triggered,
			Text:    "At the end of your turn, discard a loot card. If you can't, take 2 damage.",
			Trigger: engine.AtEndOfYourTurn(),
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				hand := c.G.Players[c.Controller].Hand
				if a[0] < 0 || a[0] >= len(hand) {
					c.Do(engine.DealDamage(2, engine.You))
					return
				}
				c.G.DiscardFromHand(c.Controller, hand[a[0]])
			}, engine.HandQuestion(cardkit.AskDiscardWhich))},
		},
	},
}
