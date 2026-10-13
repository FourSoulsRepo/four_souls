# UC-11: Look at a card

**Module:** `frontend/src/components`, `cards.go`
**Status:** Implemented
**Actors:** Player
**Goal:** The player reads any card on screen, with or without card images
**Preconditions:** A card is on screen

## Main scenario (happy path)

1. The card shows its image.
2. The player rests the mouse on it.
3. A large copy shows on the side away from the pointer.
4. The mouse leaves the card; the large copy goes away.

## Alternative scenarios and errors

* **1a. No image:** the app has no image file, or it fails to load (A-11).
  * A blank card shows instead, coloured by its deck.
  * Small sizes show name, type, rewards and stats.
  * The large copy adds the effects and keyword boxes.
  * Text icons show as short labels, e.g. ↷ for {Tap Effect}.
* **1b. Unknown card:** the app has no data for it; its ID shows as the name.
* **2a. Touch:** the player holds a finger on the card for about half a second.
  * The large copy shows in the middle; the next tap closes it (B-06).
  * A short tap does not zoom.
* **1c. Face-down card:** shows its deck name; it cannot be zoomed.
* **3a. Badges (15.23):** the large copy shows the card's badges.
  * Counters, changed monster stats, and a character's health and attack.

## Postconditions

* Nothing changes in the game; looking is local.
