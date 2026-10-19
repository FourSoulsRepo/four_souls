# 15. Window layers

* **Status:** Accepted
* **Date:** 2026-10-19
* **Authors:** @HardDie

---

## Context

1. 19.26: a zoomed card showed under the discard window.
   1. The zoom had z-index 10; the window had 12.
   2. The player could not read the card they discarded.
2. Each z-index was picked alone, in its own CSS file.
   1. 19 values in 6 files, from 1 to 1001.
   2. Nothing said which window must cover which.
   3. A new window got a number that looked free.
3. Other pairs were wrong too.
   1. The action menu's backdrop was above both prompts by chance.
   2. A dragged hand card was above the modal overlays.
4. The owner asked for a list, like the animation list (ADR 012).

## Considered options

1. **One ordered list of CSS variables, checked by a test**
   1. Won.
   2. One file holds every number; the order reads at a glance.
   3. A test ties it to a wiki table with a reason per layer.
   4. Bare numbers fail the test, so no layer hides elsewhere.
2. **A TypeScript enum, set as inline styles**
   1. Lost.
   2. Most layers are CSS rules; inline styles would move them.
3. **Fix the one pair, add a comment**
   1. Lost.
   2. The next window would pick a number blind again.

## Decision

Use option 1.

1. `frontend/src/layers.css` lists every layer, top first.
   1. One variable per layer: `--z-<id>`.
   2. Numbers fall by 10, so a layer fits between two others.
2. Code uses `z-index: var(--z-<id>)`.
   1. A number that orders parts of one element is allowed.
   2. It carries `/* layer: local */` on its line.
3. `docs/wiki/Layers.md` has one row per layer, same order.
   1. Each row says why it covers the rows below.
4. Vitest `layer inventory` fails on:
   1. A bare z-index.
   2. Numbers out of order.
   3. A layer unused, or not defined.
   4. A wiki table that differs from `layers.css`.
5. Pairs a player meets get e2e tests in `layers.spec.ts`.
6. The order follows what the player is looking at.
   1. Pointer feedback on top: the aim arrow and its note.
   2. Game events next: effects never hide.
   3. Then the card zoom, over every window.
   4. Then modal overlays, then windows, then the table.

## Consequences

### Positive

1. The zoom shows over every window (19.26).
2. A new window must say where it goes and why.
3. Wrong pairs fail a test, not a play test.

### Negative and risks

1. The test reads source text; a z-index set in an odd way slips by.
   1. It matches `z-index:` and `zIndex`, which covers today's code.
2. Stacking contexts still matter.
   1. A layer inside a positioned parent with a z-index stacks locally.
   2. The wiki page says so; the e2e pairs catch the rest.

### Neutral

1. Effects stay above the overlays, as before.
2. The aim arrow drops from 1000 to the top of a short list.
