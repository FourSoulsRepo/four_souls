# Step 18. Play-test

Goal: real players play real games on every set done so far.
Moved from step 8 (8.4–8.6): the Base Game alone is too few cards.

Ideas: B-04, B-05.

---

### [ ] 18.1 New-user check

1. Goal: users can start the app (B-04, B-05).
2. Tasks:
   1. Give a new user only the release archive.
   2. Watch them start it on macOS, Windows, Linux.
   3. Fix `README.txt` and the wiki where they get stuck.
3. Done when:
   1. A new user starts the app with only the README.

### [ ] 18.2 Play-test and fixes

1. Goal: real games with real players.
2. Tasks:
   1. Play full games with 2, 3, 4 players.
   2. Use every set done in step 14.
   3. Collect issues; fix blockers.
   4. New rule disputes go to the open list.
3. Done when:
   1. **Owner** approves the release.

### [ ] 18.3 Docs check

1. Goal: docs match the released app.
2. Tasks:
   1. Every released feature has a use case.
   2. README is short and current.
   3. Accepted ADRs match the code.
3. Done when:
   1. Use case index has no "Planned" item for released features.

### [ ] 18.4 Multi-instance review

1. Goal: the "several copies" solution (UC-16) is reviewed and kept or reworked.
2. Tasks:
   1. Review `internal/instance`, `main.go`, `app.go` from commit 62a610f.
   2. It was built outside the roadmap; it has no ADR yet.
   3. Check it on macOS (`open -n`), Windows, and Linux.
   4. Check edge cases: crash, all 16 slots busy, two copies hosting.
   5. Check settings seeding: no copy takes another copy's seat.
   6. **Owner** decides: keep, simplify, or replace.
   7. Write an ADR for the result; update UC-16 if behavior changes.
3. Done when:
   1. ADR accepted; UC-16 matches the code.

### [x] 18.5 Sound loudness tool

1. Goal: a tool measures every sound's loudness (AU-04, ADR 022).
2. Tasks:
   1. Work in the separate repo `FourSoulsRepo/tools` (`../FourSoulsRepo/tools`).
   2. Never add this tool or its dependency to the game repo (owner).
   3. New command: `go run ./loudness`, next to `scraper`.
   4. MP3 decoding: `github.com/hajimehoshi/go-mp3`.
      1. Add it to that repo's `NOTICE` (Apache-2.0).
   5. Measure M200 as ADR 022 says.
      1. K-weighting filter coefficients for the file's sample rate.
      2. Use the BS.1770 formulas, as libebur128 does.
      3. Mix to mono first; the sounds are mono anyway.
      4. Pad with silence, so a short sound fills the window.
      5. Also report the sample peak in dBFS.
   6. `loudness check -game ../four_souls`
      1. Reads every `*.mp3` in `pkg/card_db/images/assets/sounds/`.
      2. Prints one row per file: M200, peak, gain needed.
      3. Exit code 1 when any file is off -20 ±1.5 LU or peaks above -1 dBFS.
      4. A file list after the flags checks only those files.
   7. `loudness level -game ../four_souls <file>…`
      1. Needs `ffmpeg`; add it to that repo's README requirements.
      2. Gain to -20 LUFS M200, then a limiter at -1 dBFS.
      3. Re-measure and adjust, up to 3 passes; the limiter lowers loudness.
      4. Writes mono, 44.1 kHz, 64 kbps MP3 in place.
      5. Prints before and after values.
   8. Tests in the tools repo.
      1. A 997 Hz sine at full scale reads -3.01 LUFS.
      2. A 400 ms window matches ffmpeg `ebur128` momentary within 0.1 LU.
      3. A sound shorter than the window still reads right.
   9. Document both commands in that repo's README.
   10. Commit in the tools repo after its own `go vet` and tests pass.
3. Done when:
   1. `check` prints the four values of ADR 022 (200 ms column) within 0.2 LU.
   2. `check` exits 1 on today's sounds.

### [x] 18.6 Level the sounds

1. Goal: every sound plays at one loudness (AU-04, ADR 022).
2. Tasks:
   1. Needs 18.5.
   2. Run `loudness level` on the four files in `images/assets/sounds/`.
   3. Commit them in the private submodule repo `FourSoulsRepo/card_db`.
   4. Bump the submodule pointer in the game repo.
   5. `frontend/src/sounds.ts`: `volume` from 0.6 to 1.0.
      1. Levelled files are 3–10 dB quieter than before.
      2. One volume for every sound; no per-sound gain.
   6. New `CLAUDE.md` rule, appended under Workflow:
      1. Every sound sits at -20 LUFS M200 (ADR 022).
      2. A new or changed sound passes `loudness check` before its commit.
      3. Fix one with `loudness level` in the tools repo.
   7. Wiki: a "Sound loudness" section in `docs/wiki/Card-data-tools.md`.
      1. What M200 is, the target, both commands.
   8. Wiki: the sound paragraph in `docs/wiki/Client.md` names ADR 022.
   9. ADR 022: status Proposed → Accepted; fix any drift.
   10. Gates: `make lint`, `make frontend e2e`.
3. Done when:
   1. `loudness check` passes on every sound.
   2. **Owner** plays a game and hears no sound jump out.
4. Owner notes (2026-10-27)
   1. Owner marks 18.5 and 18.6 solved.
   2. The images and tools repos are pushed.

### [x] 18.7 Loot draw sound

1. Goal: each looted card makes a sound (AU-05).
2. Tasks:
   1. Source: owner's `Card_Deal_1.wav`; source unknown.
      1. So it lives only in the private images, as `loot-draw.mp3`.
      2. Mono MP3, levelled with `loudness level` (ADR 022).
      3. Amend the submodule commit; bump the pointer.
   2. `App.Sounds()` gains `lootDraw`; `''` without the file.
   3. Each `looted` plays it as its card flies.
      1. One sound per card; each card has its own fly.
      2. Every cause: start of turn, monster reward, card effects.
      3. A card taken from the loot discard plays it too (owner).
      4. With animations off, nothing plays, as other sounds.
   4. `NOTICE`, UC-12, wiki `Client.md` name the sound.
   5. Tests: Vitest for playback and sounds; e2e for the played file.
   6. Gates: `make lint`, `make frontend e2e`.
3. Done when:
   1. `loudness check` passes on every sound.
   2. **Owner** hears one deal per looted card.

### [x] 18.8 Player death sound

1. Goal: a player's death makes a sound (AU-06).
2. Tasks:
   1. Source: owner's `player_died.wav`; source unknown.
      1. So it lives only in the private images, as `player-died.mp3`.
      2. Mono MP3, levelled with `loudness level` (ADR 022).
      3. Amend the submodule commit; bump the pointer.
   2. `App.Sounds()` gains `playerDied`; `''` without the file.
   3. It plays with the skull of a player's `died` event.
      1. The engine sends it only when the death resolves (R-DEATH-17).
      2. A death still on the stack can be prevented; it is silent.
      3. A monster's death is silent.
      4. Each death resolves on its own, so each plays its own sound.
      5. With animations off, nothing plays, as other sounds.
   4. `NOTICE`, UC-12, wiki `Client.md` name the sound.
   5. Tests: Vitest for playback and sounds; e2e for the played file.
   6. Gates: `make lint`, `make frontend e2e`.
3. Done when:
   1. `loudness check` passes on every sound.
   2. **Owner** hears it once per player death, never on a prevented one.

### [x] 18.9 Treasure gain sound

1. Goal: each gained treasure flies in with a sound (AU-07).
2. Tasks:
   1. Source: owner's `power up1.wav`; source unknown.
      1. So it lives only in the private images, as `treasure-gain.mp3`.
      2. Mono MP3, levelled with `loudness level` (ADR 022).
      3. Amend the submodule commit; bump the pointer.
   2. `App.Sounds()` gains `treasureGain`; `''` without the file.
   3. Each gained treasure flies to its player, one at a time.
      1. Shop slot buy: `purchased` flies from the slot.
      2. The engine's `purchased` gains `prev`: the slot's object.
      3. The slot shows empty while it flies, as a dead monster's.
      4. Deck top buy, monster reward, card effects: `gained_treasure`.
      5. They fly from the treasure deck.
      6. Each fly plays the sound as it leaves.
   4. Silent, no fly: a loot trinket, a starting item, a stolen item.
   5. With animations off, nothing plays, as other sounds.
   6. `NOTICE`, UC-12, wiki `Client.md`, `Animations.md` name it.
   7. Tests: Vitest for playback and sounds; e2e for the played file.
   8. Gates: `make lint`, `make frontend e2e`.
3. Done when:
   1. `loudness check` passes on every sound.
   2. **Owner** hears one sound per gained treasure, and sees it fly.

### [x] 18.10 e2e waits name the state they need

1. Goal: a card that asks one more input does not break a test (TS-08).
2. Tasks:
   1. `turn.ts`: `untilAttackRoll(pages, press, attacker)`.
      1. It runs `settle` until the attack roll waits for its throw.
      2. `attackRollWaits` reads the store: top roll is the attack's.
      3. The attack has its target; the Throw button shows.
   2. Every "Attack, then wait for any Throw" uses it: ten places.
   3. Wiki `Frontend-tests.md` names it; CLAUDE.md 5.16.10.
   4. Gates: `make lint`; the changed specs in both browsers.
3. Done when:
   1. The changed specs pass; no test waits for any Throw after Attack.
   2. The 19.85 Bob's Brain test passes through the new wait.
4. Result (2026-10-27)
   1. 126 of 126 in the eight changed specs, chromium and webkit.
5. Owner notes (2026-10-27)
   1. Owner marks 18.10 solved.

### [x] 18.11 Attack and buy pick on the table

1. Goal: attack, buy and slot picks happen on the table, no dialog (V-53).
2. Tasks:
   1. Three choices skip the choice dialog:
      1. What to attack (`Purpose.attackTarget`, R-ATK-02).
      2. What to buy (`Purpose.purchase`, R-SHOP-02).
      3. A deck attack's slot (`Purpose.monsterSlot`, R-ATK-08).
   2. Their options glow on the table, the rest dims, as today.
      1. Monsters, shop items, attackable items and characters.
      2. The monster deck or the treasure deck for a deck option.
      3. A deck attack: the drawn card stays face up on the deck (19.9).
   3. A click on a glowing option picks it; no arrow, no confirm.
   4. No cancel: the attack or purchase is declared already.
   5. An option not on the table keeps the dialog (e.g. an empty slot).
   6. The banner says what the game waits for.
      1. Title: the choice's question, e.g. "Choose what to attack".
      2. Hint: a new line, "Click a glowing card".
      3. Other players see "Waiting for Ann" and the question, as today.
   7. Out of scope (owner, 2026-10-27):
      1. A click on a monster before Attack does nothing, as today.
      2. The rules choose the target after the priority window.
   8. Docs: UC-13 (5, 5a), UC-21, wiki `Client.md`; V-53 marked done.
   9. Tests:
      1. Vitest: which choices skip the dialog; the banner hint.
      2. e2e: helpers pick glowing table options without a dialog.
      3. e2e: specs that used these dialogs pick on the table.
   10. Gates: `make lint`, `make frontend e2e`.
3. Done when:
   1. Attack, buy and slot picks show no dialog.
   2. **Owner** checks the screenshots of the three picks.
4. Result (2026-10-27)
   1. `make lint`, `make frontend`: pass.
   2. `make e2e`: 376 passed, 6 skipped, chromium and webkit.
   3. Screenshots: `attack-pick-table`, `turn-actions-buy`, `deck-attack-pick-slot`.
   4. `fd65095` held only these result lines: a test command stashed the work.
   5. The next commit brings the code; an internal build of `fd65095` has none of it.
5. Owner notes (2026-10-27)
   1. Owner marks 18.11 solved.

### [x] 18.12 Random seats

1. Goal: the host can seat the players in a random order (UC-37).
2. Tasks:
   1. Host screen: "Random seats" in Detailed setup, remembered.
   2. The server shuffles the seats at the start, from the game seed.
   3. The table's setup summary says "random seats".
   4. UC-37, UC-03, `CHANGELOG.md`, wiki `Server.md` name it.
   5. Tests: Go shuffle test; Vitest for the summary text.
   6. Gates: `make lint`, `make frontend e2e`.
3. Done when:
   1. A play-test game with Random seats starts in a random order.
   2. **Owner** approves the option's place in Detailed setup.
4. Result (2026-10-28)
   1. Implemented; committed with the 18.12 change on `step-18-play-test`.
   2. `make lint`, `make frontend`, `go test ./internal/...`: pass.
   3. `make e2e`, run 1: 379 passed, 6 skipped; `info.spec.ts:528` (15.90, zoom shows "Cain") failed; alone it passed twice.
   4. `make e2e`, run 2: 379 passed, 6 skipped; `languages.spec.ts:102` (16.4) failed: Eden's starting item came before the first-roll prompt; alone it passed.
   5. `make e2e`, run 3: 378 passed, 6 skipped; `languages.spec.ts:102` (19.98) failed in chromium; `stack.spec.ts:68` (19.100) failed in webkit.
   6. Not solved: the gate is red on all three full runs; the failures are flaky tests, not in the seat code (19.98, 19.99, 19.100).
   7. Fixed in 19.98, 19.99, 19.100 (other commits on this branch).
   8. `make e2e`, run 4 (2026-10-28): 384 passed, 6 skipped, exit 0.
   9. `make lint`, `make frontend`, `go test ./internal/...`: pass.
   10. `wails build -m -tags embed` (no card images): builds.
   11. Done-when 1 is covered by the Go shuffle test; no manual game was played.
5. Owner notes (2026-10-28)
   1. Owner asks to mark 18.12 solved.

### [x] 18.13 Bonus soul flies

1. Goal: a bonus soul flies from its card on the board to its player.
2. Tasks:
   1. Sound: the same `soulGain` sound as any soul (owner, 2026-10-28).
   2. A `gained_soul` with text `bonus soul` flies from its card on the board.
      1. The card is `obj-<id>` in the view before the update; no engine change.
      2. Reverses the 19.73 rule that a bonus soul does not fly.
      3. A bonus soul not on the board before: no flight, silent.
   3. With animations off, nothing plays, as other sounds.
   4. UC-12 (1ag), wiki `Client.md`, `Animations.md` name it.
   5. Tests: Vitest for playback; e2e for the played file.
   6. Gates: `make lint`, `make frontend e2e`.
3. Done when:
   1. **Owner** sees a bonus soul fly to its player, and hears the soul sound.

### [x] 18.14 Lost Soul sound

1. Goal: the Lost Soul makes the soul sound as it becomes a soul in play.
2. Tasks:
   1. Sound: the same `soulGain` sound as any soul (owner, 2026-10-28).
   2. A `gained_soul` with no `prev` for a `Lost Soul Card` plays it.
      1. Nothing flies: the Lost Soul is already on the player's side.
      2. A `sound` effect plays it and holds a beat.
      3. Other souls with no `prev` stay silent.
   3. With animations off, nothing plays, as other sounds.
   4. UC-12 (1ag), wiki `Client.md`, `Animations.md` name it.
   5. Tests: Vitest for playback; e2e for the played file.
   6. Gates: `make lint`, `make frontend e2e`.
3. Done when:
   1. **Owner** hears the sound when the Lost Soul becomes a soul.
4. Result (2026-10-28)
   1. Implemented; `make lint`, `make frontend`, `go test ./...`: pass.
   2. `make e2e`: 384 passed; two webkit tests failed once (19.16, aim) and passed on rerun.
   3. `go build -tags embed .`: builds.
5. Owner notes (2026-10-28)
   1. Owner asks to mark 18.13 and 18.14 solved.

### [x] 18.15 Counter kinds on the table

1. Goal: a player can tell each counter apart, and each card counts only its own kind.
2. Tasks:
   1. Colours (owner, 2026-10-29): done.
      1. General counter: green.
      2. Golden counter: gold.
      3. A new counter kind: the owner picks its colour.
   2. One badge per kind, not one total: done.
      1. The badge shows the kind's colour and count.
      2. The zoomed card names each kind, bottom right: done (owner, 2026-10-29).
   3. Each card counts only its own kind.
      1. Bum-Bo! levels up from general counters only.
      2. Tech X's "Remove 3 counters" pays from general counters only.
   4. Eve's line is a test case (owner, 2026-10-29): done, `TestEveGoldenTrinketLine`.
      1. Eve has Tech X and Bum-Bo!.
      2. Golden Trinket puts a gold counter on Tech X.
      3. Eve dies and loses Golden Trinket.
      4. Eve gets Golden Trinket back in hand.
      5. Golden Trinket puts a gold counter on Bum-Bo!.
      6. Expected: Tech X and Bum-Bo! each show one gold counter; Bum-Bo! is level 0.
   5. Owner answers (2026-10-29):
      1. Tech X's "Remove 3 counters" does not use gold counters.
      2. A gold counter on Bum-Bo! does nothing to Bum-Bo!.
      3. A gold counter on Bum-Bo! or Tech X does count for Golden Trinket: it gains that card's abilities.
   6. Engine tests: `cards/r/golden_trinket_counters_test.go`, `cards/r/golden_trinket_bum_test.go`.
      1. Counters do not mix (Tech X, Bum-Bo!).
      2. Two gold counters work at once: Golden Trinket gains both cards' abilities.
      3. Bum-Bo! adds no ability; one activation uses one ability the player picks.
      4. Gold on Bum Friend and Tech X: both ↷ abilities are the trinket's, and they share its one tap.
      5. Bum-Bo!'s statics reach Golden Trinket and read the trinket's own level.
      6. Gaining ¢ levels Bum-Bo! itself, even with gold on it; the trinket does not level.
   7. Open question, wait for the owner:
      1. Q-27: does Golden Trinket take Bum-Bo!'s gain-¢ replacement? The engine says no; `TestGainCentsReplacementGained` is skipped.
      2. Q-28: may an item use two gained ↷ abilities in one turn? The engine says no.
      3. Both are in `pkg/rules_engine/docs/rules/open-questions.md`.
   8. Client: one badge per kind, coloured; the zoom lists the kinds (`counters.ts`, `Card.tsx`, `CardZoom.tsx`, `Card.css`).
      1. Tests: `counters.test.ts`; e2e "a gold counter is its own badge, in gold (18.15)".
      2. Wiki `Client.md` names the colours.
   9. Gates: `make lint`, `make frontend e2e`, `go test ./...` in `pkg/rules_engine`.
3. Done when:
   1. **Owner** plays Eve's line and reads each counter's colour and count.
4. Result (2026-10-29)
   1. Implemented; `make lint`, `make frontend`, `go test ./...`: pass.
   2. `make e2e`: 390 passed, 6 skipped.
   3. Q-27 and Q-28 stay open in `open-questions.md`.
5. Owner notes (2026-10-29)
   1. Owner asks to mark 18.15 solved.

### [ ] 18.16 Server speed under load

1. Goal: the server handles more games at once, with no slow responses.
2. Tasks:
   1. Profile the server under the load tool (`docs/wiki/Profiling.md`).
   2. Cut the per-step costs the profile shows, keeping the wire format and the record format.
      1. Protocol: encode each message once.
      2. Record: each line is written and flushed by its own goroutine per game; each step is still flushed.
      3. Rollback: one `rollback_last` per update, not the list of targets.
      4. Discard piles: the top card and the size; the full pile is not sent.
      5. Blocked reasons: the turn actions in each update; a card's reasons only when its owner clicks it (`explain`, `explained`).
      6. Checksum: a hash chain over the steps, with the whole state every 100 steps (ADR 025, engine 0.5.0).
   3. Records made before engine 0.5.0 cannot replay (`CanReplay`).
   4. Find the cause of closed connections at 128 games.
      1. The server closes a client whose 64-message send queue is full, without a log line.
      2. Log each such close, then rerun.
   5. Docs: ADR 025, ADR 005, UC-22, `docs/wiki/Client.md`, `docs/wiki/Rules-engine.md`, `docs/wiki/Server.md`.
   6. Gates: `make lint`, `make frontend e2e`, `go test` for the changed modules.
3. Done when:
   1. 64 games at once pass, and the slowest response stays under 1 s.
   2. 128 games at once pass, with no connection closed by the server.
   3. **Owner** reads the before and after profiles.
4. Result (2026-10-29)
   1. Local, same seed, 64 games: median response 17.7 → 6.1 ms; p95 60 → 22 ms; max 349 → 132 ms.
   2. Over the LAN from another computer, 64 games: median 15 ms; max 41 ms (traced build).
   3. 128 games over the LAN: 40 of 128 games closed by `1006`; not yet explained.
   4. Gates: `make lint`, `make frontend`, `make e2e` (390 passed, 6 skipped), `go test` for the changed modules: pass.
   5. One full `make e2e` run failed `rollback.spec.ts` 17.4 in both browsers (Undo, not Ask to undo); it passed alone and in the next full run. Not yet logged as a bug.
   6. `TestStopAtAgainstServer` failed once (closed connection); it passed 5 of 5 reruns.
   7. Open: 128 games; the profile after the change.

### [ ] 18.17 One-CPU public server

1. Goal: a public server on one CPU runs 100 to 200 games at once.
2. Tasks:
   1. Run the tagged server with `GOMAXPROCS=1` and `-trace_dir`.
   2. Load it with `cmd/server_load`: 64, 128, then 200 games.
   3. Record `cpu.pprof`, `allocs.pprof`, `block.pprof` and `trace.out` at each load.
   4. Find the top hot spots by flat and cumulative CPU time.
   5. Cut the top costs; keep the wire format and the record format.
   6. Find each closed connection's cause (18.16.4).
   7. Check memory: `stats.csv` stays flat while games end.
   8. Write the findings: what was tried, what helped, what did not.
   9. Docs: `docs/wiki/Profiling.md`, `docs/wiki/Load-test.md`, `docs/wiki/Server.md`.
   10. Gates: `make lint`, `go test` for the changed modules.
3. Done when:
   1. 200 games at once pass on `GOMAXPROCS=1`; the slowest response is under 1 s.
   2. No connection is closed by the server at 200 games.
   3. **Owner** reads the before and after profiles.

4. Result (2026-10-30, in progress)
   1. Profile, `GOMAXPROCS=1`, 2 games (run 1): 38% in `madvise`, 25% in JSON encode.
   2. Allocations: about 390 GB in 380 s; 87% under `protocol.Encode`.
   3. `GOGC=400`: no gain; 1.30 ms per step against 1.40 ms.
   4. `encoding/json/v2` `MarshalerTo` on `Args`: slower (15.4 µs against 9.4 µs) and other bytes; reverted.
   5. Typed fast path for `Args` numbers: same bytes, no gain; reverted.
   6. fastjson: intent decode 387 ns to 91 ns; the update encode is not in its reach.
   7. easyjson (generated, `omitempty` only): update encode 36.7 µs to 27.9 µs; `omitzero` fails to generate.
   8. Stopped by 19.109: a game loops forever under the load tool.
   9. Benchmark: `internal/protocol/bench_test.go` `BenchmarkEncodeUpdate`, baseline 9.4 µs, 48 allocs.
5. Result (2026-10-30, round 2)
   1. 200 games at once pass on `GOMAXPROCS=1`: median 84 ms, max 214 ms; `-ban` default in place.
   2. The 128 failures were the send queue: 18 of 128 games cut off at 64 queued messages.
   3. Send queue 64 to 256 (`internal/server/ws.go`): 128 and 200 pass, no cut-off.
   4. Writes are batched: queued messages go in one write; pongs take the same lock.
   5. Each close logs its cause: full queue, failed write, unexpected read end.
   6. Load tool: the last stage is `-max-games` itself; `-ban` bans the loop cards.
   7. Measured, no gain kept:
      1. `GOGC=400`: 1.30 ms per step against 1.40 ms.
      2. `encoding/json/v2` `MarshalerTo` for `Args`: 15.4 µs against 9.4 µs; other bytes; reverted.
      3. Typed `Args` (ints, strings, `CardValue`, `PlayerValue`): same bytes, no CPU gain; kept, tested.
      4. Batched writes: 97 µs to 91.5 µs per step, about 6%.
      5. fastjson: intent decode 387 ns to 91 ns; intents are a small share.
      6. easyjson (generated, `omitempty` only): update encode 36.7 µs to 27.9 µs; `omitzero` does not generate.
   8. On loopback 82% of CPU is in write syscalls; the load client shares the Mac, and tracing adds cost.
   9. Open:
      1. 19.97 and 19.109: endless loops; the load tool bans their cards.
      2. Ping max 6 s at 200 games; not yet explained.
      3. About 290 MB heap after 200 finished games; rooms stay until the empty limit.
6. Plan (2026-10-30, owner answers)
   1. LAN baseline, step 0: the server here, the load from a second laptop.
      1. Server: `192.168.1.137:4810`, `GOMAXPROCS=1`, tagged build.
      2. Load: `-max-games 200 -seed 11`; the default `-ban` applies.
      3. Result: the stage table, the slowest response, and the profile.
   2. Record batching: one flush per second and on every shutdown.
      1. Graceful shutdown writes every queued line before the files close.
      2. A game end still writes its end line at once.
      3. Wiki `Profiling.md` and ADR for the records name the rule: a crash loses up to a second.
      4. Benchmark: CPU per step against step 0 on the same seed.
   3. Encoder, one piece at a time, each with a benchmark and a row here.
      1. Benchmark: a real game's updates, not a synthetic one.
      2. Pieces, in order: `Args` (done, typed), events, view, seats, allowed intents.
      3. Wire bytes stay the same; a golden test per piece.
   4. Not now (owner): coalescing queued updates; reconsider later.
   5. Engine keeps the loop cards (19.97, 19.109); the load tool bans them.
7. Log of steps (2026-10-30)
   1. Step 0, LAN baseline (committed state: send queue 256, batched writes, typed `Args`).
      1. Server: this laptop, `GOMAXPROCS=1`, tagged, 192.168.1.137:4810; load: second laptop, seed 11, bans default.
      2. 200 games pass: response median 81.8 ms, p95 117 ms, max 835 ms; no game failed; no notices.
      3. Ping median 81 ms, p95 130 ms, max 8.6 s (not explained).
      4. Handshake: median 313 ms, p95 1.37 s, max 1.41 s (800 connections).
      5. Profile, `cpu.pprof` over 648 s, 347 s of CPU: 62% writes (syscalls); 18% engine `View` (`eachPart`, `facesOf`); 7% JSON.
   2. Step 1, record batching: flush once a second and at every stop.
      1. Code: `pkg/record/record.go`; tests `pkg/record/flush_test.go` (stop writes every line, flush is readable, tick flushes).
      2. ADR 007 item 3.1 updated in place (flush cadence; the format stays).
      3. Local A/B, same seed 11, 64 games, two runs each: step 0 97 and 102 µs per step; step 1 88 and 93 µs per step (about 9% less).
      4. Record size, 126 files: 69 MB at step 0, 30 MB at step 1 (fewer flush blocks).
      5. Every record closes cleanly after SIGTERM (`gzip -t` passes on all 126 files).
      6. Gates: `make lint` exit 0; `go test` for `pkg/record` and `internal/server` (record subset) pass.
      7. Not yet: LAN run for step 1.
   3. Step 2, encoder: baseline and first piece.
      1. Baseline, real updates (`BenchmarkEncodeRealUpdate`, `internal/server/bench_test.go`): 10.1 µs per update, 16 allocs.
      2. Bytes per update: about 3.2 KB; the view is 88%, events 4%, allowed under 1%.
      3. Profile: reflection in `encoding/json/v2` struct encoding is most of the encode time.
      4. Piece 1, card view: `internal/protocol/encode_card.go`; test `encode_card_test.go`.
         1. Hand-written: 168 ns against 334 ns per card (`BenchmarkAppendCard`, `BenchmarkMarshalCard`).
         2. Equal bytes to `encoding/json` on every optional field and on escaped strings.
         3. Not wired yet: it gains only when its parents (player, slot, view) are written by hand too.
      5. Next: the player view, then the slots, the stack, the turn, the prompt; each with a row here.
   3. Step 1 LAN run (2026-10-30, second laptop, 200 games, seed 11, bans default).
      1. 200 games pass. Response median 75.2 ms (step 0: 81.8), p95 113 ms (117), max 491 ms (835).
      2. Ping median 77 ms (81); ping max 10.4 s (8.6 s): still not explained.
      3. Handshake median 309 ms (313); max 1.53 s (1.41 s).
      4. Not a clean A/B on the network: one 128-game stage hit 935 ms; the loopback A/B above is the CPU figure.
   4. Step 2 encoder, pieces and benchmark (`BenchmarkEncodeRealUpdate`, 3 s, same updates).
      1. Baseline: 10.1 µs, 9.5 KB allocated, 16 allocs.
      2. Piece 1, card view: 168 ns against 334 ns (`BenchmarkAppendCard`).
      3. Piece 2, player and slot: equal bytes; `encode_view.go`, `encode_view_test.go`.
      4. Piece 3, update envelope and view wired into `Encode`: 7.9 µs, 18.5 KB, 66 allocs (buffer not sized yet).
      5. Piece 4, preallocated 4 KiB buffer: 6.8 µs, 7.9 KB, 55 allocs.
      6. Piece 5, turn, deck and discard arrays, prompt: 5.4 µs, 6.8 KB, 41 allocs.
      7. Tests: `TestEncodeUpdateMatchesJSON` (real game updates), `TestEncodeFullUpdateMatchesJSON` (every optional field).
      8. Still by `encoding/json`: events, allowed, seats, stack, the optional view parts, rollback fields.
   5. Step 2, pieces after the card: `Turn`, arrays, prompt, events, intents, seats, empty stack, envelope in one buffer.
      1. Real updates: 10.1 µs → 4.5 µs per update; 16 → 29 allocs; 9.5 KB → 7.1 KB allocated.
      2. Gates: `make lint` exit 0; `internal/protocol`, `internal/server` tests pass; real-update equality test passes.
      3. Loopback A/B, 64 games, seed 11, two runs each: step 1 89.8 and 87.1 µs per step; step 2 83.3 and 93.7 µs per step (no clear change; noise about ±5 µs).
      4. Profile: `protocol.Encode` 1.41 s → 0.60 s of server CPU; syscalls 84% of CPU in both runs.
      5. Engine `View` is 5% to 7%; next candidate, rules engine, owner's call.
      6. Still by `encoding/json`: non-empty stacks, rollbacks, pauses, aims, blocked, the rare view parts.
   6. Next candidates, owner decides:
      1. Write window: hold each client's frames up to 1 ms, one write for a burst of steps (no data dropped; a little latency).
      2. The engine's view: build the shared part once per step instead of once per seat.
   7. Step 3a, write window: rejected (owner asked, 2026-10-30).
      1. A 1 ms timer per write: write syscalls 68 → 45 µs per step; `kevent` time 2 s → 8 s; total CPU no better; adds up to 1 ms.
      2. `runtime.Gosched` before the batch: write syscalls 15–23 µs per step; reads rise to 55 s; total CPU no better.
      3. Both reverted; the code is the step 2 writer (batch of the queued messages).
   8. Step 3b, shared view (owner asked, 2026-10-30).
      1. `ViewBase` (once per step) and `ViewFor` (per viewer) in `pkg/rules_engine/view.go`; `View` = both.
      2. Room: `broadcast` builds the base once; each client's update adds its viewer part.
      3. Viewer-dependent parts: hands, treasure peek, the prompt's choices, ordering.
      4. Tests: `pkg/rules_engine`, `internal/protocol`, `internal/server` pass; `make lint` exit 0.
      5. Loopback A/B, 64 games, seed 11, two runs each: step 2 80.7 and 82.5 µs per step; step 3b 55.9 and 60.5 µs per step (about −29%).
      6. Engine `View` cum: 4.0–4.5 s → 0.9–1.2 s of CPU.
      7. 64-game response median: 21 ms → 15 ms (both runs).
      8. No send-queue cut-off; two broken-pipe lines at game end (client closed), as before.
      9. Not yet: LAN run for step 3b.
   9. Step 3b LAN run (2026-10-30, second laptop, 200 games, seed 11, bans default).
      1. 200 games pass; no failure, no reconnect, no notice.
      2. Response median 45.7 ms (step 0: 81.8; step 1: 75.2); p95 71 ms (117; 113); max 318 ms (835; 491).
      3. Ping median 46 ms (81; 77); ping max 7.0 s (8.6 s; 10.4 s): still not explained.
      4. Handshake median 118 ms (313; 309); p95 1.22 s (1.37 s; 1.43 s).
      5. Stage 64: response median 14.6 ms; stage 128: 27.7 ms.
