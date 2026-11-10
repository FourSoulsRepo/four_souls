# UC-40: Rooms on the table

**Module:** `frontend`, `internal/server`, `rules_engine`
**Status:** Implemented
**Actors:** Host, active player
**Goal:** A game with Requiem plays with the room deck (R-SETUP-16)
**Preconditions:** The host creates a game with Requiem (UC-03)

## Main scenario (happy path)

1. The host opens "Detailed" and ticks "Room deck".
   1. The box shows only while Requiem is ticked.
   2. It is off by default; the host screen remembers it.
2. The server checks the sets have rooms.
3. The game starts with one room slot and the room deck.
4. Every player sees the room on the board, landscape.
   1. The room deck and its discard are under it.
5. On their turn, the active player uses a room's ability (R-CARD-36).
   1. A usable room glows; a click uses it, as an item's.
6. A monster died this turn: the end phase asks about the room (R-TURN-12).
7. The active player picks the room: it goes into discard.
8. The next room comes from the room deck into the slot.
9. The table screen shows "room deck" in its options.
10. The record header says `room_deck`.

## Alternative scenarios and errors

* **1a. Requiem is not ticked:** no box; the game has no rooms.
* **2a. The sets have no rooms (every room banned):** the server turns it off.
* **5a. Another player's turn:** the room shows why it can't be used.
* **5b. A room that came in the end phase:** usable from the next turn (R-TURN-18).
* **6a. No monster died this turn:** no question; the room stays.
* **7a. The player picks "Keep":** the room stays.
* **7b. A room that can't be put into discard (Social Goals):** not offered.
* **Card images missing:** rooms show as landscape text cards.

## Postconditions

* One room lies on the room slot, unless the room deck ran out.
