# Frontend tests

The frontend is tested like the Go code (ADR 008). Three layers, from fast to slow:

| Layer | Tool | Where | Run |
|---|---|---|---|
| Unit tests | Vitest | `frontend/src/**/*.test.ts` | `make frontend` (with typecheck and lint) |
| Protocol check | Vitest | reads `internal/protocol/testdata/golden.jsonl` | `make frontend` |
| UI tests | Playwright | `e2e/tests/*.spec.ts` | `make e2e` |

Frontend changes gate on `make frontend e2e`, next to the Go gate.

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
* `?anim=off` — start with animations off (until the player changes the setting).
* With `?fake`, `window.fourSouls.store` is the store: tests feed it server messages with `store.getState().receive(env)`, e.g. a burst of updates.
* `?screen=cards` — the card gallery: every card in hand and table size. The fake has a card with an image, one whose image fails, and blank ones.

Every new screen gets a Playwright test. Find elements the way a player does: by role and visible text (`getByRole('button', { name: 'Save' })`), not by CSS classes.

`e2e/tests/game.ts` has `startGame(browser, names)`: a real game with one page per player, and `expectFits(page)`, which fails when the table scrolls or leaves the window. Tests share one server, so each test uses its own host name.

The fake bridge also reads the real card text from `pkg/card_db/data` (no images), so tables show real names, like a public clone does.

Network play runs against the real Go server: `e2e/playwright.config.ts` starts `cmd/server` on port 4799, next to Vite. The fake `Host` binding points at it (`?server=4799`), so one page "hosts" and others join it by address. Use separate browser contexts for separate players: each has its own `localStorage`, like a separate app.
