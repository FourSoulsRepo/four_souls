# 7. Match record format

* **Status:** Accepted
* **Date:** 2026-10-10
* **Authors:** @HardDie

---

## Context

1. Every match is recorded, finished or not (RP-01, RP-07, RP-08).
2. Records replay exactly from a seed and the applied intents (A-08).
3. A record holds all hands; players get it only after the end (RP-04, RP-05).
4. Outdated clients still show the cards of a record (RP-10).
5. Records stay small and readable for debugging (RP-09).
6. `pkg/record` is its own module; it should not need the engine (A-12).

## Considered options

1. **gzip-compressed JSON lines**
   1. Won. Stdlib only; one line per step; `zcat` shows it.
2. **zstd**
   1. Lost. Smaller, but a new dependency for files of a few hundred KB.
3. **The engine's Save of the whole state per step**
   1. Lost. Far larger; the seed and intents already rebuild the state.
4. **Intents and events as engine types in pkg/record**
   1. Lost. The module would follow the engine version; old files might not load.
5. **Intents and events as raw JSON**
   1. Won. The record module reads every version; the replayer decodes them.

## Decision

Use options 1 and 5.

1. File
   1. Extension `.fsrec`; gzip stream of JSON lines.
   2. Name: start time and game ID, e.g. `2026-10-10-120000-g3.fsrec`.
2. Lines
   1. First a header: format, app, engine and protocol versions, game ID, start time.
   2. The header also holds the setup: seed, players, set codes, bonus souls, characters.
   3. Then seats (nicknames) and the text of every card of the sets (RP-10).
   4. One `step` line per applied intent: intent, every event (hands too), checksum (RP-12).
   5. Rejected intents are never recorded (RP-03).
   6. Last an `end` line: finished or not, winners.
   7. A saved game loaded again adds a gzip part with a `resume` line (N-11).
3. Writing
   1. Each line is flushed: a crash leaves a readable record without `end`.
   2. A server that stops writes `end` with finished false (RP-08).
4. Server
   1. Records go to the configured folder; a hosted game uses the user's config folder.
   2. Retention: 30 days on a dedicated server, forever when hosted; 0 keeps them (RP-07).
   3. Old records are deleted at start and every hour.
   4. `GET /record?game=…&token=…` gives a finished game's record to its players only (RP-05).
5. Replay
   1. Rebuild the engine setup from the header and the set codes.
   2. Apply the steps' intents; every checksum must match (`engine.Replay`).

## Consequences

### Positive

1. Any record replays exactly, or shows the first step that differs.
2. Records are readable with standard tools.

### Negative and risks

1. Card text in every header repeats data; a few dozen KB compressed.
2. A record cut off mid-line loses that last line.

### Neutral

1. Replays and the record list in the app come with step 10.
