package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Imp (Basic Monster Card)
//
//	Each time the attacking player rolls an attack roll of 3, this gains +1{DC} till end of turn.
var imp = engine.CardDef{
	Ref:     "imp",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      5,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Roll: true}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time the attacking player rolls an attack roll of 3, this gains +1 DC till end of turn.",
		Trigger: attackRollOnThis(3),
		Effects: []engine.Effect{thisGains(engine.StatMonsterDC, 1)},
	}},
}
