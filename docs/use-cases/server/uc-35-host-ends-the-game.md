# UC-35: The host ends the game

**Module:** `internal/server`, `frontend`
**Status:** Implemented
**Actors:** Host; the other players
**Goal:** The host stops a game nobody will finish, on any server
**Preconditions:** A game is running, not over yet (UC-03, UC-19)

## Main scenario (happy path)

1. The host opens **Menu** in the game, then **End game**.
2. The menu asks "End the game?"; the host presses **End game**.
3. The server stops the game; its record ends unfinished.
4. Everyone lands on the server's game list.
5. The others see "The game ended: The host ended the game."
6. The game is gone from the list.

## Alternative scenarios and errors

* **1a. A player who is not the host:** the menu has no **End game**.
  * Sent anyway, `end` is refused with `not_host`.
* **2a. The host presses Back:** the menu shows again; the game goes on.
* **3a. A player is away:** they get the notice when they come back.
  * Within 30 minutes; after that they just see the game list.
* **1b. On the app's own server:** the server keeps running.
  * The host stays on its game list and may start a new game.

## Postconditions

* The game cannot go on; unlike a save, nothing is kept to load.
* Its record is kept.
