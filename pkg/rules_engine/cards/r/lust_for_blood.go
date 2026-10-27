package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The card's own lines; their words are in en.json (ADR 021).
var (
	lustForBloodWhichPlayerMust = engine.NewKey("card.lust_for_blood.which_player_must_be_attacked")
)

// Lust For Blood! (Good Event Card)
//
//	Choose a player. You must make an additional attack against them. They have 4+{DC} for the attack.
//	When that player dies this turn, they give the active player the item they would destroy for the death penalty.
var lustForBlood = engine.CardDef{
	Ref:     "lust_for_blood",
	Kind:    engine.EventCard,
	Copies:  1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Choose a player. You must make an additional attack against them. They have 4+ DC for the attack. When that player dies this turn, they give the active player the item they would destroy for the death penalty.",
			Trigger: engine.WhenThisEntersPlay(),
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				others := othersInOrder(c)
				if a[0] < 0 || a[0] >= len(others) {
					return
				}
				p := others[a[0]]
				c.G.Boosts = append(c.G.Boosts, engine.Boost{Source: c.Source, Stat: engine.StatPlayerAttackDC, Player: p, Amount: 4})
				c.G.ForceAttacks(1, false)
				c.G.Turn.MustAttack = c.G.Players[p].Character
				c.G.GivePenaltyItem(p, c.G.Turn.Active)
			}, engine.Question{Text: lustForBloodWhichPlayerMust, Options: func(c *engine.Ctx, _ []int) []engine.Msg {
				return playerLabels(c.G, othersInOrder(c))
			}})},
		},
	},
}
