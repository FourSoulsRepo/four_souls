package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Peep (Boss Card)
//
//	When this dies, search the monster deck for a card named The Bloat and put it in a monster slot not being attacked, then shuffle the monster deck.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var peep = engine.CardDef{
	Ref:     "peep",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
}
