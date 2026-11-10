package challenge

import engine "github.com/FourSoulsRepo/rules_engine"

// How the Krampus Stole Christmas (Challenge Card)
//
//	If you would gain rewards from a monster or gain a soul, you may instead forego the reward or soul and put counters on the present of the respective type (¢, loot, treasure, or soul) equal to the amount of that resource foregone.
//	If a present is reduced to 0 counters, destroy it. Put the top 3 cards of the monster deck in minion slots.
//	The first time there are 5 or more souls controlled by all players, each present is flipped, becomes an item, and is gained by a player of the active player's choice.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
// A question or option it asks is a key, engine.NewKey("card.how_the_krampus_stole_christmas_norm.<name>"),
// with its words and a note in frontend/src/locales/en.json (ADR 021).
var howTheKrampusStoleChristmasNorm = engine.CardDef{
	Ref:    "how_the_krampus_stole_christmas_norm",
	Kind:   engine.ChallengeCard,
	Copies: 1,
}
