# UC-29: A game ends by an idle limit

**Module:** `internal/server`, `frontend`
**Status:** Implemented
**Actors:** The server; the players of a stuck or endless game
**Goal:** Abandoned and endless games do not live on the server forever
**Preconditions:** A server runs a game, or keeps a table over a restart (UC-28)

## Main scenario (happy path)

1. Every player of a game leaves, or the server restarts and nobody comes back.
2. 30 minutes pass with no player connected.
3. The server ends the game; its record ends unfinished.
4. A player who comes back within 30 minutes more gets a notice.
5. The app shows "The game ended: nobody played it for too long" and the game list.

## Alternative scenarios and errors

* **1a. A game runs for 24 hours:** it ends even while players play.
  * Connected players see "It ran for 24 hours, the longest a game may run."
* **2a. Server downtime:** it counts toward the 30 minutes.
* **1b. A table kept over a restart:** it closes after 30 minutes with nobody at it.
* **1c. The operator changes the limits:** `-empty_minutes` and `-game_hours`; 0 turns one off.

## Postconditions

* The game is gone from the server's list; its record is kept.
