package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Baby Haunt (Passive Treasure Card)
//
//	{Curse Effect}Monsters have +1{DC} on your turn.
//	When you die, before paying penalties, give this to another player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var babyHaunt = engine.CardDef{
	Ref:    "baby_haunt",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
