package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Spider Webs (Room Card)
//
//	At the end of the turn, the active player deactivates their character.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var spiderWebs = engine.CardDef{
	Ref:    "spider_webs",
	Kind:   engine.RoomCard,
	Copies: 1,
}
