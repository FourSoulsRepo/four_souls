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
7. Every item gets a group in "Bugs by system" below.
   1. The group is where the cause was, not where it showed.
   2. Add the item to its row and update the count.
   3. A new kind of cause gets a new row at the end.

## Bugs by system

| ID | Group                                    | System     | Items                                                  | Count |
|----|------------------------------------------|------------|--------------------------------------------------------|-------|
| A  | Rules: a card or rule coded wrong        | Engine     | 19.22, 19.23, 19.25, 19.27, 19.36                      | 5     |
| B  | Events and view: data never sent         | Engine     | 19.2, 19.9, 19.11, 19.39                               | 4     |
| C  | Card data: wrong values in `card_db`     | Card data  | 19.24                                                  | 1     |
| D  | View data sent, never drawn              | Client     | 19.5, 19.29, 19.30                                     | 3     |
| E  | Playback, animations, history            | Client     | 19.1, 19.3, 19.4, 19.6, 19.13, 19.41                   | 6     |
| F  | Layout and visuals                       | Client     | 19.15, 19.26, 19.40                                    | 3     |
| G  | Prompts: the choice is unclear           | Client     | 19.16, 19.17, 19.19, 19.42                             | 4     |
| H  | Lobby and rejoin                         | Server     | 19.32                                                  | 1     |
| I  | e2e test logic: the app was right        | e2e        | 19.7, 19.8, 19.10, 19.18, 19.21, 19.28, 19.31, 19.34, 19.37 | 9 |
| J  | e2e under load, or cause not found       | e2e        | 19.12, 19.14, 19.20, 19.33, 19.35, 19.38, 19.43, 19.44 | 8     |

Totals by system (44 items):

| System    | Groups  | Count |
|-----------|---------|-------|
| Engine    | A, B    | 9     |
| Card data | C       | 1     |
| Client    | D–G     | 16    |
| Server    | H       | 1     |
| e2e       | I, J    | 17    |

---

### [x] 19.1 Roll results in the history

1. Goal: the history shows the roll result the game used.
2. Report (owner, 2026-10-14)
   1. Play test; record `2026-10-14-173137-g1_Gaijin_Oleg.fsrec`.
   2. Polycephalus: "the next attack roll this turn has -2".
   3. The history did not show the -2; it looked ignored.
3. Tasks:
   1. Find the cause in the match record.
   2. Show the die and the result when they differ.
   3. The strip's die shows the result.
4. Done when:
   1. Vitest: a changed roll shows both numbers; an unchanged one, one.
   2. e2e test: a changed attack roll in the strip and the panel.
   3. Use case UC-15 updated.
5. Cause
   1. Engine is right: Polycephalus's trigger adds a one-shot boost.
   2. The next attack roll uses it up (`useNextAttackRoll`).
   3. Record: die 6, `roll_resolved` 4; the 4 hit.
   4. The history used `stack_resolved`: the die before changes.
   5. `roll_resolved`, the result, was never shown.
   6. Same for every roll change, e.g. a +1: "rolled 5" for a 6.
6. Result (2026-10-14)
   1. A changed roll reads "Roll — attack roll (rolled 6 → 4)".
   2. The strip's die and line show the result: "rolls 4".
   3. An unchanged roll reads as before: "(rolled 6)".
   4. Engine, protocol and records unchanged.
   5. Screenshots: `.cache/screenshots/19.1/before/` and `after/`.
7. Owner notes (2026-10-14)
   1. Owner marks 19.1 solved.

### [x] 19.2 Heal popup after 1-Up

1. Goal: the heal popup shows the HP a player got back.
2. Report (owner, 2026-10-14)
   1. Play test; record `2026-10-14-173137-g1_Gaijin_Oleg.fsrec`.
   2. Lazarus died on his turn with 1-Up.
   3. The popup said +1 HP; his HP went up by 2.
3. Tasks:
   1. Find the cause in the match record.
   2. Every HP a death prevention gives back gets a `healed` event.
   3. 1-Up shows one popup with the full amount.
4. Done when:
   1. Engine test: 1-Up heals in one event for the full amount.
   2. Engine test: Guppy's Collar emits a `healed` event for 1 HP.
   3. Old records still replay.
5. Cause
   1. Step 196: Leech's combat damage leaves Lazarus at 0 of 2 HP.
   2. Step 198: 1-Up resolves; `healed` amount 1.
   3. `PreventYourDeath` sets HP to 1 with no `healed` event.
   4. 1-Up's heal to full then emits only the rest: +1.
   5. Broken Ankh and Guppy's Collar show no heal at all.
6. Result (2026-10-14)
   1. `PreventYourDeath` emits `healed` for the HP it gives back.
   2. New block `PreventYourDeathHealFull`: one `healed` to full.
   3. 1-Up uses it: popup "+2 HP" for Lazarus at 0 of 2 HP.
   4. Broken Ankh and Guppy's Collar now show "+1 HP".
   5. Events are not in the checksum; the owner's record replays.
   6. No frontend change; the popup is short, so no screenshots.
7. Owner notes (2026-10-14)
   1. Owner marks 19.2 solved.

### [x] 19.3 Deck attack: the monster lands before the roll

1. Goal: on a deck attack the monster is in its slot before the dice roll.
2. Report (owner, 2026-10-14)
   1. A player attacks the top of the monster deck.
   2. The dice roll plays first.
   3. Only then the monster appears in the chosen slot.
3. Tasks:
   1. Find the cause in the code.
   2. A card put into a slot flies there and stays until the board redraws.
   3. Effects after it, like the attack roll, play over the placed card.
4. Done when:
   1. Vitest: a slot reveal gives a place effect before the dice.
   2. e2e test: the placed card is on the slot while the dice plays.
   3. Use case updated.
5. Cause
   1. Engine order is right: `placeRevealed` puts the card in the slot.
   2. `startAttack` then rolls; both land in one update.
   3. Playback plays an update's effects, then draws its board.
   4. So the dice plays before the board shows the monster in the slot.
6. Result (2026-10-14)
   1. New `place` effect: the card flies from its deck to the slot.
   2. It stays there, drawn by `Effects`, until the board redraws.
   3. Order now: reveal, the card lands in its slot, the attack roll.
   4. Shop refills get the same fly into their slot.
   5. Engine, protocol and records unchanged.
   6. Screenshots: `.cache/screenshots/19.3/before/` and `after/`.
7. Owner notes (2026-10-14)
   1. Owner marks 19.3 solved.

### [x] 19.4 Deck attack shows the monster large twice

1. Goal: a deck attack shows the revealed monster large once.
2. Report (owner, 2026-10-14)
   1. Found while fixing 19.3.
   2. The monster shows large when it is drawn.
   3. It shows large again when it goes into the slot.
3. Tasks:
   1. Find the cause in the code.
   2. A card shown large in the update before is not shown again.
   3. It still flies into its slot (19.3).
4. Done when:
   1. Vitest: a card revealed in the update before gets no second reveal.
   2. e2e test: one large reveal for a deck attack.
   3. Use case updated.
5. Cause
   1. `attackDeck` emits `card_revealed` for the drawn card.
   2. `placeRevealed` puts it in the slot; `putInSlot` emits it again.
   3. Playback shows every `card_revealed` large.
6. Result (2026-10-14)
   1. A card the update before revealed gets no second large reveal.
   2. A deck attack: one large reveal, then the card flies to its slot.
   3. Other reveals and refills show large as before.
   4. Engine, protocol and records unchanged.
   5. Screenshots: `.cache/screenshots/19.4/before/` and `after/`.
7. Owner notes (2026-10-14)
   1. Owner marks 19.4 solved.

### [x] 19.5 Revealed hands show as backs

1. Goal: a revealed opponent's hand shows its cards face up.
2. Report (owner, 2026-10-14)
   1. A player controls Guppy's Eye.
   2. "Each other player plays with their hand revealed."
   3. The opponents' hands still show only card backs.
3. Tasks:
   1. Find the cause in the code.
   2. An opponent's hand the view sends shows face up, fanned.
   3. Each card zooms on hover, like other table cards.
   4. A hidden hand stays backs with its count.
4. Done when:
   1. e2e test: a revealed hand shows its cards at the top and sides.
   2. e2e test: the hand zooms a card on hover.
   3. Use case UC-12 updated.
5. Cause
   1. Engine is right: `View` sends the hand under `StatRevealHand`.
   2. `OpponentArea` always draws `HiddenHand` from `hand_size`.
   3. It never reads `player.hand`; no screen ever drew it.
6. Result (2026-10-14)
   1. A hand the view sends shows face up: `RevealedHand`.
   2. Same fan, place and count as the backs; every card shows.
   3. At the top the cards hang past the edge, turned to their owner.
   4. Hover a card to read it in full.
   5. Engine, protocol and records unchanged.
   6. Screenshots: `.cache/screenshots/19.5/before/` and `after/`.
7. Owner notes (2026-10-15)
   1. Owner marks 19.5 solved.

### [x] 19.6 Deck attack: the monster lands in its slot twice

1. Goal: a card that flies into its place does not appear there again.
2. Report (owner, 2026-10-15)
   1. A player attacks the top of the monster deck.
   2. The monster flies to its slot; the dice roll plays.
   3. Then the same monster appears in the slot again.
   4. Not the first animation clash: 19.3 and 19.4 were the same kind.
   5. Owner asks for a rule and a check against animation clashes.
3. Tasks:
   1. Find the cause in the code.
   2. One card gets one arrival animation per update.
   3. List every animation in one place, with what it owns.
   4. A check that finds two animations on one card.
   5. A rule for new animations; research how others solve it.
4. Done when:
   1. Vitest: the clash check finds a double arrival and a double reveal.
   2. e2e test: a deck attack's monster does not pop in after it lands.
   3. Every e2e test fails on an animation clash.
   4. ADR and wiki page for the rule; use case updated.
5. Cause
   1. The client has two sources of animation.
   2. Events: the `place` effect flies the monster to its slot (19.3).
   3. The view: `Card` pops in whenever a new `objectId` mounts.
   4. The slot's card gets a new ID, so it popped in after it landed.
   5. Looted and discarded cards did the same after their fly.
   6. Nothing listed the animations or checked two on one card.
6. Result (2026-10-15)
   1. Effects claim the card they bring or show: `animations.ts`.
   2. `Card` does not pop in a card an effect brought that update.
   3. A clash is logged and skipped; e2e tests fail on any.
   4. Every animation has an `anim: <id>` tag.
   5. `docs/wiki/Animations.md` lists them; Vitest checks the list.
   6. ADR 012 holds the rule and the research.
   7. Engine, protocol and records unchanged.
   8. Screenshots: `.cache/screenshots/19.6/before/` and `after/`.
7. Owner notes (2026-10-15)
   1. Owner marks 19.6 solved.
   2. Claims stay on cards; zones get them with a second source.

### [x] 19.7 Buy is lit, then the purchase fails

1. Goal: e2e `settle` answers only a state its done check saw.
2. Report (agent, 2026-10-15)
   1. Found in `make e2e` while fixing 19.6; chromium.
   2. Test: "buy and attack icons are lit only when the engine allows them".
   3. `actions.spec.ts:149`: "the game never settled".
   4. Tia had Buy lit and pressed B at step 10.
   5. The engine answered `purchase_failed`, amount 10.
   6. No choice dialog opened; the test waited until it timed out.
   7. Passed twice when rerun alone: the deal may matter.
3. Tasks:
   1. Find the cause: engine, client, or the test.
   2. Find which dealt card makes Buy allowed but unaffordable.
4. Done when:
   1. Cause found and written here.
   2. A test that forces the race, fails before the fix, passes after.
5. Cause
   1. The test helper `settle`, not the engine or the client.
   2. Buy needs no cents: a failed payment is a rule (R-SHOP-04).
   3. Players start with 3 cents; every price is 10.
   4. So every purchase on turn 1 fails; no card in the deal matters.
   5. `purchase_failed` comes only after the purchase choice is answered.
   6. So the dialog opened, and `settle` picked its first option.
   7. Each `settle` loop runs `done()`, then `answer()`.
   8. `done()` saw no dialog: Tia's update had not come yet.
   9. It came while `answer()` read the pages, one by one.
   10. `answer()` found the dialog and answered it.
   11. The purchase failed; no dialog again; `done()` never held.
   12. Timing decides it: reruns pass.
6. Result (2026-10-15)
   1. `settle` keeps each page's step from before `done()`.
   2. `answer()` makes no move when a page's step changed since.
   3. The next loop checks `done()` against the new state first.
   4. e2e test: updates come 150 ms late; `done()` is slow.
   5. It failed 2 of 2 before the fix and passes 3 of 3 after.
   6. Engine, client, protocol and records unchanged.
7. Owner notes (2026-10-15)
   1. Owner marks 19.7 solved.

### [x] 19.8 Intent waits forever after a pass

1. Goal: a sent intent always ends: by an update or a refusal.
2. Report (agent, 2026-10-15)
   1. Found in `make e2e` while fixing 19.6; webkit.
   2. Test: "a resolved loot card goes to the history strip (15.22)".
   3. `info.spec.ts:99`: "the game never settled".
   4. Both pages were at step 21, turn 2, Ada's turn.
   5. Ada's page still had `sent: 19`; the other page `sent: 20`.
   6. Last event: `passed_priority` by player 1.
   7. Ada showed "Your turn" with no buttons.
   8. Passed 4 of 4 fresh runs alone: timing may matter.
3. Tasks:
   1. Find the cause: client `sent`, server, or the test.
   2. Find why `sent` stays set after later updates arrive.
4. Done when:
   1. Cause found and written here.
   2. A test that reproduces it.
5. Cause
   1. The test, not the client or the server.
   2. `sent` was not stuck: Ada had `sent: 19` at step 21.
   3. The client blocks only while `sent === game.step`.
   4. Turn 1 was Ben's; the test kept his page as `active`.
   5. Ben's loot card ended his turn early.
   6. Such cards: The Fool, `EndYourTurn` cards, damage that kills him.
   7. Turn 2 was Ada's; Ben's page showed "Their turn".
   8. `settle` waited for `canEndTurn(active)` on Ben's page: never true.
   9. `answer()` never ends a turn, so nothing moved on.
   10. The dump lists only the `Actions` region's buttons.
   11. End turn sits in the `Board` region, so it looked like no buttons.
   12. Rare: only deals where that loot card is the first playable one.
6. Fix
   1. The wait also ends when the turn moved on or the game ended.
   2. Same check as `backToMe` in `actions.spec.ts`.
   3. No test reproduces it: that needs a rigged hand.
   4. The owner chose no rig hook for e2e (2026-10-15).
7. Owner notes (2026-10-15)
   1. Owner marks 19.8 solved.

### [x] 19.9 Deck attack: the drawn card hides during the slot choice

1. Goal: the drawn card stays in view until the player places it.
2. Report (owner, 2026-10-15)
   1. A player attacks the top of the monster deck.
   2. The new card shows for a few seconds, then hides.
   3. The player then picks the slot to put it in.
   4. A player who missed the reveal does not know the card.
   5. The slot choice may depend on the card.
3. Tasks:
   1. Find the cause in the code.
   2. Show the drawn card while the slot choice is open.
   3. It hides once the card is placed.
4. Done when:
   1. e2e test: the drawn card is visible during the slot choice.
   2. Use case updated.
5. Cause
   1. Engine order is right: `attackDeck` reveals, then asks for the slot.
   2. The card waits outside the table in `Attack.Revealed`.
   3. `View` never sent it; only the `card_revealed` event named it.
   4. The client showed that event large for a moment, then nothing.
   5. The slot dialog lists only "slot 1", "slot 2".
6. Result (2026-10-15)
   1. The view sends the waiting card as `revealed`, to every viewer.
   2. The board shows it face up on the monster deck.
   3. The slot dialog shows it above the slot buttons.
   4. It leaves the deck when it flies to the chosen slot.
   5. New view field: protocol stays version 1 (ADR 006).
   6. Records unchanged; they hold intents, not views.
   7. Screenshots: `.cache/screenshots/19.9/before/` and `after/`.
7. Owner notes (2026-10-15)
   1. Owner marks 19.9 solved.

### [x] 19.10 e2e `settle` cannot answer a target aimed from a card

1. Goal: `settle` answers every prompt a player can answer.
2. Report (agent, 2026-10-15)
   1. Found in a loop of the 19.8 test; webkit, run 54 of 80.
   2. Same error: "the game never settled".
   3. Ada played Dagaz: choose one, then choose a player.
   4. Ada showed "Choose a target"; no dialog, empty stack.
   5. Her loot was not played yet: still in hand, "Loot 1".
3. Tasks:
   1. Find the cause.
   2. Teach `settle` to answer that prompt.
4. Done when:
   1. Cause found and written here.
   2. The 19.8 test passes in a long loop.
5. Cause
   1. Since 15.34 a target prompt of one's own card aims from it.
   2. The dialog hides; the targets glow `hl-aim` on the table.
   3. `answer()` picks table targets only while a `.choice` dialog shows.
   4. Then Pass is disabled and Attack does not apply: no move.
   5. `playLoot` aims before the play, so most cards never hit this.
   6. A modal card (Dagaz) asks its mode first, its target after.
   7. That target prompt comes as a fresh aim: nobody clicks it.
6. Fix
   1. `answer()` clicks the first uncovered `hl-aim` target.
   2. It runs when no dialog shows, before Pass.
7. Owner notes (2026-10-15)
   1. Owner marks 19.10 solved.

### [x] 19.11 Look at the top card: only its name shows

1. Goal: the player sees the card they decide about.
2. Report (owner, 2026-10-15)
   1. Golden Horseshoe and Cain's Eye trigger at turn start.
   2. The screen shows only the card's name in the buttons.
   3. The card itself never shows while the player decides.
3. Tasks:
   1. Find the cause in the code.
   2. Show the card in the dialog while the choice is open.
   3. Only the chooser sees it; the deck top stays hidden to others.
4. Done when:
   1. Engine test: the chooser's view has the card; others' do not.
   2. e2e test: the card shows in the dialog.
   3. Use case updated.
5. Cause
   1. Both use `cardkit.LookMayBottom`: one question, then `LookAt`.
   2. The card is only text in the options: "keep X on top".
   3. `LookAt` runs after the answer; the look dialog comes too late.
   4. The view has no field for the card a question is about.
   5. Same for Sack Head and fsp2 `lookMayDiscard`.
   6. Options also showed the raw ref: "keep a_dime on top".
6. Result (2026-10-15)
   1. New `Question.Subject`: the card a question is about.
   2. The view sends it as `subject`, only to the chooser.
   3. The choice dialog shows it above the buttons.
   4. Buttons name it by title: "keep A Dime!! on top".
   5. Set for `LookMayBottom`, Sack Head, `lookMayDiscard`.
   6. Also for fsp2 curse and boss reveals ("Give X to which player?").
   7. New view field: protocol stays version 1 (ADR 006).
   8. Screenshots: `.cache/screenshots/19.11/before/` and `after/`.
7. Owner notes (2026-10-15)
   1. Owner marks 19.11 solved.

### [ ] 19.12 A whole random game hangs: Escape never returns

1. Goal: a whole random game never freezes a page.
2. Report (agent, 2026-10-15)
   1. New e2e `games.spec.ts`: random moves, every set, chromium.
   2. 3 of 12 games hung until the 60-minute test timeout.
   3. Two 4-player games, one 2-player game.
   4. Each hung in `page.keyboard.press('Escape')`.
   5. That runs after a move found nothing to press.
   6. One hung after turn 30; no progress line came after it.
   7. The 60-second stall check never fired: the loop was blocked.
   8. A key press waits for the page to take input.
   9. So the page likely stopped running its event loop.
   10. The error snapshot shows only an old finished page.
   11. Seeds and match records of the three were lost.
   12. The spec logged a seed only on a win.
   13. The next run's start deleted the server's records folder.
3. Tasks:
   1. Make the spec record a hang: page state, screenshot, last events.
   2. Bound every wait in the loop, so a hang fails in seconds.
   3. Rerun until a hang repeats; find the cause.
   4. Fix it in the client or the server.
4. Done when:
   1. Cause found and written here.
   2. `make e2e-games` passes five rounds with no hang.
5. Findings so far
   1. A rerun at 11:10 lost three games: every page reloaded.
   2. The owner's own Wails build rewrote `frontend/wailsjs`.
   3. The e2e Vite server watched it and reloaded every page.
   4. Not the cause of the first hangs: no build ran then.
   5. e2e Vite now watches no files (`FOUR_SOULS_E2E`).
   6. The spec now fails a hang in 30 s and keeps its evidence.
   7. Evidence: seed, record, screenshots, store, a frozen page's stack.
   8. Cause still unknown; waiting for a hang to repeat.
   9. More than 20 games since, with the diagnostics on: no hang.
   10. Paused by the owner (2026-10-15); see `docs/wiki/Frontend-tests.md`.

### [x] 19.13 The look dialog shows one card twice

1. Goal: the look dialog shows each card it is about once.
2. Report (agent, 2026-10-15)
   1. Found by the new `games.spec.ts`; 2 players, seed 371908.
   2. Record `2026-10-15-071253-g1_Rnd2r0Ann_Rnd2r0Bob.fsrec`.
   3. Bob's page logged a React error: two children with key 361.
   4. 361 is a Library Card Bob looked at twice.
   5. Step 2314: Forever Alone, "look at the top card of a deck".
   6. Step 2316: More Options looks at the top 2 treasures: 361, 430.
3. Tasks:
   1. Find the cause.
   2. Show each looked-at card once.
4. Done when:
   1. Cause found and written here.
   2. Vitest: the same card looked at twice in one update shows once.
5. Cause
   1. The server views have no duplicate: a replay checked each list.
   2. Steps 2315 and 2316 are server auto-passes (`after`, N-05).
   3. Their events join step 2314's in one update.
   4. That update has `looked_at` 361 twice.
   5. `lookedAt` keeps every event, so 361 is listed twice.
   6. `LookDialog` keys its cards by object id: React warns.
6. Result (2026-10-15)
   1. `lookedAt` keeps the first look at each object in one update.
   2. The history still logs both looks; each was a real look.
   3. Engine, server, protocol and records unchanged.
   4. Vitest in `history.test.ts`.
   5. No screenshots: the state needs that one seed's deal.
7. Owner notes (2026-10-15)
   1. Owner marks 19.13 solved.
   2. The animation inventory could not catch it: no animation ran.
   3. Idea TS-06: every e2e test fails on a console error.

### [ ] 19.14 e2e: a guest's Join click can miss under load

1. Goal: `startGame` joins every guest, however busy the lobby is.
2. Report (agent, 2026-10-15)
   1. `make e2e`, webkit: N-16 test in `stack.spec.ts` failed.
   2. Oti clicked Join on "Wan's game · 1/2 players".
   3. The page stayed on the list; no error; Join still enabled.
   4. About 20 other test games were starting in the same lobby.
   5. Passed 3 of 3 separate reruns.
   6. Evidence: `.cache/e2e-fails/20261015-120743-gate/` (not in git).
   7. Next full run: the same miss in `table.spec.ts` (15.24), webkit.
   8. Evidence: `.cache/e2e-fails/20261015-121014-gate2/`.
   9. Full runs with e2e Vite watching no files: 2 of 4 failed.
   10. Full runs with watching on: 2 of 2 passed.
   11. Too few runs to blame the no-watch setting; check it first.
3. Tasks:
   1. Find why the click did not join: a re-rendered row or a lost reply.
   2. Make `startGame` join reliably, or fix the lobby if it drops it.
4. Done when:
   1. Cause found and written here.
   2. `make e2e` passes several runs in a row.
5. Findings so far
   1. Gate run for 19.13: the same miss in `info.spec.ts` (15.47), webkit.
   2. Again with e2e Vite watching no files.
   3. The rerun passed; its run overwrote the evidence.

### [x] 19.15 A full hand: the end cards are hard to play

1. Goal: every card of a big hand can be hovered and played.
2. Report (owner, 2026-10-15)
   1. A hand of 12 cards.
   2. The cards at the sides are almost impossible to play.
3. Tasks:
   1. Find the cause in the fan layout.
   2. Keep every card's top in reach, however big the hand.
4. Done when:
   1. e2e test: a 12-card hand; each card shows a top to point at.
   2. Screenshots before and after.
5. Cause
   1. `fan` dropped each card by `spread² × 4px`, with no limit.
   2. 12 cards: the end spread is ±5.5, so a drop of about 121 px.
   3. The hand box is 104 px high; the table clips below it.
   4. So the end cards sank off screen.
   5. Their raised neighbours covered what was left.
   6. e2e measure: the first card had 0 px² the pointer could hit.
6. Result (2026-10-15)
   1. The end cards drop at most 24 px (`fanDrop`), any hand size.
   2. Hands of up to 5 cards keep their old shape.
   3. Every card of a 12-card hand shows its top and name.
   4. Screenshots: `.cache/screenshots/19.15/before/` and `after/`.
7. Owner notes (2026-10-15)
   1. Owner marks 19.15 solved.

### [x] 19.16 Put back in any order: the final order is unclear

1. Goal: the player sees and sets where each card goes.
2. Report (owner, 2026-10-15)
   1. Look at the top 5 cards; put them back in any order.
   2. The screens are confusing.
   3. In the end the player does not know which card is where.
3. Tasks:
   1. Find the cause in the code.
   2. Show the whole order while the player builds it.
   3. After the choice, show the cards in their new order.
4. Done when:
   1. Cause found and written here.
   2. e2e test: the dialog shows each card's place in the deck.
   3. Screenshots before and after.
5. Cause
   1. `PutBackInOrder` asks one question per card.
   2. Each has the same text: "Which card goes next, from the top?".
   3. The dialog shows only the cards left, never the placed ones.
   4. No position number: card 3 of 5 looks like card 1.
   5. The last card is asked too, with one option.
   6. Then `LookAt` shows the cards in the old order, not the new one.
   7. Same code in Sleight of Hand and fsp2 `reorderADeck`.
6. Result (2026-10-15)
   1. Owner's pick: arrange all cards, then confirm.
   2. New engine prompt `PromptOrder`, answered by `IntentOrder`.
   3. The intent holds every card, top first.
   4. An `Ask` question with `Order` opens it; `do` reads `c.Order`.
   5. `PutBackInOrder`, Sleight of Hand and `reorderADeck` use it.
   6. The view sends the cards as `ordering`, only to the chooser.
   7. `OrderDialog`: numbered places, "1 · top" first; cards below.
   8. A click places a card; a click on a placed card takes it back.
   9. "Put back" sends the order once every card has a place.
   10. No "You look at" dialog after it.
   11. Protocol 2: an older client cannot answer the new prompt.
   12. Card tests unchanged: test drivers answer an order by card refs.
   13. Screenshots: `.cache/screenshots/19.16/before/` and `after/`.
7. Owner notes (2026-10-15)
   1. Owner marks 19.16 solved.

### [x] 19.17 Two "instead" effects: the order choice is unclear

1. Goal: the player sees what each order of effects gives.
2. Report (owner, 2026-10-15)
   1. Bum-Bo! and Counterfeit Penny in play; the player gains ¢.
   2. The game asks which effect goes first.
   3. The question is confusing.
   4. The owner wants Bum-Bo! to level up by more, and cannot tell how.
3. Tasks:
   1. Find the cause in the code.
   2. Name the card of each option.
   3. Say what each order does.
4. Done when:
   1. Cause found and written here.
   2. e2e test: the dialog names both cards and each result.
   3. Screenshots before and after.
5. Cause
   1. Both cards are replacements ("instead"), not triggers.
   2. `processQueue` asks `ChooseReplacement` (R-ABIL-33).
   3. Title: "Choose which happens first".
   4. Each option is only the card's rule text; no card name.
   5. Nothing says what each order does.
   6. Order matters: the first one rewrites the gain.
   7. Penny first: gain 3¢ → 4¢, then Bum-Bo! levels up by 4.
   8. Bum-Bo! first: it levels up by 3; no ¢ left for the Penny.
   9. The player must work this out from two card texts.
   10. Poker Chip (fsp2) is the third such card; it rolls a die.
6. Result (2026-10-15)
   1. Owner's pick: show what each order gives.
   2. The engine plays each order on a copy of the game (`previewReplacement`).
   3. The choice keeps the results as `Previews`.
   4. The view sends them as each option's `outcome`.
   5. Each option is now the card itself: image and name.
   6. Title: "Gaining 10¢: which card applies first?".
   7. Under each card: "→ Bum-Bo! levels up by 11".
   8. A card that then does nothing: "→ Counterfeit Penny does nothing".
   9. Three or more cards: "Then you pick the next card".
   10. Poker Chip: "→ Poker Chip rolls for 4¢"; the roll is unknown.
   11. No protocol bump: older clients ignore the new field.
   12. Screenshots: `.cache/screenshots/19.17/before/` and `after/`.
7. Owner notes (2026-10-15)
   1. Owner marks 19.17 solved.
   2. Trigger order on the stack is a separate bug: 19.19.

### [ ] 19.18 e2e: a played Dagaz waits for a player target

1. Goal: the stack test plays any loot card the deal gives.
2. Report (agent, 2026-10-15)
   1. Found by `make e2e` during 19.17.
   2. `stack.spec.ts`: "a played loot card shows on the stack".
   3. Chromium; the hand had Dagaz.
   4. The stack stayed empty; the status said "You: Choose a target".
   5. Seen again in 15.28, Chromium and WebKit; the panel said "Choose a target".
3. Tasks:
   1. Answer a target picked on the table, not only in a dialog.
4. Done when:
   1. The test passes with Dagaz in the hand.
5. Cause
   1. Dagaz asks its mode, then a player to protect.
   2. A player target is picked on the table (15.25), not in a dialog.
   3. The test's loop clicks only `.choice-option` buttons.
   4. So it never answers the target; the card never reaches the stack.
   5. Not caused by 19.17; the engine and client behave as designed.
   6. Evidence kept outside the repo before the rerun.
   7. Seen again in 15.69, WebKit: "every player sees whom the stack waits for".
   8. Same loop in that test; the panel said "Choose a target".
   9. Evidence: `.cache/e2e-evidence-15.69/`.
   10. Seen again in 15.35, WebKit, the same test, both tries.
   11. Records: one played Dagaz, the other Lil Battery (an item target).
   12. Evidence: `.cache/e2e-evidence-15.35/`.

### [x] 19.19 Trigger order: "on the stack first" is unclear

1. Goal: the player knows which trigger resolves first.
2. Report (owner, 2026-10-15)
   1. Several triggers fire at once; the player orders them.
   2. The title: "Choose which goes on the stack first".
   3. It does not say that the first one in resolves last.
3. Tasks:
   1. Say in the title that the first one resolves last.
   2. Show each option as its card, with name and ability text.
4. Done when:
   1. e2e test: the dialog says "resolves last" and names each card.
   2. Screenshots before and after.
5. Cause
   1. `placeNextTrigger` asks `ChooseTriggerOrder` (R-ABIL-15).
   2. The client title: "Choose which goes on the stack first".
   3. The stack resolves last in, first out; the title hides that.
   4. Each option label is `abilityText`: "card_ref: text".
   5. The client shows the raw ref, not the card's name.
   6. No preview as in 19.17: a trigger resolves later.
   7. Responses, rolls and new triggers may come between.
6. Result (2026-10-15)
   1. Title: "Choose which goes on the stack first; it resolves last".
   2. Each option shows its card, the card's name, and the ability text.
   3. The client splits the label "card_ref: text" into card and text.
   4. Engine, protocol and records unchanged.
   5. Vitest and e2e test in `actions.spec.ts`.
   6. Screenshots: `.cache/screenshots/19.19/before/` and `after/`.
7. Owner notes (2026-10-15)
   1. Owner marks 19.19 solved.

### [ ] 19.20 e2e: a cancelled Dagaz shows twice in the hand

1. Goal: the aim test passes every run.
2. Report (agent, 2026-10-15)
   1. Found by `make screenshots` during 15.55.
   2. `aim.spec.ts`: "a card placed first aims after its mode …".
   3. Chromium; after the right-click cancel.
   4. The hand had two "Dagaz" buttons; the test expects one.
   5. Fails about 1 run in 5, also without the 15.55 change.
   6. Evidence: `.cache/e2e-evidence/15.55-aim-*`.
   7. Seen again in `make e2e` during 15.64: `.cache/e2e-evidence/15.64-aim/`.
3. Tasks:
   1. Find why the hand shows Dagaz twice after the cancel.
   2. Fix the client or the test, whichever is wrong.
4. Done when:
   1. The test passes ten runs in a row.

### [x] 19.21 e2e: Eden's starting choice hides the soul crystals

1. Goal: the 15.56 soul test passes on every deal.
2. Report (agent, 2026-10-15)
   1. Found by `make e2e` during 15.6.
   2. `table.spec.ts`: "soul cards hide behind the crystals …".
   3. Chromium; hovering the second own crystal timed out.
   4. "Choose your starting item" stood over the own area.
   5. Evidence: `.cache/e2e-evidence/15.6-souls/`.
3. Cause
   1. The deal gave the viewer Eden (R-SETUP-09).
   2. Eden's choice opens at the bottom edge, as agreed in 15.16.
   3. So the dialog is right; the test expects a plain board.
4. Tasks:
   1. Ban Eden in that test, with the reason (CLAUDE.md 5.16).
5. Done when:
   1. The test passes with the ban.
6. Result (2026-10-15)
   1. The test bans Eden; a comment says why.
   2. Eden is the only character with a starting choice.
   3. Three separate runs pass.
   4. No screen change; no screenshots.

### [x] 19.22 The D10 puts an event in a slot; it never resolves

1. Goal: an event put in a monster slot by a card resolves at once.
2. Report (owner, 2026-10-19)
   1. Play test; record `2026-10-18-075512-g4_Hitokiri-Battosai_KING_jеня.fsrec`.
   2. A monster died; its slot was empty.
   3. Another player's D10 put the monster deck's top card there.
   4. It was Gold Chest; it never resolved, up to the end of the game.
3. Cause
   1. Record step 1188: D10 resolves; only `card_revealed` gold_chest.
   2. No `entered_play` event; no roll; the chest stays in the slot.
   3. `CoverMonsterSlot` (zones.go) never calls `enterMonsterSlot`.
   4. `PlaceFromDeck` and `placeRevealed` do call it.
   5. Same bug for every `coverASlot` card in g2 (Dad's Key).
4. Tasks:
   1. `CoverMonsterSlot` runs `enterMonsterSlot` on the new card.
5. Done when:
   1. Go test: the D10 puts Gold Chest in a slot; it rolls and pays.
   2. After it, the chest is in the discard; the covered card is back.
6. Result (2026-10-19)
   1. `CoverMonsterSlot` runs `enterMonsterSlot`, like `PlaceFromDeck`.
   2. An event rolls at once, goes to discard; the covered card returns.
   3. A curse goes to a player, as from a deck attack.
   4. Test `TestTheD10Event`: fails before the fix, passes after.
   5. Old records replay apart from such a placement on.
   6. Engine only; no screen change, no screenshots.
7. Owner notes (2026-10-19)
   1. Owner marks 19.22 solved.

### [x] 19.23 Deck attack reveals Stoney; the fight goes on

1. Goal: a revealed monster that can't be attacked cancels the attack.
2. Report (owner, 2026-10-19)
   1. Play test; the game is not remembered.
   2. A deck attack revealed Stoney ("This can't be attacked").
   3. The attack went on: attack rolls against Stoney.
3. Cause
   1. Record `2026-10-15-190346-g1_Oleg_jеня.fsrec`, step 9.
   2. `card_revealed` stoney, then `attack_target` stoney.
   3. `placeRevealed` (combat.go) starts the attack on any monster.
   4. Only the target list checks `Unattackable`.
   5. Rule R-ATK-10 was never coded.
   6. Same for Death's Head (g2).
4. Tasks:
   1. `placeRevealed`: an `Unattackable` card stays; the attack ends.
5. Done when:
   1. Go test: a deck attack reveals Stoney; no roll; the attack is over.
   2. Stoney stays in the slot; the turn goes on.
6. Result (2026-10-19)
   1. `placeRevealed`: an `Unattackable` card stays in the slot.
   2. The attack ends; the active player gets priority back.
   3. Covers Stoney (b2) and Death's Head (g2).
   4. Test `TestStoneyFromTheDeck`: fails before the fix, passes after.
   5. Old records replay apart from such a reveal on.
   6. UC-13 gets alternative 5e.
   7. Engine only; no screen change, no screenshots.
7. Owner notes (2026-10-19)
   1. Owner marks 19.23 solved.

### [x] 19.24 Wrong rewards: Ring of Flies, The Lamb, Mom's Dead Hand

1. Goal: the three monsters pay what their printed cards show.
2. Report (owner, 2026-10-19)
   1. Ring of Flies should give 3¢.
   2. The Lamb should give 3¢ and 1 soul.
   3. Mom's Dead Hand should give 4¢.
3. Cause
   1. The card images show those rewards.
   2. foursouls.com reward tables add extra icons.
   3. Ring of Flies: an extra 1x Loot.
   4. The Lamb: 1x Soul listed twice.
   5. Mom's Dead Hand: an extra 1x Treasure.
   6. `scraper build` copied them into `data/b2.json`.
   7. `cardgen` stubs took the rewards and souls from there.
   8. Same kind of site error as Famine.
4. Tasks:
   1. Overrides in `pkg/card_db/overrides/b2.json`; rebuild `data/b2.json`.
   2. Engine cards get the printed rewards; The Lamb gets `Soul: 1`.
5. Done when:
   1. Go test: engine rewards match the card_db reward box, every card.
   2. Go test: The Lamb's kill gives 1 soul.
6. Result (2026-10-19)
   1. Overrides set the printed reward boxes; `data/b2.json` rebuilt.
   2. The rebuild changed only these three cards.
   3. Ring of Flies: 3¢; Mom's Dead Hand: 4¢.
   4. The Lamb: 3¢ and `Soul: 1`.
   5. Test `TestRewardsMatchCardDB`: every monster, fixed rewards.
   6. It failed on Ring of Flies and Mom's Dead Hand only.
   7. `TestSoulsMatchCardDB` failed on The Lamb; both pass now.
   8. `TestTheLamb` expects soul value 2: its 1 and the stolen soul.
   9. Old records replay apart from a kill of one of these cards on.
   10. Engine and data only; no screen change, no screenshots.
   11. Why: the site's "Potential Rewards" lists what a card may steal.
   12. Owner's finding; the page does not mark those rows.
   13. Tools repo: `scraper build` now stops on such a monster.
   14. It names the card until an override sets its `rewards`.
7. Owner notes (2026-10-19)
   1. Owner marks 19.24 solved.
### [x] 19.25 Krampus revealed mid-turn: no forced attack

1. Goal: "when this enters play" works for every card entering play.
2. Report (agent, 2026-10-19)
   1. Found in the code while adding Requiem's Keeper's Box.
   2. Krampus: "The active player must attack this once each turn if able."
   3. Revealed mid-turn (a slot refill), it does not force an attack.
   4. It forces one only from the next start of turn.
3. Cause:
   1. Only events and trinkets emit the "entered play" event.
   2. A monster, shop item or treasure entering play never did.
   3. So Krampus's "when this enters play" ability never triggers.
4. Tasks:
   1. Monsters and shop items entering a slot emit it; not during setup.
   2. Treasure gained from a deck emits it.
   3. A purchase does not: the item was in play in the shop already.
   4. Shop items' "enters play" abilities work in the shop (R-CARD-05, S-YUG).
5. Done when:
   1. Test: Krampus revealed mid-turn must be attacked that turn.
   2. Test: an item gained from the deck triggers "when this enters play".
6. Result (2026-10-19)
   1. Test `TestKrampusRevealedMidTurn`; engine tests `TestGainedItemEntersPlay`, `TestShopItemEntersPlay`.
   2. Engine only; no screen change, no screenshots.
7. Owner notes (2026-10-19)
   1. Owner marks 19.25 solved.
### [x] 19.26 Card zoom under the discard window

1. Goal: a window that must be on top is always on top.
2. Report (owner, 2026-10-19)
   1. The discard window asks which loot cards to discard.
   2. Hovering a card there zooms it under the window.
   3. The player cannot read the card they discard.
   4. Other windows also show under ones they should cover.
3. Cause:
   1. `.card-zoom` has z-index 10; `.choice` has 12.
   2. 19 z-index values in 6 CSS files, picked one by one.
   3. No list says which layer must be above which.
   4. Also wrong: the action backdrop (15) is above both choices.
   5. Also wrong: the dragged hand card (30) is above the overlays.
4. Tasks:
   1. One ordered list of layers: CSS variables in `layers.css`.
   2. Every z-index uses a variable, or is tagged local.
   3. Wiki page `Layers.md`: a row per layer, in order.
   4. Vitest test: code and wiki agree; no bare numbers.
   5. ADR 015 for the rule.
   6. Zoom above every window; effects above the zoom.
5. Done when:
   1. e2e test: a zoom from the discard window is on top.
   2. Vitest `layer inventory` passes.
6. Result (2026-10-19)
   1. `frontend/src/layers.css`: 14 layers, top first.
   2. Every z-index names one; four inner ones tagged local.
   3. Zoom moved above every window, below effects.
   4. Card menu and ability dialog above the prompts, on purpose now.
   5. Dragged hand card moved under the overlays.
   6. Rules and reasons: wiki `Layers.md`; decision: ADR 015.
   7. Vitest `layer inventory`; e2e `layers.spec.ts`.
   8. The e2e test failed before the fix; passes after.
   9. Use case UC-11 gets 3e: a card in a window.
   10. Screenshots: `.cache/screenshots/19.26/before/` and `after/`.
7. Owner notes (2026-10-19)
   1. Owner marks 19.26 solved.

### [x] 19.27 Cards to the bottom: no order prompt

1. Goal: cards put on a deck's bottom go there in random order.
2. Report (owner, 2026-10-19)
   1. Players are asked to order cards going to the bottom.
   2. No card reads the bottom of a deck; the prompt seems useless.
   3. Decision: always random; marked questionable (Q-22).
3. Cause:
   1. The Moon, The Emperor, The Hermit, Flush!, Eden ask the order.
   2. Black Candle says "in any order" but keeps deck order.
   3. Bobby Bomb, Mom's Bottle of Pills keep deck order.
   4. Deck order leaks hidden cards: everyone saw what went down.
4. Tasks:
   1. One engine helper puts cards on the bottom in random order.
   2. Every multi-card "to the bottom" uses it.
   3. Remove the order questions and Eden's lowest-card choice.
   4. Q-22 in `open-questions.md`, basis House.
5. Done when:
   1. Card tests: no order question; the cards end at the bottom.
   2. `make lint test` pass.
6. Result (2026-10-19)
   1. `Game.ToBottomRandom` shuffles the cards onto the bottom.
   2. Used by The Moon, The Emperor, The Hermit, Flush!, Eden.
   3. Also Black Candle, Ouija Board, Bobby Bomb, Mom's Bottle of Pills.
   4. Black Candle's "in any order" is House-overridden (Q-22).
   5. `ChooseLowest` is never asked; kept so purpose numbers stay.
   6. One-card "you may put it on the bottom" is unchanged.
   7. Tests: Eden, Flush!, The Moon, The Emperor, The Hermit.
   8. No screen change, so no screenshots.
7. Owner notes (2026-10-19)
   1. Owner marks 19.27 solved.

### [x] 19.28 e2e: the 4-player table test hangs on a hand back

1. Goal: the 4-player table test passes every run.
2. Report (agent, 2026-10-19)
   1. Found by `make e2e` during 15.29; Chromium.
   2. `table.spec.ts`: "four players sit left, top and right in turn order".
   3. `expectOpponents` waits on `.hidden-back` nth(2) of the top player.
   4. The element was gone; `boundingBox` waited until the 30 s timeout.
   5. Passed five reruns in a row.
   6. Evidence: `.cache/e2e-evidence/15.29-run4/`.
   7. Seen again in `make screenshots`: `.cache/e2e-evidence/15.29-shots2/`.
3. Tasks:
   1. Find why the top player's hand backs change during the check.
   2. Fix the client or the test, whichever is wrong.
4. Done when:
   1. The test passes ten runs in a row.
5. Cause
   1. Came with 15.29: the test now bans Eden (its choice covers the decks).
   2. Bans route the host's socket through the test runner; updates come later.
   3. The checks then ran while the start's effects still moved cards.
   4. A hand back or character measured a moment before was gone.
   5. Under heavy load the 4-player start alone also passed 30 s.
   6. The client is right; the test checked too early.
6. Result (2026-10-19)
   1. Board-layout tests wait until each page has played its updates (`settled`).
   2. The 4-player test is marked slow: 90 s.
   3. The test passed ten runs in a row (single runs and full suites).

### [x] 19.29 Bonus souls show no counters

1. Goal: a bonus soul's counters read on its card.
2. Report (owner, 2026-10-19)
   1. Soul of Wrath and the like show no counters.
   2. The progress to the 6th counter is hard to track.
3. Cause
   1. The engine puts the counters on and the view sends them.
   2. The board draws bonus souls without passing `counters`.
4. Tasks:
   1. Bonus souls on the board get the counter badge.
5. Done when:
   1. e2e test: a bonus soul with counters shows the badge.
   2. Screenshots: `.cache/screenshots/19.29/before/` and `after/`.
6. Result (2026-10-19)
   1. The board passes each bonus soul's counters to its card.
   2. The green counter badge shows, as on other cards (15.23).
   3. e2e test: `table.spec.ts`, "a bonus soul shows its counters".
   4. UC-12 2d updated.
   5. Screenshots: `.cache/screenshots/19.29/before/` and `after/`.
7. Owner notes (2026-10-19)
   1. Owner marks 19.29 solved.

### [x] 19.30 There's Options: the top treasure card is not shown

1. Goal: There's Options lets its owner see the top treasure card.
2. Report (owner, 2026-10-19)
   1. There's Options does not work.
   2. Its owner does not see the top treasure card.
3. Cause
   1. The engine is right: the view sends `treasure_top` (`PlayerView`).
   2. Only to the owner, only on their turn.
   3. The client never reads `treasure_top`; the deck stays face down.
4. Tasks:
   1. The treasure deck shows the top card face up to its owner.
   2. A click on the deck still buys from it.
5. Done when:
   1. e2e test: the owner sees the top card face up on the deck.
   2. Screenshots: `.cache/screenshots/19.30/before/` and `after/`.
6. Result (2026-10-19)
   1. The treasure deck shows `treasure_top` face up to its owner.
   2. The deck keeps its click, so buying from it works as before.
   3. Engine unchanged; its test already covers who gets the card.
   4. e2e test: `table.spec.ts`, "There's Options shows its owner …".
   5. UC-18 1a added.
   6. Screenshots: `.cache/screenshots/19.30/before/` and `after/`.
7. Owner notes (2026-10-19)
   1. Owner marks 19.30 solved.

### [x] 19.31 e2e: a reused host name breaks two tests

1. Goal: e2e tests never clash on player names.
2. Report (owner, 2026-10-19)
   1. Tests failed again on names; seen in 19.30 ("Wen").
   2. Asked for a fix once and for all.
3. Cause
   1. All tests share one server lobby.
   2. Guests join by the row text "<host>'s game".
   3. Two games with one host name: strict-mode failure in both.
   4. Finished games stay listed, so `--repeat-each` always clashes.
   5. The only guard was a rule to grep names by hand.
4. Tasks:
   1. Each lobby row carries its game id (`data-game`).
   2. Guests find the row by the host's game id, not its name.
   3. `startGame` and `lobby.spec.ts` use one helper.
   4. Docs drop the unique-names rule.
5. Done when:
   1. Two tests with the same names pass together.
   2. `--repeat-each 3` passes on the table tests.
   3. `make lint frontend e2e` pass.
6. Result (2026-10-19)
   1. Each Games row has `data-game` with its game id.
   2. `gameId(page)` reads the host's game id; `lobbyGame` finds its row.
   3. `startGame` and every `lobby.spec.ts` join use them.
   4. Names may now repeat; the unique-names rule is gone.
   5. `table` and `lobby` specs passed with `--repeat-each 3` (168 runs).
   6. Wiki `Frontend-tests.md` updated.
   7. No screen change, so no screenshots.
7. Owner notes (2026-10-19)
   1. Owner marks 19.31 solved.

### [x] 19.32 Dedicated server: a finished game takes the player back

1. Goal: after a finished game, the player can open a new game.
2. Report (owner, 2026-10-19)
   1. Dedicated server; the game finished; the player left it.
   2. They connected again to create a new game.
   3. The server always took them back to the finished game.
3. Cause
   1. The client keeps the seat token in `rejoin` (N-08).
   2. Its next hello to that server sends the token.
   3. `Hub.hello` attaches a token to any started game with it.
   4. A finished game stays in `h.games` while the server runs.
   5. So the hello lands in the finished game, every time.
   6. A hosted game is not hit: leaving stops the app's server.
4. Tasks:
   1. Hello never attaches a token to a finished game.
   2. The player keeps the token; the record download still works.
5. Done when:
   1. Go test: hello after the game is over gets the lobby.
   2. Go test: a new game with that token is joined, not the old one.
   3. `make lint test` pass.
6. Result (2026-10-19)
   1. `Hub.hello` skips a started game that is over.
   2. Room `over` is set at every game end, records or not.
   3. Before, it was set only when the server kept records.
   4. Go test: `TestFinishedGameLetsGo` in `hub_test.go`.
   5. It fails without the fix: the hello got seat 0 of the old game.
   6. UC-07 3c added; wiki `Server.md` updated.
   7. No screen change, so no screenshots.
7. Owner notes (2026-10-19)
   1. Owner marks 19.32 solved.

### [x] 19.33 e2e: a Ready click in the lobby never returns

1. Goal: `startGame` never hangs on the Ready click.
2. Report (agent, 2026-10-20)
   1. Seen once in `make e2e` during the 15.36 fix.
   2. Test: "curses sit at the far left …" (15.55), chromium, 4 players.
   3. Playwright clicked Ready; the click never returned; 30 s timeout.
   4. Rerun alone 8 times: all passed.
   5. Evidence: `.cache/e2e-evidence/19.33-ready-click/`.
   6. The error snapshot shows a page of another test, not this one.
3. Cause
   1. Not found yet.
   2. It hangs before the table, so the 15.36 presence code is not mounted.
4. Tasks:
   1. Find why the click never returns; fix the client or the test.
5. Done when:
   1. The cause is named here; a test covers it if it is in the code.
6. Owner notes (2026-10-20)
   1. Owner marks 19.33 solved; likely the same cause as 19.35.
   2. Not seen since the e2e runs use three workers.
   3. If it shows up again, this item reopens.

### [x] 19.34 e2e: the thrown die test misses its sound

1. Goal: the 15.62 throw test checks the sound every run.
2. Report (agent, 2026-10-20)
   1. Seen in `make e2e` during the 15.36 fix, webkit.
   2. "a thrown die falls, bounces and lands …": no sound was played.
   3. Rerun 10 times: 2 failed; 10 of 10 passed on the 15.35 commit.
   4. Evidence: `.cache/e2e-evidence/15.36-throw/`.
3. Cause
   1. The sound is a real-time timer: 28% of 1600 ms after the start.
   2. The test seeks the paused throw to 75%, then plays it.
   3. The rest ends in 400 ms; the effect unmounts and clears the timer.
   4. So the sound plays only if the test's own steps took long enough.
   5. In the app the throw is never sought, so it always clatters.
4. Tasks:
   1. The test waits for the sound while the throw is paused.
5. Done when:
   1. The test passes 20 runs in a row on webkit and chromium.
6. Result (2026-10-20)
   1. The test polls for the sound before it plays the throw on.
   2. Repeated 20 times on webkit and chromium: 40 of 40 passed.
   3. No app change; no screenshots.

### [x] 19.35 e2e: game starts time out in a busy full run

1. Goal: a full `make e2e` passes on a busy machine.
2. Report (agent, 2026-10-20)
   1. Seen in most full runs during 15.74 and 15.75.
   2. Mostly chromium `table.spec.ts`, around tests 103 – 107: 30 s timeouts.
   3. Also the 5 s Seats check in `startGame`, and `shell.spec.ts` once.
   4. Each failed test passes alone, 8 of 8.
   5. Load average 20 – 35 during the runs; chromium logs CVDisplayLink errors.
   6. The error snapshots show another test's page, as in 19.33.
   7. Evidence: `.cache/e2e-evidence/15.74-run*`, `15.75*`, `19.35-ab/`.
3. Cause
   1. Not found yet; not the 15.74 or 15.75 changes.
   2. A/B, two full runs each: the 15.73 commit failed 4, then 0.
   3. The current code failed 2, then 1, in other tests.
   4. `startGame` contexts stay open until the worker ends.
   5. One Go server and one Vite server serve every test of a browser.
   6. Maybe the same cause as 19.33.
4. Tasks:
   1. Find what the slow tests wait for; fix the tests or the client.
5. Done when:
   1. The cause is named here; three full runs in a row pass.
6. Owner answers (2026-10-20)
   1. Every test ran at once; no PC handles that.
   2. Give the tests room: a sensible number of workers.
7. Result (2026-10-20)
   1. Cause: 5 workers (half of 10 cores), each keeping 2 – 4 pages per test.
   2. 46 browser processes were alive at once; game starts hit 30 s.
   3. `workers: 3` in `e2e/playwright.config.ts`.
   4. Makefile `E2E_WORKERS ?= 3` for `e2e` and `screenshots`; `make` lists it.
   5. Three full runs in a row: 228 passed each, in 4.6 – 4.9 min.
   6. As fast as with 5 workers, which took 4.5 – 5 min.
   7. Wiki `Frontend-tests`.
   8. 19.33 may share this cause; left open until it is seen again.
8. Owner notes (2026-10-20)
   1. Owner marks 19.35 solved.

### [x] 19.36 XX Judgement destroys no soul

1. Goal: "destroy a soul you control" always destroys one, if able.
2. Report (owner, 2026-10-20)
   1. Play test; record `2026-10-20-104330-g1_Hitokiri-Battosai_jеня_KING.fsrec`.
   2. The Miser played XX Judgement on The Forgotten.
   3. The Forgotten kept every soul.
3. Cause
   1. Replay to step 1633: options `the_bone`, `lost_soul`, `blastocyst`.
   2. The player picked The Bone, a soul made from the starting item.
   3. A starting item is Eternal; it stayed Eternal as a soul.
   4. `DestroyObject` refuses Eternal objects; the card did nothing.
   5. Same hole: Wizoob, Curse of Loss, Forget Me Now.
4. Owner answers (2026-10-20)
   1. A soul that "loses all abilities" loses Eternal too (The Bone).
   2. Souls that can't be destroyed stay; they are not offered.
   3. The Lost's own soul is one of them.
5. Tasks:
   1. `Eternal` is false for a card with no abilities.
   2. Destroy-a-soul choices list only souls that can be destroyed.
6. Done when:
   1. Tests: a Bone soul is destroyed; an Eternal soul is not offered.
   2. The record's case: Judgement on The Forgotten destroys a soul.
7. Result (2026-10-20)
   1. `Eternal` is false for a card with no abilities (`combat.go`).
   2. New `Game.DestroyableSouls`: souls that are not Eternal.
   3. XX Judgement, Wizoob, Curse of Loss, Forget Me Now use it.
   4. The record replayed: the same pick now destroys The Bone.
   5. Ruling: Q-23 in `open-questions.md`, House, in "Rulings to revisit".
   6. Tests: `TestTheBoneSoulCanBeDestroyed`, `TestXXJudgementSkipsEternalSouls`.
   7. Engine only; no screen change, no screenshots.
8. Owner notes (2026-10-20)
   1. Owner marks 19.36 solved.

### [x] 19.37 e2e: the 15.77 shot test misses its start by a pixel

1. Goal: the 15.77 test checks the shot's start without flaking.
2. Report (agent, 2026-10-20)
   1. `make e2e` gate for 15.42, chromium.
   2. `playback.spec.ts`: "damage from a card no longer on screen…".
   3. The shot's first tracked box was 30.5 px from the player; limit < 30.
   4. Reruns: 1 of 16 failed with 15.42; 2 of 12 failed without it.
   5. So it is older than 15.42; the sky is not on the game table.
   6. Evidence: `.cache/e2e-fails/20261020-15.42-gate/` and `-rerun/` (not in git).
   7. Again in the 15.80 gate, webkit: y off by 31.7 px.
   8. Evidence: `.cache/e2e-evidence/15.80-run2/` (not in git).
   9. Again in the 19.40 gate, webkit: y off by 30.5 px.
   10. Evidence: `.cache/e2e-evidence/19.40-run1/` (not in git).
3. Tasks:
   1. Find why the first tracked box is late: sampling or the shot's start.
   2. Fix the shot if it starts off the player; else the test.
4. Done when:
   1. Cause found and written here.
   2. The test passes 30 repeats in a row, both browsers.
5. Cause
   1. The shot starts on the player; the test measures it late.
   2. `trackBoxes` sampled only on animation frames.
   3. The test pauses the shot and jumps it to 150 ms for its picture.
   4. Playwright can see the shot before the browser draws a frame.
   5. Then the pause wins: the first box is at 150 of 450 ms.
   6. With `ease-in`, that is about 0.106 of a ~290 px flight.
   7. So the miss is the same each time: 30.55 px twice, 31.67 once.
   8. A busy machine (webkit, full run) makes the race likelier.
6. Result (2026-10-20)
   1. Test only; the shot is unchanged.
   2. `trackBoxes` takes each first box when the element joins the page.
   3. A `MutationObserver` runs before the test can pause the shot.
   4. Every `trackBoxes` test gains it, 15.76's roll shot too.
   5. 15.77 alone: 30 repeats per browser, 60 of 60 passed.
   6. `make lint frontend e2e` passed.
7. Owner notes (2026-10-20)
   1. Owner marks 19.37 solved.

### [ ] 19.38 e2e: game starts time out again in a full run

1. Goal: a full `make e2e` passes on the first run again.
2. Report (agent, 2026-10-20)
   1. `make e2e` gate for 15.79, webkit, 3 workers.
   2. Three tests failed in `startGame`, before the table.
   3. `info.spec.ts:12`: the 5 s Seats check.
   4. `playback.spec.ts:89`, `table.spec.ts:264`: 30 s, waiting for Nickname.
   5. Load average 17 after the run, 34 at the end of the next.
   6. The three alone, 4 repeats each: 12 of 12 passed.
   7. The next full run passed: 238 passed.
   8. Evidence: `.cache/e2e-evidence/15.79-run1/` (not in git).
   9. Same signs as 19.35, solved with 3 workers.
   10. Again in the 15.80 gate: chromium `prompt.spec.ts:199`, 120 s.
   11. It waited for the stack's attack target after the throw button.
   12. The error snapshot shows another test's game (Max, Nia).
   13. Alone, 6 repeats per browser: 12 of 12 passed.
   14. Evidence: `.cache/e2e-evidence/15.80-run1/` (not in git).
   15. Again in the 19.42 gate, webkit: the 5 s Seats check.
   16. `actions.spec.ts:676`, `table.spec.ts:191`; load average 13.
   17. The two alone, 3 repeats each: 6 of 6 passed.
   18. Evidence: `.cache/e2e-evidence/19.42-run1/` (not in git).
   19. Again in the third 19.42 gate run, chromium, two tests.
   20. `table.spec.ts:724`: 30 s, waiting for Join a game.
   21. `table.spec.ts:802`: the 5 s check for the table.
   22. Load average 34 during the run; the run took 10.3 min.
   23. Evidence: `.cache/e2e-evidence/19.42-run3/` (not in git).
3. Tasks:
   1. Find what else loads the machine; maybe fewer workers for webkit.
4. Done when:
   1. Cause written here; three full runs in a row pass.

### [x] 19.39 Many dice show as one thrown die

1. Goal: a roll of several dice throws every die at once.
2. Report (owner, 2026-10-20)
   1. Classic Roller: "rolls 2 dice instead".
   2. The screen shows one thrown die; it is confusing.
   3. Mutant Spider should throw 4 dice the same way.
3. Cause
   1. The engine draws all the dice in `throwDie` (`ManyDice`).
   2. `dice_rolled` carries only the first die's value.
   3. So the client throws one die.
   4. Same for Cursed Lil Haunt's 2 dice.
4. Tasks:
   1. `dice_rolled` carries every die: `dice`.
   2. The client throws them side by side, together.
5. Done when:
   1. Go test: Classic Roller's and Mutant Spider's events carry all dice.
   2. Vitest: the effect keeps every value.
   3. e2e test: a 2-dice roll throws two dice.
   4. Wiki `Animations.md` row updated.
6. Result (2026-10-20)
   1. `dice_rolled` of many dice carries them all: `dice`.
   2. The throw drops every die in a row, each a bit later.
   3. Classic Roller throws 2; Mutant Spider 4; Cursed Lil Haunt 2.
   4. A single die throws as before.
   5. Records gain the field; older records still play.
   6. Screenshots: `.cache/screenshots/19.39/before/` and `after/`.
7. Owner review (2026-10-20)
   1. One clatter for many dice sounds wrong.
   2. Each die plays its own clatter as it lands.
   3. Dice land at small, uneven gaps, not all at once.
   4. It should feel like a real throw with a friend.
8. Result after review (2026-10-20)
   1. Dice start at small uneven gaps: `diceDelays` in `dice.ts`.
   2. Each die clatters as it lands: 4 dice, 4 clatters.
   3. Screenshot mid-throw: `after/screen-dice-many-falling.png`.
9. Owner notes (2026-10-20)
   1. Owner marks 19.39 solved.

### [x] 19.40 Thrown dice show gaps between their sides

1. Goal: a thrown die looks like one solid cube.
2. Report (owner, 2026-10-20)
   1. Seen while reviewing 19.39.
   2. The sides have gaps between them.
   3. The die looks like image on image, not a 3D object.
3. Cause
   1. Each cube face has round corners: `border-radius: 14%`.
   2. The texture's faces are full squares.
   3. So every edge and corner of the cube has a hole.
4. Tasks:
   1. Cube faces meet edge to edge.
5. Done when:
   1. e2e test: cube faces have no round corners.
   2. Screenshot: no gaps at the edges.
6. Result (2026-10-20)
   1. Cube faces are square; the sides meet edge to edge.
   2. The bevel's light and shade still soften the edges.
   3. Flat dice (stack icons, last roll) keep round corners.
   4. Screenshots: `.cache/screenshots/19.40/before/` and `after/`.
7. Owner notes (2026-10-20)
   1. Owner marks 19.40 solved.

### [x] 19.41 A resolved card shows on the stack again after its fly

1. Goal: a card that flew from the stack to the discard is gone.
2. Report (owner, 2026-10-20)
   1. A card on the stack resolves and flies to the discard.
   2. Then it shows on the stack again, hides, and shows again.
   3. The fly already said the card moved.
3. Cause
   1. The large played card hides only while its fly plays (15.74).
   2. The stack row never hides.
   3. More effects of the same update often follow the fly.
   4. Records: `discarded` then `gained_cents` (A Penny), or `looted` (A Sack).
   5. Between and after them, the old view shows the card again.
   6. It goes only when the update's view is drawn.
4. Tasks:
   1. A card that left the stack stays hidden until the next view.
   2. Both its large card and its stack row.
5. Done when:
   1. Vitest: the fly names the stack item it leaves.
   2. e2e test: after the fly, while a popup plays, the card stays hidden.
6. Result (2026-10-20)
   1. The fly names the stack item it leaves: `leaves`.
   2. The store keeps it in `left` until the next view is shown.
   3. The large card and the stack row stay hidden meanwhile.
   4. Popups and loot flies after it play over an empty spot.
   5. Screenshots: `.cache/screenshots/19.41/before/` and `after/`.
7. Owner notes (2026-10-20)
   1. Owner marks 19.41 solved.

### [x] 19.42 A "no" answer shows the No! card

1. Goal: a word answer shows as a plain button, never as a card.
2. Report (owner, 2026-10-20)
   1. P.H.D. and Mom's Box ask the player to change something.
   2. The "no" button shows the No! card above it.
   3. The card has nothing to do with the choice.
3. Cause
   1. The client draws an option as a card when its label is a card ref.
   2. Rule: `isCard(label)` in `actionsFor` (`actions.ts`).
   3. No! has the ref `no`; the answer word is also `no`.
   4. So every yes/no question shows No! on its "no" option.
   5. A scan of engine option words against 696 refs: `no` is the only clash.

4. Tasks:
   1. The decline word is `No`, not `no` (owner, 2026-10-20).
   2. A later translation may replace the word with a token.
5. Done when:
   1. Go test: P.H.D.'s and Mom's Box's answers are no card refs.
   2. e2e test: a word answer is a button, with no card.
6. Result (2026-10-20)
   1. Every "no" answer reads `No`.
   2. P.H.D., Mom's Box, every "may" question, Isaac's Blessing, Telepathy.
   3. The client is unchanged; `No` is no card ref.
   4. Records keep answers by index; older records still play.
   5. Screenshots: `.cache/screenshots/19.42/before/` and `after/`.
7. Owner notes (2026-10-20)
   1. Owner marks 19.42 solved.
   2. Committed with `make e2e` failing on 19.38, 19.43, 19.44 (owner).

### [ ] 19.43 e2e: the 15.74 fly test misses the fly in a full run

1. Goal: the test sees the fly in a busy full run too.
2. Report (agent, 2026-10-20)
   1. Second `make e2e` gate run for 19.42.
   2. chromium `playback.spec.ts:868`, the stack row case.
   3. `effect-fly` never became visible in 5 s.
   4. Load average 25 at the end of the run.
   5. Alone, 6 repeats per case and browser: 24 of 24 passed.
   6. Evidence: `.cache/e2e-evidence/19.42-run2/` (not in git).
3. Cause
   1. Not found yet.
   2. Guess: the short fly ends between two checks under load.
4. Tasks:
   1. Find the cause in the evidence; fix the test or the app.
5. Done when:
   1. Cause written here; three full runs in a row pass.

### [ ] 19.44 e2e: the V-11 test gives up on a target prompt

1. Goal: the V-11 stack test passes with any loot card it plays.
2. Report (agent, 2026-10-20)
   1. Third `make e2e` gate run for 19.42, webkit.
   2. `stack.spec.ts:109`: the stack shows no "Waiting for" line.
   3. The stack said "Choose a target"; the loot waited for its target.
   4. The test answers only `.choice` dialogs, for 10 tries of 200 ms.
   5. Load average 34 during the run.
   6. Alone, 3 repeats per browser: 6 of 6 passed.
   7. Evidence: `.cache/e2e-evidence/19.42-run3/` (not in git).
3. Cause
   1. Not found yet; the loot card is not known.
   2. Maybe a table target, which the test never picks.
   3. Maybe the dialog came later than 2 s under load.
4. Tasks:
   1. Find the card in the match record; fix the test or the app.
5. Done when:
   1. Cause written here; three full runs in a row pass.

