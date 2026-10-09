package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Box! (One-Use Treasure Card)
//
//	{Tap Effect}Destroy this. If you do, you may play any number of additional loot cards till end of turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var box = engine.CardDef{
	Ref:    "box",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
