# UC-07: Reconnect after a lost connection

**Module:** `internal/server`, `frontend`
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
* **1a. In the app:** a "Connection lost" dialog keeps the table on screen.
  * **Reconnect** says hello again with the seat's token.
  * **Leave** goes back to the main menu.
* **2d. The others, in the app:** a "Waiting for <name>" dialog shows the votes.
  * **Wait** or **Go on without them**; a vote can be changed.
* **1b. The app was closed:** the app keeps the last game's server and seat token in its settings.
  * Joining that server again sends the token; the player is back in their seat.
  * The welcome names the seat, so the app knows which player it is.
* **2e. A kicked seat rolls:** the server throws its die at once (15.60).
* **3c. The game is over:** the token takes the player to the lobby, not back (19.32).
  * The player keeps the token; the record download still works.
  * A new game with that token is the one a later hello comes back to.
* **1c. The app reconnects by itself (N-20):** at once, then after 1, 2, 4, 8, then every 10 seconds.
  * A small "Reconnecting…" badge shows; the table stays in view.
  * The "Connection lost" dialog replaces it after 30 seconds.
  * A failed retry's error shows only in that dialog.

## Postconditions

* The player plays on from the same seat, or the game goes on without them.
