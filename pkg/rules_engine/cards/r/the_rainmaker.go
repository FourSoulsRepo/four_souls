package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The Rainmaker (Boss Card)
//
//	Each time a player rolls a ❶, they loot 1.
//	When this dies, search the monster deck for a curse card and put it in a monster slot not being attacked, then shuffle the monster deck.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theRainmaker = engine.CardDef{
	Ref:     "the_rainmaker",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 3}},
}
