# Step 10. Replays

Goal: watch saved matches in the app.

Ideas: RP-02, RP-04, RP-10, RP-11, M-01.

Rules for this step (owner, 2026-11-16):

1. Work order: 10.5, 10.6, 10.1, 10.2, 10.3.
2. Decisions: ADR 031; read it before any sub-step.
3. The step runs on the `step-18-play-test` branch.
4. Playing a record never calls the rules engine.
5. Only format 2 records play; no migrations, no old formats.

---

### [ ] 10.1 Records screen

1. Goal: manage saved matches (M-01).
2. Tasks:
   1. List own saved records and hosted autosaves.
   2. Watch, export, delete.
   3. "Open a file" for a record anywhere else.
   4. A format 1 record says "made before replays" (ADR 031, 1.5).
   5. A record newer than the app says "update the app" (ADR 031, 1.6).
   6. Binding: `List`, `Open`, `Delete`, `Export` in `bindings/records`.
3. Done when:
   1. Use case: manage records.

### [ ] 10.2 Replay player

1. Goal: play a record on the table (RP-11).
2. Tasks:
   1. Reuse the table and event playback from step 7.
   2. Play, pause, speed.
   3. Step one event forward or back.
   4. Jump to any turn.
   5. Steps come from the record's views; never the engine (ADR 031, 2).
   6. Forward plays effects; back and jump show the view at once (ADR 031, 4).
   7. Speeds 0.5×, 1×, 2×, 4×; autoplay pauses as ADR 031, 4.6.
   8. A slider over all steps.
   9. No action controls, no intents, no aims (N-15).
   10. Binding: `Step` in `bindings/records`.
3. Done when:
   1. A full Base Game record plays to the end.
   2. Use case: watch a record.

### [ ] 10.3 Hands and unknown cards

1. Goal: RP-04, RP-10.
2. Tasks:
   1. Toggle: show or hide hands.
   2. Unknown cards use the card text stored in the record.
   3. Hands show by default (ADR 031, 4.7).
3. Done when:
   1. A record with a card the client lacks still plays.

### [-] 10.4 Engine check

1. Dropped for now (owner, 2026-11-16).

### [x] 10.5 Replay ADR

1. Goal: decide how records play without the engine (owner, 2026-11-16).
2. Tasks:
   1. Measure record size with views per step.
   2. ADR 031; ADR 007 items 2.20, 5.8, 5.9.
3. Done when:
   1. ADR 031 is Accepted.

### [ ] 10.6 Views in the record

1. Goal: every new record plays without the engine (ADR 031, 1).
2. Tasks:
   1. `record.Format` 2; `Step` lines get `view`.
   2. The server writes the judge view after each step.
   3. It uses the protocol's view encoder, not a second `json.Marshal`.
   4. `internal/replay`: read a record, serve its live steps.
   5. It refuses format 1 and newer formats with a clear error.
3. Done when:
   1. A server game's record has a view per step.
   2. That view equals the engine's judge view after the step.
