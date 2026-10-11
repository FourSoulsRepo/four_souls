# Step 17. Rollback

Goal: a mistaken event can be undone by the other players.

Ideas: N-12.

---

### [ ] 17.1 Rollback ADR

1. Goal: clear limits for rollback (N-12).
2. Tasks:
   1. **Owner**: may hidden info (draws, peeks) be rolled back?
   2. **Owner**: may a dice roll be rolled back?
   3. **Owner**: in 3–4 player games, one other player or all?
   4. How far back: last event only, or more.
3. Done when:
   1. ADR accepted.

### [ ] 17.2 Engine undo

1. Goal: the engine returns to the state before an event.
2. Tasks:
   1. Keep the state before each event (engine save, step 4).
   2. Undo restores it exactly, RNG included.
   3. Fuzz check: undo then redo gives the same state.
3. Done when:
   1. Tests cover undo of a loot play, an attack, a purchase.

### [ ] 17.3 Server and client

1. Goal: players use rollback in a game.
2. Tasks:
   1. Intent "roll back"; refused for the event's author.
   2. History entry shows who rolled back what.
   3. The record keeps the rollback (RP-03).
   4. Button near the history entry, for other players only.
   5. Use case: roll back a mistake.
3. Done when:
   1. e2e: the author sees no button; another player rolls back.
