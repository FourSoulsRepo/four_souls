package b2

import (
	"reflect"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestChampionBelt(t *testing.T) {
	tb := itemTable(t, items("champion_belt"))
	if n := len(tb.G.Allowed(0)); n == 0 {
		t.Fatal("nothing allowed")
	}
	tb.Attack("fly", 6)
	tb.Attack("leech", 6) // a second attack
	if tb.G.Turn.BonusAttacksUsed != 1 {
		t.Errorf("bonus attacks used = %d, want 1", tb.G.Turn.BonusAttacksUsed)
	}
}

// The view names Champion Belt's +1 ATK and its condition once the first
// attack roll hits, beside Meat!'s +1 to attack rolls (V-41, 15.54).
func TestChampionBeltSourceDuringFirstRoll(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", "champion_belt", "meat"), seat("cain")},
		Monsters: items("big_spider"),
	}, Set)
	cards := func() []string {
		var out []string
		for _, s := range tb.G.View(engine.Viewer{Kind: engine.ViewSpectator}).Players[0].Sources {
			out = append(out, s.Stat+" "+string(s.Card)+" "+s.Lasts+" "+s.When)
		}
		return out
	}
	if got, want := cards(), []string{"attack_roll meat in_play "}; !reflect.DeepEqual(got, want) {
		t.Errorf("before the attack: sources = %q, want %q", got, want)
	}
	tb.G.ForceRolls(6)
	attackRoll(tb, "big_spider")
	tb.Pass(0)
	tb.Pass(1) // the roll resolves: a hit, its damage on the stack
	want := []string{"atk champion_belt in_play first attack roll", "attack_roll meat in_play "}
	if got := cards(); !reflect.DeepEqual(got, want) {
		t.Errorf("first attack roll hit: sources = %q, want %q", got, want)
	}
}

// The attack's odds count Champion Belt's +1 ATK before the first roll
// is made, and name it with its condition (15.94).
func TestChampionBeltInTheOdds(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", "champion_belt"), seat("cain")},
		Monsters: items("big_spider"),
	}, Set)
	attackRoll(tb, "big_spider")
	odds := tb.G.View(engine.Viewer{Kind: engine.ViewSpectator}).Attack.Odds
	if odds == nil {
		t.Fatal("no odds during the attack")
	}
	if want := tb.G.PlayerATK(0) + 1; odds.ATK != want {
		t.Errorf("odds ATK = %d, want %d: the first roll counts the belt", odds.ATK, want)
	}
	found := false
	for _, s := range odds.Sources {
		found = found || (s.Card == "champion_belt" && s.Stat == "atk" && s.When == "first attack roll")
	}
	if !found {
		t.Errorf("sources %+v, want Champion Belt's +1 ATK for the first attack roll", odds.Sources)
	}
}
