package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Holy Bony (Holy/Charmed Monster Card)
//
//	Each time a player rolls a ❻, they may look at the top 3 cards of a deck and put them back in any order.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var holyBony = engine.CardDef{
	Ref:     "holy_bony",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
}
