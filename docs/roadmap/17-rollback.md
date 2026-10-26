# Step 17. Rollback

Goal: a mistaken event can be undone by the other players.

Ideas: N-12.

Rules for this step (owner, 2026-10-26):

1. Sub-steps run one after another, each in a fresh subagent.
2. Stop only for a question the owner must answer.
3. The owner reviews the whole step at its end.
4. UI sub-steps leave screenshots for that review.

Owner notes (2026-10-26):

1. Owner reviewed step 17 and marks it solved.

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

### [x] 17.5 Rollback over a server restart

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
   3. Server test: an old restart file loads with no history.
   4. e2e: Godmode rolls back past a restart; a cut shows its row.

### [x] 17.6 Soft: roll back a played card

1. Goal: in Soft, a played card can be undone until a deck card shows (N-12, item 24).
2. Tasks:
   1. Played card: a loot card played, or an item's ability activated.
   2. Not a purchase; a purchase keeps today's barriers.
   3. Undoable while no new card came from the treasure, monster or loot deck.
   4. And no dice was rolled or changed, and no damage was dealt.
   5. Still blocking: a turn end, a death, a seen hand card (owner, 2026-10-26).
   6. Already true: `loot_played` and `activated` are no barriers.
   7. Change: a private `moved_to_deck` is no barrier; no deck card shows.
   8. ADR 020 item 4: append the rule in plain words; UC-31 likewise.
   9. A play gets its own history row, "Tom plays A Penny!" (owner, 2026-10-26).
   10. That row carries the ↶ button, also while the card is on the stack.
   11. It changes the 15.88 history rule; update its docs.
   12. Fixes 19.69: Godmode ↶ on the result row left the card on the stack.
3. Done when:
   1. Server tests: a played loot card and an activated item are undone.
   2. Server tests: each blocking event after the play stops it.
   3. Engine test: a private `moved_to_deck` is no barrier.
   4. e2e: a played loot card is undone by a Soft vote.

### [x] 17.7 One Undo button

1. Goal: one Undo button replaces the buttons on history rows (N-12, item 25).
2. Tasks:
   1. Undo sits under the History title; hidden in Disabled mode.
   2. One press: back to before the latest target in `rollback_steps`.
   3. A play and the answers it asks first are one move.
   4. Soft: the author sees "Ask to undo"; others "Undo", which opens the vote.
   5. Godmode: anyone, at once; many presses go far back.
   6. Greyed with the reason: barrier, open vote, game over, spectator.
   7. The tooltip names the move it undoes.
   8. Remove the ↶ buttons on rows and the step-to-row mapping (17.3, 17.6).
   9. Remove the 17.6 play rows; the history stays clean (owner, 2026-10-26).
   10. Keep the 17.6 engine rule (`moved_to_deck`) and its server tests.
   11. 15.88 item 7 and UC-15 2e: the play rows are gone again.
   12. ADR 020: append the change; item 6.2 "any step" is now press by press.
   13. UC-30, UC-31 updated; screenshots redone.
   14. 19.69: its row fix is replaced by the button; note it there.
3. Done when:
   1. Vitest: the button's target, label, grey reasons, hidden when Disabled.
   2. e2e: Godmode presses twice and goes back two moves.
   3. e2e: a Soft press opens the vote; the author's press asks.
   4. e2e: a played card that asks a target is undone by one press.
   5. e2e: the button is greyed after a dice roll in Soft.
