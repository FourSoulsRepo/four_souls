# UC-03: Host a game

**Module:** `internal/server`, `host.go`
**Status:** In progress
**Actors:** Host (a player)
**Goal:** A game waits in the lobby for friends to join
**Preconditions:** The app runs; or a dedicated server is reachable

## Main scenario (happy path)

1. The host starts hosting in the app.
2. The app runs the server and shows its LAN and virtual LAN addresses.
3. The host connects to it like any client and says hello with a nickname.
4. The host creates a game: 2–4 seats, Base Game cards.
5. The server opens a table and seats the host at seat 1.
6. Friends join (UC-04); the host sees their nicknames appear.
7. Everyone marks ready; the game starts when every seat is taken.

## Alternative scenarios and errors

* **4a. Seats outside 2–4:** the server refuses with `bad_setup`.
* **4b. Unknown card set:** the server refuses with `bad_setup`.
* **4c. The host lacks a card set:** refused with `missing_sets` (CD-03).
* **7a. Someone leaves before the start:** their seat is free again.
* **7b. The host leaves:** the table stays while anyone sits at it.

## Postconditions

* The game runs; every seat gets its own view (ADR 006).
