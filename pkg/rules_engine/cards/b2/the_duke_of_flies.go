package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Duke Of Flies (Boss Card)
//
//	Each time this would take damage, the active player rolls-
//	1: Prevent that damage.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theDukeOfFlies = engine.CardDef{
	Ref:     "the_duke_of_flies",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      3,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
}
