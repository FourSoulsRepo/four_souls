package rulesengine

import "testing"

// Every target kind has words for its reason (15.33).
func TestTargetReason(t *testing.T) {
	for k := TargetPlayer; k <= TargetItemOrSoul; k++ {
		if _, ok := targetWords[k]; !ok {
			t.Errorf("target kind %d has no words", k)
		}
	}
	if got := targetReason(Choose(TargetMonsterOrPlayer)); got != "Pick a monster or a player." {
		t.Errorf("reason %q", got)
	}
	where := ChooseWhere(TargetPlayer, func(*Game, ObjectID, Chosen) bool { return true })
	if got := targetReason(where); got != "Pick a player. Only the glowing ones fit." {
		t.Errorf("filtered reason %q", got)
	}
}
