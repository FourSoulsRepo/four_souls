# UC-22: See why I can't act

**Module:** `frontend`, `rules_engine`
**Status:** Implemented
**Actors:** Player
**Goal:** The player learns why a card or action does nothing now
**Preconditions:** The game runs (UC-12)

## Main scenario (happy path)

1. The player hovers a hand card that does not glow.
2. The card's zoom shows a red line under it.
   1. E.g. "No loot plays left", or "Loot is played on your turn".
3. The player clicks the card anyway.
4. A short note shows where they clicked, with the same words.
5. The note fades after 2.5 s.

## Alternative scenarios and errors

* **1a. A card in play:** an item or the character says why the same way.
  * E.g. "Used; it recharges on your turn", "Not enough ¢".
  * A card with no activated ability says nothing.
* **1b. Buy or Attack:** a dim icon takes the click and shows the note.
  * E.g. "No attacks left", "Not your turn", "Wait for the stack to resolve".
  * Its tooltip says the same.
* **1c. End turn:** on own turn, a dim End turn says why on a click.
  * E.g. "You must attack first".
* **2a. Usable with priority:** "Wait for your priority".
* **2b. Another prompt is open for the player:** "Answer the prompt first".
* **3a. Updates still playing:** nothing shows; the click waits like any other.
* **4a. Aiming:** a click on a wrong target shows the aim's note (UC-13, 2e).

## Postconditions

* Nothing is sent; the reasons come with each update from the server.
