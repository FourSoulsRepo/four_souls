package rulesengine

// Barrier reports whether an event stops a Soft rollback: a step with
// it cannot be undone, nor any step before it (ADR 020, item 4). Dice,
// damage and death, hidden information seen, the end of a turn, and
// the game end are barriers. An event only one player sees (Private)
// is always one.
func Barrier(e Event) bool {
	return e.Private || barrierKinds[e.Kind]
}

// barrierKinds sorts every event kind: true is a barrier. A new kind
// must be added here; TestBarrierKinds fails until it is sorted.
var barrierKinds = map[EventKind]bool{
	// Dice (item 4.2): a roll and what follows from it.
	EvDiceRolled:         true,
	EvRollChanged:        true,
	EvRandomPick:         true,
	EvRollWouldResolve:   true,
	EvRollResolved:       true,
	EvRollRepeated:       true,
	EvAttackRollRepeated: true,
	EvFirstRollTied:      true,
	EvFirstPlayer:        true,
	// Damage (item 4.3); a death follows damage, and a death without
	// damage ("kill a player") is one too.
	EvDamaged:    true,
	EvDied:       true,
	EvObjectDied: true,
	// Hidden information seen (item 4.4): a card from a deck comes
	// into view, or one player looks at a card.
	EvCardRevealed:   true,
	EvLookedAt:       true,
	EvGainedTreasure: true, // the treasure deck's top, or a searched deck
	EvCharacterDealt: true,
	// The turn ends (item 4.5).
	EvTurnEnded:      true,
	EvTurnEndedEarly: true,
	EvTurnSkipped:    true,
	// The game (item 4.6).
	EvGameStarted: true,
	EvGameWon:     true,

	// Not barriers: what is seen anyway, and steps on the way.
	EvActivated:        false,
	EvAttackDeclared:   false,
	EvAttackEnded:      false,
	EvAttackTarget:     false,
	EvBoosted:          false,
	EvCancelled:        false,
	EvCounters:         false,
	EvCursed:           false,
	EvDamagePending:    false, // on the stack; not dealt yet
	EvDeactivated:      false,
	EvDeathPending:     false, // on the stack; not dead yet
	EvDestroyed:        false,
	EvDiscarded:        false,
	EvEndOfTurn:        false, // end-of-turn triggers; the turn goes on
	EvEnteredPlay:      false,
	EvExpanded:         false, // the new slots' cards are card_revealed
	EvExtraTurn:        false,
	EvAbilityFizzle:    false,
	EvFlipped:          false,
	EvGainedCents:      false,
	EvGainedItem:       false, // from the shop, or with gained_treasure
	EvGainedSoul:       false,
	EvGaveCard:         false, // private when the card is hidden
	EvGaveCents:        false,
	EvHealed:           false,
	EvLootDone:         false,
	EvLootPlayed:       false, // the player shows their own card
	EvLooted:           false, // private when drawn; else from the discard
	EvLostCents:        false,
	EvMovedToDeck:      false, // private when hidden
	EvPaid:             false,
	EvPassed:           false,
	EvPenaltyPaid:      false,
	EvPhase:            false,
	EvPrevented:        false,
	EvPurchaseDeclared: false,
	EvPurchaseFailed:   false,
	EvPurchased:        false, // the refilled slot is card_revealed
	EvRecharged:        false,
	EvReplaced:         false,
	EvRewardsGained:    false,
	EvShielded:         false,
	EvStackAdded:       false,
	EvStackResolved:    false,
	EvStartOfTurn:      false,
	EvStole:            false,
	EvTurnStarted:      false,
}
