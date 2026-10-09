# UC-04: Join a game

**Module:** `internal/server`
**Status:** In progress
**Actors:** Player
**Goal:** The player sits at a friend's game and plays
**Preconditions:** The player knows the server's address

## Main scenario (happy path)

1. The player enters the address and connects.
2. The app says hello: protocol version, nickname, card sets.
3. The server welcomes the player and lists its games.
4. The player picks a game that has a free seat.
5. The server seats the player and shows the table to everyone at it.
6. The player marks ready; the game starts when everyone is ready.

## Alternative scenarios and errors

* **2a. Different protocol version:** `client_outdated` or `server_outdated`, then the connection closes.
* **2b. Empty or too long nickname:** refused with `bad_name`.
* **4a. The game is full:** refused with `room_full`.
* **4b. The game has started:** refused with `started`.
* **4c. The player lacks a card set of the game:** refused with `missing_sets` (CD-03).
* **5a. The nickname is taken at this table:** the server adds " (2)".
* **6a. The connection drops during the game:** the player says hello again with their token and gets the seat back (N-08).

## Postconditions

* The player has a seat and sees the game from it.
