package r

import "testing"

func TestFloorSpikes(t *testing.T) {
	tb := roomTable(t, "floor_spikes")
	tb.Attack("leech", 1)
	if got := tb.G.Players[0].Damage; got < 2 {
		t.Errorf("Isaac's damage = %d, want the leech's 1 and the spikes' 1", got)
	}
}
