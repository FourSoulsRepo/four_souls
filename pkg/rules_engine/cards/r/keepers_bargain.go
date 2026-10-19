package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Keeper’s Bargain (Eternal Treasure Card)
//
//	{Tap Effect}Each other player may choose to gain 1¢. Gain 1¢ + 1¢ for each player who did.
//	Shop items you purchase cost 3¢ less.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var keepersBargain = engine.CardDef{
	Ref:     "keepers_bargain",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
}
