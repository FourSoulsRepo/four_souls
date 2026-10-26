# UC-31: Undo a step by vote (Soft)

**Module:** `internal/server`, `rules_engine`, `frontend`
**Status:** Implemented
**Actors:** The step's author, the other players
**Goal:** A misclick is undone when the other players agree (N-12)
**Preconditions:** The game was created with the rollback mode Soft (ADR 020)

## Main scenario (happy path)

1. A player makes a mistake, e.g. attacks the wrong monster.
2. A player opens the History panel.
   1. Its title shows "Soft rollback".
   2. One Undo button under the title targets the latest move (17.7).
   3. Barriers: dice, damage and death, hidden cards seen, a turn's end.
   4. A lone pass or throw is no move.
3. Another player presses Undo while the mistake is the latest move.
   1. That opens a vote; their own vote counts as yes.
4. The game waits; every table shows the vote.
   1. It names who asks, the move, and how many steps it undoes.
   2. Voters: the connected players other than the author.
   3. Each voter answers "Yes, undo" or "No"; they may change it.
5. More than half of the voters say yes.
   1. 2 players: the opponent; 3: both others; 4: two of three.
6. The game rolls back as in UC-30, steps 4 to 6.

## Alternative scenarios and errors

* **2a. The author's own move:** the button reads "Ask to undo".
* **3a. The author asks:** the vote opens with no yes yet.
* **3b. The author sends a rollback anyway:** the server refuses it.
* **3c. 2 players:** the opponent is the only voter; the rollback happens at once.
* **3d. No voter is connected:** the server refuses the ask.
* **4a. A player acts during the vote:** the server refuses the intent.
* **4b. The opener cancels:** the vote ends; a row says so.
* **4c. The response timer is on:** the vote gets the same seconds.
* **4d. Its time runs out:** the vote is declined.
* **4e. A player leaves:** the game pauses and the vote ends.
* **5a. The majority is out of reach:** the vote is declined at once.
* **5b. A declined step:** it may be asked again.
* **6a. A barrier came after the move:** Undo is greyed and names it.
  * "Nothing to undo since a die was rolled"; the server refuses it.
* **3e. A played loot card or a used item (17.6, 17.7):** Undo targets the play.
  * "Undo: Bo plays A Penny!"; also while the card is on the stack.
  * Its result came with a pass, so the play stays the latest move.
  * It stays undoable until a new card comes from a deck.
  * Or a die is rolled or changed, or damage is dealt.
  * A turn end, a death, and a seen hand card block it too.
  * A card put from a hand on the loot deck does not (ADR 020, item 4.22).
  * A purchase keeps its barriers.
* **4f. While a vote is open:** Undo is greyed, "A rollback vote is open".
* **3f. A played card that asks a target:** the play and the pick are one move.
  * One press puts the card back in the hand, not at its target prompt.
* **3g. A move with no history line (17.8):** Undo says "Undo: Bo's choice".
  * The author's button says "Ask to undo: your choice".
  * The vote names it the same way.

## Postconditions

* The record keeps a `rollback` line only for a vote that passed.
* A vote that ended with no rollback leaves a history row.
