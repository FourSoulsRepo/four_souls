# UC-24: Follow fights and played cards

**Module:** `frontend/src/components/Card.tsx`, `frontend/src/components/StackPanel.tsx`
**Status:** Implemented
**Actors:** Player
**Goal:** Every player follows a fight and another player's loot card
**Preconditions:** The game runs (UC-12); it is Ann's turn

## Main scenario (happy path)

1. Ann attacks Pooter.
2. Ann's character and Pooter get a red frame on every screen.
3. Ann throws the attack roll.
   1. Every player sees the die fall and the last roll by the stack.
   2. The stack names Pooter as the target.
4. The attack ends; both frames go.
5. Ann plays Bomb!.
6. On Bo's screen, Bomb! shows large in the middle for a moment.
   1. "Ann plays" stands above it.
7. Bomb! then stays under the stack, with "Ann plays" beside it.
8. Bomb! resolves; its damage plays while the card is still there.
9. The card leaves when the table shows the resolved stack.

## Alternative scenarios and errors

* **1a. Deck attack or no target yet:** no frame until the target is set.
* **1b. Ann attacks Bo:** both characters get the frame.
* **5a. Ann's own screen:** no large card; she played it herself.
* **5b. Animations off:** no pop; the card still stays under the stack.
* **7a. Bo responds with his own loot card:** Ann's card stays.
  * Bo's screen shows only other players' cards.

## Postconditions

* No frame and no played card stay after the fight and the stack end.
