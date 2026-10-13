package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Cursed Tumor (Cursed Monster Card)
//
//	{Curse Effect}Each time a player rolls a ➃, each player takes 1 damage.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var cursedTumor = engine.CardDef{
	Ref:     "cursed_tumor",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Roll: true}},
}
