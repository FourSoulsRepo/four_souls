package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Cursed Fatty (Cursed Monster Card)
//
//	{Curse Effect}Each time a player rolls a ➄, they discard a loot card.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var cursedFatty = engine.CardDef{
	Ref:     "cursed_fatty",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      2,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
}
