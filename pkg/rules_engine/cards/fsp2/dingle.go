package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dingle (Boss Card)
//
//	Each time the active player deals damage to this, they roll-
//	1-2: They take 1 damage.
//	3-4: Each player takes 1 damage.
//	5-6: This takes 1 damage.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var dingle = engine.CardDef{
	Ref:     "dingle",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      3,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 6}},
}
