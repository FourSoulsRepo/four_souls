package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Magic Skin (Active Treasure Card)
//
//	{Tap Effect}Gain +1 treasure, put a counter on this, then give this to another player.
//	At the start of your turn, roll-
//	If the roll is less than the number of counters on this, destroy it and all other items you control.
var magicSkin = engine.CardDef{
	Ref:    "magic_skin",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Gain +1 treasure, put a counter on this, then give this to another player.",
			Costs:   []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{engine.GainTreasure(1), engine.AddCounters(1), cardkit.GiveThisAway},
		},
		{
			Kind:    engine.Triggered,
			Text:    "At the start of your turn, roll- If the roll is less than the number of counters on this, destroy it and all other items you control.",
			Trigger: engine.AtStartOfYourTurn(),
			Effects: []engine.Effect{engine.Roll(magicSkinTable())},
		},
	},
}

func magicSkinTable() engine.RollTable {
	var t engine.RollTable
	for r := 1; r <= 6; r++ {
		t = t.Results(r, r, engine.EffectFunc(func(c *engine.Ctx) {
			o := c.G.Object(c.Source)
			if o.Zone.Kind != engine.ZoneInPlay || r >= o.CountersOf("") {
				return
			}
			c.G.DestroyObject(o.Controller, c.Source)
			// The player of the trigger; the target stays fixed (S-YUG).
			for _, id := range itemsOf(c.G, c.Controller) {
				c.G.DestroyObject(c.Controller, id)
			}
		}))
	}
	return t
}
