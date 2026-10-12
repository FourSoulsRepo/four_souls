# UC-14: Follow the stack

**Module:** `frontend/src/stack.ts`, `frontend/src/components`
**Status:** Implemented
**Actors:** Player
**Goal:** The player always sees why the game waits and can respond
**Preconditions:** The game runs (UC-12)

## Main scenario (happy path)

1. Someone plays a loot card or uses an ability.
2. The stack panel shows it on every player's table.
   1. The top item, which resolves next, comes first.
   2. Rolls show their result; damage shows its target.
3. The action bar says who may respond: "Bo may respond".
4. The player who may respond sees "Respond, or pass".
5. That player passes or acts (UC-13).
6. The item resolves; the panel shows what is left, or hides.

## Alternative scenarios and errors

* **3a. Skip all (N-06):** any player presses "Skip all" while the stack is not empty.
  * The server passes for them on everything now on the stack.
  * The button turns into "Cancel skip".
* **3b. Response timer (N-07):** the host set one in the detailed setup.
  * The action bar counts down the seconds left.
  * When it runs out, the server answers for the player.
* **3c. Choices and discards:** the bar says "Ann is choosing" or "Ann is discarding".
* **2c. Whom it waits for (V-11):** the panel says "Waiting for Bo, Cy".
  * "Passed: Ann" lists who passed, auto-skipped players too.
  * Every client names the same players; the viewer is "You".
* **2d. A prompt on top:** the panel says "Bo: Choose a target".
  * A discard reads "Bo: Discard 2 loot cards".

## Postconditions

* Nobody waits without seeing who and why.
