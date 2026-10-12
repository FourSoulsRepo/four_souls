# UC-17: Join a game found on the LAN

**Module:** `internal/lan`, `internal/server`, `lan.go`, `frontend`
**Status:** Implemented
**Actors:** Player
**Goal:** The player joins a friend's game without typing an address (N-13)
**Preconditions:** A friend hosts a game from the app on the same network

## Main scenario (happy path)

1. The player opens **Join a game**.
2. The app searches the network for about a second.
3. "On your network" lists each waiting game.
   1. Host nickname, taken and free seats, address.
4. The player presses **Join** next to a game.
5. The app connects to that address and lists its games (UC-04).
6. The player joins the friend's game.

## Alternative scenarios and errors

* **3a. Nothing found:** "No games found. Type the host address above."
  * **Search again** repeats the search.
* **3b. Another protocol version:** the game shows "needs app X"; Join is off.
* **3c. The game started or is full:** the host no longer answers; it is not listed.
* **3d. A dedicated server:** never listed; join it by address (UC-04).
* **2a. macOS asks for Local Network access:** "no" hides LAN games.
* **2b. The router blocks broadcasts:** nothing is found; the address field works.

## Postconditions

* The player sits at the friend's table, as after UC-04.
