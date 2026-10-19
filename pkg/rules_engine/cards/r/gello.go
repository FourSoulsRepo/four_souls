package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Gello (Eternal Treasure Card)
//
//	{Tap Effect}Choose a monster being attacked. Heal that monster to full {HP}, then deal damage equal to the number of {HP} healed in this way to another monster. If it's not your turn, cancel the attack and the active player may make an additional attack this turn.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var gello = engine.CardDef{
	Ref:     "gello",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
}
