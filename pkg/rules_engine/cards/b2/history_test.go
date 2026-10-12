package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// resolvedItems are the stack items that resolved, as the history strip
// gets them (15.22).
func resolvedItems(events []engine.Event) []engine.StackItem {
	var out []engine.StackItem
	for _, e := range events {
		if e.Kind == engine.EvStackResolved && e.Item != nil {
			out = append(out, *e.Item)
		}
	}
	return out
}

func TestHistorySplitsARollFromItsResult(t *testing.T) {
	tb := lootTable(t, "blank_rune")
	tb.G.ForceRolls(1)
	items := resolvedItems(tb.Play(0, "blank_rune"))
	var texts []string
	for _, it := range items {
		texts = append(texts, it.Text)
	}
	want := []string{"Roll.", "", "1: Each player gains 1¢."}
	if len(texts) != len(want) {
		t.Fatalf("resolved texts %q, want %q", texts, want)
	}
	for i := range want {
		if texts[i] != want[i] {
			t.Errorf("resolved item %d says %q, want %q", i, texts[i], want[i])
		}
	}
}

func TestHistoryNamesTheAttackTarget(t *testing.T) {
	tb := lootTable(t)
	items := resolvedItems(tb.Attack("fly", 6))
	var roll, damage bool
	for _, it := range items {
		if !it.Attack {
			continue
		}
		if it.TargetCard != "fly" {
			t.Errorf("attack item %+v names target %q, want fly", it, it.TargetCard)
		}
		roll = roll || it.Kind == engine.StackRoll
		damage = damage || it.Kind == engine.StackDamage
	}
	if !roll || !damage {
		t.Errorf("want an attack roll and combat damage, got %+v", items)
	}
}

func TestStackAndHistoryShowOnlyTheUsedAbility(t *testing.T) {
	tb := monsterTable(t, []engine.CardRef{"fly"}, seat("the_forgotten", "the_bone"), seat("cain"))
	tb.G.Object(tb.Find(0, "the_bone")).Counters = []engine.Counter{{Count: 7}}
	want := "Remove 2 counters from this: Deal 1 damage to a monster or player."
	tb.Start(0, "the_bone", 2)
	tb.Choose("fly")
	stack := tb.G.View(engine.Viewer{Kind: engine.ViewSpectator}).Stack
	if len(stack) == 0 || stack[len(stack)-1].Text != want {
		t.Fatalf("stack %+v, want The Bone's third ability on top", stack)
	}
	var found bool
	for _, it := range resolvedItems(tb.Settle()) {
		if it.Ability.Card == "the_bone" {
			found = true
			if it.Text != want {
				t.Errorf("history says %q, want %q", it.Text, want)
			}
		}
		if it.Kind == engine.StackDamage && it.TargetCard != "fly" {
			t.Errorf("The Bone's damage names %q, want fly", it.TargetCard)
		}
	}
	if !found {
		t.Error("The Bone's ability never resolved")
	}
}

func TestStackShowsOnlyTheRollUntilItsResult(t *testing.T) {
	tb := lootTable(t, "blank_rune")
	in := engine.Intent{Player: 0, Kind: engine.IntentPlayLoot, Objects: []engine.ObjectID{tb.Find(0, "blank_rune")}}
	if _, err := tb.G.Apply(in); err != nil {
		t.Fatal(err)
	}
	stack := tb.G.View(engine.Viewer{Kind: engine.ViewSpectator}).Stack
	if len(stack) == 0 || stack[len(stack)-1].Text != "Roll." {
		t.Errorf("stack %+v, want Blank Rune on top saying Roll.", stack)
	}
}
