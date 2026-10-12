# Step 15. Visual refresh

Goal: a clearer table, using fouronlinesouls.net as a second reference.

Ideas: V-09, V-12, V-13, V-15 – V-18.

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
   4. Simple but informative: nothing needed is hidden (V-18).
   5. Every action gets instant feedback (V-15).
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

### [ ] 15.4 Hearthstone-style card play

1. Goal: cards are played by dragging them to a place or target (V-16).
2. Tasks:
   1. Study Hearthstone screenshots: dragging, aiming, highlights.
   2. Drag a card onto the table or onto its target.
   3. Aiming arrow; valid targets glow, others stay plain.
   4. Click card, then click target, still works.
   5. Works with mouse and touch (B-06).
3. Done when:
   1. e2e test: a bomb dragged onto a monster damages it.
   2. e2e test: the same play by two clicks.

### [ ] 15.5 Dice throw

1. Goal: dice look real and the player throws them (V-17).
2. Tasks:
   1. **Owner** decides what happens if nobody presses "Throw".
   2. "Throw" button for the player who rolls; others watch.
   3. Fall, bounce, and roll animation ends on the engine's value.
   4. Only `transform` and `opacity`; follows animation speed (ST-05).
3. Done when:
   1. e2e screenshot of a die in the air.
   2. The rolled value matches the engine in every test.
