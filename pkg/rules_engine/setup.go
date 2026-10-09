package rulesengine

import (
	"errors"
	"fmt"
)

// Setup describes a new game.
type Setup struct {
	Seed       uint64
	Players    int
	Sets       []CardSet
	BonusSouls bool // play with 3 random bonus souls (R-SETUP-06)
	// Characters, when set, gives seat i the character Characters[i]
	// instead of a random one (character picks, tests).
	Characters []CardRef
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
	for i := range s.Players {
		g.Players = append(g.Players, Player{ID: PlayerID(i)})
	}

	var characters, bonus []CardRef
	for _, d := range idx.defs {
		if d.Outside {
			continue
		}
		deck, hasDeck := d.Kind.Deck()
		for range max(d.Copies, 1) {
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
	for d := range deckCount {
		g.ShuffleDeck(d) // R-SETUP-01
	}
	if err := g.fillSlots(); err != nil {
		return nil, nil, err
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
	g.Turn = Turn{Active: g.firstPlayer(), Number: 1, Step: StepRecharge}
	g.emit(Event{Kind: EvGameStarted, Player: g.Turn.Active})
	g.run()
	return g, g.takeEvents(), nil
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
		g.Rooms[i].Cards = append(g.Rooms[i].Cards, nid)
	}
	g.emit(Event{Kind: EvCardRevealed, Player: NoPlayer, Object: nid, Card: o.Card})
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
		o.Role, o.Charged = RoleItem, true // starts charged (R-SETUP-08)
		g.Players[i].InPlay = append(g.Players[i].InPlay, item)
	}
	return nil
}

// firstPlayer is the player whose character says they go first (Cain),
// or else the winner of the roll.
func (g *Game) firstPlayer() PlayerID {
	for _, pl := range g.Players {
		if g.def(pl.Character).GoesFirst {
			g.emit(Event{Kind: EvFirstPlayer, Player: pl.ID, Text: "character"})
			return pl.ID
		}
	}
	return g.rollForFirst()
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
	g.ask(Choice{Purpose: ChooseStartingItem, Player: p, Rule: "R-SETUP-09", Objects: top}, g.labels(top))
}

// chooseStartingItem puts the chosen card into play as an eternal item;
// the others go to the bottom of the treasure deck.
func (g *Game) chooseStartingItem(c Choice, i int) {
	for _, id := range c.Objects {
		g.Decks[TreasureDeck] = remove(g.Decks[TreasureDeck], id)
	}
	for j, id := range c.Objects {
		if j == i {
			continue
		}
		nid := g.move(id, DeckZone(TreasureDeck), NoPlayer)
		g.Decks[TreasureDeck] = append([]ObjectID{nid}, g.Decks[TreasureDeck]...)
	}
	nid := g.move(c.Objects[i], Zone{Kind: ZoneInPlay}, c.Player)
	o := g.Object(nid)
	o.Role, o.Charged, o.Eternal = RoleItem, true, true
	g.Players[c.Player].InPlay = append(g.Players[c.Player].InPlay, nid)
	g.emit(Event{Kind: EvGainedTreasure, Player: c.Player, Object: nid, Card: o.Card, Text: "starting item"})
}

// rollForFirst: each player rolls a D6; the lowest goes first; ties among
// the lowest roll again (R-SETUP-11, GS-06).
func (g *Game) rollForFirst() PlayerID {
	contenders := make([]PlayerID, len(g.Players))
	for i := range contenders {
		contenders[i] = PlayerID(i)
	}
	for len(contenders) > 1 {
		lowest := 7
		var low []PlayerID
		for _, p := range contenders {
			r := g.RNG.D6()
			g.emit(Event{Kind: EvDiceRolled, Player: p, Amount: r, Text: "first player"})
			switch {
			case r < lowest:
				lowest, low = r, []PlayerID{p}
			case r == lowest:
				low = append(low, p)
			}
		}
		contenders = low
	}
	g.emit(Event{Kind: EvFirstPlayer, Player: contenders[0]})
	return contenders[0]
}
