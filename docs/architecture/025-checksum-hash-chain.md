# 25. Checksum: hash chain with anchors

* **Status:** Accepted
* **Date:** 2026-10-29
* **Authors:** @HardDie

---

## Context

1. Every applied step stores a checksum (RP-12).
2. Replay checks each one against the rebuilt state.
3. The checksum hashed the whole game as JSON, every step.
4. That was about a third of server CPU at 64 games.
5. Replay must still catch any state that differs from the record.

## Considered options

1. **Keep the full-state checksum**
   1. Lost. Too costly per step.
2. **Hash chain on every step only**
   1. Lost. A hidden state change with no event goes unseen.
3. **Hash chain on every step, full state every 100 steps**
   1. Won. Cheap per step. A divergence is seen within 100 steps.
4. **Incremental hash of the state per part**
   1. Lost. Needs change tracking in every engine path.

## Decision

Use option 3.

1. Each step: `chain = hash(chain, intent, events, RNG)`.
   1. The chain and the step count are game state.
2. Every 100th applied step: the checksum also covers the whole state.
   1. `checksum = mix(chain, hash(state))`.
3. Records from engines before 0.5.0 cannot replay.
   1. `CanReplay` says so, as for earlier changes.
4. Records keep one checksum field, as before.

## Consequences

### Positive

1. Per-step checksum cost falls to the size of one step.
2. Full-state checks still run, once per 100 steps.

### Negative and risks

1. A divergence in the state is reported up to 100 steps late.
2. Records from before 0.5.0 no longer replay.
3. A new state field must still be saved with `omitempty` (N-20).

### Neutral

1. The checksum values change, so tests that pin them were updated.
