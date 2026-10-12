package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Skeleton Key (Passive Treasure Card)
//
//	You may attack up to 99 times during your turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var skeletonKey = engine.CardDef{
	Ref:    "skeleton_key",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
