package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// I Can’t Believe It’s Not Butter Bean (Active Treasure Card)
//
//	{Tap Effect}Cancel the ↷ or $ ability of an item or a loot card being played that targets an item or dice roll you control.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var iCantBelieveItsNotButterBean = engine.CardDef{
	Ref:    "i_cant_believe_its_not_butter_bean",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
