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

// 22.2: a silent damage change fails. 19.178: a shuffle needs its
// shuffled event.
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
	events := g.takeEvents()
	if err := g.CheckEvents(snap, nil); err == nil || !strings.Contains(err.Error(), "shuffled") {
		t.Fatalf("a silent shuffle passed the check: %v", err)
	}
	if len(events) != 1 || events[0].Kind != EvShuffled || events[0].Text != "loot deck" {
		t.Fatalf("want one shuffled event for the loot deck, got %+v", events)
	}
	snap = g.TakeSnapshot()
	g.ShuffleDeck(LootDeck)
	if err := g.CheckEvents(snap, g.takeEvents()); err != nil {
		t.Fatalf("a shuffle with its event: %v", err)
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

// 22.10: a silent change of a player's cents fails, and so does one
// whose events move another amount; a passed_priority names no cents.
func TestCheckEventsPlayerCents(t *testing.T) {
	g := eventCheckGame(t)
	snap := g.TakeSnapshot()
	g.Players[1].Cents += 2
	passed := Event{Kind: EvPassed, Player: 1}
	if err := g.CheckEvents(snap, []Event{passed}); err == nil || !strings.Contains(err.Error(), "player 1 cents") {
		t.Fatalf("a silent cents change passed the check: %v", err)
	}
	gained := Event{Kind: EvGainedCents, Player: 1, Amount: 1}
	if err := g.CheckEvents(snap, []Event{gained}); err == nil || !strings.Contains(err.Error(), "player 1 cents") {
		t.Fatalf("a gained_cents of 1¢ explained 2¢: %v", err)
	}
	gained.Amount = 2
	if err := g.CheckEvents(snap, []Event{gained}); err != nil {
		t.Fatalf("a gained_cents of 2¢: %v", err)
	}
	snap = g.TakeSnapshot()
	g.GiveCents(1, 0, 2)
	if err := g.CheckEvents(snap, g.takeEvents()); err != nil {
		t.Fatalf("GiveCents: %v", err)
	}
}

// 22.10: a player's damage and soul value need an event of their kind.
func TestCheckEventsPlayerDamageAndSouls(t *testing.T) {
	g := eventCheckGame(t)
	snap := g.TakeSnapshot()
	g.Players[0].Damage++
	passed := Event{Kind: EvPassed, Player: 0}
	if err := g.CheckEvents(snap, []Event{passed}); err == nil || !strings.Contains(err.Error(), "player 0 player damage") {
		t.Fatalf("silent player damage passed the check: %v", err)
	}
	if err := g.CheckEvents(snap, []Event{{Kind: EvDamaged, Player: 0, Amount: 1}}); err != nil {
		t.Fatalf("damaged: %v", err)
	}
	g.Players[0].Damage = 0
	// 19.193: a soul stolen in place names the player it left.
	id := giveItem(g, 0, g.Object(g.Decks[TreasureDeck][0]).Card)
	g.Object(id).Role = RoleSoul
	snap = g.TakeSnapshot()
	g.GainControl(1, id)
	events := g.takeEvents()
	if err := g.CheckEvents(snap, events); err != nil {
		t.Fatalf("a stolen soul: %v", err)
	}
	events[0].By = 0
	if err := g.CheckEvents(snap, events); err == nil || !strings.Contains(err.Error(), "player 0 souls") {
		t.Fatalf("a soul lost with no event naming its player passed: %v", err)
	}
}

// 22.10: object 0 is a real card; an event with no object does not
// name it, so its silent move fails. An event with its card does.
func TestCheckEventsObjectZero(t *testing.T) {
	var g *Game
	for seed := uint64(1); g == nil; seed++ { // a deal that leaves object 0 in its deck
		h, _, err := NewGame(Setup{Seed: seed, Players: 2, Sets: fuzzSets})
		if err != nil {
			t.Fatal(err)
		}
		if h.Object(0).Zone.Kind == ZoneDeck {
			g = h
		}
	}
	g.takeEvents()
	card, d := g.Object(0).Card, g.Object(0).Zone.Deck
	discard := func() ObjectID {
		g.Decks[d] = remove(g.Decks[d], 0)
		return g.discard(0, d)
	}
	h, err := g.Clone()
	if err != nil {
		t.Fatal(err)
	}
	snap := g.TakeSnapshot()
	discard()
	passed := Event{Kind: EvPassed, Player: 0} // Object unset: 0
	if err := g.CheckEvents(snap, []Event{passed}); err == nil || !strings.Contains(err.Error(), "moved from 0") {
		t.Fatalf("a silent move of object 0 passed the check: %v", err)
	}
	g = h
	snap = g.TakeSnapshot()
	nid := discard()
	named := Event{Kind: EvDiscarded, Player: NoPlayer, Object: nid, Card: card}
	if err := g.CheckEvents(snap, []Event{passed, named}); err != nil {
		t.Fatalf("a discard of object 0 with its event: %v", err)
	}
}

// 22.10: a max HP that falls and rises again in one step still explains
// the damage its fall removed: checkMaxHP logs it.
func TestCheckEventsMaxHPFellAndRose(t *testing.T) {
	g := eventCheckGame(t)
	g.Players[0].Damage = 2
	snap := g.TakeSnapshot()
	g.Players[0].Damage = 1
	if err := g.CheckEvents(snap, nil); err == nil || !strings.Contains(err.Error(), "player 0 player damage") {
		t.Fatalf("a silent heal passed: %v", err)
	}
	snap = g.TakeSnapshot()
	g.logHPFell(hpKey{p: 0}, 1, 0)
	g.Players[0].Damage = 0
	if err := g.CheckEvents(snap, nil); err != nil {
		t.Fatalf("damage a lower max HP removed: %v", err)
	}
}
