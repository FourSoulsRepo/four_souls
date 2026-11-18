package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// TestXIXTheSunGoesToTheBottom: the card goes to the bottom of the loot
// deck. The extra turn is the file xix_the_sun/takes-an-extra-turn.json.
func TestXIXTheSunGoesToTheBottom(t *testing.T) {
	tb := lootTable(t, "xix_the_sun")
	tb.Play(0, "xix_the_sun")
	if c := tb.G.Object(tb.G.Decks[engine.LootDeck][0]).Card; c != "xix_the_sun" {
		t.Errorf("bottom of the loot deck is %s", c)
	}
}
