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
4. There is no "Skip all" (15.84); a player who can act passes each item.
5. The server never passes for a player who can act, unless they wait (7).
6. With a response timer, everyone sees when the waiting player's time ends (N-07).
7. On another's turn a player may press "Wait for my turn" (N-16).
8. The server then passes for them until the player before them ends the turn.
9. In that end window they are asked again: their last chance to act.

## Alternative scenarios and errors

* **6a. The response timer runs out:** the server answers for the player: it passes, or ends their turn, or discards their first cards, or takes the first option of a choice.
* **7a. Choices and discards:** they still ask the waiting player.
* **7b. Cancel wait:** the player is asked again from the next pass on.
* **1a. A client misses an update:** it asks for a resync and gets the whole view again.
* **6b. A roll waits for its throw:** when the timer runs out, the server throws for the roller (15.60).
* **7c. A roll waits for its throw:** "Wait for my turn" never throws.
  * The server only passes for a player; the roller is always asked to throw.
* **7d. Another player targets them (15.84):** the wait ends at once.
  * It counts the player, their character, items, souls, curses, hand cards.
  * Damage to them or their cards counts too.
  * Only items put on the stack after the press count.
  * They are asked as usual and may answer, e.g. play a card before it is stolen.
  * A card that picks its victim as it resolves counts too (15.85).
  * E.g. XV. The Devil, Joker; Dead Bird only for the player who rolled.

## Postconditions

* The game moves on as long as somebody can act.
