# UC-27: Join a server by its domain

**Module:** `frontend`
**Status:** Implemented
**Actors:** Player
**Goal:** The player reaches a server by its domain name alone
**Preconditions:** The server has a domain, maybe behind an HTTPS proxy

## Main scenario (happy path)

1. The player types `games.example.com` in **Host address**.
2. The app tries `wss://games.example.com/ws` on port 443.
3. The proxy answers; the connection is encrypted.
4. Joining goes on as in UC-04.

## Alternative scenarios and errors

* **2a. No HTTPS on the domain:** the app tries `ws://games.example.com:4774/ws`.
  * The wait for the first try is at most 5 seconds.
* **2b. Neither answers:** "Could not connect to the server.".
* **1a. A typed scheme:** only that scheme is tried.
  * `wss://` or `https://`: port 443 unless typed.
  * `ws://` or `http://`: port 4774 unless typed.
* **1b. A port, an IP, or a LAN name:** plain `ws://` only, as before.
  * LAN names: no dot, or `.local`, `.lan`, `.home`, `.internal`, `.home.arpa`.
* **1c. Not an address:** the app shows how one looks (UC-04, 1b).
* **5a. The record download:** `https://` when the game ran over `wss://`.

## Postconditions

* The player is connected; **Recent** keeps the address as typed.
