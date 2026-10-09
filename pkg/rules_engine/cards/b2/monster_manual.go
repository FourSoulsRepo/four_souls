package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Monster Manual (Active Treasure Card)
//
//	{Tap Effect}Choose a monster. The active player must attack that monster this turn if able.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var monsterManual = engine.CardDef{
	Ref:    "monster_manual",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
