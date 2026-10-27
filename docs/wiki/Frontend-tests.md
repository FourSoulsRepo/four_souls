# Frontend tests

The frontend is tested like the Go code (ADR 008). Three layers, from fast to slow:

| Layer | Tool | Where | Run |
|---|---|---|---|
| Unit tests | Vitest | `frontend/src/**/*.test.ts` | `make frontend` (with typecheck and lint) |
| Protocol check | Vitest | reads `internal/protocol/testdata/golden.jsonl` | `make frontend` |
| UI tests | Playwright | `e2e/tests/*.spec.ts` | `make e2e` |

Frontend changes gate on `make frontend e2e`, next to the Go gate. `make e2e` runs every test in Chromium and in WebKit (the engine of macOS's WKWebView, and the nearest check for Linux's WebKitGTK).

It runs five Playwright workers at once (`workers` in `e2e/playwright.config.ts`; `make e2e E2E_WORKERS=N` for another number). Each test opens two to four app pages. Before 19.38 the pages of past tests stayed open, and five workers slowed a laptop into 30 s game-start timeouts (19.35); three became the default. Since each test closes its pages, a sweep showed no swapping up to eight workers; past three, flaky tests failed more often, and once they were fixed (19.14 – 19.51) the default became five: about 3.7 minutes for a full run instead of 6.

The automatic fixture in `e2e/tests/fixtures.ts` closes every context a test opened when the test ends. Before that, `startGame`'s pages stayed open until the worker ended: about 175 live game pages per browser by the end of a run, each with a WebSocket and a looping glow animation. The laptop ran out of memory and swapped, and whichever test ran then timed out (19.38). A test that opens its own contexts needs no cleanup of its own.

## Keep it simple

The frontend follows the same KISS rules as the backend:

1. The server is the truth. The client computes no rules; it shows the view and sends intents.
2. Logic goes into plain `.ts` functions or the store, with Vitest tests. Components stay thin: they read the store and draw.
3. One Zustand store (`frontend/src/store.ts`). Tests drive it with `useApp.getState()`, no browser needed.
4. Only `frontend/src/bridge/` talks to Wails. Unit tests mock the bridge with `vi.mock('./bridge', …)`.
5. No layers "for later". A new library needs a reason and a `NOTICE` line.

## Unit tests

A test sits next to its file: `store.ts` and `store.test.ts`. They run in Node, without a DOM. `translations.test.ts` is the one that touches files: `make i18n` runs it on the real language files ([Client](Client#translation-checks)); the checks it calls are tested on test tables in `locale.test.ts`, so no trial language sits in `locales/`.

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

`e2e/tests/game.ts` has `startGame(browser, names)`: a real game with one page per player, and `expectFits(page)`, which fails when the table scrolls or leaves the window. Tests share one server and its game list, so guests join the host's game by its id (`lobbyGame`, 19.31); player names may repeat.

Specs import `test` and `expect` from `e2e/tests/fixtures.ts`, not from `@playwright/test`. Its automatic fixture fails a test when any app page logged an animation clash (ADR 012, [Animations](Animations)); `fixtures.spec.ts` fails a spec that imports `test` from elsewhere. Types like `Page` still come from `@playwright/test`.

A game opens before turn 1, with the first-player roll (15.43). `startGame` throws it for everyone with `throwFirstRoll`, which only presses Throw and then waits until every page is idle; so a test starts in turn 1 with nobody's action spent. `{ roll: false }` leaves the roll to the test (`first-roll.spec.ts`). A fake update a test sends counts its step from the real one (`game.step + 1`): a page acts once per step, so a step it already acted on would show its buttons disabled.

A random deal can change what a test sees. `startGame(browser, names, { bans })` leaves cards out of the game (GS-13): it adds `host_bans` to the host's `create` message through Playwright's WebSocket routing, since the host screen has no ban list yet. A test that checks a fresh monster's stats passes `bans: statMonsters` (Delirium and Stoney, 14.10); `TestStatMonsterBans` in `pkg/rules_engine/cards` proves that list is complete for the Base Game.

Bans are a last resort, not a way to silence a failure. A test bans only the fewest cards it needs, and a comment next to the bans says what each card changes that the test checks. A test never bans a card just because it failed on that card: such a failure may be a real engine or client bug, so find the cause first.

The fake bridge also reads the real card text from `pkg/card_db/data` (no images), so tables show real names, like a public clone does.

`e2e/tests/turn.ts` plays prompts through the UI like a careful player: it answers choices with the first option (picking on the table when it can), clicks the first glowing target when its own card aims with no dialog (15.34, 19.10), discards the first cards, passes when only responding, and attacks when a card forces it before the turn can end. `settle(pages, press, done)` runs it until `done()` holds; `click` and `tap` choose mouse or touch. It never makes a move on a page whose state changed after `done()` began, so a dialog the test waits for is not answered before `done()` sees it (19.7). A `done()` that waits for the active page must also accept a new turn: a card may end the turn (19.8). When a game stalls, its error lists each page's prompt, buttons and dialog.

Text that comes from `en.json` is found with `tr(key)` from `e2e/tests/locale.ts`, not the English words, so the tests also pass with `LOCALE=qps` ([Client](Client#languages)). `tr('history.loot.other', { player: 'Bo', card: 'A Penny!' })` is the whole line. `trRe(key, values)` is a pattern: a value left out matches any text, and a plural line with no `n` matches any of its forms. `lineRe` makes a pattern of a line built of lines (`tr('stack.attack', { line: tr('stack.damage', { n: 2 }) })`), `whole` anchors one, and `anyOf` joins several. Text the server sends is a message too, so `tr(key)` finds it: a refusal is `tr('refuse.loot_your_turn')`, a roll label `tr('event.attack_roll')`. A fake update a test sends gives server text as messages, `{ key: 'ask.yes' }`. Player names stay as given. Card names come from `cardName(ref)`: card_db's name in the tests' language (`b2.json`, as the fake bridge), else English, so `make screenshots LOCALE=ru` passes with Russian names (16.7). A ref is the card's id, `id@N` for a later version. With a language file for `LOCALE`, `tr` reads it too, as the app does. `trIn(lang, key)` gives a line in another language, for a page switched to it (`languages.spec.ts`). Before a UI change is done, `make e2e` passes and so do `make screenshots LOCALE=qps` and `LOCALE=ru`.

Match dialogs and buttons by exact name where a question could contain the same word (`{ name: 'Discard', exact: true }`).

Network play runs against the real Go server: `e2e/playwright.config.ts` starts `cmd/server` next to Vite, one per browser project (Chromium on port 4799, WebKit on 4797; `e2e/tests/ports.ts`), so the two runs never see each other's games. Tests reach theirs with `serverPort()` and `appURL()` (`e2e/tests/server.ts`). The fake `Host` binding points at it (`?server=PORT`), so one page "hosts" and others join it by address. Use separate browser contexts for separate players: each has its own `localStorage`, like a separate app.

## Whole random games

`e2e/tests/games.spec.ts` plays one whole game each for 2, 3 and 4 players, with every card set, by random moves through the UI. It clicks only what a player could: menus, dialogs, glowing targets, Attack, Buy, Pass and End turn. Attacks weigh most, so games move toward a winner; a player mostly passes on others' stack. A game fails when:

* it has no winner after `GAME_TURNS` turns (default 300), or no step for a minute;
* any page logs a console error or warning, or throws;
* the server sends an `error` message or closes a socket;
* a page reloads, which drops its game;
* one round of moves takes over 30 seconds.

It is slow, so `make e2e` skips it. Run it with `make e2e-games` (Chromium only). Every 10 turns a test prints the turn, step, time and each player's souls; at the end, its seed, time, turns, steps and moves. `GAME_SEED=N` repeats the moves' choices; the server's deal stays random.

A failed game leaves what is needed to find the cause. Its match record (the deal and every move) is copied to `.cache/e2e-games/`, since the next run deletes the test servers' folders. Each page gets a screenshot and its store in the test's report; a page that stops answering gives the script stack it is stuck in, through a DevTools session opened with the page (19.12).

The Vite server that Playwright starts watches no files (`FOUR_SOULS_E2E=1`, `frontend/vite.config.ts`). Otherwise a Wails build or an edit elsewhere would reload every test page in the middle of a game.

### Results so far (2026-10-15)

Measured on one Mac, Chromium. Some runs played three games at once, which makes each game a little slower. Random play makes the spread wide.

| Players | Games | Average | Range | Average turns |
|---|---|---|---|---|
| 2 | 13 | 410 s (6.8 min) | 109–835 s | 37 |
| 3 | 10 | 643 s (10.7 min) | 122–1244 s | 46 |
| 4 | 8 | 878 s (14.6 min) | 307–1547 s | 58 |

So one round of `make e2e-games` (one game per player count, one after another) takes about 30 minutes, and longer when a 4-player game drags. Some games stay near a winner for dozens of turns: cards take souls away again. Steps also get slower as a game grows long (about 9 steps per second early in a 4-player game, under 5 after turn 80); nobody has looked at why yet.

### Open items

* **19.12** (`docs/roadmap/19-bugs.md`): in one serial run, 3 of 12 games hung for an hour in `page.keyboard.press('Escape')`, as if a page stopped taking input. Their seeds and records were lost. The cause is unknown: no page reload (the owner confirmed no build ran then), and no endless loop found by reading the code. The spec now catches a hang in 30 seconds and keeps its evidence, but more than 20 games since have not hung again. Next step: run rounds until one hangs, then read its stack and record.
* **19.13**: the look dialog lists one card twice when one update looks at it twice. The server sends the events of its auto-passes in the same update (`after` in `internal/server/flow.go`), and `lookedAt` in `frontend/src/history.ts` keeps every event. React logs a duplicate key, and the spec fails on it. Cause found; not fixed yet.

### Working on it

* Run one setup: `cd e2e && FULL_GAMES=1 npx playwright test games --project chromium -g "4 players"`.
* More games: add `--repeat-each N`; `--workers 3` plays three at once. Each repeat hosts under its own names (`Rnd<players>r<repeat><Name>`), so its progress lines stay apart.
* Progress lines start with the host's name, so games played at once stay apart.
* A failed game's record lands in `.cache/e2e-games/`. Replay it in Go with `engine.Replay` or by applying each step's intent to `engine.NewGame`; the header's `setup` (seed, sets, characters) rebuilds the deal.
* `e2e/test-results/` and the test servers' folders are emptied when a run starts: copy anything you need before the next run.
* A Wails build during a run no longer breaks it: e2e Vite watches no files. Before that setting, the owner's build at 11:10 reloaded every test page and failed three games.
