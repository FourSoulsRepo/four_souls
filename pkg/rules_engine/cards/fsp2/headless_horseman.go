package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Headless Horseman (Boss Card)
//
//	The first time this would die each turn, prevent death. This heals 2{HP} and gains +1{DC} and -1{ATK} till end of turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var headlessHorseman = engine.CardDef{
	Ref:     "headless_horseman",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      3,
	ATK:     2,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
}
