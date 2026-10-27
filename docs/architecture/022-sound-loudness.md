# 22. Sound loudness

* **Status:** Proposed
* **Date:** 2026-10-27
* **Authors:** @HardDie

---

## Context

1. The game has four sounds from different sources (AU-03).
   1. Hearthstone hit, cast and impact (15.50); a dice roll (15.69).
   2. Their styles cannot match, so their loudness must (AU-04).
2. Measured on 2026-10-27, loudest window in LUFS:

| Sound | 100 ms | 200 ms | 400 ms | Sample peak |
|---|---|---|---|---|
| attack-hit | -9.9 | -12.4 | -15.5 | +0.1 dBFS |
| dice-roll | -26.8 | -27.8 | -30.8 | -8.0 dBFS |
| effect-cast | -15.3 | -15.8 | -17.1 | -5.6 dBFS |
| effect-impact | -8.5 | -10.5 | -11.3 | +1.2 dBFS |

3. The spread is about 17 LU; a listener notices 2–3 LU.
   1. The dice is far too quiet.
   2. effect-impact is the loudest and clips.
4. Peak level does not show how loud a sound feels.
   1. BS.1770 (EBU R128) weights frequencies like the ear.
   2. Its integrated loudness needs at least 400 ms of sound.
   3. The dice lasts 335 ms and reads -70 LUFS: useless.
5. For short sounds the loudest window is the right number.
   1. The window length changes the order of short hits.
   2. The owner compared 200 ms and 400 ms levelled sets by ear.
   3. 200 ms sounded even (owner, 2026-10-27).

## Considered options

1. **Loudest 200 ms window, checked by a tool in `FourSoulsRepo/tools`**
   1. Won. Matches the owner's ear; keeps audio code out of the game.
2. **Loudest 400 ms window (EBU momentary)**
   1. Lost. Under-rates short hits; the levelled set sounded uneven.
3. **Integrated loudness**
   1. Lost. Fails on sounds shorter than 400 ms.
4. **Check as a Go test in this repo**
   1. Lost. Needs an MP3 decoder dependency here (owner: tools repo only).
5. **Per-sound volume in `sounds.ts`**
   1. Lost. Hides the problem; the files stay uneven.

## Decision

Use option 1.

1. Measure: "M200".
   1. BS.1770 K-weighting at the file's sample rate.
   2. Mean square over a sliding 200 ms window.
   3. The sound is padded with silence, so short sounds fill the window.
   4. M200 is the loudest window, in LUFS.
2. Target: every sound at -20 LUFS M200, ±1.5 LU.
3. Sample peak at most -1 dBFS after decoding.
4. The tool: `go run ./loudness` in `FourSoulsRepo/tools`.
   1. `check` measures every sound and fails on any out of range.
   2. `level` fixes a sound: gain, limiter, MP3 re-encode.
   3. MP3 decoding: `hajimehoshi/go-mp3`, only in the tools repo.
5. Files stay mono MP3 in `images/assets/sounds/` (15.50).
6. `sounds.ts` plays every sound at one volume.
7. A new or changed sound passes `loudness check` before it is committed.

## Consequences

### Positive

1. No sound jumps out or gets lost.
2. A new sound is checked with one command.
3. The game repo gains no audio code or dependency.

### Negative and risks

1. The check is not in `make test`; it relies on the rule.
   1. The sounds live in the private images, so CI could not run it anyway.
2. M200 is our own window, not a broadcast standard.
   1. At 400 ms the same code must match ffmpeg `ebur128`.
3. Loud transients need a limiter; a big boost can dull a sound.
4. Levelling to -20 LUFS makes the game quieter overall.
   1. `sounds.ts` volume goes from 0.6 to 1.0 to make up for it.

### Neutral

1. `level` needs ffmpeg to encode MP3; `check` needs only Go.
