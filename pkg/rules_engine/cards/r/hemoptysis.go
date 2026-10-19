package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Hemoptysis (Eternal Treasure Card)
//
//	{Tap Effect}Put a counter on this.
//	{Paid Effect}Remove a counter from this:
//	Choose a monster. Its {ATK} becomes 1.
//	{Paid Effect}Remove 2 counters from this:
//	Deal 1 damage to up to 2 monsters or players.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var hemoptysis = engine.CardDef{
	Ref:     "hemoptysis",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
}
