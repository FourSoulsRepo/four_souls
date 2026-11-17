# Step 22. Test reliability

Goal: fewer bugs reach play tests; e2e tests fail only for a real reason.

Ideas: TS-07, TS-08, TS-09, TS-10, TS-11, TS-12.

## Why (owner, 2026-11-16)

1. Step 19 had 163 items by 2026-11-16.
   1. 71 were e2e test faults; the app was right (groups I, J).
   2. 18 were state the engine never sent, or the client never drew.
   3. 3 showed only in long games (19.135, 19.137, 19.138).
2. No card bug was found by a random e2e deal.
   1. Card bugs came from play tests, records and the load tool.
   2. Random deals stay in fuzz, `make sims` and `make e2e-games`.

## Rules for this step (owner, 2026-11-16)

1. Branch: the step 18 branch, `step-18-play-test`; not its own.
2. A failed test is investigated every time, as before (CLAUDE.md 5.16).
   1. No retry rule; no "fails twice" threshold.
3. Each gap found by a check here is a step 19 item.
   1. Logged before its fix, under the step 19 rules.
   2. Fixed game-wide (CLAUDE.md 5.24, 5.26).
4. Gates: `make lint`; `make frontend e2e` for frontend or e2e changes.
5. Sub-steps run in fresh subagents, one at a time (CLAUDE.md 5.21).

### [x] 22.1 ADR: fixed deals in e2e

1. Goal: one ADR fixes how a test gets its deal.
2. Tasks:
   1. ADR 032: a test-only `seed` create option, a seed per test.
3. Done when:
   1. ADR 032 is Accepted and in `INDEX.md`.

### [x] 22.2 Fuzz: every state change has an event

1. Goal: a state change with no event fails the fuzz run (TS-09).
2. Tasks:
   1. Snapshot each object's zone, damage and charge before a step.
   2. After the step, every changed object must be named in its events.
   3. Run it where `fuzz_test.go` calls `CheckInvariants`.
   4. A change silent on purpose goes on a list with its reason.
   5. Hidden zones: a hidden move still needs its hidden event.
   6. Log each gap found as a 19.N item; do not fix it here.
   7. Wiki: the check in `docs/wiki/Fuzzing.md`.
3. Done when:
   1. A test with a silent `DeckToBottom` fails the check.
   2. `make fuzz` runs the check; the gaps found are logged.

### [x] 22.3 Fix the event gaps from 22.2

1. Goal: every gap logged by 22.2 is fixed, game-wide.
2. Tasks:
   1. One 19.N item per kind of gap; fix each under the step 19 rules.
   2. The client plays each new event, or lists why not.
   3. Remove each fixed gap from the 22.2 list.
   4. Fix 19.177 (redo replay differs); `make fuzz-sets` fails on it.
   5. Snapshot max HP; "max HP fell" matches only a real fall.
3. Done when:
   1. The list holds only changes silent on purpose.
   2. `make fuzz` passes.

### [x] 22.4 The client reads every view field

1. Goal: a field the server sends and the client never reads fails a test (TS-07).
2. Tasks:
   1. Go test: every JSON field of the view types is in `protocol.ts`.
   2. Vitest: every view field in `protocol.ts` is read in `frontend/src`.
   3. Reads in `protocol.ts` itself and in tests do not count.
   4. A field unread on purpose is listed with its reason.
   5. Log each other unread field as a 19.N item; do not fix it here.
   6. Wiki: the checks in `docs/wiki/Frontend-tests.md`.
3. Done when:
   1. Removing the read of `copy` in `table.ts` fails the Vitest check.
   2. `make test frontend` runs both checks.

### [x] 22.5 Fix the unread fields from 22.4

1. Goal: every field logged by 22.4 is drawn, or listed with a reason.
2. Tasks:
   1. One 19.N item per field; fix each under the step 19 rules.
   2. A screen change gets before and after screenshots (19 rule 4).
   3. Owner (2026-11-16): draw all eight, 19.180 – 19.187.
   4. Forced attacks (19.180 – 19.182): a mark and the turn banner.
      1. The monster that must be attacked gets a mark on the table.
      2. The banner names it and the attacks left: "You must attack: X (2 left)".
      3. A forced attack on the deck marks the monster deck the same way.
      4. Reuse the existing banner and mark styles.
   5. Deck-only extra attack (19.183): the attack count says it is deck only.
   6. Next attack helpers and sharers (19.184, 19.185): shown for everyone.
      1. Also check that a running attack shows its helpers (`AttackView`).
   7. Roll result (19.186): the stack row shows the number.
   8. Kicked seat (19.187): the player shows as gone for good.
   9. A new window or mark follows ADR 015 (layers) and ADR 012 (animations).
   10. Text goes through i18n keys; `make i18n` passes.
3. Done when:
   1. The 22.4 list holds only fields unread on purpose.
   2. `make frontend e2e` passes.

### [x] 22.6 e2e: animations off by default

1. Goal: e2e pages start at speed "off" (TS-10).
2. Tasks:
   1. `startGame` sets every page to "off" (`setSpeed`, `game.ts`).
   2. An option turns animations on for a test that checks them.
   3. Find every test that checks an animation, sound or flight; give it the option.
   4. `games.spec.ts` keeps its own `GAME_SPEED`.
   5. Wiki: the default and the option in `docs/wiki/Frontend-tests.md`.
3. Done when:
   1. `make e2e` passes, both browsers.
   2. `make screenshots` still shows the animation shots.

### [x] 22.7 e2e: a fixed deal per test

1. Goal: each e2e test sees the same deal every run (TS-11, ADR 032).
2. Tasks:
   1. Server: `-test_seeds` flag, `Config.TestSeeds`, create option `seed`.
   2. Without `TestSeeds`, a non-zero `seed` is refused; Go test for both.
   3. Playwright's server command passes `-test_seeds`.
   4. `startGame` sends a hash of the test's title path as the seed.
   5. Option `seed` overrides it; `seed: 0` asks for a clock seed.
   6. `games.spec.ts` uses `GAME_SEED` or `seed: 0`, as now.
   7. A test that fails with its new deal: find why, as CLAUDE.md 5.16 says.
   8. Wiki: `docs/wiki/Frontend-tests.md` and the server flags page.
3. Done when:
   1. Two runs of one test log the same seed and the same deal.
   2. `make e2e` passes, both browsers.

### [x] 22.8 e2e: named waits on the store

1. Goal: waits check the exact state, not a sign of it (TS-08).
2. Tasks:
   1. List every wait in `e2e/tests` that checks a button, text or class.
   2. Where it stands for game state, use a named wait in `turn.ts`.
   3. A named wait reads the store: step, stack, prompt, whose turn.
   4. `settle` answers any prompt on the way, table targets too.
   5. Wiki: the named waits in `docs/wiki/Frontend-tests.md`.
3. Done when:
   1. No game-state wait checks a DOM sign alone.
   2. `make e2e` passes, both browsers.

### [x] 22.9 Long games stay fast

1. Goal: a whole game fails if a page grows or slows with its length (TS-12).
2. Tasks:
   1. `games.spec.ts` reads each page's element count at the game's end.
   2. It reads the page's script time per step, early and late in the game.
   3. Measure today's values first, over a few games and sets.
   4. Limits: about twice today's values; write them in the spec with that note.
   5. Over a limit, the test fails and names the page and the value.
   6. Wiki: the limits in `docs/wiki/Frontend-tests.md`.
3. Done when:
   1. `make e2e-games` passes with the limits.
   2. A test page with a growing list over the limit fails the check.

### [ ] 22.10 The event check covers players and whole games

1. Goal: a player's state change with no event fails too (owner, 2026-11-16).
2. Tasks:
   1. Snapshot each player's damage, cents and soul count before a step.
   2. A changed player must be named in an event (`Player` field).
   3. 19.2 was this kind: a player's HP set with no `healed`.
   4. Run the check in `TestSimulations` (`make sims`) too.
   5. Fix the object 0 hole: an unset `Object` must not name object 0.
   6. Log each gap found as a 19.N item; fix them here, game-wide.
   7. Wiki: `docs/wiki/Fuzzing.md`.
3. Done when:
   1. A silent change of a player's cents fails the check.
   2. `make fuzz` and `make sims` pass with the check.
