package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// The card's own lines; their words are in en.json (ADR 021).
var (
	momsHeartWhoWins = engine.NewKey("card.moms_heart.who_wins")
)

// Mom’s Heart! (Epic Boss Card)
//
//	When a player gains this soul, choose a player who controls the most souls or tied for the most. That player wins.
var momsHeart = engine.CardDef{
	Ref:    "moms_heart",
	Kind:   engine.MonsterCard,
	Copies: 1,
	HP:     8,
	DC:     4,
	ATK:    2,
	Soul:   2,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "When a player gains this soul, choose a player who controls the most souls or tied for the most. That player wins.",
		Trigger: engine.Trigger{On: engine.EvGainedSoul, Match: func(_ *engine.Game, self engine.ObjectID, e engine.Event) bool { return e.Object == self }},
		Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
			if ps := mostSouls(c.G); a[0] >= 0 && a[0] < len(ps) {
				c.G.WinGame(ps[a[0]])
			}
		}, engine.Question{Text: momsHeartWhoWins, Options: func(c *engine.Ctx, _ []int) []engine.Msg {
			var out []engine.Msg
			for _, p := range mostSouls(c.G) {
				out = append(out, playerLabel(c.G, p))
			}
			return out
		}})},
	}},
}
