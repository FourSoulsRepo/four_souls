# UC-15: Read the history

**Module:** `frontend/src/history.ts`, `frontend/src/components/InfoPanels.tsx`
**Status:** Implemented
**Actors:** Player
**Goal:** The player never gets lost (V-03, V-05, V-06)
**Preconditions:** The game runs (UC-12)

## Main scenario (happy path)

1. The history panel sits on the left of the board.
2. Each resolved stack item adds one entry, with an icon.
3. Each turn starts with a "Turn N: <name>" line.
4. The newest entry is at the bottom; the list scrolls.
5. Hovering an entry with a card zooms that card.
6. The turn badge shows the turn number, the active player and the step.
7. The active player's area shows the loot plays, attacks and buys left.

## Alternative scenarios and errors

* **2a. Hidden cards:** another player's loot reads "Bo loots a card".
  * The viewer's own loot names the card.
* **1a. Hide:** the Hide button folds the panel into a History button.
* **5a. Touch:** tapping an entry zooms its card; the next tap closes it.
* **7a. Counters:** a card with counters shows their number.
* **4a. Copy:** history text can be selected and copied (V-10).

## Postconditions

* The history holds what this player saw since joining the table.
* Planned: the whole match after a reconnect, with replays (step 10).
