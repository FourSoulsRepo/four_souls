package rulesengine

import (
	"fmt"
	"slices"
	"strings"
)

// ExpectedMonster checks one card in a slot (24.5). Card names it, as
// a setup ref does: the first such card in slot order (monster, final
// boss, minion, room slots). With Slot (from 1), Card must be on top of
// that monster slot. Dead true: no slot holds the card, and it is in a
// discard or a player's souls.
type ExpectedMonster struct {
	Slot    int     `json:"slot,omitempty"`
	Card    CardRef `json:"card"`
	Damage  *int    `json:"damage,omitempty"`
	HP      *int    `json:"hp,omitempty"` // what is left: HP minus damage
	ATK     *int    `json:"atk,omitempty"`
	Evasion *int    `json:"evasion,omitempty"`
	Dead    *bool   `json:"dead,omitempty"`
}

// slotCards are the cards on top of the monster, final boss, minion and
// room slots, in that order, then the active bonus souls; slotted is
// how many of them are in slots.
func (g *Game) slotCards() (refs []CardRef, ids []ObjectID, slotted int) {
	for _, rows := range append(g.MonsterRows(), g.Rooms) {
		for _, slot := range rows {
			if id, ok := slot.TopOf(); ok {
				refs, ids = append(refs, g.Object(id).Card), append(ids, id)
			}
		}
	}
	slotted = len(refs)
	for _, id := range g.BonusSouls {
		refs, ids = append(refs, g.Object(id).Card), append(ids, id)
	}
	return refs, ids, slotted
}

// checks gathers mismatches; each names what it checked.
type checks []string

func (c *checks) add(format string, args ...any) { *c = append(*c, fmt.Sprintf(format, args...)) }

func (c *checks) int(name string, want *int, got int) {
	if want != nil && *want != got {
		c.add("%s is %d, expected %d", name, got, *want)
	}
}

func (c *checks) bool(name string, want *bool, got bool) {
	if want != nil && *want != got {
		c.add("%s is %v, expected %v", name, got, *want)
	}
}

// cards compares a list of cards with an expected one in any order;
// nil expects nothing (not checked), an empty list no cards.
func (c *checks) cards(name string, want, got []CardRef) {
	if want == nil {
		return
	}
	a, b := slices.Clone(want), slices.Clone(got)
	slices.Sort(a)
	slices.Sort(b)
	if !slices.Equal(a, b) {
		c.add("%s are %v, expected %v", name, refList(got), refList(want))
	}
}

// counters compares counters on cards: refs and ids are the cards a
// ref may name, in order; who is "player 0 " or "".
func (c *checks) counters(who string, want []SituationCounter, refs []CardRef, ids []ObjectID, g *Game) {
	for _, w := range want {
		i := slices.Index(refs, w.Card)
		if i < 0 {
			c.add("%scounters: %q is not there, options %v", who, w.Card, refList(refs))
			continue
		}
		label := "counters"
		if w.Name != "" {
			label = w.Name + " counters"
		}
		if got := g.Object(ids[i]).CountersOf(w.Name); got != w.Count {
			c.add("%s%s on %s are %d, expected %d", who, label, w.Card, got, w.Count)
		}
	}
}

// refList prints refs as a list, [] when empty.
func refList(refs []CardRef) string {
	s := make([]string, len(refs))
	for i, r := range refs {
		s[i] = string(r)
	}
	return "[" + strings.Join(s, " ") + "]"
}

// checkPlayer compares one player with x.
func (g *Game) checkPlayer(c *checks, x ExpectedPlayer) {
	if x.Player < 0 || x.Player >= len(g.Players) {
		c.add("no player %d", x.Player)
		return
	}
	pl := g.Players[x.Player]
	who := fmt.Sprintf("player %d ", x.Player)
	c.int(who+"cents", x.Cents, pl.Cents)
	c.int(who+"hand", x.Hand, len(pl.Hand))
	c.int(who+"damage", x.Damage, pl.Damage)
	c.int(who+"souls", x.Souls, g.SoulValue(pl.ID))
	c.bool(who+"dead", x.Dead, pl.Dead)
	c.int(who+"hp", x.HP, g.PlayerHP(pl.ID))
	c.int(who+"atk", x.ATK, g.PlayerATK(pl.ID))
	var items, souls, hand, off []CardRef
	refs, ids := []CardRef{"character"}, []ObjectID{pl.Character}
	if !g.Object(pl.Character).Charged {
		off = append(off, "character")
	}
	for _, id := range pl.InPlay {
		o := g.Object(id)
		refs, ids = append(refs, o.Card), append(ids, id)
		switch o.Role { //nolint:exhaustive // only items and souls are listed
		case RoleItem:
			items = append(items, o.Card)
			if !o.Charged {
				off = append(off, o.Card)
			}
		case RoleSoul:
			souls = append(souls, o.Card)
		}
	}
	for _, id := range pl.Hand {
		hand = append(hand, g.Object(id).Card)
	}
	c.cards(who+"items", x.Items, items)
	c.cards(who+"hand_cards", x.HandCards, hand)
	c.cards(who+"soul_cards", x.SoulCards, souls)
	c.cards(who+"deactivated", x.Deactivated, off)
	c.counters(who, x.Counters, refs, ids, g)
}

// checkMonster compares one card in a slot with x.
func (g *Game) checkMonster(c *checks, x ExpectedMonster) {
	if x.Card == "" {
		c.add("monsters: a check needs card")
		return
	}
	refs, ids, slotted := g.slotCards()
	refs, ids = refs[:slotted], ids[:slotted]
	name := fmt.Sprintf("monster %s", x.Card)
	id := ObjectID(0)
	found := false
	if x.Slot != 0 {
		if x.Slot < 1 || x.Slot > len(g.Monsters) {
			c.add("no monster slot %d, slots 1 to %d", x.Slot, len(g.Monsters))
			return
		}
		name = fmt.Sprintf("monster slot %d", x.Slot)
		top, ok := g.Monsters[x.Slot-1].TopOf()
		got := CardRef("")
		if ok {
			got = g.Object(top).Card
		}
		if got != x.Card {
			c.add("%s is %q, expected %q", name, got, x.Card)
			return
		}
		id, found = top, true
	} else if i := slices.Index(refs, x.Card); i >= 0 {
		id, found = ids[i], true
	}
	if x.Dead != nil {
		dead := !found && g.leftPlay(x.Card)
		c.bool(name+" dead", x.Dead, dead)
	}
	if x.Damage == nil && x.HP == nil && x.ATK == nil && x.Evasion == nil {
		return
	}
	if !found {
		c.add("%s is in no slot, options %v", name, refList(refs))
		return
	}
	c.int(name+" damage", x.Damage, g.Object(id).Damage)
	c.int(name+" hp", x.HP, g.HP(id))
	c.int(name+" atk", x.ATK, g.MonsterATK(id))
	c.int(name+" evasion", x.Evasion, g.Evasion(id))
}

// leftPlay: a card of this ref is in a discard or a player's souls.
func (g *Game) leftPlay(card CardRef) bool {
	for d := range deckCount {
		for _, id := range g.Discards[d] {
			if g.Object(id).Card == card {
				return true
			}
		}
	}
	for _, pl := range g.Players {
		for _, id := range pl.InPlay {
			if o := g.Object(id); o.Card == card && o.Role == RoleSoul {
				return true
			}
		}
	}
	return false
}

// checkDecks compares deck tops (top first) and whole discards, any
// order, by deck name.
func (g *Game) checkDecks(c *checks, decks, discards map[string][]CardRef) {
	for _, m := range []map[string][]CardRef{decks, discards} {
		for name := range m {
			if !slices.Contains(deckNames[:], name) {
				c.add("%q is no deck, options %v", name, deckNames)
			}
		}
	}
	for d := range deckCount {
		if want, ok := decks[d.String()]; ok {
			var got []CardRef
			for _, id := range g.DeckTop(d, len(want)) {
				got = append(got, g.Object(id).Card)
			}
			if !slices.Equal(got, want) {
				c.add("%s deck top is %v, expected %v", d, refList(got), refList(want))
			}
		}
		if want, ok := discards[d.String()]; ok {
			var got []CardRef
			for _, id := range g.Discards[d] {
				got = append(got, g.Object(id).Card)
			}
			c.cards(d.String()+" discards", want, got)
		}
	}
}

// checkBottoms compares the cards at the bottom of decks, as many as
// want names, in any order: several cards go to the bottom in random
// order (Q-22).
func (g *Game) checkBottoms(c *checks, bottoms map[string][]CardRef) {
	for name := range bottoms {
		if !slices.Contains(deckNames[:], name) {
			c.add("%q is no deck, options %v", name, deckNames)
		}
	}
	for d := range deckCount {
		want, ok := bottoms[d.String()]
		if !ok {
			continue
		}
		deck := g.Decks[d]
		got := []CardRef{}
		for _, id := range deck[:min(len(want), len(deck))] {
			got = append(got, g.Object(id).Card)
		}
		slices.Reverse(got) // bottom last, as the top is top first
		a, b := slices.Clone(want), slices.Clone(got)
		slices.Sort(a)
		slices.Sort(b)
		if !slices.Equal(a, b) {
			c.add("%s deck bottom is %v, expected %v in any order", d, refList(got), refList(want))
		}
	}
}

// checkRolls compares the values the rolls resolved as (roll_resolved
// events), in order; nil is not checked, an empty list expects none.
func checkRolls(c *checks, want []int, events []Event) {
	if want == nil {
		return
	}
	got := []int{}
	for _, e := range events {
		if e.Kind == EvRollResolved {
			got = append(got, e.Amount)
		}
	}
	if !slices.Equal(got, want) {
		c.add("resolved_rolls are %v, expected %v", got, want)
	}
}

// checkGame compares the turn and the game's end.
func (g *Game) checkGame(c *checks, e Expectation) {
	c.int("active", e.Active, int(g.Turn.Active))
	c.int("loot_plays", e.LootPlays, g.Turn.LootPlays)
	c.int("attacks", e.Attacks, g.Turn.Attacks)
	c.bool("over", e.Over, g.Over)
	if e.Winners != nil {
		got := make([]int, len(g.Winners))
		for i, p := range g.Winners {
			got[i] = int(p)
		}
		a, b := slices.Clone(e.Winners), slices.Clone(got)
		slices.Sort(a)
		slices.Sort(b)
		if !slices.Equal(a, b) {
			c.add("winners are %v, expected %v", got, e.Winners)
		}
	}
	refs, ids, _ := g.slotCards()
	c.counters("", e.Counters, refs, ids, g)
}

// ExpectedPrompt checks the open question once a step settled (24.6):
// Player is who is asked, Options every option by plain label
// (Game.Plain), in any order. A field left out is not checked.
type ExpectedPrompt struct {
	Player  *int     `json:"player,omitempty"`
	Options []string `json:"options,omitempty"`
}

// ExpectedEvent is an event kind that must happen (24.6), for Player
// when set. None: no such event may happen.
type ExpectedEvent struct {
	Kind   EventKind `json:"kind"`
	Player *int      `json:"player,omitempty"`
	None   bool      `json:"none,omitempty"`
}

// asks reports whether e checks a prompt, so settling stops at it.
func (e *Expectation) asks() bool { return e != nil && e.Prompt != nil }

// checkPrompt compares the open question with x.
func (g *Game) checkPrompt(c *checks, x *ExpectedPrompt) {
	if x == nil {
		return
	}
	w := g.Prompt()
	if w.Kind != PromptChoose && w.Kind != PromptOrder {
		c.add("prompt: no question is open, the game waits for %s", waitingFor(w.Kind))
		return
	}
	c.int("prompt player", x.Player, int(w.Player))
	if x.Options == nil {
		return
	}
	got := g.PlainOptions()
	a, b := slices.Clone(x.Options), slices.Clone(got)
	slices.Sort(a)
	slices.Sort(b)
	if !slices.Equal(a, b) {
		c.add("prompt options are %q, expected %q", got, x.Options)
	}
}

// stackItemNames name stack items that are not a card's loot, ability
// or trigger in a resolved check.
var stackItemNames = map[StackKind]string{
	StackRoll: "dice_roll", StackDamage: "damage", StackDeath: "death",
	StackPenalty: "death_penalty", StackDeathStep: "death_step",
}

// resolvedName names a resolved item: its card for a loot card, an
// ability or a trigger; else its kind, attack_roll for an attack roll.
func resolvedName(e Event) string {
	if e.Item == nil {
		return string(e.Card)
	}
	switch {
	case e.Item.Kind == StackRoll && e.Item.Attack:
		return "attack_roll"
	case e.Item.Kind == StackLoot || e.Item.Kind == StackAbility || e.Item.Kind == StackTrigger:
		return string(e.Card)
	}
	return stackItemNames[e.Item.Kind]
}

// checkResolved compares the items that resolved with want, in order:
// each name must resolve after the one before it; other items may
// resolve between them. nil is not checked; an empty list expects none.
func (g *Game) checkResolved(c *checks, want []string, events []Event) {
	if want == nil {
		return
	}
	var got []string
	for _, e := range events {
		if e.Kind == EvStackResolved {
			got = append(got, resolvedName(e))
		}
	}
	kinds := []string{"attack_roll"}
	for _, n := range stackItemNames {
		kinds = append(kinds, n)
	}
	slices.Sort(kinds)
	for _, n := range want {
		if _, ok := g.cards.find(CardRef(n)); !ok && !slices.Contains(kinds, n) {
			c.add("resolved: %q is no card and no stack item, items %v", n, kinds)
			return
		}
	}
	if !inOrder(want, got) || len(want) == 0 && len(got) > 0 {
		c.add("resolved %v, expected %v in this order", got, want)
	}
}

// inOrder reports whether want is a subsequence of got: each item
// after the one before it.
func inOrder(want, got []string) bool {
	at := 0
	for _, n := range want {
		for at < len(got) && got[at] != n {
			at++
		}
		if at == len(got) {
			return false
		}
		at++
	}
	return true
}

// eventKinds are every event kind, sorted, for errors.
func eventKinds() []string {
	out := make([]string, 0, len(barrierKinds))
	for k := range barrierKinds {
		out = append(out, string(k))
	}
	slices.Sort(out)
	return out
}

// checkEvents checks that each expected kind happened (for its player,
// when set), or with None that it did not.
func checkEvents(c *checks, want []ExpectedEvent, events []Event) {
	for _, x := range want {
		if _, ok := barrierKinds[x.Kind]; !ok {
			c.add("events: %q is no event kind, options %v", x.Kind, eventKinds())
			continue
		}
		who := ""
		if x.Player != nil {
			who = fmt.Sprintf(" for player %d", *x.Player)
		}
		n := 0
		for _, e := range events {
			if e.Kind == x.Kind && (x.Player == nil || int(e.Player) == *x.Player) {
				n++
			}
		}
		switch {
		case x.None && n > 0:
			c.add("events: %d %q%s, expected none", n, x.Kind, who)
		case !x.None && n == 0:
			c.add("events: no %q%s, expected one", x.Kind, who)
		}
	}
}
