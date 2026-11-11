package challenge

import engine "github.com/FourSoulsRepo/rules_engine"

// Trick (Challenge Card)
//
//	At the start of each turn, each player votes trick or treat. If treat wins, flip this. Each non-active player that votes treat loots x, where x is equal to the number of souls the active player controls.
//	At the start of each attack against a monster, the active player rolls-
//	1-2: Damage you would take is increased by 1 till end of turn.
//	3-4: Reveal cards from the top of the monster deck until you reveal a non-event monster card. The monster being attacked gains its abilities till end of turn. Put any revealed cards on the bottom of the deck.
//	5-6: The rewards of that monster become 1¢ till end of turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
// A question or option it asks is a key, engine.NewKey("card.trick_treat_comp.<name>"),
// with its words and a note in frontend/src/locales/en.json (ADR 021).
var trickTreatComp = engine.CardDef{
	Ref:    "trick_treat_comp",
	Kind:   engine.ChallengeCard,
	Copies: 1,
	Back:   "treat_comp",
}

// treatComp is the back of trick_treat_comp: Treat.
//
//	At the start of each turn, each player votes trick or treat. If trick wins, flip this. Each non-active player that votes treat loots x, where x is equal to the number of souls the active player controls.
//	At the start of each attack against a monster, the active player rolls-
//	1-2: Add 3¢ to the rewards of that monster till end of turn.
//	3-4: Add loot 1 to the rewards of that monster till end of turn.
//	5-6: Add +1 treasure to the rewards of that monster till end of turn.
var treatComp = engine.CardDef{
	Ref:     "treat_comp",
	Kind:    engine.ChallengeCard,
	Outside: true,
}
