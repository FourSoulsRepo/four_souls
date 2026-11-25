# UC-49: See a card under an item

**Module:** `pkg/rules_engine/view.go`, `frontend/src/components/PlayerArea.tsx`, `frontend/src/components/CardZoom.tsx`
**Status:** Implemented
**Actors:** Player, spectator, record viewer
**Goal:** Everyone sees which card lies under an item
**Preconditions:** The game runs (UC-12); Ann has Friendly Ball in play

## Main scenario (happy path)

1. Ann uses Friendly Ball's first ↷ on Fly.
2. Fly leaves its slot and goes under the ball (R-ZONE-13).
3. Fly lies behind the ball, face up, a little higher (19.244).
   1. Its top band shows above the ball: the name and a little art.
   2. Every player sees it, and so does a spectator.
   3. The ball lies on top and keeps its place in the items row.
   4. The band may reach onto the board; it shows over it.
   5. In a side column the ball moves down by the band.
   6. So the band never hides the character's stats.
4. A hover on Fly's band zooms Fly.
5. A hover on the ball zooms the ball.
   1. Fly stands beside the zoom, smaller, captioned "Under it".
6. Ann uses the second ↷: Fly goes back to a monster slot.
7. No band shows above the ball any more.

## Alternative scenarios and errors

* **3a. Several cards under one card:** each lies a band higher.
  * The first one put under lies lowest, at the top.
* **3b. Text cards, without the private images:** the same band.
  * The band shows the top of Fly's drawn face: its name.
* **4a. Touch:** a long press on the band zooms Fly; on the ball, the ball and Fly.
* **3c. A record or replay:** the band shows as it did in the game.
* **3d. A full items row:** it still scrolls sideways; no band is clipped.

## Postconditions

* Nobody has to remember what a card holds under it.
