package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Stoney (Basic Monster Card)
//
//	Monsters have +1{DC}.
//	This can't be attacked.
//	When another monster dies, this dies.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var stoney = engine.CardDef{
	Ref:     "stoney",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      0,
	ATK:     0,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
}
