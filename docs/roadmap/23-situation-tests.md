# Step 23. Situation tests

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

## File layout

1. A card's tests: `cards/<set>/testdata/situations/<card>/<case>.json`.
   1. `<card>` is the card ref, as its `<card>_test.go`.
   2. `<case>` says what is checked, e.g. `pays-5-cents.json`.
2. Two or more cards of one set together: `.../situations/interactions/`.
   1. Named after the cards, e.g. `lucky-foot-bomb.json`.
3. Cards of several sets: `cards/testdata/situations/interactions/`.
4. Engine rules, test sets only: `testdata/situations/<area>/`.
   1. `<area>` is the rule ID area in lower case: `atk`, `stack`, `abil`.
   2. The 6 files today move there (23.2).
5. A resolved disputed case (R-03) goes where its card or rule goes.
   1. Its `source` names the ruling.
6. A folder holds only `.json` files; no deeper levels.

### [x] 23.1 ADR: situation files as card tests

1. Goal: one ADR fixes when a test is a file and when Go.
2. Tasks:
   1. Written in the main conversation (5.21.2).
   2. The owner answers and the file layout above.
   3. The `source` field and its forms.
   4. "Simple" and "complicated", with examples of each.
3. Done when:
   1. The ADR is Accepted and in `INDEX.md`.

### [x] 23.2 A runner per set, and the `source` field

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

### [x] 23.3 Setup: dice, decks, counters, turn

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
   4. `.Hand =` 6: `hand`; mid-test changes wait for `steps` (23.4).
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

### [x] 23.4 Several steps

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

### [x] 23.5 Richer checks

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

### [x] 23.6 Checks on prompts and order

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

### [ ] 23.7 Dump a Go test's setup as a file

1. Goal: a tool writes `setup` and `steps` from a Go test run.
2. Tasks:
   1. `enginetest.Table` records its setup, rolls, decks and steps.
   2. An env var writes them to a file per test.
   3. It never writes `expect` (rule 3).
   4. Hand-set state it can't express is listed as a TODO.
3. Done when:
   1. A dumped b2 test, with `expect` from its Go checks, passes.

### [ ] 23.8 Pilot: move one group

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

### [ ] 23.9 Docs

1. Goal: a simple new card test starts as a file.
2. Tasks:
   1. `Adding-a-card.md`: simple is a file, complicated is Go.
   2. `Situation-payload.md`: every new field and the file layout.
   3. `CLAUDE.md` Commands: how to run one set's files.
3. Done when:
   1. A simple new card from the guide needs no Go test.

### [x] 23.10 Priority windows

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
   5. With a stack, `priority` works as now (23.4).
   6. Rewrite the Cain file into a real window, e.g. the loot step.
   7. A file per window, expected values from the rules.
   8. Wiki: `Situation-payload.md`; ADR 033 if it drifts.
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
