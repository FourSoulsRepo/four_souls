# 12. Animation ownership

* **Status:** Accepted
* **Date:** 2026-10-15
* **Authors:** @HardDie

---

## Context

1. Three bugs in a row were two animations on one card.
   1. 19.3: the dice played before the monster reached its slot.
   2. 19.4: a deck attack showed the monster large twice.
   3. 19.6: the monster flew to its slot, then popped in there again.
2. The client has two sources of animation.
   1. Events: playback turns an update's events into effects (7.6).
   2. The view: `Card` pops in when a new `objectId` mounts.
   3. Neither knew what the other did to a card.
3. Each fix patched one pair of animations.
   1. Nothing listed every animation.
   2. Nothing failed when a new one clashed with an old one.
   3. A replaced animation could stay behind unnoticed.
4. The owner asked for a rule, a check, and research.

## Research

1. Turn-based and card games split the model from the view.
   1. The model resolves at once; only the animations queue.
   2. Hearthstone is described this way: state first, visuals catch up.
   3. Unity and indie devlogs: logic emits events, the view queues them.
2. A diff of two states gives the end, not the order.
   1. So games animate from an ordered event stream.
   2. Or they turn the diff into events first.
   3. Animating both the events and the diff draws things twice.
3. VCMI's battle view tracks running animations in one list.
   1. The display waits on that list, not on the game state.
4. State machines (XState) model one actor's states well.
   1. Good for a modal, a cart, a queue: idle, playing, showing.
   2. They order steps; they do not see two systems on one element.
5. Shared-element tools (Motion `layoutId`, FLIP) hand an element over.
   1. They fit one element moving between two places in one tree.
   2. Our effects draw copies in an overlay, outside the table tree.
6. Sources:
   1. [VCMI battle interface](https://wiki.vcmi.eu/Battle_interface)
   2. [Unity: decoupling state and graphics](https://discussions.unity.com/t/decoupling-state-graphics-while-sequentially-executing-graphic-updates/247744)
   3. [Devlog: separating the core game loop](https://chademm.itch.io/rootless/devlog/1121926/separating-the-core-game-loop)
   4. [itch.io: animation queue and Hearthstone](https://itch.io/post/6780327)
   5. [GameDev.net: animations in client-server turn games](https://nene.gamedev.net/forums/topic/698237-how-to-design-animations-and-turn-timers-in-the-clientserver-architecture-of-a-turn-based-game)
   6. [XState](https://github.com/statelyai/xstate)

## Considered options

1. **Claims, an inventory, and an e2e clash check**
   1. Won. It finds the bug class itself, not one pair.
   2. Small: one module, one test, one fixture; no library.
2. **A state machine for all animations (XState)**
   1. Lost. Playback is already a queue with fixed steps.
   2. Its states would not show two systems on one card.
   3. A new library is against ADR 008's short list.
3. **Only events animate; drop the view's pop-in**
   1. Lost. Cards no event brings would appear with no motion.
   2. Kept as the rule's spirit: events win over the view.
4. **Motion `layoutId` hand-over between overlay and table**
   1. Lost. Effects live in an overlay and in `.effects-trail`.
   2. Moving them into the table tree is a large rewrite.
5. **Patch each pair, as before**
   1. Lost. It found clashes only in play tests.

## Decision

Use option 1.

1. Events own animations; the view follows.
   1. An effect that brings or shows a card claims it first.
   2. `Card`'s pop-in plays only for cards no effect claimed.
2. Claims live in `frontend/src/animations.ts`.
   1. Jobs: `arrive` (to its place) and `show` (large).
   2. A beat is one update's effects plus its view.
   3. Two `arrive`s of a card in one beat clash.
   4. The same animation may claim again: two copies of a card.
   5. A `show` clashes in the same or the next beat.
3. A clash is logged, kept in `clashes`, and the second one skipped.
4. Every animation site has an `anim: <id>` tag.
   1. `docs/wiki/Animations.md` has one row per id.
   2. Vitest `animation inventory` fails on a missing tag or row.
5. Every e2e test fails on a clash.
   1. Specs import `test` from `e2e/tests/fixtures.ts`.
   2. Its auto fixture reads `window.fourSouls.clashes`.
   3. `fixtures.spec.ts` fails a spec that imports elsewhere.

## Consequences

### Positive

1. A clash fails a test with both animation names.
2. One page lists every animation and what starts it.
3. A leftover animation fails the inventory test.

### Negative and risks

1. Claims cover cards only, not zones or the stack.
   1. A clash on a zone still needs a play test to find.
2. Only effects that e2e tests play are checked for clashes.
3. Vitest now loads CSS (`css: true`) to read the keyframes.
4. A looted or discarded card no longer pops after its fly.

### Neutral

1. Hover and focus CSS transitions are not in the inventory.
2. The rule for new animations is in `Animations.md`.
