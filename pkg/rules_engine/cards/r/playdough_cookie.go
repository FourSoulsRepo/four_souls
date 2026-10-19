package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Playdough Cookie (Passive Treasure Card)
//
//	At the start of your turn, you may choose a non-eternal item. Put a counter on that item or remove one from it. Otherwise, gain 1¢.
var playdoughCookie = engine.CardDef{
	Ref:    "playdough_cookie",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "At the start of your turn, you may choose a non-eternal item. Put a counter on that item or remove one from it. Otherwise, gain 1¢.",
		Trigger: engine.AtStartOfYourTurn(),
		Effects: []engine.Effect{changeACounterOr(engine.GainCents(1))},
	}},
}
