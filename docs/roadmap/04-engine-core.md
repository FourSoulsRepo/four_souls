# Step 4. Engine core

Goal: the general rules of Four Souls, without real cards.
Engine is tested with small fake cards.

Ideas: A-05 – A-08, N-04, LR-04, LR-07, TS-05.

---

### [x] 4.1 Engine design ADR

1. Goal: agree on the design before code.
2. Tasks:
   1. State model: players, zones, card instances, counters.
   2. Flow: intent → validate → pending actions → rewrite → apply → events.
   3. Stack and priority model.
   4. Effect blocks and custom card hooks.
   5. Public API: new game, apply intent, allowed actions, view.
   6. Inputs from ADR 001 and the rules text (3.7).
3. Done when:
   1. **Owner** accepts the ADR.

### [x] 4.2 State and zones

1. Goal: the game state types.
2. Tasks:
   1. Players: HP, attack, coins, souls, items, hand.
   2. Zones: decks, discards, shop slots, monster slots.
   3. Only ordered slices; never decide by map order (A-08).
   4. Seeded RNG lives in the state; shuffle and dice use it.
   5. Deep copy and a stable checksum of the state.
3. Done when:
   1. Same seed gives the same shuffles in tests.

### [x] 4.3 Game setup and turn structure

1. Goal: a turn loop with no card effects.
2. Tasks:
   1. Setup: characters, starting items, decks, starting loot and coins.
   2. Start phase: recharge, then loot 1.
   3. Action phase: loot play, purchase, attack limits.
   4. End phase; next player.
   5. Win check: 4 souls.
3. Done when:
   1. A game with fake cards reaches a winner.

### [x] 4.4 Stack and priority

1. Goal: responses work like the rules say.
2. Tasks:
   1. Push, respond, pass, resolve.
   2. Dice rolls are stack items.
   3. All players pass in order before resolve.
   4. Adding to the stack while it resolves (R-07).
3. Done when:
   1. Tests cover nested responses.

### [x] 4.5 Pending actions and rewriting

1. Goal: the MTG Arena "whiteboard" pattern (A-05).
2. Tasks:
   1. Actions are written as pending before they apply.
   2. Hooks can prevent, replace, or cancel them.
   3. Ordering of several hooks is defined by a rule ID.
3. Done when:
   1. Prevent-damage and replace tests pass with fake cards.

### [x] 4.6 Combat and death

1. Goal: attacks, damage, death, rewards.
2. Tasks:
   1. Attack a monster or the monster deck.
   2. Attack rolls vs dice value; damage both ways.
   3. Monster death: rewards, souls.
   4. Player death: penalty, end of turn effects.
   5. Death cancels combat (R-07).
3. Done when:
   1. Combat tests match the rules text IDs.

### [x] 4.7 Effect blocks and card definitions

1. Goal: cards as data (A-06).
2. Tasks:
   1. Card definition: card ref → list of effect blocks.
   2. First blocks: loot, gain coins, deal damage, heal, roll table.
   3. Costs: pay coins, tap, destroy self.
   4. Targets: player, monster, item, choice.
   5. Triggers: on destroy, on death, on roll, on turn start.
   6. Custom Go hook for cards that do not fit.
3. Done when:
   1. Fake cards built only from blocks pass tests.
   2. Wiki: "Effect blocks" reference.

### [x] 4.8 Allowed actions and intents

1. Goal: what each player may do now (N-04).
2. Tasks:
   1. Allowed actions per player, incl. item abilities.
   2. Intent validation with a reason and rule ID.
   3. Rejected intents change nothing.
3. Done when:
   1. Every allowed action is accepted when sent.
   2. Every other action is rejected with a reason.

### [x] 4.9 View filter

1. Goal: one function builds every view (A-07).
2. Tasks:
   1. Views: player, viewer, judge.
   2. Hidden: other hands, deck order, RNG seed.
   3. Judge sees all hands; never the seed.
3. Done when:
   1. Tests prove no view leaks hidden data.

### [x] 4.10 Events and checksums

1. Goal: a step log for records and animations.
2. Tasks:
   1. Each applied step emits ordered events.
   2. Events are filtered by the same view rules.
   3. Each step carries a state checksum (RP-12).
3. Done when:
   1. Replaying events rebuilds the same checksum.

### [ ] 4.11 Situation payload and runner

1. Goal: payload tests (LR-04, LR-07).
2. Tasks:
   1. Format: title, setup, action, optional expect.
   2. Runner: allowed or not, reasons, rule IDs.
   3. Folder of payload files run by `go test`.
   4. Convert resolved cases from 3.8 into payloads.
3. Done when:
   1. Payload tests run in CI.
   2. Wiki: "Situation payload" reference.

### [ ] 4.12 Determinism and fuzz harness

1. Goal: catch bugs no one thought of (A-08, TS-05).
2. Tasks:
   1. Same seed and intents, run twice: equal checksums.
   2. Fuzz mode "all cards": random situations.
   3. Fuzz mode "focused": given cards always in play.
   4. Invariants: no panic, valid state, progress, no leaks.
   5. CI runs only the seed corpus (smoke run).
3. Done when:
   1. Local fuzz runs for 10 minutes without failures.
   2. Wiki: how to run fuzzing locally.

### [x] 4.13 Purchasing

1. Goal: buying items (R-SHOP); missing from the first plan of step 4.
2. Tasks:
   1. Declare a purchase; priority passes before the choice (R-SHOP-02).
   2. Choose a shop item or the top of the treasure deck.
   3. Price fixed at declaration time; 10¢ by default (R-SHOP-03).
   4. Pay and gain, or fail when unable to pay (R-SHOP-04).
   5. One purchase per turn by default (R-SHOP-05); refill the shop.
3. Done when:
   1. Tests cover buying from the shop and the deck, and failing to pay.
