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
