# UC-26: Roll for the first player

**Module:** `rules_engine`, `internal/server`, `frontend`
**Status:** Implemented
**Actors:** Every seated player; spectators watch
**Goal:** Everyone sees each die of the first-player roll and who goes first
**Preconditions:** The game was just created (UC-03, UC-04); nobody plays Cain

## Main scenario (happy path)

1. The table opens before turn 1.
   1. The turn badge says "Before turn 1 · First player roll".
   2. No seat is active; no turn buttons show.
2. The stack panel lists every player in seat order.
   1. Its title says "Lowest roll goes first".
3. The first seat gets the Throw button (UC-13).
   1. The others see "Ann is throwing…".
4. The player throws; the die falls on every screen.
5. Their die shows by their name, in their dice look.
6. Each next seat throws the same way.
7. After the last die, the lowest roll wins.
8. A banner names the winner: "Bo goes first — rolled 1".
9. Turn 1 starts with that player.
10. The history keeps every die and the winner (UC-15).

## Alternative scenarios and errors

* **7a. The lowest roll ties:** a banner says "Tie at 2: Ann and Bo roll again".
  1. Only the tied players throw again, in seat order.
  2. The panel title says "Tie: roll again, lowest goes first".
* **1a. Someone plays Cain:** no roll; the banner says "Bo goes first — Cain".
* **3a. A seat was kicked:** the server throws for it at once (UC-06).
* **3b. The player is idle:** the response timer throws for them.
* **6a. Someone plays Eden:** her choice comes after the roll.
* **Reconnect during the roll:** the panel shows the dice of this round (UC-07).

## Postconditions

* Turn 1 belongs to the winner of the roll, or to Cain's player.
* The match record holds each throw; a replay draws the same dice.
