package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Gaper (Basic Monster Card)
//
//	When this dies on an attack roll of 6, the active player must make an additional attack on the monster deck.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var gaper = engine.CardDef{
	Ref:     "gaper",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 3}},
}
