package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Friendly Ball (Active Treasure Card)
//
//	{Tap Effect}Put a monster not being attacked under this if there are no cards under this.
//	{Tap Effect}Put a monster from under this in a monster slot not being attacked. The active player must make an additional attack on it this turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var friendlyBall = engine.CardDef{
	Ref:    "friendly_ball",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
