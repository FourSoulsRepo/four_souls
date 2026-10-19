# UC-03: Host a game

**Module:** `internal/server`, `host.go`, `frontend`
**Status:** In progress
**Actors:** Host (a player)
**Goal:** A game waits in the lobby for friends to join
**Preconditions:** The app runs; or a dedicated server is reachable

## Main scenario (happy path)

1. The host starts hosting in the app.
2. The app runs the server and shows its LAN and virtual LAN addresses.
3. The host connects to it like any client and says hello with a nickname.
4. The host creates a game: 2–4 seats and the card sets to play with.
5. The server opens a table and seats the host at seat 1.
6. Friends join (UC-04); the host sees their nicknames appear.
7. Everyone marks ready; the game starts when every seat is taken.

## Alternative scenarios and errors

* **4a. Seats outside 2–4:** the server refuses with `bad_setup`.
* **4b. Unknown card set:** the server refuses with `bad_setup`.
* **4c. The host lacks a card set:** refused with `missing_sets` (CD-03).
* **7a. Someone leaves before the start:** their seat is free again.
* **7b. The host leaves:** the table stays while anyone sits at it.
* **1a. In the app:** Main menu, **Host a game**.
  * Nickname and players (2–4); "Detailed setup" shows the rest (GS-10).
  * Detailed: response timer, bonus souls.
  * Draft and bans wait for the tournament tools (roadmap, Future).
  * Host bans need the card list; they come with the card component.
  * Detailed also has deck ratios, on by default (UC-20).
* **2a. The port is taken:** the app shows the error and stays on the host screen.
* **6a. The table screen:** the host sees "Friends join at" and the addresses.
* **6b. Copy an address:** the addresses can be selected and copied (V-10).
* **6c. LAN players:** while a seat is free, the app answers LAN searches (UC-17).
* **1b. Password (N-14):** the host may set one, up to 64 characters.
  * Friends need it to join; the game list shows a lock.
  * It is never shown again, logged, recorded or saved.
* **7c. The host leaves in the app:** the app stops its server; friends lose the connection.
* **1c. Card sets (CD-02):** the host ticks the sets, oldest first.
  * The Base Game is always in; its box can't be unticked.
  * The server adds the Base Game if a client leaves it out.
  * The game list shows each game's sets, e.g. "Base Game, Gold Box".
  * Gold Box and Four Souls+ are core sets: ticked by default, optional.
* **1d. The last setup:** the host screen starts from the last hosted game.
  * Players, card sets, detailed setup, response timer, bonus souls.
  * Kept in the app settings; the password is never kept (N-14).
  * The deck ratios box too (UC-20).
* **1e. On a dedicated server:** the player creates the game there (UC-19).
* **4d. A friend creates on the host's app:** refused with `not_host` (N-17).
  * The friend's app offers no "Create a game" there.

## Postconditions

* The game runs; every seat gets its own view (ADR 006).
