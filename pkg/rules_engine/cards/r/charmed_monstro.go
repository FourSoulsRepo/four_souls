package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Charmed Monstro (Boss Card)
//
//	If the active player would gain this soul, they instead choose another player. That player gains this soul.
var charmedMonstro = engine.CardDef{
	Ref:     "charmed_monstro",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      3,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 2}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "If the active player would gain this soul, they instead choose another player. That player gains this soul.",
		Trigger: engine.AfterThisRewards(),
		Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
			if ps := othersOfActive(c.G); a[0] >= 0 && a[0] < len(ps) && c.G.Object(c.Source).Zone.Kind == engine.ZoneOutside {
				c.G.GainSoul(c.Source, ps[a[0]])
			}
		}, engine.Question{Text: "Who gains Charmed Monstro's soul?", Options: func(c *engine.Ctx, _ []int) []string {
			return playerLabels(c.G, othersOfActive(c.G))
		}})},
	}},
}
