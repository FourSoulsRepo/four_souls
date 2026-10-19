package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Mama Gurdy (Boss Card)
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var mamaGurdy = engine.CardDef{
	Ref:     "mama_gurdy",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      7,
	DC:      3,
	ATK:     2,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Roll: true}},
}
