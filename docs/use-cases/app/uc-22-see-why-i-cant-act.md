# UC-22: See why I can't act

**Module:** `frontend`, `rules_engine`
**Status:** Implemented
**Actors:** Player
**Goal:** The player learns why a card or action does nothing now
**Preconditions:** The game runs (UC-12)

## Main scenario (happy path)

1. The player clicks one of their own cards that does not glow.
   1. E.g. a loot card in hand, or an item with an activated ability.
   2. The app asks the server why; nothing else is sent.
2. A short note shows at the card, with the reason.
   1. E.g. "No loot plays left", or "Loot is played on your turn".
3. The note fades after 2.5 s.
4. Hovering the card later shows the same reason as a red line under it.
   1. It holds until the game moves on.

## Alternative scenarios and errors

* **1a. Another player's card, or a shop item, monster or deck:** a click does nothing.
  * No reason is asked for, and none is shown.
* **1b. Buy or Attack:** a dim icon takes the click and shows the note.
  * E.g. "No attacks left", "Not your turn", "Wait for the stack to resolve".
  * Its tooltip says the same. These reasons come with each update.
* **1c. End turn:** on own turn, a dim End turn says why on a click.
  * E.g. "You must attack first".
* **2a. Usable with priority:** "Wait for your priority".
* **2b. Another prompt is open for the player:** "Answer the prompt first".
* **3a. Updates still playing:** the click waits like any other.
* **4a. Aiming:** a click on a wrong target shows the aim's note (UC-13, 2e).
* **5a. The game moves on:** the reason clears; a new click asks again.

## Postconditions

* The server has said why the card can't be used, at the step it was asked at.
* The turn action reasons come with every update, as before.
