# UC-37: Seat the players in a random order

**Module:** `internal/server`, `frontend`
**Status:** Implemented
**Actors:** Host (a player)
**Goal:** Nobody gets a seat by connecting first
**Preconditions:** The host creates a game (UC-03) with the detailed setup

## Main scenario (happy path)

1. The host ticks "Random seats" in the detailed setup.
2. The table's setup summary says "random seats".
3. Friends join and take a free seat, as before.
4. Everyone marks ready.
5. The server shuffles the seats once, when the game starts.
6. The game starts with the players in that random order.

## Alternative scenarios and errors

* **4a. Random seats off:** the players keep the order they joined.
* **1a. The last setup:** the box starts as it was when last hosted (UC-03).

## Postconditions

* Every player holds one seat; the seats are in a random order.
