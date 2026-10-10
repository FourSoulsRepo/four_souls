# 8. UI libraries

* **Status:** Proposed
* **Date:** 2026-10-10
* **Authors:** @HardDie

---

## Context

1. ADR 002 picked React + TS with a DOM table.
2. The client needs three things React lacks.
   1. A store for server views and the event queue.
   2. Animations: card moves, flips, damage, dice.
   3. Drag and drop that works with touch.
3. The client never changes game state on its own.
   1. The store holds what the server sent.
   2. Plus local UI state: screen, zoom, settings.
4. Events play one at a time (V-08).
   1. Each animation must report when it ends.
5. Cards move between zones: hand, table, stack, discard.
   1. Moving a DOM node between parents breaks plain CSS transitions.
   2. A shared-element ("FLIP") animation solves it.
6. Linux uses WebKitGTK.
   1. Animate only `transform` and `opacity` (ADR 002).
   2. HTML5 native drag is weak there and on touch.
7. Few libraries: each one goes into `NOTICE`.
8. Candidates checked on 2026-10-10 (npm).

| Library | Latest | Last release | License |
|---|---|---|---|
| zustand | 5.0.15 | 2026-08 | MIT |
| jotai | 3.0.1 | 2026-09 | MIT |
| motion | 14.1.0 | 2026-10 | MIT |
| @dnd-kit/core | 6.3.1 | 2024-12 | MIT |
| @dnd-kit/react | 0.5.0 | 2026-09 | MIT |
| @atlaskit/pragmatic-drag-and-drop | 4.0.0 | 2026-09 | Apache-2.0 |

## Considered options

1. **Store: Zustand**
   1. Won. Tiny, one hook, no providers.
   2. Plain functions can write it (the bridge, the event queue).
   3. Selectors keep re-renders small.
2. **Store: React only (`useReducer`, context)**
   1. Lost. Context re-renders every reader on each update.
   2. The bridge would need hand-made subscriptions.
3. **Store: Redux Toolkit**
   1. Lost. More code for the same job.
4. **Store: Jotai**
   1. Lost. Atoms suit scattered state.
   2. Ours is one server view; a single store fits better.
5. **Animation: Motion (ex Framer Motion)**
   1. Won. `layout` and `layoutId` move a card between zones.
   2. Animations return promises; the queue awaits them.
   3. Uses transform and opacity; hardware-accelerated.
   4. Has pointer-based drag; touch works.
6. **Animation: Web Animations API only**
   1. Lost. Zone-to-zone moves need hand-written FLIP code.
   2. Kept for small effects where Motion is overkill.
7. **Animation: GSAP**
   1. Lost. Not an open-source license.
8. **Animation: react-spring**
   1. Lost. No shared-element moves between parents.
9. **Drag and drop: Motion drag + own drop zones**
   1. Won. No extra library.
   2. Pointer events; mouse and touch alike.
   3. Drop target found with `document.elementFromPoint`.
   4. Our drops are few: play area, targets, order lists.
10. **Drag and drop: dnd-kit**
    1. Lost. `@dnd-kit/core` has had no release since 2024-12.
    2. `@dnd-kit/react` is still 0.x.
    3. Two drag systems would fight over the same card.
11. **Drag and drop: Pragmatic drag and drop**
    1. Lost. Built on native HTML5 drag.
    2. Touch and WebKitGTK support are weak.

## Decision

Use options 1, 5 and 9.

1. Store: Zustand.
   1. One store for the connection, view and event queue.
   2. Small slices for UI state and settings.
   3. Only the bridge and the queue write server data.
2. Animation: Motion.
   1. Animate `transform` and `opacity` only.
   2. Card moves use `layoutId`.
   3. The queue awaits each animation's promise.
   4. Speed setting (ST-05) scales every duration.
3. Drag and drop: Motion drag.
   1. A small hook finds the drop zone under the pointer.
   2. A drop only sends an intent the server allows.
   3. Click or tap to play works too, without drag.
4. Versions
   1. Pin exact versions in `package.json`.
   2. Pick a release at least two weeks old.
   3. Add each to `NOTICE` in the same commit.
5. Nothing else for now.
   1. No router: the screen lives in the store.
   2. No component kit: own CSS.

## Consequences

### Positive

1. Two runtime libraries, both MIT and widely used.
2. Card moves between zones come almost for free.
3. Mouse and touch share one drag code path.

### Negative and risks

1. Motion is the biggest piece (~750 KB unpacked).
   1. Tree-shaking keeps the bundle much smaller.
   2. A desktop app loads it from disk anyway.
2. Motion ships major versions often.
   1. Pinned versions; upgrade on purpose.
3. Own drop zones: we maintain them.
   1. Keyboard drag is not covered; click to play is.
4. Layout animations on WebKitGTK must be checked early (7.12).

### Neutral

1. The website (A-03) can reuse the same stack.
2. A move to Wails v3 or mobile does not touch these choices.
