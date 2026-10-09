package b2

import "testing"

func TestBigSpider(t *testing.T) {
	tb := slayTable(t, "big_spider", seat("isaac"), seat("cain"))
	killWithSixes(t, tb)
	if n := tb.G.Turn.DeckAttacks; n != 1 {
		t.Errorf("deck attacks = %d, want 1", n)
	}
}
