package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Bony (Basic Monster Card)
//
//	Each time the attacking player rolls an attack roll of 1, this gains +2{ATK} till end of turn.
var bony = engine.CardDef{
	Ref:     "bony",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time the attacking player rolls an attack roll of 1, this gains +2 ATK till end of turn.",
		Trigger: attackRollOnThis(1),
		Effects: []engine.Effect{thisGains(engine.StatMonsterATK, 2)},
	}},
}
