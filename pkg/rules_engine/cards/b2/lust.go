package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Lust (Boss Card)
//
//	Each time this takes combat damage, it deals 1 damage to the attacking player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var lust = engine.CardDef{
	Ref:     "lust",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      4,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
}
