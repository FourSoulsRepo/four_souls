package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Corny (Basic Monster Card)
//
//	When this dies, deal 2 damage to a non-active player chosen at random.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var corny = engine.CardDef{
	Ref:       "corny",
	Kind:      engine.MonsterCard,
	Copies:    1,
	HP:        1,
	DC:        4,
	ATK:       1,
	Rewards:   []engine.Reward{{Kind: engine.RewardCents, Amount: 4}},
	ThreePlus: true,
}
