package r

import (
	"strings"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Bobby Bomb (Passive Treasure Card)
//
//	When you gain or purchase this, reveal the top 10 cards of the loot deck. Put each card with "Bomb" in its name in your hand and the rest on the bottom of the loot deck.
//	At the start of your turn, put a card with "Bomb" in its name from the loot discard into your hand.
var bobbyBomb = engine.CardDef{
	Ref:    "bobby_bomb",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When you gain or purchase this, reveal the top 10 cards of the loot deck. Put each card with \"Bomb\" in its name in your hand and the rest on the bottom of the loot deck.",
			Trigger: cardkit.WhenYouGainThis(),
			Effects: []engine.Effect{revealTenTakeNamed(isBomb)},
		},
		{
			Kind:    engine.Triggered,
			Text:    "At the start of your turn, put a card with \"Bomb\" in its name from the loot discard into your hand.",
			Trigger: engine.AtStartOfYourTurn(),
			Effects: []engine.Effect{lootDiscardToHand(isBomb)},
		},
	},
}

// isBomb: a loot card with "Bomb" in its name. Card IDs follow the names
// (bomb, troll_bombs, mega_bomb).
func isBomb(card engine.CardRef) bool { return strings.Contains(string(card), "bomb") }
