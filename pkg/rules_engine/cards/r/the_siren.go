package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The Siren (Boss Card)
//
//	The active player may attack other players. Attacked players have 4+{DC}.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theSiren = engine.CardDef{
	Ref:     "the_siren",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      3,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
}
