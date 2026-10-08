# 1. Lessons from similar card game projects

* **Status:** Proposed
* **Date:** 2026-10-08
* **Authors:** @HardDie

---

## Context

1. Four Souls has many interacting card effects.
2. Others have built digital card games with the same problems.
3. We studied them before writing code.
4. No digital Four Souls exists besides Tabletop Simulator tables.
5. Sources studied:
   1. MTG Arena Game Rules Engine (official dev posts).
   2. XMage, Forge, mage-go (open-source Magic engines).
   3. SabberStone, Brimstone (Hearthstone simulators).
   4. boardgame.io (board game framework).
   5. Factorio, OpenTTD, Widelands (desync write-ups).
   6. Four Souls Tabletop Simulator tables, incl. the official one.

## Considered options

1. **Per-card code for every card**
   1. Lost. Forge found it very hard to maintain as cards grew.
2. **General rules engine + card effects as data blocks**
   1. Won. MTG Arena, XMage, mage-go all go this way.
3. **Client-side rules with lockstep sync**
   1. Lost. Desyncs are hard to find; clients can cheat.
4. **Authoritative server; clients only render**
   1. Won. XMage and our N-02 idea.
5. **Wails v3 (beta since 2026-08-02)**
   1. Lost for now. v2 is the stable line; v3 needs a real port.
6. **Wails v2**
   1. Won. Already scaffolded at v2.16.

## Decision

Use options 2, 4 and 6.
Apply the lessons below when the roadmap is written.

1. Rules engine
   1. Engine knows only general rules: phases, priority, stack, death.
   2. Pending actions are written down before they apply.
   3. Card effects rewrite them: prevent, replace, cancel.
   4. Source: MTG Arena "whiteboard" pattern.
2. Cards
   1. Cards are data built from reusable effect blocks.
   2. Only unique cards get their own code.
   3. A generator makes card stubs from scraped data.
   4. Source: mage-go `genset`.
3. Testing
   1. Each card is written test-first.
   2. Most bugs come from card combinations; test those.
   3. Card status report: not started, implemented, tested.
   4. Fuzz tests on random situations with invariants.
   5. Source: XMage bug reports, SabberStone, mage-go.
4. Determinism
   1. RNG seed stays on the server.
   2. Never rely on Go map iteration order.
   3. Same game run twice must give the same result.
   4. Each recorded step stores a state checksum.
   5. Source: boardgame.io, Widelands.
5. Hidden information
   1. One view-filter function builds every view.
   2. Views: player, viewer, judge.
   3. Source: boardgame.io `playerView`.
6. Priority and auto-pass
   1. Auto-skip uses the full allowed-actions list.
   2. That list includes item and character abilities.
   3. Phase stops are worth considering.
   4. Source: MTG Arena player complaints.
7. Animations and timers
   1. The server can run ahead of client animations.
   2. Response timer gets a fixed animation allowance.
   3. Players can change animation speed.
   4. Source: MTG Arena guides, Legends of Runeterra patch 0.9.1.
8. WebAssembly
   1. CI builds the engine for `js/wasm` from day one.
   2. Standard Go wasm is large; TinyGo breaks `encoding/json`.
   3. Keep engine dependencies minimal.

## Consequences

### Positive

1. Known pitfalls are avoided before the first line of engine code.
2. Card work scales: most cards are data, not code.
3. Replays and tests stay reliable through determinism.

### Negative and risks

1. Effect blocks need upfront design before the first card.
2. Some Four Souls cards may not fit the blocks.
3. Lessons come from Magic and Hearthstone; Four Souls differs.

### Neutral

1. Status is Proposed until the roadmap adopts each lesson.
2. Sources:
   1. https://magic.wizards.com/news/mtg-arena/on-whiteboards-naps-and-living-breakthrough
   2. https://pcgamesn.com/magic-the-gathering-arena/mtg-arena-snappiness-flow
   3. https://feedback.wizards.com/forums/918667-mtg-arena-bugs-product-suggestions/suggestions/51384970-game-automatically-passes-priority-ignoring-on-boa
   4. https://slightlymagic.net/forum/viewtopic.php?p=184527
   5. https://slightlymagic.net/forum/viewtopic.php?p=12850
   6. https://pkg.go.dev/github.com/benprew/mage-go
   7. https://github.com/s13n4/SabberStone
   8. https://github.com/HearthCode/Brimstone
   9. https://lightrun.com/answers/boardgameio-boardgame-io-random-api
   10. https://widelands.org/wiki/DebuggingDesyncs
   11. https://factorio.com/blog/post/fff-188
   12. https://invenglobal.com/articles/10550/patch-091-speeds-up-animations-introduces-jinx-and-garen-themed-boards-to-legends-of-runeterra
   13. https://v3.wails.io/blog/wails-v3-beta/
   14. https://wazero.io/languages/tinygo/
   15. https://steamcommunity.com/sharedfiles/filedetails/?id=2501791757
