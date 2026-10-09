package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Pandora’s Box (Soul Treasure Card)
//
//	{Tap Effect}Destroy this. If you do, roll-
//	1: Gain 1¢.
//	2: Gain 6¢.
//	3: Kill a monster.
//	4: Loot 3.
//	5: Gain 9¢.
//	6: This becomes a soul. Gain it.
var pandorasBox = engine.CardDef{
	Ref:    "pandoras_box",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Soul:   1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:  engine.Activated,
			Text:  "↷: Destroy this. If you do, roll- 1: Gain 1¢. 2: Gain 6¢. 3: Kill a monster. 4: Loot 3. 5: Gain 9¢. 6: This becomes a soul. Gain it.",
			Costs: []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{engine.DestroyThis(engine.Roll(engine.RollTable{}.
				Results(1, 1, engine.GainCents(1)).
				Results(2, 2, engine.GainCents(6)).
				Results(3, 3, killAMonster).
				Results(4, 4, engine.Loot(3)).
				Results(5, 5, engine.GainCents(9)).
				Results(6, 6, engine.EffectFunc(func(c *engine.Ctx) {
					c.G.SoulFromDiscard(engine.TreasureDeck, "pandoras_box", c.Controller)
				}))))},
		},
	},
}
