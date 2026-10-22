package rulesengine

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"maps"
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
	RoleRoom                  // R-CARD-33
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
	// Eternal is set when an effect makes the object eternal, e.g. Eden's
	// starting item (R-ABIL-18).
	Eternal bool `json:"eternal,omitempty"`
	// CopyOf is the card this object is a copy of; CopyThisTurn ends the
	// copy at the end of the turn.
	CopyOf       CardRef `json:"copy_of,omitempty"`
	CopyThisTurn bool    `json:"copy_this_turn,omitempty"`
	// CopyTillYourTurn ends the copy at the start of its controller's next
	// turn (Undefined).
	CopyTillYourTurn bool `json:"copy_till_your_turn,omitempty"`
	// Gained is a card whose abilities this object has too (R-ABIL-38:
	// Glitch); GainedTill ends it: GainEndOfTurn or GainTillYourTurn.
	Gained     CardRef `json:"gained,omitempty"`
	GainedTill int     `json:"gained_till,omitempty"`
	// ATKSet, if not 0, is the object's ATK plus 1 ("Its ATK becomes 1":
	// Hemoptysis).
	ATKSet int `json:"atk_set,omitempty"`
	// Flipped: the back face is up (R-MECH-33); the object acts as its
	// card's Back. Counters, damage and charge stay (R-MECH-53).
	Flipped bool `json:"flipped,omitempty"`
	// LostEternal: an effect took Eternal away ("This loses eternal":
	// Birthright).
	LostEternal bool `json:"lost_eternal,omitempty"`
	// NoAbilities: the card "loses all abilities" (a soul made by
	// BecomeSoul, The Chest). Other treasure souls keep theirs (R-CARD-38).
	NoAbilities bool `json:"no_abilities,omitempty"`
	// StartingItem marks a character's starting item (R-CARD-26), also
	// one picked at the start (Eden).
	StartingItem bool `json:"starting_item,omitempty"`
	// Under lists cards put under this object; they are outside the game
	// (R-ZONE-13: Friendly Ball).
	Under []ObjectID `json:"under,omitempty"`
	// KilledBy is who dealt the damage that brought it to 0 HP, plus 1
	// (0: nobody yet); KilledOn is the attack roll, if it was combat.
	KilledBy int `json:"killed_by,omitempty"`
	KilledOn int `json:"killed_on,omitempty"`
	// Chosen is a number the card's ability chose and keeps, not a
	// counter (Crystal Ball's guess); 0 if none.
	Chosen int `json:"chosen,omitempty"`
	// DoubleRewards doubles a dying monster's rewards (Dinga).
	DoubleRewards bool `json:"double_rewards,omitempty"`
	// ShareRewards also gain a dead card's rewards (Sibling Rivalry).
	ShareRewards []PlayerID `json:"share_rewards,omitempty"`
	// HitsThisTurn counts the times it took damage this turn.
	HitsThisTurn int `json:"hits_this_turn,omitempty"`
}

// Killer is the player who brought the object to 0 HP, or NoPlayer.
func (o *Object) Killer() PlayerID { return PlayerID(o.KilledBy - 1) }

// CountersOf returns how many counters named name the object has.
func (o *Object) CountersOf(name string) int {
	for _, c := range o.Counters {
		if c.Name == name {
			return c.Count
		}
	}
	return 0
}

// addCounters adds n counters (or removes, if negative), never below 0.
func (o *Object) addCounters(name string, n int) {
	for i, c := range o.Counters {
		if c.Name == name {
			o.Counters[i].Count = max(c.Count+n, 0)
			if o.Counters[i].Count == 0 {
				o.Counters = append(o.Counters[:i], o.Counters[i+1:]...)
			}
			return
		}
	}
	if n > 0 {
		o.Counters = append(o.Counters, Counter{Name: name, Count: n})
	}
}

// Shield prevents damage to Target: the whole next instance, or up to
// Amount of it when Amount > 0 (R-MECH-46). It lasts till end of turn.
type Shield struct {
	Target Target `json:"target"`
	Amount int    `json:"amount,omitempty"`
	// Cap > 0 reduces the next instance to at most Cap instead.
	Cap int `json:"cap,omitempty"`
	// Source is the object that made the shield (Host Hat reacts).
	Source ObjectID `json:"source,omitempty"`
}

// Boost is a "till end of turn" stat change (R-TURN-13). It applies to
// Object, or to Player when Object is 0.
type Boost struct {
	Stat   Stat     `json:"stat"`
	Player PlayerID `json:"player"`
	Object ObjectID `json:"object,omitempty"`
	Amount int      `json:"amount"`
	// Source is the object whose ability gave the boost, so a client can
	// show why a stat changed (V-41); 0 in states from before 15.54.
	Source ObjectID `json:"source,omitempty"`
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
	// CopyNextLoot: the next non-trinket loot p plays this turn is
	// copied (Blank Card).
	CopyNextLoot bool `json:"copy_next_loot,omitempty"`
	// SkipTurns: the player skips this many of their next turns (Famine).
	SkipTurns int `json:"skip_turns,omitempty"`
	// Activations and LootPlayed count this turn's activated abilities
	// and loot cards played (Shadow of Famine).
	Activations int `json:"activations,omitempty"`
	LootPlayed  int `json:"loot_played,omitempty"`
	// TimesDamaged counts this turn's damage ("the first time you take
	// damage each turn").
	TimesDamaged int `json:"times_damaged,omitempty"`
	// ExtraLootPlays come from abilities and last until the turn ends.
	ExtraLootPlays int `json:"extra_loot_plays,omitempty"`
	// Rolls and LastRoll are this turn's resolved rolls ("your first dice
	// roll each turn", "twice in a row": R-DICE-08).
	Rolls    int `json:"rolls,omitempty"`
	LastRoll int `json:"last_roll,omitempty"`
	// LastAttack is this turn's last attack roll ("twice in a row on an
	// attack roll": R-DICE-09).
	LastAttack int `json:"last_attack,omitempty"`
	// KilledBy is who dealt the damage or effect that killed the player
	// this turn, plus 1 (0: nobody).
	KilledBy int `json:"killed_by,omitempty"`
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
	// Purchase is a declared purchase (R-SHOP).
	Purchase PurchaseState `json:"purchase"`

	// Choice is the open question of a choose prompt.
	Choice *Choice `json:"choice,omitempty"`
	// Throw is the roll waiting for its die to be thrown (PromptThrow).
	Throw *Throw `json:"throw,omitempty"`
	// FirstRoll is the first-player roll before turn 1 (15.43).
	FirstRoll *FirstRoll `json:"first_roll,omitempty"`
	// Activating is an ability being put on the stack (targets pending).
	Activating *Activation `json:"activating,omitempty"`
	// FreeLoot is the loot card being played without a loot play
	// (Distant Admiration); 0 for none.
	FreeLoot ObjectID `json:"free_loot,omitempty"`
	// PendingTriggers wait to go on the stack (R-ABIL-14).
	PendingTriggers []PendingTrigger `json:"pending_triggers,omitempty"`
	// Boosts and Shields last till end of turn (R-TURN-13). A shield
	// prevents the next damage its target would take (R-MECH-46).
	// Compost: the next loot comes from the top of the loot discard.
	Compost bool `json:"compost,omitempty"`
	// ExtraTurn: the active player takes another turn after this one.
	ExtraTurn bool `json:"extra_turn,omitempty"`

	Boosts  []Boost  `json:"boosts,omitempty"`
	Shields []Shield `json:"shields,omitempty"`
	// Lingering abilities last till end of turn (Linger).
	Lingering []Lingering `json:"lingering,omitempty"`
	// Markers change numbers in effect text till end of turn (Magic Marker).
	Markers []Marker `json:"markers,omitempty"`
	// ManyDice: the next roll is made with Dice dice and Chooser picks the
	// result (Mutant Spider, R-DICE-11).
	ManyDice []ManyDice `json:"many_dice,omitempty"`

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
	// Banned lists the cards the host left out of the game (GS-13).
	Banned []CardRef `json:"banned,omitempty"`

	// AskSeq numbers Ask effects; asks holds the waiting ones (not saved:
	// see Asking.Key).
	AskSeq int `json:"ask_seq,omitempty"`
	asks   map[int]askEffect
	// maxHP is the max HP at the last checkMaxHP (R-MECH-21; not saved).
	maxHP map[hpKey]int

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
	c.asks = maps.Clone(g.asks)
	c.maxHP = maps.Clone(g.maxHP)
	return &c, nil
}

// Checksum is a stable hash of the whole state (RP-12).
func (g *Game) Checksum() (uint64, error) {
	// Encode straight into the hash: no copy of the state is kept.
	h := fnv.New64a()
	if err := json.NewEncoder(h).Encode(g); err != nil {
		return 0, fmt.Errorf("checksum: %w", err)
	}
	return h.Sum64(), nil
}

// inPlay lists the objects in play whose abilities work, in a fixed
// order: each player's character and play area in seat order, then the
// top of every monster and room slot. Souls are left out: a soul only
// counts for its value (R-CARD-18).
func (g *Game) inPlay() []ObjectID {
	var out []ObjectID
	for _, pl := range g.Players {
		out = append(out, pl.Character)
		for _, id := range pl.InPlay {
			if g.Object(id).Role != RoleSoul || g.soulKeepsAbilities(id) {
				out = append(out, id)
			}
		}
	}
	// Shop items' abilities do not work in the shop (R-CARD-05).
	for _, rows := range [][]Slot{g.Monsters, g.Rooms} {
		for _, s := range rows {
			if top, ok := s.TopOf(); ok && g.Object(top).Zone.Kind == ZoneInPlay {
				out = append(out, top)
			}
		}
	}
	return out
}

// ManyDice is a pending "roll several dice, one is the result".
type ManyDice struct {
	Chooser PlayerID `json:"chooser"`
	Dice    int      `json:"dice"`
	// Values and StackID are set once the roll is made.
	Values  []int `json:"values,omitempty"`
	StackID int   `json:"stack_id,omitempty"`
}

// soulKeepsAbilities: a treasure that became a soul keeps its abilities
// unless it says it loses them (R-CARD-38: Experimental Treatment).
func (g *Game) soulKeepsAbilities(id ObjectID) bool {
	o := g.Object(id)
	return o.Role == RoleSoul && !o.NoAbilities && g.kindOf(id) == TreasureCard
}

// How long a gained card's abilities last (Object.GainedTill).
const (
	GainEndOfTurn    = 1
	GainTillYourTurn = 2
)

// facesOf lists the cards whose abilities an object has: the card it
// acts as, then any it gained (R-ABIL-38: Glitch, Golden Trinket).
func (g *Game) facesOf(id ObjectID) []CardRef {
	out := []CardRef{g.CardOf(id)}
	o := g.Object(id)
	if o.Gained != "" {
		out = append(out, o.Gained)
	}
	if d, ok := g.cards.find(out[0]); ok && d.GainsAbilities != nil {
		out = append(out, d.GainsAbilities(g, id)...)
	}
	return out
}

// GainAbilities gives an object the abilities of card till the end of
// this turn, or till the start of its controller's next turn.
func (g *Game) GainAbilities(id ObjectID, card CardRef, till int) {
	o := g.Object(id)
	o.Gained, o.GainedTill = card, till
}
