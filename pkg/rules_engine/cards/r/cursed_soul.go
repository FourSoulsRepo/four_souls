package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Cursed Soul (Soul Treasure Card)
//
//	At the start of your turn, take 1 damage and put a counter on this. Then, if this has 6+ counters, it becomes a soul and loses all abilities.
//	At the end of your turn, give this to another player.
var cursedSoul = engine.CardDef{
	Ref:    "cursed_soul",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Soul:   1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "At the start of your turn, take 1 damage and put a counter on this. Then, if this has 6+ counters, it becomes a soul and loses all abilities.",
			Trigger: engine.AtStartOfYourTurn(),
			Effects: []engine.Effect{engine.DealDamage(1, engine.You), engine.AddCounters(1), soulAt(6)},
		},
		{
			Kind:    engine.Triggered,
			Text:    "At the end of your turn, give this to another player.",
			Trigger: engine.AtEndOfYourTurn(),
			Effects: []engine.Effect{cardkit.GiveThisAway},
		},
	},
}

// soulAt: "if this has n+ counters, it becomes a soul and loses all
// abilities".
func soulAt(n int) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) {
		if c.G.Object(c.Source).CountersOf("") >= n {
			c.Do(engine.BecomeSoul())
		}
	})
}
