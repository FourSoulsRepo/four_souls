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
