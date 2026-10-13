package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Guppy’s Eye (Passive Treasure Card)
//
//	Each other player plays with their hand revealed.
//	-Guppy- The first player to control 2 or more Guppy items gains the Soul of Guppy.
var guppysEye = engine.CardDef{
	Ref:    "guppys_eye",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Guppy:  true,
	Statics: []engine.Static{{Stat: engine.StatRevealHand, Amount: 1, Applies: func(g *engine.Game, self engine.ObjectID, p engine.PlayerID, _ engine.ObjectID) bool {
		return p != g.Object(self).Controller // "Each other player plays with their hand revealed."
	}}},
}
