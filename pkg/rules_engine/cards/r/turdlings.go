package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Turdlings (Boss Card)
//
//	This only takes combat damage on attack rolls of 2.
//	When the attacking player makes their second attack roll this turn, after combat damage, cancel the attack.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var turdlings = engine.CardDef{
	Ref:     "turdlings",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      2,
	ATK:     2,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
}
