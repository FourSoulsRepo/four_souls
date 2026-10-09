package rulesengine

// EventKind names something that happened.
type EventKind string

// Events emitted so far; the list grows with the engine.
const (
	EvGameStarted      EventKind = "game_started"
	EvFirstPlayer      EventKind = "first_player"
	EvDiceRolled       EventKind = "dice_rolled"
	EvTurnStarted      EventKind = "turn_started"
	EvPhase            EventKind = "phase"
	EvRecharged        EventKind = "recharged"
	EvLooted           EventKind = "looted"
	EvDiscarded        EventKind = "discarded"
	EvGainedCents      EventKind = "gained_cents"
	EvHealed           EventKind = "healed"
	EvPassed           EventKind = "passed_priority"
	EvTurnEnded        EventKind = "turn_ended"
	EvGameWon          EventKind = "game_won"
	EvCardRevealed     EventKind = "card_revealed"
	EvCharacterDealt   EventKind = "character_dealt"
	EvCounters         EventKind = "counters"     // Amount added, negative if removed
	EvRollChanged      EventKind = "roll_changed" // Amount is the new result
	EvBoosted          EventKind = "boosted"      // till end of turn; Text names the stat
	EvShielded         EventKind = "shielded"     // the next damage will be prevented
	EvPrevented        EventKind = "prevented"    // damage was prevented
	EvStole            EventKind = "stole"        // Player took Amount¢ from Text's player
	EvLookedAt         EventKind = "looked_at"    // private: Player saw Card
	EvGaveCard         EventKind = "gave_card"    // a hand card changed hands
	EvMovedToDeck      EventKind = "moved_to_deck"
	EvDamagePending    EventKind = "damage_pending" // damage went on the stack, aimed at Player or Object
	EvDeathPending     EventKind = "death_pending"  // a death went on the stack
	EvEnteredPlay      EventKind = "entered_play"
	EvTurnEndedEarly   EventKind = "turn_ended_early"
	EvPenaltyPaid      EventKind = "penalty_paid"
	EvExtraTurn        EventKind = "extra_turn"
	EvCursed           EventKind = "cursed"             // Player gained the curse Object
	EvExpanded         EventKind = "expanded"           // Amount slots of kind Text were added
	EvRandomPick       EventKind = "random_pick"        // Text is the picked option
	EvRollWouldResolve EventKind = "roll_would_resolve" // a roll of Amount tries to resolve (R-DICE-05)
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
	// Source is the object that caused the event, when that matters
	// (the shield that prevented damage).
	Source ObjectID `json:"source,omitempty"`
	// StackID is the stack item the event is about (a roll).
	StackID int `json:"stack_id,omitempty"`
}
