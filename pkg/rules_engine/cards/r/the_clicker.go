package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The Clicker (One-Use Treasure Card)
//
//	{Tap Effect}Destroy this. If you do, choose a character card from outside the game. Replace your character with it and your starting item with the chosen card's starting item.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theClicker = engine.CardDef{
	Ref:    "the_clicker",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
