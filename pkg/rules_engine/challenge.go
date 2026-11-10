package rulesengine

import (
	"fmt"
	"slices"
)

// The difficulties of a challenge (ADR 028, Decision 4).
const (
	DifficultyNormal = "normal"
	DifficultyHard   = "hard"
	DifficultyUltra  = "ultra"
)

// The kinds of a present (ADR 028, Decision 5; R-KRAMPUS-05).
const (
	PresentCents    = "cents"
	PresentLoot     = "loot"
	PresentTreasure = "treasure"
	PresentSoul     = "soul"
)

const (
	challengePlayers = 2 // R-COOP-01
	timerStart       = 8 // R-COOP-02
)

// ReasonTimer is why a challenge is lost when its timer reaches 0
// (R-COOP-09).
var ReasonTimer = NewKey("event.challenge_lost.timer")

// ChallengeDef is one challenge (ADR 028, Decision 3). Its cards are in
// the card set Set; its rules are the abilities of its boss and rule
// cards.
type ChallengeDef struct {
	// ID names the challenge in options, records and saves.
	ID string
	// Set is the code of the card set that holds its cards; the caller
	// puts that set into Setup.Sets.
	Set string
	// Levels are its difficulties, easiest first.
	Levels []ChallengeLevel
	// Setup does the challenge's own setup, after the co-op setup
	// (R-CHAL-01); nil for none.
	Setup func(g *Game, difficulty string) error
}

// ChallengeLevel is one difficulty of a challenge: its final boss and
// its rule card.
type ChallengeLevel struct {
	Difficulty string
	// Boss is the final boss card; "" for a challenge without one.
	Boss CardRef
	// Card is the rule card; its abilities work all game.
	Card CardRef
}

// Level returns the level of the difficulty d.
func (c *ChallengeDef) Level(d string) (ChallengeLevel, bool) {
	i := slices.IndexFunc(c.Levels, func(l ChallengeLevel) bool { return l.Difficulty == d })
	if i < 0 {
		return ChallengeLevel{}, false
	}
	return c.Levels[i], true
}

// ChallengeState is the challenge being played (ADR 028, Decision 5);
// Game.Challenge is nil in a regular game. Every field is omitempty, so
// regular saves do not change (ADR 019).
type ChallengeState struct {
	ID         string `json:"id,omitempty"`
	Difficulty string `json:"difficulty,omitempty"`
	// Timer counts the rounds left (R-COOP-02, R-COOP-08).
	Timer int `json:"timer,omitempty"`
	// Card is the rule card's object, outside the game; its abilities
	// work (inPlay).
	Card ObjectID `json:"card,omitempty"`
	// Presents lie outside the game, face down (R-KRAMPUS-05).
	Presents []Present `json:"presents,omitempty"`
	// Holidays: Happy Holidays! has happened (R-KRAMPUS-10).
	Holidays bool `json:"holidays,omitempty"`
	// Lost and Reason say how the game ended: lost, and why.
	Lost   bool `json:"lost,omitempty"`
	Reason Msg  `json:"reason,omitzero"`
	// Starter is the starting player: a round ends before their next
	// regular turn (R-TURN-15, Q-34).
	Starter PlayerID `json:"starter,omitempty"`
	// TurnFrom is, per player, the first object ID made in their latest
	// turn; RoundFrom the same for the turn before it (R-COOP-11, Q-39).
	// 0: no such turn yet.
	TurnFrom  []ObjectID `json:"turn_from,omitempty"`
	RoundFrom []ObjectID `json:"round_from,omitempty"`
	// RoomTurn is the turn whose co-op room rule has run (R-COOP-11).
	RoomTurn int `json:"room_turn,omitempty"`
}

// Present is a present of a challenge: its kind and the face-down
// treasure card outside the game, with its counters (R-KRAMPUS-05).
// Destroyed: it went to the treasure discard (R-KRAMPUS-08); Object is
// then the discarded card. Flipped: it is face up, without counters
// (R-KRAMPUS-10); once a player gains it, Object is that item.
type Present struct {
	Kind      string   `json:"kind,omitempty"`
	Object    ObjectID `json:"object,omitempty"`
	Destroyed bool     `json:"destroyed,omitempty"`
	Flipped   bool     `json:"flipped,omitempty"`
}

// Presents and their counters (ADR 028, Decision 5; R-KRAMPUS-05 to
// R-KRAMPUS-10). Presents are not in play: no card ability sees them;
// only these calls, made by a challenge's own cards, touch them (Q-40).

// AddPresent puts the top treasure card outside the game, face down, as
// the present of kind with n counters (R-KRAMPUS-05). It reports false
// if the treasure deck and its discard are empty.
func (g *Game) AddPresent(kind string, n int) bool {
	c := g.Challenge
	if c == nil {
		return false
	}
	id, ok := g.drawTop(TreasureDeck)
	if !ok {
		return false
	}
	nid := g.move(id, Zone{Kind: ZoneOutside}, NoPlayer)
	g.Object(nid).addCounters("", n)
	c.Presents = append(c.Presents, Present{Kind: kind, Object: nid})
	return true
}

// present returns the present of kind, or nil.
func (g *Game) present(kind string) *Present {
	if g.Challenge == nil {
		return nil
	}
	for i := range g.Challenge.Presents {
		if p := &g.Challenge.Presents[i]; p.Kind == kind {
			return p
		}
	}
	return nil
}

// PresentLive reports whether the present of kind is still there, face
// down: neither destroyed nor flipped (Q-33).
func (g *Game) PresentLive(kind string) bool {
	p := g.present(kind)
	return p != nil && !p.Destroyed && !p.Flipped
}

// PresentOf returns the kind of the present whose card is id, or "".
func (g *Game) PresentOf(id ObjectID) string {
	if g.Challenge == nil {
		return ""
	}
	for _, p := range g.Challenge.Presents {
		if p.Object == id {
			return p.Kind
		}
	}
	return ""
}

// PresentCounters returns the counters on the live present of kind; 0
// if it is gone.
func (g *Game) PresentCounters(kind string) int {
	if !g.PresentLive(kind) {
		return 0
	}
	return g.Object(g.present(kind).Object).CountersOf("")
}

// AddPresentCounters puts n counters on the live present of kind, or
// removes -n; never below 0, and nothing on a present that is gone
// (R-KRAMPUS-09). It returns the change made; a change emits
// present_counters, by p (NoPlayer: the game).
func (g *Game) AddPresentCounters(kind string, n int, p PlayerID) int {
	if !g.PresentLive(kind) || n == 0 {
		return 0
	}
	id := g.present(kind).Object
	o := g.Object(id)
	before := o.CountersOf("")
	o.addCounters("", n)
	d := o.CountersOf("") - before
	if d != 0 {
		g.emit(Event{Kind: EvPresentCounters, Player: p, Object: id, Amount: d, Text: kind})
	}
	return d
}

// DestroyPresent puts the live present of kind into the treasure
// discard, face up (R-KRAMPUS-08). It reports whether it did, and
// whether every present is now destroyed (R-KRAMPUS-09).
func (g *Game) DestroyPresent(kind string) (done, all bool) {
	if !g.PresentLive(kind) {
		return false, false
	}
	pr := g.present(kind)
	old := pr.Object
	pr.Object, pr.Destroyed = g.discard(old, TreasureDeck), true
	g.emit(Event{Kind: EvPresentDestroyed, Player: NoPlayer, Object: pr.Object, Card: g.Object(pr.Object).Card, Prev: old, Text: kind})
	for _, p := range g.Challenge.Presents {
		if !p.Destroyed {
			return true, false
		}
	}
	return true, true
}

// FlipPresents turns every live present face up, without counters
// (R-KRAMPUS-10, Q-40); it marks Holidays, once a game. It reports
// false if it happened already.
func (g *Game) FlipPresents() bool {
	c := g.Challenge
	if c == nil || c.Holidays {
		return false
	}
	c.Holidays = true
	n := 0
	for i := range c.Presents {
		if p := &c.Presents[i]; !p.Destroyed && !p.Flipped {
			p.Flipped = true
			g.Object(p.Object).Counters = nil
			n++
		}
	}
	g.emit(Event{Kind: EvPresentsFlipped, Player: NoPlayer, Amount: n})
	return true
}

// PresentToGive returns the card of the flipped present of kind that
// nobody has gained yet, or 0.
func (g *Game) PresentToGive(kind string) ObjectID {
	p := g.present(kind)
	if p == nil || !p.Flipped || g.Object(p.Object).Zone.Kind != ZoneOutside {
		return 0
	}
	return p.Object
}

// GivePresent makes player to gain the flipped present of kind as an
// item, charged, like any gained item (R-KRAMPUS-10, Q-40).
func (g *Game) GivePresent(kind string, to PlayerID) {
	id := g.PresentToGive(kind)
	if id == 0 {
		return
	}
	nid := g.move(id, Zone{Kind: ZoneInPlay}, to)
	g.enterAsItem(to, nid)
	g.present(kind).Object = nid
	card := g.Object(nid).Card
	g.emit(Event{Kind: EvPresentGiven, Player: to, Object: nid, Card: card, Prev: id, Text: kind})
	g.emit(Event{Kind: EvGainedItem, Player: to, Object: nid, Card: card})
	g.emit(Event{Kind: EvEnteredPlay, Player: to, Object: nid, Card: card})
}

// Forgo records that p forgoes a monster's reward or a soul for n
// counters on the present of kind (R-KRAMPUS-07): the counters go on,
// and forgone tells the table.
func (g *Game) Forgo(kind string, n int, p PlayerID) {
	g.emit(Event{Kind: EvForgone, Player: p, Amount: n, Text: kind})
	g.AddPresentCounters(kind, n, p)
}

// SoulsInPlay is the players' total soul value together, character
// soul values too (R-CARD-21, R-CARD-35, Q-30).
func (g *Game) SoulsInPlay() int {
	n := 0
	for _, pl := range g.Players {
		n += g.SoulValue(pl.ID)
	}
	return n
}

// checkChallenge checks the setup of a challenge (R-COOP-01, R-COOP-02,
// ADR 028 Decision 4).
func checkChallenge(s Setup) error {
	c := s.Challenge
	if s.Players != challengePlayers {
		return fmt.Errorf("setup: challenge %s is for %d players, not %d (R-COOP-01)", c.ID, challengePlayers, s.Players)
	}
	if s.BonusSouls {
		return fmt.Errorf("setup: challenge %s is played without bonus souls (R-COOP-02)", c.ID)
	}
	if !slices.ContainsFunc(s.Sets, func(set CardSet) bool { return set.Code == c.Set }) {
		return fmt.Errorf("setup: challenge %s needs the card set %q", c.ID, c.Set)
	}
	if _, ok := c.Level(s.Difficulty); !ok {
		return fmt.Errorf("setup: challenge %s has no difficulty %q", c.ID, s.Difficulty)
	}
	return nil
}

// startChallenge puts the timer and the rule card out, then runs the
// challenge's own setup (R-COOP-02, R-CHAL-01).
func (g *Game) startChallenge(c *ChallengeDef, difficulty string) error {
	l, _ := c.Level(difficulty)
	if d, ok := g.cards.find(l.Card); !ok || d.Kind != ChallengeCard {
		return fmt.Errorf("setup: challenge %s: %q is no challenge card", c.ID, l.Card)
	}
	if l.Boss != "" {
		if d, ok := g.cards.find(l.Boss); !ok || d.Kind != MonsterCard {
			return fmt.Errorf("setup: challenge %s: %q is no monster", c.ID, l.Boss)
		}
	}
	g.Challenge = &ChallengeState{
		ID: c.ID, Difficulty: difficulty, Timer: timerStart,
		Card:      g.newObject(l.Card, Zone{Kind: ZoneOutside}, NoPlayer),
		TurnFrom:  make([]ObjectID, len(g.Players)),
		RoundFrom: make([]ObjectID, len(g.Players)),
	}
	if c.Setup != nil {
		if err := c.Setup(g, difficulty); err != nil {
			return fmt.Errorf("setup: challenge %s: %w", c.ID, err)
		}
	}
	return nil
}

// StartChallenge makes a table built by BuildTable a challenge at the
// difficulty, for situations and tests: the timer, the rule card, and
// the challenge's own setup. Rounds count from the active player.
func (g *Game) StartChallenge(c *ChallengeDef, difficulty string) error {
	if len(g.Players) != challengePlayers || g.Challenge != nil {
		return fmt.Errorf("setup: challenge %s needs a 2-player table that is no challenge yet", c.ID)
	}
	if _, ok := c.Level(difficulty); !ok {
		return fmt.Errorf("setup: challenge %s has no difficulty %q", c.ID, difficulty)
	}
	if err := g.startChallenge(c, difficulty); err != nil {
		return err
	}
	g.Challenge.Starter = g.Turn.Active
	g.turnTaken(g.Turn.Active)
	g.events = nil
	return nil
}

// WinChallenge ends the challenge: both players win together
// (R-COOP-10, R-CHAL-15).
func (g *Game) WinChallenge() {
	if g.Challenge == nil || g.Over {
		return
	}
	g.Over, g.Winners = true, nil
	for _, pl := range g.Players {
		g.Winners = append(g.Winners, pl.ID)
	}
	g.Waiting = Prompt{Kind: PromptGameOver, Player: NoPlayer}
	g.emit(Event{Kind: EvChallengeWon, Player: NoPlayer})
}

// LoseChallenge ends the challenge: both players lose together, for
// the reason given (R-COOP-09, R-COOP-10, R-CHAL-15).
func (g *Game) LoseChallenge(reason Msg) {
	if g.Challenge == nil || g.Over {
		return
	}
	g.Over, g.Winners = true, nil
	g.Challenge.Lost, g.Challenge.Reason = true, reason
	g.Waiting = Prompt{Kind: PromptGameOver, Player: NoPlayer}
	g.emit(Event{Kind: EvChallengeLost, Player: NoPlayer, Label: reason})
}

// roundEndsBefore: in a challenge, a regular turn of p starts a new
// round, so the round before it ends (R-TURN-15, R-TURN-16, Q-34).
func (g *Game) roundEndsBefore(p PlayerID) bool {
	return g.Challenge != nil && p == g.Challenge.Starter
}

// tickTimer ticks the challenge timer down by 1; at 0 the challenge is
// lost at once (R-COOP-08, R-COOP-09, R-CHAL-13). It reports a loss.
func (g *Game) tickTimer() bool {
	c := g.Challenge
	c.Timer = max(c.Timer-1, 0)
	g.emit(Event{Kind: EvChallengeTimer, Player: NoPlayer, Amount: c.Timer})
	if c.Timer == 0 {
		g.LoseChallenge(ReasonTimer.Msg())
		return true
	}
	return false
}

// turnTaken notes that p's regular turn starts now, or that it was
// skipped, which counts as taken (R-TURN-16), for the co-op room rule.
func (g *Game) turnTaken(p PlayerID) {
	c := g.Challenge
	if c == nil {
		return
	}
	c.RoundFrom[p], c.TurnFrom[p] = c.TurnFrom[p], ObjectID(len(g.Objects))
}

// coopRooms is the co-op room rule after the room step: a room whose
// top card was in its slot when the active player's previous turn
// started goes to discard; the slot refills (R-COOP-11, Q-38, Q-39).
// It runs once per turn and reports whether a room went.
func (g *Game) coopRooms() bool {
	c := g.Challenge
	if c == nil || c.RoomTurn == g.Turn.Number {
		return false
	}
	c.RoomTurn = g.Turn.Number
	from := c.RoundFrom[g.Turn.Active]
	if from == 0 {
		return false // the first round: no previous turn
	}
	gone := false
	for i := range g.Rooms {
		top, ok := g.Rooms[i].TopOf()
		if !ok || top >= from || g.Object(top).Zone.Kind != ZoneInPlay {
			continue
		}
		g.DiscardRoom(top)
		gone = gone || g.Object(top).Zone.Kind != ZoneInPlay
	}
	return gone
}
