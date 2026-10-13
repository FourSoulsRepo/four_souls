# Step 14. More sets

Goal: official sets beyond the Base Game.
Order: Gold Box, Four Souls+, Requiem, Warp Zone, promos.

Ideas: CD-02, GS-07, GS-08, CD-11, A-16, CD-12, V-32, GS-12, GS-13, GS-01.

---

### [ ] 14.1 Data and images per set

1. Goal: card data for the next set.
2. Tasks:
   1. Run the scrapers from step 2 (`FourSoulsRepo/tools`).
   2. Images into the image repo, own folder.
3. Done when:
   1. The set validates in `card_db`.
4. Progress (owner's one-time jump ahead of step 15):
   1. Gold Box V2 (`g2`): 64 cards, 68 with copies; done.
   2. Four Souls+ V2 (`fsp2`): 90 cards, 98 with copies; done.
   3. TTS bags: "Kickstarter" is Gold Box; "Four Souls+".
   4. Two TTS name typos fixed in the tools' name matching.
   5. Reprints and taken IDs: ADR 010.
   6. g2: 6 reprints; `pills_g2` renamed.
   7. fsp2: 10 reprints; 6 cards renamed, e.g. `isaac_fsp2`.
   8. Every copy now has its own image, Base Game too.
   9. Base Game IDs did not change.
   10. Steps 14.2 and 14.3 for these sets are still open.

### [ ] 14.2 New mechanics

1. Goal: engine support for the set's new rules.
2. Tasks:
   1. Research new rules (step 3 process).
   2. Rules text with new IDs.
   3. New effect blocks where needed.
3. Done when:
   1. **Owner** resolved the set's open cases.
4. Progress:
   1. Gold Box: rulings from S-OFF, S-YUG, S-DIS, Reddit, X; 2026-10-13.
   2. New rules: R-DICE-08, R-MECH-52, R-ATK-19.
   3. Open case Q-10 (Fruitcake); **Owner**: keep the proposal for now.
   4. New blocks are built with their first card in 14.3, tested there.
   5. Blocks: max hand size, recharge limit, roll history.
   6. Blocks: next attack roll, till-end-of-turn abilities, free purchase.
   7. Blocks: reprint copies in the engine's deck setup (ADR 010).
   8. Four Souls+: rulings from S-OFF, S-YUG, S-DIS, Reddit, X; 2026-10-13.
   9. New rules: R-DICE-09 (attack rolls in a row), R-DICE-10, R-DICE-11.
   10. Q-11 (PHD) resolved by the **Owner**: once per roll.
   11. Q-12 (Euthanasia) open; **Owner**: keep the proposal for now.
   12. Blocks: attack-roll repeats, roll change once, several dice, stored roll.
   13. Blocks: revealed hands, attack redirect, win by a card, play lock.
   14. Blocks: forced attacks, cents-gain hooks, free loot play, curse giving.

### [ ] 14.3 Cards per set

1. Goal: every card of the set tested.
2. Tasks:
   1. Stubs, batches, interaction tests (step 5 process).
   2. Focused fuzz on the new cards.
3. Done when:
   1. Status report shows the set at 100%.
4. Progress:
   1. Gold Box: 58 of 58 cards tested; registered in `cards.Sets()`.
   2. Shared helpers moved to `cards/cardkit`; b2 keeps its names.
   3. Engine: reprints, death ward, CancelAttack, Linger, NextAttackRoll.
   4. Engine: roll history, free shop item, recharge limit, hand cap.
   5. Engine: a character's soul value counts (R-CARD-35).
   6. Tests that mean "the Base Game" now use `b2.Set` only.
   7. Four Souls+: 80 of 80 cards tested; registered in `cards.Sets()`.
   8. Engine: copy a ↷ ability, several dice, revealed hands, attack redirect.
   9. Engine: win by a card, free loot play, lingering push, damage dealer.
   10. Simulations found two bugs with copies (Placebo, Diplopia); fixed, tested.

### [ ] 14.4 Setup options

1. Goal: options that come with more sets.
2. Tasks:
   1. Bonus Souls and Rooms in detailed setup.
   2. Deck ratios from Ed's recommendations (GS-07).
   3. Remove multiplayer-only cards in 2-player games (GS-08).
   4. Deck ratios: the collection check still compares whole sets (A-16).
   5. The deck takes every chosen set's copies of a card (ADR 010).
3. Done when:
   1. Each option has a server test and a use case.

### [ ] 14.5 Copy art on the table

1. Goal: each copy shows its own image (CD-12).
2. Tasks:
   1. The engine gives each card instance its print and copy index.
   2. Views and the protocol carry it; records too.
   3. The client picks the image from `Card.Prints`.
   4. A missing copy image falls back to the first one.
3. Done when:
   1. Two coins of one card show different art in a game.
   2. Use case and e2e test.

Repeat 14.1–14.3 for each set.
Alt art (CD-11) comes after all sets.

### [x] 14.6 Host picks the card sets

1. Goal: games can use Gold Box and Four Souls+ (CD-02).
2. Tasks:
   1. Host screen: a box per card set; the Base Game is always in.
   2. The server adds the Base Game if a client leaves it out.
   3. The game list shows each game's sets.
3. Done when:
   1. Server test, e2e test, UC-03 updated.
4. Result:
   1. Before 14.6 the app sent every set it had data for.
   2. So hosting failed between 14.1 and the registration of both sets.
   3. **Owner:** Gold Box and Four Souls+ are core sets, ticked by default.
   4. The host screen remembers the last setup in the settings.
   5. IDs `a_penny_6`, `dice_shard_3`, `soul_heart_2` became `a_penny`, `dice_shard`, `soul_heart`.
   6. Old saved games with those IDs no longer load.

### [ ] 14.7 Switch reprint art

1. Goal: a player picks which print of a card they see (V-32).
2. Tasks:
   1. Only characters and eternal items with reprints.
   2. Right click, or long tap on touch, cycles the art.
   3. Needs the per-copy art from 14.5.
   4. **Owner** decides: only prints from enabled sets, or every print.
   5. **Owner** decides: only the viewer sees it, or everyone.
   6. Card text and stats never change; it is only the image.
3. Done when:
   1. e2e test: a right click on a reprinted character changes its image.
   2. A card without reprints ignores the click.
   3. Use case.

### [ ] 14.8 House rule: no damage loops at 0 HP

1. Goal: the host may stop damage loops at 0 HP for a match (GS-12).
2. Tasks:
   1. Detailed setup: a checkbox, off by default.
   2. A hint next to it explains the rule in plain words.
   3. Hint text: damage to a player at 0 HP no longer fires "take damage" effects.
   4. Hint example: Forever Alone with Gurdy Jr. stops looping.
   5. Engine: a setup flag; R-MECH-51 skips the triggers when it is on.
   6. The flag goes into the record header and saved games.
   7. The game list and the table show that the house rule is on.
3. Done when:
   1. Engine test: Forever Alone with Gurdy Jr. loops when off, not when on.
   2. Server test: the flag reaches the game.
   3. e2e test: the checkbox and its hint.
   4. Use case for the option; Q-13 updated.

### [ ] 14.9 Host bans cards

1. Goal: the host leaves chosen cards out of the game (GS-13).
2. Tasks:
   1. Detailed setup: pick cards to ban, from the chosen sets.
   2. Banned cards are never in a deck and never dealt.
   3. The ban list goes into the record header and saved games.
   4. The game list and the table show the ban list.
3. Done when:
   1. Engine test: a banned card is in no deck.
   2. Server test: the ban list reaches the game.
   3. e2e test: ban a card in the host screen.
   4. Use case for the option.

### [ ] 14.10 Tests deal without stat-changing monsters

1. Goal: e2e tests stop failing on unlucky deals.
2. Found in 15.3 (2026-10-13):
   1. The 15.23 badge test expects no badge on a fresh monster.
   2. Delirium and Stoney give other monsters +1 DC.
   3. Dealt into a slot, they put a dice badge on the first monster.
   4. A probe: 3 of 120 deals failed the check (2.5%).
3. Tasks:
   1. Depends on 14.9; use its bans, no test-only deal code.
   2. Tests that check monster stats ban such monsters.
   3. Find the list in card statics: each monster that changes others.
   4. Base Game: Delirium, Stoney; then each new set.
4. Done when:
   1. The 15.23 test bans them and checks the untouched monster.
   2. Many deals in a row show no stat badge on a fresh monster.

### [ ] 14.11 Draft pick screen

1. Goal: normal games can use draft picking (GS-01).
2. Tasks:
   1. The server part is done in 6.6 (UC-05); this is the client.
   2. Detailed setup: picking (random or draft) and draft size.
   3. Pick screen: each player sees only their own offers.
   4. A click on a character picks it.
   5. Others see who has picked, not what.
   6. Without card images, the offers show as text cards.
   7. No bans here; ban rounds stay with the tournament tools.
3. Done when:
   1. e2e test: host a draft game, each player picks, the game starts.
   2. UC-05 updated; use case for the pick screen.
