package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mega Fatty (Boss Card)
//
//	Each time this deals combat damage, it heals 1{HP}.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var megaFatty = engine.CardDef{
	Ref:     "mega_fatty",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      3,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
}
