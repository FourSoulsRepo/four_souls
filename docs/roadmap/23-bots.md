# Step 23. Bots

Goal: standard and hard bots fill the free seats of a network game.

Ideas: N-21, LR-06.
Decision: ADR 033; every sub-step follows it.

Rules for this step (owner, 2026-11-16):

1. Branch: `bots` (was `step-21-bots`).
   1. Numbered 21 until the merge with main (2026-11-17); main's steps 21 and 22 came first, so it is 23.
   2. Its sub-steps 21.N are 23.N; older commits and branch names keep 21.
2. Two lines in two worktrees (owner, 2026-11-16); each runs its sub-steps one after another, each in a fresh subagent (5.21).
   1. Line 1, bot: this branch; 23.5, 23.6, 23.15.
   2. Line 2, server: worktree on branch `step-21-bots-server` from 23.4; 23.7, 23.8, 23.9.
   3. At most two subagents at once; each line edits only its own lines in shared files.
   4. Time-per-move numbers are taken while the other line runs no tests.
   5. Both done: `git merge --no-ff step-21-bots-server`; then 23.10 and 23.13, one at a time.
   6. The main agent plans and checks; it writes no feature code.
   7. It checks each result against ADR 033 before the next one starts.
3. A bot never reads the real game.
   1. `internal/bot` gets games only from `World(i)` (ADR 033, Decision 4).
   2. Never fix a weak bot by a look at hidden cards.
4. No card names in bot code (5.24).
   1. A card the bot plays badly is fixed in the score or the policy, for every such card.
5. A number that differs from ADR 033 is changed in the ADR in the same commit.
6. The game must stay the same for humans.
   1. A table without bots behaves as before; its tests stay unchanged.
7. Gates (owner, 2026-11-16).
   1. Before a sub-step commit: `go build ./...` in the touched modules and `go test` of the touched packages only.
   2. At the end of each line, once: `make build vet test wasm` and Go lint of the four modules, plain and `embed`.
   3. Line 1 end also: full `make bots`, time per move, B/op and allocs/op.
   4. After the merge: all of the above plus `make sims` and `make fuzz-sets FUZZTIME=30s`.
   5. A failure at an end run is fixed in its own commit, named by its sub-step.
   6. Frontend sub-steps also: `make frontend e2e`.
   7. Arena while working: `BOT_GAMES=20` on the needed line-up.
8. UI sub-steps leave screenshots in `.cache/screenshots/21.N/`; never committed.
9. Order: 23.15 runs right after 23.6, in line 1 (rule 2).
10. Code first, check later (owner, 2026-11-16); for 23.6, 23.15 and 23.16, wins over rule 7.7.
   1. A sub-step is code, unit tests and docs; before its commit only the checks of rule 7.1.
   2. Arena in a sub-step: at most one short run, 20 games, one table size; only to see the bot does not fail or hang.
   3. No tuning of weights or thresholds in loops; a missed arena gate does not stop the sub-step.
   4. A miss goes into its "Done when" as "not reached: <number>"; the sub-step is committed.
   5. Stop only on: no build, a failing unit test, a refused allowed intent, or an architecture choice.
   6. Tuning once, after 23.16: full arena run; every miss in one list to the lead.
   7. The "Done when" gates still hold: line 1 is not merged until they pass at its end.
11. Order: 23.16 runs after 23.15, in line 1, before the merge.
12. All code first, one check phase after (owner, 2026-11-17); stronger than rules 7 and 10.
   1. The remaining code (23.10, 23.12, 23.13, 23.14, bug 19.197) is written and committed before any arena, e2e or measurement.
   2. Before each commit only: `go build ./...` and `go test` of the touched packages (no `-race`, no `-count`); for the frontend, typecheck and Vitest of the touched files.
   3. No arena, e2e, sims, fuzz, lint or benchmarks until the check phase; its contents come from the lead.
   4. Stop only on a failed build or an architecture choice; a small choice is made and written as one line under its sub-step.
   5. `step-21-bots-ui` may merge into this branch before its e2e; the e2e runs in the check phase.
13. Review fixes without subagents (owner, 2026-11-17); stronger than CLAUDE.md 5.21 and rules 7, 10 and 12: one session does every item in a row; before a commit only `go build ./...` and `npm run typecheck`; every test, e2e, arena and screenshot runs once at the end.

---

### [x] 23.1 Plan

1. Goal: the decisions, before any code.
2. Tasks:
   1. Idea N-21; LR-06 changed.
   2. ADR 033.
   3. Sub-step 11.4 moved here.
3. Done when:
   1. The owner's answers are in ADR 033: votes, challenges, pace, names, memory.

### [x] 23.2 Fast clone

1. Goal: a game copy cheap enough to search on (ADR 033, Decision 3).
2. Tasks:
   1. `Game.Clone` by hand-written deep copy; no JSON.
   2. Test: copy and game have equal `Checksum` and `Save`, on simulated games of every set, and of each challenge.
   3. Test: a change to the copy never changes the game; play 50 steps on the copy, check the game's checksum.
   4. Test by reflection: every field of `Game` and its nested types is in a list of copied fields; a new field fails it.
   5. Benchmark `BenchmarkClone` on a 4-player game with every set.
   6. Note in `docs/wiki/Rules-engine.md`: what a new state field must add.
3. Done when:
   1. The benchmark shows under 0.3 ms; the number is in the commit message.
   2. `make sims` and `make fuzz-sets FUZZTIME=30s` pass.

### [x] 23.3 Imagined game

1. Goal: the game as one seat may imagine it (ADR 033, Decision 2).
2. Tasks:
   1. `Game.Imagine(p PlayerID, seed uint64, known Known) (*Game, error)` in the engine.
   2. `Known`: cards of other hands, tops and bottoms of decks (Decision 2.7 – 2.9).
   3. Test: a known card that is hidden there is dealt there; one that is not is skipped.
   2. List first what a seat cannot see: read `view.go` and every place it hides a card.
   3. Pool hidden cards by kind, shuffle, deal back; zone sizes stay.
   4. Keep what the seat looks at now: `Prompt.Look`, shown cards, peeks.
   5. Reseed the copy's RNG from `seed`.
   6. Test: `ViewFor(p)` of the copy equals the game's, at every step of simulated games, all sets, 2–4 players, each challenge.
   7. Test: `Allowed(p)` of the copy equals the game's, when `p` must answer.
   8. Test: `CheckInvariants` passes on the copy; a random playout of 200 steps on it never panics.
   9. Test: other hands and deck tops differ between two seeds.
   10. Test: the real game's checksum is unchanged.
   11. Section "Imagined game" in `docs/wiki/Rules-engine.md`.
3. Done when:
   1. Every test above passes; the engine still builds for `js/wasm`.

### [x] 23.4 Bot core and arena

1. Goal: a bot that plays whole games by a cheap policy, and a place to measure it.
2. Tasks:
   1. Package `internal/bot`: `Level`, `Input`, `Choose` (ADR 033, Decision 4).
   2. Candidates for discards and orders (Decision 5).
   3. Level `random`: a random candidate.
   4. Playout policy (Decision 6), also usable as a level of its own in tests.
   5. `cmd/bot_arena` and `make bots` (Decision 18.1 – 18.3).
   6. The arena gives each seat `World(i)` from `Imagine`; never the game.
   7. Test: no refused intent in 30 policy games, every set, 2–4 players.
   8. Test: `internal/bot` does not import `internal/server` (a lint rule or an import test).
   9. Wiki page `Bots.md` in the developers group: layout, levels, arena.
3. Done when:
   1. Policy against policy, 50 games each, within 150 game turns: 2 players 96%, 3 players 100%, 4 players 92%; median 42 / 50 / 58 turns.
   2. No loop guard hit, no refused intent.
   3. The 95% gate is checked on the standard level in 23.5 (lead, 2026-11-16); policy line-ups are a reference only.

### [x] 23.5 Score and standard level

Gates by the latest numbers: against random 86%, 17.5 ms a move; finished ×2 96%, ×3 98% (tuning, before 7937f18; not measured after it), ×4 100% (after 7937f18).

1. Goal: the standard bot (ADR 033, Decisions 7 and 8).
2. Tasks:
   1. `Score(g, p)` with its weights in one file.
   2. Standard: one world, apply, settle, score, pick; one mistake in ten.
   3. Unit tests on built situations: it takes a soul it can take; it does not attack at hopeless odds; it does not discard its only loot play for nothing.
   4. Test: equal seeds give equal moves.
   5. Benchmark: time per move.
3. Done when:
   1. `make bots`: standard wins at least 70% of 2-player games against random.
   2. Under 100 ms per move on average (lead, 2026-11-16; was 50 ms).
   3. `make bots`: standard×2, ×3 and ×4 each finish 95% of games with a winner within 150 turns; the finish gate is switched on for standard.
   4. Short runs, 20 games each (2026-11-16, machine shared with line 2): standard wins 18 of 20 against random (90%); finished ×2 100% (median 30 turns), ×3 100% (median 42), ×4 90% (median 61; seed 15 is a table The Harbingers locks, seed 9 a long turn of real plays, no loop); deaths a kill 1.20, 1.44, 2.16; 11–41 ms a move under load.
   5. Full run at the end of line 1 (2026-11-17, `make bots`, 50 games a line-up, machine shared): standard wins 43 of 50 (86%) against random. Finished: ×2 96% (48 of 50, median 39 turns); ×3 not reached: 94% (47 of 50, median 60; seeds 10, 23, 35 at the cap); ×4 not reached: 90% (45 of 50, median 66; seeds 10, 11, 15, 18, 42). Deaths a kill ×2 1.69, ×3 1.86, ×4 2.07 (policy 2.22, 2.32, 2.38). Turns with no attack: standard 47 – 54%, policy 17 – 18%. Seed 13, 2 players: standard×2 not finished at 150 turns (souls 2 and 3; 2.97 deaths a kill); hard×2 won in 12 turns.
   6. Time per move: `BenchmarkStandardMove` 6.7 – 7.5 ms, 2.9 – 3.3 MB, 42 000 – 46 000 allocations; in the arena, 12 games side by side, 10 – 25 ms.
   7. Tuning (2026-11-17), partial, not measured in full: attack bar 0 on an empty stack (other moves 1), 12 dice draws, `wItem` 15. Bar 0 for every move with 12 draws: standard×3 98% (49 of 50, median 40), but seed 28 hit the turn guard (items used in a loop), so other moves got 1 back. Attacks only, before `wItem` 15: standard×3 98% (49 of 50, median 31, max 150; seed 11 at the cap), 1.83 deaths a kill, 28% of turns with no attack, 1.7 purchases a seat, 73 ms a move (machine shared). ×4, seed 13, random, flat, hard and `wItem` 15 not measured.
   8. Check phase, before 7937f18 (2026-11-17, 30 games, 12 cores, machine shared): standard×4 29 of 30 finished (96.7%, median 44 turns, max 110); seed 28 hit the game guard (30000 intents at turn 110, bug 19.199); 1.93 deaths a kill; 34.8 ms a move.
   9. After 7937f18 (2026-11-17, e5797b5, 30 games, seeds 1–30 with seed 28, 12 cores, no other job; load average 2.6 at the start, up to 15.8 from the arena itself): standard×4 30 of 30 finished (100%), median 49 turns, max 127; no loop guard, nothing refused; 1.79 deaths a kill; 31% of turns with no attack; 17.5 ms a move; 5 min 33 s.

### [x] 23.6 Hard level

1. Goal: the hard bot (ADR 033, Decision 9).
2. Tasks:
   1. Worlds, playouts to the end of the turn, halving of candidates.
   2. Budget by time and by playout count; `ctx` ends it with the best so far.
   3. The leader's weight in the score (Decision 7.6).
   4. Unit tests on built situations: it answers a stack item that would win the game for another seat, when it can; it hits the leader rather than the last seat, all else equal.
   5. Test: a cancelled `ctx` returns an allowed intent within 50 ms.
   6. Race test: `go test -race ./internal/bot`.
   7. Built situation (owner via lead): 3 seats; hard holds one answer that can cancel a monster kill. The weak seat (0 souls, no items) attacks a monster with no soul while the strong seat (3 souls, items) moves later: hard passes. Same table, the strong seat attacks a monster with a soul: hard answers. If the score cannot do it without a plan over several turns, no planner: report what fails.
   8. Test: with `GOMAXPROCS=1` another goroutine gets the core at least every 20 ms while hard thinks.
   9. Measure playouts per second, B/op and allocs/op of a hard move; write them into ADR 033 Decision 9.8.
   10. Arena: one short run of 20 games at one table size, to show hard neither crashes nor hangs (owner, 2026-11-16). The sweep at 50, 150 and 500 playouts a move and one hard among three standard wait for the tuning after 23.16.
3. Done when:
   1. `make bots`: hard wins at least 60% of 2-player games against standard.
   2. One hard among three standard wins over 30%.
   3. A missed gate: tune, write the numbers into ADR 033, or stop and tell the owner.
   4. Short run, 20 games, 2 players, 150 playouts a move (2026-11-16): hard wins 14 of 20 against standard (70%); all 20 finished (median 25 turns, max 58); no refused intent, no loop guard; 1.36 deaths a kill; 0.31 s a hard move under load.
   5. One hard among three standard (end of line 1, 50 games, 150 playouts): hard wins 32 of 50 (64%); 48 finished (96%), median 49 turns; 1.69 deaths a kill.
   6. The built situations pass on 8 worlds each at 150 playouts: the win answered, the weak seat's kill passed, the strong seat's soul answered (3 souls and 2), the leader hit.
   7. Full run at the end of line 1 (2026-11-17, 50 games, 2 players): hard wins against standard 32 of 50 (64%) at 50 playouts a move, 34 of 50 (68%) at 150, 41 of 50 (82%) at 500; finished 100%, 98%, 100%. Hard×4, 20 games: all finished, median 48 turns, 1.25 deaths a kill. Turns with no attack, same games at 150: hard 35%, standard 52%.
   8. Time: `BenchmarkHardMove` 0.43 – 0.61 s, 196 – 198 MB, 3.8 – 3.9 million allocations a move (machine loaded); `BenchmarkPlayout` 2.7 – 2.9 ms; in the arena 0.19 s a move at 50 playouts, 0.56 s at 150, 1.67 s at 500 (ADR 033, Decision 9.8.6).
   9. Check phase, before 7937f18 (2026-11-17, 150 playouts, 12 cores, machine shared): hard wins 12 of 20 (60%) against standard, 20 of 20 finished (median 23 turns, max 47), 1.34 deaths a kill, 0.64 s a hard move; hard×4 not run (time).
   10. After 7937f18 (2026-11-17, e5797b5, 150 playouts, 12 cores, no other job; load average 8.3 at the start from the run before, up to 19.0 from the arena itself):
       1. Hard against standard, 20 games: hard wins 12 of 20 (60%, exactly at the gate); 20 of 20 finished, median 23 turns, max 47; 1.34 deaths a kill; no loop guard, nothing refused; hard 0.27 s a move, standard 13 ms; turns with no attack hard 27%, standard 21%; 7 min 4 s.
       2. Hard×4, 4 games (lead: only whether they finish with no loop or guard): 4 of 4 finished with a winner, median 80 turns, max 115; no loop guard, nothing refused; 2.14 deaths a kill; 0.21 s a move; 31% of turns with no attack. The 4 games ran side by side, so the slowest game took the whole 27 min 55 s (load average 7.4 at the start, 2.7 – 8.5 during).

### [x] 23.7 Bots in the lobby

1. Goal: the host seats bots at a table (ADR 033, Decision 10).
2. Tasks:
   1. `protocol`: `Seat.Bot`, `Seat.BotN`, message `bot`, the allowed levels in the hello answer (`welcome` → `bots`), version + 1.
   2. Hub: add, change and remove a bot; host only; refusals with translated messages.
   3. Start rule: all seats taken, every human ready; a table with no human closes.
   4. The game list counts bot seats as taken.
   5. A host who leaves: the next human becomes the host, or the table closes.
   6. Restart file of a lobby table keeps its bots (N-20).
   7. Dedicated server flag `-bots=off|standard|all`, default `all`; the app's own server always allows every level (lead, 2026-11-16).
      1. `-bot_think` moves to 23.10: nothing reads it before the bot driver and the thinker limit.
   8. A table with bots does not start yet: one guarded refusal, removed in 23.8.
   9. Tests in `hub_test.go` for each rule above.
   10. `docs/wiki/Server.md`: the flag; ADR 006 list of messages.
3. Done when:
   1. A table of 1 human and 3 bots shows as full and waits only for the human.

### [x] 23.8 Bots in the room

1. Goal: bots play a running game (ADR 033, Decisions 11 – 14).
2. Tasks:
   1. `roomSeat` level; a bot seat is never away.
   2. Bot driver `bots.go`: thinker, step number, cancel, forced answer.
   3. Each hard seat's `Memory` is fed on every broadcast with that seat's view and events.
   3. Pace: 0.7 s and 1–2 s, scaled by the room's `second`.
   4. Bans, draft picks, start-of-game choices.
   5. Tests: a game of 1 human and 3 bots reaches its end; the human only passes and ends turns.
   6. Tests: no bot move while paused, during a rollback vote, after the game end.
   7. Tests: a rollback of a bot's move drops its open thinker; the bot thinks again.
   8. Tests: a kick vote and a Soft vote count only humans.
   9. Tests: a refused bot intent falls back to the forced answer and is logged.
   10. `go test -race ./internal/server`.
   11. No room request over the 100 ms stall limit comes from a bot.
3. Done when:
   1. `TestBotsGame` plays 20 seeded games of 1 human and 3 bots, standard and hard, to the end.

### [x] 23.9 Saves, restarts, records

1. Goal: a game with bots survives a save and a server restart (ADR 033, Decision 16).
2. Tasks:
   1. `savedSeat` and `record.Seat` get the level and number.
   2. A loaded table seats its bots; the game starts when the humans are back.
   3. A kept game goes on after a restart; bots think again.
   4. "No player connected" counts humans only.
   5. Record format: a reader of the old format still reads new files, and the other way; ADR 007 notes the field.
   6. Tests next to `restart_test.go`, `records_test.go` and the save tests.
3. Done when:
   1. A restart in the middle of a bot's turn: the game goes on with no notice to players on the same engine.

### [~] 23.10 Server limits and load

Not [x] (2026-11-17): Done when 1 asks for 20 tables of hard bots on one core; only 3 were run (lead), see 3.1.6.

1. Goal: bots never starve a public server (ADR 033, Decision 17).
2. Tasks:
   1. One limit on hard thinkers for the whole hub.
   2. `cmd/server_load`: tables with bots; report room stalls and move times.
   3. Wiki `Load-test.md`: how to run it with bots, and the numbers measured.
   4. Engine: `checkMaxHP` takes 57% of a standard move (23.5 profile); cache it, or recompute only when an HP bonus changes.
   5. Measure the longest bot turn in seconds with the server's pauses (Decision 13): 40 or more paid abilities at 1–2 s each keep humans waiting over a minute (23.5 arena).
   6. Room latency (lead, 2026-11-16): 182 ms was measured under the load of two lines.
      1. Measure again on a quiet machine, and separately with `GOMAXPROCS=1`.
      2. Numbers taken under load never go into the wiki.
3. Done when:
   1. 20 tables of 1 human and 3 hard bots: no room stall over 100 ms; every bot answers before the response timer.
      1. Check phase (2026-11-17, e5797b5 with a temporary log, not committed): `GOMAXPROCS=1 GOMEMLIMIT=512MiB`, `server_load -players 4 -bots LEVEL -start N -max-games N -game-timeout 8m`; each run lasted 8 minutes (lead), so unfinished games count as failed; no other job on the machine (load average 1.3 – 3.2). Numbers taken under load earlier were removed (rule 2.6.2).
      2. 3 tables of 1 client and 3 hard bots: response median 5.1 ms, p99 15 ms, max 31 ms; no room request over 100 ms in the server log. Longest bot move 3.1 s; bot turn median 47 s, max 1 min 25 s. Peak RSS 34 MB; GC 24 795 runs, 0.6 s of pauses, 6% of the CPU. The server took 99% of the core.
      3. Hard bots: 854 moves; 488 (57%) found no thinker slot within the budget and played standard. Every search played at least 386 playouts (median 1032) and tried every candidate; 71% chose a pass or an end of turn.
      4. 20 tables of 1 client and 3 standard bots: response median 1.2 ms, p99 58 ms, max 140 ms; no room request over 100 ms in the server log, so the 140 ms was most likely a wait for the core. Longest bot move 293 ms; bot turn median 21 s, max 1 min 2 s. Peak RSS 95 MB; GC 1751 runs, 68 ms of pauses, 1% of the CPU. The server took 39% of the core.
      5. Every bot answered far before the response timer (60 s): longest move 3.1 s.
      6. Not met as written: 20 tables of hard bots were not run (lead: 3 hard tables, 20 standard tables). Three hard tables already keep one core 99% busy, so 20 would only queue for slots and play mostly standard. The step stays open for this gate: the owner decides whether 3 hard tables on one core replace it.
4. Decisions (rule 12, 2026-11-17):
   1. `checkMaxHP`: one pass lists the cards in play and their HP statics; each player and monster sums its boosts and those statics (`hpBonus`, the same sum as `bonus`); the last max HP is updated in place, objects that left the check are deleted. No cache across steps: a static's `Applies` and `Value` read the game, so a cache would need every change to invalidate it.
   2. Thinker limit: one channel of `max(1, GOMAXPROCS/2)` slots in the hub; a hard thinker that gets none within its budget plays the move as standard (no slot), so it answers in about twice the budget.
   3. `-bot_think`: a Go duration string (`Config.BotThink`, config `bot_think`, ENV `FOUR_SOULS_BOT_THINK`), 100 ms to 10 s; compose and `.env.example` carry it.
   4. Hard's allocations, the simple part: `eachPart`, behind every stat read, walked a new `inPlay()` slice and a new `facesOf` slice for each card; it now walks them in place (`eachInPlay`, `eachFace`; `inPlay` and `facesOf` build their lists from them, so the order is defined once). The playout copy (about 240 KB) is not reused: that needs a copy into an existing game, which the lead ruled out.
   5. Measure in the check phase: tasks 2.2 (`cmd/server_load` with bots), 2.3 (wiki Load-test), 2.5 (longest bot turn), 2.6 (latency, `GOMAXPROCS=1`, `GOMEMLIMIT`), and the Done when; also the effect of the `checkMaxHP` change and hard's allocations.
   6. Check phase (lead, 2026-11-17): measure the share of hard moves with ≤ 2 candidates and the time spent on them (639e595 plays the full budget); if notable, cut the budget for few candidates.
      1. Measured in the 1-CPU load test (machine shared): 348 of 702 searches (50%) had ≤ 2 candidates and took 529 s of the 1066 s of search. Again on a quiet machine (8 minutes, Done when 1.3): 170 of 316 (54%), 257 of 477 s.
      2. Lead (2026-11-17): not cut; idea N-22 in `docs/ideas.md`.
   7. A hard bot with no thinker slot within its budget plays that move as standard (55be6e6); written in ADR 033, Decision 17.3.1.

### [~] 23.11 Lobby screen

Not [x] (2026-11-18): the full `make e2e` passed in Chromium (287, check phase); the full WebKit run did not pass (19.219, 19.220); 19.197 is solved.

1. Goal: the host adds bots in the app (ADR 033, Decision 10).
2. Tasks:
   1. First, before any new work: `make frontend` and `make i18n` on the frontend edits of 23.7 (`protocol.ts`, `protocol.test.ts`, locales).
   2. `TableScreen`: on a free seat, the host sees "Add bot" with the two levels.
   3. On a bot seat: its name, level, robot icon; the host changes the level or removes it.
   4. Other players see the bot and its level, no buttons.
   5. Buttons hidden when the server does not allow bots.
   6. Every string in every locale; `make i18n`.
   7. Unit tests in `lobby.test.ts`; e2e: host a table, add 3 bots, start alone.
   8. Works without the private images.
   9. The lobby never shows a bot as away or disconnected: `Seat.Connected` is false for a bot.
3. Done when:
   1. Use case "Play with bots" passes as an e2e test.
      1. UC-47; `e2e/tests/bots.spec.ts`: the host adds 3 bots, changes and removes one, starts alone, and the bots throw their first-player dice.
   2. waits: WebKit e2e (libavif13, owner), 19.197.
      1. Lead (2026-11-17): not closed, and the branch not merged, until a full `make e2e` passes, WebKit too.
      2. WebKit (2026-11-17, libavif13 installed with the owner's leave; 5 workers): run 1 218 passed, 92 failed, 12 skipped, 4 did not run (30.6 min, load up to 50 from another session's tests); run 2 295 passed, 18 failed, 12 skipped, 1 did not run (24.2 min).
      3. Every failure but one passed again with 1 worker: they time out under load (19.220). `playback.spec.ts:1975` fails in WebKit every time, also alone (19.219).
      4. No failure was in the bot code: all 5 tests of `bots.spec.ts` passed in both WebKit runs. 19.197 is solved (owner).
      5. Evidence: `.cache/e2e-evidence/23-webkit/` (not in git).

### [~] 23.12 Bots at the game table

Not [x] (2026-11-18): it waits for a passing full WebKit e2e, as 23.11; Chromium passed in full; `bots.spec.ts` passed in both WebKit runs (23.11, 3.2).

1. Goal: everyone sees which seats are bots (N-21).
2. Tasks:
   1. Seat name: translated "Bot N (level)" with the robot icon, everywhere a seat name shows: table, history, stack, pause and vote windows, game end.
   2. A bot seat never shows "away" or disconnected: `Seat.Connected` is false for a bot.
   3. e2e: one human and three bots play a whole game; reuse the helpers of `games.spec.ts`.
   4. `make e2e-games` gets a variant with bots.
   5. Check the pace by eye on screenshots or a recording; note it in the commit.
   6. Layers and animations stay as they are; a new element follows ADR 012 and ADR 015.
3. Done when:
   1. `make frontend e2e` passes; screenshots in `.cache/screenshots/21.12/`.
   2. waits: check phase (e2e incl. WebKit). Owner (2026-11-17): code first, tests later; pace: check phase.
      1. WebKit (2026-11-17): `bots.spec.ts`, 5 of 5 passed in both full runs, the game table test (`:155`) among them; the full run did not pass (23.11, 3.2).
   3. Before the stop: `bots.spec.ts`, Chromium, 4 of 4 passed twice; English screenshots in `.cache/screenshots/21.12/en/`; qps not taken.
4. Decisions (agent, 2026-11-17):
   1. The robot goes before a bot name wherever a line names it: `namePieces` finds the names in the finished text, so no translated line changes.
   2. Also in the prompt banner, the stack's wait lines, start-of-game waits, a look at a hand, choice options, dice labels and banners, not only the places listed.
   3. The turn banner is drawn by hand (`turn.ts`); it builds the same robot span; no new animation, no new layer.
   4. The robot keeps to the name's first word; the rest of the name may wrap.
   5. The right column's name keeps clear of the corner Menu by a float spacer: a long name wraps beside and under it; a short name looks as before. The first screenshot showed "Bot 3 (Standard)" under Menu.
   6. The whole-game e2e with bots lives in `games.spec.ts` only (`make e2e-games`): at the bots' real pace it is too slow for `make e2e`; `bots.spec.ts` plays until a bot's turn instead.
   7. No server flag to speed the bots up for tests.
   8. Every game of `games.spec.ts` now fails when the table pauses: none of them has a player who leaves.
   9. The bots game seats standard, hard, standard.

### [x] 23.13 Bots in challenges

1. Goal: a bot as the second player of a challenge (ADR 033, Decision 15).
2. Tasks:
   1. Team score for challenge games.
   2. The hub lets a bot sit at a challenge table.
   3. Arena: each challenge at each difficulty, bot with bot; wins and losses printed.
   4. Test: no refused intent in challenge games.
   5. e2e: start Krampus with one bot.
   6. `docs/wiki/Challenges.md`: a line on bots.
3. Done when:
   1. Two hard bots win Krampus on normal more often than two random bots; both numbers in the commit message.
4. Decisions (rule 12, 2026-11-17):
   1. The hub already let a bot sit at a challenge table: no hub change; `TestChallengeBot` seats a hard bot at Krampus and starts it.
   2. Boss HP weight `wBossHP` 2 (four times `wHP`); hard scores a challenge by the team sum too, without the holding value.
   3. The policy treats the partner like itself on a curse choice.
   4. Arena: a `-challenge` flag (2-seat line-up, every challenge and difficulty), not in `make bots`; no rooms or deck ratios; no gates but refusals and guards.
   5. Test for refused intents: `TestChallenges`, random seats whole games, standard and hard two turns (8 playouts), so plain `go test` stays fast.
   6. The e2e (Krampus with one bot) is written after `step-21-bots-ui` merges, as the bot seat's UI is there; the Done-when numbers come in the check phase.
   7. Check phase (2026-11-17): the e2e is `bots.spec.ts` "Krampus starts with one bot as the partner": the challenge table, the bot's robot, and the bot's own turn on the turn badge; it passed in the full Chromium run.
5. Check phase (2026-11-17, `bot_arena -games 10 -challenge`, hard at 150 playouts, 12 cores, machine shared):
   1. Two hard bots win Krampus normal 1 of 10, hard 0 of 10, ultra 0 of 10; Trick or Treat normal 1 of 10, hard 0 of 10; ultra not run (time). Two random bots win 0 of 10 at every challenge and difficulty. So hard wins Krampus normal more often than random, by one game of ten.
   2. How they lose: on the timer. Krampus normal games end at turn 16 (median 16, max 16), the timer's 8 rounds of 2 seats; hard median 14, max 16; with 0.64 – 0.96 deaths a kill, not on deaths. Ultra ends by turn 8 – 9, before the timer: there the loss comes earlier (its reason not read out).
   3. No refused intent, no loop guard, no panic. Hard 0.29 – 0.62 s a move.

### [~] 23.14 Docs and review

Not [x] (2026-11-18): the full gate run of task 5 (WebKit e2e did not pass: 23.11, 3.2) and the owner's game. Task 6 done.

1. Goal: the step is ready for the owner.
2. Tasks:
   1. Use case "Play with bots" in `docs/use-cases/` and its index.
   2. Wiki: `Bots.md` final; a players' page on playing with bots; `Home.md` and `_Sidebar.md`.
   3. `README`, changelog.
   4. ADR 033 matches the code: every number that moved.
   5. `cardgen -report` unchanged; `make build vet test lint wasm frontend e2e bots`.
   6. `wails build -m -tags embed` without the private images.
      1. Done (2026-11-17, libgtk-3-dev and libwebkit2gtk-4.0-dev installed with the owner's leave): wails v2.17.0, at b1546638; built in 26 s; `build/bin/FourSouls` 17 211 144 bytes (17M). It changed only file modes in `frontend/wailsjs`, reverted.
   7. Rule for new card sets in `Bots.md` and the developer wiki: after cards are added, `make bot-values`, then `make bots` (ADR 033, Decision 20).
3. Done when:
   1. **Owner** plays a game against three hard bots and approves the merge.
4. Waits: check phase, owner's game.
5. Done now (2026-11-17, non-UI parts):
   1. `Bots.md`: sections "A new card set" and "In the server" (flags, thinker slots, pace, memory), the arena's and `bot_values`' flags, hard's time budget and the numbers before 23.10 marked as such.
   2. `Server.md`: the stale line-1 paragraph replaced; `Dedicated-server-on-a-VPS.md`, the compose page and `deploy/docker` already carried `-bots` and `-bot_think` alike: no change.
   3. `Adding-a-card.md`: `make bot-values`, then `make bots`, in the card steps, the fan set and the checklist.
   4. ADR 033 checked Decision by Decision against the code: 9.8.2.1, 9.8.4, 9.9.1, 17.4, 18.3, 18.6.2, 20.4.4, 20.10 corrected or completed; numbers that wait for the check phase say so.
   5. `README.md` and `CHANGELOG.md` (Unreleased): bots in network games, the server flags; the lobby's part as one line.
6. After the UI merge:
   1. Task 1: UC-47 "Play with bots" comes from `step-21-bots-ui`; only its index entry and links are checked here, no second use case.
   2. Task 2: the players' wiki page on playing with bots, in the players' group of `Home.md` and `_Sidebar.md`; `Bots.md` and `README.md` link to it.
   3. Tasks 5 and 6: the check phase.
   4. Done (2026-11-17): tasks 1 and 2. UC-47 is in the index; it links the players' page and `Bots.md`. `Playing-with-bots.md` (adding bots in the lobby, standard against hard, the robot at the table) is in the players' group of `Home.md` and `_Sidebar.md`; `Bots.md` and `README.md` link to it. 23.14 stays open for tasks 5 and 6 and the owner's game.
7. Decisions (rule 12, 2026-11-17):
   1. No new developer or admin wiki page: the server part of the bots goes into `Bots.md` ("In the server") and links to `Server.md` and the operators' options table, so Home and `_Sidebar` stay as they are now; only the Home line of Bots names its new sections.
   2. The measure keeps 6 draws a loot play (`lootDraws`) after standard's 12, as it keeps the flat values: the table stays as it was measured (ADR 033, 20.4.4.1).

### [x] 23.15 Memory of the hard bot

1. Goal: the hard bot remembers cards it saw (ADR 033, Decision 19).
2. Tasks:
   1. First list every way a seat sees a card that then hides: read `view.go`, `events.go`, `FilterEvents`.
   2. `world.Memory` (ADR 033, Decision 19.2) with `Observe(view, events)`; fed with the seat's own view and events only.
   3. Learn and forget rules of Decision 19.7 – 19.9.
   4. Missing event data: add the field in the engine, filtered for other seats; test the filter.
   5. `Memory.Known()` goes into `Imagine`; the arena keeps one per hard seat.
   6. Test: in simulated games of every set, every remembered card is truly where the memory says, at every step.
   7. Test: a memory built from seat A's feed never names a card seat A's views and events did not show.
   8. Test: a shuffle empties the deck's memory.
   9. Unit tests on built situations: it remembers a looked-at deck top; a card shown from a hand; a discard taken back.
3. Done when:
   1. Test 6 finds no wrong memory in 100 games.
   2. `make bots`: hard with memory is not weaker than hard without; both numbers in the commit message.
   3. Test 6 (`TestMemoryIsTrue`): 100 games of every set and challenge, 49 449 steps, every seat checked at every step; 30 622 remembered cards out of sight checked, none wrong.
   4. The ways a seat sees a card, what it forgets, and the engine fields: ADR 033, Decision 19.14 – 19.17.
   5. Short run, 20 games, 2 players, 150 playouts a move, seeds 1–20 (2026-11-16): hard with memory wins 18 of 20 against standard (90%); without memory (23.6, same seeds) 14 of 20. All 20 finished (median 18 turns, max 70); no refused intent, no loop guard; 1.16 deaths a kill; 0.48 s a hard move under load.
   6. Full run at the end of line 1 (2026-11-17): hard with memory wins 34 of 50 (68%) against standard at 150 playouts; hard without memory was not run again. `TestMemoryStress` (800 games of 2000 steps) runs from `make sims`.

### [~] 23.16 Card values

Not [x] (2026-11-17): Done when 2 not reached, standard with the table wins 46 of 100 (46%) against flat; the 55% gate is the owner's decision (3.11).

1. Goal: each loot card and item has a measured value in the score (ADR 033, Decision 20).
2. Tasks:
   1. First look for a way tests already put a card into a hand or into play (`enginetest`, card tests); use it if there is one.
   2. Else one helper in `enginetest`, not in the `Game` API: it moves a card from its deck to a seat's hand or play; `CheckInvariants` after it.
      1. Depguard: `internal/server` and `internal/bot` never import it.
   3. `make bot-values` (`cmd/bot_values`): sampled positions from saved arena games; each card is put into a clone.
   4. Loot: in positions where the card can be played, stack positions included: mean of max(0, play − nothing), common seeds, `settledScore` units.
   5. Item: a clone with the item and one without, policy playout to the end of the seat's next turn, common seeds; the `settledScore` difference.
   6. Output: a generated data file, card id → value, n, standard error; never written by hand.
   7. One pass, with the table off (flat loot and item values); no iterations.
   8. Fixed seed: the same cards give the same file; a mode for cards missing from the table only.
   9. `cardValue` (score.go) reads the table, shrunk to the mean of the card's kind by n (one constant); an unknown card gets its kind's mean (loot, active item, passive item).
   10. Same honesty: own hand by card; other hands by the mean.
   11. The answer-holding bonus of 23.6 stays a multiplier on this base (Decision 7.6).
   12. `TestBotRules` stays: no card names in bot code; the table is generated data only.
   13. Test: every card of every set is in the table, or its fallback works.
   14. Test: on a forced discard the bot drops a cheap card, not an expensive one.
3. Done when:
   1. A full run takes at most 30 minutes on all cores of the work machine; else fewer samples a card, and the error says so.
   2. 200 games of 2 players: standard with the table wins at least 55% against standard without it.
   3. Time per move does not grow.
   4. Helper: no test path put a card into a running game's hand; `enginetest.PlaceFromDeck` (one helper, quiet, `CheckInvariants` after); depguard `no-enginetest` keeps it out of `internal/server` and `internal/bot`.
   5. Full run (2026-11-17, 12 cores, i7-8750H, machine shared): 7 min 55 s; 36 flat games to 60 turns give 1934 positions in 4 min 55 s; 329 cards, 32 samples each (3 loot cards fewer, `magic_marker` none: falls back to the loot mean).
   6. Error: median standard error loot 0.79 (max 5.7), active items 4.4 (max 6.4), passive items 3.8 (max 6.7); kind means loot 5.09, active 4.43, passive 7.31.
   7. Short run, 20 games of 2 players: standard with the table wins 11 of 20 (55%) against flat; 100% finished, median 30 turns; 12.7 ms a move for both levels together. The 200-game gate and time per move: item 10.
   8. Scale fixes (lead, 2026-11-17): loot in hand at `wHand` 0.5, items scaled to a mean of 12 (horizon), shrink by precision, active items also by their uses. Full run 11 min: 1934 positions in 4 min 48 s; 32 samples a loot card, 64 an item; median standard error loot 0.79, active 3.05, passive 2.70; scaled means loot 5.09 (in hand 2.55), active 10.32, passive 12.90; τ loot 9.69, active 3.28, passive 5.88.
   9. Short run after the scale fixes, 20 games of 2 players: standard wins 13 of 20 (65%) against flat; 100% finished, median 34 turns, max 89; 20.5 ms a move for both levels together (12.7 ms before the fixes; the machine was shared, measured at the tuning).
   10. Full run at the end of line 1 (2026-11-17): 200 games of 2 players, not reached: standard with the table wins 108 of 200 (54%) against flat; 99% finished, median 34 turns; 1.39 deaths a kill. Time per move did not grow: standard 19.6 ms, flat 23.8 ms in the same games. A seat's game: loot plays standard 27.9, flat 25.6; purchases standard 3.8, flat 6.6.
   11. Check phase (2026-11-17, after the scale fixes and the tuning, 12 cores, machine shared): 100 games of 2 players, not reached: standard with the table wins 46 of 100 (46%) against flat; 100% finished, median 26 turns, max 114; 1.48 deaths a kill. Standard 25.8 ms a move, flat 28.2 ms. A seat's game: loot plays standard 23.6, flat 23.0; purchases standard 3.4, flat 4.5. Not tuned: the gate (55%) is the owner's decision; 108 of 200 (54%) before the fixes, 46 of 100 after.
4. Decisions (rule 12, 2026-11-17):
   1. `flat` keeps its own values (lead): `wItemFlat` 12 and `wLootFlat` 0.75, as before 23.16; the measure (`make bot-values`) scores with them too, so the table is measured as it was, and a tuning of `wItem` or `wLoot` moves neither.
