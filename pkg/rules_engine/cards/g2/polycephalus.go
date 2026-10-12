package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Polycephalus (Boss Card)
//
//	Each time this takes damage, the next attack roll this turn has -2.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var polycephalus = engine.CardDef{
	Ref:     "polycephalus",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      3,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 3}},
}
