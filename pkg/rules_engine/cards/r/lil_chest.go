package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Lil Chest (Passive Treasure Card)
//
//	Each time you kill a monster, roll-
//	1-2: Gain 2¢.
//	3-4: Loot 2, then put a loot card from your hand on top of the loot deck.
//	5-6: Recharge each item you control.
var lilChest = engine.CardDef{
	Ref:    "lil_chest",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "Each time you kill a monster, roll- 1-2: Gain 2¢. 3-4: Loot 2, then put a loot card from your hand on top of the loot deck. 5-6: Recharge each item you control.",
		Trigger: engine.Trigger{On: engine.EvDied, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
			return youKilledMonster(g, g.Object(self).Controller, e)
		}},
		Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
			Results(1, 2, engine.GainCents(2)).
			Results(3, 4, engine.Loot(2), handToLootTop(1)).
			Results(5, 6, engine.RechargeItemsOf(engine.You)))},
	}},
}
