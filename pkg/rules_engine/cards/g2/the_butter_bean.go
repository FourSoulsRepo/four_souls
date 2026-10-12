package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Butter Bean (Active Treasure Card)
//
//	{Tap Effect}Roll-
//	6: Recharge this.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theButterBean = engine.CardDef{
	Ref:    "the_butter_bean",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
