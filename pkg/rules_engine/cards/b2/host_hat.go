package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Host Hat (Active Treasure Card)
//
//	{Tap Effect}Prevent the next 1 damage you would take this turn. When you prevent damage this way, deal 1 damage to another player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var hostHat = engine.CardDef{
	Ref:    "host_hat",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
