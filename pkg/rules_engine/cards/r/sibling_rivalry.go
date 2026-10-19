package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Sibling Rivalry (Eternal Treasure Card)
//
//	{Tap Effect}Choose a non-active player. The next time the active player declares an attack this turn, the chosen player must make an attack roll after each attack roll the active player makes for the attack. If that monster dies this attack, the chosen player also gains the rewards.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var siblingRivalry = engine.CardDef{
	Ref:       "sibling_rivalry",
	Kind:      engine.TreasureCard,
	Copies:    1,
	Eternal:   true,
	ThreePlus: true,
	Outside:   true,
	Tap:       true,
}
