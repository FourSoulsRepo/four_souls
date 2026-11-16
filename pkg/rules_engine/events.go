package rulesengine

// EventKind names something that happened.
type EventKind string

// Events emitted so far; the list grows with the engine.
const (
	EvGameStarted EventKind = "game_started"
	EvFirstPlayer EventKind = "first_player"
	EvDiceRolled  EventKind = "dice_rolled"
	EvTurnStarted EventKind = "turn_started"
	EvPhase       EventKind = "phase"
	EvRecharged   EventKind = "recharged"
	EvLooted      EventKind = "looted"
	// EvLootDone: Player looted Amount cards in one go ("each time a
	// player loots": Laser Eye).
	EvLootDone EventKind = "loot_done"
	// EvGaveCents: Player gave Amount¢ to player By-1 (Social Goals).
	EvGaveCents        EventKind = "gave_cents"
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
	EvRewardsGained    EventKind = "rewards_gained" // the active player got Object's rewards
	EvTurnSkipped      EventKind = "turn_skipped"
	EvCursed           EventKind = "cursed"             // Player gained the curse Object
	EvExpanded         EventKind = "expanded"           // Amount slots of kind Text were added
	EvRandomPick       EventKind = "random_pick"        // Text is the picked option
	EvRollWouldResolve EventKind = "roll_would_resolve" // a roll of Amount tries to resolve (R-DICE-05)
	// EvFirstRollTied: the lowest first-player roll, Amount, tied; Text
	// lists the tied players, who roll again, as "0,2" (15.43).
	EvFirstRollTied EventKind = "first_roll_tied"
	// EvTopChoice: Player looked at a deck's top card and kept it there
	// or put it on the bottom. Text is "<deck> deck top" or "<deck> deck
	// bottom"; Private, so only Player sees Card (19.90).
	EvTopChoice EventKind = "top_choice"
	// EvDeckBottom: Amount cards that were in a deck went to its bottom
	// (ToBottomRandom). Text is "<deck> deck"; public, no card, Player
	// NoPlayer: the board shows them go under the deck (19.147).
	EvDeckBottom EventKind = "deck_bottom"
	// EvChallengeTimer: the challenge timer ticked down to Amount
	// (R-COOP-08).
	EvChallengeTimer EventKind = "challenge_timer"
	// EvChallengeWon and EvChallengeLost end a challenge; a loss has its
	// reason as Label (R-COOP-10, ADR 028).
	EvChallengeWon  EventKind = "challenge_won"
	EvChallengeLost EventKind = "challenge_lost"
	// EvPresentCounters: the present of kind Text got Amount counters,
	// negative if removed; Player is who forwent, or NoPlayer
	// (R-KRAMPUS-07, R-KRAMPUS-12).
	EvPresentCounters EventKind = "present_counters"
	// EvPresentDestroyed: the present of kind Text went to the treasure
	// discard as Card (R-KRAMPUS-08).
	EvPresentDestroyed EventKind = "present_destroyed"
	// EvPresentsFlipped: Amount presents flipped face up (Happy
	// Holidays!, R-KRAMPUS-10).
	EvPresentsFlipped EventKind = "presents_flipped"
	// EvPresentGiven: Player gained the present of kind Text, Card, as
	// an item (R-KRAMPUS-10).
	EvPresentGiven EventKind = "present_given"
	// EvForgone: Player forwent a reward or a soul for Amount counters on
	// the present of kind Text (R-KRAMPUS-07).
	EvForgone EventKind = "forgone"
	// EvMinionArrived: Card came from the monster deck into a new minion
	// slot (R-CHAL-09).
	EvMinionArrived EventKind = "minion_arrived"
	// EvDeckOrder: the cards of a deck changed places out of sight: it
	// was shuffled (a search, a discard shuffled in), its top cards put
	// back in an order, or a card moved within it. Text is "<deck>
	// deck"; public, no card. Who remembers where its cards lie forgets
	// them (ADR 030, Decision 19.8).
	EvDeckOrder EventKind = "deck_order"
)

// Event is one thing that happened, in order. Hidden cards (e.g. a looted
// card) are filtered per viewer by the view filter (step 4.9).
type Event struct {
	Kind   EventKind `json:"kind"`
	Player PlayerID  `json:"player"`
	Object ObjectID  `json:"object,omitempty"`
	Card   CardRef   `json:"card,omitempty"`
	Amount int       `json:"amount,omitempty"`
	// Text is a code for the client and for card triggers, never shown
	// as it is: "combat", "attack", "starting item" (ADR 021).
	Text string `json:"text,omitempty"`
	// Label is what a stack item is, on its events and on a roll's
	// dice_rolled; the picked option of a random_pick (ADR 021).
	Label Msg `json:"label,omitzero"`
	// Private means only Player sees Card (e.g. a looted card), and
	// whoever SeenBy names (the giver of a hand card, V-14).
	Private bool       `json:"private,omitempty"`
	SeenBy  []PlayerID `json:"seen_by,omitempty"`
	// Prev is the object an event moved away, when it got a new ID.
	Prev ObjectID `json:"prev,omitempty"`
	// Source is the object that caused the event, when that matters
	// (the shield that prevented damage).
	Source ObjectID `json:"source,omitempty"`
	// StackID is the stack item the event is about (a roll).
	StackID int `json:"stack_id,omitempty"`
	// Dice are all the dice of a roll of several dice, the first is
	// Amount (R-DICE-11, 19.39).
	Dice []int `json:"dice,omitempty"`
	// By is the player who dealt damage, plus 1 (0: a monster or nobody).
	// On a death it is the killer, for the history (15.88). On a look
	// it is the owner of the hand the card is in; 0 for a deck card
	// (19.71).
	By int `json:"by,omitempty"`
	// Item is the stack item that resolved, with its used ability's
	// words, for the history strip (15.22, S-10).
	Item *StackItem `json:"item,omitempty"`
	// Sources are what changed a resolved roll, for the history (15.94).
	Sources []StatSourceView `json:"sources,omitempty"`
	// Deck places a deck card a look shows: its deck and its depth, 0
	// the top. A seat's memory reads it (ADR 030, Decision 19.11); it is
	// hidden with Card from those who do not see the card.
	Deck *DeckPlace `json:"deck,omitempty"`
}

// DeckPlace is where a card lies in a deck: Depth cards above it.
type DeckPlace struct {
	Deck  DeckKind `json:"deck"`
	Depth int      `json:"depth"`
}

// Texts of events that name a deck or a discard. The engine writes them
// and a seat's memory (internal/world) reads them, so both use these:
// the client's history and its board read the same words.
const (
	// TextFromDiscard is the Text of EvLooted for a card taken from a
	// discard pile into a hand, in sight of all.
	TextFromDiscard = "from the discard"
	// DeckTextTop and DeckTextBottom follow "<deck> deck" in an
	// event's Text for a card on the deck's top or bottom (EvTopChoice,
	// and EvMovedToDeck to the bottom).
	DeckTextTop    = " top"
	DeckTextBottom = " bottom"
)

// DeckText is how an event names deck d: "<deck> deck", then where (""
// for the deck itself, DeckTextTop or DeckTextBottom).
func DeckText(d DeckKind, where string) string {
	return d.String() + " deck" + where
}

// DeckOfText reads a Text DeckText wrote with where; false for any
// other text, a deck's text with another where too.
func DeckOfText(text, where string) (DeckKind, bool) {
	for d := range deckCount {
		if text == DeckText(d, where) {
			return d, true
		}
	}
	return 0, false
}
