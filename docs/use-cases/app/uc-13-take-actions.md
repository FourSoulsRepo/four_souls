# UC-13: Take actions

**Module:** `frontend/src/actions.ts`, `frontend/src/components`
**Status:** Implemented
**Actors:** Player
**Goal:** The player plays a whole turn with mouse or touch
**Preconditions:** The game runs (UC-12); the player must act

## Main scenario (happy path)

1. The table highlights what the player may do now.
   1. Yellow: cards to play or use.
   2. The stack panel: Pass.
   3. Icons by own name: Buy and Attack; lit when usable, dim when not.
   4. A big End turn button at the board's right edge.
2. The player drags a loot card onto the board to play it.
3. Others may respond; the player passes or acts again.
4. The player presses Attack.
5. The monsters that can be attacked glow blue; the player clicks one.
6. Rolls and responses play out.
7. The player presses End turn; the next player's turn starts.

## Alternative scenarios and errors

* **2a. Click or tap:** a card opens its action menu: Play, Use, Cancel.
* **2b. Abilities:** a character or item with abilities shows "Use" or "Use ability N".
* **5a. A choice dialog:** every choice also shows as a dialog.
  * Card options show as cards (V-07); text options as buttons.
  * A question shows its text as the title.
* **5b. Pick a player:** a "Choose <name>" button in that player's area.
* **3a. Discard to hand size:** pick exactly N hand cards, then Discard.
* **1a. Not the player's move:** the bar says "Waiting for <name>".
* **1b. Updates still playing:** buttons stay but are disabled.
* **2c. Refused:** the bar shows the engine's reason and rule, e.g. "(R-TURN-06)".
* **4a. Keys:** B buys, A attacks; each icon's tooltip names its key.
* **4b. Pressed:** a pressed icon shows at once, before the server answers.
* **4c. Choosing:** while the player picks what to buy or attack, that icon stays pressed.
* **7a. Their turn:** the big button is grey and says "Their turn".
* **7b. Not yet:** on own turn, End turn stays dim until the engine allows it.
* **7c. Nothing left:** End turn glows green when only passing or ending is left.
* **7d. Keys:** E ends the turn; Space passes.
* **2d. Several abilities (15.24):** the card opens a dialog with the card large and one line per ability.
  * The player clicks a line to use that ability.
  * Lines not usable now stay greyed; Cancel or a click outside closes it.
* **4d. What is left (15.24):** on own turn the Buy and Attack icons carry a count badge.
  * Loot plays left show by the player's name as "Loot N".
* **2e. Aim (15.25):** a card that needs a target is aimed by two clicks.
  * The first click starts aiming; the card stays lifted in the hand.
  * The valid targets glow red: monsters, players, items, stack items.
  * With a mouse, a red arrow points from the card to the pointer.
  * The second click, on a target, plays the card at it; no dialog asks.
  * A click on empty space, a right click, or Esc cancels; nothing is sent.
  * Such a card opens no menu and is not dragged.
  * The arrow is a chain of red blocks with a big head (15.32).
  * Over a valid target a red ring marks the head.
  * A click on the aiming card again does nothing (15.33).
  * A click on a card, pile, or stack item that is no target keeps aiming.
  * A short note above the click says what the card may target.
  * E.g. "Pick a monster or a player."; it fades after 2.5 s.
  * A player is aimed at by the character card, as a monster is.
  * The character glows red; the player's area stays plain.
  * A card that may target itself, as The D20, picks itself on that click.
* **2f. Place, then aim (15.34):** a card with no target is dragged onto the board; no arrow.
  * A card that asks its mode first is dragged too; the mode dialog asks.
  * Then the target is aimed from the card, as in 2e; no dialog asks.
  * The card stays lifted in the hand until the target is picked.
  * Esc, a right click, or empty space cancels the whole play.
  * The card stays in hand; nothing goes on the stack.
  * A card with no valid target cannot be played; the engine decides.
  * A target that is not on the table keeps the dialog.

## Postconditions

* Only intents the server allows are offered; the server still checks each one.
