package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Bone (Eternal Treasure Card)
//
//	{Tap Effect}Put a counter on this.
//	{Paid Effect}Remove 1 counter from this:
//	Add +1 to a dice roll.
//	{Paid Effect}Remove 2 counters from this:
//	Deal 1 damage to a monster or player.
//	{Paid Effect}Remove 5 counters from this:
//	This becomes a soul and loses all abilities.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theBone = engine.CardDef{
	Ref:     "the_bone",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Soul:    1,
	Eternal: true,
	Outside: true,
	Tap:     true,
}
