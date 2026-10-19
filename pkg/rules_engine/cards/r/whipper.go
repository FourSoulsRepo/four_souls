package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Whipper (Basic Monster Card)
//
//	Each time a player declares an attack, the active player rolls-
//	1-2: This deals 1 damage to that player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var whipper = engine.CardDef{
	Ref:     "whipper",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 5}},
}
