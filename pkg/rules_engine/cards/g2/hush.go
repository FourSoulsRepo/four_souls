package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Hush (Epic Boss Card)
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var hush = engine.CardDef{
	Ref:     "hush",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      8,
	DC:      3,
	ATK:     2,
	Soul:    2,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 2}},
}
