# Step 15. Visual refresh

Goal: a clearer table, using fouronlinesouls.net as a second reference.

Ideas: V-09, V-12, V-13.

---

### [ ] 15.1 Study the reference

1. Goal: a list of what to borrow (V-09).
2. Tasks:
   1. Open https://fouronlinesouls.net/ and start a game with a bot.
   2. Screenshots go to `.cache/screenshots/`, never into the repo.
   3. Compare with our table: layout, card sizes, panels, stack, prompts.
   4. List ideas to borrow; ideas only, no images, styles, or code.
   5. Compare with Hearthstone again; avoid the poker-table look (V-13).
3. Done when:
   1. **Owner** picks the ideas to apply.

### [ ] 15.2 Rework the table

1. Goal: apply the chosen ideas.
2. Tasks:
   1. Change layout and panels for 2, 3, 4 players.
   2. Keep V-01 – V-08 and touch support (B-06).
   3. Update `docs/wiki/Client.md`.
3. Done when:
   1. `make frontend e2e` passes.
   2. `make screenshots` shows every player count.

### [ ] 15.3 Attack animation

1. Goal: attacks look like Hearthstone (V-12).
2. Tasks:
   1. The attacker's card jumps at the monster, hits, and returns.
   2. Damage shows on the card that was hit.
   3. Only `transform` and `opacity`; follows animation speed (ST-05).
3. Done when:
   1. e2e screenshot of an attack in progress.
