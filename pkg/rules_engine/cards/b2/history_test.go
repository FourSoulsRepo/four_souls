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
