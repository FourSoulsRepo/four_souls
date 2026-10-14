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

### [ ] 19.7 Buy is lit, then the purchase fails

1. Goal: Buy is lit only when the player can buy.
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
   2. A test that plays the failing deal.
5. Cause
   1. Not found yet.

### [ ] 19.8 Intent waits forever after a pass

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
   1. Not found yet.

### [ ] 19.9 Deck attack: the drawn card hides during the slot choice

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
