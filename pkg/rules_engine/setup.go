package rulesengine

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// Setup describes a new game.
type Setup struct {
	Seed       uint64
	Players    int
	Sets       []CardSet
	BonusSouls bool // play with 3 random bonus souls (R-SETUP-06)
	// Rooms plays with the room deck (R-SETUP-16); without it, room
	// cards stay out of the game.
	Rooms bool
	// Characters, when set, gives seat i the character Characters[i]
	// instead of a random one (character picks, tests).
	Characters []CardRef
	// Banned cards are left out of the game: no copy goes into a deck,
	// and none is dealt (GS-13).
	Banned []CardRef
	// DeckRatios builds the treasure, loot and monster decks from the
	// official count per group, 100 cards each (R-SETUP-17, GS-07).
	DeckRatios bool
}

// Defaults of the official rules.
const (
	startingLoot  = 3 // R-SETUP-10
	startingCents = 3 // R-SETUP-10
	defaultHand   = 10
	defaultSouls  = 4
	activeBonus   = 3
)

// NewGame sets up a game (R-SETUP) and runs it until the first prompt.
func NewGame(s Setup) (*Game, []Event, error) {
	if s.Players < 2 || s.Players > 4 {
		return nil, nil, fmt.Errorf("setup: %d players; the game is for 2 to 4 (GS-09)", s.Players)
	}
	idx, err := newCardIndex(s.Sets...)
	if err != nil {
		return nil, nil, fmt.Errorf("setup: %w", err)
	}
	g := &Game{RNG: NewRNG(s.Seed), cards: idx, MaxHand: defaultHand, WinSouls: defaultSouls}
	for _, set := range s.Sets {
		g.Sets = append(g.Sets, set.Name)
	}
	for _, b := range s.Banned {
		d, ok := idx.find(b)
		if !ok {
			return nil, nil, fmt.Errorf("setup: banned card %q is in none of the sets", b)
		}
		if d.Outside {
			// A starting item comes with its character; ban the character.
			return nil, nil, fmt.Errorf("setup: banned card %q starts outside the game", b)
		}
	}
	g.Banned = slices.Clone(s.Banned)
	for i := range s.Players {
		g.Players = append(g.Players, Player{ID: PlayerID(i)})
	}

	var characters, bonus []CardRef
	for i, d := range idx.defs {
		if !inDecks(d, s) {
			continue
		}
		deck, hasDeck := d.Kind.Deck()
		if hasDeck && s.DeckRatios && d.Kind != RoomCard {
			continue // added by addRatioDecks
		}
		for range idx.copies(i) {
			switch {
			case hasDeck:
				g.AddToDeck(deck, d.Ref)
			case d.Kind == CharacterCard:
				characters = append(characters, d.Ref)
			case d.Kind == BonusSoulCard:
				bonus = append(bonus, d.Ref)
			}
		}
	}
	if s.DeckRatios {
		if err := g.addRatioDecks(s); err != nil {
			return nil, nil, err
		}
	}
	for d := range deckCount {
		g.ShuffleDeck(d) // R-SETUP-01
	}
	if err := g.fillSlots(); err != nil {
		return nil, nil, err
	}
	if len(g.Decks[RoomDeck]) > 0 {
		g.Rooms = append(g.Rooms, Slot{}) // R-SETUP-12
		id, _ := g.drawTop(RoomDeck)
		g.putInSlot(id, RoomSlot, 0)
	}
	if s.BonusSouls {
		g.pickBonusSouls(bonus)
	}
	if len(s.Characters) > 0 {
		if len(s.Characters) != s.Players {
			return nil, nil, fmt.Errorf("setup: %d characters for %d players", len(s.Characters), s.Players)
		}
		for _, c := range s.Characters {
			if d, ok := idx.find(c); !ok || d.Kind != CharacterCard {
				return nil, nil, fmt.Errorf("setup: %q is not a character", c)
			}
			if slices.Contains(s.Banned, c) {
				return nil, nil, fmt.Errorf("setup: %q is banned", c)
			}
		}
		characters = s.Characters
	}
	if err := g.dealCharacters(characters, len(s.Characters) == 0); err != nil {
		return nil, nil, err
	}
	for p := range g.Players {
		g.loot(PlayerID(p), startingLoot)
		g.Players[p].Cents = startingCents
	}
	for p := range g.Players {
		// Start-of-game abilities resolve at once, without priority (R-SETUP-09).
		if n := g.def(g.Players[p].Character).StartingChoice; n > 0 {
			g.enqueue(Action{Kind: ActChooseStartingItem, Player: PlayerID(p), Amount: n})
		}
	}
	g.chooseFirstPlayer()
	g.run()
	return g, g.takeEvents(), nil
}

// startGame starts turn 1 with the first player p.
func (g *Game) startGame(p PlayerID) {
	g.Turn = Turn{Active: p, Number: 1, Step: StepRecharge}
	g.emit(Event{Kind: EvGameStarted, Player: p})
}

// fillSlots reveals 2 shop items and 2 monsters (R-SETUP-03 to R-SETUP-05).
func (g *Game) fillSlots() error {
	for i := range 2 {
		id, ok := g.drawTop(TreasureDeck)
		if !ok {
			return errors.New("setup: the treasure deck is too small")
		}
		g.Shop = append(g.Shop, Slot{})
		g.putInSlot(id, ShopSlot, i)
	}
	for i := range 2 {
		g.Monsters = append(g.Monsters, Slot{})
		placed := false
		for tries := len(g.Decks[MonsterDeck]); tries >= 0; tries-- {
			id, ok := g.drawTop(MonsterDeck)
			if !ok {
				break
			}
			if g.kindOf(id) == EventCard {
				// Events go to the bottom during setup (R-SETUP-05).
				nid := g.move(id, DeckZone(MonsterDeck), NoPlayer)
				g.Decks[MonsterDeck] = append([]ObjectID{nid}, g.Decks[MonsterDeck]...)
				continue
			}
			g.putInSlot(id, MonsterSlot, i)
			placed = true
			break
		}
		if !placed {
			return errors.New("setup: the monster deck has too few monsters")
		}
	}
	return nil
}

// putInSlot puts a card on top of a slot, in play (R-ZONE-02).
func (g *Game) putInSlot(id ObjectID, kind SlotKind, i int) ObjectID {
	return g.putInSlotFrom(id, kind, i, 0)
}

// putInSlotFrom is putInSlot whose card_revealed names prev: the covered
// object that left under a new ID. The client then shows no flight from
// the deck for a card that was in the slot already (19.72).
func (g *Game) putInSlotFrom(id ObjectID, kind SlotKind, i int, prev ObjectID) ObjectID {
	nid := g.move(id, Zone{Kind: ZoneInPlay, Slot: kind, Index: i}, NoPlayer)
	o := g.Object(nid)
	o.Charged = true
	switch kind {
	case ShopSlot:
		o.Role = RoleItem
		g.Shop[i].Cards = append(g.Shop[i].Cards, nid)
	case MonsterSlot:
		o.Role = RoleMonster
		if g.kindOf(nid) == EventCard {
			o.Role = RoleEvent
		}
		g.Monsters[i].Cards = append(g.Monsters[i].Cards, nid)
	case RoomSlot:
		o.Role = RoleRoom
		if g.Turn.Step >= StepEndTriggers {
			g.Turn.LateRooms = append(g.Turn.LateRooms, nid) // R-TURN-18
		}
		g.Rooms[i].Cards = append(g.Rooms[i].Cards, nid)
	}
	g.emit(Event{Kind: EvCardRevealed, Player: NoPlayer, Object: nid, Card: o.Card, Prev: prev})
	// A card entering a slot enters play; an event does so in
	// enterMonsterSlot. Setup cards (turn 0) trigger nothing (19.25).
	if g.Turn.Number > 0 && o.Role != RoleEvent {
		g.emit(Event{Kind: EvEnteredPlay, Player: NoPlayer, Object: nid, Card: o.Card})
	}
	return nid
}

func (g *Game) kindOf(id ObjectID) CardKind {
	d, ok := g.cards.find(g.Object(id).Card)
	if !ok {
		panic(fmt.Sprintf("rulesengine: card %s has no definition", g.Object(id).Card))
	}
	return d.Kind
}

// pickBonusSouls sets aside 3 random bonus souls outside the game (R-SETUP-06).
func (g *Game) pickBonusSouls(pool []CardRef) {
	g.RNG.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
	for _, ref := range pool[:min(activeBonus, len(pool))] {
		g.BonusSouls = append(g.BonusSouls, g.newObject(ref, Zone{Kind: ZoneOutside}, NoPlayer))
	}
}

// dealCharacters gives each player a random character and its starting
// item (R-SETUP-07, R-SETUP-08, R-CARD-26).
func (g *Game) dealCharacters(pool []CardRef, random bool) error {
	if len(pool) < len(g.Players) {
		return fmt.Errorf("setup: %d characters for %d players", len(pool), len(g.Players))
	}
	if random {
		g.RNG.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
	}
	for i := range g.Players {
		p := PlayerID(i)
		ch := g.newObject(pool[i], Zone{Kind: ZoneInPlay}, p)
		g.Object(ch).Role = RoleCharacter // deactivated at start (R-SETUP-08)
		g.Players[i].Character = ch
		g.emit(Event{Kind: EvCharacterDealt, Player: p, Object: ch, Card: pool[i]})

		def, _ := g.cards.find(pool[i])
		if def.StartingItem == "" {
			continue
		}
		item := g.newObject(def.StartingItem, Zone{Kind: ZoneInPlay}, p)
		o := g.Object(item)
		o.Role, o.Charged, o.StartingItem = RoleItem, true, true // starts charged (R-SETUP-08)
		g.Players[i].InPlay = append(g.Players[i].InPlay, item)
	}
	return nil
}

// chooseFirstPlayer starts the game with the player whose character
// says they go first (Cain), or else starts the first-player roll.
func (g *Game) chooseFirstPlayer() {
	for _, pl := range g.Players {
		if g.def(pl.Character).GoesFirst {
			g.emit(Event{Kind: EvFirstPlayer, Player: pl.ID, Card: g.Object(pl.Character).Card, Text: "character"})
			g.startGame(pl.ID)
			return
		}
	}
	all := make([]PlayerID, len(g.Players))
	for i := range all {
		all[i] = PlayerID(i)
	}
	g.FirstRoll = &FirstRoll{Round: 1, Players: all}
}

// askStartingItem shows p the top n treasure cards to pick a starting
// item from (Eden).
func (g *Game) askStartingItem(p PlayerID, n int) {
	deck := g.Decks[TreasureDeck]
	n = min(n, len(deck))
	if n == 0 {
		return
	}
	top := make([]ObjectID, 0, n)
	for i := len(deck) - 1; i >= len(deck)-n; i-- {
		top = append(top, deck[i])
	}
	g.ask(Choice{Purpose: ChooseStartingItem, Player: p, Rule: "R-SETUP-09", Objects: top}, g.Labels(top))
}

// chooseStartingItem puts the chosen card into play as an eternal item;
// the others go to the bottom of the treasure deck in random order (Q-22).
func (g *Game) chooseStartingItem(c Choice, i int) {
	for _, id := range c.Objects {
		g.Decks[TreasureDeck] = remove(g.Decks[TreasureDeck], id)
	}
	var rest []ObjectID
	for j, id := range c.Objects {
		if j == i {
			continue
		}
		nid := g.move(id, DeckZone(TreasureDeck), NoPlayer)
		g.Decks[TreasureDeck] = append(g.Decks[TreasureDeck], nid)
		rest = append(rest, nid)
	}
	g.ToBottomRandom(TreasureDeck, rest)
	nid := g.move(c.Objects[i], Zone{Kind: ZoneInPlay}, c.Player)
	o := g.Object(nid)
	o.Role, o.Charged, o.Eternal, o.StartingItem = RoleItem, true, true, true
	g.Players[c.Player].InPlay = append(g.Players[c.Player].InPlay, nid)
	g.emit(Event{Kind: EvGainedTreasure, Player: c.Player, Object: nid, Card: o.Card, Text: "starting item"})
}

// FirstRoll is the first-player roll in progress (R-SETUP-11, GS-06,
// 15.43). Each player in it throws a D6, in seat order; the lowest goes
// first; ties among the lowest throw again.
type FirstRoll struct {
	// Round counts from 1; each tie starts the next.
	Round int `json:"round"`
	// Players still roll, in seat order.
	Players []PlayerID `json:"players"`
	// Rolls are this round's dice so far, one per player in order.
	Rolls []int `json:"rolls,omitempty"`
}

// askFirstThrow asks the next player of the first-player roll to throw.
// The game has not started: turn 0, nobody has priority.
func (g *Game) askFirstThrow() {
	fr := g.FirstRoll
	g.Throw = &Throw{First: true}
	g.Waiting = Prompt{Kind: PromptThrow, Player: fr.Players[len(fr.Rolls)], Text: eventFirstPlayerRoll.Msg()}
}

// throwFirst draws one die of the first-player roll. After the last die
// of a round, the lowest goes first, or the tied lowest roll again.
func (g *Game) throwFirst(p PlayerID) {
	fr := g.FirstRoll
	r := g.RNG.D6()
	fr.Rolls = append(fr.Rolls, r)
	g.emit(Event{Kind: EvDiceRolled, Player: p, Amount: r, Text: "first", Label: eventFirstPlayerRoll.Msg()})
	if len(fr.Rolls) < len(fr.Players) {
		return
	}
	lowest := slices.Min(fr.Rolls)
	var low []PlayerID
	for i, v := range fr.Rolls {
		if v == lowest {
			low = append(low, fr.Players[i])
		}
	}
	if len(low) > 1 {
		g.FirstRoll = &FirstRoll{Round: fr.Round + 1, Players: low}
		seats := make([]string, len(low))
		for i, q := range low {
			seats[i] = strconv.Itoa(int(q))
		}
		g.emit(Event{Kind: EvFirstRollTied, Player: NoPlayer, Amount: lowest, Text: strings.Join(seats, ",")})
		return
	}
	g.FirstRoll = nil
	g.emit(Event{Kind: EvFirstPlayer, Player: low[0], Amount: lowest})
	g.startGame(low[0])
}

// CharactersOutside lists the character cards outside the game that a
// player could get: the game's characters with a copy not in play, not
// banned, and not 3+ in a 2-player game (R-ZONE-15: The Clicker).
func (g *Game) CharactersOutside() []CardRef {
	var out []CardRef
	for i, d := range g.cards.defs {
		if d.Kind != CharacterCard || d.Outside || g.cards.borrowed[i] || slices.Contains(g.Banned, d.Ref) || (d.ThreePlus && len(g.Players) < 3) {
			continue
		}
		used := 0
		for _, pl := range g.Players {
			if g.Object(pl.Character).Card == d.Ref {
				used++
			}
		}
		if used < g.cards.copies(i) {
			out = append(out, d.Ref)
		}
	}
	return out
}

// ReplaceCharacter puts card into play as p's character, charged (Q-18);
// the old character and p's starting item go outside the game, and the
// new character's starting item enters play under p (The Clicker).
// Damage stays with the player (R-MECH-19); counters on the old
// character go with it (R-MECH-13).
func (g *Game) ReplaceCharacter(p PlayerID, card CardRef) {
	if d, ok := g.cards.find(card); !ok || d.Kind != CharacterCard {
		return
	}
	pl := &g.Players[p]
	old := pl.Character
	oldItem := g.def(old).StartingItem
	for _, id := range append([]ObjectID(nil), pl.InPlay...) {
		o := g.Object(id)
		if o.Role == RoleItem && (o.StartingItem || (oldItem != "" && o.Card == oldItem)) {
			pl.InPlay = remove(pl.InPlay, id)
			g.move(id, Zone{Kind: ZoneOutside}, NoPlayer)
			g.emit(Event{Kind: EvDiscarded, Player: p, Card: o.Card, Prev: id, Text: "outside the game"})
		}
	}
	g.move(old, Zone{Kind: ZoneOutside}, NoPlayer)
	ch := g.newObject(card, Zone{Kind: ZoneInPlay}, p)
	o := g.Object(ch)
	o.Role, o.Charged = RoleCharacter, true
	pl.Character = ch
	g.emit(Event{Kind: EvCharacterDealt, Player: p, Object: ch, Card: card, Prev: old, Text: "replaced"})
	if item := g.def(ch).StartingItem; item != "" {
		id := g.newObject(item, Zone{Kind: ZoneInPlay}, p)
		io := g.Object(id)
		io.Role, io.Charged, io.StartingItem = RoleItem, true, true
		pl.InPlay = append(pl.InPlay, id)
		g.emit(Event{Kind: EvGainedTreasure, Player: p, Object: id, Card: item, Text: "starting item"})
	}
}
