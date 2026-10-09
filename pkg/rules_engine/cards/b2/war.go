package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// War (Boss Card)
//
//	Each time this takes damage, it gains +1{ATK} till end of turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var war = engine.CardDef{
	Ref:     "war",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      3,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 8}},
}
