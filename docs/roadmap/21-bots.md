# Step 21. Bots

Goal: standard and hard bots fill the free seats of a network game.

Ideas: N-21, LR-06.
Decision: ADR 030; every sub-step follows it.

Rules for this step (owner, 2026-11-16):

1. Branch: `step-21-bots`.
2. Sub-steps run one after another, each in a fresh subagent (5.21).
   1. The main agent plans and checks; it writes no feature code.
   2. It checks each result against ADR 030 before the next one starts.
3. A bot never reads the real game.
   1. `internal/bot` gets games only from `World(i)` (ADR 030, Decision 4).
   2. Never fix a weak bot by a look at hidden cards.
4. No card names in bot code (5.24).
   1. A card the bot plays badly is fixed in the score or the policy, for every such card.
5. A number that differs from ADR 030 is changed in the ADR in the same commit.
6. The game must stay the same for humans.
   1. A table without bots behaves as before; its tests stay unchanged.
7. Gates before each commit: `make build vet test lint wasm`.
   1. Frontend sub-steps also: `make frontend e2e`.
   2. Bot strength sub-steps also: `make bots`, with its numbers in the commit message.
   3. `make sims` and `make fuzz-sets FUZZTIME=30s`: once at the end of a batch that changed `pkg/rules_engine`, not per commit (owner, 2026-11-16).
8. UI sub-steps leave screenshots in `.cache/screenshots/21.N/`; never committed.
9. Order: 21.15 runs right after 21.6, before 21.7.

---

### [x] 21.1 Plan

1. Goal: the decisions, before any code.
2. Tasks:
   1. Idea N-21; LR-06 changed.
   2. ADR 030.
   3. Sub-step 11.4 moved here.
3. Done when:
   1. The owner's answers are in ADR 030: votes, challenges, pace, names, memory.

### [x] 21.2 Fast clone

1. Goal: a game copy cheap enough to search on (ADR 030, Decision 3).
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

### [x] 21.3 Imagined game

1. Goal: the game as one seat may imagine it (ADR 030, Decision 2).
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

### [x] 21.4 Bot core and arena

1. Goal: a bot that plays whole games by a cheap policy, and a place to measure it.
2. Tasks:
   1. Package `internal/bot`: `Level`, `Input`, `Choose` (ADR 030, Decision 4).
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
   3. The 95% gate is checked on the standard level in 21.5 (lead, 2026-11-16); policy line-ups are a reference only.

### [ ] 21.5 Score and standard level

1. Goal: the standard bot (ADR 030, Decisions 7 and 8).
2. Tasks:
   1. `Score(g, p)` with its weights in one file.
   2. Standard: one world, apply, settle, score, pick; one mistake in ten.
   3. Unit tests on built situations: it takes a soul it can take; it does not attack at hopeless odds; it does not discard its only loot play for nothing.
   4. Test: equal seeds give equal moves.
   5. Benchmark: time per move.
3. Done when:
   1. `make bots`: standard wins at least 70% of 2-player games against random.
   2. Under 50 ms per move on average.
   3. `make bots`: standard×2, ×3 and ×4 each finish 95% of games with a winner within 150 turns; the finish gate is switched on for standard.

### [ ] 21.6 Hard level

1. Goal: the hard bot (ADR 030, Decision 9).
2. Tasks:
   1. Worlds, playouts to the end of the turn, halving of candidates.
   2. Budget by time and by playout count; `ctx` ends it with the best so far.
   3. The leader's weight in the score (Decision 7.6).
   4. Unit tests on built situations: it answers a stack item that would win the game for another seat, when it can; it hits the leader rather than the last seat, all else equal.
   5. Test: a cancelled `ctx` returns an allowed intent within 50 ms.
   6. Race test: `go test -race ./internal/bot`.
3. Done when:
   1. `make bots`: hard wins at least 60% of 2-player games against standard.
   2. One hard among three standard wins over 30%.
   3. A missed gate: tune, write the numbers into ADR 030, or stop and tell the owner.

### [ ] 21.7 Bots in the lobby

1. Goal: the host seats bots at a table (ADR 030, Decision 10).
2. Tasks:
   1. `protocol`: `Seat.Bot`, `Seat.BotN`, message `bot`, bots flag in the hello answer, version + 1.
   2. Hub: add, change and remove a bot; host only; refusals with translated messages.
   3. Start rule: all seats taken, every human ready; a table with no human closes.
   4. The game list counts bot seats as taken.
   5. A host who leaves: the next human becomes the host, or the table closes.
   6. Restart file of a lobby table keeps its bots (N-20).
   7. Dedicated server flag `-bots`; the app's own server always allows them.
   8. A table with bots does not start yet: one guarded refusal, removed in 21.8.
   9. Tests in `hub_test.go` for each rule above.
   10. `docs/wiki/Server.md`: the flag; ADR 006 list of messages.
3. Done when:
   1. A table of 1 human and 3 bots shows as full and waits only for the human.

### [ ] 21.8 Bots in the room

1. Goal: bots play a running game (ADR 030, Decisions 11 – 14).
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

### [ ] 21.9 Saves, restarts, records

1. Goal: a game with bots survives a save and a server restart (ADR 030, Decision 16).
2. Tasks:
   1. `savedSeat` and `record.Seat` get the level and number.
   2. A loaded table seats its bots; the game starts when the humans are back.
   3. A kept game goes on after a restart; bots think again.
   4. "No player connected" counts humans only.
   5. Record format: a reader of the old format still reads new files, and the other way; ADR 007 notes the field.
   6. Tests next to `restart_test.go`, `records_test.go` and the save tests.
3. Done when:
   1. A restart in the middle of a bot's turn: the game goes on with no notice to players on the same engine.

### [ ] 21.10 Server limits and load

1. Goal: bots never starve a public server (ADR 030, Decision 17).
2. Tasks:
   1. One limit on hard thinkers for the whole hub.
   2. `cmd/server_load`: tables with bots; report room stalls and move times.
   3. Wiki `Load-test.md`: how to run it with bots, and the numbers measured.
3. Done when:
   1. 20 tables of 1 human and 3 hard bots: no room stall over 100 ms; every bot answers before the response timer.

### [ ] 21.11 Lobby screen

1. Goal: the host adds bots in the app (ADR 030, Decision 10).
2. Tasks:
   1. `TableScreen`: on a free seat, the host sees "Add bot" with the two levels.
   2. On a bot seat: its name, level, robot icon; the host changes the level or removes it.
   3. Other players see the bot and its level, no buttons.
   4. Buttons hidden when the server does not allow bots.
   5. Every string in every locale; `make i18n`.
   6. Unit tests in `lobby.test.ts`; e2e: host a table, add 3 bots, start alone.
   7. Works without the private images.
3. Done when:
   1. Use case "Play with bots" passes as an e2e test.

### [ ] 21.12 Bots at the game table

1. Goal: everyone sees which seats are bots (N-21).
2. Tasks:
   1. Seat name: translated "Bot N (level)" with the robot icon, everywhere a seat name shows: table, history, stack, pause and vote windows, game end.
   2. A bot seat never shows "away".
   3. e2e: one human and three bots play a whole game; reuse the helpers of `games.spec.ts`.
   4. `make e2e-games` gets a variant with bots.
   5. Check the pace by eye on screenshots or a recording; note it in the commit.
   6. Layers and animations stay as they are; a new element follows ADR 012 and ADR 015.
3. Done when:
   1. `make frontend e2e` passes; screenshots in `.cache/screenshots/21.12/`.

### [ ] 21.13 Bots in challenges

1. Goal: a bot as the second player of a challenge (ADR 030, Decision 15).
2. Tasks:
   1. Team score for challenge games.
   2. The hub lets a bot sit at a challenge table.
   3. Arena: each challenge at each difficulty, bot with bot; wins and losses printed.
   4. Test: no refused intent in challenge games.
   5. e2e: start Krampus with one bot.
   6. `docs/wiki/Challenges.md`: a line on bots.
3. Done when:
   1. Two hard bots win Krampus on normal more often than two random bots; both numbers in the commit message.

### [ ] 21.14 Docs and review

1. Goal: the step is ready for the owner.
2. Tasks:
   1. Use case "Play with bots" in `docs/use-cases/` and its index.
   2. Wiki: `Bots.md` final; a players' page on playing with bots; `Home.md` and `_Sidebar.md`.
   3. `README`, changelog.
   4. ADR 030 matches the code: every number that moved.
   5. `cardgen -report` unchanged; `make build vet test lint wasm frontend e2e bots`.
   6. `wails build -m -tags embed` without the private images.
3. Done when:
   1. **Owner** plays a game against three hard bots and approves the merge.

### [ ] 21.15 Memory of the hard bot

1. Goal: the hard bot remembers cards it saw (ADR 030, Decision 19).
2. Tasks:
   1. First list every way a seat sees a card that then hides: read `view.go`, `events.go`, `FilterEvents`.
   2. `world.Memory` (ADR 030, Decision 19.2) with `Observe(view, events)`; fed with the seat's own view and events only.
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
