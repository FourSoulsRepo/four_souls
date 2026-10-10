# UC-07: Reconnect after a lost connection

**Module:** `internal/server`
**Status:** In progress
**Actors:** A player who lost the connection; the other players
**Goal:** The player comes back to their seat, or the others go on without them
**Preconditions:** A game is running (UC-04)

## Main scenario (happy path)

1. A player's connection drops.
2. The game pauses for everyone; the others see who is away (N-08).
3. The player's app connects again and says hello with its token.
4. The server gives the seat back and sends the full current view.
5. The pause ends; the game goes on where it stopped.

## Alternative scenarios and errors

* **2a. The others vote:** each connected player votes to wait longer or to kick; a vote can be changed.
* **2b. More than half vote to kick:** the away players are kicked; their seats stay empty, never a bot.
* **2c. A kicked seat's turn:** the server answers its prompts the minimal way (pass or end the turn, first cards, first option).
* **3a. A kicked player comes back:** the server refuses the seat and closes the connection.
* **3b. The connection drops at a table before the start:** the seat is freed (UC-04).

## Postconditions

* The player plays on from the same seat, or the game goes on without them.
