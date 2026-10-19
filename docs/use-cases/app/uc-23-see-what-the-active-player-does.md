# UC-23: See what the active player does

**Module:** `frontend/src/prompt.ts`, `frontend/src/components/PlayerArea.tsx`
**Status:** Implemented
**Actors:** Player, spectator
**Goal:** Everyone else knows what the active player is doing now
**Preconditions:** The game runs (UC-12); it is Ann's turn

## Main scenario (happy path)

1. Ann holds priority on her turn.
2. Her seat shows "Ann is thinking" under her name.
   1. Her name is in her seat color.
3. Ann attacks Pooter.
4. Her seat shows "Ann attacks Pooter".
   1. Pooter is in red, as on the stack.
5. The attack ends; the line follows what Ann does next.

## Alternative scenarios and errors

* **2a. Ann's own screen:** her seat shows no such line.
  * Her prompt banner says what to do (UC-21).
* **2b. The game waits for another player:** no line on Ann's seat.
  * Unless one of the cases below goes on.
* **3a. A play:** "Ann plays Bomb!" while her loot card is on top of the stack.
* **3b. A use:** "Ann uses The D6" while her ability is on top.
* **3c. Attack declared, no target yet:** "Ann declares an attack".
* **3d. Deck attack:** "Ann attacks the monster deck" until the card is placed.
* **3e. Attack on a player:** the player is named: "Ann attacks Bo".
* **3f. Another player responds in the attack:** the line still names the attack.
* **3g. Purchase declared:** "Ann is buying".
* **3h. Death penalty:** "Ann pays the death penalty".
* **3i. Game over:** no line.

## Postconditions

* Every seat but Ann's own reads the same line.
* The line uses only public state, so spectators can show it too.
