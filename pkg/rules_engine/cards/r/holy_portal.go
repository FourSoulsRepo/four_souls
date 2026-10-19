package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Holy Portal (Holy/Charmed Monster Card)
//
//	Each time a monster dies, the active player may attack an additional time this turn.
var holyPortal = engine.CardDef{
	Ref:     "holy_portal",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 5}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time a monster dies, the active player may attack an additional time this turn.",
		Trigger: engine.WhenAMonsterDies(),
		Effects: []engine.Effect{engine.ActivePlayerAttacksAgain()},
	}},
}
