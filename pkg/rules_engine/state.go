package rulesengine

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
)

// PlayerID is a seat; seats are in turn order (R-TURN-01).
type PlayerID int

// NoPlayer marks objects controlled by the game, not by a player.
const NoPlayer PlayerID = -1

// ObjectID names one object. A card that changes zone becomes a new
// object with a new ID (R-ZONE-01); old IDs are never reused.
type ObjectID int

// CardRef names a card version, e.g. "the_d6" or "the_d6@2" (ADR 004).
type CardRef string

// DeckKind is one of the decks (R-ZONE-04).
type DeckKind int

// The decks. Rooms come with later sets.
const (
	TreasureDeck DeckKind = iota
	LootDeck
	MonsterDeck
	RoomDeck
	deckCount
)

// ZoneKind is where an object is (R-ZONE).
type ZoneKind int

// The zones.
const (
	ZoneGone    ZoneKind = iota // the object moved on and became a new object
	ZoneDeck                    // hidden (R-ZONE-04)
	ZoneDiscard                 // public (R-ZONE-07)
	ZoneHand                    // hidden from other players (R-ZONE-08)
	ZoneInPlay                  // public (R-ZONE-02)
	ZoneCovered                 // under another card in a slot (R-ZONE-10)
	ZoneStack                   // on the stack (R-STACK-01)
	ZoneOutside                 // outside the game (R-ZONE-13)
)

// SlotKind is a row of slots on the table.
type SlotKind int

// The slot rows.
const (
	ShopSlot SlotKind = iota
	MonsterSlot
	RoomSlot
)

// Zone tells exactly where an object is.
type Zone struct {
	Kind ZoneKind `json:"kind"`
	Deck DeckKind `json:"deck,omitempty"` // for decks and discards
	Slot SlotKind `json:"slot,omitempty"` // for objects in or under a slot
	// Index is the slot number for slot objects.
	Index int `json:"index,omitempty"`
}

// Role is what an object in play is (R-CARD).
type Role int

// The roles of objects in play.
const (
	RoleNone      Role = iota // not in play
	RoleCharacter             // R-CARD-25
	RoleItem                  // R-CARD-01
	RoleMonster               // R-CARD-09
	RoleEvent                 // R-CARD-09
	RoleSoul                  // R-CARD-18
	RoleCurse                 // R-ABIL-20
)

// Object is one card in one zone.
type Object struct {
	ID         ObjectID  `json:"id"`
	Card       CardRef   `json:"card"`
	Zone       Zone      `json:"zone"`
	Role       Role      `json:"role,omitempty"`
	Controller PlayerID  `json:"controller"`
	Charged    bool      `json:"charged"`
	Damage     int       `json:"damage,omitempty"`
	Counters   []Counter `json:"counters,omitempty"`
}

// Counter is a counter on an object; Name is "" for a generic counter
// (R-MECH-12).
type Counter struct {
	Name  string `json:"name,omitempty"`
	Count int    `json:"count"`
}

// Player is one seat at the table.
type Player struct {
	ID        PlayerID   `json:"id"`
	Character ObjectID   `json:"character"`
	Cents     int        `json:"cents"`
	Hand      []ObjectID `json:"hand"`
	// InPlay holds items, souls and curses the player controls.
	InPlay []ObjectID `json:"in_play"`
	Damage int        `json:"damage,omitempty"` // tied to the player (R-MECH-19)
	Dead   bool       `json:"dead,omitempty"`   // died this turn (R-DEATH-17)
}

// Slot is one table slot; the last card is on top and in play, the rest
// are covered (R-ZONE-10).
type Slot struct {
	Cards []ObjectID `json:"cards"`
}

// Game is the whole game state. It holds no Go maps, so it encodes to
// JSON the same way every time (A-08).
type Game struct {
	RNG     RNG      `json:"rng"`
	Players []Player `json:"players"`
	Objects []Object `json:"objects"` // index = ObjectID

	// Decks and discards: the last element is the top card.
	Decks    [deckCount][]ObjectID `json:"decks"`
	Discards [deckCount][]ObjectID `json:"discards"`

	Shop     []Slot `json:"shop"`
	Monsters []Slot `json:"monsters"`
	Rooms    []Slot `json:"rooms,omitempty"`

	Outside []ObjectID `json:"outside"`
	// BonusSouls are the active bonus souls, outside the game (R-SETUP-06).
	BonusSouls []ObjectID `json:"bonus_souls,omitempty"`

	// Queue holds pending actions that cards may rewrite (A-05).
	Queue []Action `json:"queue,omitempty"`
	// Attack is the attack in progress (R-ATK).
	Attack AttackState `json:"attack"`

	// Choice is the open question of a choose prompt.
	Choice *Choice `json:"choice,omitempty"`

	// Stack: the last item is on top (R-STACK-02).
	Stack    []StackItem `json:"stack"`
	StackSeq int         `json:"stack_seq"`

	Turn     Turn       `json:"turn"`
	Priority Priority   `json:"priority"`
	Waiting  Prompt     `json:"waiting"`
	MaxHand  int        `json:"max_hand"`  // default 10 (R-TURN-11)
	WinSouls int        `json:"win_souls"` // default 4 (R-WIN-01)
	Over     bool       `json:"over,omitempty"`
	Winners  []PlayerID `json:"winners,omitempty"`

	// Sets lists the card sets in play; definitions are looked up by name
	// when a saved game is loaded.
	Sets []string `json:"sets"`

	cards  cardIndex
	events []Event
	// forcedRolls lets tests decide dice results; never set in games.
	forcedRolls []int
}

// newObject adds a card as a new object and returns its ID.
func (g *Game) newObject(card CardRef, zone Zone, controller PlayerID) ObjectID {
	id := ObjectID(len(g.Objects))
	g.Objects = append(g.Objects, Object{ID: id, Card: card, Zone: zone, Controller: controller})
	return id
}

// Object returns the object with that ID.
func (g *Game) Object(id ObjectID) *Object {
	if id < 0 || int(id) >= len(g.Objects) {
		panic(fmt.Sprintf("rulesengine: no object %d", id))
	}
	return &g.Objects[id]
}

// move puts the card of id into a new zone as a new object (R-ZONE-01).
// The old object is gone; callers update the zone lists.
func (g *Game) move(id ObjectID, to Zone, controller PlayerID) ObjectID {
	old := g.Object(id)
	card := old.Card
	old.Zone = Zone{Kind: ZoneGone}
	return g.newObject(card, to, controller)
}

// Clone returns a deep copy of the game.
func (g *Game) Clone() (*Game, error) {
	data, err := json.Marshal(g)
	if err != nil {
		return nil, fmt.Errorf("clone: %w", err)
	}
	var c Game
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("clone: %w", err)
	}
	c.cards = g.cards // definitions are not part of the saved state
	return &c, nil
}

// Checksum is a stable hash of the whole state (RP-12).
func (g *Game) Checksum() (uint64, error) {
	data, err := json.Marshal(g)
	if err != nil {
		return 0, fmt.Errorf("checksum: %w", err)
	}
	h := fnv.New64a()
	if _, err := h.Write(data); err != nil {
		return 0, fmt.Errorf("checksum: %w", err)
	}
	return h.Sum64(), nil
}
