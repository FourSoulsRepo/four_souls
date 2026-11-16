package main

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/FourSoulsRepo/four_souls/internal/bot"
)

// tiny are the flags of a short run: one short game of each size, a few
// cards, few samples.
func tiny(out string, more ...string) []string {
	return append([]string{"-games", "1", "-cap", "4", "-every", "3", "-samples", "3", "-draws", "1", "-out", out}, more...)
}

// TestRunRepeats: a short run writes a table that parses and holds the
// asked cards; the same seed writes the same file; -missing keeps a
// full table as it is and measures nothing (Decision 20.7).
func TestRunRepeats(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "values.json")
	if err := os.WriteFile(out, []byte(`{"Cards":[]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	cardsFlag := []string{"-cards", "a_penny,bomb,breakfast,the_battery"}
	if err := run(context.Background(), tiny(out, cardsFlag...), io.Discard); err != nil {
		t.Fatal(err)
	}
	first, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	f, err := bot.ParseValues(first)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]string{}
	for _, e := range f.Cards {
		got[e.Key] = e.Kind
	}
	want := map[string]string{"a_penny": bot.KindLoot, "bomb": bot.KindLoot, "breakfast": bot.KindPassive, "the_battery": bot.KindActive}
	for k, kind := range want {
		if got[k] != kind {
			t.Errorf("%s: kind %q, want %q, in\n%s", k, got[k], kind, first)
		}
	}
	if err := run(context.Background(), tiny(out, cardsFlag...), io.Discard); err != nil {
		t.Fatal(err)
	}
	second, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(first, second) {
		t.Errorf("the same seed wrote another file:\n%s\n%s", first, second)
	}
	// A table of every card: -missing measures nothing, writes nothing.
	var all []measured
	for _, c := range allCards() {
		all = append(all, measured{card: c})
	}
	full, err := encode(merge(bot.ValueFile{}, all, config{seed: 1}))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(out, full, 0o600); err != nil {
		t.Fatal(err)
	}
	var log bytes.Buffer
	if err := run(context.Background(), tiny(out, "-missing"), &log); err != nil {
		t.Fatal(err)
	}
	third, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(full, third) || !bytes.Contains(log.Bytes(), []byte("0 cards to measure")) {
		t.Errorf("-missing on a full table: %s", log.String())
	}
}

// TestAllCards: loot cards once as loot, trinkets also as items,
// treasure as items; no starting item.
func TestAllCards(t *testing.T) {
	kinds := map[string][]string{}
	for _, c := range allCards() {
		kinds[string(c.ref)] = append(kinds[string(c.ref)], c.kind)
	}
	if k := kinds["a_penny"]; len(k) != 1 || k[0] != bot.KindLoot {
		t.Errorf("a_penny: %v", k)
	}
	if k := kinds["the_battery"]; len(k) != 1 || k[0] != bot.KindActive {
		t.Errorf("the_battery: %v", k)
	}
	if _, ok := kinds["the_d6"]; ok {
		t.Error("a starting item is measured; it is in no deck")
	}
}
