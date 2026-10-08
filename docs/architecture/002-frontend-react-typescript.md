# 2. Frontend: React + TypeScript with a DOM table

* **Status:** Accepted
* **Date:** 2026-10-08
* **Authors:** @HardDie

---

## Context

1. Wails v2 ships templates for six frontends.
   1. React, Vue, Svelte, Preact, Lit, Vanilla.
   2. Each in JS or TS.
2. The game table holds at most a few dozen cards.
3. The UI is text-heavy.
   1. Card text, blank fallback cards, tooltips.
   2. History panel, judge log, menus, settings.
4. Visual style: flat Hearthstone-like cards (V-01).
   1. Drag and drop and cute animations wanted.
   2. No particle effects.
5. The server is authoritative (N-02).
   1. The client only renders state and plays events.
6. Agents will write most of the code.
7. The website (A-03) may reuse UI components.
8. Mobile may come later (B-06).
   1. Only Wails v3 supports mobile, still experimental.
9. What others used:
   1. Hearthstone, MTG Arena, Runeterra: Unity; not an option.
   2. Duelyst: Cocos2d-JS canvas inside Electron.
   3. Board Game Arena: plain DOM + Dojo for hundreds of games.
   4. boardgame.io: React, server-authoritative state.
   5. Web solitaire games: PixiJS for many moving sprites.

## Considered options

1. **DOM rendering (HTML/CSS)**
   1. Won. Text, tooltips, menus come for free.
   2. Enough for dozens of cards.
2. **Canvas/WebGL (PixiJS, Phaser)**
   1. Lost. Built for thousands of objects and particles.
   2. Text is hard; menus still need the DOM.
   3. WebKitGTK may silently fall back to slow software WebGL.
3. **React + TS**
   1. Won. Largest ecosystem, stable API.
   2. Agents know it best.
   3. Same model as boardgame.io.
   4. Already scaffolded; in Wails v2 and v3 templates.
4. **Svelte + TS**
   1. Lost. Less code, but Svelte 5 syntax is new.
   2. Agents often mix old and new syntax.
5. **Vue + TS**
   1. Lost. Smaller game-UI ecosystem than React.
6. **Preact, Lit, Vanilla**
   1. Lost. Fewer libraries; Preact and Lit not in v3 templates.

## Decision

Use options 1 and 3.

1. Frontend is React + TypeScript (Vite), as scaffolded.
2. The table is drawn with the DOM.
3. Animations
   1. Animate only `transform` and `opacity`.
   2. Drag and drop for playing cards.
   3. No particles, no canvas layer.
4. The frontend reaches Wails only through a thin bridge layer.
5. Touch-friendly from the start.
   1. Every hover action also works by tap or long press.
   2. Drag and drop uses pointer events, so touch works.
6. Test on Linux (WebKitGTK) early.

## Consequences

### Positive

1. Fast, simple UI with no game engine to learn.
2. Agents produce reliable React code.
3. Components can be reused by the website.
4. A move to Wails v3 or mobile touches only the bridge.

### Negative and risks

1. React re-renders need care during animations.
2. WebKitGTK on Linux animates worse than Chromium or WebKit.
3. Rich effects would need a canvas later; out of scope now.

### Neutral

1. Animation and drag-and-drop libraries are picked in the roadmap.
2. Sources:
   1. https://wails.io/blog/wails-v2-released
   2. https://v3.wails.io/guides/dev/frontend-frameworks
   3. https://v3.wails.io/guides/mobile
   4. https://gigazine.net/gsc_news/en/20230112-open-duelyst
   5. https://en.doc.boardgamearena.com/Create_a_game_in_BGA_Studio:_Complete_Walkthrough
   6. https://html5gamedevs.com/topic/32121-dom-element-vs-pixitext-performance
   7. https://v2.tauri.app/develop/debug/linux-graphics/
   8. https://bugs.webkit.org/show_bug.cgi?id=265048
