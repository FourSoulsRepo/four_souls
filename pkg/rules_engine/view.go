package rulesengine

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

// ViewerKind is who looks at the game (A-07, SP-01, SP-02).
type ViewerKind int

// The viewer kinds.
const (
	ViewPlayer    ViewerKind = iota // a seated player: sees their own hand
	ViewSpectator                   // public information only
	ViewJudge                       // sees every hand, never deck order or the RNG
)

// Viewer is one person looking at the game.
type Viewer struct {
	Kind   ViewerKind
	Player PlayerID // for ViewPlayer
}

// CardView is a visible card.
type CardView struct {
	ID      ObjectID `json:"id"`
	Card    CardRef  `json:"card"`
	Role    Role     `json:"role,omitempty"`
	Charged bool     `json:"charged"`
	Damage  int      `json:"damage,omitempty"`
	// Counters are the counters on the card, e.g. charge counters.
	Counters []Counter `json:"counters,omitempty"`
	// Stats are a monster's live and printed stats (C-03); nil for
	// other cards.
	Stats *MonsterStats `json:"stats,omitempty"`
	// Abilities are the words of each ability an activate intent may
	// pick, by its Choice; only on a played card with two or more
	// activated abilities, so a client can offer the lines (C-05).
	Abilities []string `json:"abilities,omitempty"`
	// Souls is the soul value the card adds to its owner, on a soul in
	// play or a player's character, so a client can tie crystals to it
	// (V-43).
	Souls int `json:"souls,omitempty"`
}

// MonsterStats are a monster's stats now and as printed, so a client can
// mark the ones that changed (C-03).
type MonsterStats struct {
	HP         int `json:"hp"`
	DC         int `json:"dc"`
	ATK        int `json:"atk"`
	PrintedHP  int `json:"printed_hp"`
	PrintedDC  int `json:"printed_dc"`
	PrintedATK int `json:"printed_atk"`
}

// PlayerView is what a viewer sees of one player.
type PlayerView struct {
	ID        PlayerID   `json:"id"`
	Character CardView   `json:"character"`
	HP        int        `json:"hp"`
	ATK       int        `json:"atk"`
	Cents     int        `json:"cents"`
	Souls     int        `json:"souls"`
	Dead      bool       `json:"dead,omitempty"`
	InPlay    []CardView `json:"in_play"`
	HandSize  int        `json:"hand_size"`
	Hand      []CardView `json:"hand,omitempty"` // only when the viewer may see it
	// LootPlays, Attacks and Purchases are what the player may still do
	// now, bonuses included; attacks and purchases only on their turn.
	LootPlays int `json:"loot_plays"`
	Attacks   int `json:"attacks"`
	Purchases int `json:"purchases"`
	// TreasureTop is the top treasure card, when this player may look at
	// it now (There's Options) and the viewer may see it.
	TreasureTop *CardView `json:"treasure_top,omitempty"`
}

// SlotView is a table slot: the top card and the covered ones below.
type SlotView struct {
	Top     *CardView  `json:"top,omitempty"`
	Covered []CardView `json:"covered,omitempty"`
}

// GameView is everything one viewer may see.
type GameView struct {
	Players    []PlayerView          `json:"players"`
	Shop       []SlotView            `json:"shop"`
	Monsters   []SlotView            `json:"monsters"`
	DeckSizes  [deckCount]int        `json:"deck_sizes"`
	Discards   [deckCount][]CardView `json:"discards"`
	BonusSouls []CardView            `json:"bonus_souls,omitempty"`
	Stack      []StackItem           `json:"stack"`
	Turn       Turn                  `json:"turn"`
	Waiting    Prompt                `json:"waiting"`
	Over       bool                  `json:"over,omitempty"`
	Winners    []PlayerID            `json:"winners,omitempty"`
	// WinSouls is the soul value that wins the game (R-WIN-01), so a
	// client can show how far each player still has to go (V-30).
	WinSouls int `json:"win_souls"`
	// Revealed is the monster-deck card a deck attack drew, while the
	// attacker picks its slot (R-ATK-08); it is public (19.9).
	Revealed *CardView `json:"revealed,omitempty"`
	// Subject is the card the open question is about, e.g. the deck top
	// a player looks at; only for the chooser (19.11).
	Subject *CardView `json:"subject,omitempty"`
	// Ordering are the cards of an open order prompt, in their order
	// now; only for the chooser (19.16).
	Ordering []CardView `json:"ordering,omitempty"`
	// Options says what each option of a choose prompt points at, one
	// per Waiting.Options; only for whoever sees those options.
	Options []OptionView `json:"options,omitempty"`
	// Activating is the card whose mode or targets the chooser picks
	// now, so a client can aim from it (15.34, T-13); only for them.
	Activating *ObjectID `json:"activating,omitempty"`
	// Passed lists who passed in a row in the open priority window,
	// first passer first; Responding lists who may still respond, the
	// waiting player first (V-11). Both are empty outside a window.
	Passed     []PlayerID `json:"passed,omitempty"`
	Responding []PlayerID `json:"responding,omitempty"`
}

// OptionView is what one option of a choose prompt points at, so a
// client can let the player pick it on the table. An option with none
// of these set is plain text.
type OptionView struct {
	Object   ObjectID `json:"object,omitempty"` // a card on the table or in the chooser's hand
	Player   PlayerID `json:"player"`           // a player, when IsPlayer
	IsPlayer bool     `json:"is_player,omitempty"`
	Deck     bool     `json:"deck,omitempty"`     // the monster or treasure deck
	StackID  int      `json:"stack_id,omitempty"` // a stack item
	// Outcome is what the option leads to, when the engine can tell
	// (replacement order, 19.17).
	Outcome *OutcomeView `json:"outcome,omitempty"`
}

// OutcomeView is what an order of replacements gives (19.17).
// Actions are what then happens; Unused are the cards whose effect no
// longer applies; Again means the player orders the rest next.
type OutcomeView struct {
	Actions []OutcomeAction `json:"actions,omitempty"`
	Unused  []CardRef       `json:"unused,omitempty"`
	Again   bool            `json:"again,omitempty"`
}

// OutcomeAction is one action of an outcome: its kind, player, amount,
// and the card it acts on, if any (Bum-Bo! levels up).
type OutcomeAction struct {
	Kind   ActionKind `json:"kind"`
	Player PlayerID   `json:"player"`
	Amount int        `json:"amount"`
	Card   CardRef    `json:"card,omitempty"`
}

// sees reports whether the viewer may see player p's hand.
func (v Viewer) sees(p PlayerID) bool {
	return v.Kind == ViewJudge || (v.Kind == ViewPlayer && v.Player == p)
}

// View builds what a viewer sees. It is the only way state reaches
// clients (A-07): hidden hands, deck contents and order, and the RNG
// never leave the engine.
func (g *Game) View(v Viewer) GameView {
	out := GameView{Turn: g.Turn, Over: g.Over, Winners: g.Winners, WinSouls: g.WinSouls, Stack: append([]StackItem(nil), g.Stack...)}
	for i := range out.Stack {
		out.Stack[i] = g.shown(out.Stack[i])
	}
	for _, pl := range g.Players {
		pv := PlayerView{
			ID: pl.ID, Character: g.playedView(pl.Character), HP: g.PlayerHP(pl.ID), ATK: g.PlayerATK(pl.ID),
			Cents: pl.Cents, Souls: g.SoulValue(pl.ID), Dead: pl.Dead, HandSize: len(pl.Hand),
			LootPlays: g.lootPlaysFor(pl.ID),
		}
		if pl.ID == g.Turn.Active {
			pv.Attacks, pv.Purchases = g.attacksLeft(), g.purchasesLeft()
		}
		pv.Character.Souls = g.characterSouls(pl.Character)
		for _, id := range pl.InPlay {
			cv := g.playedView(id)
			cv.Souls = g.soulsOf(id)
			pv.InPlay = append(pv.InPlay, cv)
		}
		if v.sees(pl.ID) || g.bonus(StatRevealHand, pl.ID, 0) > 0 {
			for _, id := range pl.Hand {
				pv.Hand = append(pv.Hand, g.cardView(id))
			}
		}
		if v.sees(pl.ID) {
			if top := g.DeckTop(TreasureDeck, 1); len(top) > 0 && g.peeksTreasure(pl.ID) {
				cv := g.cardView(top[0])
				pv.TreasureTop = &cv
			}
		}
		out.Players = append(out.Players, pv)
	}
	out.Shop = g.slotViews(g.Shop)
	out.Monsters = g.slotViews(g.Monsters)
	for d := range deckCount {
		out.DeckSizes[d] = len(g.Decks[d])
		for _, id := range g.Discards[d] {
			out.Discards[d] = append(out.Discards[d], g.cardView(id))
		}
	}
	for _, id := range g.BonusSouls {
		out.BonusSouls = append(out.BonusSouls, g.cardView(id))
	}
	if g.Attack.Revealed != 0 {
		cv := g.cardView(g.Attack.Revealed)
		out.Revealed = &cv
	}
	out.Waiting = g.Waiting
	if g.Waiting.Kind == PromptChoose && !v.sees(g.Waiting.Player) {
		// Options may name hidden cards (e.g. which loot card to discard).
		out.Waiting.Options = nil
	}
	if g.Waiting.Kind == PromptChoose && v.sees(g.Waiting.Player) {
		out.Options = g.optionViews()
		out.Activating = g.activatingSource()
		if g.Choice != nil && g.Choice.Subject != 0 {
			cv := g.cardView(g.Choice.Subject)
			out.Subject = &cv
		}
	}
	if g.Waiting.Kind == PromptOrder && v.sees(g.Waiting.Player) && g.Choice != nil {
		for _, id := range g.Choice.Objects {
			out.Ordering = append(out.Ordering, g.cardView(id))
		}
	}
	out.Passed, out.Responding = g.responders()
	return out
}

// responders splits the players of the open priority window into those
// who passed in a row and those still to answer, in turn order.
// Priority goes round once; the window closes when all passed
// (R-PRIO-05).
func (g *Game) responders() (passed, responding []PlayerID) {
	if g.Waiting.Kind != PromptPriority || !g.Priority.Open {
		return nil, nil
	}
	n := len(g.Players)
	p := g.Priority.Holder
	for range g.Priority.Passes {
		p = PlayerID((int(p) + n - 1) % n)
		passed = append([]PlayerID{p}, passed...)
	}
	p = g.Priority.Holder
	for range n - g.Priority.Passes {
		responding = append(responding, p)
		p = g.next(p)
	}
	return passed, responding
}

// activatingSource is the card whose mode or targets the open choice
// asks for, if a play or use asks it.
func (g *Game) activatingSource() *ObjectID {
	a := g.Activating
	if a == nil || g.Choice == nil || a.Player != g.Choice.Player {
		return nil
	}
	if p := g.Choice.Purpose; p != ChooseMode && p != ChooseTarget {
		return nil
	}
	src := a.Source
	return &src
}

// optionViews maps the open choice's options to what they point at.
func (g *Game) optionViews() []OptionView {
	c := g.Choice
	if c == nil {
		return nil
	}
	out := make([]OptionView, len(g.Waiting.Options))
	for i := range out {
		out[i].Player = NoPlayer
	}
	set := func(i int, f func(o *OptionView)) {
		if i < len(out) {
			f(&out[i])
		}
	}
	switch c.Purpose {
	case ChooseAttackTarget, ChoosePurchase, ChoosePenaltyItem, ChoosePenaltyLoot, ChooseStartingItem, ChooseLowest, ChooseRecharge, ChooseRedirect:
		for i, id := range c.Objects {
			set(i, func(o *OptionView) {
				o.Object = id
				if p, ok := g.characterOf(id); ok && c.Purpose == ChooseAttackTarget {
					o.Player, o.IsPlayer = p, true // attacking a player (R-ATK-20)
				}
			})
		}
		if c.Deck {
			set(len(c.Objects), func(o *OptionView) { o.Deck = true })
		}
	case ChooseTarget:
		for i, t := range c.Targets {
			set(i, func(o *OptionView) { *o = targetView(t) })
		}
	case ChooseMonsterSlot:
		for i, s := range c.Slots {
			if top, ok := g.Monsters[s].TopOf(); ok {
				set(i, func(o *OptionView) { o.Object = top })
			}
		}
	case ChooseCursed, ChoosePassCurse:
		for i, p := range c.Slots {
			set(i, func(o *OptionView) { o.Player, o.IsPlayer = PlayerID(p), true })
		}
	case ChooseReplacement:
		for i, r := range c.Replacements {
			set(i, func(o *OptionView) { o.Object, o.Outcome = r.Object, g.outcomeView(c, i) })
		}
	case ChooseTriggerOrder, ChooseMode, ChooseAnswer, ChooseCopyAbility, ChooseDie, ChooseOrder:
		// Text options.
	}
	return out
}

// outcomeView is the preview of option i of a replacement choice.
func (g *Game) outcomeView(c *Choice, i int) *OutcomeView {
	if i >= len(c.Previews) {
		return nil
	}
	p := c.Previews[i]
	out := &OutcomeView{Again: p.Again}
	for _, a := range p.Actions {
		oa := OutcomeAction{Kind: a.Kind, Player: a.Player, Amount: a.Amount}
		if a.Object != 0 {
			oa.Card = g.CardOf(a.Object)
		}
		out.Actions = append(out.Actions, oa)
	}
	for _, r := range p.Unused {
		out.Unused = append(out.Unused, g.CardOf(r.Object))
	}
	return out
}

// targetView is what a chosen target points at on the table.
func targetView(t Chosen) OptionView {
	o := OptionView{Player: NoPlayer}
	switch t.Kind {
	case TargetPlayer, TargetOtherPlayer:
		o.Player, o.IsPlayer = t.Player, true
	case TargetDiceRoll, TargetStackAbility:
		o.StackID = t.StackID
	case TargetMonster, TargetMonsterOrPlayer, TargetItem, TargetCurse, TargetYourItem,
		TargetYourHandCard, TargetItemOrSoul:
		o.Object = t.Object
	}
	return o
}

func (g *Game) cardView(id ObjectID) CardView {
	o := g.Object(id)
	cv := CardView{
		ID: id, Card: o.Card, Role: o.Role, Charged: o.Charged, Damage: o.Damage,
		Counters: append([]Counter(nil), o.Counters...),
	}
	if o.Role == RoleMonster {
		d := g.def(id)
		cv.Stats = &MonsterStats{
			HP: g.HP(id), DC: g.Evasion(id), ATK: g.MonsterATK(id),
			// The printed dice clamp as Evasion does, so only a bonus
			// marks it as changed.
			PrintedHP: d.HP, PrintedDC: min(max(d.DC, 1), 6), PrintedATK: d.ATK,
		}
	}
	return cv
}

// playedView is the view of a character or an item in play: a card
// view with the words of its abilities, when it has two or more to
// pick from (C-05).
func (g *Game) playedView(id ObjectID) CardView {
	cv := g.cardView(id)
	refs := g.AbilitiesOf(id)
	texts := make([]string, len(refs))
	activated := 0
	for i, ref := range refs {
		d, ok := g.cards.find(ref.Card)
		if !ok || ref.Index < 0 || ref.Index >= len(d.Abilities) {
			continue
		}
		if a := d.Abilities[ref.Index]; a.Kind == Activated {
			texts[i] = a.Text
			activated++
		}
	}
	if activated > 1 {
		cv.Abilities = texts
	}
	return cv
}

func (g *Game) slotViews(slots []Slot) []SlotView {
	out := make([]SlotView, len(slots))
	for i, s := range slots {
		for j, id := range s.Cards {
			cv := g.cardView(id)
			if j == len(s.Cards)-1 {
				out[i].Top = &cv
			} else {
				out[i].Covered = append(out[i].Covered, cv) // covered cards are public (R-ZONE-12)
			}
		}
	}
	return out
}

// FilterEvents returns the events a viewer may see; a private event (such
// as a looted card) shows its card only to its player and judges.
func FilterEvents(events []Event, v Viewer) []Event {
	out := make([]Event, 0, len(events))
	for _, e := range events {
		if e.Private && !v.sees(e.Player) && !slices.ContainsFunc(e.SeenBy, v.sees) {
			e.Card, e.Object = "", 0
		}
		out = append(out, e)
	}
	return out
}

// peeksTreasure: it is p's turn and p controls a card that lets them
// look at the top treasure card.
func (g *Game) peeksTreasure(p PlayerID) bool {
	if g.Turn.Active != p {
		return false
	}
	for _, id := range g.Players[p].InPlay {
		if g.Object(id).Role == RoleItem && g.def(id).PeeksTreasure {
			return true
		}
	}
	return false
}

// shown is a stack item with its display fields filled: the used
// ability's words and the cards of its source and target (S-07, 15.22).
func (g *Game) shown(it StackItem) StackItem {
	it.Text = g.usedText(it)
	it.SourceCard = g.publicCard(it.Source)
	if !it.Target.IsPlayer {
		it.TargetCard = g.publicCard(it.Target.Object)
	}
	return it
}

// publicCard is the card of an object anyone may see, or "".
func (g *Game) publicCard(id ObjectID) CardRef {
	if id <= 0 || int(id) >= len(g.Objects) {
		return ""
	}
	o := g.Object(id)
	if o.Zone.Kind == ZoneDeck || o.Zone.Kind == ZoneHand {
		return ""
	}
	return o.Card
}

// usedText is the words of the ability a stack item uses: the chosen
// mode's when it has one, else the whole ability's (S-07). Rolls say
// only what the item does (15.22): the item that rolls says the words up
// to the roll, its result trigger only the line for its result. The card
// shows every result: its image, or its text without images.
func (g *Game) usedText(it StackItem) string {
	d, ok := g.cards.find(it.Ability.Card)
	if !ok || it.Ability.Index < 0 || it.Ability.Index >= len(d.Abilities) {
		return ""
	}
	a := d.Abilities[it.Ability.Index]
	text := a.Text
	if it.Mode >= 0 && it.Mode < len(a.Modes) && a.Modes[it.Mode].Text != "" {
		text = a.Modes[it.Mode].Text
	}
	if it.RollResult > 0 {
		return rollResultText(text, it.RollResult)
	}
	return rollStartText(text)
}

// rollWord is where a roll starts in an ability's words: "roll-" or
// "rolls-"; rollHead is one result line's head: "1:", "1-2:", "1 or 6:".
var (
	rollWord = regexp.MustCompile(`(?i)\brolls?-`)
	rollHead = regexp.MustCompile(`(\d)(?:-(\d)| or (\d))?:\s`)
)

// rollStartText is what an ability does before its roll's result: the
// words up to the roll. "Roll- 1: Gain 1¢. 2: …" says "Roll." (15.22).
func rollStartText(text string) string {
	m := rollWord.FindStringIndex(text)
	if m == nil {
		return text
	}
	return text[:m[1]-1] + "."
}

// rollResultText is the line of an ability's words for one roll result:
// "1: Each player gains 1¢." A roll without a table keeps the words
// after the roll; a result with no line says it did nothing.
func rollResultText(text string, result int) string {
	m := rollWord.FindStringIndex(text)
	if m == nil {
		return text
	}
	rest := strings.TrimSpace(text[m[1]:])
	heads := rollHead.FindAllStringSubmatchIndex(rest, -1)
	// A head starts the words or follows the end of a line.
	heads = slices.DeleteFunc(heads, func(h []int) bool {
		return h[0] > 0 && !strings.HasSuffix(strings.TrimSpace(rest[:h[0]]), ".")
	})
	if len(heads) == 0 {
		return rest
	}
	for i, h := range heads {
		lo := digit(rest[h[2]:h[3]])
		hi, or := lo, -1
		if h[4] >= 0 {
			hi = digit(rest[h[4]:h[5]])
		}
		if h[6] >= 0 {
			or = digit(rest[h[6]:h[7]])
		}
		if (result >= lo && result <= hi) || result == or {
			end := len(rest)
			if i+1 < len(heads) {
				end = heads[i+1][0]
			}
			return strings.TrimSpace(rest[h[0]:end])
		}
	}
	return fmt.Sprintf("%d: Nothing happens.", result)
}

// digit is a one-digit number; rollHead matches only those.
func digit(s string) int {
	return int(s[0] - '0')
}
