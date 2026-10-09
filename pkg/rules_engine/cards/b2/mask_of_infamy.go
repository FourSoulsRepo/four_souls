package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mask Of Infamy (Boss Card)
//
//	While this is at 1{HP}, it has +2{DC}.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var maskOfInfamy = engine.CardDef{
	Ref:     "mask_of_infamy",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      4,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
}
