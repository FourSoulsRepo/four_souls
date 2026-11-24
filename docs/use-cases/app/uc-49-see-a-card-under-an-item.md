# UC-49: See a card under an item

**Module:** `pkg/rules_engine/view.go`, `frontend/src/components/PlayerArea.tsx`, `frontend/src/components/CardZoom.tsx`
**Status:** Implemented
**Actors:** Player, spectator, record viewer
**Goal:** Everyone sees which card lies under an item
**Preconditions:** The game runs (UC-12); Ann has Friendly Ball in play

## Main scenario (happy path)

1. Ann uses Friendly Ball's first ↷ on Fly.
2. Fly leaves its slot and goes under the ball (R-ZONE-13).
3. Fly peeks out from under the ball, face up (19.244).
   1. On the ball's right, a little turned, as on a real table.
   2. Every player sees it, and so does a spectator.
   3. The ball lies on top; the items row keeps room for Fly.
4. A hover on the peeking part of Fly zooms Fly.
5. A hover on the ball zooms the ball.
   1. Fly stands beside the zoom, smaller, captioned "Under it".
6. Ann uses the second ↷: Fly goes back to a monster slot.
7. Nothing peeks out from under the ball any more.

## Alternative scenarios and errors

* **3a. Several cards under one card:** each peeks out a step further.
  * The first one put under lies lowest, furthest out.
* **3b. Text cards, without the private images:** the same peek.
  * The peek shows the right part of Fly's drawn face.
* **4a. Touch:** a long press on the peek zooms Fly; on the ball, the ball and Fly.
* **3c. A record or replay:** the peek shows as it did in the game.

## Postconditions

* Nobody has to remember what a card holds under it.
