# UC-25: See what others aim at

**Module:** `frontend`, `internal/server`
**Status:** Implemented
**Actors:** A player watching another player's move
**Goal:** The player sees live what another player aims at or looks at
**Preconditions:** A game runs; the server has presence on (ADR 016)

## Main scenario (happy path)

1. Another player starts aiming with a card.
2. An arrow runs from that card on this screen.
   1. A hidden hand card shows as a back lifted from its slot.
3. The arrow follows the other player's pointer.
4. The card under the arrow head gets a red frame.
5. The other player picks a target or cancels.
6. The arrow, the lifted back, and the frame go.

## Alternative scenarios and errors

* **1a. The active player hovers a card:** that card gets a red frame here.
* **1b. A waiting player hovers a card:** nothing is shown to the others.
* **3a. The pointer is over no card:** the head points at that spot of the board.
* **4a. The aimed card is gone from this screen:** no frame shows.
* **5a. A message is lost:** the arrow goes after 5 s without news.
* **5b. The other player drops out:** the server clears their arrow.
* **Server with presence off:** nothing is sent or shown.

## Postconditions

* No arrow or frame of another player stays on the table.
* No hidden card was named on the wire (A-16).
