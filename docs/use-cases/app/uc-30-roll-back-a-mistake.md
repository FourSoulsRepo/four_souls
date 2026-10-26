# UC-30: Roll back a mistake (Godmode)

**Module:** `internal/server`, `pkg/record`, `frontend`
**Status:** Implemented
**Actors:** Any player of the game
**Goal:** A misclick is undone; the game goes on from before it (N-12)
**Preconditions:** The game was created with the rollback mode Godmode (ADR 020)

## Main scenario (happy path)

1. A player makes a mistake, e.g. attacks the wrong monster.
2. A player opens the History panel.
   1. Its title shows "Godmode".
   2. The first row of each step has a ↶ button.
3. The player presses ↶ on the row of the mistake.
   1. Any player may: the author too.
   2. No vote; it happens at once.
4. The server rebuilds the game as it was before that step.
   1. The rows of that step and of every step after it turn grey.
   2. A row says who rolled back: "Bo rolls back 3 steps".
5. Every table jumps to the restored view, with no animation.
6. The game goes on from there.
   1. A roll made again gives the same value (N-12, item 7).
   2. "Wait for my turn" is off for everyone.

## Alternative scenarios and errors

* **2a. Disabled (the default):** no buttons; the server refuses a rollback.
* **3a. An undone row:** it has no button; the server refuses it.
* **3b. A step from before a load or a server restart:** no button; refused.
* **3c. The game is over:** no buttons.
* **3d. The game is paused for a missing player:** the rollback is refused.
* **4a. The replay differs from the game:** the rollback is refused; play goes on.
* **6a. Soft mode, with a vote and barriers:** planned in 17.4.

## Postconditions

* The record keeps every step and a `rollback` line (ADR 007).
* A replay of the record leaves the undone steps out.
* Step numbers never repeat.
