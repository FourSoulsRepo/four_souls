package rulesengine

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

// TestBarrierKinds: every event kind is sorted into barrier or not
// (ADR 020, item 4.9). A new kind fails until barrierKinds and the
// list here name it. testdata/event_kinds.txt is the list of every
// kind (TestEventKindsListed).
func TestBarrierKinds(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "event_kinds.txt"))
	if err != nil {
		t.Fatal(err)
	}
	all := strings.Fields(string(data))
	// The Soft barriers, pinned: changing one is a rules decision.
	want := []string{
		"attack_roll_repeated", "card_revealed", "character_dealt", "damaged", "dice_rolled", "died",
		"first_player", "first_roll_tied", "gained_treasure", "game_started", "game_won", "looked_at",
		"object_died", "random_pick", "roll_changed", "roll_repeated", "roll_resolved", "roll_would_resolve",
		"turn_ended", "turn_ended_early", "turn_skipped",
	}
	var got []string
	for _, k := range all {
		barrier, ok := barrierKinds[EventKind(k)]
		if !ok {
			t.Errorf("event kind %s is not sorted into barrierKinds (barrier.go)", k)
		}
		if barrier {
			got = append(got, k)
		}
	}
	for k := range barrierKinds {
		if !slices.Contains(all, string(k)) {
			t.Errorf("barrierKinds names %s, which is no event kind", k)
		}
	}
	if !slices.Equal(got, want) {
		t.Errorf("barriers\ngot  %v\nwant %v", got, want)
	}
}

// TestBarrierPrivate: an event only one player sees is a barrier,
// whatever its kind (item 4.4), but for a hand card put on the loot
// deck: no deck card shows (item 4.22).
func TestBarrierPrivate(t *testing.T) {
	for _, c := range []struct {
		e    Event
		want bool
	}{
		{Event{Kind: EvLooted, Private: true}, true},
		{Event{Kind: EvLooted, Text: "from the discard"}, false},
		{Event{Kind: EvMovedToDeck, Private: true}, false},
		{Event{Kind: EvMovedToDeck}, false},
		{Event{Kind: EvGaveCard, Private: true}, true},
		{Event{Kind: EvLookedAt, Private: true}, true},
		{Event{Kind: EvLootPlayed}, false},
		{Event{Kind: EvDiceRolled}, true},
		{Event{Kind: EvAttackDeclared}, false},
	} {
		if got := Barrier(c.e); got != c.want {
			t.Errorf("Barrier(%s, private %v) = %v", c.e.Kind, c.e.Private, got)
		}
	}
}

// TestBarrierHandToDeck: the private emits as the engine makes them. A
// hand card put on the loot deck is no barrier: no deck card shows and
// nobody else sees it (ADR 020, item 4.22; N-12, item 24). A loot draw
// and a hand card given to another player still are.
func TestBarrierHandToDeck(t *testing.T) {
	g := toAction(t, 2)
	for _, c := range []struct {
		name string
		do   func()
		kind EventKind
		want bool
	}{
		{"hand to deck", func() { g.HandToDeckTop(0, g.Players[0].Hand[0]) }, EvMovedToDeck, false},
		{"loot", func() { g.loot(0, 1) }, EvLooted, true},
		{"give", func() { g.GiveHandCard(0, 1, g.Players[0].Hand[0]) }, EvGaveCard, true},
		{"look", func() { g.LookAt(0, g.Players[1].Hand[0]) }, EvLookedAt, true},
	} {
		g.events = nil
		c.do()
		ev := g.takeEvents()
		i := slices.IndexFunc(ev, func(e Event) bool { return e.Kind == c.kind })
		if i < 0 || !ev[i].Private {
			t.Fatalf("%s: events %+v, want a private %s", c.name, ev, c.kind)
		}
		if got := slices.ContainsFunc(ev, Barrier); got != c.want {
			t.Errorf("%s: barrier %v, want %v: %+v", c.name, got, c.want, ev)
		}
	}
}
