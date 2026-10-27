package rulesengine

import "testing"

// Every target kind has words for its reason (15.33).
func TestTargetReason(t *testing.T) {
	for k := TargetPlayer; k <= TargetItemOrSoul; k++ {
		if _, ok := refusePick[k]; !ok {
			t.Errorf("target kind %d has no words", k)
		}
	}
	if got := targetReason(Choose(TargetMonsterOrPlayer)); got.String() != "refuse.pick.monster_or_player" {
		t.Errorf("reason %v", got)
	}
	where := ChooseWhere(TargetPlayer, func(*Game, ObjectID, Chosen) bool { return true })
	if got := targetReason(where); got.String() != "refuse.pick.player.only" {
		t.Errorf("filtered reason %v", got)
	}
}
