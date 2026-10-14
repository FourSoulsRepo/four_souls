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
