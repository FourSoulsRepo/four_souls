package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// The card's own lines; their words are in en.json (ADR 021).
var (
	holyGreedlingGain3 = engine.NewKey("card.holy_greedling.gain_3")
)

// Holy Greedling (Holy/Charmed Monster Card)
//
//	Each time a monster dies, the active player may gain 3¢.
var holyGreedling = engine.CardDef{
	Ref:     "holy_greedling",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Roll: true}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time a monster dies, the active player may gain 3¢.",
		Trigger: engine.WhenAMonsterDies(),
		Effects: []engine.Effect{cardkit.May(holyGreedlingGain3, engine.GainCents(3))},
	}},
}
