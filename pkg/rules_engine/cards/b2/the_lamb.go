package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Lamb (Epic Boss Card)
//
//	When this dies, the active player may choose another player. They give you a soul they control.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
// Not generated: ATK is "6!".
var theLamb = engine.CardDef{
	Ref:     "the_lamb",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      6,
	DC:      3,
	Soul:    2,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 3}},
}
