# UC-14: Follow the stack

**Module:** `frontend/src/stack.ts`, `frontend/src/components`
**Status:** Implemented
**Actors:** Player
**Goal:** The player always sees why the game waits and can respond
**Preconditions:** The game runs (UC-12)

## Main scenario (happy path)

1. Someone plays a loot card or uses an ability.
2. The stack panel shows it on every player's table, in the board's left column.
   1. The top item, which resolves next, comes first.
   2. Rolls show their result; damage shows its target.
   3. Each row: a square icon and one short line.
   4. Icons: card art, a die, a blood drop, or a skull.
   5. Players show in their seat color; targets in red.
3. The prompt banner says whom it waits for: "Waiting for Bo" (UC-21).
4. The player who may respond sees "Respond, or pass".
5. That player passes or acts (UC-13).
6. The item resolves; the panel shows what is left, or "Nothing on the stack".

## Alternative scenarios and errors

* **3a. No "Skip all" (15.84):** the button is gone.
  * A player who can act passes each item, or waits for their turn (UC-06).
* **3b. Response timer (N-07):** the host set one in the detailed setup.
  * The prompt banner counts down the seconds left.
  * When it runs out, the server answers for the player.
* **3c. Choices and discards:** the banner says "Waiting for Ann" and what for.
* **2c. Whom it waits for (V-11):** the panel says "Waiting for Bo, Cy".
  * "Passed: Ann" lists who passed, auto-skipped players too.
  * Every client names the same players; the viewer is "You".
* **2d. A prompt on top:** the panel says "Bo: Choose a target".
  * A discard reads "Bo: Discard 2 loot cards".
* **2e. Hover a row:** its card shows large.
  * Below it: the used ability's words, highlighted.
  * Then the owner's name, in their color.
  * On touch, a tap shows the same until the next tap.
* **2f. Attacks (15.22):** a roll names the monster: "Attack: Ann rolls 4 against Fly".
  * Damage names its dealer: "Fly deals 1 damage to Ann".
  * Hover shows the monster, and this line under it.
* **2g. A roll's result (15.22):** the result item has only the rolled line.
  * "1: Each player gains 1¢." instead of the whole card text.
  * The item that rolls says only "Roll."; the card shows every result.
* **2h. Several abilities (The Bone):** only the used ability's words show.
* **2i. A player takes damage (15.47):** the line names the card that deals it.
  * Boom Fly: "Boom Fly deals 1 damage to Ann", one line per player.
  * No card to name: "Ann takes 1 damage".
* **2j. A roll waits for its throw (15.61):** the roll's row shows "?".
  * The roller's panel says "Throw the dice".
  * The others read "Ann: Throw the dice for the attack roll" in the panel.
  * The others see Ann's die shaking where it will land.
  * The banner says "Throw the dice", or "Waiting for Ann".
* **2k. The last roll (15.63):** a small die by the "Stack" title shows it.
  * It shows the final value, after changes and rerolls.
  * Hovering it names the roller and the roll: "Ann rolled 4 — attack roll".
  * It stays until the next roll; before the first one there is none.
* **2l. A prompt of the viewer's (15.28):** the panel says what to do.
  * The same words as the banner: "Choose a target".
  * Pass glows when passing answers the prompt.
* **2m. An open attack or buy (15.91):** a row on top of the stack says it.
  * It is not a stack item: no icon, a dashed edge.
  * "Bo attacks…", then "Bo attacks Pooter" once aimed.
  * A monster deck attack: "Bo attacks the monster deck".
  * A buy: "Bo buys…".
  * Every player sees the same row.
  * It goes when the attack or the buy ends.

## Postconditions

* Nobody waits without seeing who and why.
