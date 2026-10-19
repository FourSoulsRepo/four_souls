package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Blood Puppy (Passive Treasure Card)
//
//	Each time a player dies, this levels up.
//	{LV4 Effect}You have +1{ATK}.
//	{LV8 Effect}At the start of your turn, deal 1 damage to a monster or player.
//	{LV10 Effect}At the start of your turn roll-
//	1-3: You die.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var bloodPuppy = engine.CardDef{
	Ref:    "blood_puppy",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
