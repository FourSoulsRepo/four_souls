package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Gemini (Boss Card)
//
//	While this is at 1{HP}, it has +1{ATK}.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var gemini = engine.CardDef{
	Ref:     "gemini",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 5}},
}
