package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Keeper’s Sack (Paid Treasure Card)
//
//	When you gain or purchase this, gain 14¢.
//	{Paid Effect}Pay 4¢: Give this to another player.
//	At the end of your turn, discard a loot card. If you can't, take 2 damage.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var keepersSack = engine.CardDef{
	Ref:     "keepers_sack",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
}
