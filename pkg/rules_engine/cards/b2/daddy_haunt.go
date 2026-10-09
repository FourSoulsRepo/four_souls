package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Daddy Haunt (Passive Treasure Card)
//
//	{Curse Effect}If you would take any amount of damage, take that much damage +1 instead.
//	When you die, before paying penalties, give this to another player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var daddyHaunt = engine.CardDef{
	Ref:    "daddy_haunt",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
