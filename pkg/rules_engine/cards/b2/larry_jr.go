package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Larry Jr. (Boss Card)
//
//	While this is at 2{HP} or less, it has +1{DC}.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var larryJr = engine.CardDef{
	Ref:     "larry_jr",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      3,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 6}},
}
