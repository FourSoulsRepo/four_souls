package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mulliboom (Basic Monster Card)
//
//	When this dies, the active player deals 3 damage to a player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
// Not generated: ATK is "4!".
var mulliboom = engine.CardDef{
	Ref:     "mulliboom",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      2,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 6}},
}
