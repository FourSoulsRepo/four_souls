# UC-21: Know when the game waits for me

**Module:** `frontend/src/prompt.ts`, `frontend/src/components`
**Status:** Implemented
**Actors:** Player
**Goal:** The player always knows whether the game waits for them, and for what
**Preconditions:** The game runs (UC-12)

## Main scenario (happy path)

1. The game asks Ann something.
2. Ann's banner, bottom right, turns gold and its glow breathes.
   1. Its title says what to do, in a few words.
   2. A hint below adds detail.
3. On her turn it says "Your turn".
   1. The hint lists what is left: "Play a card, attack, buy, or end the turn".
   2. With nothing left: "Nothing left to do: end the turn".
4. Every other player's banner says "Waiting for Ann".
   1. Her name is in her seat color.
   2. The hint says what for: "Playing their turn".
5. Ann answers (UC-13); the banners follow the next prompt.

## Alternative scenarios and errors

* **3a. Respond:** "Respond, or pass".
  * The hint names the stack's top item.
  * Pass in the stack panel glows.
* **3b. Act out of turn:** "You may act"; others read "May act, or pass".
* **3c. A choice:** the title is its question: "Choose a target".
  * Options on the table glow; every other card dims.
  * The stack panel says the same (UC-14).
* **3d. Discard:** "Discard 1 loot card".
* **3e. Put cards back:** the order question, e.g. "Choose the card for the very bottom".
* **3f. Throw (15.61):** "Throw the die", hint "For the attack roll".
* **3g. Own aim (15.25):** "Choose a target".
  * The hint names the card and says Esc cancels.
  * Only the targets and the aiming card glow; the rest dims.
* **3h. Response timer (N-07):** the banner counts down the seconds.
* **3i. Paused (N-08):** the banner says "Paused".
* **3j. Game over:** no banner; the end window shows.

## Postconditions

* No screen looks idle while it waits for its player.
* Everyone else knows whom the game waits for, and what for.
