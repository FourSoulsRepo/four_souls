package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Pride (Boss Card)
//
//	When an attack is declared on this, the active player chooses a player. That player discards 2 loot cards.
var pride = engine.CardDef{
	Ref:     "pride",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    1,
	HP:      2,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 5}},
	Abilities: []engine.Ability{
		{
			Kind: engine.Triggered,
			Text: "When an attack is declared on this, the active player chooses a player. That player discards 2 loot cards.",
			Trigger: engine.Trigger{On: engine.EvAttackTarget, Match: func(_ *engine.Game, self engine.ObjectID, e engine.Event) bool {
				return e.Object == self
			}},
			Effects: []engine.Effect{engine.Ask(forcedDiscard(2), forcedDiscardQuestions(2)...)},
		},
	},
}
