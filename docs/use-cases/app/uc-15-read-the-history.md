# UC-15: Read the history

**Module:** `frontend/src/history.ts`, `frontend/src/components/InfoPanels.tsx`, `frontend/src/components/HistoryStrip.tsx`
**Status:** Implemented
**Actors:** Player
**Goal:** The player never gets lost (V-03, V-05, V-06)
**Preconditions:** The game runs (UC-12)

## Main scenario (happy path)

1. The history panel sits in the board's info column, under the stack.
2. Each resolved stack item adds one entry, with an icon.
3. Each turn starts with a "Turn N: <name>" line.
4. The newest entry is at the bottom; the list scrolls.
5. Hovering an entry with a card zooms that card.
6. The turn badge shows the turn number, the active player and the step.
7. The active player's area shows the loot plays, attacks and buys left.
8. The strip shows resolved stack items as icons, newest on top.
   1. The icons are the stack panel's: card art, a die, a blood drop, a skull.
   2. A sword badge marks an attack; a bolt marks a played loot card.
   3. A line marks the first item of each turn.
   4. Hovering an icon shows the card large, the used ability, and the owner.
9. The strip runs along the board's left edge (15.22).
   1. Its History button opens the full history panel.
   2. The panel starts hidden, so the stack has the room.

## Alternative scenarios and errors

* **2a. Hidden cards:** another player's loot reads "Bo loots a card".
  * The viewer's own loot names the card.
* **1a. Hide:** the Hide button, or History again, closes the panel.
* **5a. Touch:** tapping an entry zooms its card; the next tap closes it.
* **7a. Counters:** a card with counters shows their number.
* **4a. Copy:** history text can be selected and copied (V-10).
* **2b. A card changes hands (V-14):** a steal or a give adds a line.
  * The taker reads "You get Bomb from Bo".
  * The giver reads "Ann gets your Bomb".
  * Everyone else reads "Ann gets a card from Bo".

## Postconditions

* The history holds what this player saw since joining the table.
* Planned: the whole match after a reconnect, with replays (step 10).
