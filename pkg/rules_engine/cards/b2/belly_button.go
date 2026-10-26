package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The card's own lines; their words are in en.json (ADR 021).
var (
	bellyButtonRechargeYourCharacter = engine.NewKey("card.belly_button.recharge_your_character")
)

// Belly Button (Passive Treasure Card)
//
//	You may play an additional loot card on your turn.
//	Each time you take damage, you may recharge your character.
var bellyButton = engine.CardDef{
	Ref:     "belly_button",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Statics: []engine.Static{engine.YouHave(engine.StatLootPlays, 1)},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time you take damage, you may recharge your character.",
			Trigger: engine.WhenYouTakeDamage(),
			Effects: []engine.Effect{may(bellyButtonRechargeYourCharacter, engine.EffectFunc(func(c *engine.Ctx) { c.G.Recharge(c.G.Players[c.Controller].Character) }))},
		},
	},
}
