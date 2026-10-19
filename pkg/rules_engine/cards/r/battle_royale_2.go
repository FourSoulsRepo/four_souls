package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Battle Royale!!!! (Room Card)
//
//	The active player may attack other players. Attacked players have 4+{DC}.
//	When a player dies, if that player was attacked this turn, that player gives the active player the item they would destroy for the death penalty.
var battleRoyale2 = engine.CardDef{
	Ref:       "battle_royale_2",
	Kind:      engine.RoomCard,
	Copies:    1,
	Statics:   []engine.Static{playersCanBeAttacked(4)},
	Abilities: []engine.Ability{attackedPlayersGivePenaltyItem},
}
