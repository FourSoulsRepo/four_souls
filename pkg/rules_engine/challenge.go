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
// treasure card, with its counters (R-KRAMPUS-05, 14.22).
type Present struct {
	Kind   string   `json:"kind,omitempty"`
	Object ObjectID `json:"object,omitempty"`
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
