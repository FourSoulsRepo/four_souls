package rulesengine

// EventKind names something that happened.
type EventKind string

// Events emitted so far; the list grows with the engine.
const (
	EvGameStarted    EventKind = "game_started"
	EvFirstPlayer    EventKind = "first_player"
	EvDiceRolled     EventKind = "dice_rolled"
	EvTurnStarted    EventKind = "turn_started"
	EvPhase          EventKind = "phase"
	EvRecharged      EventKind = "recharged"
	EvLooted         EventKind = "looted"
	EvDiscarded      EventKind = "discarded"
	EvGainedCents    EventKind = "gained_cents"
	EvHealed         EventKind = "healed"
	EvPassed         EventKind = "passed_priority"
	EvTurnEnded      EventKind = "turn_ended"
	EvGameWon        EventKind = "game_won"
	EvCardRevealed   EventKind = "card_revealed"
	EvCharacterDealt EventKind = "character_dealt"
)

// Event is one thing that happened, in order. Hidden cards (e.g. a looted
// card) are filtered per viewer by the view filter (step 4.9).
type Event struct {
	Kind   EventKind `json:"kind"`
	Player PlayerID  `json:"player"`
	Object ObjectID  `json:"object,omitempty"`
	Card   CardRef   `json:"card,omitempty"`
	Amount int       `json:"amount,omitempty"`
	Text   string    `json:"text,omitempty"`
	// Private means only Player sees Card (e.g. a looted card).
	Private bool `json:"private,omitempty"`
	// Prev is the object an event moved away, when it got a new ID.
	Prev ObjectID `json:"prev,omitempty"`
}
