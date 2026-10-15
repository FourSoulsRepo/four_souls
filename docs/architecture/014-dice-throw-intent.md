# 14. Dice throw intent

* **Status:** Accepted
* **Date:** 2026-10-15
* **Authors:** @HardDie

---

## Context

1. Dice should look real; the player throws them (V-17, 15.5).
2. Before, the engine drew a roll's value the moment the roll was made.
   1. The client only showed a die with a value it already had.
   2. The roller took no part in the roll.
3. Nobody may see a value early, the roller included.
4. Records must replay exactly (ADR 007).
5. Owner decisions (2026-10-15):
   1. A new throw intent from the roller.
   2. The engine draws the value only when the intent arrives.
   3. No timer by default; the response timer throws for an idle player.
   4. Kicked seats throw at once; bots will too (step 11).

## Considered options

1. **Throw intent; the value is drawn when it arrives**
   1. Won. No value exists before the throw, so nothing can leak.
   2. The record keeps the intent; a replay draws the same value.
2. **Value drawn at once; the throw only starts the animation**
   1. Lost. The value sits in the state before the throw.
   2. A bug in a view could show it early.
3. **Client-side throw with a client-chosen value**
   1. Lost. Clients could cheat; replays would need the value.

## Decision

Use option 1.

1. Prompt and intent
   1. `PromptThrow` names the roller; its `Text` is the roll's label.
   2. The roller answers with `IntentThrow`; it has no arguments.
   3. `Allowed` lists only the throw while it is open.
2. Rolls on the stack
   1. A roll goes on the stack with no value (`Roll` 0).
   2. Its push queues `ActThrow` for the roller.
   3. Queued actions run before anyone gets priority.
   4. So nobody can respond to a roll before its throw.
3. The throw
   1. `IntentThrow` draws the value (`d6`); forced rolls still decide it.
   2. Then `dice_rolled` is emitted with the roll's stack ID.
   3. Many dice (Mutant Spider) are drawn with the throw.
   4. A roll that left the stack before its throw is never thrown.
   5. A game-controlled roll (no player) is thrown at once.
4. Rerolls
   1. A reroll clears the value and queues a throw.
   2. The roll's own player throws again; the event says `reroll`.
5. Covered rolls
   1. Attack, reward, effect (`Roll`), and many-dice rolls.
   2. The first-player roll stays automatic until 15.43.
6. Tests
   1. `enginetest` helpers and `settle` throw for the roller.
   2. The fuzz harness throws through `Allowed`.
7. Records
   1. Engine version 0.2.0 is the first with throws.
   2. `engine.CanReplay(version)` refuses records of older engines.
   3. The record format stays 1: old records stay readable.

## Consequences

### Positive

1. No roll value exists before the roller throws.
2. The client can show a real throw: button, waiting, animation.
3. Replays stay exact: the intent is recorded, the value comes from the seed.

### Negative and risks

1. Every roll needs one more intent; a game has more steps.
2. Records and saves of engine 0.1.x cannot replay or load.
3. Until the client knows the throw (15.61), the app waits at every roll.

### Neutral

1. The server's forced answer throws for kicked seats and idle players.
2. Its tests and protocol docs come with 15.60.
