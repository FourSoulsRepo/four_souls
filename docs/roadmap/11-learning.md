# Step 11. Learning

Goal: players learn the rules inside the app.

Ideas: LR-01 – LR-06.

---

### [ ] 11.1 Rules viewer

1. Goal: our rules text in the app (LR-02).
2. Tasks:
   1. Render the engine rules files.
   2. Anchors by rule ID; links from reasons jump there.
   3. Search by text.
3. Done when:
   1. Every rule ID opens its section.

### [ ] 11.2 Sandbox

1. Goal: build a situation on the table (LR-01).
2. Tasks:
   1. Place cards, set HP, coins, souls, stack.
   2. Pick one action; engine answers allowed or not.
   3. Answer explains why, with rule links.
3. Done when:
   1. Use case: check a situation.

### [ ] 11.3 Share situations

1. Goal: payload and link (LR-04).
2. Tasks:
   1. Copy payload; paste payload.
   2. Copy link; base URL set in config until the site exists.
   3. Optional title and expected result.
3. Done when:
   1. A copied payload re-opens the same situation.

### [ ] 11.4 Simple bot

1. Goal: a practice opponent (LR-06).
2. Tasks:
   1. Random legal moves.
   2. Runs in process; never in network games.
3. Done when:
   1. Bot games finish without errors.

### [ ] 11.5 Practice games

1. Goal: two practice modes (LR-05).
2. Tasks:
   1. Free mode: real game vs bots, with hints.
   2. Seeded mode: known draws; only planned moves allowed.
   3. Seeded scenario file format.
3. Done when:
   1. One seeded scenario plays start to end.

### [ ] 11.6 Tutorial lessons

1. Goal: learn by doing (LR-03).
2. Tasks:
   1. Lesson format: setup, steps, hints.
   2. First lessons: first attack, buying, the stack, dying.
3. Done when:
   1. A new player finishes all lessons.
   2. Use case: complete a lesson.
