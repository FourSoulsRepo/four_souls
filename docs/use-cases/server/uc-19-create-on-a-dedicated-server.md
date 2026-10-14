# UC-19: Create a game on a dedicated server

**Module:** `internal/server`, `frontend`
**Status:** Implemented
**Actors:** Player
**Goal:** A game waits on a dedicated server for friends to join
**Preconditions:** A dedicated server runs (UC-02); the player knows its address

## Main scenario (happy path)

1. The player opens **Join a game** and connects to the server (UC-04).
2. The server's `welcome` says this player may create games.
3. The Games list shows **Create a game** under the games.
4. The player clicks it; the setup screen says "On the server …".
5. Players, card sets, password and detailed setup, as in UC-03.
6. The player clicks **Create**; the server opens the table.
7. Friends connect to the same server and join it (UC-04).

## Alternative scenarios and errors

* **3a. A game hosted in another player's app:** no **Create a game**.
  * Only the host's own app opens games there (N-17).
  * A friend's `create` or `load` gets `not_host`.
* **3b. This app hosts the server:** no **Create a game**; the host screen does it.
* **4a. Back:** the Games list again; the connection stays.
* **5a. No nickname field:** the nickname went with hello already.
* **5b. No "Continue a saved game":** saved games are in the Games list.
* **6a. A bad setup or missing set:** the error shows on the setup screen.

## Postconditions

* The player sits at seat 1 of the new game on the server.
* Leaving it keeps the connection; the server keeps running.
