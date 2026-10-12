package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Begotten (Basic Monster Card)
//
//	While this is at 1{HP}, it has +1{ATK}.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var begotten = engine.CardDef{
	Ref:     "begotten",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
}
