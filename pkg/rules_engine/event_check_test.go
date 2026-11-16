package rulesengine

import (
	"slices"
	"strings"
	"testing"
)

func eventCheckGame(t *testing.T) *Game {
	t.Helper()
	g, _, err := NewGame(Setup{Seed: 1, Players: 2, Sets: fuzzSets})
	if err != nil {
		t.Fatal(err)
	}
	g.takeEvents()
	return g
}

// 22.2: a card put on the bottom of its deck with no event fails the
// check; the same move with its top_choice passes.
func TestCheckEventsSilentDeckToBottom(t *testing.T) {
	g := eventCheckGame(t)
	snap := g.TakeSnapshot()
	g.DeckToBottom(LootDeck, g.DeckTop(LootDeck, 1)[0])
	err := g.CheckEvents(snap, g.takeEvents())
	if err == nil || !strings.Contains(err.Error(), "deck order") {
		t.Fatalf("a silent DeckToBottom passed the check: %v", err)
	}
	snap = g.TakeSnapshot()
	g.KeepOrBottom(0, LootDeck, g.DeckTop(LootDeck, 1)[0], true)
	if err := g.CheckEvents(snap, g.takeEvents()); err != nil {
		t.Fatalf("KeepOrBottom: %v", err)
	}
}

// 22.2: a move needs an event that names the new object, or the old
// one after the move; an event before the move does not count.
func TestCheckEventsSilentMove(t *testing.T) {
	g := eventCheckGame(t)
	id := g.Players[0].Hand[0]
	snap := g.TakeSnapshot()
	g.LookAt(1, id) // names the hand card before it moves
	g.Players[0].Hand = remove(g.Players[0].Hand, id)
	g.discard(id, LootDeck)
	if err := g.CheckEvents(snap, g.takeEvents()); err == nil || !strings.Contains(err.Error(), "moved from") {
		t.Fatalf("a silent discard passed the check: %v", err)
	}
	snap = g.TakeSnapshot()
	g.DiscardFromHand(0, g.Players[0].Hand[0])
	if err := g.CheckEvents(snap, g.takeEvents()); err != nil {
		t.Fatalf("DiscardFromHand: %v", err)
	}
}

// 22.2: a silent damage change fails; a shuffle is silent on purpose.
func TestCheckEventsDamageAndShuffle(t *testing.T) {
	g := eventCheckGame(t)
	top, _ := g.Monsters[0].TopOf()
	snap := g.TakeSnapshot()
	g.Object(top).Damage++
	if err := g.CheckEvents(snap, nil); err == nil || !strings.Contains(err.Error(), "damage") {
		t.Fatalf("silent damage passed the check: %v", err)
	}
	snap = g.TakeSnapshot()
	g.ShuffleDeck(LootDeck)
	if err := g.CheckEvents(snap, nil); err != nil {
		t.Fatalf("a shuffle: %v", err)
	}
}

// Every silent change says why; a gap names its 19.N item.
func TestSilentChangesHaveReasons(t *testing.T) {
	for _, s := range SilentChanges() {
		if s.Name == "" || s.Reason == "" || s.Match == nil {
			t.Errorf("silent change %q needs a name, a reason and a match", s.Name)
		}
		if s.Bug != "" && !strings.HasPrefix(s.Bug, "19.") {
			t.Errorf("silent change %q: bug %q is no 19.N item", s.Name, s.Bug)
		}
	}
}

// 19.174: an empty deck that takes back its discard has one
// moved_to_deck that counts the cards; without it the check fails.
func TestCheckEventsDeckRefill(t *testing.T) {
	g := eventCheckGame(t)
	ids := slices.Clone(g.Decks[LootDeck])
	g.Decks[LootDeck] = nil
	for _, id := range ids {
		g.discard(id, LootDeck)
	}
	snap := g.TakeSnapshot()
	g.drawTop(LootDeck)
	events := g.takeEvents()
	i := slices.IndexFunc(events, func(e Event) bool { return e.Kind == EvMovedToDeck })
	if i < 0 || events[i].Amount != len(ids) || events[i].Text != "loot deck" || events[i].Object != 0 {
		t.Fatalf("want one moved_to_deck of %d cards, no object; events %+v", len(ids), events)
	}
	if err := g.CheckEvents(snap, events); err != nil {
		t.Fatalf("a refill: %v", err)
	}
	ids = slices.Clone(g.Decks[LootDeck])
	g.Decks[LootDeck] = nil
	for _, id := range ids {
		g.discard(id, LootDeck)
	}
	snap = g.TakeSnapshot()
	g.drawTop(LootDeck)
	if err := g.CheckEvents(snap, nil); err == nil || !strings.Contains(err.Error(), "loot discard -> loot deck") {
		t.Fatalf("a silent refill passed the check: %v", err)
	}
}

// 22.3: a damage drop passes as "max HP fell" only when the max HP fell
// at least as much.
func TestCheckEventsMaxHPFell(t *testing.T) {
	g := eventCheckGame(t)
	top, _ := g.Monsters[0].TopOf()
	g.Object(top).Damage = 1
	snap := g.TakeSnapshot()
	g.Object(top).Damage = 0
	if err := g.CheckEvents(snap, nil); err == nil || !strings.Contains(err.Error(), "damage") {
		t.Fatalf("a silent heal passed as a lower max HP: %v", err)
	}
	var fell SilentChange
	for _, s := range SilentChanges() {
		if s.Name == "max HP fell" {
			fell = s
		}
	}
	inPlay := Zone{Kind: ZoneInPlay}
	c := Change{What: "damage", To: inPlay, Before: 2, After: 1, MaxBefore: 3, MaxAfter: 2}
	if !fell.Match(g, c, nil) {
		t.Error("max HP 3 -> 2 explains damage 2 -> 1")
	}
	c.After = 0
	if fell.Match(g, c, nil) {
		t.Error("max HP 3 -> 2 does not explain damage 2 -> 0")
	}
	c.After, c.MaxAfter = 1, 3
	if fell.Match(g, c, nil) {
		t.Error("the same max HP explains no damage drop")
	}
}
