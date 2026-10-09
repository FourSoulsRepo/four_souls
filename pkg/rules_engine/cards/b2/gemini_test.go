package b2

import "testing"

func TestGemini(t *testing.T) {
	tb := slayTable(t, "gemini", seat("isaac"), seat("cain"))
	g, _ := tb.G.Monsters[0].TopOf()
	tb.G.Object(g).Damage = 2
	if a := tb.G.MonsterATK(g); a != 2 {
		t.Errorf("ATK at 1 HP = %d, want 2", a)
	}
}
