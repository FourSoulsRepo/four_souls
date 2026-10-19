package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The Harbingers/The Beast! (Outside Card)
//
//	The Harbingers: The active player must attack this once each turn if able.
//	If this would die, instead put a counter on this, the active player gains +1 treasure, cancels their attack and this heals to full {HP}.
//	If this has 4+ counters, flip it.
//	The Beast!: When this is flipped to this side or attacked, each player rolls-
//	Each player who rolls the lowest or tied for the lowest dies.
//	-Indomitable- Indomitable monsters can't be covered. If they would be put in a monster slot, expand monster slots by 1 and put them in that new slot.
var theHarbingers = engine.CardDef{
	Ref:         "the_harbingers",
	Kind:        engine.MonsterCard,
	Copies:      1,
	HP:          4,
	DC:          4,
	ATK:         1,
	Outside:     true,
	Indomitable: true,
	Back:        "the_beast",
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "The active player must attack this once each turn if able.",
			Trigger: engine.Trigger{On: engine.EvStartOfTurn, Match: func(*engine.Game, engine.ObjectID, engine.Event) bool { return true }},
			Effects: []engine.Effect{mustAttackThis},
		},
		{
			Kind:    engine.Triggered,
			Text:    "The active player must attack this once each turn if able.",
			Trigger: engine.WhenThisEntersPlay(),
			Effects: []engine.Effect{mustAttackThis},
		},
		{
			Kind: engine.Triggered,
			Text: "If this would die, instead put a counter on this, the active player gains +1 treasure, cancels their attack and this heals to full HP. If this has 4+ counters, flip it.",
			Trigger: engine.Trigger{On: engine.EvDeathPending, Match: func(_ *engine.Game, self engine.ObjectID, e engine.Event) bool {
				return e.Object == self
			}},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				g := c.G
				g.CancelStackItem(c.EventStack)
				g.AddCountersTo(c.Source, 1)
				c.Do(engine.GainTreasure(1))
				c.Do(engine.CancelAttack())
				g.HealObject(c.Source, g.Object(c.Source).Damage)
				if g.Object(c.Source).CountersOf("") >= 4 {
					g.Flip(c.Source)
				}
			})},
		},
	},
}

// theBeast is The Harbingers' back face.
var theBeast = engine.CardDef{
	Ref:         "the_beast",
	Kind:        engine.MonsterCard,
	HP:          6,
	DC:          4,
	ATK:         2,
	Soul:        2,
	Outside:     true,
	Indomitable: true,
	Rewards:     []engine.Reward{{Kind: engine.RewardLoot, Amount: 6}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When this is flipped to this side, each player rolls- Each player who rolls the lowest or tied for the lowest dies.",
			Trigger: engine.WhenFlippedToThisSide(),
			Effects: []engine.Effect{lowestRollDies},
		},
		{
			Kind:    engine.Triggered,
			Text:    "When this is attacked, each player rolls- Each player who rolls the lowest or tied for the lowest dies.",
			Trigger: whenAttackedThis(),
			Effects: []engine.Effect{lowestRollDies},
		},
	},
}
