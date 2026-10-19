package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Soul Of Envy (Bonus Soul Card)
//
//	The first time a player controls their 3rd soul, the active player chooses a player who controls the fewest souls or tied for fewest. That player gains this soul.
var soulOfEnvy = engine.CardDef{
	Ref:    "soul_of_envy",
	Kind:   engine.BonusSoulCard,
	Copies: 1,
	Soul:   1,
	Abilities: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "The first time a player controls their 3rd soul, the active player chooses a player who controls the fewest souls or tied for fewest. That player gains this soul.",
		Trigger: engine.Trigger{On: engine.EvGainedSoul, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
			return g.Object(self).CountersOf(firstTime) == 0 && g.SoulValue(e.Player) >= 3
		}},
		Effects: []engine.Effect{fewestGetsThis(func(g *engine.Game, p engine.PlayerID) int { return g.SoulValue(p) })},
	}},
}
