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
3. The middle shows the board, in one frame (F-01).
   1. Left: the turn, the stack and the history.
   2. Middle: monster slots above shop slots, then bonus souls.
   3. Right edge: monster, treasure and loot decks, each with its discard.
   4. The board keeps one shape; it scales with the window only (V-23).
   5. Nothing lies over the board (F-02).
4. Each player area shows HP, ATK, cents, souls and hand size.
5. The active player's name is in the accent color.
6. The player's own area reads at a glance (15.17).
   1. Left: name, coins as a big number, and souls.
   2. The character sits centered under the board; items to its right.
   3. Health and attack are big numbers on the character card.
   4. The hand fans out at the bottom edge, partly off screen.
   5. Hover a card to read it in full (V-18).
   6. The coins show the game's own coin (15.38).

## Alternative scenarios and errors

* **2a. Used cards:** a character or item that is not charged is dimmed.
* **2b. Damage:** a red counter on the card shows damage taken.
* **2c. Dead player:** the area fades and says "dead".
* **3a. Empty deck or discard:** a dashed outline keeps its place.
* **3b. Covered slot:** a "+N" badge shows the cards under the top one.
* **3c. More than four slots:** the slots share the row; the board stays the same.
* **1a. Bans or a draft first:** "Setting up the match…" until 7.14.
* **1b. Events play first (V-08):** each update's events play one at a time.
  * Looted and discarded cards fly between deck, hand and discard.
  * Dice rolls show the result; cents, healing and damage pop up.
  * A new turn shows a banner; a revealed card shows large.
  * Then the table shows the update; new cards pop in.
* **1c. Many updates at once:** effects play faster; a long queue skips them.
  * The table always ends on the latest update.
* **1d. Animations off:** updates show at once.
* **1e. Dice look (V-19):** a roll shows the face of the player's dice texture.
  * No texture (public clone) or a value past 1–6: a plain die with the number.
* **1f. Card backs (V-26):** each deck shows its own back.
  * An opponent's hand is one loot back with its card count.
  * No back images (public clone): plain dark backs with the deck name.
* **1g. No coin image (public clone):** a drawn ¢ coin by the cents.

## Postconditions

* The table fits 1280×720 and larger windows.
* Planned: the stack panel (7.8), history and turn panels (7.9).
* No screen scrolls as a whole, at any window size (V-10).
  * Long lists and panels scroll on their own.
* A mouse drag selects nothing on the table or in menus.
