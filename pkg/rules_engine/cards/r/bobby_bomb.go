package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Bobby Bomb (Passive Treasure Card)
//
//	When you gain or purchase this, reveal the top 10 cards of the loot deck. Put each card with "Bomb" in its name in your hand and the rest on the bottom of the loot deck.
//	At the start of your turn, put a card with "Bomb" in its name from the loot discard into your hand.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var bobbyBomb = engine.CardDef{
	Ref:    "bobby_bomb",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
