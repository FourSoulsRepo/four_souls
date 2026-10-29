package rulesengine

import (
	"reflect"
	"slices"
	"testing"
)

// reasonOf is the reason Blocked gives for kind and id, or "".
func reasonOf(bs []Block, kind IntentKind, id ObjectID) string {
	for _, b := range bs {
		if b.Kind == kind && b.Object == id {
			return b.Reason.String()
		}
	}
	return ""
}

// Blocked says why a card or turn action can't be used now (15.29).
func TestBlockedReasons(t *testing.T) {
	g := newAbilityGame(t)
	a := g.Turn.Active
	b := g.next(a)
	loot := g.Players[a].Hand[0]
	bag := giveItem(g, a, "coin_bag")
	theirs := giveItem(g, b, "coin_bag")
	if bs := g.Blocked(a); reasonOf(bs, IntentPlayLoot, loot) != "" || reasonOf(bs, IntentAttack, 0) != "" {
		t.Fatalf("a can loot and attack; blocked %+v", bs)
	}
	// The other player: not their turn, no loot play, and no priority.
	bs := g.Blocked(b)
	for _, k := range []IntentKind{IntentAttack, IntentPurchase, IntentEndTurn} {
		if got := reasonOf(bs, k, 0); got != "refuse.not_your_turn" {
			t.Errorf("b's intent %d: reason %q", k, got)
		}
	}
	if got := reasonOf(bs, IntentPlayLoot, g.Players[b].Hand[0]); got != "refuse.loot_your_turn" {
		t.Errorf("b's loot: reason %q", got)
	}
	if got := reasonOf(bs, IntentActivate, theirs); got != "refuse.wait_priority" {
		t.Errorf("b's item: reason %q", got)
	}
	// Used up: the loot play, the attacks, the item's ↷.
	g.Turn.LootPlays, g.Turn.Attacks = 0, 0
	g.Object(bag).Charged = false
	waiting := g.Waiting
	bs = g.Blocked(a)
	if !reflect.DeepEqual(g.Waiting, waiting) {
		t.Errorf("Blocked changed the prompt: %+v, was %+v", g.Waiting, waiting)
	}
	want := map[string]string{
		"loot":   "refuse.no_loot_plays",
		"attack": "refuse.no_attacks",
		"item":   "refuse.cost.used",
	}
	got := map[string]string{
		"loot":   reasonOf(bs, IntentPlayLoot, loot),
		"attack": reasonOf(bs, IntentAttack, 0),
		"item":   reasonOf(bs, IntentActivate, bag),
	}
	for k, w := range want {
		if got[k] != w {
			t.Errorf("%s: reason %q, want %q", k, got[k], w)
		}
	}
	// With a roll on the stack, the turn actions wait for it.
	g.push(StackItem{Kind: StackRoll, Roll: 3, Controller: a})
	if got := reasonOf(g.Blocked(a), IntentPurchase, 0); got != "refuse.wait_stack" {
		t.Errorf("purchase with a stack: reason %q", got)
	}
}

// checkBlocked: what Blocked lists is never allowed, every hand card is
// one or the other, and each block has words.
func checkBlocked(t *testing.T, g *Game, p PlayerID) {
	t.Helper()
	allowed := g.Allowed(p)
	bs := g.Blocked(p)
	for _, b := range bs {
		if b.Reason.IsZero() {
			t.Fatalf("player %d: block %+v has no reason", p, b)
		}
		if slices.ContainsFunc(allowed, func(in Intent) bool {
			return in.Kind == b.Kind && (len(in.Objects) == 0 || in.Objects[0] == b.Object)
		}) {
			t.Fatalf("player %d: %+v is both allowed and blocked", p, b)
		}
	}
	for _, id := range g.Players[p].Hand {
		played := slices.ContainsFunc(allowed, func(in Intent) bool {
			return in.Kind == IntentPlayLoot && in.Objects[0] == id
		})
		if played == (reasonOf(bs, IntentPlayLoot, id) != "") {
			t.Fatalf("player %d: hand card %d must be allowed or blocked, not both or neither", p, id)
		}
	}
}

// BlockedTurn lists only the turn actions, and BlockedCard only the card
// asked about: the same reasons Blocked gives, split for the wire.
func TestBlockedSplit(t *testing.T) {
	g := newAbilityGame(t)
	b := g.next(g.Turn.Active)
	theirs := giveItem(g, b, "coin_bag")
	for _, bl := range g.BlockedTurn(b) {
		if bl.Object != 0 {
			t.Errorf("turn list has card %d", bl.Object)
		}
	}
	if got := reasonOf(g.BlockedCard(b, theirs), IntentActivate, theirs); got == "" {
		t.Errorf("the card asked about has no reason")
	}
	for _, bl := range g.BlockedCard(b, theirs) {
		if bl.Object != theirs {
			t.Errorf("card list has %d, want only %d", bl.Object, theirs)
		}
	}
	if got := reasonOf(g.BlockedCard(b, theirs), IntentAttack, 0); got != "" {
		t.Errorf("a turn action in the card list: %q", got)
	}
}
