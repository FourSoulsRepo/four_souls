package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// The card's own lines; their words are in en.json (ADR 021).
var (
	theMirrorPayHpDown = engine.NewKey("card.the_mirror.pay_hp_down_to_1")
)

// The Mirror (Room Card)
//
//	At the start of the turn, the active player may pay {HP} until they have 1{HP}. If they do, each time a monster dies this turn, they gain +1 treasure.
var theMirror = engine.CardDef{
	Ref:    "the_mirror",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "At the start of the turn, the active player may pay HP until they have 1 HP. If they do, each time a monster dies this turn, they gain +1 treasure.",
		Trigger: atStartOfTheTurn(),
		Effects: []engine.Effect{cardkit.May(theMirrorPayHpDown, engine.EffectFunc(func(c *engine.Ctx) {
			if n := c.G.PlayerHP(c.Controller) - 1; n > 0 {
				c.G.PayHP(c.Controller, n)
			}
			c.Do(engine.Linger(0)) // at 1 HP it costs nothing (S-YUG)
		}))},
	}},
	Lingering: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time a monster dies this turn, gain +1 treasure.",
		Trigger: engine.WhenAMonsterDies(),
		Effects: []engine.Effect{engine.GainTreasure(1)},
	}},
}
