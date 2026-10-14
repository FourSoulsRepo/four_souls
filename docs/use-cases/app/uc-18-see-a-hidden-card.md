# UC-18: See a hidden card

**Module:** `frontend/src`, `pkg/rules_engine`
**Status:** Implemented
**Actors:** Player
**Goal:** The player sees the cards an effect lets them look at
**Preconditions:** A game runs; the player uses a "look at" effect

## Main scenario (happy path)

1. The player uses an effect that looks at hidden cards.
   1. Example: Forever Alone, "Look at the top card of a deck".
2. The player picks the deck, if asked.
3. The effect resolves.
4. A "You look at" dialog shows the cards in the middle (V-07).
5. The player may zoom each card as usual (UC-11).
6. The dialog stays through later updates.
7. The player presses OK; the dialog closes.

## Alternative scenarios and errors

* **4a. Other players:** they get no card and no dialog.
  * Their history says "Ann looks at a card".
* **4b. A choice comes with it:** the dialog sits in the middle.
  * The choice stays at the bottom edge, below it.
* **7a. History:** the looker's history keeps "You look at <card>".
  * Its zoom shows the card again.
* **2a. A choice about the card (19.11):** e.g. Cain's Eye, Golden Horseshoe.
  * "Put the top card on the bottom?" shows the card above its buttons.
  * The buttons name the card by its title.
  * Only the chooser sees the card; others see no dialog.

## Postconditions

* The cards stay where they were; only the looker saw them.
