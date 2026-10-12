# Step 18. Play-test

Goal: real players play real games on every set done so far.
Moved from step 8 (8.4–8.6): the Base Game alone is too few cards.

Ideas: B-04, B-05.

---

### [ ] 18.1 New-user check

1. Goal: users can start the app (B-04, B-05).
2. Tasks:
   1. Give a new user only the release archive.
   2. Watch them start it on macOS, Windows, Linux.
   3. Fix `README.txt` and the wiki where they get stuck.
3. Done when:
   1. A new user starts the app with only the README.

### [ ] 18.2 Play-test and fixes

1. Goal: real games with real players.
2. Tasks:
   1. Play full games with 2, 3, 4 players.
   2. Use every set done in step 14.
   3. Collect issues; fix blockers.
   4. New rule disputes go to the open list.
3. Done when:
   1. **Owner** approves the release.

### [ ] 18.3 Docs check

1. Goal: docs match the released app.
2. Tasks:
   1. Every released feature has a use case.
   2. README is short and current.
   3. Accepted ADRs match the code.
3. Done when:
   1. Use case index has no "Planned" item for released features.

### [ ] 18.4 Multi-instance review

1. Goal: the "several copies" solution (UC-16) is reviewed and kept or reworked.
2. Tasks:
   1. Review `internal/instance`, `main.go`, `app.go` from commit 62a610f.
   2. It was built outside the roadmap; it has no ADR yet.
   3. Check it on macOS (`open -n`), Windows, and Linux.
   4. Check edge cases: crash, all 16 slots busy, two copies hosting.
   5. Check settings seeding: no copy takes another copy's seat.
   6. **Owner** decides: keep, simplify, or replace.
   7. Write an ADR for the result; update UC-16 if behavior changes.
3. Done when:
   1. ADR accepted; UC-16 matches the code.
