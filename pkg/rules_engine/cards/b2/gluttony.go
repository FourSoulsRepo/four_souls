package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Gluttony (Boss Card)
//
//	Each time this takes combat damage on an attack roll of 6, deal 1 damage to the player to the active player's left.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var gluttony = engine.CardDef{
	Ref:     "gluttony",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      3,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
}
