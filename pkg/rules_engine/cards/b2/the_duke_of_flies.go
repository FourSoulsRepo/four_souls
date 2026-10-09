package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Duke Of Flies (Boss Card)
//
//	Each time this would take damage, the active player rolls-
//	1: Prevent that damage.
var theDukeOfFlies = engine.CardDef{
	Ref:     "the_duke_of_flies",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    1,
	HP:      4,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
	Abilities: []engine.Ability{
		{
			Kind: engine.Triggered,
			Text: "Each time this would take damage, the active player rolls- 1: Prevent that damage.",
			Trigger: engine.Trigger{On: engine.EvDamagePending, Match: func(_ *engine.Game, self engine.ObjectID, e engine.Event) bool {
				return e.Object == self
			}},
			Effects: []engine.Effect{engine.Roll(engine.RollTable{}.Results(1, 1, engine.PreventNextDamage(engine.This)))},
		},
	},
}
