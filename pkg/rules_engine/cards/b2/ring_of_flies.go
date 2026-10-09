package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Ring Of Flies (Basic Monster Card)
//
//	Each time the attacking player rolls an attack roll of 3, they must steal a loot card from from another player at random.
var ringOfFlies = engine.CardDef{
	Ref:     "ring_of_flies",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}, {Kind: engine.RewardCents, Amount: 3}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time the attacking player rolls an attack roll of 3, they must steal a loot card from another player at random.",
			Trigger: attackRollOnThis(3),
			Effects: []engine.Effect{engine.Ask(ringSteal, ringPlayer, ringCard)},
		},
	},
}

var ringPlayer = engine.Question{Text: "Steal from which player?", Options: func(c *engine.Ctx, _ []int) []string {
	var out []string
	for _, p := range otherPlayers(c) {
		if len(c.G.Players[p].Hand) > 0 {
			out = append(out, playerLabel(c.G, p))
		}
	}
	return out
}}

// withHands lists the other players who have loot cards.
func withHands(c *engine.Ctx) []engine.PlayerID {
	var out []engine.PlayerID
	for _, p := range otherPlayers(c) {
		if len(c.G.Players[p].Hand) > 0 {
			out = append(out, p)
		}
	}
	return out
}

var ringCard = engine.Question{Random: true, Text: "Which card?", Options: func(c *engine.Ctx, a []int) []string {
	if a[0] < 0 {
		return nil
	}
	return make([]string, len(c.G.Players[withHands(c)[a[0]]].Hand))
}}

func ringSteal(c *engine.Ctx, a []int) {
	if a[0] < 0 || a[1] < 0 {
		return
	}
	from := withHands(c)[a[0]]
	c.G.GiveHandCard(from, c.Controller, c.G.Players[from].Hand[a[1]])
}
