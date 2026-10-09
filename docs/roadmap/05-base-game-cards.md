# Step 5. Base Game cards

Goal: every Base Game card works and is tested.
Card batches are sized for one agent session each.

Ideas: R-05, TS-01 – TS-04.

---

### [x] 5.1 Card registry and stub generator

1. Goal: a stub for every card (TS-03).
2. Tasks:
   1. Engine registry: card ref → definition.
   2. Generator reads `card_db` Base Game data.
   3. Each stub holds the card text and a TODO marker.
   4. Generator never overwrites finished cards.
3. Done when:
   1. Every Base Game card has a stub.

### [x] 5.2 Card status report

1. Goal: see progress (TS-02).
2. Tasks:
   1. States: not started, implemented, tested.
   2. Report per set and per card type.
   3. Generated Markdown file in the engine docs.
3. Done when:
   1. Report regenerates with one command.

### [x] 5.3 Characters and starting items

1. Goal: all Base Game characters playable.
2. Tasks:
   1. Each character test-first (TS-01).
   2. Starting items and eternal items.
3. Done when:
   1. Report shows characters as tested.

### [x] 5.4 Loot cards

1. Goal: all Base Game loot cards.
2. Tasks:
   1. Coins, bombs, batteries, pills, tarot, trinkets.
   2. Split into batches of about 20 cards.
3. Done when:
   1. Report shows loot as tested.

### [x] 5.5 Treasure items

1. Goal: all Base Game treasures.
2. Tasks:
   1. Passive, active, paid, and roll items.
   2. Batches of about 20 cards.
3. Done when:
   1. Report shows treasures as tested.

### [x] 5.6 Monsters, bosses, events, curses

1. Goal: the whole Base Game monster deck.
2. Tasks:
   1. Stats, rewards, souls, "Indomitable"-style tags.
   2. Good events, bad events, curses.
   3. Batches of about 20 cards.
3. Done when:
   1. Report shows monsters as tested.

### [x] 5.7 Bonus souls

1. Goal: Base Game bonus souls, if the set has them.
2. Tasks:
   1. Check what the Base Game includes.
   2. Implement and test.
3. Done when:
   1. Report shows bonus souls as tested or not in set.

### [ ] 5.8 Interaction tests

1. Goal: card combinations work (TS-04).
2. Tasks:
   1. Combos named in rulings (R-07, digests from step 3).
   2. Payload tests for every resolved case (R-03).
   3. New disputed cases go to the open list for the owner.
3. Done when:
   1. All known combos have tests.

### [ ] 5.9 Full-game simulations

1. Goal: whole games without errors.
2. Tasks:
   1. Random legal moves for 2, 3, 4 players.
   2. Thousands of seeded games run locally.
   3. Fuzz "all cards" with the full Base Game.
3. Done when:
   1. Report shows 100% tested for the Base Game.
   2. Simulations finish with no invariant failures.

### [ ] 5.10 Guide for adding cards

1. Goal: anyone who clones the repo can add a card or a fan set (ADR 005).
2. Tasks:
   1. Wiki page "Adding a card": from stub to passing test.
   2. Examples: a loot card, an item, a monster, a custom hook.
   3. How to add a whole fan set as its own package.
   4. Every effect block documented with a one-line example.
3. Done when:
   1. A new contributor adds a test card by following only the guide.
