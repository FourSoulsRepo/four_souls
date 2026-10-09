package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Holy Keeper Head (Holy/Charmed Monster Card)
//
//	Each time a player rolls a ❹, they gain 2¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
// Not generated: reward "Roll- Gain x" Coin.
var holyKeeperHead = engine.CardDef{
	Ref:    "holy_keeper_head",
	Kind:   engine.MonsterCard,
	Copies: 1,
	HP:     2,
	DC:     4,
	ATK:    1,
}
