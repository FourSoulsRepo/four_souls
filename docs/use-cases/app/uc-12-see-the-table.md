# UC-12: See the table

**Module:** `frontend/src/screens/GameScreen.tsx`, `frontend/src/components`
**Status:** In progress
**Actors:** Player
**Goal:** The player sees the whole game state at a glance (V-01, V-02)
**Preconditions:** The game started (UC-03, UC-04)

## Main scenario (happy path)

1. The server sends an update after every step.
2. The app draws the table from the update's view.
   1. The player's own area is at the bottom: character, items, hand.
   2. The others sit in turn order, clockwise.
   3. Two players: the other one on top.
   4. Three players: left and right.
   5. Four players: left, top, right.
3. The middle shows the board.
   1. Monster deck, monster slots, monster discard, bonus souls.
   2. Treasure deck, shop slots, treasure discard.
   3. Loot deck and loot discard on the side.
4. Each player area shows HP, ATK, cents, souls and hand size.
5. The active player's area has a red frame.

## Alternative scenarios and errors

* **2a. Used cards:** a character or item that is not charged is dimmed.
* **2b. Damage:** a red counter on the card shows damage taken.
* **2c. Dead player:** the area fades and says "dead".
* **3a. Empty deck or discard:** a dashed outline keeps its place.
* **3b. Covered slot:** a "+N" badge shows the cards under the top one.
* **1a. Bans or a draft first:** "Setting up the match…" until 7.14.
* **1b. Events play first (V-08):** each update's events play one at a time.
  * Looted and discarded cards fly between deck, hand and discard.
  * Dice rolls show the result; cents, healing and damage pop up.
  * A new turn shows a banner; a revealed card shows large.
  * Then the table shows the update; new cards pop in.
* **1c. Many updates at once:** effects play faster; a long queue skips them.
  * The table always ends on the latest update.
* **1d. Animations off:** updates show at once.

## Postconditions

* The table fits 1280×720 and larger windows.
* Planned: the stack panel (7.8), history and turn panels (7.9).
