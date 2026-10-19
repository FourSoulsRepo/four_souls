package r

import (
	"strings"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Mom’s Bottle of Pills (Passive Treasure Card)
//
//	When you gain or purchase this, reveal the top 10 cards of the loot deck. Put each card named Pills! in your hand and the rest on the bottom of the deck.
//	At the start of your turn, put a card named Pills! from the loot discard into your hand.
var momsBottleOfPills = engine.CardDef{
	Ref:    "moms_bottle_of_pills",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When you gain or purchase this, reveal the top 10 cards of the loot deck. Put each card named Pills! in your hand and the rest on the bottom of the deck.",
			Trigger: cardkit.WhenYouGainThis(),
			Effects: []engine.Effect{revealTenTakeNamed(isPills)},
		},
		{
			Kind:    engine.Triggered,
			Text:    "At the start of your turn, put a card named Pills! from the loot discard into your hand.",
			Trigger: engine.AtStartOfYourTurn(),
			Effects: []engine.Effect{lootDiscardToHand(isPills)},
		},
	},
}

// isPills: a card named Pills!; every such card's ID starts with "pills"
// (pills, pills_2, pills_g2, pills_r).
func isPills(card engine.CardRef) bool { return strings.HasPrefix(string(card), "pills") }
