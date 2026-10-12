package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Lost (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
//	At the end of your turn, recharge up to 1 character.
var theLost = engine.CardDef{
	Ref:          "the_lost",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           1,
	ATK:          1,
	Soul:         1,
	StartingItem: "holy_mantle",
	Tap:          true,
	Abilities: []engine.Ability{
		extraLoot,
		{
			Kind:    engine.Triggered,
			Text:    "At the end of your turn, recharge up to 1 character.",
			Trigger: engine.AtEndOfYourTurn(),
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				if ps := livingPlayers(c.G, engine.NoPlayer); a[0] >= 0 && a[0] < len(ps) {
					c.G.Recharge(c.G.Players[ps[a[0]]].Character)
				}
			}, engine.Question{Text: "Recharge which character?", Options: func(c *engine.Ctx, _ []int) []string {
				var out []string
				for _, p := range livingPlayers(c.G, engine.NoPlayer) {
					out = append(out, string(c.G.Object(c.G.Players[p].Character).Card))
				}
				return append(out, "none")
			}})},
		},
	},
}
