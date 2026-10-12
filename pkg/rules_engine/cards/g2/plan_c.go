package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Plan C (One-Use Treasure Card)
//
//	{Tap Effect}Destroy this. If you do, kill a monster. At the end of the turn, kill all players.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var planC = engine.CardDef{
	Ref:    "plan_c",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
