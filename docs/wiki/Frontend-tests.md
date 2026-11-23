# Frontend tests

The frontend is tested like the Go code (ADR 008). Its layers, from fast to slow:

| Layer | Tool | Where | Run |
|---|---|---|---|
| Unit tests | Vitest | `frontend/src/**/*.test.ts` | `make frontend` (with typecheck and lint) |
| Protocol check | Vitest | reads `internal/protocol/testdata/golden.jsonl` | `make frontend` |
| Update fields | Go test, Vitest | `view_fields_test.go`, `viewFields.test.ts` | `make test`, `make frontend` |
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

The inventory tests (`animations.test.ts`, `layers.test.ts`, `i18n.test.ts`) list the source files with `import.meta.glob`. Vitest's on-disk module cache (`fsModuleCache`, on by default since Vitest 5) would keep an old file list, so a new file went unchecked; `vite.config.ts` turns that cache off (19.124).

## Every update field is read

The client must read every field the server sends in a game update, or say why not (TS-07, 22.4). Before this check, the engine sent a field and no screen drew it, and nothing failed: Guppy's Eye's revealed hands (19.5), bonus soul counters (19.29), There's Options' top card (19.30). Two tests keep it so:

1. `TestViewFieldsInProtocolTS` (`internal/protocol/view_fields_test.go`, `make test`) walks `protocol.Update` by reflection: the view, the events, the prompt, the seats, every struct inside, by their `json` names (`-` skipped; `engine.Msg` is a text, not walked). Next to it, it reads the same path in `frontend/src/protocol.ts`, starting at `interface Update`. Each field Go sends must be declared there. A small reader of its own does that: one member per line, inline object types (`look?: { deck?: number }`), and a type that names one interface of the file (`(CardView | null)[]`).
2. `viewFields.test.ts` (`make frontend`) parses `protocol.ts` with the TypeScript compiler and lists every field reachable from `Update`, as `Interface.field` (`Prompt.look.deck` for an inline type). Each must be read in app code under `frontend/src`, or be in its `unread` list with a reason.

A read is the field's name in a property access (`c.copy`, `c?.copy`), an element access with a literal (`c['copy']`) or a destructuring (`{ copy }`). Reads in `protocol.ts`, in tests and in the dev-only fakes (`fakeBridge.ts`, `fakeReplay.ts`) do not count. Writing a field in an object literal is not a read.

The check goes by name, not by type. A common name such as `id`, `kind`, `card` or `amount` counts as read when any object's field of that name is read; `Seat.connected` passes on the lobby's reads alone. So a pass means "the name is read somewhere", not "this field is drawn". It catches a field whose name nothing reads, as most new fields are; a precise check would need the type checker over the whole app.

The `unread` list has two kinds of entries. A field unread on purpose says why: the engine's own bookkeeping (`Turn.entered`), or data that reaches the player another way (`Turn.bonus_attacks_used`: the attacks left are on the player). A gap names its step 19 item (`19.180: …`) and stays until it is fixed; 22.5 drew the eight 22.4 found (19.180 – 19.187), so only fields unread on purpose are left. An entry whose field is gone or now read fails the test, so the list stays current.

Scope: the update message only. Lobby, setup and rollback request messages are not checked. A new field in a Go view type fails the Go test until `protocol.ts` has it, then the Vitest check until the app reads it or lists why not.

## UI tests

Vite serves the frontend with fake Wails bindings (`?fake`, `frontend/src/fakeBridge.ts`). It fakes each bound struct where Wails puts it, `window.go.<package>.<Struct>`, e.g. `window.go.assets.Assets` (ADR 030). The fake keeps settings in `localStorage`, so a page reload acts like an app restart.

URL flags, dev server only:

* `?fake` — fake Wails bindings.
* `?nosplash` — skip the 5-second fan-game notice. Players always see it (L-01); one test checks it.
* `?firstrun` — start without a nickname, as on the first start; otherwise the fake starts as "Player".
* `?anim=off` — start with animations off (until the player changes the setting).
* With `?fake`, `window.fourSouls.drop()` fails the connection as the network would, to test reconnecting.
* With `?fake`, `window.fourSouls.store` is the store: tests feed it server messages with `store.getState().receive(env)`, e.g. a burst of updates.
* `?screen=cards` — the card gallery: every card in hand and table size. The fake has a card with an image, one whose image fails, and blank ones.

Every new screen gets a Playwright test. Find elements the way a player does: by role and visible text (`getByRole('button', { name: 'Save' })`), not by CSS classes.

`e2e/tests/game.ts` has `startGame(browser, names)`: a real game with one page per player, and `expectFits(page)`, which fails when the table scrolls or leaves the window. Tests share one server and its game list, so guests join the host's game by its id (`lobbyGame`, 19.31); player names may repeat. A test that saves and continues a game finds its save by the players' names, so it names them with `repeatNames(names)`: from the second `--repeat-each` repeat on, each name gets the repeat's number, and repeats running at once never see each other's saves (19.152).

Specs import `test` and `expect` from `e2e/tests/fixtures.ts`, not from `@playwright/test`. Its automatic fixture fails a test when any app page logged an animation clash (ADR 012, [Animations](Animations)); `fixtures.spec.ts` fails a spec that imports `test` from elsewhere. Types like `Page` still come from `@playwright/test`.

A game opens before turn 1, with the first-player roll (15.43). `startGame` throws it for everyone with `throwFirstRoll`, which only presses Throw and then waits until every page is idle; so a test starts in turn 1 with nobody's action spent. It throws with animations off, then gives each page its speed back: animated, each throw waits for the dice before it and a tie throws again, which took up to 25 s of a 30 s test under a full run (19.120). The speed changes only in the page's store, not in its saved settings. `{ roll: false }` leaves the roll to the test (`first-roll.spec.ts`). A fake update a test sends counts its step from the real one (`game.step + 1`): a page acts once per step, so a step it already acted on would show its buttons disabled.

Every page of `startGame` plays with animations off (TS-10, 22.6): just before Ready it sets each page's speed to "off" with `setSpeed`, in the store only, and the pages stay so after the first-player roll. Many past failures were checks made while an effect still played (19.37, 19.122, 19.132); with no effects the table shows each update at once. `{ animations: true }` keeps each page's own speed instead: its saved setting ("normal" in the fake), or the url's `&anim=` flag (`&anim=fast`). A test that checks a motion, a sound, a flight, a thrown die, a glow, a banner or any other effect passes it; without it the effect never plays and the check would test nothing. That includes every test in `playback.spec.ts`, and the "on" half of a test run both ways (`stats.spec.ts`, `choice.spec.ts`). A test that checks the "off" look still passes `&anim=off` in its url, as before. `games.spec.ts` passes `animations: true`, then sets its own `GAME_SPEED`. The default covers `startGame` only: a page a test opens itself (lobby, menus, settings, the card gallery, a store-fed table) starts at its saved speed, and passes `&anim=off` when it wants none; such tests mostly check the screens and their own motions. A page that reloads in a game also goes back to its saved speed.

Each test sees the same deal on every run (TS-11, ADR 032). The test servers run with `-test_seeds`, and `startGame` sends a seed in the host's `create` options: `testSeed()`, an FNV-1a hash of the test's title path (`test.info().titlePath`: file, describe blocks, title). So a failure comes back on its first rerun, and both browsers get the same deal. The repeat index is not in the hash: `--repeat-each` plays one deal again and again, so it finds timing faults, not deal faults. `{ seed: N }` picks another deal; `{ seed: 0 }` asks for a random one (the server's clock seed). `startGame` logs the seed when it starts (`<title path>: seed N`) and adds it as the test's `seed` annotation; the record header's `setup` names it too, with the dealt characters. Renaming a test changes its deal. An engine change that draws randomness in another order changes every deal at once: a test that then fails assumed its deal, and is fixed like any test fault (CLAUDE.md 5.16).

The options go into the `create` message inside the host page: an init script wraps `WebSocket.prototype.send` and adds them to a message that starts with `{"type":"create"` (`createOptions`, `game.ts`). Before 22.7, Playwright's WebSocket routing added the bans, which sent every message of the host page through the test runner; with a seed in every test, that would have slowed every host page and set it apart from its guests. `startGame` then reads the table's options from the host's store and checks that the server kept the bans and the seed.

A deal can still hold a card that changes what a test sees. `startGame(browser, names, { bans })` leaves cards out of the game (GS-13): it adds `host_bans` to the host's `create` options the same way, since the host screen has no ban list yet. A test that checks a fresh monster's stats passes `bans: statMonsters` (Delirium and Stoney, 14.10); `TestStatMonsterBans` in `pkg/rules_engine/cards` proves that list is complete for the Base Game. Bans from before fixed deals stay with their reasons: they keep a test right when its deal changes.

Bans are a last resort, not a way to silence a failure. A test bans only the fewest cards it needs, and a comment next to the bans says what each card changes that the test checks. A test never bans a card just because it failed on that card: such a failure may be a real engine or client bug, so find the cause first.

The fake bridge also reads the real card text from `pkg/card_db/data` (no images), so tables show real names, like a public clone does.

`e2e/tests/turn.ts` plays prompts through the UI like a careful player: it answers choices with the first option (picking on the table when it can), picks what to attack, what to buy and a deck attack's slot on the table, where no dialog asks (`tablePicks`, 18.11; a board card before a player; a card counts when nothing covers its middle, by a hit test with no wait, as a trial click waited for the card to stand still and under load gave up on the monster slots for the monster deck, 19.197), clicks the first glowing target when its own card aims with no dialog (15.34, 19.10), discards the first cards, passes when only responding, and attacks when a card forces it before the turn can end. `settle(pages, press, done)` runs it until `done()` holds; `click` and `tap` choose mouse or touch. It never makes a move on a page whose state changed after `done()` began, so a dialog the test waits for is not answered before `done()` sees it (19.7). A `done()` that waits for the active page must also accept a new turn: a card may end the turn (19.8). When a game stalls, its error lists each page's prompt, buttons and dialog.

A test waits for the exact state it needs, not for a sign of it (18.10). After an Attack click, `untilAttackRoll(pages, press, attacker)` runs `settle` until the store shows the attacker's throw prompt for the attack roll, that roll on the stack not yet thrown, the attack aimed at its target, and the throw allowed. Waiting for any Throw button stopped too early when a card rolled first: Bob's Brain rolls when the attack is declared, before the target is picked (19.85). A trigger of the target, such as Pride's, resolves before the first roll goes on the stack (R-ATK-25, 19.151); `settle` answers its prompts on the way, so the wait still holds at the first roll. The check does not need the roll on top. A test that reads the latest move at that throw (Undo) bans Pride: its discard is then the latest move. A new kind of wait gets its own named helper in `turn.ts` the same way. `untilStackWaits(pages)` waits until every page shows a stack item waiting for a response (V-11); a resolving card's question does not count, as its prompt line reads differently for the asked player (19.126).

Text that comes from `en.json` is found with `tr(key)` from `e2e/tests/locale.ts`, not the English words, so the tests also pass with `LOCALE=qps` ([Client](Client#languages)). `tr('history.loot.other', { player: 'Bo', card: 'A Penny!' })` is the whole line. `trRe(key, values)` is a pattern: a value left out matches any text, and a plural line with no `n` matches any of its forms. `lineRe` makes a pattern of a line built of lines (`tr('stack.attack', { line: tr('stack.damage', { n: 2 }) })`), `whole` anchors one, and `anyOf` joins several. Text the server sends is a message too, so `tr(key)` finds it: a refusal is `tr('refuse.loot_your_turn')`, a roll label `tr('event.attack_roll')`. A fake update a test sends gives server text as messages, `{ key: 'ask.yes' }`. Player names stay as given. Card names come from `cardName(ref)`: card_db's name in the tests' language (`b2.json`, as the fake bridge), else English, so `make screenshots LOCALE=ru` passes with Russian names (16.7). A ref is the card's id, `id@N` for a later version. The tests read every language file, as the app does, so `tr` gives `LOCALE`'s line and `trIn(lang, key)` a line in another language, for a page switched to it (`languages.spec.ts`, `settings.spec.ts`). A key shared by two screens is found with the key that screen uses, even where the English is the same: `history.mode_god`, not `host.rollback_god` (12.12). Before a UI change is done, `make e2e` passes and so do `make screenshots LOCALE=qps` and `LOCALE=ru`.

The rules website has its own spec, `site.spec.ts`: it opens the page `make site` built, served by `vite preview` ([Website wasm](Website-wasm#tests)). The site takes its language from the link, and a built page has no `qps`, so the spec reads English lines with `trIn('en', …)`.

Match dialogs and buttons by exact name where a question could contain the same word (`{ name: 'Discard', exact: true }`).

Network play runs against the real Go server: `e2e/playwright.config.ts` starts `cmd/server` next to Vite, one per browser project (Chromium on port 4799, WebKit on 4797; `e2e/tests/ports.ts`), so the two runs never see each other's games. Tests reach theirs with `serverPort()` and `appURL()` (`e2e/tests/server.ts`). The fake `Host` binding (`window.go.host.Host`) points at it (`?server=PORT`), so one page "hosts" and others join it by address. Use separate browser contexts for separate players: each has its own `localStorage`, like a separate app.

### Named waits (22.8)

A wait that stands for game state reads the store, never a button, a text or a class (TS-08). Whose turn it is, what the stack holds, which prompt waits, whether an attack roll waits for its throw: each is in the page's store (`window.fourSouls.store`), and a button or a glow only follows it. Waiting for "a Throw button shows" stopped at the wrong roll (19.85); waiting for "a target glows" or "a dialog shows" missed a card that asked one more question (19.10, 19.18, 19.21, 19.58). The rule: wait for the state through the store, then assert the DOM. A check of what the screen shows (the thing under test) stays a DOM check.

`seen(page)` reads one snapshot: the step, whose turn, the prompt (kind, purpose, text key), the intents the page may send now, the stack (kind, controller, attack roll, thrown), the attack's target, and `ready` (the latest state, nothing sent for this step, no rollback vote open). A `Check` is a condition on it; `holds(page, check)` tests it now, `until(page, check)` polls until it holds and answers nothing, and `settle(pages, press, () => holds(page, check))` answers prompts on the way. The engine's enums are mirrored at the top of `turn.ts` (`Intent`, `Ask`, `Purpose`, `Item`), in Go's order.

| Name | Kind | Holds when |
|---|---|---|
| `may(Intent.x)` | check | the page is ready and the engine allows that intent: `attack`, `purchase`, `pass`, `endTurn`, `throw`, `playLoot` |
| `asked(Ask.x, purpose?)` | check | the game waits for the page's own player on a prompt of that kind (and purpose) |
| `question(key)` | check | the page's player is asked a choice whose question has that text key (`ask.one_on_top`) |
| `stackHolds(n)` | check | the stack holds n items |
| `turnFree` | check | the stack is empty and the page's player may end the turn |
| `attackRoll` | check | the attacker's attack roll waits for its throw (above) |
| `stackWaits` | check | the stack holds an item and waits for a response (V-11) |
| `canEndTurn(page)` | holds | `may(Intent.endTurn)` |
| `untilPick(pages, press, page, purpose)` | settle | the page's player is asked to pick: an attack's target, what to buy, a room |
| `untilAttackRoll(pages, press, attacker)` | settle | `attackRoll` on the attacker's page |
| `untilStackWaits(pages)` | poll | `stackWaits` on every page |
| `untilResynced(page)` | poll | the page reconnected and shows the server's answer (19.103) |
| `activePage(pages, press)` | settle | a page may end the turn; then its End turn button is enabled |

`endTurn(page)` reads from the store whether the button glows (nothing allowed but ending the turn or passing), and waits for the button to show the same before it presses it; a glowing button ends the turn at once, any other asks first (V-46).

A spec may add its own small store check next to these, for state only it needs: a challenge's counters, a flipped card in play, a boost for the next roll. A new kind of wait used by more than one spec goes into `turn.ts`.

## Whole random games

`e2e/tests/games.spec.ts` plays one whole game each for 2, 3 and 4 players, with every card set, by random moves through the UI. It clicks only what a player could: menus, dialogs, glowing targets, Attack, Buy, Pass and End turn. Attacks weigh most, so games move toward a winner; a player mostly passes on others' stack. A game fails when:

* it has no winner after `GAME_TURNS` turns (default 300), or no step for a minute;
* any page logs a console error or warning, or throws;
* the server sends an `error` message or closes a socket;
* a page reloads, which drops its game;
* one round of moves takes over 30 seconds;
* a page grows or slows with the game's length (below, 22.9).

It plays one game of one player and three bots too (23.12, ADR 033): `startGame`'s `bots` option seats a bot of each listed level after the players, through the host's Add bot buttons, so the bots play at their real pace (Decision 13: 0.7 s an answer, 1 to 2 s an action) and the game is slower than one of humans. The random player plays as in the other games; the bots play by themselves. A bot that wins must show its robot in the game end's title. In every game here the table must never pause: nobody leaves, and a bot is never away. Run it alone with `-g "bots"`.

It also plays the Krampus (14.26) and Trick or Treat (14.34) challenges ([Challenges](Challenges)) once each at each difficulty, Normal, Hard and Ultra Hard, for 2 players with every card set. Random play rarely wins one: the timer runs out or the presents are destroyed in a few minutes. The random moves answer the challenges' own dialogs like any other: the forgo and flip questions, the fear pick, the attack-start throws and the trigger order. That is a valid end. The game must reach game over with "Challenge won!" or "Challenge lost" and its reason on both pages; the test logs the result and the timer left. Each test logs its seed when it starts, so a game that never ends still names it.

It is slow, so `make e2e` skips it. Run it with `make e2e-games` (Chromium only). Every 10 turns a test prints the turn, step, time and each player's souls; at the end, its seed, time, turns, steps and moves. Its deals stay random (`seed: 0`, ADR 032): finding surprises is its job. `GAME_SEED=N` repeats the moves' choices and the deal (the same `N` as the server seed). The pages play at the "fast" animation speed: every effect still plays and is claimed, so the clash check holds, at half its time (19.137). `GAME_SPEED=normal` plays at the normal pace.

The host bans the cards the load tool bans by default (`-ban` in `cmd/server_load`): Blue Baby, whose Forever Alone can repeat forever (19.97), and Smelter and Portable Slot Machine, which loop with other cards (19.109). A random player could repeat such a loop and the turn would never end. Keep the list in `games.spec.ts` the same as the load tool's, and drop a ban when its item is fixed.

A failed game leaves what is needed to find the cause. Its match record (the deal and every move) is copied to `.cache/e2e-games/`, since the next run deletes the test servers' folders. Each page gets a screenshot and its store in the test's report; a page that stops answering gives the script stack it is stuck in, through a DevTools session opened with the page (19.12).

The Vite server that Playwright starts watches no files (`FOUR_SOULS_E2E=1`, `frontend/vite.config.ts`). Otherwise a Wails build or an edit elsewhere would reload every test page in the middle of a game.

### Long games stay fast (22.9)

Some bugs showed only after thousands of steps: sounds that kept a media player each (19.135), and a history strip that drew every tile of the game (19.138). So each game measures every page as it plays (`e2e/tests/growth.ts`) and fails a page over a limit. The failure names the page, the value, the limit and the game with its seed, e.g. `Rnd4r0Bob: elements at the end 1850, limit 1700; 4 players finish a game with every card set, seed 856041`. Each test also logs one line per page with its values.

A page is read every 100 steps and at the game's end:

* **Elements**: `document.getElementsByTagName('*').length` at the end, game-over dialog open.
* **Script time per step**: Chromium's own `ScriptDuration` counter (DevTools `Performance.getMetrics`), through the DevTools session each page already opens. Its change over a window, divided by the steps in that window, is the page's script time per step. It counts the page's own tasks: updates, renders, animation frames. It does not count the test's `page.evaluate` reads, so the bot's own polling does not show in it. The early window is steps 200–500, after the deal; the late one is the game's last 300 steps. A game too short for two windows apart (under about 800 steps) gets only the early one.
* **Live sounds**: the most Audio elements that held a media player at once (`fourSouls.sounds.live()`, 19.135). The test build has no sound files, so each play fails and is released at once; it reads 0 or 1. It catches a release that stops working.

Limits, about twice the values measured on 2026-11-17: 15 games over every set, 2–4 players and the challenges, up to 10547 steps; one Mac, Chromium, one game at a time, "fast":

| Measure | Measured | Limit |
|---|---|---|
| Elements at the end | 580–869 | 1700 |
| Script ms a step, early or late | 5.6–12.9 | 25 |
| Growth: late / early | 0.64–1.19 | 2 |
| Live sounds | 0–1 | 4 |

The growth ratio is the main check for "slows with its length": machine load moves both windows alike, so it holds on a slower or busier machine. The absolute limit is a backstop for a page slow from the start; it depends on the machine, so raise it there rather than chase a slow laptop. 19.138 was a growth of about 5 (12 ms a step early, 60 late) at 4700 elements.

`e2e/tests/growth.spec.ts` checks the check (in `make e2e`, Chromium): a plain page whose list grows each step fails both the element and the growth limit; one that stays flat passes.

### Results so far (2026-10-15)

Measured on one Mac, Chromium. Some runs played three games at once, which makes each game a little slower. Random play makes the spread wide.

| Players | Games | Average | Range | Average turns |
|---|---|---|---|---|
| 2 | 13 | 410 s (6.8 min) | 109–835 s | 37 |
| 3 | 10 | 643 s (10.7 min) | 122–1244 s | 46 |
| 4 | 8 | 878 s (14.6 min) | 307–1547 s | 58 |

So one round of `make e2e-games` (one game per player count, one after another) takes about 30 minutes, and longer when a 4-player game drags. Some games stay near a winner for dozens of turns: cards take souls away again. Steps also got slower as a game grew long (about 9 steps per second early in a 4-player game, under 5 after turn 80). The cause was the app: the history strip drew every headline of the game on each update; it now holds the newest 40 (19.138).

### Open items

* **19.12** (`docs/roadmap/19-bugs.md`): in one serial run, 3 of 12 games hung for an hour in `page.keyboard.press('Escape')`, as if a page stopped taking input. Their seeds and records were lost. The cause is unknown: no page reload (the owner confirmed no build ran then), and no endless loop found by reading the code. The spec now catches a hang in 30 seconds and keeps its evidence, but more than 20 games since have not hung again. Next step: run rounds until one hangs, then read its stack and record.
* **19.13**: the look dialog lists one card twice when one update looks at it twice. The server sends the events of its auto-passes in the same update (`after` in `internal/server/flow.go`), and `lookedAt` in `frontend/src/history.ts` keeps every event. React logs a duplicate key, and the spec fails on it. Cause found; not fixed yet.

### Working on it

* Run one setup: `cd e2e && FULL_GAMES=1 npx playwright test games --project chromium -g "4 players"`.
* The bots game alone: `cd e2e && FULL_GAMES=1 npx playwright test games --project chromium -g "bots"`.
* More games: add `--repeat-each N`; `--workers 3` plays three at once. Each repeat hosts under its own names (`Rnd<players>r<repeat><Name>`), so its progress lines stay apart.
* Progress lines start with the host's name, so games played at once stay apart.
* A failed game's record lands in `.cache/e2e-games/`. Replay it in Go with `engine.Replay` or by applying each step's intent to `engine.NewGame`; the header's `setup` (seed, sets, characters) rebuilds the deal.
* `e2e/test-results/` and the test servers' folders are emptied when a run starts: copy anything you need before the next run.
* A Wails build during a run no longer breaks it: e2e Vite watches no files. Before that setting, the owner's build at 11:10 reloaded every test page and failed three games.
