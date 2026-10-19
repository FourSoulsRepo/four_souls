package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Ultra Greed! (Epic Boss Card)
//
//	After each attack roll the active player makes against this, each other player in turn order makes an attack roll against this.
//	The player who kills this gains its rewards.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var ultraGreed = engine.CardDef{
	Ref:     "ultra_greed",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      8,
	DC:      5,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 40}},
}
