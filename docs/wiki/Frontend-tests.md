# Frontend tests

The frontend is tested like the Go code (ADR 008). Three layers, from fast to slow:

| Layer | Tool | Where | Run |
|---|---|---|---|
| Unit tests | Vitest | `frontend/src/**/*.test.ts` | `make frontend` (with typecheck and lint) |
| Protocol check | Vitest | reads `internal/protocol/testdata/golden.jsonl` | `make frontend` |
| UI tests | Playwright | `e2e/tests/*.spec.ts` | `make e2e` |

Frontend changes gate on `make frontend e2e`, next to the Go gate. `make e2e` runs every test in Chromium and in WebKit (the engine of macOS's WKWebView, and the nearest check for Linux's WebKitGTK).

## Keep it simple

The frontend follows the same KISS rules as the backend:

1. The server is the truth. The client computes no rules; it shows the view and sends intents.
2. Logic goes into plain `.ts` functions or the store, with Vitest tests. Components stay thin: they read the store and draw.
3. One Zustand store (`frontend/src/store.ts`). Tests drive it with `useApp.getState()`, no browser needed.
4. Only `frontend/src/bridge/` talks to Wails. Unit tests mock the bridge with `vi.mock('./bridge', …)`.
5. No layers "for later". A new library needs a reason and a `NOTICE` line.

## Unit tests

A test sits next to its file: `store.ts` and `store.test.ts`. They run in Node, without a DOM.

## UI tests

Vite serves the frontend with fake Wails bindings (`?fake`, `frontend/src/fakeBridge.ts`). The fake keeps settings in `localStorage`, so a page reload acts like an app restart.

URL flags, dev server only:

* `?fake` — fake Wails bindings.
* `?nosplash` — skip the 5-second fan-game notice. Players always see it (L-01); one test checks it.
* `?firstrun` — start without a nickname, as on the first start; otherwise the fake starts as "Player".
* `?anim=off` — start with animations off (until the player changes the setting).
* With `?fake`, `window.fourSouls.drop()` fails the connection as the network would, to test reconnecting.
* With `?fake`, `window.fourSouls.store` is the store: tests feed it server messages with `store.getState().receive(env)`, e.g. a burst of updates.
* `?screen=cards` — the card gallery: every card in hand and table size. The fake has a card with an image, one whose image fails, and blank ones.

Every new screen gets a Playwright test. Find elements the way a player does: by role and visible text (`getByRole('button', { name: 'Save' })`), not by CSS classes.

`e2e/tests/game.ts` has `startGame(browser, names)`: a real game with one page per player, and `expectFits(page)`, which fails when the table scrolls or leaves the window. Tests share one server, so each test uses its own host name.

Specs import `test` and `expect` from `e2e/tests/fixtures.ts`, not from `@playwright/test`. Its automatic fixture fails a test when any app page logged an animation clash (ADR 012, [Animations](Animations)); `fixtures.spec.ts` fails a spec that imports `test` from elsewhere. Types like `Page` still come from `@playwright/test`.

A random deal can change what a test sees. `startGame(browser, names, { bans })` leaves cards out of the game (GS-13): it adds `host_bans` to the host's `create` message through Playwright's WebSocket routing, since the host screen has no ban list yet. A test that checks a fresh monster's stats passes `bans: statMonsters` (Delirium and Stoney, 14.10); `TestStatMonsterBans` in `pkg/rules_engine/cards` proves that list is complete for the Base Game.

Bans are a last resort, not a way to silence a failure. A test bans only the fewest cards it needs, and a comment next to the bans says what each card changes that the test checks. A test never bans a card just because it failed on that card: such a failure may be a real engine or client bug, so find the cause first.

The fake bridge also reads the real card text from `pkg/card_db/data` (no images), so tables show real names, like a public clone does.

`e2e/tests/turn.ts` plays prompts through the UI like a careful player: it answers choices with the first option (picking on the table when it can), clicks the first glowing target when its own card aims with no dialog (15.34, 19.10), discards the first cards, passes when only responding, and attacks when a card forces it before the turn can end. `settle(pages, press, done)` runs it until `done()` holds; `click` and `tap` choose mouse or touch. It never makes a move on a page whose state changed after `done()` began, so a dialog the test waits for is not answered before `done()` sees it (19.7). A `done()` that waits for the active page must also accept a new turn: a card may end the turn (19.8). When a game stalls, its error lists each page's prompt, buttons and dialog.

Match dialogs and buttons by exact name where a question could contain the same word (`{ name: 'Discard', exact: true }`).

Network play runs against the real Go server: `e2e/playwright.config.ts` starts `cmd/server` next to Vite, one per browser project (Chromium on port 4799, WebKit on 4797; `e2e/tests/ports.ts`), so the two runs never see each other's games. Tests reach theirs with `serverPort()` and `appURL()` (`e2e/tests/server.ts`). The fake `Host` binding points at it (`?server=PORT`), so one page "hosts" and others join it by address. Use separate browser contexts for separate players: each has its own `localStorage`, like a separate app.
