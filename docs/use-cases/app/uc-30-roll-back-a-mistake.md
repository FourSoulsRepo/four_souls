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
   2. Under the title sits one Undo button (17.7).
   3. Its tooltip names the latest move: "Undo: Bo attacks Pin".
3. The player presses Undo, once per move, until the mistake is undone.
   1. Any player may: the author too.
   2. No vote; it happens at once.
4. The server rebuilds the game as it was before that move.
   1. The rows of that move and of every move after it turn grey.
   2. A row says who rolled back: "Bo rolls back 3 steps".
5. Every table jumps to the restored view, with no animation.
6. The game goes on from there.
   1. A roll made again gives the same value (N-12, item 7).
   2. "Wait for my turn" is off for everyone.

## Alternative scenarios and errors

* **2a. Disabled (the default):** no Undo button; the server refuses a rollback.
* **3a. Undone moves:** the next press undoes the move before them.
* **3b. Nothing since a load:** Undo is greyed, "Nothing to undo yet".
* **3c. The game is over:** Undo is greyed, "The game is over".
* **3d. The game is paused for a missing player:** Undo is greyed; the server refuses it.
* **4a. The replay differs from the game:** the rollback is refused; play goes on.
* **6a. Soft mode, with a vote and barriers:** UC-31.
* **3e. A dedicated server restarts:** the history stays if it replays to the kept game.
* **3f. The replay after a restart differs:** the history ends there.
  * A row says "The server restarted; no rollback before this".
  * Undo greys there: "Nothing to undo since the server restarted".
* **3g. A played card (17.7, 19.69):** one press undoes the play.
  * On the stack or resolved; the tooltip says "Undo: Bo plays A Penny!".
  * The card goes back to the hand, or the item is ready again.
  * A mode or target the play asked first is part of the move.
* **3h. A spectator:** Undo is greyed, "Spectators cannot undo".
* **3i. A greyed Undo is clicked:** a note at the pointer says why.
* **3j. A choice window covers the table:** "View board" first, then Undo.
* **3k. An attack:** its declaration and its target are two moves.

## Postconditions

* The record keeps every step and a `rollback` line (ADR 007).
* A replay of the record leaves the undone steps out.
* Step numbers never repeat.
