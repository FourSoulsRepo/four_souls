package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Haunt (Boss Card)
//
//	Every other time this takes damage each turn, it gains +1{DC} till end of turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theHaunt = engine.CardDef{
	Ref:     "the_haunt",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
}
