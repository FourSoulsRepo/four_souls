# 34. Situation files as card tests

* **Status:** Accepted
* **Date:** 2026-11-17
* **Authors:** @HardDie
* **Renumbered:** from 033 at the merge with bots (owner, 2026-11-18).

---

## Context

1. The situation payload exists since 4.11 (LR-04, LR-07).
   1. The sandbox and the rules website read it.
   2. `pkg/rules_engine/testdata/situations/` holds 6 test files.
2. Card tests are Go: about 930 tests in `pkg/rules_engine/cards/`.
   1. They build tables from `engine.SituationSetup` already.
   2. About 380 are one action and a plain state check.
3. New tests stayed in Go because the format is too narrow.
   1. One action only; no forced dice, deck tops or counters.
   2. Checks are only cents, hand size, damage, souls, dead.
   3. No runner in the set packages; no guide in `Adding-a-card.md`.
4. A test's expected value comes from the rules (CLAUDE.md 5.25).
   1. Go tests name the source in a comment.
   2. JSON has no comments.
5. Owner answers, 2026-11-17:
   1. Move a test only when a file says exactly the same.
   2. Go and JSON side by side for one card is fine.
   3. A simple test is a file; a complicated one is Go.
   4. Files are grouped by purpose, never one big folder.
   5. No website examples from test files for now.

## Considered options

1. **Simple tests as files, grouped by purpose; Go for the rest**
   1. Won. Readable without Go; the format already exists.
   2. Go keeps what a fixed format can't say.
2. **Every card test as a file**
   1. Lost. The format would grow into a programming language.
3. **Keep all tests in Go**
   1. Lost. The format stays a toy; tests stay unreadable to the owner.
4. **Copy tests to files, keep the Go ones too**
   1. Lost (owner). Two copies drift apart.
5. **All files in one folder**
   1. Lost (owner). The purpose of a test is lost.

## Decision

Use option 1.

1. Simple and complicated
   1. Simple: the format states its setup, steps and checks.
   2. No engine state set by hand outside `setup`.
   3. Complicated, so Go: loops that compute values.
   4. Go: a check against a value read earlier in the test.
   5. Go: views, hidden cards, records, replay, undo.
   6. Go: helpers and engine APIs, not a game action.
   7. Go: random runs, fuzz, sims.
   8. A missing setup field or check is added to the format first.
   9. Go only when adding it would not serve other tests.
   10. A table of fixed cases: one file per case (owner, 2026-11-18).
2. Example of simple: `TestRazorBlade`.
   1. Razor Blade in play, activate it on the other player.
   2. That player has 1 damage.
3. Example of complicated: `TestRulingDamageToDeadPlayer`.
   1. It compares damage with a value read before the play.
4. Moving
   1. A test moves only when its file checks exactly the same.
   2. Its Go test is deleted in the same change.
   3. A Go test with simple and complicated checks may split.
   4. Each part then keeps the checks it can say.
5. File layout, relative to `pkg/rules_engine/`:
   1. One card: `cards/<set>/testdata/situations/<card>/<case>.json`.
   2. Several cards of one set: `cards/<set>/testdata/situations/interactions/`.
   3. Cards of several sets: `cards/testdata/situations/interactions/`.
   4. Engine rules, test sets only: `testdata/situations/<topic>/`.
6. Names
   1. `<card>`: the card ref, as in `<card>_test.go`.
   2. `<topic>`: the rule topic in lower case, as `R-ATK` is `atk`.
   3. `<case>`: what is checked, in kebab case: `pays-5-cents.json`.
   4. An interaction: the cards' refs, e.g. `lucky-foot-bomb.json`.
   5. A folder holds only `.json` files; no deeper levels.
   6. The runner fails on a file outside this layout.
7. A resolved disputed case (R-03) goes where its card or rule goes.
8. The `source` field
   1. A test file needs `source`: a non-empty list of strings.
   2. A sandbox or website payload needs none.
   3. Each item names one basis for the expected values.
9. `source` forms
   1. `card <ref>`: the card's own text.
   2. `R-<TOPIC>-<NN>`: our rule text, `pkg/rules_engine/docs/rules/`.
   3. `Q-<NN>`: a resolved question in `open-questions.md`.
   4. `R-07 #<N>`: an old ruling re-checked in `open-questions.md`.
   5. An `https://` link: a ruling outside the repo.
   6. The runner checks the form; a rule or question ID must exist.
   7. House and Community rulings keep their label in `open-questions.md`.
   8. An open question may be cited only by a file whose title says it waits.
   9. Such a file checks the engine's current answer (owner, 2026-11-25: Q-27, Q-28).
10. Format compatibility
   1. Old payloads keep working; website links hold them (12.7).
   2. `action` stays; `steps` is added beside it (24.4).
   3. Every new field is optional for a payload.
   4. A file has `action` or `steps`, never both.
   5. A step may carry its own `expect`; the top one is checked at the end.
   6. `expect` checks are named as in `setup`; one left out is not checked (24.5).
   7. A card list in `expect` names every card, in any order.
   8. A deck in `expect.decks` compares only its top cards, top first.
   9. `prompt` checks the open question: who is asked, every option (24.6).
   10. A prompt check stops its step at the question it checks.
   11. `resolved` names stack items in order; others may come between.
   12. `events` names event kinds that must, or must not, happen.
   13. A room works as an item in `activate` steps and `setup.stack`.
   14. `setup.phase` and `setup.declared` start in a real priority window (24.10).
   15. Phases: `start`, `loot`, `action` (default), `end`.
   16. `declared`: an attack or purchase, before its target or item.
   17. Everyone passes: the game moves on as in a game (R-STACK-07).
   18. Open action phase: `priority` on an empty stack is the active player's (R-TURN-09).
   19. `declared` may come with `setup.stack`: responses resolve before the choice.
   20. `expect.deck_bottoms` names a deck's bottom cards, in any order (24.11).
   21. Any order: several cards go to the bottom at random (Q-22).
   22. `expect.resolved_rolls`: each roll's final value, in order (R-DICE-06).
   23. A card in a slot sets no duty; the setup names it (owner, 2026-11-23).
   24. `setup.must_attack`, `must_attacks`, `must_attack_deck`: this turn's duties (19.235).
   25. `setup.must_attack_one_of`: an extra attack owed on one of these cards (19.241).
   26. It names cards; a covered one counts once it comes up (19.241).
   27. `setup.random`: the next random picks, by option label (24.30).
   28. Picks after the list are random from `seed`, as dice after `rolls`.
   29. `setup.game`: the game type; `regular` by default (24.34).
   30. `challenge` starts one on the built table, as a game does.
   31. `team` is refused until Q-03 is resolved.
   32. `"card": "player"` in counters: counters on the player (R-MECH-13).
   33. `expect.presents`: a challenge's presents.
11. `expect` comes from the rules, never from a run (5.25).
   1. A tool may write `setup` and `steps`; never `expect` (24.7).
   2. `enginetest` dumps a Go test with `SITUATION_DUMP` (24.7).
   3. The dump adds `title`, `sets`, an empty `source`, and `todo`.
   4. `todo` lists what it can't say; the runner rejects it until removed.
   5. Hand-set state before a first direct `G.Apply` is read from the save origin (24.11).
12. Runner
   1. `enginetest.RunSituations` runs a set's tree in its `go test`.
   2. One subtest per file, named by its path in the tree.
   3. `enginetest.RunRuleSituations` runs the engine's topic tree.
   4. Rule and question IDs come from package `docs` (`pkg/rules_engine/docs`).
   5. Only tests import it; the wasm build stays the same.
   6. It reads a test file strictly: an unknown field fails (24.3).
   7. A payload outside a tree is read as before.
   8. An unknown field's error names its line, path and nearest field.
   9. A mismatch names its step, or the end for the top `expect`.
   10. `enginetest.RunInteractionSituations` runs the cross-set tree with every set (24.22).
   11. A test file fails on rolls or picks left unused (24.32).
13. A new card's simple tests are files from now on.
   1. `Adding-a-card.md` says how (24.9).

## Consequences

### Positive

1. The owner reads a test without Go.
2. A folder shows what its tests are for.
3. Every expected value names its source in the file.
4. Tests, the sandbox and the website share one format.

### Negative and risks

1. Two test forms to learn for card work.
   1. The rule above picks one per test.
2. The format grows with each new kind of check.
   1. A field is added only when several tests need it.
3. A refactor of the engine can't rename refs in JSON files.
   1. The runner fails on an unknown ref, with the options.
4. The move costs time; the pilot (24.8) measures it first.

### Neutral

1. The 6 root files move into topic folders (24.2).
2. Website examples from test files may come later; not now.
