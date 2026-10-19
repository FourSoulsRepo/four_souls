package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Might For The Meek (Room Card)
//
//	Players who control the fewest souls or tied for fewest have +1{ATK} and may attack an additional time on their turn.
var mightForTheMeek = engine.CardDef{
	Ref:    "might_for_the_meek",
	Kind:   engine.RoomCard,
	Copies: 1,
	Statics: []engine.Static{
		playersWith(engine.StatPlayerATK, 1, fewestOf(soulCount)),
		playersWith(engine.StatAttacks, 1, fewestOf(soulCount)),
	},
}
