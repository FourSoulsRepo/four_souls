# UC-11: Look at a card

**Module:** `frontend/src/components`, `cards.go`
**Status:** Implemented
**Actors:** Player
**Goal:** The player reads any card on screen, with or without card images
**Preconditions:** A card is on screen

## Main scenario (happy path)

1. The card shows its image.
2. The player rests the mouse on it.
3. A large copy shows right beside the card (15.27).
4. The mouse leaves the card; the large copy goes away.

## Alternative scenarios and errors

* **1a. No image:** the app has no image file, or it fails to load (A-11).
  * A blank card shows instead, coloured by its deck.
  * Small sizes show name, type, rewards and stats.
  * The large copy adds the effects.
  * Text icons show as short labels, e.g. ↷ for {Tap Effect}.
* **1b. Unknown card:** the app has no data for it; its ID shows as the name.
* **2a. Touch:** the player holds a finger on the card for about half a second.
  * The large copy shows in the middle; the next tap closes it (B-06).
  * A short tap does not zoom.
* **1c. Face-down card:** shows its deck name; it cannot be zoomed.
* **3a. Badges (15.23):** the large copy shows no badges.
  * Its words stay readable; the card on the table keeps them.
* **3b. Hand card (15.27):** the large copy rises from the card.
  * It stands on the window's bottom edge, over the card.
* **3c. Keywords (15.27):** tips stand beside the large copy.
  * One per keyword: Eternal, Trinket, Curse, Guppy.
  * Each names the keyword and says what it means.
  * They show with or without card images.
* **3d. Near an edge:** the large copy stays inside the window.
* **3e. Card in a window (19.26):** the large copy shows over it.
  * Discard, choice, look and order windows alike.
  * The player reads the card before picking it.

## Postconditions

* Nothing changes in the game; looking is local.
