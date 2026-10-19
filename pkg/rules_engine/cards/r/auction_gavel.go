package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Auction Gavel (Active Treasure Card)
//
//	{Tap Effect}Choose a shop item. Other players may bid ¢ for it, starting at 1 or more ¢. In turn order, other players may top the high bid. Bidding ends when the high bid stands. The high bidder pays you the number of ¢ they bid and steals the chosen item.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var auctionGavel = engine.CardDef{
	Ref:       "auction_gavel",
	Kind:      engine.TreasureCard,
	Copies:    1,
	ThreePlus: true,
	Tap:       true,
}
