package world

import (
	"slices"
	"strconv"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// decks is the number of decks a view counts.
const decks = len(engine.GameView{}.DeckSizes)

// Memory is what one seat remembers of cards it saw go where it cannot
// see them now (ADR 030, Decision 19): the cards of other hands, and
// the cards at known places of each deck. It is fed only what the
// seat's client gets, the seat's view and its filtered events, step by
// step; Known gives it to Imagine through For.
//
// When unsure, it forgets: a card leaves a hand unseen, a deck is
// shuffled or changes in a way no event tells, or a count does not add
// up (Decision 19.9). Test 6 of 21.15 checks it against real games.
//
// A Memory is for one goroutine at a time; it is not saved.
type Memory struct {
	seat    engine.PlayerID
	started bool
	// hands are the remembered cards of each seat's hand; the seat's own
	// stays empty. sizes are the hand sizes of the last view.
	hands [][]handCard
	sizes []int
	// deck is each deck as remembered, the bottom first: one place per
	// card of the last view's deck, "" where the card is unknown.
	deck [decks][]engine.CardRef
}

// handCard is a remembered card in a hand: its object, when the seat
// saw it, else 0.
type handCard struct {
	id   engine.ObjectID
	card engine.CardRef
}

// NewMemory is an empty memory of seat p, for the start of a game.
func NewMemory(p engine.PlayerID) *Memory {
	return &Memory{seat: p}
}

// Reset forgets everything: after a load or a rollback, the steps in
// between were not observed. The next Observe starts over.
func (m *Memory) Reset() {
	*m = Memory{seat: m.seat}
}

// Known is what the memory holds now, for Imagine (Decision 2.7). The
// result shares nothing with the memory.
func (m *Memory) Known() engine.Known {
	var k engine.Known
	if !m.started {
		return k
	}
	k.Hands = make([][]engine.CardRef, len(m.hands))
	for q, hand := range m.hands {
		for _, h := range hand {
			k.Hands[q] = append(k.Hands[q], h.card)
		}
	}
	for d, cards := range m.deck {
		half := len(cards) / 2
		for i := len(cards) - 1; i >= half; i-- {
			k.Tops[d] = append(k.Tops[d], cards[i])
		}
		for i := range half {
			k.Bottoms[d] = append(k.Bottoms[d], cards[i])
		}
		k.Tops[d] = trimUnknown(k.Tops[d])
		k.Bottoms[d] = trimUnknown(k.Bottoms[d])
	}
	return k
}

// trimUnknown drops the unknown places at the end of a list.
func trimUnknown(cards []engine.CardRef) []engine.CardRef {
	for len(cards) > 0 && cards[len(cards)-1] == "" {
		cards = cards[:len(cards)-1]
	}
	if len(cards) == 0 {
		return nil
	}
	return cards
}

// Observe takes one step as the seat's client gets it: the seat's view
// after the step (Game.ViewFor) and the step's events filtered for the
// seat (FilterEvents). It must see every step, in order; the first call
// only takes the view as it is.
func (m *Memory) Observe(v engine.GameView, events []engine.Event) {
	if !m.started || len(v.Players) != len(m.hands) {
		m.start(v)
		return
	}
	s := step{m: m, moved: make([]int, len(m.hands))}
	for _, e := range events {
		s.event(e)
	}
	s.endDraws()
	s.finish(v)
}

// start takes a view with nothing remembered.
func (m *Memory) start(v engine.GameView) {
	m.started = true
	m.hands = make([][]handCard, len(v.Players))
	m.sizes = make([]int, len(v.Players))
	for i, pl := range v.Players {
		m.sizes[i] = pl.HandSize
	}
	for d := range m.deck {
		m.deck[d] = make([]engine.CardRef, v.DeckSizes[d])
	}
	m.see(v)
}

// step is one Observe: the memory and what the step's events left in
// doubt.
type step struct {
	m *Memory
	// lost marks the decks whose memory the step spoiled.
	lost [decks]bool
	// moved counts cards into each hand, less cards out of it.
	moved []int
	// draws are the looted events of a run not yet told by loot_done.
	draws []engine.Event
}

// event takes one event.
func (s *step) event(e engine.Event) {
	if e.Kind == engine.EvLooted && e.Text != engine.TextFromDiscard {
		// A draw from the loot deck, or a card taken out of it
		// (DeckToHand): loot_done tells which.
		s.addToHand(e.Player, 1)
		s.draws = append(s.draws, e)
		return
	}
	if e.Kind == engine.EvLootDone {
		s.lootDone(e)
		return
	}
	s.endDraws()
	switch e.Kind {
	case engine.EvLookedAt:
		s.lookedAt(e)
	case engine.EvTopChoice:
		s.topChoice(e)
	case engine.EvCardRevealed:
		if d, ok := deckNamed(e.Text, ""); ok && e.Card != "" {
			s.set(d, 0, e.Card) // the top shown to all (RevealTop)
		}
	case engine.EvMovedToDeck:
		s.movedToDeck(e)
	case engine.EvDeckBottom, engine.EvDeckOrder:
		if d, ok := deckNamed(e.Text, ""); ok {
			s.lose(d)
		}
	case engine.EvDiscarded:
		s.discarded(e)
	case engine.EvLootPlayed:
		s.leaveHand(e.Player, e.Card)
	case engine.EvGaveCard:
		s.gaveCard(e)
	case engine.EvLooted: // from the discard: in sight of all
		s.addToHand(e.Player, 1)
		s.m.learn(e.Player, e.Object, e.Card)
	default:
	}
}

// lookedAt learns a card the seat looked at: in a deck, or in a hand.
func (s *step) lookedAt(e engine.Event) {
	if e.Player != s.m.seat || e.Card == "" {
		return
	}
	if e.Deck != nil {
		s.set(int(e.Deck.Deck), e.Deck.Depth, e.Card)
		return
	}
	if e.By > 0 {
		s.m.learn(engine.PlayerID(e.By-1), e.Object, e.Card)
	}
}

// topChoice: a player looked at a deck's top card and kept it there or
// put it on the bottom.
func (s *step) topChoice(e engine.Event) {
	mine := e.Player == s.m.seat && e.Card != ""
	if d, ok := deckNamed(e.Text, engine.DeckTextTop); ok {
		if mine {
			s.set(d, 0, e.Card)
		}
		return
	}
	d, ok := deckNamed(e.Text, engine.DeckTextBottom)
	if !ok {
		return
	}
	card, ok := s.pop(d)
	if !ok {
		return
	}
	if mine {
		if card != "" && card != e.Card {
			s.lose(d)
			return
		}
		card = e.Card
	}
	s.insert(d, len(s.m.deck[d]), card)
}

// movedToDeck puts a card on a deck, under its top Amount cards, or on
// its bottom; one from a hand (HandToDeckTop) leaves that hand.
func (s *step) movedToDeck(e engine.Event) {
	if d, ok := deckNamed(e.Text, engine.DeckTextBottom); ok {
		s.insert(d, len(s.m.deck[d]), e.Card)
		return
	}
	d, ok := deckNamed(e.Text, "")
	if !ok {
		return // under a card, outside the game
	}
	s.insert(d, e.Amount, e.Card)
	if e.Private && e.Player != engine.NoPlayer {
		s.leaveHand(e.Player, e.Card)
	}
}

// discarded: the top of a deck milled, or a loot card discarded from a
// hand; other discards touch nothing remembered.
func (s *step) discarded(e engine.Event) {
	if d, ok := deckNamed(e.Text, ""); ok {
		card, ok := s.pop(d)
		if ok && card != "" && e.Card != "" && card != e.Card {
			s.lose(d)
		}
		return
	}
	if e.Player != engine.NoPlayer && e.Text == "" && e.Prev == 0 {
		s.leaveHand(e.Player, e.Card)
	}
}

// gaveCard: a hand card went from one player to another; the seat sees
// it if it gave or got it.
func (s *step) gaveCard(e engine.Event) {
	from, err := strconv.Atoi(e.Text)
	if err != nil || from < 0 || from >= len(s.moved) {
		s.m.forgetHands()
		return
	}
	s.leaveHand(engine.PlayerID(from), e.Card)
	s.addToHand(e.Player, 1)
	s.m.learn(e.Player, e.Object, e.Card)
}

// lootDone: the last Amount looted events of its player were draws from
// the top of the loot deck. Each takes the remembered top card along.
func (s *step) lootDone(e engine.Event) {
	n := e.Amount
	run := s.draws
	s.draws = nil
	if n <= 0 || len(run) != n {
		s.lose(int(engine.LootDeck))
		return
	}
	for _, d := range run {
		if d.Player != e.Player {
			s.lose(int(engine.LootDeck))
			return
		}
	}
	for _, d := range run {
		card, ok := s.pop(int(engine.LootDeck))
		if !ok {
			return
		}
		switch {
		case d.Card != "":
			if card != "" && card != d.Card {
				s.lose(int(engine.LootDeck))
				return
			}
		case card != "":
			s.m.learn(d.Player, 0, card)
		}
	}
}

// endDraws closes a run of looted events no loot_done told: cards taken
// from somewhere in the loot deck.
func (s *step) endDraws() {
	if len(s.draws) > 0 {
		s.draws = nil
		s.lose(int(engine.LootDeck))
	}
}

// addToHand counts n cards into p's hand.
func (s *step) addToHand(p engine.PlayerID, n int) {
	if p >= 0 && int(p) < len(s.moved) {
		s.moved[p] += n
	}
}

// leaveHand takes a card out of p's hand: every remembered copy of it,
// as the seat cannot tell which copy left; all of them when the card is
// hidden.
func (s *step) leaveHand(p engine.PlayerID, card engine.CardRef) {
	if p < 0 || int(p) >= len(s.moved) {
		return
	}
	s.moved[p]--
	if card == "" {
		s.m.hands[p] = nil
		return
	}
	s.m.hands[p] = slices.DeleteFunc(s.m.hands[p], func(h handCard) bool { return h.card == card })
}

// lose spoils deck d's memory for the rest of the step.
func (s *step) lose(d int) {
	if d >= 0 && d < decks {
		s.lost[d] = true
	}
}

// set learns that the card depth places below deck d's top is card. A
// different card remembered there is a doubt.
func (s *step) set(d, depth int, card engine.CardRef) {
	if d < 0 || d >= decks || s.lost[d] {
		return
	}
	cards := s.m.deck[d]
	i := len(cards) - 1 - depth
	if i < 0 || depth < 0 {
		s.lose(d)
		return
	}
	if cards[i] != "" && cards[i] != card {
		s.lose(d)
		return
	}
	cards[i] = card
}

// insert puts card (or an unknown one, "") into deck d with depth cards
// above it.
func (s *step) insert(d, depth int, card engine.CardRef) {
	if d < 0 || d >= decks || s.lost[d] {
		return
	}
	cards := s.m.deck[d]
	i := len(cards) - depth
	if i < 0 || depth < 0 {
		s.lose(d)
		return
	}
	s.m.deck[d] = slices.Insert(cards, i, card)
}

// pop takes deck d's top card; ok is false when the deck is spoilt or
// remembered empty (it was refilled, out of sight).
func (s *step) pop(d int) (card engine.CardRef, ok bool) {
	if d < 0 || d >= decks || s.lost[d] {
		return "", false
	}
	cards := s.m.deck[d]
	if len(cards) == 0 {
		s.lose(d)
		return "", false
	}
	s.m.deck[d] = cards[:len(cards)-1]
	return cards[len(cards)-1], true
}

// finish checks the step's counts against the view and takes what the
// view shows: a count that does not add up spoils that memory.
func (s *step) finish(v engine.GameView) {
	m := s.m
	for d := range m.deck {
		if s.lost[d] || len(m.deck[d]) != v.DeckSizes[d] {
			m.deck[d] = make([]engine.CardRef, v.DeckSizes[d])
		}
	}
	for q, pl := range v.Players {
		if m.sizes[q]+s.moved[q] != pl.HandSize || len(m.hands[q]) > pl.HandSize {
			m.hands[q] = nil
		}
		m.sizes[q] = pl.HandSize
	}
	m.see(v)
}

// see takes what the view shows now: the hands it shows and the
// treasure deck's top the seat peeks at.
func (m *Memory) see(v engine.GameView) {
	for q, pl := range v.Players {
		if engine.PlayerID(q) == m.seat {
			continue
		}
		if pl.Hand != nil {
			m.hands[q] = m.hands[q][:0]
			for _, c := range pl.Hand {
				m.hands[q] = append(m.hands[q], handCard{id: c.ID, card: c.Card})
			}
		}
	}
	if int(m.seat) < len(v.Players) {
		if top := v.Players[m.seat].TreasureTop; top != nil {
			cards := m.deck[engine.TreasureDeck]
			if n := len(cards); n > 0 {
				if cards[n-1] != "" && cards[n-1] != top.Card {
					cards = make([]engine.CardRef, n)
					m.deck[engine.TreasureDeck] = cards
				}
				cards[n-1] = top.Card
			}
		}
	}
}

// learn remembers card in p's hand, as object id (0: not known). A card
// seen before as the same object, or remembered there with no object,
// is not counted twice.
func (m *Memory) learn(p engine.PlayerID, id engine.ObjectID, card engine.CardRef) {
	if p == m.seat || p < 0 || int(p) >= len(m.hands) || card == "" {
		return
	}
	hand := m.hands[p]
	if id != 0 {
		if slices.ContainsFunc(hand, func(h handCard) bool { return h.id == id }) {
			return
		}
		if i := slices.IndexFunc(hand, func(h handCard) bool { return h.id == 0 && h.card == card }); i >= 0 {
			hand[i].id = id
			return
		}
	}
	m.hands[p] = append(hand, handCard{id: id, card: card})
}

// forgetHands forgets every hand: an event the memory cannot read.
func (m *Memory) forgetHands() {
	for q := range m.hands {
		m.hands[q] = nil
	}
}

// deckNamed reads "<deck> deck" plus where, as events name decks
// (engine.DeckText).
func deckNamed(text, where string) (int, bool) {
	d, ok := engine.DeckOfText(text, where)
	return int(d), ok
}
