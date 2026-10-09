package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Cursed Horf (Cursed Monster Card)
//
//	{Curse Effect}Each time a player rolls a ➁, they take 2 damage.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var cursedHorf = engine.CardDef{
	Ref:     "cursed_horf",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 3}},
}
