# Step 24. Situation tests

Was step 23; renumbered at the merge with bots (owner, 2026-11-18).

Goal: most card tests are situation files, not Go code.

Ideas: TS-13, LR-04, R-03.

## Why (2026-11-17)

1. The situation format exists since 4.11; the stack came in 12.5.
   1. `pkg/rules_engine/testdata/situations/` holds 6 files.
   2. They use the test sets only, not real cards.
2. Card tests are Go: about 930 tests in `pkg/rules_engine/cards/`.
   1. They already build the table from `engine.SituationSetup`.
   2. So the setup half of the language is shared today.
3. Why new tests stayed in Go:
   1. One action only; 159 card tests do several.
   2. No forced dice; `ForceRolls` is in about 300 lines.
   3. No deck tops (`PutOnDeck`, 83), counters (35), turn counters.
   4. Checks are only cents, hand size, damage, souls, dead.
   5. No check of prompt options (55 tests) or resolve order (76).
   6. Each set's package has no situation runner.
   7. `Adding-a-card.md` shows only the Go way.
4. About 380 card tests are one action and a plain state check.
   1. They fit the format once checks and dice are added.
5. Gains:
   1. The owner reads and reviews a test without Go.
   2. Tests are grouped by card and rule, so their purpose shows.
   3. A file names its source; CLAUDE.md 5.25 asks for that.
   4. The sandbox (LR-01) and tests share one format.

## Rules for this step

1. Branch: `step-23-situation-tests`, from `main`.
2. Old payloads keep working: website links hold them (12.7).
   1. `action` stays; `steps` is added beside it.
   2. A new field is optional.
3. An `expect` value comes from the rules, never from a run (5.25).
   1. A tool may write `setup` and `steps`; never `expect`.
4. Every new field: a line in `docs/wiki/Situation-payload.md`.
5. Unknown fields stay errors that list the options.
6. Gates: `make lint test`; `make site` when the format changes.
7. Sub-steps run in fresh subagents, one at a time (5.21).

## Owner answers (2026-11-17)

1. Move, not copy: only a test a file says exactly the same.
   1. Its Go test is deleted in the same change.
2. Go and JSON side by side for one card: yes.
   1. A test only Go can say stays in Go.
3. A new test: simple is a file; complicated is Go.
4. Website examples from test files: no, not for now.
5. Files are grouped by purpose, never one big folder.
6. A step `expect` with no `allowed` expects a refusal; keep it.
   1. Files will come from a visual constructor, not typed by hand.
7. No responses in the middle of a stack within `steps`.
   1. The goal: a player builds the table position they have now.
   2. A stack in play is set with `setup.stack` and `priority`.
8. A moved file may check more than its Go test, from the card text.
   1. More checks, not fewer.
9. Bulk move: yes, one card family per sub-step, with the dump (2026-11-18).
10. Deck bottom and roll value checks come first (24.11).
11. A table of fixed cases: one file per case (ADR 034 1.10).

## File layout

1. A card's tests: `cards/<set>/testdata/situations/<card>/<case>.json`.
   1. `<card>` is the card ref, as its `<card>_test.go`.
   2. `<case>` says what is checked, e.g. `pays-5-cents.json`.
2. Two or more cards of one set together: `.../situations/interactions/`.
   1. Named after the cards, e.g. `lucky-foot-bomb.json`.
3. Cards of several sets: `cards/testdata/situations/interactions/`.
4. Engine rules, test sets only: `testdata/situations/<area>/`.
   1. `<area>` is the rule ID area in lower case: `atk`, `stack`, `abil`.
   2. The 6 files today move there (24.2).
5. A resolved disputed case (R-03) goes where its card or rule goes.
   1. Its `source` names the ruling.
6. A folder holds only `.json` files; no deeper levels.

### [x] 24.1 ADR: situation files as card tests

1. Goal: one ADR fixes when a test is a file and when Go.
2. Tasks:
   1. Written in the main conversation (5.21.2).
   2. The owner answers and the file layout above.
   3. The `source` field and its forms.
   4. "Simple" and "complicated", with examples of each.
3. Done when:
   1. The ADR is Accepted and in `INDEX.md`.

### [x] 24.2 A runner per set, and the `source` field

1. Goal: each set's `go test` runs its own situation files.
2. Tasks:
   1. `enginetest.RunSituations(t, files, sets...)`, one subtest per file.
   2. Each set package embeds its `testdata/situations/` tree.
   3. The subtest name is the path: `<card>/<case>`.
   4. A file outside the layout fails the run.
   5. A test file needs `expect` and `source`; a sandbox payload needs neither.
   6. A mismatch prints the file, the step and every failed check.
   7. Move the 6 root files into rule area folders; run them with it.
3. Done when:
   1. A b2 file with a wrong `expect` fails `go test ./cards/b2`.
   2. A test file with no `source` fails.

### [x] 24.3 Setup: dice, decks, counters, turn

1. Goal: the setup says what Go tests set by hand today.
2. Tasks:
   1. `rolls`: the next dice results, in order (as `ForceRolls`).
   2. `decks`: cards on top of the loot, treasure, monster decks, top first.
   3. `counters` on a player's card, a monster or an item.
   4. Monster `damage` in a slot.
   5. Turn: `attacks`, `purchases`, extra loot plays, `phase` if needed.
   6. List the other hand-set fields in card tests; add or reject each.
3. Done when:
   1. A test of each field, from card text, passes.
4. Added (count: lines in `cards/**/*_test.go`, 2026-11-17):
   1. `rolls`: `ForceRolls`, 319.
   2. `decks`: `PutOnDeck`, 89.
   3. `counters`, per player and in `setup`: `AddCountersTo` 21, `AddNamedCounters` 15.
   4. Also `Object(…).Counters =`, 8; named counters by `name`.
   5. `monster_damage`: `Object(…).Damage =`, 17.
   6. `extra_loot_plays` per player: 11.
   7. `attacks`: `Turn.Attacks =` 6, `++` 4.
   8. `purchases`: `Turn.Purchases =`, 2; the pair of `attacks`.
   9. `eternal` per player: `.Eternal = true`, 4.
5. Already in the format:
   1. `Players[i].Cents =` 22, `.Damage =` 3, `.Dead =` 1.
   2. `Turn.LootPlays =` 7: `loot_plays`.
   3. Souls appended to `InPlay`, 8: `souls`.
   4. `.Hand =` 6: `hand`; mid-test changes wait for `steps` (24.4).
   5. `.Charged =` 5: `deactivated`.
   6. `Rooms`, `Character` (4, 2): set while building the setup.
6. Rejected:
   1. `phase`: no card test sets `Turn.Step`.
   2. `AddPlayerCounters` 15, `AddPresentCounters` 5: challenge only.
   3. A challenge can't start from a setup; those tests stay Go.
   4. `Boosts`, 4: an engine struct, challenge tests only.
   5. `.Flipped =`, 3: Anima Sola only.
   6. `MustAttack`, `RandomAttackMade`: 1 each, one challenge test.
   7. `Decks[]` cut or emptied, 2: one helper each.
   8. `Challenge`, 1: the sim test.
7. Files: one per field, from card text.
   1. Moved: TestIITheHighPriestess, TestChest2, TestTechX.
   2. New: Gemini, Champion Belt, There's Options, Cain, XX. Judgement.
   3. New in `r`: Soul of Wrath (bonus soul counters); `r` gets a runner.
8. A test file is read strictly: an unknown field fails.

### [x] 24.4 Several steps

1. Goal: a file plays several actions in a row.
2. Tasks:
   1. `steps`: a list of actions, each with its own `choices`.
   2. `action` is one step; a file has one of the two.
   3. A step may expect a refusal; later steps still run.
   4. Step kinds as `enginetest`: `throw`, `choose`, `reveal`.
   5. `expect` may sit on a step, checked after it.
   6. `setup.priority`: the seat that holds priority (owner, 2026-11-17).
      1. Empty stack: the players before it in turn order passed (R-PRIO-02).
      2. The engine records those passes (R-PRIO-05).
      3. With a stack, it overrides the top item's owner.
      4. Rewrite `cain/plays-loot-on-another-turn.json` with it.
      5. Its dice roll on the stack was only for priority.
3. Done when:
   1. An attack, a roll and an end turn run from one file.
   2. The Cain file uses `priority` and has no stack.
4. Done (2026-11-17):
   1. Steps settle to the open action phase.
   2. They stop early only for a next `choose` or `throw` step.
   3. `action` keeps its old settle; old payloads answer the same.
   4. Moved: TestEdensBlessing (2 files), TestDagaz, TestPolydactyly.
   5. Moved: TestBlankCard, TestGoldenRazorBlade.
   6. New: Eden's Blessing after a fly: attack, throw, end turn.

### [x] 24.5 Richer checks

1. Goal: the checks Go card tests use most.
2. Tasks:
   1. Player: `hp`, `atk`, `items` in play, `hand` cards, `deactivated`, `counters`.
   2. Monster slot: the card, `damage`, `dead`, `counters`.
   3. Discards and deck tops; `souls` by card.
   4. Turn: active player, loot plays, attacks left.
   5. Game: `over`, `winners`.
   6. A missing field means "not checked", as today.
3. Done when:
   1. Each check has a passing and a failing test.
4. Added (count: lines in `cards/**/*_test.go`, 2026-11-17):
   1. Player `hp`: `PlayerHP` 35; `atk`: `PlayerATK` 33.
   2. `items`: `hasItem` 89, `InPlay` 115.
   3. `hand_cards`: `handCards` 11, `Hand[` 5.
   4. `deactivated`: `.Charged` 76.
   5. Player `counters`: `CountersOf` 39.
   6. `soul_cards`: souls by card; `SoulValue` 56 stays `souls`.
   7. `monsters`: `Monsters[` 86, monster `.Damage` 52.
   8. Also monster `hp`, `atk`, `evasion`: `Evasion` 38, `MonsterATK` 15, `HP` 9.
   9. Top-level `counters`, as `setup.counters`: slots and bonus souls.
   10. `decks`: `DeckTop` 48; `discards`: `inDiscard` 39, `Discards[` 21.
   11. `active` 26, `loot_plays` 23, `attacks` 33 (`Turn.…` reads).
   12. `over` 29, `winners` 13.
5. Shape:
   1. Names follow `setup`; `hand` and `souls` stay counts.
   2. Card lists: every card, any order; `[]` expects none.
   3. `decks`: only the top cards, top first.
   4. `monsters[]`: `card`, optional `slot` from 1; `dead` means it left play.
   5. A counter check with `count` 0 expects none.
6. Rejected:
   1. `purchases`: no card test compares `Turn.Purchases`.
   2. `extra_loot_plays`: 4 compares.
   3. `MustAttacks` 15, `AttackRolls` 18: mid-attack state; Go for now.
   4. No own check for `.Character` (17): mostly its `.Charged`, now `deactivated`.
   5. Curses 6, `Eternal` 5, shields 3, shop slots 4: too few.
7. Files:
   1. Moved: TestXIStrength, TestTheHabit, TestTheShovel, TestTarotCloth.
   2. Moved: TestGlassCannon (2 files), TestSmelter.
   3. Moved in `r`: TestMegaSatanWins, TestRedChampions, TestSumptorium.
   4. Soul of Wrath's file also checks `soul_cards`.

### [x] 24.6 Checks on prompts and order

1. Goal: a file checks what a player is asked and what resolves first.
2. Tasks:
   1. `prompt`: who is asked, and the options by plain label.
   2. `resolved`: cards in the order they resolved.
   3. `events`: kinds that must or must not happen, with their player.
   4. An `activate` step may name a room in a room slot (owner, 2026-11-17).
      1. Only the active player may use it (R-CARD-33, R-CARD-36).
      2. Move the Bomb Bum and Beggar tests to files.
      3. A `setup.stack` `activate` item may name a room too, as an item.
      4. Its owner is the active player, who activated it.
      5. Same rules as a stacked item ability: costs paid, nothing triggers.
3. Done when:
   1. A Lucky Foot and Bomb! file checks the order.
   2. A Bomb Bum file activates the room.
   3. A file starts with a room's ability on the stack and resolves it.
4. Added (count: lines in `cards/**/*_test.go`, 2026-11-17):
   1. `prompt.options`: `tb.Options()` 40, `PlainOptions` 1.
   2. `prompt.player`: `G.Prompt()` 41; most are loops that pass.
   3. `resolved`: `ResolvedAt` 4, `FirstAttackRoll` 7, `resolvedItems` 4.
   4. `events`: `enginetest.Has` 9, `Kind == engine.Ev` 38.
5. Shape:
   1. `prompt`: `player` and `options`, any order; a field left out is not checked.
   2. A prompt check stops its step at the question, as a next `choose` does.
   3. The top-level `prompt` does that for the last step and for `action`.
   4. `resolved`: an ordered subsequence; other items may resolve between.
   5. A card names its loot, ability or trigger; other items by kind.
   6. Kinds: `dice_roll`, `attack_roll`, `damage`, `death` and the rest.
   7. A full list was rejected: Bomb!'s damage item would be in every file.
   8. `events`: `kind`, optional `player`, `none: true` for "must not".
   9. Event kinds are the engine's names; an unknown one lists them all.
   10. A step's checks read its own events; the top one every step's.
6. Rejected:
   1. Roll values: `resolvedRoll` 13; `resolved` names items, not values.
   2. Event fields past kind and player (`By`, `Card`, `Amount`): too varied.
   3. Option order (`opts[0]`, 5): the display order; the set is checked.
   4. Prompt text keys (1) and `Look` (3): too few.
7. Rooms work as items (owner, 2026-11-17):
   1. An `activate` step finds a room in a room slot too.
   2. The engine refuses another player by R-CARD-36, as in a game.
   3. A stacked `activate` uses the `trigger` kind's slot lookup.
   4. Its owner must be the active player; else an error names R-CARD-36.
8. Files:
   1. `b2`: `interactions/lucky-foot-bomb.json` checks the order.
   2. Moved: TestTheBattery, TestBatteryBumOffersOnlyActiveItems (prompt).
   3. Moved: TestXXITheWorld (events).
   4. Moved in `r`: TestBombBum, TestBeggar, TestRedVise.
   5. New in `r`: Bomb Bum refused to Cain; Beggar from the stack.
   6. New in `r`: Red Vise spares a player with 2 loot cards.
   7. `stack/three-stacked-events.json` also checks `resolved`.

### [x] 24.7 Dump a Go test's setup as a file

1. Goal: a tool writes `setup` and `steps` from a Go test run.
2. Tasks:
   1. `enginetest.Table` records its setup, rolls, decks and steps.
   2. An env var writes them to a file per test.
   3. It never writes `expect` (rule 3).
   4. Hand-set state it can't express is listed as a TODO.
3. Done when:
   1. A dumped b2 test, with `expect` from its Go checks, passes.
4. Done (2026-11-18):
   1. `SITUATION_DUMP=<dir>` makes each `Table` write a draft at its test's end.
   2. Path: `<dir>/<package folder>/<TestName>.json`; relative from the repo root.
   3. It writes `title`, `sets`, `setup`, `steps`, an empty `source`; no `expect`.
   4. `rolls`: every die thrown, forced or random, from the events.
   5. Hand-set state before step 1 goes into setup where a field says it.
   6. A rebuild from the dumped setup is compared with the test's table.
   7. What still differs is a `todo` line; so is a write between steps.
   8. Direct `G.Apply` calls: read from the engine's applied-intent log by reflect.
   9. `todo` is no situation field; the author removes it.
   10. Moved: TestDinga, dumped, `expect` by hand from the card text.
   11. Wiki: `Situation-payload.md`, "Dump a Go test".
   12. Lint: `enginetest` may import `os`; only tests import it (depguard nolint).

### [x] 24.8 Pilot: move one group

1. Goal: learn the cost on real tests before any bulk move.
2. Tasks:
   1. Move simple tests of `cards/b2/interactions_test.go` and one card family.
   2. Each `expect` from the Go test's checks and its source comment.
   3. Count what fit, what needed a format change, what stayed Go.
   4. A rule bug found on the way is a step 19 item.
   5. Report the numbers to the owner; the bulk move waits for a yes.
3. Done when:
   1. The moved files pass; the Go tests they replace are gone.
   2. The report is in this file.
4. Group (2026-11-18): `interactions_test.go` and the tarot loot cards.
   1. Tarot: 0. The Fool to XXI. The World, one mechanic each.
   2. Why: 24 tests of one card kind; plain plays, rolls, deck looks.
5. Counts: 35 tests looked at (11 interactions, 24 tarot).
   1. Moved: 21 (8 interactions, 13 tarot), as 23 files.
   2. Split: 2, a file each; the rest stays Go.
   3. `TestRulingDamageToDeadPlayer`: damage read before the play stays Go.
   4. `TestXIXTheSun`: the loot deck bottom stays Go.
   5. Stayed Go, complicated: 4.
   6. Table of cases: `TestXWheelOfFortune` (ADR 034 1.3).
   7. Loop and a recharge between steps: `TestRulingTwoOfClubsTwice`.
   8. Views: `TestXVIIITheMoonShowsTheLook`, `TestXIITheHangedManTellsEveryoneTheChoice`.
   9. Needed a format change: 8.
   10. A roll's resolved value: `TestITheMagician`, `TestRulingRollAboveSix`.
   11. A deck bottom check: IV, IX, XVIII (`testOneOnTop`), `TestXIITheHangedMan`.
   12. Event fields (`By`, `Amount`, `Text`, `Private`): The World, The Moon's bottom.
6. Effort
   1. A plain test: 2-4 minutes; dump, title, source, expect, a mutation run.
   2. A test with a death or a turn end: 5-10 minutes.
   3. Steps settle into the next turn; dead players are back by then.
   4. So such a check reads a step's `died` event instead.
   5. `action` stops at the first window; a turn change needs `steps`.
   6. Every file failed with one expected value changed, then passed.
7. Drafts: 22 of 25 files kept the dumped setup and steps.
   1. Conjoined Fatty: more 6s in `rolls`, so a going-on attack would kill it.
   2. The Stars: a known treasure on the deck, so `items` can name it.
   3. Judgement's eternal soul: set after a direct `G.Apply`; the draft lacked it.
   4. Every draft needed `title`, `source`, `expect` by hand; never `todo` left.
8. Checks added from the card text (owner answer 8):
   1. Monster Manual tapped; Glass Cannon took Cain's Breakfast.
   2. Bomb! left 1 damage on Lust; Bum Friend put the Dime on the loot deck.
   3. Steamy Sale: no item gained; The Tower on 3 kills the Leech too.
   4. The Fool in the loot discard; Justice leaves Cain as he was.
   5. The Devil leaves Cain The D6; Judgement leaves Monstro.
   6. Judgement asks Cain, the soul's controller; The Sun is not discarded.
9. No rules bug found; `cardgen -report` unchanged.
10. Recommendation for the bulk move:
   1. Yes, one card family per sub-step, with the dump.
   2. About two thirds fit today (23 of 35 with splits).
   3. First add a deck bottom check and a roll value check: 6 more here.
   4. Event fields and views stay Go.
   5. Owner: may a table of fixed cases become a file per case?
   6. The Wheel of Fortune would be 6 files.

### [ ] 24.9 Docs

1. Goal: a simple new card test starts as a file.
2. Tasks:
   1. `Adding-a-card.md`: simple is a file, complicated is Go.
   2. `Situation-payload.md`: every new field and the file layout.
   3. `CLAUDE.md` Commands: how to run one set's files.
3. Done when:
   1. A simple new card from the guide needs no Go test.

### [x] 24.10 Priority windows

1. Goal: a file starts in any moment where a player may respond (owner, 2026-11-18).
2. Why:
   1. A file always starts in the open action phase, nothing declared.
   2. There the active player can't pass (R-TURN-09).
   3. So `priority` on an empty stack there is a state no game reaches.
   4. Other windows have an empty stack and priority passes around.
   5. Start phase (R-TURN-03), loot step (R-TURN-04), end phase (R-TURN-10).
   6. A declared attack (R-ATK-02) or purchase (R-SHOP-02).
3. Tasks:
   1. `setup.phase`: `start`, `loot`, `action` (default) or `end`.
   2. `setup.declared`: `attack` or `purchase`, before the target or item.
   3. `priority` on an empty stack works in each window above.
   4. In the open action phase, nothing declared: an error citing R-TURN-09.
   5. With a stack, `priority` works as now (24.4).
   6. Rewrite the Cain file into a real window, e.g. the loot step.
   7. A file per window, expected values from the rules.
   8. Wiki: `Situation-payload.md`; ADR 034 if it drifts.
4. Done when:
   1. Each window has a file where a non-active player acts.
   2. Empty-stack `priority` in the open action phase fails, naming R-TURN-09.
5. Done (2026-11-18):
   1. `phase: start`: after the recharge and start-of-turn triggers (R-TURN-03).
   2. `phase: loot`: after the active player looted 1 (R-TURN-04).
   3. `phase: end`: after the end-of-turn triggers (R-TURN-10).
   4. Start and loot: no loot play, attack or purchase yet (R-TURN-05, R-TURN-07).
   5. Their `loot_plays`, `attacks`, `purchases` are errors.
   6. `declared`: used up as on a declaration; nothing paid, nothing triggers.
   7. `declared` needs the action phase, an empty stack, one left.
   8. An `attack_roll` stack item needs the action phase (R-ATK-01).
   9. Everyone passes: the game moves on (R-STACK-07), as in a game.
   10. Steps settle to the open action phase; from `end`, the next player's.
   11. Empty-stack `priority` in the open action phase: an error, R-TURN-09.
   12. Engine: `spendAttack`, `spendPurchase` split from the declarations.
   13. Cain: the loot step; Bomb Bum (`r`): the loot step too.
   14. New in `b2`: Mr. Boom (start), Jawbone (end, purchase).
   15. New in `b2`: Monster Manual (declared attack).
6. Follow-up (owner, 2026-11-18): `declared` with a stack.
   1. A response to a declaration is a real table: allow it.
   2. E.g. Isaac declared an attack; Cain's Bomb! waits on the stack.
   3. The stack resolves first; then the target or item is chosen.
   4. Replaces Done item 7's "an empty stack" for `declared`.
   5. A file per kind: a declared attack and a declared purchase.
   6. Done (2026-11-18):
      1. `declared` with `setup.stack` is allowed; the stack resolves first.
      2. Then the active player picks the target or item (R-STACK-07).
      3. An `attack_roll` item with `declared` is an error.
      4. Attack: R-ATK-02, R-ATK-11; purchase: R-SHOP-01.
      5. New in `b2`: Bomb! (declared attack), Jawbone (declared purchase).

### [x] 24.11 Deck bottom and roll value checks

1. Goal: the checks the pilot (24.8) missed most.
2. Tasks:
   1. `expect.decks`: the bottom cards too, bottom last.
   2. A check of the value a roll resolved as.
   3. Shape from the Go tests: `testOneOnTop`, `resolvedRoll`, `ResolvedAt`.
   4. A passing and a failing test for each.
   5. Move the pilot's 6 tests that wait for these (24.8 item 6).
   6. The Emperor, The Hermit, The Moon, The Hanged Man.
   7. The Magician, `TestRulingRollAboveSix`.
   8. The Wheel of Fortune: one file per outcome (owner answer 11).
   9. The dump: a hand-set field before a direct `G.Apply` (24.8 item 7.3).
   10. Wiki `Situation-payload.md`; ADR 034 if it drifts.
3. Done when:
   1. The 7 tests above are files; their Go tests are gone.
4. Done (2026-11-18):
   1. `expect.deck_bottoms`: a deck's bottom cards by deck name.
   2. Compared in any order: several cards go to the bottom at random (Q-22).
   3. So not "bottom last"; no rule gives that order.
   4. `expect.resolved_rolls`: each roll's final value, in order.
   5. Read from `roll_resolved` events; attack rolls too; `[]` expects none.
   6. Moved: The Emperor, The Hermit, The Moon (`testOneOnTop` is gone).
   7. Moved: TestXIITheHangedMan, TestITheMagician.
   8. Moved: TestRulingRollAboveSix, as `interactions/spoon-bender-book-of-belial-roll-above-six.json`.
   9. Its two items start on `setup.stack`, Spoon Bender on top (owner answer 7).
   10. Its Go check never ran: the stack was empty by then.
   11. Wheel of Fortune: 6 files, `x_wheel_of_fortune/on-<n>-….json`.
   12. Also XIX. The Sun's Go bottom check: now in its file; Go test gone.
   13. Stay Go: The Moon's and The Hanged Man's view and event-field tests.
   14. Dump: hand-set state before a first direct `G.Apply` is in setup.
   15. Read from the engine's save origin, the state before its first intent.
   16. Unlocked across sets (count only): about 19 roll and 10 bottom tests.
   17. No rules bug found; `cardgen -report` unchanged.

## Bulk move rules (owner, 2026-11-18)

1. One card family per sub-step, in a fresh agent (CLAUDE.md 5.21).
2. The way of 24.8: dump, then `title`, `source`, `expect` by hand.
3. Owner answers 1, 8, 11; ADR 034 Decision 1 and 4.
4. Each file fails with one expected value changed, then passes.
5. No format or engine change in a bulk sub-step.
   1. A test that needs one stays Go; count it and name the change.
6. A rules bug: a step 19 item first (CLAUDE.md 5.17); not fixed there.
7. Rerun `go run ./cmd/cardgen -report`; no card drops from "tested".
8. Report under the sub-step: moved, split, stayed Go (why), changes needed.
9. Families from the 2026-11-18 count of card kinds in `_test.go` files.
10. "Other": test files with no card kind (rules, helpers, combos).

### [x] 24.12 Bulk: b2 loot cards

1. Goal: Wildcard, Trinket, Nickel, Bomb, Battery, Butter Bean (about 33 tests).
2. Counts (2026-11-18): 26 tests left in the family's `_test.go` files.
   1. 22 of Trinket, Nickel, Bomb, Battery, Butter Bean cards.
   2. 4 tarot (Wildcard): the rest moved in 24.8 and 24.11.
   3. So fewer than 33: the count was taken before those moves.
3. Moved: 18 tests, as 19 files.
   1. TestANickel, TestADime, TestCounterfeitPenny.
   2. TestBomb, TestGoldBomb, TestBloodyPenny.
   3. TestBrokenAnkh: 2 files, one per roll (owner answer 11).
   4. TestButterBean, TestGuppysHairball, TestSwallowedPenny.
   5. TestCainsEye, TestGoldenHorseshoe, TestPurpleHeart.
   6. TestCurvedHorn, TestMegaBattery.
   7. TestLilBattery, TestLilBatteryOffersOnlyActiveItems, TestLilBatteryNeedsASpentItem.
4. Split: 0.
5. Stayed Go, complicated (ADR 034 1.5): 8.
   1. Views: TestCainsEyeShowsTheTopCardToItsOwner.
   2. Event fields and views: TestCainsEyeTellsEveryoneTheChoice.
   3. View ATK mid-attack: TestCurvedHornShownForTheFirstRoll, TestCurvedHornGoneAfterAMiss.
   4. Tarot, as in 24.8 and 24.11: The Hanged Man, The Moon (2), The World.
6. Needed a format change: 0.
7. Drafts: 15 of 19 files kept the dumped setup.
   1. Butter Bean: A Dime!! starts on `setup.stack`, `priority` 1 (owner answer 7).
   2. Its draft had a mid-stack play and a hand-set extra loot play.
   3. Cain's Eye, Golden Horseshoe, Purple Heart: two known deck tops in `setup.decks`.
   4. So the file names the bottomed card and the next top.
   5. Lil Battery's refusal: the draft had no step; one `play` added.
   6. A one-step draft became `action`: it stops before the turn's end.
   7. So Broken Ankh's file still sees 1 HP before the end-of-turn heal.
8. Checks added from the card text (owner answer 8):
   1. Coins, batteries, bombs: the card is in the loot discard.
   2. Bomb! gives Isaac the Fly's 1¢; Gold Bomb!! leaves Isaac unhurt.
   3. Bloody Penny: `resolved` checks the loot before the death penalty.
   4. Broken Ankh on 6: `turn_ended_early`, no `died`; on 5: the Ankh is lost.
   5. Cain's Eye: Isaac loots the card under the bottomed one.
   6. Horseshoe, Purple Heart: the next card is on top.
   7. Curved Horn: 3 rolls resolved, ATK back to 1, 2 loot from the Fatty.
   8. Hairball: 3 rolls, the Leech dies; Swallowed Penny: 1 damage, the Fly dies.
   9. Lil Battery: every option, as Battery Bum's file; the refusal keeps the card.
9. Go tests that did not check what they claimed: 1.
   1. TestBloodyPenny's comment: the loot resolves before the penalty (R-DEATH-13).
   2. Its checks (hand 1, 2¢) hold in either order; the file checks the order.
10. No rules bug found; `cardgen -report` unchanged.

### [x] 24.13 Bulk: b2 passive treasures

1. Goal: Passive Treasure (about 58 tests).
2. Counts (2026-11-18): 58 tests in 50 `_test.go` files of the family.
   1. Eden's Blessing, Polydactyly, Tarot Cloth, The Habit: moved before.
3. Moved: 52 tests, as 52 files.
   1. TestBabyHaunt, TestBellyButton, TestBobsBrain, TestBreakfast, TestDinner.
   2. TestBrimstone, TestIpecac, TestMeat, TestSynthoil, TestEmptyVessel.
   3. TestBumBo, TestCambionConception, TestChampionBelt, TestTheDeadCat.
   4. TestChargedBaby, TestCheeseGrater, TestEyeOfGreed, TestShinyRock.
   5. TestDadsLostCoin, TestDadsLostCoinDeclined, TestDadsLostCoinTwice.
   6. TestCurseOfTheTower, TestDaddyHaunt, TestDryBaby, TestFannyPack.
   7. TestDarkBum, TestMomsCoinPurse, TestMomsPurse, TestThePolaroid.
   8. TestDeadBird, TestFinger, TestMomsRazor, TestSpiderMod, TestTheRelic.
   9. TestMomsBox, TestMomsBoxAnswersAreNotCards.
   10. TestGoatHead, TestStarterDeck, TestMonstrosTooth, TestRestock.
   11. TestGreedsGullet, TestShadowTakesThePenalty, TestSuicideKing.
   12. TestSacredHeart, TestSteamySale, TestTheMidasTouch, TestTrinityShield.
   13. TestTheBlueMap, TestTheCompass, TestTheMap, TestTheD10, TestTheD10Event.
4. Split: 2, as 3 files.
   1. TestBumBoWithCounterfeitPenny: a file per order (owner answer 11).
   2. Its option outcomes are a view; they stay Go.
   3. TestGuppysCollar: the file checks the prevented death.
   4. The heal's amount stays Go: TestGuppysCollarHealsBack1HP.
5. Stayed Go, complicated (ADR 034 1.5): 4, all views.
   1. TestChampionBeltSourceDuringFirstRoll, TestChampionBeltInTheOdds.
   2. TestDeadBirdShowsTheCardToTwo, TestTheresOptions.
6. Needed a format change: 0 whole tests.
   1. The Guppy's Collar split keeps an event field: `Amount`.
7. Drafts: 37 of 55 files kept the dumped setup and steps.
   1. Baby Haunt, Breakfast, Dinner: no step; an attack was added.
   2. Trinity Shield: Jawbone starts on `setup.stack`, `priority` 1 (owner answer 7).
   3. Random deck cards fixed in `setup.decks`, so the file names them.
   4. Blue Map, Compass, Map, Suicide King, Goat Head, Mom's Box.
   5. Also Cambion Conception, Cheese Grater, Spider Mod, The D10.
   6. A check before a later event: the step stops at a throw or question.
   7. Baby Haunt, Empty Vessel, both Bum-Bo! orders, Mom's Box answers.
8. Checks added from the card text (owner answer 8):
   1. Every attack: the monster died, its reward, the roll values.
   2. Baby Haunt: a 2 misses and a 3 hits the Fly.
   3. Cheese Grater: the top card is discarded; the next is on top.
   4. The D10: the Fly is back on top when the Dip dies (R-ZONE-11).
   5. Compass: Cain's loot step draws the new top card.
   6. Deaths: `died`, the penalty's item and loot discards.
   7. Guppy's Collar: `healed`, `turn_ended_early`, no `died`.
   8. Trinity Shield: Cain's character is used, the Sack is not.
9. Choices made here:
   1. Champion Belt: the turn's `attacks` is 0, then a second attack is allowed.
   2. It replaces the engine's `BonusAttacksUsed` count.
   3. Mom's Box answers: options in any order; the order is the display (24.6).
   4. The D10: no covered check; slot 1's card and the Fly's return say it.
   5. Monstro's Tooth: a random player; one Breakfast in the treasure discard.
10. Go tests that did not check what they claimed: 1.
   1. TestCheeseGrater counted 1 treasure discard, not the top card.
   2. The file names the discarded card and the next top.
11. Trinity Shield's refusal cites R-ABIL-12, a static ability; it fits.
12. No rules bug found; `cardgen -report` unchanged.

### [x] 24.14 Bulk: b2 other treasures

1. Goal: Active, Paid, Eternal, Soul Treasure (about 60 tests).
2. Counts (2026-11-18): 60 tests in 49 `_test.go` files of the family.
   1. Blank Card, Glass Cannon, Golden Razor Blade: no Go test left.
   2. Razor Blade, Smelter, The Battery, Tech X, The Shovel: the same.
3. Moved: 54 tests, as 61 files.
   1. TestBatteryBum, TestBloodLust, TestBloodLustOnAMonster.
   2. TestBookOfBelial: 5 files, one per case (owner answer 11).
   3. TestBookOfSin, TestBoomerang, TestBumFriend, TestCompost, TestChaos.
   4. TestContractFromBelow, TestDecoy, TestDiplopia, TestDonationMachine.
   5. TestFlush, TestForeverAloneSteals, TestForeverAloneDiscardsThenLoots.
   6. TestForeverAloneRechargesOnDamage, TestGodhead, TestGuppysHead.
   7. TestGuppysPaw, TestHostHat, TestIncubusSwaps, TestIncubusMayNotSwap.
   8. TestIncubusLootsAndPutsBack, TestJawbone, TestLazarusRags.
   9. TestLuckyFoot, TestMiniMush, TestMomsBra, TestModelingClay.
   10. TestMonsterManual, TestNo, TestMysterySack, TestMrBoom.
   11. TestPayToPlay, TestPortableSlotMachine, TestPotatoPeeler.
   12. TestPandorasBox, TestRemoteDetonator: 2 files each, one per case.
   13. TestSackHead, TestSackOfPennies, TestSleightOfHand, TestSpoonBender.
   14. TestTheChest, TestTheD100, TestTheD20, TestTheCurse.
   15. TestTheBoneAddsToARoll, TestTheBoneSoulCanBeDestroyed.
   16. TestTheD6, TestTheD6NeedsARoll, TestThePoop, TestTwoOfClubs.
   17. TestYumHeart: 2 files; see item 7.
4. Split: 3, as 4 files.
   1. TestCrystalBall: the loot 3 is a file.
   2. Its view dice and `Chosen` stay Go: TestCrystalBallShowsTheGuess.
   3. TestPlacebo: the use is a file.
   4. `AbilitiesOf` stays Go: TestPlaceboBorrowsNoEternalAbility.
   5. TestTheBone: 2 files; a counter set by hand starts the second.
   6. The soul's refusal stays Go: TestTheBoneSoulHasNoAbilities.
   7. It waits on 19.226: the engine cites R-ABIL-09.
5. Stayed Go, complicated (ADR 034 1.5): 3.
   1. Event fields (`Card`, `Private`): TestForeverAloneLooksAtADeck.
   2. Views and event fields: TestSackHeadTellsEveryoneTheChoice.
   3. Views, `Allowed`, bad intents: TestSleightOfHandAsksOneOrder.
6. Needed a format change: 0.
7. Go tests that did not check what they claimed: 2.
   1. TestYumHeart: "shields last past the end of the turn".
   2. The attack had used the shield up; no shield was left to end.
   3. `yum_heart/shield-ends-with-the-turn.json` checks it for real.
   4. TestPotatoPeeler counted 1 card per discard, not the top cards.
   5. The file names the three top cards in their discards.
8. Drafts: 31 of 65 files kept the dumped setup and steps.
   1. Rolls on the stack (11): `priority` 0 instead of Cain's pass.
   2. No!: the Jawbone steal starts on `setup.stack` (owner answer 7).
   3. Random deck cards fixed in `setup.decks` (15), so the file names them.
   4. Bum Friend, Compost, Crystal Ball, Flush, Incubus, Lazarus' Rags.
   5. Also Potato Peeler, Sack Head, Sleight of Hand, The Curse, The D20.
   6. Also Two of Clubs, Forever Alone, Portable Slot Machine, Sack of Pennies.
   7. A prompt check splits a step: Decoy, Diplopia, Donation Machine.
   8. Monster Manual: the refused end of turn is a step of its own.
   9. The Poop: a Bomb! on Isaac proves the shield.
   10. It replaces the Go check of `len(Shields)`.
   11. New setups: The Bone's second file, Yum Heart's second file.
9. Checks added from the card text (owner answer 8):
   1. Every ↷ use: the item is spent; roll items: `resolved_rolls`.
   2. Compost: A Dime!! back on top, the Bomb! under it.
   3. Contract From Below, Remote Detonator, The D20: the treasure discard.
   4. Flush: Dip and Pooter refill the slots.
   5. Lazarus' Rags: the penalty can't destroy the new Dinner.
   6. Pandora's Box: on 5 the box is in the discard; on 6 it is not.
   7. The Chest: Cain gains no treasure (R-MECH-48, R-MECH-03).
   8. The Curse: the milled card, the loot, the deck tops.
   9. Guppy's Paw, Mom's Bra, Host Hat: the monster dies, its reward.
   10. The Bone's soul: Judgement puts it in the treasure discard.
10. Refusals: rules from the rules text.
   1. Battery Bum, The Bone without counters: R-ABIL-07 fits.
   2. The D6 with no roll: R-ABIL-40 fits.
   3. Monster Manual's end of turn: R-ABIL-12 does not fit (19.225).
   4. Its Go test checked no rule; the file checks `allowed` only.
11. Bugs logged: 19.225 and 19.226; not fixed here.
12. Helpers `rollTable`, `resolvedRoll`, `handCards` moved to `common_test.go`.
   1. Other cards' tests use them; their files are gone.
13. `cardgen -report` unchanged.

### [x] 24.15 Bulk: b2 monsters

1. Goal: Basic and Cursed Monster (about 46 tests).
2. Counts (2026-11-23): 46 tests in 44 `_test.go` files of the family.
   1. 38 Basic, 6 Cursed Monster cards; Stoney has 3 tests.
   2. Dinga: no Go test left (24.7).
3. Moved: 46 tests, as 46 files; every `_test.go` of the family is gone.
   1. `testPlainMonster`, 14: Cod Worm, Clotty, Conjoined Fatty, Dip.
   2. Also Fatty, Fat Bat, Fly, Leech, Pale Fatty, Pooter.
   3. Also Red Host, Spider, Squirt, Trite: `<card>/dies-for-<reward>.json`.
   4. TestBigSpider, TestBlackBony, TestBoomFly, TestDankGlobin.
   5. TestDople, TestEvilTwin, TestGreedling, TestHanger, TestHopper.
   6. TestHorf, TestKeeperHead, TestLeaper, TestMomsDeadHand, TestMomsEye.
   7. TestMomsHand, TestMulligan, TestMulliboom, TestPortal, TestPsyHorf.
   8. TestRageCreep, TestRingOfFlies, TestSwarmOfFlies, TestWizoob.
   9. Cursed: Gaper, Keeper Head, Fatty, Psy Horf, Horf, Mom's Hand.
   10. TestStoney, TestStoneyCovered, TestStoneyFromTheDeck.
4. Split: 0. Stayed Go: 0. Needed a format change: 0.
5. Engine state read in Go, said by a game action in the file:
   1. Big Spider's `DeckAttacks`: the next attack offers only the deck.
   2. Portal's `MustAttacks`: `attacks` 1, an end of turn refused by R-ATK-03.
   3. Hanger's shop slot: a purchase offers the new slot's Breakfast.
   4. Mulligan's slot: two refill cards are both in play.
   5. Hopper's `AttackRolls`: `resolved_rolls` 6, 5, 5.
   6. Stoney covered: a deck attack puts the Pooter on it, then kills it.
   7. The Go test used `CoverMonsterSlot` and `KillObject`.
   8. Stoney from the deck: no roll, no `attack_target`; then an end of turn.
6. Drafts: 39 of 46 files kept the dumped steps.
   1. 37 got monster or treasure deck tops in `setup.decks`.
   2. So the refill and a gained treasure are named.
   3. Steps added: Big Spider, Hanger, Portal, Stoney from the deck.
   4. A prompt check splits a step: Mom's Dead Hand, Stoney.
   5. One-step Mystery Sack drafts became `action` (4).
7. Checks added from the card text (owner answer 8):
   1. Every kill: the monster died, its reward, the slot's refill.
   2. Plain monsters: no other reward kind, no soul.
   3. Black Bony: a reward roll of 6 loots 6; Boom Fly: 4¢.
   4. Dank Globin: both cards in the loot discard; Wizoob: Monstro too.
   5. Cursed: the roll's own result (4¢, loot 1) beside the curse.
   6. Mom's Hand: one attack roll; Mom's Dead Hand: the steal options.
   7. Dople, Evil Twin, Mulliboom: Cain's `died`; Isaac unhurt.
8. Go tests that did not check what they claimed: 1.
   1. TestMomsEye had no check (`_ = tb`).
   2. The file checks Isaac's `looked_at` and Cain's hand.
9. Bugs logged: 19.233; not fixed here.
   1. Mystery Sack's result resolves after Cursed Horf's trigger; no order asked.
   2. The Cursed Horf file leaves out the loot and its discard.
10. `cardgen -report` unchanged.

### [x] 24.16 Bulk: b2 events, curses, souls

1. Goal: Good and Bad Event, Curse, Bonus and Lost Soul (about 25 tests).
2. Counts (2026-11-23): 25 tests in 25 `_test.go` files of the family.
   1. 26 cards: 12 Good Event, 5 Bad Event, 5 Curse, 3 Bonus Soul, 1 Lost Soul.
   2. Chest (2nd): no Go test left (24.3); Soul of Wrath is in `r`.
3. Moved: 23 tests, as 23 files; their `_test.go` files are gone.
   1. TestAmbush, TestCursedChest, TestDarkChest, TestDarkChest2.
   2. TestDevilDeal, TestGoldChest, TestGoldChest2, TestSecretRoom.
   3. TestICanSeeForever, TestShopUpgrade, TestXLFloor, TestWeNeedToGoDeeper.
   4. TestMegaTrollBomb, TestTrollBombs, TestLostSoul.
   5. TestCurseOfAmnesia, TestCurseOfGreed, TestCurseOfLoss.
   6. TestCurseOfPain, TestCurseOfTheBlind.
   7. TestSoulOfGluttony, TestSoulOfGreed, TestSoulOfGuppy.
4. Split: 2, a file each.
   1. TestChest: the file checks the cents, the discard, the Fly uncovered.
   2. The discarded event's `Card` and `Prev` stay Go: TestChestDiscardedEvent.
   3. TestGreedEvent: the file checks the choice, the cents, `lost_cents`.
   4. The event's `Amount` stays Go: TestGreedEventLostCentsAmount.
5. Stayed Go whole: 0. Needed a format change: 0.
   1. The splits keep event fields (`Card`, `Prev`, `Amount`), as 24.13.
6. Engine state read in Go, said by a game action in the file:
   1. Ambush's `MustAttackDeck`: `attacks` 2, an end of turn refused by R-ATK-03.
   2. Then an attack offers only the monster deck.
   3. Shop Upgrade's shop slots: a purchase offers both new slots' items.
   4. Curses (`hasCurse`): no curse check (24.5); the curse's effect says it.
   5. Also `discards.monster` is `[]`: the curse did not go to discard.
   6. Soul of Guppy's condition: Cain is active with no Guppy item.
   7. So turn order asks Cain first; Isaac, with 2, gains the soul.
   8. Bonus souls: `setup.bonus_souls`; A Penny! is the action that checks them.
7. Drafts: 15 of 23 kept the dumped setup.
   1. Soul of Gluttony and Soul of Greed: a real game, no `Table`, no draft.
   2. Deck tops fixed in `setup.decks`, so the file names them (6).
   3. Curse of Amnesia, Gold Chest, I Can See Forever, Shop Upgrade.
   4. Also XL Floor, We Need To Go Deeper.
   5. Curse of Greed: Cain's 5¢, set between steps, is in `setup`.
   6. Soul of Guppy: no step; Cain active and his A Penny! added.
   7. One-step reveal drafts became `action` (13).
   8. Greed!: a prompt check splits the reveal from the choice.
8. Checks added from the card text (owner answer 8):
   1. Every event: it is in the monster discard; roll events: `resolved_rolls`.
   2. Shop Upgrade, XL Floor, We Need To Go Deeper: one more attack.
   3. XL Floor: the Pooter fills the new slot.
   4. Curse of Greed: no loss at the end of Isaac's turn.
   5. Curse of the Blind: DC 2 on Isaac's turn; Curse of Pain: no damage then.
   6. Curse of Loss: Monstro and the curse in the monster discard, Cain `died`.
   7. Dark Chest, Troll Bombs: `died` for Isaac only; Mega Troll Bomb: both.
   8. Greed!: Cain is the only choice; Isaac loses nothing.
   9. I Can See Forever: the 5 other cards on top, in the chosen order.
   10. Lost Soul: no item, not in the loot discard.
   11. Bonus souls: the other player gains none.
9. Go tests that did not check what they claimed: 1.
   1. TestSoulOfGluttony always skipped: its loop moved only Cain's few cards.
   2. The file gives Cain 10 loot cards in setup and checks the soul.
10. Refusal: Ambush's end of turn cites R-ATK-03; its text fits.
11. No rules bug found; `cardgen -report` unchanged.
12. Helpers `hasCurse`, `handCards` removed from `common_test.go`: no use left.

### [x] 24.17 Bulk: b2 bosses and characters

1. Goal: Boss, Epic Boss, Character (about 49 tests).
2. Counts (2026-11-23): 49 tests in 46 `_test.go` files of the family.
   1. 31 Boss, 4 Epic Boss, 11 Character cards.
   2. Bosses: one test each; Cain 2 tests, Eden 3.
3. Moved: 35 tests, as 35 files; their `_test.go` files are gone.
   1. Plain bosses: Gurdy, Little Horn, Monstro: `<card>/dies-for-<reward>.json`.
   2. TestCarrionQueen, TestChub, TestConquest, TestDaddyLongLegs, TestDarkOne.
   3. TestDeath, TestDelirium, TestEnvy, TestFamine, TestGemini, TestGluttony.
   4. TestGreed, TestGurdyJr, TestLarryJr, TestLust, TestMaskOfInfamy.
   5. TestMegaFatty, TestMom, TestPestilence, TestPeep, TestPin, TestPride.
   6. TestRagMan, TestSatan, TestScolex, TestSloth, TestTheBloat.
   7. TestTheHaunt, TestTheDukeOfFlies, TestTheLamb, TestWar, TestWrath.
4. Split: 11, a file each: every character's ↷ test.
   1. Blue Baby, Cain, Eden, Eve, Isaac, Judas, Lazarus.
   2. Also Lilith, Maggy, Samson, The Forgotten.
   3. File: `<character>/extra-loot-play-in-response.json`.
   4. The active one's ↷ starts on `setup.stack`; the other responds.
   5. It checks 2 loot plays, both deactivated, a second ↷ refused.
   6. Stays Go: the other player's `ExtraLootPlays` (`testExtraLootOnAnotherTurn`).
   7. No check says it (24.5 rejected `extra_loot_plays`).
   8. No step can use it: a play in response is mid-stack (owner answer 7).
   9. Stays Go: the starting item, from a real game start (`newGame`).
   10. Cain's and Eden's start order and Eden's choice stay Go too.
5. Stayed Go whole, complicated (ADR 034 1.5): 3.
   1. Game start, events, views, replay: TestCainGoesFirstWithoutARoll.
   2. Also TestEdenBeforeCain, TestEdensAtOnce.
6. Needed a format change: 0 whole tests.
   1. The splits keep an `extra_loot_plays` check and the game start.
7. Engine state read in Go, said by a game action in the file:
   1. `AttackRolls`: `resolved_rolls` (Carrion Queen, Chub, Mega Fatty, Pin, Duke).
   2. `MustAttacks` (Conquest, Envy): `attacks` 1, an end of turn refused by R-ATK-03.
   3. `Evasion`, `MonsterATK` (Larry Jr., Mask, Gemini): the attack's target question.
   4. The step stops there; then the rolls show the DC.
   5. Mom's `len(Monsters)`: two refill cards in slots, the third on the deck.
   6. Daddy Long Legs' `Boosts`: a 4 then misses; the Fly has DC 3.
   7. Famine: two `end_turn` steps; Cain is active again.
8. Drafts: 26 of 46 files kept the dumped steps.
   1. Monster or treasure deck tops in `setup.decks` (35), so cards are named.
   2. Characters: no draft kept; a response is `setup.stack` (owner answer 7).
   3. Gurdy Jr.: the roll of 1 starts as an `attack_roll` on `setup.stack`.
   4. A prompt check splits a step: Death, Satan!, Delirium.
   5. Also Gemini, Larry Jr., Mask of Infamy: no step in the draft.
   6. Satan! starts at 1 HP: one 6 kills it and a player.
   7. Conquest, Envy: a refused `end_turn` step added.
9. Setups made to prove the card text:
   1. Dark One, War: Isaac 3 HP survives the 2 damage; it was "2 or dead".
   2. Gluttony: Cain 4 HP; a 5 among the 6s deals him nothing.
   3. Daddy Long Legs: a Fly beside it, a 4 after the 1.
   4. Larry Jr.: a 3 misses, a 4 hits; Mask of Infamy: a 5 misses.
10. Checks added from the card text (owner answer 8):
   1. Every kill: the monster died, its reward, its soul; most: the refill.
   2. Delirium: the Fly back to DC 2; the deck's 5 top cards, Delirium 5th.
   3. Death, Satan!: both players offered; only Cain `died`.
   4. Pride: its trigger resolves before the first attack roll (R-ATK-25).
   5. Pride, Scolex, Sloth: the cards in the loot discard.
   6. Rag Man: no soul, 3 loot, the Pooter still on the deck.
   7. Gurdy Jr.: its trigger, then Mr. Boom, then the roll (R-STACK-02).
   8. Greed: Isaac 1 damage; The Bloat: Isaac 2; The Lamb: Cain has no soul.
11. Go tests that did not check what they claimed: 2.
   1. TestGurdyJr checked damage ≥ 1; the roll of 1 alone dealt that.
   2. The file checks 2 damage on 3 HP and the resolve order.
   3. TestDaddyLongLegs checked only that some boost existed.
   4. The file checks the 4 missing and the Fly's DC 3.
12. Runner note, no change here: a `monsters` check with only `card` checks nothing.
   1. `dead: false` also passes for a card in a deck.
   2. Mom's file checks `damage` 0, so the card must be in a slot.
13. No rules bug found; `cardgen -report` unchanged.
14. Helpers `testPlainMonster`, `killWithSixes` removed: no use left.
   1. `testExtraLoot` is now `testExtraLootOnAnotherTurn`.

### [x] 24.18 Bulk: b2 other

1. Goal: b2 test files with no card kind (about 57 tests).
2. Counts (2026-11-23): 57 tests looked at.
   1. 33 in files of no card: history 5, interactions 2, pride 3, three_plus 1, yuggy 22.
   2. 22 of card kinds no family had (24.12-24.17 left them).
   3. 1 Cent to 4 Cent, Pill/Rune, Dice Shard/Soul Heart, One-Use Treasure.
   4. Also Holy/Charmed Monster; Charged Baby's test sat in `holy_moms_eye_test.go`.
   5. 2 in `godhead_test.go`, reported by no sub-step.
   6. Not looked at again: Go tests reported "stays Go" before.
3. Moved: 47 tests, as 55 files; their Go tests are gone.
   1. Yuggy's rulings, 22: 17 in `interactions/`, 5 in card folders.
   2. Card folders: Chaos Card (eternal), Mom's Shovel, Guppy's Paw (dead).
   3. Also Soul of Gluttony (once per game), The D20 (a shop item).
   4. TestPrideDeckAttackTriggersBeforeTheRoll, TestPrideResolvesBeforeTheAttackRoll: `pride/`.
   5. TestAttackWithNoTargetTriggerRollsAtOnce: `fly/` (R-ATK-25 shown with the Fly).
   6. Coins: TestAPenny, TestTwoCents, TestThreeCents, TestFourCents.
   7. Runes and pills: TestBlankRune (2 files), TestEhwaz.
   8. TestPills, TestPills2, TestPills3: 3 files each (owner answer 11).
   9. TestDiceShard, TestSoulHeart, TestBox, TestChaosCard, TestTheD4, TestMomsShovel.
   10. Holy: Dinga, Dip, Keeper Head, Squirt, Mom's Eye (2 tests).
   11. TestChargedBabyOffersOnlySpentActiveItems: `charged_baby/`.
4. Split: 1. TestGodheadAsksForDiceValues.
   1. `godhead/offers-1-or-6.json`: the options 1 and 6, then a 1 resolves.
   2. The options' key (`ask.dice`) stays Go: no check reads a key.
5. Stayed Go, complicated (ADR 034 1.5, 1.6): 7, and 2 from 24.8.
   1. Views and stack item fields (`Text`, `TargetCard`): the 5 history tests.
   2. Real game start, decks counted: TestThreePlusCards.
   3. Engine API (`cardkit.DieFaces`): TestDieFacesAreDiceValues.
   4. The 2 of `interactions_test.go`: Go since 24.8, not judged again.
6. Needed a format change: 0.
7. Engine state read in Go, said by a game action in the file:
   1. Placebo's `AbilitiesOf`: Placebo copies the shop's Mystery Sack; 3 gains 4¢.
   2. The D20's shop slot: a purchase offers the treasure deck's Dinner.
   3. Guppy's Paw's `Shields`: no `shielded` event.
   4. Pride before the roll: a `throw` step stops at the first throw.
   5. Its step checks `resolved` Pride and no roll yet.
   6. The roll at once: no `passed_priority` between the pick and the throw.
8. Drafts: 34 of 55 files kept the dumped setup and steps.
   1. A response is `setup.stack` and `priority` 1 (owner answer 7): 9.
   2. Chaos Card (3), Sloth, Decoy, Glass Cannon, Godhead, Dice Shard.
   3. Curved Horn: the first attack roll starts on `setup.stack`.
   4. Deck tops in `setup.decks`, so cards are named: Ehwaz, Pills (2), The D4, The D20.
   5. Pride's deck attack also answers the slot (`ask.slot_card 1 fly`).
   6. No draft: Placebo (no step), Soul of Gluttony (a real game).
9. Checks added from the card text (owner answer 8):
   1. Every loot card in the loot discard; roll cards: `resolved_rolls`.
   2. Ehwaz: Pooter and Spider fill the slots; Fly and Leech in the discard.
   3. Dice Shard: Cain, the roll's controller, throws the reroll.
   4. Soul Heart: `prevented`, the Fly dies; Box!: each play allowed.
   5. Chaos Card, The D4, Mom's Shovel: the treasure discard; The D4 the new items.
   6. Mom's Shovel recharged on Isaac's next turn; the soul is Monstro.
   7. The Dead Cat: 8 counters left; Yum Heart: no `died`.
   8. Curved Horn: 2 damage on Big Spider, then the second 6 kills it.
   9. Chaos Card and XIII. Death: the treasure discard has both items.
   10. Sloth: no attack roll resolves; Isaac gains its soul and 1¢.
   11. Decoy: Isaac has the Decoy; Glass Cannon and The D20: no roll.
   12. Holy monsters: the Sack's own result beside the trigger.
10. Go tests that did not check what they claimed: 1.
   1. TestYuggyBonusSoulOncePerGame always skipped, as TestSoulOfGluttony (24.16).
   2. `soul_of_gluttony/gained-once-per-game.json` checks it for real.
   3. Isaac gains it, Chaos Card destroys it, a play leaves 10 cards: no soul.
11. Sources: Yuggy's Reddit comment IDs have no thread ID, so no link.
   1. The files cite the card text and the R-rules made from them.
   2. R-ABIL-34, R-ABIL-35, R-DICE-07, R-CARD-24, R-CARD-34, R-MECH-51.
12. Refusal: Modeling Clay's tapped copy cites R-ABIL-07; it fits (with R-ABIL-10).
13. No rules bug found; `cardgen -report` unchanged.
14. Helpers:
   1. `attackRoll`, `ownRoll` moved to `common_test.go`: other tests use them.
   2. Removed, no use left: `inDiscard`, `testGain`, `slayTable`, `bonusGame`.
   3. Also `runSome`, `rollTable`, `resolvedRoll`.
15. `open-questions.md`: Q-29's and S-YUG's test references point at the files.

### [x] 24.19 Bulk: fsp2

1. Goal: every fsp2 test (about 89 tests).
2. Counts (2026-11-23): 95 tests in 80 `_test.go` files of `cards/fsp2`.
   1. `common_test.go` holds helpers only.
   2. fsp2 gets its runner: `situations_test.go`, with its first files.
   3. Every file names `"sets": ["b2", "fsp2"]`: fsp2 plays with the Base Game.
3. Moved: 82 tests, as 91 files; 67 `_test.go` files are gone.
   1. Tables, one file per case (owner answer 11): Angel Room 3.
   2. Also Black Rune, Cursed Eye, Gimpy, Head of Krampus 2 each.
   3. Also Holy Chest, Pills! (3-4 / 5-6), Poker Chip 2 each.
   4. Every other card test: one file in its card folder.
   5. Poker Chip through Diplopia: `poker_chip/copied-by-diplopia.json`.
   6. No `interactions/`: no test is about two fsp2 cards.
   7. No cross-set file; `cards/testdata/situations/` stays unwired.
4. Split: 9, a file each; the rest stays Go.
   1. Test1Up: the heal's `Amount` stays Go (`Test1UpHealsOnce`).
   2. TestBossRushOwnSlot: `CheckInvariants` stays Go (engine API).
   3. TestMutantSpider: the dice event and the prompt text stay Go.
   4. TestRedCandle, TestTelepathyForDummies: view dice stay Go.
   5. Bum-Bo, Dark Judas, Guppy, Whore of Babylon: `<character>/extra-loot-play.json`.
   6. Their starting item stays Go: `testStartingItem` (a real game start).
5. Stayed Go whole: 4.
   1. Save and load: TestGimpySaveAtDiscard (ADR 034 1.5).
   2. Views: TestGuppysEye (1.5).
   3. Event field `Object`: TestCurseOfImpulseDeactivatedEvents (24.6).
   4. TestGetOutOfJailCard: needs a format change.
   5. Cain acts with priority over a stack, after the card resolved.
   6. A step can't stop there (owner answer 7); no setup field holds the lock.
6. Engine state or API in Go, said by a game action in the file:
   1. `DealDamageTo`: Bomb! in hand, or on `setup.stack` with `priority`.
   2. Abaddon, Gimpy, Globin, Isaac!, Lard, Smart Fly, Tumor.
   3. Dark Arts, The Left Hand: Cain starts with 1 damage; Bomb! kills him.
   4. Mama Haunt: Cain misses the Leech on his own turn and dies.
   5. `KillObject`: an attack (Swarmer, The Fallen, Mom's Heart!).
   6. Mom's Heart! starts with 7 damage; one 6 kills it.
   7. Headless Horseman: three hits, then a 3 misses for 1, two 4s kill.
   8. `DeckAttacks`: Blastocyst, Gold Key: the next attack offers only the deck.
   9. Gold Key: three deck attacks on events after the Fly, each allowed.
   10. `MustAttacks`: Boss Rush!, Krampus: an end of turn refused by R-ATK-03.
   11. `PlayerATK`, `Evasion`, `HP`: a prompt check stops at the target question.
   12. Nerve Ending, Widow, Polyphemus, The Wiz: then misses and hits show them.
   13. `Lingering` (Red Candle): Cain's 3 on his turn deals no damage.
   14. `FreeLoot`, `ExtraLootPlays` (Distant Admiration): later plays and a refusal.
   15. Responses to a waiting roll or loot: `setup.stack` (owner answer 7).
   16. 20/20, Libra, Euthanasia, Pills! (cancel).
7. Needed a format change: 1 whole test (5.4); the splits keep Go-only fields.
8. Drafts: most kept the dumped setup; each got `title`, `source`, `expect`.
   1. Deck tops in `setup.decks`, so refills, loot and treasure are named.
   2. `ForceRolls` after a direct call: all dice in `setup.rolls`.
   3. No draft (no step): Guppy's Eye, Nerve Ending, Widow, Polyphemus, ? Card refused.
9. Setups changed to prove the card text:
   1. Brain, Roundy: Breakfast for HP, not Lard; one fsp2 card per file.
   2. Flaming Hopper: no Lard; it took no damage.
   3. Bony keeps two Lards: 6 HP shows 1, then 3 damage; throw steps check each.
   4. Forget Me Now: Cain's Fatty soul is Widow, a soul card (R-CARD-18).
   5. Dark Arts: one file; a 6 gives 2¢, the Sack gives nothing on a 6.
   6. The Go test compared two tables; the file needs no second one.
10. Checks added from the card text (owner answer 8):
   1. Every kill: the monster died, its reward, its soul, the refill.
   2. Roll cards and attacks: `resolved_rolls`; used items: `deactivated`.
   3. Loot cards, destroyed items, events: their discard pile.
   4. Black Candle: all 5 other cards at the bottom; Cain's Curse of Pain hurts him.
   5. Curse of Blood Lust: Isaac may end the turn he used his attack.
   6. Cursed Eye on 1: no `damaged` for Isaac; Monstro still in its slot.
   7. Daddy Long Legs: Cain gets the Fly's 1¢.
   8. Divorce Papers: Cain keeps 4¢ and two cards; the card is destroyed.
   9. Hourglass: a `deactivated` event for Cain; Isaac still loots 1.
   10. Mom's Eye Shadow: the Fly untouched; Monstro II: Maggy unhurt.
   11. Mom's Heart!: Isaac's 2 souls; Poker Chip via Diplopia: one roll only.
   12. Rubber Cement: both slots refilled; Tape Worm: Isaac's own damage.
   13. Telepathy for Dummies: the start roll and the loot step's card.
11. Go tests that did not check what they claimed: 2.
   1. TestKrampus: Isaac ended his own turn first; that passed only by 19.235.
   2. The file starts in Isaac's end phase; Cain's end is refused.
   3. TestBlackCandle: 2 cards not in the top 5, not the bottom.
   4. The file names all 5 cards at the bottom.
   5. Weak, not wrong: 5 refusals accepted any rule.
   6. Blood Lust, both Krampus tests: the files check R-ATK-03.
12. Refusals: rules from the rules text.
   1. No loot play left: R-CARD-08 fits (Distant Admiration, 4 files).
   2. Second ↷ of a character: R-ABIL-07 fits.
   3. ? Card with no ↷ item cites R-ABIL-04; R-ABIL-40 fits (19.236).
   4. Its file checks `allowed: false` only.
13. Bugs logged: 19.235 and 19.236; not fixed here.
   1. 19.235: a setup with Krampus in a slot skips its duty that turn.
14. Every file failed with each expected value changed, then passed.
   1. An in-memory run changed each `expect` leaf: 1104 changes.
   2. Only swaps in any-order lists passed, as they should.
15. Helpers: `testCharacter` is now `testStartingItem`.
   1. Removed, no use left: `table`, `lootTable`, `inDiscard`, `foe`.
   2. Also `settleAny`, `top`, `topID`, `attackRolls`, `declareAttack`.
   3. Also `testStats`, `sackFirst`; the files spell the Sack's label.
16. `cardgen -report` unchanged.

### [x] 24.20 Bulk: g2

1. Goal: every g2 test (about 71 tests).
2. Counts (2026-11-23): 77 tests in 58 `_test.go` files of `cards/g2`.
   1. `common_test.go` holds helpers only.
   2. g2 gets its runner: `situations_test.go`, with its first files.
   3. Every file names `"sets": ["b2", "g2"]`: g2 plays with the Base Game.
3. Moved: 71 tests, as 82 files; 52 `_test.go` files are gone.
   1. Tables, one file per case (owner answer 11): I Am Error! 3.
   2. Also Crooked Penny, Gaper, Infamy, Joker, Trap Door! 2 each.
   3. Also Lord of the Pit, The Missing Page, Pills! (3 / 5) 2 each.
   4. TestHolyMantleOnYourTurn: 2 files; see item 7.
   5. Every other test: one file in its card folder.
   6. No `interactions/`: a second g2 card is only a prop.
   7. E.g. The Butter Bean rolls for Fruitcake and The Missing Page.
   8. No cross-set file; `cards/testdata/situations/` stays unwired.
4. Split: 4, the characters' `testCharacter`.
   1. Apollyon, Azazel, The Keeper, The Lost: `<character>/extra-loot-play.json`.
   2. Their starting item stays Go: `testStartingItem` (a real game start).
   3. The Lost also: `the_lost/recharges-a-character-at-end-of-turn.json`.
5. Stayed Go whole: 2.
   1. Save and load: TestPlanCSaveAtKill (ADR 034 1.5).
   2. TestCurseOfTinyHands: needs a format change.
   3. A step can't answer a hand-size discard, and `prompt` can't check it.
   4. Max hand size shows only in that discard (R-TURN-11).
6. Engine calls in Go, said by a game action in the file:
   1. `DealDamageTo` on a monster: Bomb! steps (Begotten, Fistula, Gurglings).
   2. The first check is at the target question, at full HP.
   3. `DealDamageTo` on a player: an attack that misses (Holy Card, Holy Mantle).
   4. Another player's turn: `phase: loot`, `priority` 1, as 24.10.
   5. Boil at 1 HP: `setup.monster_damage`.
   6. `MustAttacks`, `MustAttackDeck` (Gaper, I Am Error!): `attacks` 1.
   7. Then an end of turn refused by R-ATK-03; an attack offers only the deck.
   8. `Turn.Attacks` after a cancel: a second attack is allowed.
   9. `Evasion`, `MonsterATK`, `HP`, `PlayerATK`: a prompt check stops at the target question.
   10. `Lingering` (Lemon Mishap): Cain's 1 on his turn leaves Monstro unhurt.
   11. `len(Shop)` (Ouija Board, Void): a purchase offers the new shop items.
   12. Responses: `setup.stack` (owner answer 7).
   13. Lord of the Pit, The Bible: an `attack_roll`; I Can't Believe…: Lil Battery.
   14. Fruitcake (Q-10): the three activations, Fruitcake on top.
7. Setups changed to prove the card text:
   1. Holy Mantle's "next death only": two Bombs on the Lost on Isaac's turn.
   2. The Go test's second damage came after its turn had ended.
   3. Breakfast for 3 HP: Imp (was Placenta, a g2 card), Knight, Polycephalus, Pills! on 5.
   4. Knight: throw steps check each 1; the order of the two damages is no rule.
   5. Credit Card after a purchase: 23¢, so Breakfast can be paid.
8. Checks added from the card text (owner answer 8):
   1. Every kill: the monster died, its reward, its soul, the refill.
   2. Roll cards and attacks: `resolved_rolls`; used items: `deactivated`.
   3. Loot cards and events: their discard pile.
   4. Charger, Hush, The Cage: stats at the target question, then a kill.
   5. Fistula: DC back to 2 on Cain's turn; Death's Head: not offered as a target.
   6. Curse of Fatigue: the options; The Lost: hers, with `hp` 1, `souls` 1.
   7. Ouija Board: the other 3 at the bottom, 2 attacks; Trap Door!: the deck order.
   8. Steven on 2: the roll resolves, Steven takes no damage; on 6: no kill.
   9. Joker: `looked_at`; Plan C: the soul, 6¢, then both `died`.
9. Go tests that did not check what they claimed: 3.
   1. TestCreditCardAfterAPurchase: 3¢ left either way.
   2. With 3¢ Breakfast can't be bought at 10¢, nor would a free one cost.
   3. TestTheLost: Cain's own recharge step recharged him anyway.
   4. The file recharges Isaac, the third seat, before his turn.
   5. TestPolycephalus: "dies at 3" was wrong; the 4th roll is a 1 (R-ATK-19).
   6. Isaac died; the check only counted 4 rolls.
   7. Weak, not wrong: I Can't Believe…'s refusal accepted any rule.
   8. The file checks R-ABIL-40; Holy Mantle's recharge was a `t.Log`.
10. Refusals: rules from the rules text.
   1. I Can't Believe… on another player's item: R-ABIL-40 fits.
   2. Second ↷ of a character: R-ABIL-07 fits.
   3. End of turn with an attack owed: R-ATK-03 fits.
11. Choice made here: Lord of the Pit with no attack is allowed.
   1. "Cancel an attack" names no target, so R-ABIL-40 does not apply.
   2. It does nothing: 1 attack left.
   3. Reversed by 19.237 (owner, 2026-11-23): refused by R-ABIL-40.
12. Every file failed with each expected value changed, then passed.
   1. An in-memory run changed each `expect` leaf: 813 changes.
13. No rules bug found.
14. Helpers removed, no use left: `table`, `lootTable`, `inDiscard`, `foe`.
   1. Also `declareAttack`, `settleAny`, `fight`, `topID`, `engineRef`.
   2. Also `attackIntent`, `attackRolls`, `testStats`, `beanFirst`.
   3. `testCharacter` is now `testStartingItem`; the files spell the Bean's label.
15. Old step 19 entries name removed Go tests; left as history.
16. `cardgen -report` unchanged.

### [x] 24.21 Bulk: r rooms

1. Goal: Room (about 59 tests).
2. Counts (2026-11-23): 59 tests in 57 `_test.go` files of the family.
   1. 61 Room cards; Beggar, Bomb Bum, Red Champions, Red Vise: files since 24.5-24.6.
   2. Spider Webs and Tax For The Mighty have 2 tests each.
3. Moved: 58 tests, as 59 files; 56 `_test.go` files are gone.
   1. TestAllHallowsEve, TestAngelicIntervention, TestBattleRoyale, 2, 3.
   2. TestBlackChampions, TestBlessingOfGluttony, Greed, Steam, The Inner Eye, The Sack.
   3. TestBlindRage, TestBloodDonation, TestBloodLustR, TestBloodMoney, TestBumBoIsLoose.
   4. TestButterFingers, TestChallengeRoom, TestConjoinedTwin, TestDevilBeggar.
   5. TestDiceRoom1 to 6, TestEdensBlessingR, TestEquality, TestEternalChest.
   6. TestFloorSpikes, TestFortuneTeller, TestGreedLooms, TestGus, TestHaunted.
   7. TestHeavyIsTheHead, TestIAmErrorR, TestIsaacsBlessing, TestKeyMaster, TestLaserEye.
   8. TestMaggysBlessing, TestMightForTheMeek, TestPityForThePoor, TestPlanetarium.
   9. TestRestockMachine, TestSamsonsBlessing, TestShadowOfFamine, TestShadowOfWar.
   10. TestShellGame, TestSlotMachine, TestSocialGoals, TestSpiderWebs, TestSplashDamage.
   11. TestStaticShock, TestSpoilsOfWar, TestTheMirror, TestWhiteChampions.
   12. TestTaxForTheMighty: 2 files, one per case (owner answer 11).
   13. TestTaxForTheMightyForcedAttack: `tax_for_the_mighty/forced-attack-unpaid-may-end-the-turn.json`.
   14. Social Goals: in its own folder; Auction Gavel is a prop (as 24.20).
   15. An `interactions/` file would drop it from cardgen's "tested".
4. Split: 0.
5. Stayed Go: 1. TestSpiderWebsDeactivatedEvent: the event's `Object` (24.6, as 24.19).
   1. The Spider Webs file checks the event's kind and player only.
6. Needed a format change: 0.
7. Engine state or API in Go, said by a game action in the file:
   1. `CanGainSouls` (All Hallows' Eve): Isaac kills Chub, Cain Little Horn; no soul.
   2. `PlayerAttackDC` (Battle Royale 2, 3): a roll one below misses, the DC hits.
   3. `MonsterATK`, `Evasion`, `PlayerHP`, `PlayerATK`: a prompt check, or a turn's end.
   4. `MustAttackDeck` (Blind Rage): end refused by R-ATK-03; the attack offers only the deck.
   5. `MustAttacks` (Blood Lust, Challenge Room): `attacks` 1, end refused by R-ATK-03.
   6. `DealDamageTo` (Blood Money, Shadow Of War, Spoils Of War): a miss at 1 HP, or Bomb!.
   7. `RefillSlots` (Dice Rooms): the "enters play" trigger on `setup.stack`.
   8. `SkipTurns` (Haunted): Cain ends his turn and is active again.
   9. `GiveCents` (Social Goals): Cain bids 6¢ on Isaac's Auction Gavel.
   10. `AddNamedCounters` (Social Goals): `setup.counters`, as the draft had them.
   11. `ForceAttacks` (Tax, forced): `must_attacks` 1, `must_attack_deck` 1 (19.235).
8. Drafts: 27 of 59 files kept the dumped setup and steps.
   1. The `setup` helper's 6 Clotty and 6 Book Of Sin: only the named tops kept.
   2. No step in the draft (7): a step added; stats at the attack target question.
   3. Dice Rooms: a mid-stack pass; now a `trigger` item on `setup.stack`.
   4. Refused steps added: Blind Rage, Blood Lust, Challenge Room, Might For The Meek.
   5. Setups changed to prove the card text: see item 10.
9. Checks added from the card text (owner answer 8):
   1. Roll rooms: `resolved_rolls`; costs: cents, damage, the loot discard.
   2. Rooms that discard themselves: the room discard; the slot refills (R-SHOP-06).
   3. Kills: the monster died, its reward; Gus: the soul, Cain's helper roll.
   4. Dice Room!: at the turn's end it goes to discard, White Champions refills.
   5. Dice Rooms: the discards, the new items and hands, the refilled slots.
   6. Inner Eye: the deck order after the Pooter; the Fly back in slot 1.
   7. Might For The Meek: a third attack refused by R-TURN-07.
   8. Shadow Of War: the room's trigger resolves before the death penalty.
   9. Equality: a third player gains 1¢ too; Splash Damage: Samson unhurt.
   10. Social Goals: the goal counters, every player's 2 treasures.
   11. Tax: the `gave_cents` event; Isaac's 1¢ from the Fly.
10. Setups changed to prove the card text:
   1. Floor Spikes: the Fly and Breakfast; Isaac takes exactly 2 and lives.
   2. Gus: 2 HP left; it dies only from Cain's helper roll.
   3. Heavy Is The Head: Jawbone for Glitch, which has no ↷ (R-MECH-54).
   4. Shadow Of Famine: the second ability is Sack Of Pennies, not the room.
11. Go tests that did not check what they claimed: 2.
   1. TestFloorSpikes: damage ≥ 2, but the Leech's 2 ATK killed Isaac alone.
   2. TestShadowOfFamine: Shadow Of Famine has no activated ability to refuse.
   3. Weak, not wrong: Blind Rage, Shadow Of Famine, Tax refusals accepted any rule.
   4. The files check R-ATK-03 and R-ABIL-12.
12. Refusals: rules from the rules text.
   1. Shadow Of Famine, Tax For The Mighty: R-ABIL-12, a static ability; it fits.
   2. Blind Rage, Blood Lust, Challenge Room: R-ATK-03 fits.
   3. Might For The Meek's third attack: R-TURN-07 fits.
13. Note, no change: `attacks` leaves out a static bonus attack (as 24.13).
   1. Might For The Meek: `attacks` 0 after the first attack; the second is allowed.
14. Every file failed with each expected value changed, then passed.
   1. An in-memory run changed each `expect` leaf: 607 changes.
   2. Only changes that say another true thing passed (14).
   3. E.g. a player index to another player with the same values.
15. No rules bug found.
16. Helpers removed, no use left: `fightRoom`, `enterRoom`.
17. `cardgen -report` unchanged.

### [ ] 24.22 Bulk: r passive and eternal treasures

1. Goal: Passive and Eternal Treasure (about 55 tests).

### [ ] 24.23 Bulk: r other treasures

1. Goal: Active, Paid, Soul Treasure (about 35 tests).

### [ ] 24.24 Bulk: r monsters, bosses, events, souls

1. Goal: monsters, bosses, events, Bonus Soul, Curse (about 49 tests).

### [ ] 24.25 Bulk: r loot and characters

1. Goal: loot cards, Outside, Character (about 50 tests).

### [ ] 24.26 Bulk: r other

1. Goal: r test files with no card kind (about 55 tests).

### [ ] 24.27 Bulk: challenge

1. Goal: Challenge and Challenge Boss (about 50 tests).

### [x] 24.28 Turn checks count every attack and loot play left

1. Why (owner, 2026-11-24)
   1. `expect.attacks` reads only `Turn.Attacks`.
   2. A bonus from a card always in effect is left out (24.13, 24.21 item 13).
   3. Might For The Meek: `attacks` 0 after the first attack; a second is allowed.
   4. Owner: the check counts what the player can still do.
2. Change
   1. `attacks`: the attacks the active player has left, bonus ones too.
   2. The engine's own count: `attacksLeft` (`combat.go`).
   3. `loot_plays`: the same, by `lootPlaysFor` (`activate.go`).
   4. It then counts `extra_loot_plays` and static bonus loot plays.
   5. Same mechanism, same fix (CLAUDE.md 5.26).
   6. Setup fields keep their meaning: base counts, bonus ones on top.
   7. No rules change; the game already counts them so.
3. Files
   1. Every file whose `attacks` or `loot_plays` check changes is fixed.
   2. The new value comes from the card text, not from the run (5.25).
   3. Might For The Meek, Champion Belt and similar: check the bonus at once.
   4. List each file changed and why in the report.
4. Tests and docs
   1. Go: a pass and a fail for a static bonus attack and loot play.
   2. Wiki `Situation-payload.md`: the two checks, and the dump if it changes.
   3. ADR 034: drift edited in place, if any.
5. Done (2026-11-24)
   1. `checkGame`: `attacks` reads `attacksLeft`, `loot_plays` `lootPlaysFor(active)`.
   2. No engine game code changed; no version bump.
   3. Only one file failed with the new counts: Champion Belt, attacks twice.
6. Files changed, values from the card text:
   1. Champion Belt, attacks twice: `attacks` 1 after the first attack, 0 after the second.
   2. Champion Belt, a second time: `attacks` 0 after the bonus attack.
   3. Might For The Meek: `attacks` 1 after the Fly; 0 after the Leech kept.
   4. Skeleton Key: `attacks` 98 after one attack, 97 after two ("up to 99").
   5. Polydactyly: `loot_plays` 1 after the first A Penny!, 0 after the second.
   6. Belly Button: `loot_plays` 2; no loot played; R-TURN-05 added to `source`.
   7. Polydactyly: R-TURN-05 added to `source`.
7. Other checks searched (5.26):
   1. Player `hp`, `atk`, `souls`; monster `hp`, `atk`, `evasion`: engine counts already.
   2. Expect has no other turn count; `after` has none.
   3. The 74 files with an `attacks` or `loot_plays` check were read.
   4. `AddLootPlays` on one's own turn adds to `Turn.LootPlays`; those files keep their values.
   5. Dump: writes setup only, base counts; no change.
   6. Go card tests read `Turn.Attacks`, `Turn.LootPlays` (about 20 files).
   7. None has a bonus card; left as is.
8. Note, no change: a stacked `attack_roll` in setup uses a base attack only.
   1. With `setup.attacks` 0 it uses no bonus attack; the engine would.
   2. No file does that today.
9. Tests: `TestSituationExpectBonusCounts`, a pass and two fails.
   1. Test cards `belt` and `finger`, like Champion Belt and Polydactyly.
10. Wiki `Situation-payload.md`: the two checks and the setup meaning.
11. ADR 034: no drift; it does not say what the counts read.
12. No rules bug found.
