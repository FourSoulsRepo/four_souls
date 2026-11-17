# UC-44: Play with bots

**Module:** `frontend`, `internal/server`
**Status:** Implemented
**Actors:** Host (a player), other players
**Goal:** A game starts without enough friends: bots fill the free seats
**Preconditions:** The host opened a table (UC-03 or UC-19); the server allows bots (ADR 030, Decision 10)

## Main scenario (happy path)

1. The host's table shows "Add bot" with a button per level on each free seat: Standard and Hard.
2. The host adds a bot; the seat shows a robot, "Bot N" and its level, and is ready.
3. The host may change a bot's level, or remove the bot to free its seat.
4. The host fills every free seat with bots and marks ready.
5. Every seat is taken and every human is ready: the game starts.
6. The bots throw their own first-player dice and play their turns.

## Alternative scenarios and errors

* **1a. A player who is not the host:** sees each bot with its robot, "Bot N (level)" and Ready, and no bot buttons.
* **1b. The server allows no bots (`-bots=off`):** no bot buttons; one level only (`-bots=standard`): one button.
* **1c. A saved game's table:** no bot buttons; its seats wait for its players.
* **2a. A friend joins later:** takes a free seat only, never a bot's; the game list counts bots as players.
* **5a. A bot is never shown as away or disconnected.**

## Postconditions

* The game runs with the bots in their seats; the e2e test is `e2e/tests/bots.spec.ts`.
