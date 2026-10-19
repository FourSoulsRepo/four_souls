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
   11. Requiem (`r`): 247 cards, 269 with copies; 2026-10-19.
   12. **Owner:** Requiem goes through 14.1–14.3 in full; Warp Zone later.
   13. TTS bag "Requiem"; set order 4, released 2022-10 (ADR 010).
   14. 12 reprints: Eden, coins, Bomb!, Butter Bean!, and others.
   15. Tools: two-sided names ("The Enigma/amginE ehT") match TTS.
   16. Tools: one card listed twice by the site is one card (`r-eden`).
   17. Overrides from the printed cards: Lunch, Shell Game, Experimental Treatment.
   18. Overrides: Charmed Dip stats, Charmed Monstro's soul.
   19. Each other card is checked against its image in 14.3.
   20. The Harbingers ("Outside Card") has no engine kind yet; 14.3.
   21. Rooms are printed landscape; the client must show them so.

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
   15. Requiem: rulings from S-OFF, S-YUG, S-DIS, Reddit, X; 2026-10-19.
   16. New rules: R-CARD-36 – 40, R-SETUP-16, R-TURN-17 – 19, R-ATK-20 – 24.
   17. New rules: R-ABIL-37 – 39, R-MECH-53, R-DICE-12.
   18. **Owner:** the room deck is a host option, off by default (R-SETUP-16).
   19. **Owner:** Requiem is ticked by default in the host screen.
   20. **Owner:** Magic Marker changes only common numbers (Q-14, CD-13).
   21. **Owner:** Requiem's client work goes to 14.13 – 14.15.
   22. Open for the **Owner**: Q-15 – Q-18; the engine follows the proposals.
   23. Engine: room deck and slot, the room step, a room's "you".
   24. Engine: attack a player, attackable items and rooms.
   25. Engine: extra attack rolls, flip, voting, outside cards.
   26. **Owner** resolved Q-15 – Q-18 (2026-10-19); Lemegeton as written.
   27. Spindown Dice follows the Russian print: non-eternal, "if you do" (Q-16).

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
   11. Requiem: 19 loot cards tested; Hagalaz, Broken Remote wait for rooms.
   12. Golden Trinket waits for "gains abilities" (Glitch); Magic Marker for Q-14.
   13. Engine: doubled rewards, eternal till end of turn, Rib of Greed's protection.
   14. Engine: destruction replaced (Bag Lunch), combat flag on damage changes.
   15. Engine: a reprint borrows its card from a set not in the game (ADR 010).
   16. Engine: the D10 fix 19.22 copied from step 15; identical lines.
   17. Requiem: 27 active, paid, one-use and soul treasures tested; Red Key waits for rooms.
   18. Engine: copies till your next turn, lost eternal, cards under an item, a gained-item event.
   19. Engine: Car Battery's recharge rule, cheaper $ costs, a 2-or-3-to-1 roll replacement.
   20. Engine: a character replaced from outside the game (The Clicker, Q-18).
   21. Engine: a treasure soul keeps its abilities unless it loses them (R-CARD-38).
   22. Purchases from the treasure deck now emit "purchased" too; shop steals are "stolen".
   23. **Owner:** Berkano and Spindown Dice steal only if the item was destroyed (Q-16, Q-19).
   24. Requiem: 23 passive treasures tested; Backstabber, Punching Bag wait for attacks.
   25. Engine: statics with a computed amount; level counters by name.
   26. Engine: buying the treasure deck's top is a treasure gain; Damocles applies (S-YUG).
   27. Engine: attacks on players and on attackable items and rooms (R-ATK-20, R-MECH-04).
   28. Engine: helpers' attack rolls (R-ATK-21), rewards to the killer, Indomitable (R-ABIL-23).
   29. Requiem: Backstabber, Punching Bag, The Siren, Mega Satan, Ultra Greed, Sisters Vis tested.
   30. Requiem: Sibling Rivalry, Lust for Blood tested; Backstabber's DC fixed by an override.
   31. Requiem: 23 basic, cursed, holy and charmed monsters tested.
   32. Engine: attack rolls counted per attack (Null), two dice for every roll (Cursed Lil Haunt).
   33. Engine: AskWithRoll, a roll made after a question finds its result table.
   34. Every Requiem card has its own test file, as the status report counts them.
   35. Requiem: 11 bosses, 8 events and 4 curses tested; Dogma waits for flip.
   36. Engine: "can't gain souls" checked everywhere a soul is gained; a curse passed on.
   37. Engine: an attack cancelled after combat damage; pending damage names its stack item.
   38. Card tests draw plain cards first (Clotty, Book of Sin): no test depends on the deal.
   39. Engine: two-sided cards flip; the back face is its own definition (R-MECH-33).
   40. Back faces: The Revenant, amginE ehT, The Beast!, checked against foursouls.ru images.
   41. Requiem: Anima Sola, The Enigma, Flip, The Harbingers, Dogma tested.
   42. Requiem: 19 characters and 17 eternal items tested.
   43. Engine: gained abilities (Glitch); rerolls with 2 dice (Classic Roller); ATK set (Hemoptysis).
   44. Engine: triggered abilities as targets (Dead Weight); the penalty's ¢ loss marked (Soulbond).
   45. Requiem: Golden Trinket and the 5 bonus souls tested.
   46. Engine: active bonus souls trigger and take counters outside the game.

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

### [x] 14.9 Host bans cards

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
4. Split (owner, 2026-10-15):
   1. This step: engine and server only; no host screen yet.
   2. Host screen, game list, table, their e2e test, use case: 14.12.
5. Result (2026-10-15):
   1. `engine.Setup.Banned`: no copy goes into a deck, none is dealt.
   2. The list stays in `Game.Banned`, so saved games keep it.
   3. `host_bans` takes any card of the sets, not only characters.
   4. A starting item cannot be banned; ban its character.
   5. A ban that leaves a deck too small is refused (`bad_setup`).
   6. The record header has `setup.banned`.
   7. Tests: `TestBannedCardsAreLeftOut`, `TestHostBansAnyCard`.
   8. UC-05 updated: host bans of any card.
6. Owner notes (2026-10-15)
   1. Owner marks 14.9 solved; the host screen goes on in 14.12.

### [x] 14.10 Tests deal without stat-changing monsters

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
5. Result (2026-10-15):
   1. `startGame` takes `bans`; Playwright adds them to `create`.
   2. `statMonsters` in `e2e/tests/game.ts`: Delirium, Stoney.
   3. The 15.23 test bans them.
   4. `TestStatMonsterBans`: 400 deals, no fresh monster changed.
   5. It also fails if the Base Game gets a new such monster.
   6. Gold Box and Four Souls+: e2e tests play the Base Game only.
6. Owner notes (2026-10-15)
   1. Owner marks 14.10 solved.
   2. Rule: a test bans only the cards it needs, each with a reason.

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

### [ ] 14.12 Ban list in the host screen

1. Goal: the host picks banned cards in the app (GS-13).
2. Tasks:
   1. Engine and server are done in 14.9; this is the client.
   2. Detailed setup: pick cards to ban, from the chosen sets.
   3. Unchecking a set drops its cards from the list.
   4. Starting items are not offered; their character is.
   5. Without card images, the picker shows text cards.
   6. The game list and the table show the ban list.
3. Done when:
   1. e2e test: ban a card in the host screen; it is in no deck.
   2. e2e test: the table shows the ban list to every player.
   3. Use case for the option.

### [ ] 14.13 Rooms on the table

1. Goal: the room deck plays in the app (Requiem).
2. Tasks:
   1. Detailed setup: a "Room deck" checkbox, off by default (R-SETUP-16).
   2. The checkbox shows only when Requiem is chosen.
   3. A room row on the table; rooms are landscape cards.
   4. The active player uses a room's activated abilities.
   5. The end-phase room step: put a room into discard, or keep it.
   6. Without card images, rooms show as landscape text cards.
3. Done when:
   1. e2e test: host with rooms, use a room, change it in the end phase.
   2. Use case for rooms; UC-05 updated.

### [ ] 14.14 Attack a player or an item

1. Goal: Requiem's new attack targets work in the app.
2. Tasks:
   1. Players are offered as targets when an ability allows it (R-ATK-20).
   2. Attackable items and rooms are offered too (R-MECH-04).
   3. Extra attack rolls show whose roll it is (R-ATK-21).
   4. Aiming works as it does for monsters (15.32).
3. Done when:
   1. e2e test: The Siren lets the active player attack a player.
   2. e2e test: Punching Bag can be attacked.
   3. UC-13 updated.

### [ ] 14.15 Two-sided cards

1. Goal: a flipped card shows its other face (R-CARD-39).
2. Tasks:
   1. card_db: a back face with name, text, stats and image (ADR).
   2. Tools: the back image from the card page (`r-the_enigma_back`).
   3. The table shows the face-up side; the zoom shows both.
   4. The history names the flip.
3. Done when:
   1. e2e test: Anima Sola flips to The Revenant and back.
   2. Use case.
