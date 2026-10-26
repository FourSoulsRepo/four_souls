package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The card's own lines; their words are in en.json (ADR 021).
var (
	theClickerWhichCharacter = engine.NewKey("card.the_clicker.which_character")
)

// The Clicker (One-Use Treasure Card)
//
//	{Tap Effect}Destroy this. If you do, choose a character card from outside the game. Replace your character with it and your starting item with the chosen card's starting item.
var theClicker = engine.CardDef{
	Ref:    "the_clicker",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{{
		Kind:  engine.Activated,
		Text:  "↷: Destroy this. If you do, choose a character card from outside the game. Replace your character with it and your starting item with the chosen card's starting item.",
		Costs: []engine.Cost{engine.Tap()},
		Effects: []engine.Effect{engine.DestroyThis(engine.Ask(func(c *engine.Ctx, a []int) {
			if chars := c.G.CharactersOutside(); a[0] >= 0 && a[0] < len(chars) {
				c.G.ReplaceCharacter(c.Controller, chars[a[0]])
			}
		}, engine.Question{Text: theClickerWhichCharacter, Options: func(c *engine.Ctx, _ []int) []engine.Msg {
			var out []engine.Msg
			for _, ch := range c.G.CharactersOutside() {
				out = append(out, engine.CardOption(ch))
			}
			return out
		}}))},
	}},
}
