# UC-28: Play on after a server update

**Module:** `internal/server`, `rules_engine`, `frontend`
**Status:** Implemented
**Actors:** The server's operator; the players of its games
**Goal:** An update or restart of the dedicated server does not end the games
**Preconditions:** A dedicated server runs games (UC-02)

## Main scenario (happy path)

1. The operator stops the server: systemd, Docker, or Ctrl+C.
2. Every player gets a "restarting" message.
3. The server writes each table and game to a restart file.
4. The operator starts the new version on the same saves folder.
5. The server loads the restart files before it takes players.
6. The players reconnect with their tokens (UC-07).
7. The game goes on at the same step once every seat is back.

## Alternative scenarios and errors

* **4a. The new version has another rules engine:**
  * The game goes on if the engine reads it back exactly.
  * Each player is told once: cards may now behave differently.
* **4b. The new engine cannot read the game:**
  * The game ends; its record ends unfinished.
  * Each player who comes back within 30 minutes gets an apology.
  * The file stays in `saves/restart/` with `.lost` added.
* **4c. A lobby table:** it comes back with its seats and password.
* **4d. Bans and picks were under way:** they start over once every player is back.
* **4e. The game was over:** it is not kept.
* **1a. An app's own host stops:** its games end as before; nothing is kept.
* **2a. In the app:** "The server is restarting. Reconnecting…" shows over the table.
  * The app reconnects by itself (UC-07 1c); no dialog in a short restart.
* **6a. The first player back:** the pause shows, without a kick vote.
  * The kick vote opens 60 seconds after the server started.
* **4f. The notices in the app:** a dialog with OK, on any screen.
  * After the apology the app shows the server's game list.

## Postconditions

* Every game goes on, or its players know why it ended.
