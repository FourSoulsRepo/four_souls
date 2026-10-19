package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Tax For The Mighty (Room Card)
//
//	When any player controls a soul, players who control the most souls or tied for the most must pay each other player 2¢ to attack.
var taxForTheMighty = engine.CardDef{
	Ref:    "tax_for_the_mighty",
	Kind:   engine.RoomCard,
	Copies: 1,
	Statics: []engine.Static{playersWith(engine.StatAttackTax, 2, func(g *engine.Game, p engine.PlayerID) bool {
		return soulsInPlay(g) && mostOf(soulCount)(g, p)
	})},
}

// soulsInPlay: any player controls a soul.
func soulsInPlay(g *engine.Game) bool {
	for _, pl := range g.Players {
		if g.SoulValue(pl.ID) > 0 {
			return true
		}
	}
	return false
}
