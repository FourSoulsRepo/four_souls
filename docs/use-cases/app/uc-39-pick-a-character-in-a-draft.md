# UC-39: Pick a character in a draft

**Module:** `frontend`, `internal/server`
**Status:** Implemented
**Actors:** Host, players
**Goal:** Each player starts with a character they picked from their offers
**Preconditions:** The host chose a draft in detailed setup (UC-05)

## Main scenario (happy path)

1. The host opens "Detailed" on the host screen.
2. "Characters" is "One random each" by default.
3. The host picks "Draft" and sets "Characters offered": 2 to 5.
4. The host screen remembers both for the next game.
5. Everyone gets ready; the pick screen opens for each player.
6. Each player sees only their own offers.
7. Under each character is its starting item.
8. A click picks a character; a second click unpicks it.
9. Confirm sends the pick; it waits disabled until a card is picked.
10. The list of players shows who has picked, not what.
11. When everyone has picked, the game starts with the picks.

## Alternative scenarios and errors

* **7a. A character without a fixed item (Eden):** no item shows.
  * It says the starting item is chosen when the game starts.
  * After the pick, the player chooses it from the treasure deck (UC-26 6a).
* **7b. No card images:** characters and items show as text cards.
* **9a. The server refuses the pick:** the error shows; the player picks again.
* **5a. A player reconnects:** the pick screen opens again with the same offers.

## Postconditions

* Every player plays the character they picked.
* Each starts with that character's item, or the one they chose (Eden).
