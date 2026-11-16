# UC-46: See what the turn owes

**Module:** `frontend/src/prompt.ts`, `frontend/src/components/PlayerArea.tsx`, `frontend/src/components/Board.tsx`
**Status:** Implemented
**Actors:** Player, spectator
**Goal:** Everyone sees which attacks a card forces, and who joins the next one
**Preconditions:** The game runs (UC-12); it is Ann's turn

## Main scenario (happy path)

1. Ann plays Monster Manual on Gaper.
2. Gaper gets a dashed red frame for everyone (19.180).
   1. The fight's frame, dashed; solid once the fight starts.
   2. Its zoom says "Must be attacked this turn".
3. Ann's banner says "You must attack: Gaper" under Your turn.
4. Ann attacks Gaper; the frame and the line go.

## Alternative scenarios and errors

* **3a. A count is owed (19.181):** the banner adds the attacks left.
  * Lust for Blood: "You must attack: Bo (1 left)".
  * Envy names no target: "You must attack (1 left)".
* **2a. Attacks on the monster deck are owed (19.182):** the deck gets the frame.
  * Ambush: "You must attack: the monster deck (2 left)".
  * The deck comes first, as the engine asks for it first.
* **2b. Some attacks are only on the monster deck (19.183):** the count says so.
  * Big Spider: "Attacks 2 (1 deck only)".
  * On Ann's seat for everyone; by her loot count for Ann.
* **3b. The marked monster leaves the table:** the frame and the line go.
* **3c. Others' banners:** they see Ann's turn; the marks are for everyone.

## Postconditions

* Nobody is surprised that Ann cannot end her turn yet.
