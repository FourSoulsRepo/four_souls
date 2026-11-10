# Step 14. More sets

Goal: official sets beyond the Base Game.
Order: Gold Box, Four Souls+, Requiem, Warp Zone, promos.

Ideas: CD-02, GS-07, GS-08, CD-11, A-16, CD-12, V-32, GS-12, GS-13, GS-01, GS-11.

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
   47. Engine: the room deck as a setup option, one room slot, the room step after hand size.
   48. Engine: the active player controls rooms and uses their $ abilities (R-CARD-36).
   49. Requiem: 37 rooms tested (statics, turn triggers, $ rooms, Dice Rooms).
   50. **Owner:** Social Goals counts every player's progress (Q-20); Pity for the Poor makes every purchase free (Q-21).
   51. Engine: giving ¢, killing a player, attack tax, one ability and one loot a turn, rooms kept in play.
   52. Requiem: the other 20 rooms tested, Gus included.
   53. Engine: Magic Marker numbers (R-ABIL-39); a room from the deck onto a room slot.
   54. Requiem: Hagalaz, Broken Remote, Red Key and Magic Marker tested; the report is 235/235.
   55. Requiem registered in `cards.Sets()` and named in the host screen; new hosts get it ticked.
   56. FuzzSets and the simulations turn the room deck on in some games.
   57. Fix from the simulations: a player who can't pay Tax for the Mighty may end the turn.
   58. card_db: Betrayal! and The Harbingers have no reward box (site icons, back face's souls).
   59. New `make fuzz-new SETS=b2,r` (FuzzFocusedSets): fuzz a new set with the Base Game.
   60. Checks: 3000 simulations, 15 min of `fuzz-new SETS=b2,r` (9590 games), `make e2e`.

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
4. Progress (owner's jump ahead, 2026-10-19):
   1. Deck ratios (GS-07) done; Bonus Souls, Rooms, GS-08 still open.
   2. Ratios from S-OFF Deck Ratios: 100 cards a deck (R-SETUP-17).
   3. Groups are the site's card types; cardgen writes `CardSet.Groups`.
   4. `engine.Setup.DeckRatios` picks each group's count at random.
   5. `engine.RatioShortfall` names the groups short of cards.
   6. 2 players drop 3+ cards, so a few set mixes fit 3 players only.
   7. **Owner:** on by default; the host screen locks the box when short.
   8. The tooltip names the short groups; binding `App.DeckRatios`.
   9. Protocol `no_deck_ratios`; the server sets it when the sets are short.
   10. Only all four sets fit 2 players today; Base Game alone does not.
   11. The record header has `deck_ratios`; the table shows "deck ratios".
   12. A-16 holds: the client still gets whole sets, never the dealt decks.
   13. FuzzSets and the simulations turn the ratios on in some games.
   14. Use case UC-20; tests `TestDeckRatios`, `TestDeckRatiosOption`, e2e.
   15. Rooms done in 14.13 (2026-11-10); Bonus Souls was already in.

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

### [x] 14.11 Draft pick screen

1. Goal: normal games can use draft picking (GS-01).
2. Tasks:
   1. The server part is done in 6.6 (UC-05); this is the client.
   2. Detailed setup: picking (random or draft) and draft size.
   3. Pick screen: each player sees only their own offers.
   4. A click on a character picks it.
   5. Others see who has picked, not what.
   6. Without card images, the offers show as text cards.
   7. No bans here; ban rounds stay with the tournament tools.
   8. **Owner:** random (one character each) stays the default (2026-11-09).
   9. Draft size: the host sets it in detailed setup.
       1. Range 2–5, default 3, as the server allows (6.6).
   10. Each offer shows the character and its starting item.
       1. The server sends the item with each offer.
       2. It comes from the engine's card definition (`StartingItem`).
       3. Without card images, the item shows as a text card too.
   11. A character without a fixed starting item shows none (Eden).
       1. Its offer says the item is chosen after the pick.
       2. The player picks it at game start, as today (R-SETUP-09).
       3. Check every set for such characters, reprints included.
3. Done when:
   1. e2e test: host a draft game, each player picks, the game starts.
   2. UC-05 updated; use case for the pick screen.
   3. e2e test: each offer shows its starting item; Eden shows none.
   4. Server test: offers carry the starting items.
4. Result (owner's jump ahead, 2026-11-09):
   1. Protocol: `setup.items`, the offers' starting items, by index.
   2. Eden is the only character with a start choice, in every set.
   3. `TestEveryCharacterHasAStartingItem` fails on a new one.
   4. Detailed setup: "Characters" (random or draft), "Characters offered".
   5. The host screen remembers both (`HostSetup.picking`, `draft_size`).
   6. `PickScreen`: click to pick, Confirm to send, who has picked.
   7. Tests: `TestDraft`, `TestStartingItemsEden`, `e2e/tests/draft.spec.ts`.
   8. UC-39 added; UC-05 and UC-12 updated.
5. Owner notes (2026-11-09)
   1. Owner marks 14.11 solved.

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

### [x] 14.13 Rooms on the table

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
4. Result (owner's jump ahead, 2026-11-10):
   1. Protocol option `room_deck`; the server turns it off without rooms.
   2. Record header `room_deck`; ADR 007 updated.
   3. View: `rooms`, the room slots; presence covers them.
   4. Host screen: "Room deck" in detailed setup, only with Requiem.
   5. Board: a rooms column after the slots; room, room deck, discard.
   6. Client kind `room`: landscape in every size, text cards too.
   7. Room use and the room step need no new client code.
   8. No room back image at first: a face-down room was the plain back.
   9. Tests: `TestRoomDeckOption`, `TestRoomsOption`, `e2e/tests/rooms.spec.ts`.
   10. UC-40 added; UC-05 and the Client wiki updated.
   11. Owner: the room deck showed no back image (2026-11-10).
       1. The tools fetched only loot, treasure and monster backs.
       2. The app served only those three kinds (`backKinds`).
       3. Tools: `backs` also fetches `RoomCardBack.png` (landscape).
       4. Images: `assets/backs/room.webp`, landscape like the rooms (818×600 since 14.16).
       5. App: `backKinds` has `room`; without the images it stays plain.

### [x] 14.14 Attack a player or an item

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
4. Result (owner's jump ahead, 2026-11-10):
   1. Engine, server and view were done before: targets, odds, helper rolls.
   2. The client's pick, glow and arrow already took players and items.
   3. One bug: an attack roll on a player named the character card.
   4. "Ada rolls against Judas"; now "Ada rolls against Bo".
   5. `attackedPlayer` in `stack.ts`; the stack and the history use it.
   6. Helpers' rolls already named the helper: Throw, banner, stack line.
   7. The odds row stays the active player's during a helper's roll.
   8. Tests: `stack.test.ts`, `e2e/tests/attack-targets.spec.ts` (Siren, Punching Bag, Sibling Rivalry).
   9. UC-13: 5o, 5p, 5q; Client wiki: `attackedPlayer`.

### [x] 14.15 Two-sided cards

1. Goal: a flipped card shows its other face (R-CARD-39).
2. Tasks:
   1. card_db: a back face with name, text, stats and image (ADR).
   2. Tools: the back image from the card page (`r-the_enigma_back`).
   3. The table shows the face-up side; the zoom shows both.
   4. The history names the flip.
   5. Each side has its own name; never "Front/Back" (owner, 2026-10-27).
   6. A card that changes side plays a flip animation (owner, 2026-10-27).
3. Done when:
   1. e2e test: Anima Sola flips to The Revenant and back.
   2. Use case.
   3. ADR 023 accepted.
4. Progress:
   1. **Owner:** done on `step-18-play-test`, ahead of step 14 (2026-10-27).
   2. **Owner:** the tools fetch the back images now (2026-10-27).
5. Owner notes (2026-10-27)
   1. Owner marks 14.15 solved.

### [x] 14.16 Landscape cards at full quality

1. Goal: a landscape card is as sharp as a portrait one (owner, 2026-11-10).
2. Tasks:
   1. Tools: `resize` limits the width to 600 px, whatever the shape.
   2. A portrait card is 600×818; a room is only 600×440.
   3. The board and zoom draw a room 88/63 wider than a portrait card.
   4. So a room has fewer pixels per screen pixel; its zoom is blurry.
   5. Card images and backs: limit the short side to 600 px instead.
   6. Portrait cards stay 600×818; landscape ones become 818×600.
   7. `assets` (mats, dice) keeps its width limit.
   8. Regenerate every landscape image: the rooms and the room back.
   9. Translated room images too, if any.
3. Done when:
   1. Tools test: a landscape card resizes to 818×600.
   2. Every room image and the room back are 818×600.
   3. No portrait image changes.
4. Result (owner's jump ahead, 2026-11-10):
   1. Tools: `resizeCard` limits the short side; `images` and `backs` use it.
   2. `assets` keeps `resize`, the width limit, for mats and dice.
   3. Test `TestResizeCard`: 962×1312 to 600×818, 1312×962 to 818×600.
   4. Only Requiem has landscape cards: 61 rooms; no translated images.
   5. Reran `images -origin r` and `backs`: 62 files changed, all landscape.
   6. Every portrait image came out byte-identical.
   7. gosec flagged two temp-file writes after the change (G703); marked, with the reason.
   8. The images wait in the private repo for the owner's amend and push.
5. Owner notes (2026-11-10)
   1. Owner amended and pushed the images (`4ad7123`).
   2. Owner marks 14.16 solved.

### [x] 14.17 Challenge rules research

1. Goal: our rules text for co-op challenges and Krampus (GS-11).
2. Owner decisions (2026-11-10):
   1. A "Challenge" game mode beside the regular game.
   2. First challenge: How the Krampus Stole Christmas.
   3. Difficulties now: Normal, Hard, Ultra Hard.
   4. Co-op only: exactly 2 people, one character each.
   5. Solo (one person, two characters): later, 14.27.
   6. Competitive variant: later, 14.28.
   7. No 3–4 player challenges; official rules use 2 characters.
3. Sources (copyrighted; only in the tools repo `.cache/research/maestro/`):
   1. `solo-coop-and-challenge-rules.pdf` and `.txt`: Maestro Media.
   2. `how-the-krampus-stole-christmas.pdf` and `.txt`: Maestro Media.
   3. `pilfering-krampus-cards.png`: the 4 boss cards (page 4 of the PDF).
   4. Rulebook sections: `.cache/research/official/extended-rulebook.md`.
   5. "Solitaire/Co-op Mode" and "Four Souls Challenges" (lines 1010–1105).
   6. Fan copy: foursoulsspoiler.net/challenges (a transcription; low trust).
   7. Links: maestromedia.com/collections/boi-challenges.
4. Facts found (2026-11-10):
   1. Pilfering Krampus: HP 5 Normal, 6 Hard, 6 Ultra Hard; 4+; ATK 2.
   2. Reward +1 treasure; soul value 1, every difficulty.
   3. The competitive boss: HP 5; its own text.
   4. Presents at start (¢, loot, treasure, soul): 8/2/3/3, 12/4/3/3, 16/4/3/3.
   5. Order of the rows: Normal, Hard, Ultra Hard.
   6. "What's in the Box?": 3 minions on Normal, 4 on Hard and Ultra Hard.
   7. "Happy Holidays!": at 5, 6 or 7 souls by difficulty.
   8. The PDF: a final boss leaves its slot only if "defeated in an attack".
   9. The rulebook says "killed"; the research picks one and says why.
   10. The PDF room rule says "after step 4"; the rulebook says "step 3".
   11. foursouls.com lists the challenge cards as origin `challenge` (83 cards).
   12. Krampus: `challenge-pilfering_krampus-{norm,hard,ultra,comp}`.
   13. Its rule cards: `challenge-how_the_krampus_stole_christmas-{norm,hard,ultra,comp}`.
5. Tasks:
   1. New rules file `challenges.md`; topics COOP, CHAL, KRAMPUS.
   2. COOP: co-op setup and play; no solitaire rules yet.
   3. CHAL: final boss slot, boss armor, minion slots, timer, beating a challenge.
   4. KRAMPUS: goal, setup, presents, three global modifiers, per difficulty.
   5. Search S-YUG, S-RED, S-DIS for Krampus challenge rulings.
   6. Open cases go to `open-questions.md` with a proposal each.
   7. Known cases to write down:
      1. "Souls controlled by all players": soul value, or soul cards?
      2. Normal's forced attack: who picks the random minion, and when?
      3. "All I Want for Christmas" with a mixed reward (¢ and loot).
      4. A forgone reward for a destroyed present.
      5. "Round": when does it end with an extra or a skipped turn?
      6. Q-04 (timer on death) does not apply: challenges never tick on death.
      7. Does the regular Krampus (`fsp2`) stay in the monster deck?
      8. The PDF lets it stand in for the boss, which takes it out.
   8. Mark each basis Official, Community or House (rulings-basis rule).
6. Done when:
   1. Rules test passes: every new ID unique and well formed.
   2. The rules website shows the new file, if it lists every topic file.
7. Result (2026-11-10):
   1. New `pkg/rules_engine/docs/rules/challenges.md`.
   2. Rules R-COOP-01–12, R-CHAL-01–15, R-KRAMPUS-01–16.
   3. Topics added to the rules README; the website lists the file last.
   4. Open cases Q-30 to Q-42 in `open-questions.md`; Q-04 notes Q-35.
   5. House answers Q-31–Q-33, Q-39–Q-42 are in "Rulings to revisit".
   6. Final boss leaves when "killed" (S-OFF over the PDF): Q-37.
   7. Room rule after the room step: Q-38.
   8. Only Normal has the forced minion attack; checked on the boss cards.
   9. `sources.md`: new source S-MAE, the cached Maestro Media PDFs.
   10. Community: Yuggy on events in minion slots, timer, regular boss cards.
   11. Nothing on Krampus in S-RED, S-X or S-DIS.

### [x] 14.18 ADR: challenges

1. Goal: one decision record for the challenge mode (ADR 028).
2. Written in the main conversation (CLAUDE.md 5.21.2).
3. Tasks:
   1. Engine: a challenge is a setup option, with its ID and difficulty.
   2. Engine: challenge code lives apart from card sets.
   3. Final boss slots and minion slots as new slot rows.
   4. Presents: where they live, how counters and the flip work.
   5. Timer: game state, shown to every player.
   6. Co-op end: both players win, or both lose; game over says why.
   7. Card data: `challenge` card_db set, never offered as a playable set.
   8. Protocol, record header (ADR 007), saves and restart carry the mode.
   9. Client: where the boss, minions, presents and timer go.
   10. Without card images: text cards, like every other card.
4. Done when:
   1. ADR 028 accepted by the owner; INDEX row added.
5. Result (2026-11-10):
   1. ADR 028 written in the main conversation; INDEX row added.
   2. Accepted on the owner's go-ahead to implement 14.17 – 14.26.
   3. The owner reviews it with the step; changes go back into it.
   4. A challenge is a card set plus a `ChallengeDef`.
   5. Own slot rows: `Bosses`, `Minions`; presents outside the game.

### [x] 14.19 Challenge card data

1. Goal: the Krampus challenge's 8 cards in `card_db`, with images.
2. Tasks:
   1. Tools: `cards -origin challenge`; card data from foursouls.com.
   2. Tools: `build` and `images` for the `challenge` set.
   3. Only the cards that a supported challenge uses go into the game data.
   4. Now: the 4 Pilfering Krampus and the 4 Krampus rule cards.
   5. Kinds for the boss and the rule cards (ADR 028).
   6. Check every stat against the PDF (14.17, item 4).
   7. Images into the private image repo; the owner pushes them.
   8. The game builds and runs without them (text cards).
3. Done when:
   1. `card_db` validates the set; `cardgen -report` lists it.
   2. The host screen does not offer `challenge` as a card set.
4. Result (2026-11-10):
   1. Scraped `cards -origin challenge`: 86 cards; Chrome passed on its own.
   2. New `data/challenge.json`: 8 cards, one copy each; set order 5.
   3. Bosses: `pilfering_krampus_{norm,hard,ultra,comp}`, `Challenge Boss Card`.
   4. Rule cards: `how_the_krampus_stole_christmas_{norm,hard,ultra,comp}`, `Challenge Card`.
   5. IDs: the slugs with `-` as `_`; ADR 028 Decision 1 lists them.
   6. `comp` (14.28) is in: same cards, cheap data, never dealt.
   7. The site prints no type; the tools name them by stats.
   8. Tools: `gameCards` keeps only these 8; no TTS or copies step.
   9. Stats match R-KRAMPUS-03: HP 5/6/6 (comp 5), 4+, ATK 2, +1 treasure, soul 1.
   10. Overrides from the printed cards (`overrides/challenge.json`):
       1. Normal: "minon" twice on the site; printed "minion".
       2. Hard: "5: remove" in lower case.
       3. Comp: roll 5 is the treasure present; the last line was missing.
   11. Engine: new kind `ChallengeCard`, no deck.
   12. cardgen: a `Challenge Boss Card` is a `MonsterCard` with `Outside`.
   13. Stubs in `cards/challenge/`; not in `cards.Sets()`; abilities in 14.22.
   14. Tests: never dealt, Krampus stats, no `challenge` in `CardSets` or create.
   15. Images: 8 landscape WebP files, 839×600, in `images/en/challenge/`.
   16. They wait in the private repo for the owner's amend and push.
   17. The client draws only rooms landscape; challenge cards need it (14.25).

### [x] 14.20 Engine: co-op challenge core

1. Goal: a 2-player co-op game against the game itself.
2. Tasks:
   1. `Setup` takes a challenge and a difficulty (ADR 028).
   2. A challenge needs exactly 2 players.
   3. No bonus souls; 3+ player cards are out (R-CARD-32).
   4. Timer: starts at 8; ticks down at the end of each round.
   5. The timer reaches 0: both players lose.
   6. The challenge's goal met: both players win.
   7. Soul target and "that player wins" effects do not end the game.
   8. With rooms: a room unchanged for a whole round goes to discard.
   9. View: timer, challenge, difficulty; win or loss reason.
   10. Events and history messages for the timer, win and loss.
3. Done when:
   1. Engine tests for each rule ID of COOP and CHAL timer.
   2. A test challenge with a trivial goal plays to a win and to a loss.
4. Result (2026-11-10):
   1. `challenge.go`: `ChallengeDef`, `ChallengeLevel`, `ChallengeState`, `Present`.
   2. `Setup.Challenge`, `Setup.Difficulty`; `Game.Challenge`, nil in a regular game.
   3. Setup errors on 3+ players, bonus souls, a missing set, an unknown difficulty.
   4. Bonus souls are an error, not forced off; the server turns them off (14.23).
   5. The rule card is an object outside the game; `inPlay` lists it last.
   6. Timer: ticks before the starter's next regular turn (Q-34); 0 loses at once.
   7. `WinChallenge()`, `LoseChallenge(reason)`; reason `event.challenge_lost.timer`.
   8. Events `challenge_timer`, `challenge_won`, `challenge_lost`; history lines.
   9. `checkWin` and `WinGame` do nothing in a challenge (R-CHAL-14).
   10. Room rule by object IDs: `TurnFrom`, `RoundFrom` per player (Q-39).
   11. An extra turn checks rooms against the player's previous regular turn.
   12. View: `GameView.Challenge`; players see the partner's hand.
   13. Hand events stay private; only the view shows the partner's hand.
   14. `cards.Challenges()`, `cards.FindChallenge`; Krampus is listed as `krampus`.
   15. Krampus's own setup returns an error until 14.22; its data is tested.
   16. No Russian line: the repo has no `ru.json` yet; English only.
   17. Tests: `challenge_test.go` with a test challenge; win, timer loss, rooms.
   18. `testdata/regular_save.json`: a save from before 14.20 still resumes.
   19. ADR 028 got small drift items: Decision 3.6–3.8, 5.9–5.11.

### [x] 14.21 Engine: final boss and minion slots

1. Goal: the challenge slot rows work with every card.
2. Tasks:
   1. Final boss slot: attackable; can't be covered or leave unless killed.
   2. Boss armor: killed outright or over 3 non-combat damage: takes 3.
   3. Minion slots: a new slot per minion; never refilled.
   4. An event put in a minion slot resolves, then goes to discard.
   5. Cards that target, move or cover monsters respect both rows.
   6. Search every set for such cards; one shared rule (CLAUDE.md 5.24).
   7. Attack targets and odds include both rows.
3. Done when:
   1. Engine tests for each CHAL slot rule.
   2. Fuzz runs with a challenge setup stay green.
4. Result (2026-11-10):
   1. `boss_slots.go`: `FinalBossSlot`, `MinionSlot`; `Game.Bosses`, `Game.Minions`.
   2. Both rows `omitempty`; `TestRegularSaveUnchanged` stays green.
   3. Helpers: `MonstersInPlay()`, `MonsterRows()`; `IsFinalBoss`, `IsMinion`, `Held`.
   4. Engine calls: `PutInFinalBossSlot(card)`, `PutInMinionSlot(id)`, `MinionFromDeck()`.
   5. `Held`: `DiscardMonster`, `PutIntoDeck`, `SlotToDeckBottom` do nothing to it.
   6. Boss armor: one `killObject` for `Kill` and `KillObject`; `armored` in `resolveDamage`.
   7. Before, `Kill` and `KillObject` were two kill paths; now one.
   8. A minion slot goes when its card leaves (`takeFromSlot`); later indices shift.
   9. Invariants: every slot card's zone names its row and slot number.
   10. Events in a minion slot: `finishEvent`, `afterEventAbility`, `BecomeSoul`.
   11. A minion dies as a normal monster; Krampus's rule is 14.22's.
   12. View: `GameView.Bosses`, `Minions`; protocol encoder and `protocol.ts`.
   13. Situations take `bosses` and `minions` rows (wiki Situation-payload).
   14. Audit, engine, through the helper:
       1. `inPlay`, `attackTargets`, `askRedirect`, `TargetMonster` choices.
       2. `EachMonsterTakesDamage`, `monsterDeath`, `slotAt`, situation stack.
   15. Audit, engine, monster slots only:
       1. `fillSlots`, `refillSlots`, deck attack placement, `PlaceFromOutside`.
       2. `CoverMonsterSlot`, `PlaceFromDeck`, `PlaceFromDiscard`, `FromUnder`.
       3. `PutUnder`, Indomitable, `ExpandMonsters`, situation `monsters`.
   16. Audit, cards: 75 files in b2, g2, fsp2, r, cardkit, enginetest.
       1. Most go through `cardkit.MonstersInPlay`, now the engine helper.
       2. Changed: cardkit `MonstersInPlay`, Ehwaz, Magic Marker, r `cardsInPlay`.
       3. Moves that `Held` stops: Ehwaz, Spider Mod, Void, I Am Error, Dogma.
       4. Also: Emergency Contact, Flush, the room discard effects.
       5. Kills that become 3 damage: Chaos Card, XIII Death, Tech X, The Bible.
       6. Also: Pandora's Box, Euthanasia, Soulbond, Plan C, Mother's Shadow.
       7. Slots only: The D10, Key, Red Key, Peep, Rainmaker, Dad's Note.
       8. Also: Friendly Ball, Ouija, g2 and fsp2 `freeSlots`, expand cards.
       9. Ace of Diamonds counts monster slots only (R-CHAL-02, R-CHAL-07).
       10. Rag Man, Delirium, `backOnTop` move a dead card: unchanged.
   17. Open: a max HP drop (R-MECH-21) can still bring a boss to 0 HP.
   18. Tests: `cards/challenge_slots_test.go`, R-CHAL-02 – R-CHAL-12.
   19. Fuzz: `FuzzChallengeSets`, `make fuzz-challenge`; 3 minutes green.
   20. ADR 028 got small drift items: Decision 7.9–7.13.

### [x] 14.22 Engine: How the Krampus Stole Christmas

1. Goal: the Krampus challenge plays by its rules, all 3 difficulties.
2. Tasks:
   1. Setup: Pilfering Krampus of the difficulty in a final boss slot.
   2. Setup: 4 face-down treasure cards as presents, with counters.
   3. Start of each turn: the active player rolls; counters come off.
   4. Hard and Ultra Hard: extra rolls per minions in play.
   5. "What's in the Box?": a present at 0 is destroyed; minions come.
   6. "Happy Holidays!": presents flip into items; the active player gives them.
   7. "All I Want for Christmas": forgo a reward or a soul for counters.
   8. Krampus can't be attacked or take damage below the soul count or with a minion.
   9. Normal: attack a random minion once each turn if able.
   10. No rewards or souls from minions.
   11. Goal: Krampus dies; lost: all 4 presents destroyed, or the timer.
   12. Fuzz and sims run Krampus games at every difficulty.
3. Done when:
   1. Engine tests for each KRAMPUS rule and each open case's answer.
   2. Seeded whole games end in a win and in a loss.
4. Result (2026-11-10):
   1. `cards/challenge/krampus.go`: setup, shared parts; the card files build on it.
   2. Setup: the boss via `PutInFinalBossSlot`; 4 presents via `AddPresent`.
   3. Presents API in `challenge.go`; `Present.Destroyed`, `Present.Flipped`.
   4. New `CardDef` hooks: `Guarded`, `MustAttackRandom`, `NoSpoils`.
   5. New `Replacement.May`: optional, asked yes or no (`ChooseMayReplace`).
   6. `Action.Reward` marks a monster's reward parts, rolled ones too.
   7. Boss card: protection, start-of-turn roll, Normal's forced attack.
   8. R-KRAMPUS-14 is the boss card's too: the boss card prints it.
   9. Rule card: All I Want for Christmas, What's in the Box?, Happy Holidays!, the goal.
   10. Rule card abilities are the active player's (`abilityController`).
   11. Loss reason `event.challenge_lost.presents`.
   12. Events: `present_counters`, `present_destroyed`, `presents_flipped`.
   13. Also `present_given`, `forgone`, `minion_arrived`; history lines for each.
   14. `present_destroyed`, `presents_flipped` are Soft barriers: a card shows.
   15. View: presents show destroyed, flipped, face-up card, gainer.
   16. Decisions left open by the docs:
       1. Happy Holidays! triggers on a gained soul; it uses the stack.
       2. It flips first; then the active player picks, seeing the cards.
       3. What's in the Box? destroys at resolution, if still at 0.
       4. A forgone soul is a dead card's soul (Q-32, item 3.9).
       5. A deck attack picked away for the random one stays owed (Q-31).
   17. Tests: `cards/challenge/*_test.go`, one file per card, plus `krampus_test.go`.
   18. Whole games: seeds won at Normal and Hard; random ones lost.
   19. Fuzz: `FuzzChallengeSets` flags 1–3 play Krampus; `TestSimulations` too.
   20. ADR 028 got small drift items: Decision 5.12–5.14, 6.6–6.10, 9.5.

### [x] 14.23 Server: challenge games

1. Goal: the server hosts a challenge for 2 players.
2. Tasks:
   1. Protocol option: challenge ID and difficulty (ADR 028).
   2. A challenge table has exactly 2 seats; the server checks it.
   3. Bonus souls are off in a challenge, whatever the option says.
   4. Records, saves and the seamless restart keep the mode.
   5. The game list shows the challenge and its difficulty.
   6. A list of challenges the server knows, for the host screen.
3. Done when:
   1. Server tests: create, start, save, load and restart a Krampus game.
   2. UC-03 and UC-05 updated, or a new use case.
4. Result (2026-11-10):
   1. Protocol: `Options.challenge`, `Options.difficulty`; `GameInfo` has both.
   2. All `omitempty`; protocol version unchanged; golden lines appended.
   3. `checkChallenge`: unknown challenge, bad difficulty, seats other than 2.
   4. Message keys `error.no_challenge`, `error.no_difficulty`, `error.challenge_seats`.
   5. Bonus souls forced off; a host ban of a challenge card is refused.
   6. The table's set list stays the playable sets; `challenge` is never shown.
   7. `gameSets` adds the challenge set last; `engineSetup` builds every setup.
   8. Used by start, restart, the option checks; saves and rollback follow.
   9. `cards.ChallengeSet(def)`: the set of a challenge.
   10. Record header: `challenge`, `difficulty` (ADR 007, items 2.17–2.19, 5.7).
   11. A lost challenge: record `end` finished, no winners; the room is over.
   12. Challenge list: bridge `App.Challenges()`, like `CardSets` (ADR 028, 10.7).
   13. Lobby: `challengeText` in the game list and the table summary; English.
   14. Tests: `internal/server/challenge_test.go`; record and app tests.
   15. Test helper `replaySetup` rebuilds a setup from a record header.

### [x] 14.24 Client: game mode in the host screen

1. Goal: the host picks "Regular game" or "Challenge".
2. Tasks:
   1. A "Game mode" select: Regular game (default), Challenge.
   2. Challenge: a select of challenges beside it; now only Krampus.
   3. Challenge: a difficulty select: Normal, Hard, Ultra Hard.
   4. Challenge: seats fixed at 2; bonus souls hidden.
   5. The host screen remembers the last mode, challenge and difficulty.
   6. English and Russian texts.
   7. Other detailed options stay: sets, picking, rooms, timer, rollback, bans.
3. Done when:
   1. e2e test: host a Krampus challenge; the game list shows it.
   2. Use case for hosting a challenge.
4. Result (2026-11-10):
   1. "Game mode" in the simple setup, above Players: Regular game, Challenge.
   2. Challenge: challenge and difficulty selects beside it, from `getChallenges()`.
   3. Players show 2, locked, with a line under it; Bonus souls hidden.
   4. Every other detailed option stays; the host screen has no ban list yet.
   5. `hostSetup.ts`: `pickChallenge`, `createOptions`; Vitest tests.
   6. A challenge sends `challenge` and `difficulty` without Detailed too.
   7. It sends no `no_bonus_souls`; the server plays without them.
   8. Host settings: `mode`, `challenge`, `difficulty`; Go, TS mirror, fake bridge.
   9. An old settings file loads Regular game; Go test.
   10. Item 6: English only; the interface is English only (16.6).
   11. Selects have `aria-label`s: "Challenge" is also a mode option.
   12. The joiner's game list and table line name the challenge (14.23).
   13. e2e `challenge-host.spec.ts`: Krampus at Hard; a regular game unchanged.
   14. `startGame` hooks `listed` and `seated` check a guest before Join and Ready.
   15. New UC-41; UC-03 and UC-05 point to it; wiki Client.

### [x] 14.25 Client: the challenge table

1. Goal: both players see and play the challenge.
2. Tasks:
   1. Final boss slot beside the monster slots; aiming works on it.
   2. Minion slots appear and vanish with their minions.
   3. Presents: face-down cards with their counters; flipped ones are items.
   4. Timer: a d8 that counts down, with a tooltip.
   5. The challenge's rule card for the chosen difficulty, zoomable.
   6. "All I Want for Christmas" choice in a dialog.
   7. Game over: "Challenge won" or "Challenge lost", with the reason.
   8. History lines for presents, minions, the timer.
   9. Animations and window layers follow ADR 012 and ADR 015.
   10. Without card images: text cards for the boss and the presents.
3. Done when:
   1. e2e tests: presents lose counters; minions come; Krampus is attacked.
   2. Screenshots for review; use case for playing a challenge.
4. Result (2026-11-11):
   1. Board `.has-challenge`: three rows and a challenge column, one shape.
   2. Final boss group left of the monster slots; landscape, aimable.
   3. Minion row under the monsters; small cards; empty, it shows nothing.
   4. Challenge panel in the bonus souls' column: rule card, timer, presents.
   5. Timer: an SVG d8 with the rounds left; "Rounds left"; tooltip.
   6. Presents: treasure back, counter badge, kind under it.
   7. A destroyed present's place shows ✗; a given one's ✓.
   8. `landscape(type)` in `cards.ts`: rooms, boss and rule cards, text too.
   9. New `challenge.ts`: `monsterRows`, `presents`, `timerText`, `challengeEnd`.
   10. `findCard`, `hpOf`, `slotZone` know bosses and minions.
   11. Game over: "Challenge won!" or "Challenge lost" with the reason.
   12. Effects reuse ADR 012 animations; no new one, no new layer:
       1. A new minion: `card_revealed` places it; the row takes it first.
       2. `present_counters`: a popup; `present_destroyed`, `present_given`: a fly.
       3. `challenge_won`, `challenge_lost`: the end banner.
   13. Choices needed nothing new: forgo and gifts use the choice dialog.
   14. `ask.attack_random` shows only beside another owed attack.
       1. With none, the engine attacks a random minion at once.
   15. History lines already named their cards; checked.
   16. Fake bridge and e2e `cardName` know the challenge cards.
   17. Tests: Vitest `challenge.test.ts`, `playback.test.ts`, `cards.test.ts`.
   18. e2e `challenge-table.spec.ts`: 4 tests, saved games for exact states.
   19. UC-42; ADR 028 items 11.7–11.11; wiki Client and Animations.

### [ ] 14.26 Challenge play-check and docs

1. Goal: whole Krampus games run in the app; docs match.
2. Tasks:
   1. `make e2e-games` also plays Krampus at each difficulty.
   2. Wiki: a players' page for challenges; Home and sidebar.
   3. README: one line about the challenge mode.
   4. `docs/ideas.md` GS-11: what is done, what is left.
3. Done when:
   1. `make e2e-games` passes with the Krampus games.

### [ ] 14.27 Solo challenges

1. Goal: one person plays both characters (solitaire).
2. Not now (owner, 2026-11-10): co-op comes first.
3. Open: one app holds both seats; the table follows the waiting seat.

### [ ] 14.28 Krampus: competitive variant

1. Goal: 2–4 players race to kill Pilfering Krampus.
2. Not now (owner, 2026-11-10).
3. Rules: personal present slots; the competitive boss card; first kill wins.
