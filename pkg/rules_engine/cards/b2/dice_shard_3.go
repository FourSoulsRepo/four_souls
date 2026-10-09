package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dice Shard (Dice Shard/Soul Heart Card)
//
//	Choose a dice roll. Its controller rerolls it.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var diceShard3 = engine.CardDef{
	Ref:    "dice_shard_3",
	Kind:   engine.LootCard,
	Copies: 3,
}
