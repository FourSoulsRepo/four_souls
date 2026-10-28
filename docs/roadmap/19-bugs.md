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
| A  | Rules: a card or rule coded wrong        | Engine     | 19.22, 19.23, 19.25, 19.27, 19.36, 19.46, 19.60, 19.61, 19.62, 19.64, 19.65, 19.66, 19.92, 19.94, 19.95, 19.97 | 16 |
| B  | Events and view: data never sent         | Engine     | 19.2, 19.9, 19.11, 19.39, 19.67, 19.71, 19.72, 19.73, 19.90 | 9 |
| C  | Card data: wrong values in `card_db`     | Card data  | 19.24, 19.45                                           | 2     |
| D  | View data sent, never drawn              | Client     | 19.5, 19.29, 19.30, 19.74                              | 4     |
| E  | Playback, animations, history            | Client     | 19.1, 19.3, 19.4, 19.6, 19.13, 19.41, 19.52, 19.54, 19.69 | 9  |
| F  | Layout and visuals                       | Client     | 19.15, 19.26, 19.40, 19.49, 19.50, 19.76, 19.88        | 7     |
| G  | Prompts: the choice is unclear           | Client     | 19.16, 19.17, 19.19, 19.42                             | 4     |
| H  | Lobby and rejoin                         | Server     | 19.32                                                  | 1     |
| I  | e2e test logic: the app was right        | e2e        | 19.7, 19.8, 19.10, 19.18, 19.21, 19.28, 19.31, 19.34, 19.37, 19.48, 19.58, 19.70, 19.79, 19.80, 19.81, 19.82, 19.85, 19.87, 19.100 | 19 |
| J  | e2e under load, or cause not found       | e2e        | 19.12, 19.14, 19.20, 19.33, 19.35, 19.38, 19.43, 19.44, 19.51, 19.53, 19.55, 19.56, 19.57, 19.75, 19.77, 19.89, 19.101 | 17 |
| K  | Build and release files                  | Release    | 19.47, 19.59, 19.83, 19.91                             | 4     |
| L  | Not a bug: the app follows a ruling      | None       | 19.63                                                  | 1     |
| M  | Go test setup: the code was right        | Go tests   | 19.68                                                  | 1     |
| N  | A missing control: no way to do it       | Server     | 19.84                                                  | 1     |
| O  | Reconnect: old connection state kept     | Client     | 19.78, 19.86                                           | 2     |
| P  | A missing setting: no way to change it   | Client     | 19.93                                                  | 1     |

Totals by system (94 items):

| System    | Groups  | Count |
|-----------|---------|-------|
| Engine    | A, B    | 23    |
| Card data | C       | 2     |
| Client    | D–G, O, P | 27  |
| Server    | H, N    | 2     |
| e2e       | I, J    | 34    |
| Release   | K       | 4     |
| None      | L       | 1     |
| Go tests  | M       | 1     |

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

### [x] 19.14 e2e: a guest's Join click can miss under load

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
   4. Seen again in the 19.38 worker sweep (2026-10-21), no swap.
   5. 8 workers, webkit, `first-roll.spec.ts:58`: the 5 s Seats check.
   6. 4 workers, webkit, `table.spec.ts:399`: the same; load 7.4.
   7. Evidence: `.cache/e2e-evidence/sweep-w8-184425/`, `sweep-w4-191256/`.
6. Cause (2026-10-21)
   1. A probe logged the guest's WebSocket: no `join` was ever sent.
   2. So the click never reached the Join button's handler.
   3. The list shows waiting tables first, then started games.
   4. A table above that starts moves every waiting row below up 46 px.
   5. Playwright checks the button, then clicks at its screen point.
   6. A row that moved in between left a started game's disabled Join there.
   7. More tests at once start more tables: more workers, more misses.
   8. Evidence: `.cache/e2e-evidence/19.14-probe-195313/` (not in git).
7. Result (2026-10-21)
   1. `startGame` waits for its row's Join to be enabled.
   2. Then it sends the click to that button (`dispatchEvent`), not to a point.
   3. The Join screen itself is still clicked by the lobby tests.
   4. Test only; no app change, no screenshots.
8. Owner notes (2026-10-21)
   1. Owner marks 19.14 solved.

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

### [x] 19.18 e2e: a played Dagaz waits for a player target

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
   13. 19.44 is the same: Dagaz again, found in its record.
6. Result (2026-10-21)
   1. New `answerPlay` in `turn.ts`: a mode in a dialog, or a target on the table.
   2. It clicks the first uncovered `.hl-aim` card or stack row, as `settle` does.
   3. It waits up to 10 s for the stack row, not 10 tries of 200 ms.
   4. Both stack tests use it.
   5. Test: "a played Dagaz goes on the stack once its player is picked".
   6. It gets Dagaz from a changed save; helpers moved to `saves.ts`.
   7. It fails without the table click, both browsers; passes with it.
   8. Test only; no app change, no screenshots.
7. Owner notes (2026-10-21)
   1. Owner marks 19.18 solved.

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

### [x] 19.20 e2e: a cancelled Dagaz shows twice in the hand

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
5. Cause (2026-10-21)
   1. An unlucky deal; the app is right.
   2. `giveCard` turns the active player's first hand card into Dagaz.
   3. Sometimes the deal had already given that player a Dagaz.
   4. A probe dumped the store: two Dagaz objects in the hand (319, 320).
   5. 30 repeats per browser: 2 of 60 failed, both like that.
   6. Evidence: `.cache/e2e-evidence/19.20-probe-193048/` (not in git).
6. Result (2026-10-21)
   1. The test bans Dagaz from the deal (14.10); the reason is in a comment.
   2. `giveCard` still makes the one Dagaz the test needs.
   3. 40 repeats per browser: 80 of 80 passed.
   4. Test only; no app change, no screenshots.
7. Owner notes (2026-10-21)
   1. Owner marks 19.20 solved.

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

### [x] 19.38 e2e: game starts time out again in a full run

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
   24. Again in five 8.24 gate runs on 2026-10-21, one failure each.
   25. Each a different test; once the `?fake` Settings page alone.
   26. A full run took 7.9 – 10.7 min; 3.5 – 4.9 min the day before.
   27. Two workers instead of three: 10.7 min, still a failure.
   28. Evidence: `.cache/e2e-evidence/8.24-*/` (not in git).
3. Tasks:
   1. Find what else loads the machine; maybe fewer workers for webkit.
4. Done when:
   1. Cause written here; three full runs in a row pass.
5. Cause
   1. `startGame` opens a context per player; nothing closed them.
   2. They lived until the worker ended: ~175 game pages per browser.
   3. Each keeps a WebSocket and a looping glow (end turn, prompt, hand).
   4. Headless pages are never hidden, so the glows never stop.
   5. The laptop has 16 GB; ~12 GB was in use before the run.
   6. Chromium alone: browser processes rose from 52 to 214.
   7. ~1.3M pages (~20 GB) went to swap in 4 min; load 24.
   8. Later tests of each browser ran ~2× slower than earlier ones.
   9. A failure hit whichever test ran while the machine swapped.
   10. Fewer workers did not help: the pages per run stay the same.
6. Result (2026-10-21)
   1. The automatic fixture closes every context a test opened.
   2. It runs after the clash check, in `e2e/tests/fixtures.ts`.
   3. Chromium alone: 2.8 min instead of 4.1; 19 – 26 processes.
   4. Swap: ~2,400 pages instead of ~1.3M.
   5. Three full runs in a row: 258 passed each, 5.8 – 6.3 min, load < 10.
   6. The first try failed once on 19.48, a separate bug.
   7. Wiki `Frontend-tests`.
   8. Worker sweep after 19.48, full runs, none swapped.
   9. 3 workers: 6.1 min, 4 of 4 clean.
   10. 4 workers: 4.5 – 4.7 min, 2 of 4 clean (19.14, 19.20, 19.51).
   11. 5 workers: 3.6 – 3.9 min, 3 of 4 clean (19.49, 19.50).
   12. 6 workers: 3.4 min, 1 of 1; 8 workers: 2.9 min, failed (19.14).
   13. Workers stay 3 until those items are fixed; then try 5.
   14. After 19.14 – 19.51: 5 workers, three full runs, 264 passed each.
   15. 3.9 – 4.2 min instead of 6; load at most 12; no swap.
7. Owner notes (2026-10-21)
   1. Owner marks 19.38 solved.

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

### [x] 19.43 e2e: the 15.74 fly test misses the fly in a full run

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
   3. Likely 19.38: the machine swapped; check again after it.
   4. The discard fly lasts 0.45 s; the test polls it with `toBeVisible`.
   5. A swapping page can stall longer; the poll then misses the fly.
   6. Same kind as 19.37: `trackBoxes` already records the fly.
   7. Not seen in any full run since 19.38 (3 – 8 workers).
6. Result (2026-10-21)
   1. The test waits for the fly's box from `trackBoxes`, not a poll.
   2. `trackBoxes` takes it when the fly joins the page (19.37).
   3. Not reproduced under CPU load: 80 of 80 passed before the change too.
   4. The failure needed a swapping machine (19.38); the poll window is gone.
   5. Test only; no app change, no screenshots.
4. Tasks:
   1. Find the cause in the evidence; fix the test or the app.
5. Done when:
   1. Cause written here; three full runs in a row pass.
7. Owner notes (2026-10-21)
   1. Owner marks 19.43 solved.

### [x] 19.44 e2e: the V-11 test gives up on a target prompt

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
   4. Likely 19.38: the machine swapped; check again after it.
   5. Found (2026-10-21): the same cause as 19.18, not load.
   6. Record replayed at `f289001`: Bea played Dagaz (object 334).
   7. The test picked its mode; then "Choose a player" waited.
   8. The test answers only dialogs; Dagaz's player is picked on the table.
   9. Fails whenever the deal gives the active player such a card.
6. Result (2026-10-21)
   1. Fixed with 19.18: the stack tests answer table targets.
4. Tasks:
   1. Find the card in the match record; fix the test or the app.
5. Done when:
   1. Cause written here; three full runs in a row pass.
7. Owner notes (2026-10-21)
   1. Owner marks 19.44 solved, with 19.18.

### [x] 19.45 Abyss never uncharges after use

1. Goal: Abyss taps to activate, like its printed card.
2. Report (owner, 2026-10-21)
   1. Abyss stays charged after its ability is used.
3. Cause
   1. The printed card shows the ↷ (Tap) icon.
   2. foursouls.com shows the Paid icon instead.
   3. `scraper build` copied `{Paid Effect}` into `data/r.json`.
   4. The text has no "cost:" part, so the engine card got no cost.
   5. With no cost, nothing taps it; it can be used any number of times.
   6. `TestAbyss` even activates it twice in one turn.
   7. It is the only `{Paid Effect}` text with no cost in card_db.
4. Tasks:
   1. Override in `pkg/card_db/overrides/r.json`; rebuild `data/r.json`.
   2. Engine Abyss gets the `Tap()` cost and "↷:" text.
5. Done when:
   1. Go test: a second use in one turn is refused.
   2. Go test: it works again after it recharges.
6. Result (2026-10-21)
   1. Override sets the printed `{Tap Effect}` in `overrides/r.json`.
   2. `data/r.json` rebuilt; the rebuild changed only Abyss.
   3. Engine Abyss: `Tap()` cost, `Tap: true`, "↷:" text.
   4. `Tap: true` lets the death penalty deactivate it (R-DEATH-14).
   5. `TestAbyss`: a second use is refused (R-ABIL-07, no cost).
   6. It destroys a second item after its next turn recharges it.
   7. Old records that use Abyss twice in a turn no longer replay.
   8. Engine and data only; no screen change, no screenshots.
7. Owner notes (2026-10-21)
   1. Owner marks 19.45 solved.

### [x] 19.46 A lost {HP} bonus takes HP from the bottom

1. Goal: lowering max HP removes marked damage first (R-MECH-21).
2. Report (owner, 2026-10-21)
   1. Play test; record `2026-10-20-182229-g1_Hitokiri-Battosai_jеня.fsrec`.
   2. The Empty used Abyss on Bag Lunch (with a counter) in an attack.
   3. HP became 0; the game stood still; the attack ended.
   4. HP should go down from the top, not from the bottom.
3. Cause
   1. Replay on `f289001`, step 909: max HP 3, damage 2, HP 1.
   2. Step 910: Bag Lunch destroyed; max HP 2, damage still 2, HP 0.
   3. `PlayerHP` is max minus damage; nothing applies R-MECH-21.
   4. A death goes on the stack only when damage resolves.
   5. So the player stood at 0 HP with no death (R-DEATH-01).
   6. Step 916: the attack ended at 0 HP (R-ATK-14), with no death.
   7. Same hole for monsters that lose a {HP} bonus.
4. Tasks:
   1. Any drop in max HP removes that much damage, down to 0.
   2. Players and monsters; checked after every change in `Apply`.
   3. An object left at 0 HP this way gets its death on the stack.
5. Done when:
   1. Go test: Bag Lunch destroyed at 1 HP of 3 leaves 1 HP of 2.
   2. Go test: a monster that loses a {HP} bonus keeps its HP left.
6. Result (2026-10-21)
   1. New `checkMaxHP` in `combat.go` applies R-MECH-21.
   2. It runs at the start of `Apply`, after the intent, and each `run()` step.
   3. Players and in-play objects with HP; damage never below 0.
   4. Left at 0 HP this way (max HP 0): death on the stack.
   5. The last max HP is not saved; checksums of other states stay.
   6. The record's case: Bag Lunch destroyed at 1 HP of 3 leaves 1 HP of 2.
   7. Tests: `TestAbyssOnBagLunchKeepsHP`, `TestLowerMaxHPRemovesDamage`.
   8. Test: `TestMaxHPToZeroKills`.
   9. The new test fails without the fix: HP 0, damage 2.
   10. Records with such a drop while damaged no longer replay.
   11. Wiki: "Max HP going down" in `Rules-engine.md`.
   12. Engine only; no screen change, no screenshots.
7. Owner notes (2026-10-21)
   1. Owner marks 19.46 solved.

### [x] 19.47 Windows file details miss the release version

1. Goal: macOS and Windows show the git tag's version.
2. Report (owner, 2026-10-21)
   1. Releases ship macOS and Windows binaries.
   2. `Info.plist` and the Windows details must show the tag's version.
3. Cause
   1. The build stamps `wails.json` `info.productVersion` (8.3).
   2. Wails renders `Info.plist` and `info.json` from it.
   3. Local build at 1.4.2: `Info.plist` shows 1.4.2 in both keys.
   4. The `.exe`: fixed product version 0.0.0.0; no `FileVersion` string.
   5. `info.json` set only the fixed file version and the product string.
   6. Docs said trial tags give 0.0.0; the build strips `-test` instead.
4. Tasks:
   1. `info.json`: fixed product version and `FileVersion` string too.
   2. Test: every version field uses `{{.Info.ProductVersion}}`.
   3. Fix the build action comment and `Releases.md`.
5. Done when:
   1. A local Windows build at 1.4.2 has 1.4.2 in all four fields.
   2. Go test fails if a version field drops the template.
6. Result (2026-10-21)
   1. `info.json`: fixed `product_version` and the `FileVersion` string.
   2. Local Windows build at 1.4.2: all four fields show 1.4.2.
   3. Local macOS build at 1.4.2: both `Info.plist` keys show 1.4.2.
   4. Test: `TestVersionFields` in `identity_test.go`.
   5. It fails without the fix: `product_version` and `FileVersion` empty.
   6. Build action comment and `Releases.md` now say `-test` gives 1.2.3.
   7. Release files only; no screen change, no screenshots.
7. Owner notes (2026-10-21)
   1. Owner marks 19.47 solved.

### [x] 19.48 e2e: the layers arrow test loses the other's arrow

1. Goal: the 15.36 layers test keeps the other player's arrow.
2. Report (agent, 2026-10-21)
   1. `layers.spec.ts:109`, webkit, in 3 of 7 full runs on 2026-10-21.
   2. The arrow head is not found, or is gone one line later.
   3. Never failed before 2026-10-21.
   4. Alone, 30 repeats, webkit: 3 failed; load was low.
   5. Evidence: `.cache/e2e-evidence/19.38-gate1-140102/` (not in git).
   6. Evidence: `.cache/e2e-evidence/layers-arrow-repeat-*/` (not in git).
3. Cause
   1. The test injects a presence for the other seat on its own page.
   2. The other seat's real page sends presence too.
   3. A probe logged the store: `{}` from that seat 10 – 50 ms later.
   4. Only when that seat has the turn; its sender then shares hovers.
   5. Since 15.43, `throwFirstRoll` moves each roller's mouse in the game.
   6. Before, a guest's sender saw no mouse move and stayed silent.
   7. The app is right: it is that seat's real pointer.
4. Tasks:
   1. The test drops the server's presence frames; only its own shows.
5. Done when:
   1. 30 repeats in webkit and chromium pass.
6. Result (2026-10-21)
   1. `noPresence` in `layers.spec.ts`, through `startGame`'s `open`.
   2. 30 repeats per browser: 60 of 60 passed; before, 3 of 30 failed.
   3. Test only; the app is unchanged; no screenshots.
7. Owner notes (2026-10-21)
   1. Owner marks 19.48 solved.

### [x] 19.49 e2e: a Settings click in the main menu is lost

1. Goal: the main menu opens Settings on the first click.
2. Report (agent, 2026-10-21)
   1. `settings.spec.ts:8`, chromium, 8.24 run 4, 3 workers.
   2. `playback.spec.ts:30`, webkit, worker sweep, 5 workers.
   3. Both: `?fake&nosplash`, then a click on Settings.
   4. The page still showed the main menu after the click.
   5. Load 9.6 in the sweep; no swap.
   6. Alone, both tests, 40 repeats each, both browsers: 160 of 160 passed.
   7. Evidence: `.cache/e2e-evidence/8.24-run4-125028/` (not in git).
   8. Evidence: `.cache/e2e-evidence/sweep-w5-185635/` (not in git).
3. Cause
   1. Not found yet.
   2. Not a Vite reload: e2e runs watch no files (19.12).
   3. `loadSettings` and `show` never set the screen back to the menu.
4. Tasks:
   1. Find what undoes the click, or why it never lands.
5. Done when:
   1. Cause written here; three full runs in a row pass.
6. Cause (2026-10-21)
   1. The menu is centred; its title changed size after the menu showed.
   2. Text title first, then the banner image, then text again if it failed.
   3. e2e has no images, so the banner always fails there.
   4. A probe: Settings moved 48 px (chromium) and 95 px (webkit) in ~20 ms.
   5. A click in between hit Records, which is disabled: nothing happened.
   6. With the images the menu moved too, once the banner loaded.
7. Result (2026-10-21)
   1. `Title` shows the text until the banner has loaded, then swaps once.
   2. A banner that fails never shows; the menu never moves for it.
   3. Test: "the menu buttons never move after they show (19.49)".
   4. Without the fix it failed 7 of 10; with it 10 of 10 passed.
   5. Screenshots: `.cache/screenshots/19.49/before/` and `after/`.
   6. The settled menu looks the same, with and without images.
   7. Wiki `Client`; UC-01 7a.
8. Owner notes (2026-10-21)
   1. Owner marks 19.49 solved.

### [x] 19.50 e2e: the Menu button is under a choice once

1. Goal: the 15.29 layers test passes every run.
2. Report (agent, 2026-10-21)
   1. `layers.spec.ts:79`, chromium, worker sweep, 5 workers.
   2. `onTop(menu, backdrop)` was false.
   3. The test injects a choice into a live game, like 19.48.
   4. The table said "Start of turn"; load 9.6, no swap.
   5. Evidence: `.cache/e2e-evidence/sweep-w5-185635/` (not in git).
3. Cause
   1. Not found yet.
   2. Maybe a live update or prompt raced the injected one (19.48).
4. Tasks:
   1. Find what was drawn over the Menu button.
5. Done when:
   1. Cause written here; 30 repeats in both browsers pass.
6. Cause (2026-10-21)
   1. A probe: a shop card's zoom covered the Menu button; the app was wrong.
   2. The player threw the first roll; the pointer stayed on the Throw button.
   3. The button went; the pointer was now over a shop card: its zoom opened.
   4. The card's image failed (no images in e2e) and became the drawn face.
   5. The hovered node left the page, so no leave event ever came.
   6. The pointer moved away; the zoom stayed open over everything.
   7. Players can meet it too: a hovered badge or image that is redrawn.
   8. 5 – 7 of 240 – 300 repeats failed, chromium.
7. Result (2026-10-21)
   1. `useZoom` also ends a mouse zoom on a pointer-over outside the card.
   2. Test: "a zoom closes when the pointer leaves a card whose face was redrawn".
   3. It removes the hovered node, then moves away; it failed before the fix.
   4. The layers test: 300 of 300 repeats passed.
   5. Wiki `Client`. Nothing changes on screen; no screenshots.
8. Owner notes (2026-10-21)
   1. Owner marks 19.50 solved.

### [x] 19.51 e2e: the fast burst test runs out of time

1. Goal: the burst test passes on a busy machine.
2. Report (agent, 2026-10-21)
   1. `playback.spec.ts:41`, webkit, worker sweep, 4 workers.
   2. "Turn 21" not shown in 20 s; load 6.7, no swap.
   3. The whole test took 15 – 16 s on 2026-10-20.
   4. On 2026-10-21 it took 17 – 26 s; no failure until this one.
   5. Evidence: `.cache/e2e-evidence/sweep-w4-190815/` (not in git).
3. Cause
   1. Not found yet.
   2. Twenty rolls play at full speed; the 20 s wait is close to that.
   3. Maybe the 15.43 first roll or the 19.39 dice made it longer.
4. Tasks:
   1. Measure the burst's playback; fix the app or the budget.
5. Done when:
   1. Cause written here; three full runs in a row pass.
6. Cause (2026-10-21)
   1. Measured on a quiet machine: the burst plays for 13.6 s, both browsers.
   2. As designed: a long queue skips or speeds effects; the last ones play in full.
   3. The 20 s wait left about 40% room; the 30 s test limit was close too.
7. Result (2026-10-21)
   1. The wait for "Turn 21" is 40 s; the test's limit is 90 s.
   2. Test only; no app change, no screenshots.
8. Owner notes (2026-10-21)
   1. Owner marks 19.51 solved.

### [x] 19.52 History: card-less stack items read "Roll — …"

1. Goal: the History says what a resolved stack item did, as the strip does.
2. Report (agent, 2026-10-21)
   1. Found while checking the 15.89 screenshots.
   2. Damage reads "Roll — damage", not "Bo deals 1 damage to Gaper".
   3. Also "Roll — death", "Roll — soul", "Roll — rewards", "Roll — death penalty".
   4. The strip shows the same items with the right words.
3. Cause
   1. `history.ts` builds the line as `card || 'Roll'` plus the label.
   2. It takes every stack item without a card for a roll.
   3. Damage, death, death step and penalty items have no card.
   4. Written in 7.9; its tests had only a loot card and a roll.
   5. 15.22 and 15.47 fixed the strip's words (`stackLine`), not the History.
   6. 15.88 counted `stack_resolved` as told without checking its words.
4. Tasks:
   1. A resolved item's History line uses `stackParts`, as the strip does.
   2. Rolls keep "(rolled 6 → 4)" (19.1).
5. Done when:
   1. Vitest: each stack item kind gives the strip's words in the History.
   2. e2e: a damage line in the full History names who and whom.
6. Result (2026-10-21)
   1. A resolved item without a card takes the strip's words (`stackLine`).
   2. Damage: "Bo deals 2 damage to Gaper"; a death: "Gaper dies".
   3. Penalty, death step and rolls the same; a card's item as before.
   4. A changed roll: "Attack: Bo rolls 6 → 4 against Fly" (19.1).
   5. `entriesFor` takes the update's view, to name cards as the strip does.
   6. These lines use names, not "You", as the strip does.
   7. Vitest: seven card-less items give the strip's words, none "Roll".
   8. e2e `info.spec.ts`: two damage lines name who and whom.
   9. The 19.1 test now reads the strip's roll words.
   10. Its first rerun failed on a too strict test pattern; the app was right.
   11. Evidence: `.cache/e2e-evidence/19.52-run1/` (not in git).
   12. Wiki `Client`; CHANGELOG.
   13. Screenshots: `.cache/screenshots/19.52/before/` and `after/`.
7. Owner notes (2026-10-21)
   1. Owner marks 19.52 solved.

### [ ] 19.53 e2e: the 3-player table test times out in the first roll

1. Goal: the 3-player layout test passes every run.
2. Report (agent, 2026-10-22)
   1. Found by `make e2e` during 15.92.
   2. `table.spec.ts`: "three players sit left and right; nothing covers the decks".
   3. Chromium; the 30 s test timeout ran out in `throwFirstRoll`.
   4. Passed on its rerun, alone, on both browsers.
   5. 15.92 does not touch the first roll or the layout.
   6. Evidence: `.cache/e2e-evidence/15.92-run1/`.
   7. Seen again in `make e2e` during 15.94, on webkit: `.cache/e2e-evidence/15.94-run1/`.
   8. Seen in `make screenshots` during 16.5, chromium, another test.
      1. "the board keeps its size when shop slots empty (V-23)".
      2. The first roll was a 6–6 tie; the reroll was still throwing.
      3. Evidence: `.cache/e2e-evidence/16.5-run1/`.
   9. Seen in `make e2e` during 16.7, webkit, the 2-player table test.
      1. "two players face each other; the table fits 1280x720".
      2. The 30 s timeout ran out in `settled`, after the first roll.
      3. 16.7 changes e2e text lookups only; this test is unchanged.
      4. Evidence: `.cache/e2e-evidence/16.7-run1/`.
   10. Seen in `make screenshots LOCALE=qps` during 16.7, chromium.
       1. "curses sit at the far left …" (15.55), 4 players.
       2. The 30 s timeout ran out in `throwFirstRoll`.
       3. Evidence: `.cache/e2e-evidence/16.7-shots-qps/`.
3. Tasks:
   1. Find why the first roll took too long in the full run.
   2. Fix the client or the test, whichever is wrong.
4. Done when:
   1. The test passes ten full runs in a row.

### [x] 19.54 e2e: a loot card's flight from under the stack, webkit

1. Goal: the 15.74 flight test passes every run.
2. Report (agent, 2026-10-22)
   1. Found by `make e2e` during 15.92.
   2. `playback.spec.ts`: "a resolved loot card flies to the discard from the large card under the stack".
   3. Webkit; line 926: not every box after the first hidden one was hidden.
   4. Passed on its rerun, alone, on both browsers.
   5. 15.92 does not touch the flights or the stack panel.
   6. Evidence: `.cache/e2e-evidence/15.92-run1/`.
   7. Seen again in `make e2e` during 15.94, on chromium, the "from its stack row" case.
   8. Same line 926; evidence: `.cache/e2e-evidence/15.94-run1/`.
3. Tasks:
   1. Find which box showed while it should hide, and why.
   2. Fix the client or the test, whichever is wrong.
4. Done when:
   1. The test passes ten full runs in a row.
5. Cause (2026-10-22)
   1. A client bug; the test was right.
   2. A probe logged the row's inline style each frame.
   3. The row hid by `visibility: hidden` in its style (19.41).
   4. Then the row slid 32 px up: a layout move of 15.51's `Rows`.
   5. motion rewrites a moving row's style; the `visibility` went with it.
   6. So the flown card's row showed again on the stack, sliding.
   7. Also: a leaving row fades out (15.51) after the store clears `left`.
   8. Evidence: `.cache/e2e-evidence/19.54-repeat1/` and the probe logs.
6. Result (2026-10-22)
   1. A flown row hides by the class `flown`, which motion never touches.
   2. It stays hidden while it fades out, though `left` is cleared.
   3. 12 repeats per case and browser: 48 of 48 passed (5 of 48 failed before).
   4. Not ten full runs: a full run takes 7 minutes; one full run passed.
   5. Grouped under E: playback, animations (the cause was in the client).
   6. That full run had two other flakes, both passed alone: logged as 19.55, 19.56.
   7. The next full run lost one other test, passed alone: logged as 19.57.
   8. The flight test passed in both full runs after the fix.
7. Owner notes (2026-10-22)
   1. Owner marks 19.54 solved.

### [ ] 19.55 e2e: whom the stack waits for, webkit, under load

1. Goal: the V-11 test passes every run.
2. Report (agent, 2026-10-22)
   1. Found by `make e2e` during the 19.54 fix.
   2. `stack.spec.ts`: "every player sees whom the stack waits for (V-11)".
   3. Webkit; line 150: the poll for "Passed" or an empty stack timed out.
   4. Passed 6 of 6 alone, on both browsers.
   5. Evidence: `.cache/e2e-evidence/19.54-run1/`.
3. Tasks:
   1. Find why the pass did not show in time in the full run.
   2. Fix the client or the test, whichever is wrong.
4. Done when:
   1. The test passes ten full runs in a row.

### [ ] 19.56 e2e: Chromium crashed during a full run

1. Goal: know whether the app can crash the browser.
2. Report (agent, 2026-10-22)
   1. Found by `make e2e` during the 19.54 fix.
   2. `actions.spec.ts`: "an attack waits for the roller to throw …" (15.61, 15.63).
   3. Chromium: "Target crashed" in `seen()`, a page evaluate.
   4. Passed 6 of 6 alone, on both browsers.
   5. Evidence: `.cache/e2e-evidence/19.54-run1/`.
3. Tasks:
   1. Check the run for memory use; see if it crashes again.
   2. If it repeats, find what the page did when it crashed.
4. Done when:
   1. No crash in ten full runs, or its cause is found.

### [ ] 19.57 e2e: an order dialog keeps a card unplaced, webkit, under load

1. Goal: the 19.16 order test passes every run.
2. Report (agent, 2026-10-22)
   1. Found by `make e2e` during the 19.54 fix.
   2. `actions.spec.ts`: "put them back in any order" shows each card in its place (19.16).
   3. Webkit; line 447: "Not placed" still had one card after 5 s.
   4. Passed 10 of 10 alone, on both browsers.
   5. Evidence: `.cache/e2e-evidence/19.54-run2/`.
   6. 19.53, 19.55 – 19.57: each full run loses one to four tests, each passes alone.
3. Tasks:
   1. Find why the last placement did not register in the full run.
   2. Check whether fewer e2e workers stop these load failures.
4. Done when:
   1. The test passes ten full runs in a row.

### [x] 19.58 e2e: a played "Choose one" loot card never reaches the stack

1. Goal: the 15.31 played-card test passes whatever loot is dealt.
2. Report (agent, 2026-10-22)
   1. Found by `make e2e` during 15.97.
   2. `fight.spec.ts`: "another player's loot card shows large, then by the stack (15.31)".
   3. Webkit; line 54: no `played-pop` on the other page after 10 s.
   4. Evidence: `.cache/e2e-evidence/15.97-fight-105217/`.
3. Tasks:
   1. Find what the active page waited for.
   2. Fix the client or the test, whichever is wrong.
4. Done when:
   1. The test answers a mode or a target before it checks the pop.
5. Cause
   1. The page snapshot: Yol's prompt reads "Choose one".
   2. The played card asks its mode before it goes on the stack.
   3. The test played the card, then waited; nobody picked the mode.
   4. The app was right: others see a loot card once it is on the stack.
   5. The record holds only its header: the server stopped before a flush.
6. Fix
   1. The test calls `answerPlay` after `playLoot`, as 19.18 does.
7. Owner notes (2026-10-22)
   1. Owner marks 19.58 solved.


### [x] 19.59 A build without the images submodule fails in `go mod tidy`

1. Goal: `wails build -tags embed` works in a clone without the images.
2. Report (owner, 2026-10-22)
   1. CI job `app`, linux/amd64, on `main` after the 0.2.0 merge.
   2. "Generating bindings" fails.
   3. `card_db@latest … does not contain package github.com/FourSoulsRepo/card_db/images`.
   4. The `for-poor` builds of the v0.2.0 release hit the same error.
3. Tasks:
   1. Find why a build without `cardimages` resolves that import.
   2. Fix every scripted build and the build docs.
4. Done when:
   1. A worktree without the submodule builds with `-tags embed`.
   2. CI job `app` passes.
5. Cause
   1. Since 8.29 `embed_images.go` imports the submodule's Go package.
   2. The file has `//go:build embed && cardimages`; `go build` skips it.
   3. `wails build` and `wails dev` run `go mod tidy` first.
   4. Tidy reads files of every build tag, so it needs that package.
   5. 8.29 was checked with the submodule present, so tidy passed.
   6. The 198 commits reached CI only with the release push.
6. Fix
   1. Every `wails build` passes `-m`: skip `go mod tidy`.
   2. The build action, `make internal`, and the build docs.
   3. Go has no way to hide a tagged import from tidy.

### [x] 19.60 Distant Admiration: an unplayable given card leaves a loot play

1. Goal: the free play of a given loot card never adds a loot play.
2. Report (owner, 2026-10-22)
   1. Play test; record `2026-10-22-163245-g2_Hitokiri-Battosai_Gaijin_distant_admiration_one_more_loot.fsrec`.
   2. The Fettered played two loot cards: the turn's play and the extra one.
   3. Distant Admiration took Butter Bean! from the other player.
   4. Butter Bean had no target; it stayed in hand.
   5. The player then had one more loot play.
3. Tasks:
   1. Find the cause in the match record.
   2. A free play spends no loot play, played or not.
   3. "If able": a card that cannot be played stays in hand.
4. Done when:
   1. Engine test: an unplayable given card leaves the loot plays as they were.
   2. Engine test: a played given card spends no loot play.
   3. Engine test: a cancelled free play leaves the plays as they were.
5. Cause
   1. Record steps 635 and 651: both loot plays used.
   2. Step 662: Cain gives Butter Bean!; step 663: `cancelled`.
   3. `ActPlayLootFree` added an extra loot play, then played the card.
   4. Butter Bean! had no target: its prompt offered only "cancel".
   5. The cancel kept the card in hand and the extra play stayed.
6. Fix
   1. `Game.FreeLoot` marks the card played for free.
   2. Its play spends no loot play (`spendLootPlay`).
   3. A cancel clears the mark.
   4. "If able" uses the play checks, loot plays aside (`lootPlayable`).
   5. An unplayable card stays in hand; nothing is asked.
7. Open question for the owner
   1. "You must play that loot card": should its prompts drop "cancel"?
   2. The engine part is small.
   3. The client needs care: an aim with no cancel must not end on a click.
   4. Today a cancel leaves the card in hand, at no loot play cost.
8. Owner notes (2026-10-22)
   1. Owner marks 19.60 solved.
   2. The open question in 7 is not answered yet.

### [x] 19.61 Boss Rush!: a boss put in its own slot leaves it covered

1. Goal: an event goes to discard after it resolves, even when covered.
2. Report (owner, 2026-10-22)
   1. Play test; record `2026-10-22-190326-g2_Hitokiri-Battosai_jеня_bossrush_bug.fsrec`.
   2. Boss Rush! came from the top of the monster deck.
   3. It revealed Rag Mega and The Duke of Flies.
   4. One boss went into Boss Rush's own slot.
   5. Boss Rush stayed under that boss instead of going to discard.
3. Tasks:
   1. Find the cause in the match record.
   2. A covered event still goes to discard when its abilities are done.
   3. The card that covered it stays in the slot.
4. Done when:
   1. Engine test: Boss Rush with a boss in its own slot goes to discard.
   2. Engine test: that boss stays in play; no covered card is left.
   3. Engine test: when that boss dies, no event comes back into the slot.
5. Cause
   1. Step 100: Boss Rush resolves; step 102: both bosses enter play.
   2. Its own slot counts as "not being attacked", so it can be picked.
   3. `coverMonsterSlot` marks Boss Rush covered.
   4. `finishEvent` acts only on an event in play; a covered one is skipped.
   5. Step 515: Rag Mega dies and Boss Rush is uncovered.
   6. It is back in the slot as an event with no trigger and no discard.
   7. The game ended on that step; later it would block the slot.
6. Fix
   1. `finishEvent` also takes a covered event out of its slot.
   2. The card on top stays; nothing is uncovered.
   3. The event goes to the monster discard, as R-CARD-15 says.
   4. Test `TestBossRushOwnSlot` covers the three done-when checks.
   5. Without the fix it fails: slot 1 is `[fatty boss_rush monstro]`.
7. Owner notes (2026-10-23)
   1. Owner marks 19.61 solved.

### [x] 19.62 Emergency Contact does nothing on an attack on the monster deck

1. Goal: Emergency Contact works as the owner's ruling says.
2. Report (owner, 2026-10-23)
   1. Play test; record `2026-10-23-143924-g1_Hitokiri-Battosai_jеня_emergency_curved_horn.fsrec`.
   2. Bethany attacked; The Fallen came out.
   3. She played Emergency Contact; it did nothing.
   4. The attack went on, and The Fallen killed her.
3. Tasks:
   1. Find the cause in the match record.
   2. Get the owner's ruling on deck attacks.
   3. Apply it to every card that checks the same thing.
4. Done when:
   1. Engine test: Emergency Contact cancels a deck attack on a revealed monster.
   2. That monster goes to the bottom of the monster deck.
   3. Engine test: IBS cancels such an attack too.
5. Cause
   1. Step 258: Bethany attacks the top of the monster deck.
   2. Step 259: The Fallen is revealed and becomes the target.
   3. Step 268: Emergency Contact resolves with no effect.
   4. Its code skips an attack with `Attack.Deck` set (R-ATK-05).
   5. Steps 274–278: the roll of 2 misses; she dies.
   6. The engine works as coded; the rule reading is the question.
6. Rule reading
   1. Rulebook, Attacking: an ability that looks for an attack declared on a monster.
   2. It does not interact with an attack declared on the monster deck.
   3. That holds even when a monster is revealed.
   4. Our R-ATK-05 drops "declared": any ability about an attack on a monster.
   5. Emergency Contact says "Cancel an attack on a monster": no "declared".
   6. IBS uses the same check: "Cancel your attack on a monster".
   7. No ruling found for either card: Yuggy, Reddit, Discord, X.
7. Open question for the owner
   1. A: keep R-ATK-05; both cards do nothing on a deck attack.
   2. B: narrow R-ATK-05 to "declared"; both cards cancel the attack.
   3. Under B, Emergency Contact puts the revealed monster on the deck bottom.
   4. Either answer is House; it goes to "Rulings to revisit".
8. Owner ruling (2026-10-26)
   1. Option B: both cards work once a monster is revealed.
   2. Logged as Q-24 in `open-questions.md`; basis House.
9. Fix
   1. R-ATK-05 now names abilities about an attack *declared* on a monster.
   2. Emergency Contact and IBS no longer check `Attack.Deck`.
   3. Tests `TestEmergencyContactDeckAttack`, `TestIBSCancelsDeckAttack`.
   4. Without the fix both fail: the attack goes on.
10. Owner notes (2026-10-26)
   1. Owner marks 19.62 solved.

### [x] 19.63 Soul of Pride counts a buy from the treasure deck

1. Goal: Soul of Pride counts treasure gains as the ruling says.
2. Report (owner, 2026-10-26)
   1. Play test; record `2026-10-24-124502-g2_Hitokiri-Battosai_jеня_soul_of_pride.fsrec`.
   2. A buy from the top of the treasure deck put a counter on it.
   3. The owner expected buys not to count.
3. Tasks:
   1. Find the cause in the match record.
   2. Check the rule the engine follows.
4. Done when:
   1. The cause is found; a fix only if the engine breaks a rule.
5. Cause
   1. Not a bug: the engine follows R-CARD-37.
   2. R-CARD-37 comes from Yuggy's ruling (S-YUG: Damocles; Soul of Pride).
   3. A buy from the treasure deck top gains a treasure.
   4. A buy from a shop slot gains an item, not a treasure.
   5. Record: shop buys (Champion Belt, Plan C) add no counter.
   6. Five deck-top buys and the Bum-bo gain add one each.
   7. The 6th counter gives the soul to the active player.
   8. `TestSoulOfPride` already covers the deck-top buy.
6. Fix
   1. None; no code change.
   2. "Bugs by system" gets group L for reports that are not bugs.
7. Owner notes (2026-10-26)
   1. Owner marks 19.63 solved: works as ruled.

### [x] 19.64 Hourglass deactivates a passive item

1. Goal: "deactivate an item" picks only items with a ↷ ability.
2. Report (owner, 2026-10-26)
   1. Play test; record `2026-10-24-135401-g3_Hitokiri-Battosai_jеня_dads_lost_coin.fsrec`.
   2. Hourglass deactivated Dad's Lost Coin, a passive item.
   3. The owner expects only charged active items to be offered.
3. Tasks:
   1. Find the cause in the match record.
   2. Look for a ruling on deactivating passive items.
   3. Apply it to every card that deactivates items.
4. Done when:
   1. Engine test: Hourglass offers no passive item.
   2. Engine test: Curse of Impulse leaves passive items charged.
5. Cause
   1. Step 1116: Hourglass offers 11 items; the player picks Dad's Lost Coin.
   2. Every Hourglass prompt lists passive items, Hourglass too.
   3. `chargedItems` takes every charged item; passive items are always charged.
   4. Curse of Impulse also deactivates every item, passive ones too.
   5. The death penalty already skips them: `CardDef.Tap` (R-DEATH-14).
6. Rule reading
   1. R-ABIL-10: a ↷ ability costs deactivating its object.
   2. R-DEATH-14 deactivates only what has a ↷ ability.
   3. No rule says a passive item can't be deactivated.
   4. No ruling found: rulebook, Yuggy, Reddit (2026-10-26).
   5. The owner's reading is House: Q-25 in `open-questions.md`.
7. Fix
   1. `chargedItems` offers only charged items with a ↷ ability (`CardDef.Tap`).
   2. Curse of Impulse deactivates only such items.
   3. Rule R-MECH-54 states it.
   4. Tests `TestHourglassSkipsPassiveItems`, `TestCurseOfImpulseSkipsPassiveItems`.
   5. Without the fix both fail: Hourglass offers itself and the coin.
8. Found on the way, not fixed
   1. The Revenant (Anima Sola's back) has `Tap` set but no ↷ ability.
   2. The death penalty deactivates it; Hourglass now offers it.
   3. Left for the owner: it may need its own item.
   4. Logged as 19.65.
9. Owner notes (2026-10-26)
   1. Owner marks 19.64 solved.

### [x] 19.65 The Revenant counts as having a ↷ ability

1. Goal: The Revenant side is never deactivated as a ↷ item.
2. Report (owner, 2026-10-26)
   1. Found while fixing 19.64; no play test.
   2. The owner logs it as a bug and keeps it for later.
3. Tasks:
   1. Check the flip and charge rules for two-sided items (R-MECH-53).
   2. Decide the back face's `Tap` flag.
4. Done when:
   1. Engine test: Hourglass does not offer The Revenant.
   2. Engine test: the death penalty leaves The Revenant as it was.
5. Cause
   1. `theRevenant` in `cards/r/anima_sola.go` sets `Tap: true`.
   2. The Revenant side has no ↷ ability; only Anima Sola has one.
   3. The death penalty deactivates it (R-DEATH-14).
   4. Hourglass offers it since 19.64 (R-MECH-54).
6. Rule reading
   1. R-MECH-53: a flipped object keeps its charge state.
   2. Each face has its own abilities; The Revenant has no ↷ one.
   3. So `Tap` follows the face up, like the abilities.
   4. The recharge step still charges it; it does not read `Tap`.
7. Fix
   1. `theRevenant` loses `Tap: true`.
   2. Tests `TestHourglassSkipsTheRevenant`, `TestDeathPenaltySkipsTheRevenant`.
   3. Without the fix both fail: Hourglass offers it; the penalty deactivates it.
   4. `TestTapMatchesAbilities`: a treasure's `Tap` matches a ↷ cost.
   5. Placebo counts: it copies a ↷ ability.
   6. No other card was wrong.
8. Owner notes (2026-10-26)
   1. Owner marks 19.65 solved.

### [x] 19.66 Dad's Lost Coin misses a reroll that is a 1 again

1. Goal: a changed or rerolled roll meets "would roll N" triggers again.
2. Report (owner, 2026-10-26)
   1. Play test; record `2026-10-24-135401-g3_Hitokiri-Battosai_jеня_dads_lost_coin.fsrec`.
   2. jеня rolled a 1; the owner forced a reroll with Dad's Lost Coin.
   3. The reroll was a 1 again; the coin did not trigger.
3. Tasks:
   1. Find the cause in the match record.
   2. Fix every helper that changes a roll, not only the reroll.
4. Done when:
   1. Engine test: the coin triggers on a second 1 from its own reroll.
   2. Engine test: a roll changed and changed back is checked again.
   3. Engine test: a roll set to its own value is not checked again.
5. Cause
   1. Step 969: attack roll 1 (stack 229); step 971: the coin triggers.
   2. Step 975: the reroll is a 1; step 977: it resolves, no check.
   3. `resolveTop` checks a roll only if its value differs from `Checked`.
   4. `Reroll`, `SetRoll` and `ModifyRoll` leave `Checked` as it was.
   5. So a roll back on a checked value skips "would roll" triggers.
6. Rule reading
   1. The card says "each time": every 1 counts.
   2. R-DICE-05: a change sends the roll back; a reroll is one.
   3. The owner: a change to the same value is no change (Q-26, House).
   4. Rewards are "when you roll" triggers: they fire once (R-DICE-06).
7. Fix
   1. `StackItem.setRoll` changes a roll; a new value clears `Checked`.
   2. `SetRoll` and `ModifyRoll` use it; `Reroll` clears `Checked`.
   3. Rule R-DICE-13 states it.
   4. Tests `TestDadsLostCoinTwice`, `TestRerollToTheSameValueIsCheckedAgain`.
   5. Tests `TestRollChangedBackIsCheckedAgain`, `TestRollSetToItsValueIsNotCheckedAgain`.
   6. Without the fix the first three fail.

### [x] 19.67 Curved Horn's +1 ATK shows only after the first attack roll

1. Goal: a "first attack roll" ATK bonus shows while it is still to come.
2. Report (owner, 2026-10-26)
   1. Play test; no record named.
   2. Curved Horn's +1 ATK shows only once the attack roll resolved.
   3. The owner expects it before the attack is declared.
   4. It should go away once that roll's damage is dealt, or on a miss.
3. Tasks:
   1. Find the cause in the view.
   2. Apply it to every "first attack roll" ATK card.
4. Done when:
   1. Engine test: the view shows the +1 before the attack.
   2. Engine test: it shows while the first roll's damage is on the stack.
   3. Engine test: it is gone after that damage, and after a miss.
5. Cause
   1. `FirstAttackRollATK` applies only while `AttackRolls == 1` in an attack.
   2. `AttackRolls` counts resolved attack rolls; it is 0 before the roll.
   3. So the view showed the +1 only after the first roll resolved.
   4. It stayed after a miss, and after the damage, until the next roll.
   5. The rules are right: the +1 counts only for that roll's damage.
6. Fix
   1. View only; the engine rules and its version stay.
   2. `atkGame()` (`view.go`): the game the shown ATK and its sources read.
   3. With the active player's combat damage on the stack: the game itself.
   4. Else a shallow copy with the next attack roll made, as the odds do.
   5. Covers Curved Horn, Champion Belt, Polydactyly (`FirstAttackRollATK`).
   6. Tests `TestCurvedHornShownForTheFirstRoll`, `TestCurvedHornGoneAfterAMiss`.
   7. Without the fix both fail.
   8. `TestChampionBeltSourceDuringFirstRoll` now expects the belt before the attack.
   9. UC-12 3j; wiki `Rules-engine`.
7. Owner notes (2026-10-26)
   1. Owner marks 19.67 solved.

### [ ] 19.68 LAN test fails once: UDP address already in use

1. Goal: `TestLANFindsGamesUntilTheyEnd` never fails on a busy port.
2. Report (Claude, 2026-10-26)
   1. Seen in the 17.3 gate, `make test`; no record, it is a Go test.
   2. The test failed once: UDP "address already in use".
   3. A rerun passed; 17.3 does not touch LAN code.
3. Tasks:
   1. Find which port clashes and why.
   2. Suspect: `freePort` checks a TCP port; discovery reuses it for UDP.
   3. Check other tests that bind UDP in parallel.
   4. Group M is provisional; move the item if the cause is elsewhere.
4. Done when:
   1. The cause is named here.
   2. The test runs 200 times in a row without that failure.

### [x] 19.69 Godmode ↶ on a played card's row leaves the card on the stack

1. Goal: ↶ on a played card's row undoes the play itself.
2. Report (Claude, 2026-10-26)
   1. Found by the 17.6 agent while testing A Penny! in a Soft game.
   2. A loot play or an activation has no history row (15.88).
   3. The card's row comes with its resolution, an update that starts with a pass.
   4. In Godmode ↶ on that row goes back to before the pass.
   5. So the card stays on the stack; the play is not undone.
   6. In Soft no button shows for a played card at all.
3. Tasks:
   1. Fixed by 17.6: the play gets its own history row with the button.
   2. Check that no other target step lacks a row of its own.
4. Done when:
   1. e2e: Godmode ↶ on the play row undoes the play.
5. Cause
   1. `historyKinds` skipped `loot_played` and `activated` (15.88).
   2. A row's ↶ targets the first step of the update that brought it.
   3. The play's own update had no row, so no row carried its step.
   4. The result's row came with the pass that resolved it.
6. Fix
   1. Client only: `loot_played` and `activated` give a line (17.6).
   2. "Bo plays A Penny!", "Bo uses The D6"; no headline, the strip is unchanged.
   3. Its row carries the play's step; ↶ there undoes the play.
   4. Vitest: `a play has its own row (17.6)` in `history.test.ts`.
   5. e2e: `in Godmode ↶ on a played card undoes the play` in `rollback.spec.ts`.
   6. The result's row keeps its Godmode ↶: it undoes only the resolution.
   7. Task 3.2: other target steps still have no row of their own.
   8. A play or activation that asks a mode or target: the intent itself.
   9. Its row comes with the last answer; ↶ there goes back to that prompt.
   10. The prompt offers "cancel", so the play can still be dropped.
   11. Also: some answers to a card's question, an order, an end-turn declaration.
   12. Such as Bum Friend's put-back pick (`moved_to_deck`, ADR 020 item 4.22).
   13. Left as they are; a later item may give them rows.
   14. UC-15 2e, UC-30 3g, UC-31 3e; 15.88 item 7; wiki `Client`; CHANGELOG.
   15. Screenshots: `.cache/screenshots/19.69/before/` and `after/`.
7. Changed later (17.7, 2026-10-26)
   1. The row fix is replaced by the one Undo button.
   2. The play rows are gone; Undo names the play in its tooltip.
   3. One press undoes the play, also after it resolved: a pass is no target.
   4. A play's answers to a mode or target are no target of their own.
   5. So 6.8 – 6.10 no longer apply: one press drops the play.
8. Owner notes (2026-10-26)
   1. Owner marks 19.69 solved, with step 17.

### [x] 19.70 A resolved A Penny! shows no "gains 1¢" row for 5 seconds

1. Goal: a resolved card's result row always reaches every client.
2. Report (Claude, 2026-10-26)
   1. Seen once by the 17.6 agent in a dev e2e run, Soft mode.
   2. Tom's page showed "A Penny!" resolved, but no "Uma gains 1¢" row.
   3. It showed "Waiting for Uma, You" with an empty stack for 5 seconds.
   4. Not seen again in later runs or in either gate run.
   5. Evidence: `.cache/e2e-evidence/19.70/e2e-dev2/` (not tracked).
   6. It holds the match records (`cache-e2e`) and `test-results`.
3. Tasks:
   1. Read the saved match record and test results.
   2. Find whether the row was never sent, or sent and not drawn.
   3. Group J is provisional; move the item once the cause is found.
4. Done when:
   1. The cause is named here, with a test if it is a real bug.
5. Cause
   1. Record `g2_Tom_Uma`; its last gzip part is cut off, so plain `gunzip` stops early.
   2. Every part read: steps 10 – 12 are the Penny game after the load.
   3. Uma plays Eden; her starting item is Bum-Bo! (step 5).
   4. Bum-Bo!: "If you would gain any amount of ¢, this levels up by that much instead."
   5. Step 12: A Penny! resolves; events `replaced` and `counters` on Bum-Bo!, no `gained_cents`.
   6. So Uma gained no 1¢; no "Uma gains 1¢" row could come.
   7. Tom's page agrees: Uma has 3¢, Bum-Bo! shows a counter of 1.
   8. Its history shows "A Penny!" resolved; the stack is empty.
   9. "Waiting for Uma, You" is Uma's own action time, not a stuck pass.
   10. The engine and client were right; the probe assumed a gain.
   11. A dealt card that replaces the gain breaks that check (CLAUDE.md 5.16).
   12. Record `g1_Vic_Wes` shows the plain case: `gained_cents` 1.
6. Result
   1. Not an app bug: group I, moved from J.
   2. The probe is gone: 17.7 rewrote the test.
   3. Today it checks the resolved row `✦A Penny!`, which shows with Bum-Bo! too.
   4. No other e2e test expects a real cent gain from a dealt card.
   5. No code change; no guard test, since no gap remains.
7. Owner notes (2026-10-26)
   1. Owner marks 19.70 solved; no history line for Bum-Bo! for now.

### [x] 19.71 XXI. The World: the look dialog does not say whose hand

1. Goal: a look at hands shows whose hand each card is from.
2. Report (owner, 2026-10-26)
   1. XXI. The World: "Look at each player's hand, then loot 2."
   2. The "You look at" dialog shows every card in one row.
   3. The player cannot tell which cards are whose.
3. Tasks:
   1. Engine: a `looked_at` event of a hand card names the hand's owner.
   2. It uses `By`: the owner's seat plus 1; 0 for a deck card.
   3. Client: `Looked` carries that owner.
   4. The dialog groups hand cards by owner, in event order.
   5. Each group has a label with the owner's name: "Bo's hand".
   6. A look at deck cards stays as it is, one row, no label.
   7. Other hand looks get the labels too: Mom's Eye, X-Ray Vision, Incubus.
   8. The history line stays as it is.
4. Done when:
   1. Go test: The World's events name each hand's owner.
   2. Vitest: `lookedAt` keeps the owner.
   3. e2e: the dialog shows two labelled groups, in `actions.spec.ts`.
   4. UC-18 gets the case; CHANGELOG line.
   5. Screenshots: `.cache/screenshots/19.71/before/` and `after/`.
5. Cause
   1. `Game.LookAt` emits only the player, the object and the card.
   2. The looker's client gets no other hand's object IDs.
   3. So it cannot tell where a card came from.
6. Fix
   1. Engine: `LookAt` sets `By` on a hand card's `looked_at` event.
   2. `By` is the hand's owner plus 1; a deck card keeps 0.
   3. So Mom's Eye, X-Ray Vision, Incubus, Dead Bird and Joker get it too.
   4. Client: `lookedAt` keeps it as `Looked.owner`.
   5. `lookGroups` in `history.ts` splits the cards by owner, in event order.
   6. `LookDialog` shows one row per hand under "Bo's hand"; "Your hand" for the looker's own.
   7. Deck cards stay one row with no label; the history line is unchanged.
   8. Others still see only `by` on the event; the card stays hidden.
   9. Go test: `TestXXITheWorldNamesEachHand` in `xxi_the_world_test.go`.
   10. Vitest: `keeps whose hand a card is from; a deck card has no owner (19.71)`.
   11. Vitest: `groups hand cards by owner in event order; deck cards in one group (19.71)`.
   12. e2e: `a look at hands groups the cards by whose hand they are (19.71)` in `actions.spec.ts`.
   13. UC-18 4c; wiki `Client`; CHANGELOG.
   14. Screenshots: `.cache/screenshots/19.71/before/` and `after/`.
   15. A revealed given card (fsp2 "Reveal it") is in the taker's hand.
   16. So its look now reads "Bo's hand", or "Your hand" for the taker.
7. Owner notes (2026-10-26)
   1. Owner marks 19.71 solved.

### [x] 19.72 A killed monster's covered card flies in from the monster deck

1. Goal: an uncovered card shows up in its slot, not from the deck.
2. Report (owner, 2026-10-26)
   1. A monster slot holds 2 cards: the top one covers another.
   2. The top monster dies.
   3. The card under it flies from the monster deck to the slot.
   4. Wrong: that card was in the slot already, only covered.
3. Tasks:
   1. Engine: `uncover` marks its `card_revealed` with `Prev`.
   2. `Prev` is the covered object; the card gets a new ID (R-ZONE-01).
   3. That fits `Prev`'s meaning: the object moved away under a new ID.
   4. Client: a `card_revealed` with `prev` gives no `reveal`, no `place`.
   5. The card was face up before; it needs no large show either.
   6. It appears in place with `card-pop`, like any new table card.
   7. Rooms too: a Red Key room leaving uncovers the one below.
   8. A deck attack, a refill and setup keep their `place` flight.
4. Done when:
   1. Go test: killing a covering monster emits `card_revealed` with `Prev`.
   2. Vitest: a `card_revealed` with `prev` makes no `place`, no `reveal`.
   3. e2e: the uncovered card does not fly from the deck, in `playback.spec.ts`.
   4. Wiki `Animations` row `fly`: `place` skips an uncovered card.
   5. ADR or wiki that lists event fields: `Prev` on `card_revealed`.
   6. CHANGELOG line.
5. Cause
   1. `uncover` (`combat.go`) puts the covered card back via `putInSlot`.
   2. `putInSlot` emits a plain `card_revealed`, as for a deck card.
   3. `playback.ts` makes a `place` from `deck-<kind>` for every slot `card_revealed`.
   4. So the client cannot tell an uncover from a refill.
6. Fix
   1. Engine: `putInSlotFrom` is `putInSlot` with a `Prev` on its `card_revealed`.
   2. `uncover` calls it with the covered object; other callers keep 0.
   3. So every uncover sets `Prev`: a dead monster, a room, a shop card.
   4. Client: `card_revealed` with `prev` gives no `reveal`, no `place`.
   5. The view draws the card; it pops in place (`card-pop`).
   6. A deck attack, a refill and setup keep their flight.
   7. Go test: `TestUncoverNamesTheCoveredObject` in `combat_test.go`.
   8. Vitest: `neither reveals nor places a card a dead monster uncovered (19.72)`.
   9. e2e: `a killed monster's covered card shows in its slot, not from the deck (19.72)`.
   10. It is in `playback.spec.ts`; effects played: only `deaths`.
   11. UC-12 1aa; wiki `Animations` row `fly`, `Client`, `Rules-engine`; CHANGELOG.
   12. Screenshots: `.cache/screenshots/19.72/before/` and `after/`.
   13. Before: the uncovered card shows large over the board, then flies.
   14. After: the skull, then the card is in its slot.
7. Owner notes (2026-10-26)
   1. Owner marks 19.72 solved.

### [x] 19.73 A dead monster pops into the discard or soul slot; no flight

1. Goal: a dead monster flies from its slot to where it goes.
2. Report (owner, 2026-10-26)
   1. A monster dies.
   2. It does not fly to the monster discard, nor to the player's souls.
   3. It just appears on the discard pile.
3. Tasks:
   1. Engine: a dead card's discard emits `discarded`.
   2. `ActDiscardObject` emits it: `Player` none, `Prev` the held object.
   3. Engine: `GainSoul`'s `gained_soul` sets `Prev` to the moved object.
   4. Client: on `died`, playback keeps where the card died.
   5. That is its slot zone, `monster-slot-N`, keyed by the held object.
   6. It lasts across updates: reward rolls may come in between.
   7. `discarded` with that `prev`: `fly` from that zone to `discard-monster`.
   8. `gained_soul` with that `prev`: `fly` from that zone to the player.
   9. Any other `gained_soul` flies from `obj-<prev>` if drawn; else none.
   10. A bonus soul keeps no flight; it is not a table card.
   11. The fly starts after the `death-skull` on that card.
   12. History stays as it is: a `discarded` with no player adds no line.
4. Done when:
   1. Go test: a monster with no soul emits `discarded` with `Prev`.
   2. Go test: a killed soul monster's `gained_soul` has `Prev`.
   3. Vitest: both make a `fly` from the slot, also in a later update.
   4. e2e in `playback.spec.ts`: the dead card flies to the discard.
   5. e2e: a soul monster flies to its player.
   6. Wiki `Animations` row `fly`; wiki `Rules-engine` event fields.
   7. UC-12 gets the case; CHANGELOG line.
   8. Screenshots: `.cache/screenshots/19.73/before/` and `after/`.
5. Cause
   1. `monsterDeath` moves the card to a held object outside the table.
   2. `soulStep` later discards it or makes it a soul.
   3. `ActDiscardObject` calls `discard`, which emits no event.
   4. `GainSoul` emits `gained_soul` without `Prev`.
   5. `playback.ts` has no effect for `gained_soul` at all.
   6. So the next view just draws the card in its new place.
6. Fix
   1. Engine: `ActDiscardObject` emits `discarded`, no player, `Prev` the held object.
   2. Its callers are the soul step and a soul that cannot be gained.
   3. Engine: `GainSoul`'s `gained_soul` sets `Prev` to the moved object.
   4. No trigger listens to `discarded`; `gained_soul` triggers keep working.
   5. Client: `died` keeps its slot zone by the held object (`diedAt`).
   6. The zone comes from the view on screen, before the death.
   7. `diedAt` lives in the store; a rollback or leaving clears it.
   8. `discarded` with that `prev`: `fly` from the slot to the discard.
   9. `gained_soul` with that `prev`: `fly` from the slot to the player.
   10. Other `gained_soul` with `prev`: from `obj-<prev>`; a bonus soul: none.
   11. The fly comes after `deaths` in the event order.
   12. Added: the slot shows empty while the card flies (`lifted`).
   13. History unchanged: no line for a `discarded` with no player.
   14. Go tests: `TestDeadMonsterDiscardNamesItsHeldObject`, `TestKilledSoulNamesItsHeldObject`.
   15. Vitest: `flies a dead monster from its slot to the discard or its player, also in a later update (19.73)`.
   16. Vitest: `flies another soul from its place; a bonus soul not at all (19.73)`.
   17. e2e: `a dead monster flies from its slot to the monster discard (19.73)`.
   18. e2e: `a dead soul monster flies from its slot to its player after a reward roll (19.73)`.
   19. UC-12 1ab; wiki `Animations` row `fly`, `Client`, `Rules-engine`; CHANGELOG.
   20. Screenshots: `.cache/screenshots/19.73/before/` and `after/`.
   21. Before: the skull, then the card is in the discard; no flight.
   22. After: the skull, then the card flies from its emptied slot.
7. Owner notes (2026-10-26)
   1. Owner marks 19.73 solved, with the emptied slot during the flight.

### [x] 19.74 A declared end of turn has no row on the stack

1. Goal: the stack panel shows a declared end of turn, like an attack or a buy.
2. Report (owner, 2026-10-26)
   1. An open attack or buy has a row on top of the stack (15.91).
   2. A declared end of turn has none.
   3. The others cannot see why priority came to them.
3. Tasks:
   1. Client: `declaration` in `stack.ts` returns an end row.
   2. Text: `<player> ends the turn`; the name in their color.
   3. Shown while `turn.end_declared` and `turn.step` is the action step.
   4. Gone once the end phase starts (R-TURN-10).
   5. A turn ended early by a card shows it too (`EndDeclared`).
   6. A death end (`death_end`) gets no row; it is not declared.
   7. No fight line; same `stack-declared` row and look.
   8. No engine change: the view sends `end_declared` already.
4. Done when:
   1. Vitest: the row while declared; none in the end phase; none for a death end.
   2. e2e in `stack.spec.ts`: every player sees the row after End turn.
   3. UC-14 gets the case; wiki `Client`; CHANGELOG line.
   4. Screenshots: `.cache/screenshots/19.74/before/` and `after/`.
5. Cause
   1. `declaration` checks only `view.attack` and `view.buying`.
   2. `turn.end_declared` is used only by "Wait for my turn".
6. Fix
   1. Client: `declaration` returns `<player> ends the turn` after an attack or a buy.
   2. Shown while `turn.end_declared` and `turn.step` is the action step.
   3. So a card's early end shows too; the end phase drops it.
   4. A death end sets no `end_declared`; it has no row.
   5. Same `stack-declared` row; no fight line; no engine change.
   6. The row shows only while someone may respond.
   7. The server passes for a player who cannot (N-05); then no update shows it.
   8. Vitest: `says a declared end of turn until the end phase; not a death end (19.74)`.
   9. e2e: `a declared end of turn shows on top of the stack (19.74)` in `stack.spec.ts`.
   10. It ends turns until the other player may respond, up to six.
   11. UC-14 2q; wiki `Client`; CHANGELOG.
   12. Screenshots: `.cache/screenshots/19.74/before/` and `after/`.
   13. Before: the stack says "Nothing on the stack" while the viewer may respond.
   14. After: the row "Fay ends the turn" on top, the name in Fay's color.
7. Owner notes (2026-10-26)
   1. Owner marks 19.74 solved.

### [ ] 19.75 e2e: the 15.74 loot rise test timed out once

1. Goal: find why the test ran out of time, or show it was load.
2. Report (agent, 2026-10-26, during 19.74's gate)
   1. Full `make e2e`: 1 failed, 6 skipped, 349 passed.
   2. Test: `playback.spec.ts:1015`, Chromium.
   3. "another player's loot card rises from their hand and shrinks into the stack (15.74)".
   4. 30 s test timeout, hit at `waitForTimeout(1250)` on line 1030.
   5. The steps before it had used up the 30 s.
   6. Alone, it passed 6 of 6 (`--repeat-each 3`, both browsers).
   7. The next full `make e2e` passed: 350 passed, 6 skipped.
   8. Evidence: `.cache/e2e-fails/20261026-19.74-gate/` (not in git).
3. Tasks:
   1. Find which step before line 1030 took the time.
   2. Use the saved log, test results and match records.
   3. A real slow step is a bug; fix that.
   4. Else group J: under load, cause not found.
4. Done when:
   1. The cause is written here, or the item says it was not found.
   2. A fix, if any, comes with its test.

### [x] 19.76 End turn glows as a green box, not a breathing glow

1. Goal: with nothing left to do, End turn breathes a soft green glow (H-20).
2. Report (owner, 2026-10-26)
   1. Nothing left to do: the End turn button shows a green box.
   2. Expected: a breathing green glow, like the playable hand cards.
3. Tasks:
   1. Find why the glow draws as a box.
   2. Client: draw the glow as a soft halo that breathes.
   3. Keep the `end-turn-glow` tag and its row in `Animations.md`.
4. Done when:
   1. The idle button shows a soft glow in a screenshot.
   2. An e2e test checks the idle glow.
   3. Screenshots: `.cache/screenshots/19.76/before/` and `after/`.
5. Cause
   1. `.end-turn.idle` drew a hard ring: a shadow with only a spread.
   2. It sat on a dark green border; the ring grew and shrank.
   3. On the lumpy button it read as a green frame, not a glow.
6. First fix (replaced in 8)
   1. Client: the idle look follows the hand-card glow (`hand-glow`).
   2. The button keeps its ink outline and drop shadow.
   3. A tight soft green halo on the button itself.
   4. A wide aura in `::before` behind it; only its opacity breathes.
   5. Pressed: the aura hides; reduced motion: it stays still.
   6. e2e: `idle End turn glows softly and breathes, no ring (H-20, 19.76)`.
   7. It fails on the old style: the ring shadow is found.
   8. Wiki `Animations`, `Client`, `Screenshots`; CHANGELOG.
   9. Screenshots: `.cache/screenshots/19.76/before/` and `after/`.
7. Owner notes (2026-10-26)
   1. In the app the glow is still a green box.
   2. Its edges stop hard, left and right.
8. Cause, second look
   1. Any blurred shadow of the tilted button stops in a box in the app.
   2. The app runs macOS WebKit; Playwright's WebKit draws it right.
   3. So the glow cannot be a shadow of the button or its children.
9. Fix
   1. The button sits in a slot, `.end-turn-slot` in `TurnActions.tsx`.
   2. The idle slot draws the halo in `::before`, behind the button.
   3. A radial gradient in its own box, like Throw's glow; no shadow.
   4. It cannot sit in the button: the tilted button paints it on top.
   5. Only its opacity breathes (`end-turn-glow`).
   6. The e2e test now checks: no green in the button's shadow.
   7. And a radial-gradient halo that breathes.
   8. It fails on the first fix.
   9. Screenshots: `.cache/screenshots/19.76/again/`.
10. Owner notes (2026-10-26)
   1. Owner marks 19.76 solved.

### [ ] 19.77 e2e: the 19.73 death fly test misses the fly in a screenshot run

1. Goal: the dead monster fly test passes in `make screenshots` too.
2. Report (agent, 2026-10-27)
   1. Second `make screenshots` gate run for 16.5, chromium.
   2. `playback.spec.ts:860`: `effect-fly` never became visible in 5 s.
   3. The skull showed; then `shot(ari, 'death-skull')` saved a PNG.
   4. The monster was Holy Mom's Eye; it died and reached the discard.
   5. `make e2e` passed just before, both browsers.
   6. 16.5 changes card names only; nothing in playback.
   7. Evidence: `.cache/e2e-evidence/16.5-run2/` (not in git).
3. Cause
   1. Not found yet.
   2. Guess: as 19.43, the short fly ends while the screenshot is taken.
   3. The test then polls with `toBeVisible`; `trackFlight` already records it.
4. Tasks:
   1. Find the cause in the evidence; fix the test or the app.
5. Done when:
   1. Cause written here; three `make screenshots` runs in a row pass.

### [x] 19.78 e2e: the Meat! roll test sees no thrown roll, webkit, full run

1. Goal: the 15.94 stats test passes every run.
2. Report (agent, 2026-10-27)
   1. Found by `make e2e` during 16.6.
   2. `stats.spec.ts:149`: "the open attack says what the roll needs; Meat! …".
   3. Webkit; after the Throw click the stack still said "Throw the dice".
   4. No "(+1 Meat! = N)" line came in 5 s.
   5. Chromium passed it in the same run.
   6. 16.6 changes card image paths only; nothing in rolls or the stack.
   7. Evidence: `.cache/e2e-evidence/16.6-run1/` (not in git).
   8. Passed 3 of 3 alone, one run at a time, webkit.
   9. `--repeat-each` fails it: the saved games pile up on one server.
3. Cause
   1. The same as 19.86: a throw lost with a dropped connection.
   2. The webkit record ends at the attack roll's `stack_added`; no throw.
   3. The stack is the same: Pass disabled, "Throw the dice".
   4. Client: state from the old connection is kept (group O).
4. Tasks:
   1. Find the cause in the evidence; fix the test or the app.
5. Done when:
   1. Cause written here; the test passes ten full runs in a row.
6. Fix
   1. 19.86: a `welcome` clears `sent`.
7. Owner notes (2026-10-27)
   1. Owner marks 19.78 solved, with 19.86.

### [x] 19.79 e2e: playLoot presses the main menu's Play line

1. Goal: e2e helpers find each text by the key the app shows.
2. Report (agent, 2026-10-27)
   1. Found in 16.7 by `make screenshots LOCALE=pl` with a trial file.
   2. Nine tests timed out waiting for the action menu's Play item.
   3. English and `qps` hid it: both keys read "Play" there.
3. Cause
   1. `playLoot` in `e2e/tests/turn.ts` used `tr('menu.play')`.
   2. The card's action menu shows `actions.play` (`actions.ts`).
   3. The app was right; group I.
4. Tasks:
   1. `playLoot` uses `tr('actions.play')`.
5. Done when:
   1. The nine tests pass under `LOCALE=pl`; `make e2e` still passes.

### [x] 19.80 e2e: the two 15.93 icon tests find each other's saved game

1. Goal: each saved-game test continues only its own save.
2. Report (agent, 2026-10-27)
   1. Found by `make e2e` during 16.7, webkit, both variants at once.
   2. `stats.spec.ts:77`: "effect icons go at the end of the turn, animations on/off".
   3. `saveAndContinue` failed: two "Eli, Fox" rows had a Continue button.
   4. Evidence: `.cache/e2e-evidence/16.7-run2/` (not in git).
3. Cause
   1. Both variants play as Eli and Fox on the same server.
   2. Run at once, both saves list under "Eli, Fox"; the app was right.
4. Tasks:
   1. Each variant plays under its own names.
5. Done when:
   1. The two variants pass together in `make e2e`.

### [x] 19.81 e2e: a drag across the history selects nothing in `qps`

1. Goal: the copyable-history test passes every run.
2. Report (agent, 2026-10-27)
   1. Found by `make screenshots LOCALE=qps` during 16.7, chromium.
   2. `page.spec.ts:83`: "a drag across the table selects nothing; the history and address do".
   3. After the drag across the open history, the selection was empty.
   4. 16.7 changes e2e text lookups only; this test is unchanged.
   5. Evidence: `.cache/e2e-evidence/16.7-shots-qps/` (not in git).
   6. Failed again on the rerun: `.cache/e2e-evidence/16.7-shots2-qps/`.
   7. Also at 7c47721, before 16.7: 2 of 8 with `--repeat-each 8`, `qps`.
   8. English passed 3 of 3; so a real failure, not load.
3. Cause
   1. Not found yet.
   2. The open history's list box is only 59 px tall at 1024×768.
   3. In `qps` one line wraps over most of it.
   4. A failing drag started in a line's padding; the middle hit the same line.
   5. Guess: the drag stays in one line, or the list scrolls under it.
4. Tasks:
   1. Find the cause in the evidence; fix the test or the app.
5. Done when:
   1. Cause written here; the test passes ten full runs in a row.
6. Cause (2026-10-27)
   1. Chromium starts no selection from the pixel between two text lines.
   2. That is between two rows, or two wrapped lines of one row.
   3. WebKit selects there; so a Chromium-only miss.
   4. The test's drag started 4 px below the list's top edge.
   5. Where that lands depends on the scroll and on line wraps.
   6. `qps` wraps more lines, so it hit a line edge more often.
   7. The drag also started before the loot line came.
   8. A late line scrolls the list under the mouse.
7. Fix
   1. Test only; the app is right.
   2. The drag waits for the loot line, the newest row.
   3. It starts in the middle of that row's first line.
   4. Not app CSS: margins instead of padding moved the dead pixel, not removed it.
   5. A whole-pixel line height did not remove it either.
   6. 20 runs in `qps` and 10 in English, both browsers, pass.
   7. `make e2e` and `make screenshots LOCALE=qps` pass.
   8. Evidence: `.cache/e2e-evidence/19.81-*` (not in git).

### [x] 19.82 e2e: the dice menu test fails in a screenshot run without card images

1. Goal: `make screenshots` passes in a public clone.
2. Report (agent, 2026-10-27)
   1. Found by `make screenshots` in a clean Ubuntu 24.04 container during 16.8.
   2. The clone had no private images submodule; its folder was empty.
   3. `settings.spec.ts:135`: "the dice menu picks the look of your dice".
   4. The dice preview showed 0 faces; the test wants 6.
   5. Failed again alone, with one worker; `make e2e` passes it.
   6. Evidence: `.cache/e2e-evidence/16.8-ubuntu-shots/` (not in git).
3. Cause
   1. With `SCREENSHOTS=1` the test skips its one-pixel dice route.
   2. `realAssets` then routes to `pkg/card_db/images`, which exists but is empty.
   3. Every die texture is missing, so no face image is drawn.
4. Tasks:
   1. Keep the pixel route when the dice images are missing.
5. Done when:
   1. `make screenshots` passes this test without the images submodule.
6. Fix
   1. The test always routes dice to a one-pixel image.
   2. `realAssets` now calls `route.fallback()` for a missing file.
   3. So real dice still win in a screenshot run; missing ones get the pixel.
   4. Checked with `assets/dices` moved away: old code 0 faces, new code passes.

### [x] 19.83 The `cardgen` binary is tracked in git

1. Goal: no build output in the repo.
2. Report (owner, 2026-10-27)
   1. "Why we store cardgen binary in repository?"
   2. `cardgen` in the repo root: a 5 MB macOS arm64 executable.
3. Cause
   1. `go build ./cmd/cardgen` in the root writes `./cardgen`.
   2. `.gitignore` ignored `/four_souls`, not `/cardgen`.
   3. It came in with be58723 (16.4) by a broad `git add`.
4. Tasks:
   1. Untrack `cardgen`; ignore `/cardgen`.
5. Done when:
   1. `git ls-files cardgen` prints nothing.
6. Fix
   1. `git rm --cached cardgen`; `/cardgen` in `.gitignore`.
   2. The 5 MB stays in history; `main` is not rewritten.

### [x] 19.84 The host cannot end a game on a dedicated server

1. Goal: the host ends a running game for everyone, on any server.
2. Report (owner, 2026-10-27)
   1. "When player host game on dedicated server it could not be finished."
   2. "Allow the host end the game."
3. Cause
   1. The game menu offers Save and leave, and Leave game.
   2. Leave game on the app's own server stops that server.
   3. On a dedicated server it only drops the connection.
   4. The seat stays; the game pauses for the others (N-08).
   5. Only a save or an idle limit ends it (N-11, UC-29).
   6. A save keeps the game on the server; it does not end it.
4. Tasks:
   1. Protocol: `end`, client → server; notice kind `game_ended`.
   2. Server: only seat 1 ends; the room stops like an idle end.
   3. Connected players get the notice; the others when they return.
   4. The record ends unfinished.
   5. Menu, host only: **End game**, then a confirm step.
   6. The others see "The host ended the game" and the game list.
   7. The host sees the game list, no notice.
5. Done when:
   1. Go: the host ends it; another seat cannot; an away player hears it.
   2. e2e: the host ends; both are on the game list.
   3. UC-35; wiki `Client`, `Server`; CHANGELOG line.
   4. Screenshots: `.cache/screenshots/19.84/before/` and `after/`.
6. Fix
   1. Protocol: `end`, client → server; notice kind `game_ended`.
   2. Server: `Room.end` in `limits.go`; seat 1 only, else `not_host`.
   3. A game that is over is refused: `no_running_game`.
   4. It stops the room through `endWith`, shared with the idle limits.
   5. Away players are remembered; they hear it on their next hello.
   6. Client: Menu → **End game** → "End the game?" → **End game** or **Back**.
   7. `endsGame` counts `game_ended`; everyone lands on the game list.
   8. The store's `ending` hides the notice from the host who sent it.
   9. Overlay buttons wrap to a second row; four no longer squeeze.
   10. Go: `TestHostEndsGame` in `limits_test.go`.
   11. Vitest: `reconnect.test.ts` notices; `protocol.test.ts` golden lines.
   12. e2e: `the host ends the game; everyone is back on the game list (19.84)`.
   13. It sits in `save.spec.ts`.
   14. Docs: UC-35; ADR 006 items 36, 37; wiki `Client`, `Server`; CHANGELOG.
   15. Screenshots: `.cache/screenshots/19.84/before/` and `after/`.
7. Owner notes (2026-10-27)
   1. Owner marks 19.84 solved.

### [ ] 19.85 e2e: the activity test fails when the attacker has Bob's Brain

1. Goal: `prompt.spec.ts` "the others see what the active player does" passes on every deal.
2. Report (agent, 2026-10-27)
   1. Found by `make e2e` during 18.9, in chromium.
   2. The active player was Eden with Bob's Brain.
   3. It timed out at line 277, waiting for the stack's `.stack-target`.
   4. It passed 12 of 12 on rerun with other deals.
   5. Evidence: `.cache/e2e-evidence/18.9-bobs-brain/` (not in git).
3. Cause
   1. Bob's Brain triggers when an attack is declared.
   2. Its roll waits for a throw before the attack roll does.
   3. `settle` stops at the first Throw: Bob's Brain's roll.
   4. That stack row has no target; the test wants the monster's.
   5. Test logic: the app was right (group I).
4. Tasks:
   1. Wait for the attack roll's row, not any Throw; or ban Bob's Brain, with the reason.
5. Done when:
   1. The test passes with Bob's Brain in the attacker's items.
6. Fix
   1. The evidence: the attack row was "Tarn attacks…", no target yet.
   2. Bob's Brain's roll comes before the target is chosen.
   3. `settle` now waits for a Throw and a `.stack-target` row.
   4. No ban: the test plays any deal.
   5. The test body is `attackActivity`, run twice.
   6. New run: Bob's Brain as the first item, through a saved game.
   7. It is `the attack line waits for a roll the attack triggers: Bob's Brain (19.85)`.
   8. It takes no screenshots; the plain run keeps them.
   9. With the old wait it fails as reported; with the new one, 8 of 8.
   10. No app change; no docs or screenshots.
7. Owner notes (2026-10-27)
   1. Owner marks 19.85 solved.
8. Reopened (owner, 2026-10-27)
   1. Found by a full `make e2e` run during 14.15, webkit.
   2. Test: `the attack line waits for a roll the attack triggers: Bob's Brain (19.85)`.
   3. It failed in `untilAttackRoll`: "the game never settled" (1.5 min).
   4. Turn 2, active player 0: "Your turn", nothing on the stack.
   5. The attack never started; no Throw showed.
   6. History: Bob's Brain triggered, rolled 2; Wyn died.
   7. It passed 5 of 5 on rerun, webkit, alone.
   8. The next full run did not fail it.
   9. Cause not found yet; the 14.15 change does not touch it.
   10. Evidence: `.cache/e2e-evidence/19.85-reopen-webkit/` (not in git).
   11. It holds the error context, the match records and the run log.

### [x] 19.86 e2e: a Throw click in the Meat! odds test did not land, webkit

1. Goal: `stats.spec.ts` "the open attack says what the roll needs" passes under load.
2. Report (agent, 2026-10-27)
   1. Found by `make e2e` during 18.9, in webkit, full run.
   2. Line 152: after the Throw click, the stack still said "Throw the dice".
   3. It passed 8 of 8 alone, one worker at a time.
   4. Evidence: `.cache/e2e-evidence/18.9-stats/` (not in git).
3. Cause
   1. The record has no throw: the server never got one.
   2. A refused intent clears `sent`; the button would come back.
   3. The server closes a client that reads too slowly (`hub.go`).
   4. An intent sent as the socket closes is lost.
   5. The app reconnects; the resync has the same step.
   6. `sent` still equals that step: `canAct` stays false.
   7. So the Throw button never comes back; Pass stays disabled.
   8. The probe runs did not catch the drop itself; the e2e test makes it.
   9. Client: state from the old connection is kept (group O).
   10. 19.78 is the same test and picture; likely the same cause.
4. Tasks:
   1. Read the evidence; find why the throw did not start.
5. Done when:
   1. The cause is known; the test waits for the right thing.
6. Fix
   1. A `welcome` clears `sent` (`store.ts`).
   2. The resync then shows whether the intent was played.
   3. A lost one can be sent again.
   4. e2e: `an intent lost with a dropped connection can be sent again (19.86)`.
   5. It sits in `pause.spec.ts`; it failed in both browsers before the fix.
   6. Vitest: `store.test.ts`, a welcome forgets the sent intent.
   7. Docs: UC-07 item 4a.
   8. No screen change; no screenshots.
7. Owner notes (2026-10-27)
   1. Owner marks 19.86 solved.

### [x] 19.87 e2e: the 19.84 end-game test finds two Games headings, webkit

1. Goal: `save.spec.ts` "the host ends the game; everyone is back on the game list" passes every run.
2. Report (agent, 2026-10-27)
   1. Found by a full e2e run during 19.86, webkit.
   2. Line 76: `getByRole('heading', { name: 'Games' })` matched two headings.
   3. "Games" and "Saved games": a strict mode violation.
   4. Evidence: `.cache/e2e-evidence/19.86-run1-save/` (not in git).
3. Cause
   1. The name matches as a substring without `exact: true`.
   2. "Saved games" shows when another test's save is on the server.
   3. Test logic: the app was right (group I).
4. Tasks:
   1. Match the heading with `exact: true`.
5. Done when:
   1. The test passes with a saved game on the server.
6. Fix
   1. The Games heading matches with `exact: true`.
   2. Same lookup fixed in `lobby.spec.ts` (7 places).
   3. `restart.spec.ts` had it already.
   4. Checked with a save copied into each test server's folder.
   5. Old line: the same strict mode error, both browsers.
   6. Fixed line: passes in both browsers.
   7. No app change; no docs or screenshots.
7. Owner notes (2026-10-27)
   1. Owner marks 19.87 solved.

### [x] 19.88 The join row's password field is too small

1. Goal: a guest can see and type into a locked game's password field.
2. Report (owner, 2026-10-27)
   1. The password field is too small on the client.
   2. Screenshot `screen-join-password`: about 20 px wide, next to Join.
3. Cause
   1. `.game-password` had `flex: 1`, `min-width: 0`.
   2. The row text did not shrink; the field took what was left.
   3. A long set list left it almost nothing.
4. Tasks:
   1. Client: the field keeps a fixed width; the row text wraps instead.
5. Done when:
   1. e2e: the field is at least 120 px wide with a long set list.
   2. Screenshots: `.cache/screenshots/19.88/before/` and `after/`.
6. Fix
   1. Client: `.game-password` is `flex: none`, 140 px wide.
   2. The row text wraps beside it; Join stays at the end.
   3. e2e: `a game password keeps others out (N-14)` checks the width.
   4. It fails on the old style: the field was 24 px.
   5. No use case or wiki change: the behavior is the same.
   6. Screenshots: `.cache/screenshots/19.88/before/` and `after/`.
7. Owner notes (2026-10-27)
   1. Owner marks 19.88 solved.

### [ ] 19.89 e2e: the 15.77 damage fly test timed out once, webkit, full run

1. Goal: `playback.spec.ts` "damage from a card no longer on screen flies from the player who dealt it (15.77)" passes every run.
2. Report (agent, 2026-10-27)
   1. Found by the full e2e gate during 19.88, webkit.
   2. Timeout of 30 s in `centre`, called at line 1224.
   3. Passes alone: 3 of 3 runs, webkit.
   4. Evidence: `.cache/e2e-evidence/19.88-gate-playback-webkit/` (not in git).
3. Cause
   1. Not found yet.
4. Tasks:
   1. Find what `centre` waits for at line 1224 under load.

### [x] 19.90 Others never learn the Cain's Eye choice

1. Goal: every player sees if the top card was kept or put on the bottom.
2. Report (owner, 2026-10-27)
   1. A player with Cain's Eye looks at the top loot card.
   2. They keep it on top or put it on the bottom.
   3. The other players cannot tell which.
3. Cause
   1. `cardkit.LookMayBottom` moves the card with `DeckToBottom`.
   2. `DeckToBottom` emits no event; keeping emits none either.
   3. Others only get the hidden `looked_at` line.
   4. Same effect: Golden Horseshoe (treasure), Purple Heart (monster).
4. Ruling basis
   1. At a real table the move is seen; the card face is not.
   2. So the choice is public; the card stays hidden from others.
5. Tasks:
   1. Engine: one event after the choice, private, with the card.
   2. Its text says the deck and top or bottom.
   3. No barrier: the looked_at before it already is one.
   4. Client: a history line for the chooser and for others.
   5. Others' line names no card.
6. Done when:
   1. Go: both answers emit the event; others get no card.
   2. Vitest: the history lines for the chooser and for others.
   3. e2e: another player's history shows the choice.
7. Fix
   1. Engine: `Game.KeepOrBottom` emits `top_choice` after the answer.
   2. Text: "loot deck top" or "loot deck bottom"; private, with the card.
   3. `cardkit.LookMayBottom` calls it: Cain's Eye, Golden Horseshoe, Purple Heart.
   4. No barrier: in `privateNoBarrier`; ADR 020 item 4.24.
   5. Client: the chooser reads "You put <card> on the bottom".
   6. Others read "Ann puts the top loot card on the bottom".
   7. Go: `TestCainsEyeTellsEveryoneTheChoice`, both answers.
   8. Vitest and e2e: `info.spec.ts` (19.90).
   9. Docs: UC-18 2c, ADR 020, wiki Client and Rules-engine.
   10. Screenshot: `.cache/screenshots/19.90/after/` only.
   11. Before: no line at all; the old client had no case to draw.
8. Owner notes (2026-10-27)
   1. Owner marks 19.90 solved.

### [ ] 19.91 CI: the Go lint job fails in `make i18n`

1. Goal: the Go lint job passes on CI, so the release tag builds.
2. Report (owner, 2026-10-27)
   1. CI run of `make lint TAGS=webkit2_41` failed.
   2. Every golangci-lint run printed "0 issues".
   3. `make i18n` failed: `npx vitest` could not load `frontend/vite.config.ts`.
3. Cause
   1. The Go lint job never runs `npm ci` in `frontend/`.
   2. `npx` then fetched the latest vitest, which cannot find `vite`.
4. Tasks:
   1. Go lint job: set up Node at `NODE_VERSION` and run `npm ci` in `frontend/`.
5. Done when:
   1. CI's Go lint job passes.
   2. A clean checkout runs `make i18n` after `npm ci`.

### [x] 19.92 A two-sided card's back is dealt as a character

1. Goal: every character starts front up, with its starting item.
2. Report (owner, 2026-10-28)
   1. A game started with The Enigma on its back side.
   2. The player had no starting item.
3. Cause
   1. `internal/server` `characters()` lists every `CharacterCard`.
   2. A back face is a `CharacterCard` too: `amgine_eht`, `the_revenant`.
   3. A random or drafted pick can give the back; no item comes with it.
   4. The engine's `NewGame` takes such a pick without a check.
4. Tasks:
   1. Server: the pool skips back faces.
   2. Engine: `NewGame` rejects a back face in `Setup.Characters`.
5. Done when:
   1. Go: the server pool has no back face and keeps the front.
   2. Go: `NewGame` with a back face in `Setup.Characters` fails.
   3. Docs: ADR 023 already says a back is never dealt; no change.
6. Fix
   1. Server: `characters()` skips `Outside` cards, as the engine's `inDecks` does.
   2. Engine: `NewGame` rejects an `Outside` card in `Setup.Characters`.
   3. Go: `TestBackFacesAreNotDealt` (`internal/server`).
   4. Go: `TestBackFaceIsNotDealt` (`pkg/rules_engine/cards`).
   5. Both new tests fail on the old code.
   6. No screen change: no screenshots.
7. Owner notes (2026-10-28)
   1. Owner marks 19.92 solved.

### [x] 19.93 Players cannot change the sound volume

1. Goal: each player sets the sound volume in Settings (ST-09).
2. Report (owner, 2026-10-28)
   1. Sounds play at one fixed level; a player cannot turn them down or off.
3. Cause
   1. `frontend/src/sounds.ts` plays every sound at volume 1.
   2. Settings had no volume field.
4. Tasks:
   1. Go: `Volume` (0–100, default 100) in `internal/settings`.
   2. Frontend: a Sound volume slider on the Settings screen.
   3. Frontend: `setSoundVolume` scales every sound; the store applies it.
5. Done when:
   1. Go: `TestVolume` passes; a file without a volume loads at 100.
   2. Vitest: `sounds.test.ts` plays at the player's volume; 0 is mute.
   3. e2e: `playback.spec.ts`, `pause.spec.ts` and the settings specs pass in webkit.
6. Fix
   1. `Volume` in `settings.go`, checked on load and save.
   2. `setSoundVolume(percent)` in `sounds.ts`, called from `store.ts`.
   3. A slider in `SettingsScreen.tsx`; its readout is a `span`, not `output`.
   4. UC-10 alternative 3i; wiki `Client.md` sound paragraph.
   5. Screenshots: not taken for this change.
7. Owner notes (2026-10-28)
   1. Owner asked for the volume setting and marks 19.93 solved.

### [x] 19.94 Eden's starting choice comes after the first-player roll

1. Goal: setup finishes before the first-player roll (R-SETUP-09, R-SETUP-11).
2. Report (owner, 2026-10-28)
   1. The owner read the order as: setup all characters, Eden's item, then the first player.
   2. The game rolled for the first player first, then opened Eden's choice.
3. Cause
   1. `run()` in `pkg/rules_engine/turn.go` checks `FirstRoll` before the queue.
   2. `NewGame` queues Eden's choice, so it waits for the roll.
   3. Came with 15.43 (Result, item 6); no rule or owner decision asked for it.
4. Tasks:
   1. Engine: the queue runs before the first roll.
   2. Docs: UC-26 alternative 6a and the 15.43 result note.
   3. e2e: the `first-roll.spec.ts` comment on Eden.
5. Done when:
   1. Go: Eden's choice is asked before the first throw.
   2. `make lint test` pass; `make frontend e2e` passes.
6. Fix
   1. `run()` runs the queue (Eden) before the first roll.
   2. `TestEden` expects Eden's choice first, then the throw.
   3. e2e roll helper picks Eden's option before throwing.
   4. UC-26 6a, the Rules-engine wiki and the first-roll spec comment follow.
   5. No screen change: no screenshots.
7. Result (2026-10-28)
   1. `go test ./...` and `make lint` pass.
   2. `make frontend e2e`: 380 passed, exit 0.
8. Owner notes (2026-10-28)
   1. Owner marks 19.94 solved.

### [x] 19.95 Lil Battery recharges a passive or paid item

1. Goal: "Recharge an item" picks only items with a ↷ ability.
2. Report (owner, 2026-10-28)
   1. Lil Battery offered a passive item as its target, and a paid item.
   2. The owner expects only active items, as with 19.64.
3. Tasks:
   1. Find the cause in the code.
   2. Apply the rule to every card that recharges an item.
4. Done when:
   1. Engine test: Lil Battery offers no passive or paid item.
   2. Engine test: Battery Bum and The Battery offer no passive item.
5. Cause
   1. `lil_battery.go`, `battery_bum.go` and `the_battery.go` use `TargetItem` with no filter.
   2. Every item in play is offered, passive and paid ones too.
   3. 19.64 added the `CardDef.Tap` filter for deactivating; recharge has none.
6. Fix
   1. `engine.ActiveItem`: a target filter for items with a ↷ ability.
   2. Lil Battery and Battery Bum use it; The Battery uses it with "another".
   3. Tests `TestLilBatteryOffersOnlyActiveItems`, `TestBatteryBumOffersOnlyActiveItems`.
   4. `TestTheBattery` also checks that a passive item is not offered.
   5. Wiki: `Effect-blocks.md` lists `ActiveItem`.
   6. No screen change: no screenshots.
7. Result (2026-10-28)
   1. `make lint` and `make test` pass (exit 0).
8. Owner notes (2026-10-28)
   1. Owner marks 19.95 solved.

### [x] 19.96 Magic Marker does not work properly

1. Goal: Magic Marker is left out of every game until it works.
2. Report (owner, 2026-10-28)
   1. Magic Marker is not working properly for now.
   2. Ban it from every game, not only from games where the host bans it.
3. Tasks:
   1. Leave `magic_marker` out of every deck, in every game.
   2. Keep the card's code and tests, so it can come back.
4. Done when:
   1. Engine test: no game has a `magic_marker` object (any seed).
5. Fix
   1. `alwaysBanned` in `pkg/rules_engine/ratios.go`: a list of cards no deck gets.
   2. `inDecks` checks it with the host bans (GS-13).
   3. Both deck paths (ratios and fixed decks) use `inDecks`.
   4. Test `TestMagicMarkerIsLeftOut` in `pkg/rules_engine/cards/r/magic_marker_test.go`.
   5. No screen change: no screenshots.
6. Result (2026-10-28)
   1. `make lint` and `make test` pass (exit 0).
7. Owner notes (2026-10-28)
   1. Owner marks 19.96 solved.

### [ ] 19.97 Forever Alone's activation repeats, and the turn never ends

1. Goal: a turn always ends; an activation cannot loop without limit.
2. Report (found by `cmd/server_load`, 2026-10-28)
   1. Random games on a local dedicated server, 2 to 4 players, `b2`.
   2. Two games stopped at turn 26 and turn 52 and went on for 20 minutes or more.
   3. The steps kept growing; the turn number did not.
3. Cause (from the match records, replayed)
   1. The active player activates Forever Alone, then picks an option.
   2. Each pick opens a priority window; the window passes; the card can be activated again.
   3. Only Pass and the activation are offered; End turn is never offered.
   4. The match record of that game (a server run on 2026-10-27, game g11) shows the cycle in its last 40 intents.
4. Tasks:
   1. Read Forever Alone's card code and its activations.
   2. Find why an option reopens the activation window.
   3. Decide the rule with the owner: does one activation end the card's use for the turn?
5. Done when:
   1. A game with Forever Alone out ends its turn.
   2. A test plays the activation loop and checks that End turn is offered.

### [x] 19.100 A played Dagaz waits for a play menu or aim target that does not show (19.18)

1. Goal: `stack.spec.ts` (19.18) passes in every browser, every run.
2. Report (`make e2e` on `step-18-play-test`, 2026-10-28, third run)
   1. `tests/stack.spec.ts:68`, webkit: a played Dagaz goes on the stack.
   2. `answerPlay` (`turn.ts:319`) waited for the play menu or an aim target.
   3. Neither showed within 5 s: "element(s) not found".
   4. Failed once; not seen in the two earlier runs.
   5. The wait was in `playLoot` (`turn.ts`), not `answerPlay`.
3. Cause: a new state closed the card's menu before it was used.
   1. Clicking Dagaz in hand opens its menu (`ActionMenu`).
   2. Any new game step runs `setAiming(null)` (`GameScreen.tsx`).
   3. `setAiming` also clears `selected` (`store.ts`): the menu closes.
   4. The test waited 5 s for a menu that the app had closed.
   5. Nothing reopened it: the test pressed the card once.
   6. The app was right: a new state ends a selection (15.25).
   7. Reproduced: a `setAiming(null)` right after the click fails the test.
   8. Same message as the report: `getByRole('menu').or(…)` not visible.
   9. Evidence: `.cache/e2e-evidence-19.100/`.
4. Tasks:
   1. Find what Dagaz's play shows in webkit, and when.
   2. Fix the wait to the exact state (CLAUDE.md 16.10); no ban.
5. Done when:
   1. The test passes in repeated full runs of `make e2e`.
6. Result (2026-10-28)
   1. `playLoot` presses the card again until its menu or target shows.
   2. No ban; no app change; no screenshots.
   3. A close injected after the click fails before the fix, passes after.
   4. `make lint` passes (exit 0).
   5. `make e2e` passes (380 passed, 6 skipped).
   6. The chance failure did not come back in three WebKit runs.
7. Owner notes (2026-10-28)
   1. Owner marks 19.100 solved.

### [ ] 19.101 A bomb aimed by two clicks shows no "pick a monster" note (9.2)

1. Goal: `aim.spec.ts` (9.2) passes in every browser, every run.
2. Report (`make e2e` on `step-18-play-test`, 2026-10-28, WebKit diagnostic run)
   1. `tests/aim.spec.ts:33`, webkit: a bomb aimed at a monster by two clicks.
   2. The click on a shop item, no target, did not show `.aim-why`.
   3. Expected "Pick a monster or a player."; "element(s) not found" in 5 s.
   4. Failed once in three full WebKit runs.
3. Cause: not found yet.
   1. The test checks the note right after that click.
   2. Possible: a new state ended the aim before the click, as in 19.100.
   3. Possible: the note showed and closed before the check.
4. Tasks:
   1. Find what the click did in the failed run, and when the state changed.
   2. Fix the wait to the exact state (CLAUDE.md 16.10); no ban.
5. Done when:
   1. The test passes in repeated full runs of `make e2e`.
