# UC-36: See a card flip

**Module:** `rules_engine`, `cards.go`, `frontend/src/components`
**Status:** Implemented
**Actors:** Player, spectator
**Goal:** Everyone sees which side of a two-sided card is up
**Preconditions:** A two-sided card is in play (Anima Sola, The Enigma, The Harbingers)

## Main scenario (happy path)

1. The card shows its front, named by the front alone ("Anima Sola").
2. An effect flips it (↷: Flip this item).
3. On every table the card turns edge-on and back.
4. It now shows its back: name, image, text and stats ("The Revenant").
5. The history says "The Revenant turns face up".
6. Flipped again, it turns back to its front the same way.

## Alternative scenarios and errors

* **3a. Animations off, or a rollback's table:** the side changes at once.
* **4a. No card images:** the blank card shows that side's name and text.
* **4b. Zoom:** the large copy shows the side that is up.
  * The other side shows beside it (UC-11, 3f).

## Postconditions

* Counters, damage and charge stay as they were (R-MECH-53).
* Only the side that is up counts in the game (R-MECH-33).
