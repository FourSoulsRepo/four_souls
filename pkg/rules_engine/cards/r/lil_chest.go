package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Lil Chest (Passive Treasure Card)
//
//	Each time you kill a monster, roll-
//	1-2: Gain 2¢.
//	3-4: Loot 2, then put a loot card from your hand on top of the loot deck.
//	5-6: Recharge each item you control.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var lilChest = engine.CardDef{
	Ref:    "lil_chest",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
