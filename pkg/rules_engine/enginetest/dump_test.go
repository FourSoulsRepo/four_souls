package enginetest

import (
	"encoding/json"
	"io/fs"
	"os"
	"slices"
	"strings"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/b2"
)

// readDump reads the recording of test name in this package's folder.
func readDump(t *testing.T, dir, name string) (Dump, map[string]any) {
	t.Helper()
	data, err := fs.ReadFile(os.DirFS(dir), "enginetest/"+name+".json")
	if err != nil {
		t.Fatal(err)
	}
	var d Dump
	var raw map[string]any
	if err := json.Unmarshal(data, &d); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatal(err)
	}
	return d, raw
}

func hasTodo(d Dump, part string) bool {
	return slices.ContainsFunc(d.Todo, func(s string) bool { return strings.Contains(s, part) })
}

// bookOfSin: Isaac with Book of Sin, Cain; Book of Sin's roll of 5 gives
// Isaac +1 HP till end of turn.
func bookOfSin(t *testing.T) *Table {
	t.Helper()
	tb := NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{Seat("isaac", "book_of_sin"), Seat("cain")}}, b2.Set)
	tb.G.ForceRolls(5)
	return tb
}

// The recording holds the setup, what the test set by hand that a
// setup can say, the dice, the steps, and a todo for the rest; never
// an expect.
func TestDumpRecordsATest(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(DumpEnv, dir)
	t.Run("hand_set", func(t *testing.T) {
		tb := bookOfSin(t)
		tb.G.PutOnDeck(engine.LootDeck, "bomb")
		tb.G.Turn.Attacks = 2
		tb.G.Object(tb.Find(0, "book_of_sin")).Eternal = true
		tb.G.Players[1].SkipTurns = 1 // no setup field: Cain skips his turn
		tb.Activate(0, "book_of_sin", 0)
		tb.G.Players[0].Cents = 3 // between steps
		tb.EndTurn()
		if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentEndTurn}); err != nil {
			t.Fatal(err)
		}
	})
	d, raw := readDump(t, dir, "TestDumpRecordsATest/hand_set")
	if _, ok := raw["expect"]; ok {
		t.Error("the dump has expect; a tool never writes it (23 rule 3)")
	}
	if d.Title != "Dump Records ATest: hand set" || !slices.Equal(d.Sets, []string{"b2"}) || d.Source == nil || len(d.Source) != 0 {
		t.Errorf("title %q, sets %v, source %v; want the test name, [b2], an empty list", d.Title, d.Sets, d.Source)
	}
	s := d.Setup
	if !slices.Equal(s.Rolls, []int{5}) || !slices.Equal(s.Decks["loot"], []engine.CardRef{"bomb"}) ||
		s.Attacks == nil || *s.Attacks != 2 || !slices.Equal(s.Players[0].Eternal, []engine.CardRef{"book_of_sin"}) {
		t.Errorf("setup %+v; want rolls [5], loot deck [bomb], 2 attacks, an eternal book_of_sin", s)
	}
	if s.LootPlays != nil || s.Purchases != nil || s.Players[0].Cents != 0 {
		t.Errorf("setup %+v; want only what the test changed", s)
	}
	want := []DumpStep{{Player: seatOf(0), Kind: "activate", Card: "book_of_sin"}, {Player: seatOf(0), Kind: "end_turn"}, {Player: seatOf(0), Kind: "end_turn"}}
	if got := mustJSON(t, d.Steps); got != mustJSON(t, want) {
		t.Errorf("steps %s; want %s", got, mustJSON(t, want))
	}
	for _, part := range []string{"no setup field says it: Players[1] (the dumped setup's → the test's): SkipTurns none → 1", "before step 2: Players[0] changed outside a Table helper: Cents none → 3", "step 3 comes from a direct G.Apply"} {
		if !hasTodo(d, part) {
			t.Errorf("todo %q; want one with %q", d.Todo, part)
		}
	}
}

func mustJSON(t *testing.T, v any) string {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// A dump with todo removed and source and expect added by hand runs
// as a situation file (24.7 Done).
func TestDumpRunsWithExpect(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(DumpEnv, dir)
	t.Run("book_of_sin", func(t *testing.T) {
		bookOfSin(t).Activate(0, "book_of_sin", 0)
	})
	_, raw := readDump(t, dir, "TestDumpRunsWithExpect/book_of_sin")
	delete(raw, "todo")
	raw["source"] = []string{"card book_of_sin"}
	// Book of Sin: "5-6: Gain +1 HP till end of turn"; Isaac has 2 HP.
	raw["expect"] = map[string]any{"allowed": true, "after": []any{map[string]any{"player": 0, "hp": 3, "deactivated": []string{"book_of_sin"}}}}
	data, err := json.Marshal(raw)
	if err != nil {
		t.Fatal(err)
	}
	ans, err := engine.RunSituation(data, b2.Set)
	if err != nil {
		t.Fatal(err)
	}
	if ans.Matches == nil || !*ans.Matches {
		t.Errorf("mismatch %q", ans.Mismatch)
	}
}

// Without the env var a table records nothing and writes nothing.
func TestDumpOffWithoutEnv(t *testing.T) {
	t.Setenv(DumpEnv, "")
	tb := bookOfSin(t)
	tb.Activate(0, "book_of_sin", 0)
	if tb.rec != nil {
		t.Error("a table records with no " + DumpEnv)
	}
}

// State set by hand before a direct G.Apply as the first step is in
// the setup too: the game keeps its state from before its first
// applied intent (24.11; 24.8 pilot item 7.3, XX. Judgement's soul).
func TestDumpHandSetBeforeDirectApply(t *testing.T) {
	dir := t.TempDir()
	t.Setenv(DumpEnv, dir)
	t.Run("direct", func(t *testing.T) {
		tb := bookOfSin(t)
		tb.G.Object(tb.Find(0, "book_of_sin")).Eternal = true
		tb.G.Players[1].Cents = 4
		in := engine.Intent{Player: 0, Kind: engine.IntentActivate, Objects: []engine.ObjectID{tb.Find(0, "book_of_sin")}}
		if _, err := tb.G.Apply(in); err != nil {
			t.Fatal(err)
		}
		tb.Settle()
	})
	d, _ := readDump(t, dir, "TestDumpHandSetBeforeDirectApply/direct")
	s := d.Setup
	if !slices.Equal(s.Players[0].Eternal, []engine.CardRef{"book_of_sin"}) || s.Players[1].Cents != 4 {
		t.Errorf("setup %+v; want an eternal book_of_sin and Cain's 4 cents", s)
	}
	if s.Players[0].Cents != 0 || len(s.Players[0].Deactivated) != 0 {
		t.Errorf("setup %+v; want the state before the step, not after it", s)
	}
	want := []DumpStep{{Player: seatOf(0), Kind: "activate", Card: "book_of_sin"}}
	if got := mustJSON(t, d.Steps); got != mustJSON(t, want) {
		t.Errorf("steps %s; want %s", got, mustJSON(t, want))
	}
	if hasTodo(d, "not in setup") || hasTodo(d, "no setup field says it") {
		t.Errorf("todo %q; want no hand-set state left out", d.Todo)
	}
}
