# Step 9. Spectators

Goal: viewers watch live matches.

Ideas: SP-01, SP-03, SP-05, CD-09.
Judges and viewer delay moved to Future in [../roadmap.md](../roadmap.md).

---

### [ ] 9.1 Roles in the protocol

1. Goal: viewer role.
2. Tasks:
   1. Join as viewer.
   2. Viewer view: public only.
   3. Players see only the viewer count.
3. Done when:
   1. View tests prove viewers never see hands.

### [ ] 9.2 Non-blocking broadcast

1. Goal: viewers never slow the game (SP-03).
2. Tasks:
   1. Each viewer has a bounded buffer.
   2. A slow viewer skips ahead to the latest state.
   3. New viewers get the current state first.
3. Done when:
   1. A stalled viewer does not delay players in tests.

### [-] 9.3 Viewer delay

1. Moved to Future: tournament tools.

### [ ] 9.4 Unknown cards

1. Goal: outdated clients can watch (SP-05, CD-09).
2. Tasks:
   1. Client asks the server for unknown card text.
   2. Shows blank cards with text and an "outdated" note.
3. Done when:
   1. An old client watches a game with a new card.

### [ ] 9.5 Spectator UI

1. Goal: watching in the app.
2. Tasks:
   1. Join as viewer from the join screen.
   2. Table without action controls.
   3. Judge log panel: moved to Future.
3. Done when:
   1. Use case: watch a match.
