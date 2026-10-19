package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Classic Roller (Eternal Treasure Card)
//
//	{Tap Effect}Choose a dice roll. Its controller rerolls it, but rolls 2 dice instead. They choose another player. That player chooses one of the rolls as the result.
//	At the end of your turn, recharge this.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var classicRoller = engine.CardDef{
	Ref:       "classic_roller",
	Kind:      engine.TreasureCard,
	Copies:    1,
	Eternal:   true,
	ThreePlus: true,
	Outside:   true,
	Tap:       true,
}
