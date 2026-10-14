# Step 19. Bugs

Goal: every bug report is logged, fixed, and reviewed the same way.

This step never closes; new reports are appended.

## Rules for this step (owner, 2026-10-14)

1. Every bug report gets an item here: `19.N`, the next free number.
   1. Logged before any code change.
   2. Appended at the end; numbers never change.
2. An item has the sub-step format plus:
   1. Report: who, when, and what they saw.
   2. Cause: found in the match record or the code.
3. A fix follows the roadmap rules (`README.md`).
   1. Tests come with the fix.
   2. Docs in the same change: use cases, ADRs, wiki.
   3. Gates pass: `make lint`; `make frontend e2e` for frontend changes.
4. A screen change gets before and after screenshots.
   1. They go to `.cache/screenshots/19.N/before/` and `after/`.
   2. Never committed to the repo.
   3. The review message links each screenshot file.
5. Commit after each fixed bug; stop for the owner's review.
   1. The owner marks it solved; then "Owner notes".
6. Branch: the fix goes on the step branch in work.
   1. No `step-19` branch; fixes merge with that step.

---

### [x] 19.1 Roll results in the history

1. Goal: the history shows the roll result the game used.
2. Report (owner, 2026-10-14)
   1. Play test; record `2026-10-14-173137-g1_Gaijin_Oleg.fsrec`.
   2. Polycephalus: "the next attack roll this turn has -2".
   3. The history did not show the -2; it looked ignored.
3. Tasks:
   1. Find the cause in the match record.
   2. Show the die and the result when they differ.
   3. The strip's die shows the result.
4. Done when:
   1. Vitest: a changed roll shows both numbers; an unchanged one, one.
   2. e2e test: a changed attack roll in the strip and the panel.
   3. Use case UC-15 updated.
5. Cause
   1. Engine is right: Polycephalus's trigger adds a one-shot boost.
   2. The next attack roll uses it up (`useNextAttackRoll`).
   3. Record: die 6, `roll_resolved` 4; the 4 hit.
   4. The history used `stack_resolved`: the die before changes.
   5. `roll_resolved`, the result, was never shown.
   6. Same for every roll change, e.g. a +1: "rolled 5" for a 6.
6. Result (2026-10-14)
   1. A changed roll reads "Roll — attack roll (rolled 6 → 4)".
   2. The strip's die and line show the result: "rolls 4".
   3. An unchanged roll reads as before: "(rolled 6)".
   4. Engine, protocol and records unchanged.
   5. Screenshots: `.cache/screenshots/19.1/before/` and `after/`.
7. Owner notes (2026-10-14)
   1. Owner marks 19.1 solved.

### [x] 19.2 Heal popup after 1-Up

1. Goal: the heal popup shows the HP a player got back.
2. Report (owner, 2026-10-14)
   1. Play test; record `2026-10-14-173137-g1_Gaijin_Oleg.fsrec`.
   2. Lazarus died on his turn with 1-Up.
   3. The popup said +1 HP; his HP went up by 2.
3. Tasks:
   1. Find the cause in the match record.
   2. Every HP a death prevention gives back gets a `healed` event.
   3. 1-Up shows one popup with the full amount.
4. Done when:
   1. Engine test: 1-Up heals in one event for the full amount.
   2. Engine test: Guppy's Collar emits a `healed` event for 1 HP.
   3. Old records still replay.
5. Cause
   1. Step 196: Leech's combat damage leaves Lazarus at 0 of 2 HP.
   2. Step 198: 1-Up resolves; `healed` amount 1.
   3. `PreventYourDeath` sets HP to 1 with no `healed` event.
   4. 1-Up's heal to full then emits only the rest: +1.
   5. Broken Ankh and Guppy's Collar show no heal at all.
6. Result (2026-10-14)
   1. `PreventYourDeath` emits `healed` for the HP it gives back.
   2. New block `PreventYourDeathHealFull`: one `healed` to full.
   3. 1-Up uses it: popup "+2 HP" for Lazarus at 0 of 2 HP.
   4. Broken Ankh and Guppy's Collar now show "+1 HP".
   5. Events are not in the checksum; the owner's record replays.
   6. No frontend change; the popup is short, so no screenshots.
7. Owner notes (2026-10-14)
   1. Owner marks 19.2 solved.

### [ ] 19.3 Deck attack: the monster lands before the roll

1. Goal: on a deck attack the monster is in its slot before the dice roll.
2. Report (owner, 2026-10-14)
   1. A player attacks the top of the monster deck.
   2. The dice roll plays first.
   3. Only then the monster appears in the chosen slot.
3. Tasks:
   1. Find the cause in the code.
   2. A card put into a slot flies there and stays until the board redraws.
   3. Effects after it, like the attack roll, play over the placed card.
4. Done when:
   1. Vitest: a slot reveal gives a place effect before the dice.
   2. e2e test: the placed card is on the slot while the dice plays.
   3. Use case updated.
5. Cause
   1. Engine order is right: `placeRevealed` puts the card in the slot.
   2. `startAttack` then rolls; both land in one update.
   3. Playback plays an update's effects, then draws its board.
   4. So the dice plays before the board shows the monster in the slot.
6. Result (2026-10-14)
   1. New `place` effect: the card flies from its deck to the slot.
   2. It stays there, drawn by `Effects`, until the board redraws.
   3. Order now: reveal, the card lands in its slot, the attack roll.
   4. Shop refills get the same fly into their slot.
   5. Engine, protocol and records unchanged.
   6. Screenshots: `.cache/screenshots/19.3/before/` and `after/`.
