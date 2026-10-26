package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// The card's own lines; their words are in en.json (ADR 021).
var (
	angelicInterventionCancelTheAttack = engine.NewKey("card.angelic_intervention.cancel_the_attack")
)

// Angelic Intervention (Room Card)
//
//	Each time the active player attacks the top of the monster deck, after putting it in a monster slot, they may cancel their attack.
var angelicIntervention = engine.CardDef{
	Ref:    "angelic_intervention",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "Each time the active player attacks the top of the monster deck, after putting it in a monster slot, they may cancel their attack.",
		Trigger: engine.Trigger{On: engine.EvAttackTarget, Match: func(g *engine.Game, _ engine.ObjectID, _ engine.Event) bool {
			return g.Attack.Deck
		}},
		Effects: []engine.Effect{cardkit.May(angelicInterventionCancelTheAttack, engine.CancelAttack())},
	}},
}
