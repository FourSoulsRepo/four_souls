package r

import "testing"

func TestHolyPortal(t *testing.T) {
	tb := fight(t, "holy_portal")
	tb.Attack("fly", 6)
	if tb.G.Turn.Attacks != 1 {
		t.Error("a monster died: one more attack")
	}
}
