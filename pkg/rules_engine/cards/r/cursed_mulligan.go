package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Cursed Mulligan (Cursed Monster Card)
//
//	Each time a monster dies, expand monster slots by 1.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var cursedMulligan = engine.CardDef{
	Ref:     "cursed_mulligan",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
}
