package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Fetal Haunt (Passive Treasure Card)
//
//	{Curse Effect}At the end of your turn, give another player a loot card.
//	When you die, before paying penalties, give this to another player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var fetalHaunt = engine.CardDef{
	Ref:    "fetal_haunt",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
