package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The card's own lines; their words are in en.json (ADR 021).
var (
	lodestoneWhichMonsterGains = engine.NewKey("card.lodestone.which_monster_gains_1_dc")
)

// Lodestone (Passive Treasure Card)
//
//	Each time a player rolls a ❸, choose a monster. It gains -1{DC}, till end of turn.
var lodestone = engine.CardDef{
	Ref:    "lodestone",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time a player rolls a ❸, choose a monster. It gains -1 DC, till end of turn.",
		Trigger: engine.OnRollOf(3),
		Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
			if ms := monstersInPlay(c.G); a[0] >= 0 && a[0] < len(ms) {
				c.G.Boosts = append(c.G.Boosts, engine.Boost{Source: c.Source, Stat: engine.StatMonsterDC, Player: engine.NoPlayer, Object: ms[a[0]], Amount: -1})
			}
		}, engine.Question{Text: lodestoneWhichMonsterGains, Options: func(c *engine.Ctx, _ []int) []engine.Msg {
			return labelsOf(c.G, monstersInPlay(c.G))
		}})},
	}},
}
