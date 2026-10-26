# 20. Rollback

* **Status:** Accepted
* **Date:** 2026-10-26
* **Authors:** @HardDie

---

## Context

1. A misclick cannot be undone; the game goes on with the mistake (N-12).
2. Groups differ: some want strict play, some a forgiving table.
3. A rollback must not become a free re-roll or a peek at hidden cards.
4. The engine is deterministic: a seed and the intents rebuild a game (ADR 007).
5. `Game.Save` and `Load` cover the whole state, the RNG included.
6. Every applied intent is one step with its events and a checksum.
7. A dedicated server keeps games over a restart (ADR 019).
8. A new engine may replay old steps differently (ADR 019, option 4).

## Considered options

1. **Three modes picked by the host**
   1. Won. Disabled, Soft and Godmode fit strict, casual and trusting tables.
2. **One rule for every game**
   1. Lost. Too strict for friends, too loose for strangers.
3. **Store the state before every step**
   1. Lost. Godmode reaches any step: hundreds of states per game.
4. **Rebuild an old state by replay from one base state**
   1. Won. One state per game; the engine replays in milliseconds.
5. **Advance the RNG on every step**
   1. Lost. The step count would pick the next roll.
   2. Same steps bring a good roll back; one more pass re-rolls a bad one.
6. **Drop undone steps from the record**
   1. Lost. Bug evidence would be gone; the record is append-only.
7. **A restart is a barrier in every mode**
   1. Lost. Godmode would lose its whole history on an update.
8. **A restart keeps the history if a replay checks out**
   1. Won. Same build or a small update keeps it; else it is a barrier.
9. **A rollback button on each history row (17.3 – 17.6)**
   1. Lost in 17.7. A step without a row had no button (19.69).
   2. Rows for plays cluttered the history; the owner wants one button.
10. **One Undo button for the latest move**
    1. Won in 17.7 (N-12, item 25). Press by press, with one tooltip.
11. **The app merges a play with its answers**
    1. Lost. The app knows no intents and no prompts of past steps.
12. **The server merges a play with its answers**
    1. Won. It keeps the intents and sees the prompt after each step.

## Decision

Use options 1, 4, 8, 10 and 12.

1. Modes
   1. `Options.Rollback`: `""` (Disabled), `soft`, `god`.
   2. The host picks it in Detailed setup; Disabled is the default.
   3. The record header, saves and restart files keep it.
   4. The game list and the table show the mode.
2. Disabled
   1. No rollback; no buttons.
3. Soft
   1. A target is a step after the last barrier.
   2. A target is never a lone pass or throw.
   3. The rollback goes to the state before the target.
   4. It undoes the target and every step after it.
4. Soft barriers
   1. A step with a barrier event cannot be undone, nor any step before it.
   2. Dice: `dice_rolled`, `roll_changed`, `random_pick`.
   3. Damage: `damaged`; a death follows damage.
   4. Hidden info: an event with `private`, `looked_at`, `card_revealed`.
   5. The turn ends: `turn_ended`, `turn_ended_early`.
   6. The game ends: `game_won`.
   7. Not a barrier: a shuffle; the restored RNG repeats its order.
   8. The engine owns the list: `engine.Barrier(Event) bool`.
   9. 17.4 checks the list against every event kind; a test pins it.
   10. The check (17.4) adds kinds that fit 4.2 to 4.6.
   11. Dice: also the roll's follow-ups and the first-player roll.
   12. Death: also `died` and `object_died`, for a death without damage.
   13. Hidden info: also `gained_treasure`, a deck top coming into view.
   14. Hidden info: also `character_dealt`; the turn: also `turn_skipped`.
   15. The game: also `game_started`.
   16. 17.6: a played card may be undone (N-12, item 24).
   17. Played card: a loot card played, or an item's ability activated.
   18. Undoable while no new card came from the treasure, monster or loot deck.
   19. And no die was rolled or changed, and no damage was dealt.
   20. A turn end, a death, and a seen hand card still block it.
   21. A purchase keeps its barriers.
   22. So a private `moved_to_deck` is no barrier: a hand card goes on the loot deck.
   23. Every other private event draws a deck card or shows a hand card.
5. Soft vote
   1. Only the other players can roll back; never the target's author.
   2. The author may ask; the ask opens a vote.
   3. Another player may start the vote; their own vote counts as yes.
   4. Voters: the connected seats other than the author.
   5. More than half of them must say yes.
   6. 2 players: the opponent; 3: both others; 4: two of three.
   7. No voter connected: no rollback.
   8. The game waits while a vote is open.
   9. No timer by default; the vote stays open until it is decided.
   10. With the response timer on, the vote gets the same seconds.
   11. A vote that runs out is declined.
   12. A vote ends early once the majority is reached or out of reach.
   13. Whoever opened the vote may cancel it.
   14. A declined target may be asked again; players may change their mind.
   15. Built like the kick vote (N-08).
   16. One vote at a time; a voter may change their answer.
   17. A player who leaves pauses the game; the pause ends the vote.
   18. A vote that ends with no rollback leaves a history row.
   19. A passed vote's `rollback` line names the opener.
6. Godmode
   1. Any player rolls back, the author too; no vote.
   2. A target is any step since the base (item 7), reached press by press (item 12).
   3. Barriers do not apply; only the game end does.
   4. Dice: the restored RNG repeats the roll (N-12, item 7).
7. Base state and replay
   1. The server keeps one base state per game: the start or the last restart.
   2. It also keeps the live steps since the base.
   3. Rollback to before step k: load the base, apply the live steps before k.
   4. Each checksum must match its step; else the rollback fails, the game goes on.
   5. The steps from k on leave the live list.
   6. Engine: `engine.ReplayFrom(state, steps, sets...)`.
   7. A loaded save is a new base too.
8. Steps and history
   1. Step numbers never repeat; new steps go on after the highest.
   2. A target is named by its step number.
   3. `update` gains `rollback`: who rolled back to before which step.
   4. The client greys the undone history rows and adds a rollback row.
   5. The table jumps to the restored view; no animation (ADR 012).
   6. Protocol version 5: new messages and a field older clients misread.
   7. A rollback takes a step number itself.
   8. So `update` and record step numbers stay the same.
   9. Every `update` also carries `rollback_mode` and `rollback_from`.
9. Record
   1. Undone steps stay in the record.
   2. A `rollback` line: the target step, who, and how many steps it undid.
   3. Replay: on a `rollback` line, rebuild from the start or last `resume`.
   4. ADR 007 gains the line in 17.3.
   5. 17.5: a rollback may reach behind a restart's `resume` (item 10).
   6. `Record.Live` still replays from the start.
10. Server restart (ADR 019)
   1. The restart file also keeps the base state and the live steps.
   2. On load, the new engine resumes the base (`engine.Resume`).
   3. It applies the live steps and saves the result.
   4. Same bytes as the kept state: the history stays; the base stays.
   5. Else, or if the base does not resume: the kept state is the new base.
   6. The history then shows a "server restarted" row as a barrier.
   7. An app's own host keeps no games, so this is the dedicated server only.
   8. ADR 019 gains this in 17.5.
   9. 17.5: the old per-step checksums are not compared on load.
   10. A new engine may sum otherwise; the final bytes decide.
   11. Kept live steps take the new engine's checksums.
   12. A kept step stays a Soft barrier if it was one, or is one now.
   13. A step the new engine cannot read or refuses also ends the history.
   14. A restart file without the history (before 17.5) ends it too.
   15. A cut history: every `update` carries `rollback_cut`, the restart's step.
   16. The app adds one "server restarted" row for each cut (refines 10.6).
   17. The app keeps its history over a reconnect, so the row sits at the cut.
   18. A restarted app has no old rows; the row is its first.
   19. An open Soft vote is not kept; it ends with the restart.
11. Never reseed after a rollback; that would be a free re-roll.
12. One Undo button (17.7, N-12, item 25)
    1. The History has one Undo button under its title; none in Disabled.
    2. One press undoes the latest target: the last of `rollback_steps`.
    3. Press again to go further back; there is no jump to an older row.
    4. Every update lists the targets, in Godmode too.
    5. A play and the answers it asks first are one move.
    6. Answers: a mode, a target, or the ability a loot card copies.
    7. Such an answer is no target; the play's step is.
    8. The server marks them: it sees the prompt after each live step.
    9. An attack's target is a move of its own: a priority window comes first.
    10. Soft refuses an answer as a target; Godmode still takes any live step.
    11. Soft: the author's button asks; another player's opens the vote.
    12. Greyed, with the reason: a barrier, a vote, a pause, the game end, a spectator.
    13. The app names the barrier from its events: "Nothing to undo since a die was rolled".
    14. The tooltip names the move: "Undo: Bo plays A Penny!".
    15. History rows have no buttons; the 17.6 play rows are gone.

## Consequences

### Positive

1. A misclick can be undone where the table allows it.
2. No state copies per step; memory stays one base per game.
3. Every undone step stays in the record as evidence.
4. Disabled by default: nothing changes for a game that does not opt in.

### Negative and risks

1. After a rollback the players know the next rolls and draws.
   1. Soft stops at dice and hidden info, so only Godmode shows them.
   2. Godmode is for trusting tables; that is accepted.
2. A long Godmode game replays many steps on a rollback.
   1. Still milliseconds; the fuzz tests play whole games faster.
3. A restart with a rules change ends the Godmode history there.
4. A new event kind must be sorted into the barrier list.

### Neutral

1. The protocol goes to version 5; older apps cannot join.
2. Replays (step 10) show the undone steps and the rollback.
