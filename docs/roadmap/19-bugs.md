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
