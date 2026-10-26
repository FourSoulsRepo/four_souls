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
   2. Rows after the last barrier offer a button.
   3. Barriers: dice, damage and death, hidden cards seen, a turn's end.
   4. A lone pass or throw has no button.
3. Another player presses ↶ on the row of the mistake.
   1. That opens a vote; their own vote counts as yes.
4. The game waits; every table shows the vote.
   1. It names who asks, the step, and how many steps it undoes.
   2. Voters: the connected players other than the author.
   3. Each voter answers "Yes, undo" or "No"; they may change it.
5. More than half of the voters say yes.
   1. 2 players: the opponent; 3: both others; 4: two of three.
6. The game rolls back as in UC-30, steps 4 to 6.

## Alternative scenarios and errors

* **2a. The author's own row:** it offers "Ask ↶" instead of ↶.
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
* **6a. A barrier came after the step:** no button; the server refuses it.

## Postconditions

* The record keeps a `rollback` line only for a vote that passed.
* A vote that ended with no rollback leaves a history row.
