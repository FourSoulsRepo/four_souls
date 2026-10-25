# Step 17. Rollback

Goal: a mistaken event can be undone by the other players.

Ideas: N-12.

---

### [ ] 17.1 Rollback ADR

1. Goal: clear limits for rollback (N-12).
2. Tasks:
   1. Hidden info: decided by the mode (owner, 2026-10-26).
   2. Dice may be rolled back in Godmode; the next roll keeps its value.
   3. 3–4 player games: a majority of the other players (owner, 2026-10-26).
   4. How far back: decided by the mode.
   5. Decided: the restored RNG repeats the roll (N-12).
   6. Never reseed after a rollback; that would be a free re-roll.
   7. Modes: Disabled, Soft, Godmode (N-12, item 15).
   8. Soft: list every barrier event; dice, damage, hidden info so far.
   9. The host picks the mode in Detailed setup (owner, 2026-10-26).
   10. Soft barrier too: the turn ends; death and reveals are covered.
   11. Soft: the author may ask; the others vote, like the kick vote.
   12. Which steps can be targets: never a pass or a throw alone.
   13. Record format: a rollback line; ADR 007 updated (N-12, item 21).
   14. Server restart (ADR 019): saved states kept or dropped.
   15. Godmode: no vote; any history row is a target.
   16. The default mode is Disabled (owner, 2026-10-26).
   17. Godmode: anyone may roll back, the author too (owner, 2026-10-26).
   18. The RNG is read only when the game needs a value (N-12, item 20).
3. Done when:
   1. ADR accepted.

### [ ] 17.2 Engine undo

1. Goal: the engine returns to the state before an event.
2. Tasks:
   1. Keep the state before each event (engine save, step 4).
   2. Undo restores it exactly, RNG included.
   3. Fuzz check: undo then redo gives the same state.
   4. Test: roll, roll back, roll again gives the same value.
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
