# 31. Replays without the engine

* **Status:** Accepted
* **Date:** 2026-11-16
* **Authors:** @HardDie

---

## Context

1. Owner, 2026-11-16: a replay player for match records (RP-11).
   1. Step forward, step back, autoplay, playback speed.
2. Owner: a record must play without the rules engine.
   1. A later rule fix must not stop an older record from playing.
3. Owner: only records made from this change on (format 2).
   1. No support for older record formats.
   2. No migrations between view schemas.
4. Owner: no engine check for now (RP-02.2–5).
5. Today a `step` line keeps intent, events and checksum (ADR 007).
   1. The table draws from the engine's `GameView`.
   2. Only the engine can rebuild that view from intents.
   3. `engine.CanReplay` refuses records of many older engines.
6. Measured, random 4-player games of 4000 steps, gzip:
   1. Intent, events, checksum: about 80 KB.
   2. Plus the full judge view every step: about 250 KB.
   3. Plus a view diff, full every 100 steps: about 130 KB.
   4. One judge view: about 2.7 KB of JSON.

## Considered options

1. **Full judge view in every step line**
   1. Won. Any step shows alone; back and jump are trivial.
   2. About three times today's size; still a few hundred KB.
2. **View diffs, a full view every 100 steps**
   1. Lost. Half the size, but diff code on both sides.
   2. Jumping back needs the nearest full view first.
3. **Replay through the engine, as ADR 007 item 5**
   1. Lost. Old engines' records stop playing on any rule fix.
4. **The engine's saved state per step**
   1. Lost. Larger, engine-internal, and still needs the engine.

## Decision

Use option 1.

1. Record format 2
   1. `record.Format` becomes 2.
   2. Each `step` line also holds `view`: the judge view after it.
   3. The judge view holds every hand, never deck order or RNG.
   4. The server writes the protocol encoder's JSON, as clients read it.
   5. Format 1 records cannot be watched.
      1. The records screen says they were made before replays.
   6. A record of a newer format than the app's says "update the app".
2. Playing a record
   1. The player reads the header, seats, cards and live steps.
   2. It never calls the rules engine.
   3. Rolled back steps are left out, as `Record.Live` does.
   4. Unknown cards show the header's card text (RP-10).
   5. Views and events play as they are; no migration.
3. Code
   1. `internal/replay`: reads a record, serves its live steps.
   2. `bindings/records` gains: `List`, `Open`, `Step`, `Delete`, `Export`, `PickFile`.
   3. The binding keeps the open record; the client asks for steps by number.
   4. The client feeds each step to playback as an update: events, view, seats.
   5. Replay updates allow no intents and carry no aims (N-15).
4. The player
   1. Play and pause: autoplay goes on step by step.
   2. Step forward plays the step's events as effects, as live.
   3. Step back and jump show the step's view at once, with no effects.
   4. Jump to any turn, and a slider over all steps.
   5. Speed: 0.5×, 1×, 2×, 4×; it scales effects and pauses.
   6. Autoplay waits for the effects, then a short pause.
      1. 0.8 s at 1× after a step with events.
      2. 0.2 s at 1× after a step without.
   7. Show or hide hands (RP-04); hands show by default.
   8. The bar uses the game's own look: buttons, slider, turn list.
   9. The viewer can hide the bar while watching (owner, 2026-11-16).
      1. A small button stays and brings it back.
      2. Hiding the bar never pauses the replay.

## Consequences

### Positive

1. A record plays on a later app, without the engine.
2. Rule fixes never break a format 2 record.
3. Back and jump cost one step lookup.

### Negative and risks

1. Records grow about three times.
2. A protocol bump that renames or removes a view or event field.
   1. Older records may then show parts of the table wrong.
   2. Accepted: no migrations (owner, 2026-11-16).
3. The server builds one more view per step.

### Neutral

1. Engine replay stays for tests (`engine.Replay`).
2. The checksum chain (ADR 025) stays in every step.
