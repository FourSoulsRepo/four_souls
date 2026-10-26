# Step 17. Rollback

Goal: a mistaken event can be undone by the other players.

Ideas: N-12.

Rules for this step (owner, 2026-10-26):

1. Sub-steps run one after another, each in a fresh subagent.
2. Stop only for a question the owner must answer.
3. The owner reviews the whole step at its end.
4. UI sub-steps leave screenshots for that review.

---

### [x] 17.1 Rollback ADR

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
   14. Server restart: history kept if the replay checks out (N-12, item 23).
   15. Godmode: no vote; any history row is a target.
   16. The default mode is Disabled (owner, 2026-10-26).
   17. Godmode: anyone may roll back, the author too (owner, 2026-10-26).
   18. The RNG is read only when the game needs a value (N-12, item 20).
   19. Old states come from a replay from one base state (N-12, item 22).
   20. ADR 019 updated: the restart file keeps the old base state.
3. Done when:
   1. ADR accepted: [020](../architecture/020-rollback.md).

### [x] 17.2 Engine undo

1. Goal: the engine returns to the state before an event.
2. Tasks:
   1. Rebuild the state before a step: base state, then replay (N-12, item 22).
   2. Undo restores it exactly, RNG included.
   3. Fuzz check: undo then redo gives the same state.
   4. Test: roll, roll back, roll again gives the same value.
3. Done when:
   1. Tests cover undo of a loot play, an attack, a purchase.

### [x] 17.3 Server and client

1. Goal: players use rollback in a game.
2. Tasks:
   1. Message `rollback` naming a target step (Godmode).
   2. History entry shows who rolled back what.
   3. The record keeps the rollback (RP-03).
   4. Button near the history entry; any player in Godmode.
   5. Use case: roll back a mistake.
   6. Godmode: any player rolls back to any history row at once.
   7. ADR 007: the `rollback` line (ADR 020, item 9).
   8. Protocol version 5 (ADR 020, item 8).
   9. Scope: Godmode only; Soft moved to 17.4 (main agent, 2026-10-26).
   10. `Options.Rollback` in the record header, saves and restart files.
   11. Set by the create options; the setup control is 17.4.
   12. The game list and the table show the mode (ADR 020, item 1.4).
   13. The server refuses `soft` until 17.4.
   14. A load or a restart is the new base until 17.5.
3. Done when:
   1. e2e: in Godmode another player rolls back.
   2. e2e: in Godmode the author rolls back their own step.
   3. Use case: [UC-30](../use-cases/app/uc-30-roll-back-a-mistake.md).

### [x] 17.4 Modes and the Soft vote

1. Goal: the host picks the mode; Soft rollbacks need a majority.
2. Tasks:
   1. Detailed setup: Disabled, Soft, Godmode; Disabled by default.
   2. The mode goes into the record header and saved games.
   3. Soft: the author may ask; others vote like the kick vote.
   4. A majority of the other connected players decides.
   5. Soft barriers: dice roll, damage, hidden info seen, turn end.
   6. Use case for each mode.
   7. Moved from 17.3: in Soft the author's rollback is refused; they ask.
   8. Moved from 17.3: in Soft the button shows for the other players only.
   9. Moved from 17.3: `engine.Barrier` and the server's barrier check.
   10. Done in 17.3: the mode in the record header, saves and restart files.
   11. Accept `soft` in `checkOptions` (17.3 refuses it).
3. Done when:
   1. Server tests: majority with 2, 3, 4 players; each barrier.
   2. e2e: a Soft vote passes and one is declined.
   3. Moved from 17.3: e2e: the author sees no button; another player rolls back.

### [ ] 17.5 Rollback over a server restart

1. Goal: a kept game keeps its rollback history (N-12, item 23).
2. Tasks:
   1. The restart file also keeps the old base state.
   2. On load, replay from the old base on the new engine.
   3. Same bytes as the kept state: keep the history.
   4. Otherwise the restart is a barrier; the history shows a row.
   5. ADR 019 updated.
3. Done when:
   1. Server test: same engine, rollback reaches past the restart.
   2. Server test: a changed replay makes the restart a barrier.
