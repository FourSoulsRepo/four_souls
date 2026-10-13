package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Void (Eternal Treasure Card)
//
//	{Tap Effect}Choose one:
//	•Discard your hand, then loot equal to the number of cards discarded.
//	•Put a monster not being attacked or a shop item into discard.
//	-Eternal- This can't be destroyed or put into discard.
var void = engine.CardDef{
	Ref:     "void",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
	Abilities: []engine.Ability{
		{
			Kind:  engine.Activated,
			Text:  "↷: Choose one: Discard your hand, then loot equal to the number of cards discarded. Put a monster not being attacked or a shop item into discard.",
			Costs: []engine.Cost{engine.Tap()},
			Modes: []engine.Mode{
				{Text: "Discard your hand, then loot equal to the number of cards discarded.", Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
					hand := append([]engine.ObjectID(nil), c.G.Players[c.Controller].Hand...)
					for _, id := range hand {
						c.G.DiscardFromHand(c.Controller, id)
					}
					c.Do(engine.Loot(len(hand)))
				})}},
				// The card's second option names two kinds of target; each is a mode here.
				{Text: "Put a monster not being attacked into discard.", Targets: []engine.TargetSpec{engine.ChooseWhere(engine.TargetMonster, notAttacked)}, Effects: []engine.Effect{discardTarget}},
				{Text: "Put a shop item into discard.", Targets: []engine.TargetSpec{engine.ChooseWhere(engine.TargetItem, inShop)}, Effects: []engine.Effect{discardTarget}},
			},
		},
	},
}
