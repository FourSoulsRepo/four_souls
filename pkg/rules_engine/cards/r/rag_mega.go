package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Rag Mega (Boss Card)
//
//	When an attack is declared on this, each non-active player discards a loot card.
//	When this dies, the active player kills up to 3 other players.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var ragMega = engine.CardDef{
	Ref:     "rag_mega",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      5,
	DC:      5,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
}
