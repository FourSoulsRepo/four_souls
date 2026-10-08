# Step 10. Replays

Goal: watch saved matches in the app.

Ideas: RP-02, RP-04, RP-10, RP-11, M-01.

---

### [ ] 10.1 Records screen

1. Goal: manage saved matches (M-01).
2. Tasks:
   1. List own saved records and hosted autosaves.
   2. Watch, export, delete.
3. Done when:
   1. Use case: manage records.

### [ ] 10.2 Replay player

1. Goal: play a record on the table (RP-11).
2. Tasks:
   1. Reuse the table and event playback from step 7.
   2. Play, pause, speed.
   3. Step one event forward or back.
   4. Jump to any turn.
3. Done when:
   1. A full Base Game record plays to the end.

### [ ] 10.3 Hands and unknown cards

1. Goal: RP-04, RP-10.
2. Tasks:
   1. Toggle: show or hide hands.
   2. Unknown cards use the card text stored in the record.
3. Done when:
   1. A record with a card the client lacks still plays.

### [ ] 10.4 Engine check

1. Goal: optional check against the current engine (RP-02).
2. Tasks:
   1. Re-run steps on the current engine.
   2. Mark steps that would not work now; never reject.
   3. Compare checksums to find the first difference.
   4. Viewer can turn the marks off.
3. Done when:
   1. A test record with a changed rule shows marks.
