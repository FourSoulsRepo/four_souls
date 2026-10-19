package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Anima Sola/The Revenant (Eternal Treasure Card)
//
//	{Tap Effect}Flip this item.
//	When this is flipped to this side, loot 1.
//	The Revenant: This item can be attacked.
//	When this reaches 0{HP}, the player attacking this gains its reward, then you flip it. That player may attack an additional time this turn.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var animaSola = engine.CardDef{
	Ref:     "anima_sola",
	Kind:    engine.TreasureCard,
	Copies:  1,
	HP:      2,
	DC:      3,
	ATK:     2,
	Eternal: true,
	Outside: true,
	Tap:     true,
}
