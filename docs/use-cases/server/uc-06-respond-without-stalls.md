# UC-06: Respond without stalls

**Module:** `internal/server`
**Status:** In progress
**Actors:** Players
**Goal:** The game never waits for a player who has nothing to do
**Preconditions:** A game is running (UC-03, UC-04)

## Main scenario (happy path)

1. After every change each player gets the view and what they may do (N-04).
2. When a player could only pass, the server passes for them (N-05, N-09).
3. A player with a charged character or item is asked, never skipped.
4. A player may press "Skip all": the server passes for them on the stack they see (N-06).
5. "Skip all" ends when the stack is empty, or when a new item arrives; then they are asked again.
6. With a response timer, everyone sees when the waiting player's time ends (N-07).

## Alternative scenarios and errors

* **6a. The response timer runs out:** the server answers for the player: it passes, or ends their turn, or discards their first cards, or takes the first option of a choice.
* **1a. A client misses an update:** it asks for a resync and gets the whole view again.

## Postconditions

* The game moves on as long as somebody can act.
