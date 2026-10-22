# UC-17: Join a game found on the LAN

**Module:** `internal/lan`, `internal/server`, `lan.go`, `frontend/src/screens/JoinScreen.tsx`
**Status:** Implemented
**Actors:** Player
**Goal:** The player joins a friend's game without typing an address (N-13)
**Preconditions:** A friend hosts a game from the app on the same network

## Main scenario (happy path)

1. The player opens **Play**, then **LAN** (15.97).
2. The app searches the network again and again while the screen shows.
3. "On your network" lists each game found.
   1. Host nickname, taken and free seats, address.
4. The player presses **Join** next to a game.
5. The app connects to that address and lists its games (UC-04).
6. The player joins the friend's game.

## Alternative scenarios and errors

* **3a. Nothing found:** "No games found yet. Still looking…"
  * A game a friend opens shows up by itself; there is no button.
* **3b. Another protocol version:** the game shows "needs app X"; Join is off.
* **3c. The game is full:** listed; Join is off.
* **3d. A dedicated server:** never listed; join it by address (UC-04).
* **2a. macOS asks for Local Network access:** "no" hides LAN games.
* **2b. The router blocks broadcasts:** nothing is found; the address field works.
* **3e. A locked game:** it shows a lock; the password is asked in the game list (UC-04).
* **3f. The game started:** listed "In progress"; Join is off (15.97).
  * The app held a seat there (it crashed or closed): **Rejoin** instead.
  * Rejoin takes the seat back at once, as after UC-07.
* **3g. The player types an address:** **Host address**, then **Connect**.
  * The LAN screen keeps the last typed address, also after a restart.
  * It keeps no Recent list; Online has that (UC-04).
* **3h. Back or Esc:** the Games list goes back to the LAN screen; it goes back to the menu.

## Postconditions

* The player sits at the friend's table, as after UC-04.
