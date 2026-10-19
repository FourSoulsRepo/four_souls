package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Spindown Dice (Eternal Treasure Card)
//
//	{Tap Effect}Destroy an item. If it was controlled by a player, they may steal an item from the shop.
//	At the end of your turn, roll-
//	4-6: Recharge this.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var spindownDice = engine.CardDef{
	Ref:       "spindown_dice",
	Kind:      engine.TreasureCard,
	Copies:    1,
	Eternal:   true,
	ThreePlus: true,
	Outside:   true,
	Tap:       true,
}
