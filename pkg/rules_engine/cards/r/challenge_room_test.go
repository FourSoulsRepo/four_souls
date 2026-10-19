package r

import "testing"

func TestChallengeRoom(t *testing.T) {
	tb := roomTable(t, "challenge_room")
	tb.Attack("fly", 6, 5) // the fly dies; the room rolls 5
	if tb.G.Turn.MustAttacks != 1 || tb.G.Turn.Attacks != 1 {
		t.Error("a 4-6 makes Isaac attack an additional time")
	}
}
