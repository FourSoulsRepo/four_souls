# Open questions

Cases the sources do not settle, or old rulings to re-check (R-03).
Each case: question, sources, options, status.
Resolved cases later become situation payload tests (step 4.11).

Status: **Owner** = needs the owner's decision; **Resolved**; **Parked**.

On 2026-10-09 the owner parked every open decision for later.

## Decisions for the owner

### Q-01 Mulligan

1. Question: may players redraw their starting hand?
2. Sources: no mulligan in S-OFF or S-RU; the TTS table has a house-rule button.
3. Options: (a) no mulligan, as official; (b) optional host setting.
4. Proposal: (a), per R-SETUP-13.
5. Status: **Parked** again by the owner (2026-10-09); the engine has no mulligan.

### Q-02 First player

1. Question: how is the first player chosen?
2. Sources: S-OFF suggests "the saddest player", or the lowest dice roll, or any fair method.
3. Proposal: every player rolls a D6; lowest goes first; reroll ties (GS-06).
4. Status: **Resolved** (2026-10-09): the lowest roll, as the engine does; Cain overrides it.

### Q-03 Two vs two teams (GS-05)

1. Question: how do teams play? The official rules have no team mode.
2. Open points:
   1. Turn order: alternate teams (A1, B1, A2, B2)?
   2. Win: combined souls of the team, and how many?
   3. Can teammates give each other items or loot (bartering forbids it)?
   4. Who gets rewards and souls when a teammate kills a monster?
3. Options: (a) design house rules; (b) park 2 vs 2 until later.
4. Status: **Parked** until after the MVP (2026-10-09).

### Q-04 Co-op timer on death

1. Question: does the co-op timer drop on any death, or only a death during that character's own turn?
2. Sources: S-OFF says "during their turn"; S-RU says on any death.
3. Proposal: follow S-OFF (trust order).
4. Status: **Parked** (co-op is future, GS-11).

### Q-05 Bartering in a digital game

1. Question: how do players barter (R-BARTER-01)?
2. Proposal: a "give ¢" action usable any time, even without priority; promises stay in chat or voice, unenforced (R-BARTER-04).
3. Status: **Resolved** (2026-10-09): no bartering for now; the engine keeps room for it (ADR 005).

### Q-06 Official variants as host options

1. Question: offer the official variants Mini-draft (R-SETUP-14) and Eden Only (R-SETUP-15) in detailed setup (GS-10)?
2. Note: S-OFF recommends Mini-draft for 2-player games.
3. Status: **Resolved** (2026-10-09): neither variant for now; standard setup only.

### Q-07 Choosing a dead player

1. Question: may "choose a player" pick a player who is already dead this turn?
2. Sources:
   1. Ed (2018, R-07 #4): you can damage dead players ("beating the corpse").
   2. Yuggy (S-YUG, Reddit etdouue, fmz4v53; Ed's tweet 1085299878199062528): you can't choose a dead player; only damage that picks no one (all players, the player to your left) hits them. The Isaac boss then makes you damage yourself.
3. Engine now: dead players stay valid choices (R-07 #4).
4. Options: (a) keep; (b) dead players can't be chosen when a living one could be.
5. Proposal: (b), with the source note; R-07 #4 then means untargeted damage only.
6. Status: **Resolved** (2026-10-12): (b), as Yuggy rules (R-ABIL-35).

### Q-08 The 3+ player symbol

1. Question: R-CARD-32 removes 3+ player cards in 2-player games; which Base Game cards have the symbol?
2. Sources: S-OFF (3+ Player Only Symbol); Yuggy (Reddit kh7ze7n, jdlljs9): no loot card has it; Satan does.
3. Gap (closed): card_db now has `three_plus`, read by the scraper from foursouls.com.
4. Proposal: the scraper in `FourSoulsRepo/tools` reads the symbol into card_db; the engine then removes those cards in 2-player games.
5. Status: **Resolved** (2026-10-12): (proposal) R-CARD-32 applies.
6. Base Game V2 cards with the symbol (11):
   1. Monsters: Dople, Evil Twin, Rage Creep, Satan.
   2. Treasure: Brimstone, Curse of the Tower, Golden Razor Blade, Monster Manual.
   3. Treasure: Razor Blade, Remote Detonator, Trinity Shield.

### Q-09 Two of Clubs twice

1. Question: does a second Two of Clubs this turn double the loot again?
2. Sources:
   1. No: Ed; the S-RU FAQ (row 1 below); the tournament server's practice (S-DIS, 2024-05-30).
   2. Yes: Yuggy and another rules lead, as quoted on that server.
3. Engine now: no; a replacement applies once per event (R-ABIL-32).
4. Proposal: keep "no" until Yuggy's ruling is found in writing.
5. Status: **Resolved** (2026-10-12): (a) no stacking (R-ABIL-36).
   1. The S-RU FAQ says the wording matters; the card doubles once.
6. FAQ: explain this case in the rules FAQ (11.7); players often ask.

### Q-10 Fruitcake: whose rolls are "in a row"

1. Question: "Each time a player rolls the same result twice in a row on the same turn".
   1. Does it look at each player's own rolls, or at the last two rolls of anyone?
2. Sources:
   1. No ruling in S-OFF, S-YUG, Reddit, X, or the web (searched 2026-10-13).
   2. S-DIS (2023-12-06): each pair counts; three in a row recharge twice.
   3. S-DIS: of Mutant Spider's dice only the chosen result counts.
   4. S-YUG on Euthanasia: final results after changes count.
3. Options: (a) per player; (b) any two rolls in a row.
4. Engine now: (a), as R-DICE-08.
5. Status: **Owner**, open (2026-10-13); keep (a) until a ruling is found.

## Old rulings re-checked (R-07)

| # | Ruling (Ed's Notes, 2018–2019) | Check against current rules and V2 text | Status |
|---|---|---|---|
| 1 | Monster Manual can force a second attack | V2 text: "must attack that monster this turn if able"; it grants no extra attack | Superseded |
| 2 | You cannot kill dead players | A player dies at most once per turn (R-DEATH-17) | Holds |
| 3 | Copying an item is like getting it new | Copies take abilities but not counters or state (R-MECH-08) | Holds |
| 4 | You can damage dead players | 0-HP objects take damage that is not marked (R-MECH-16); only untargeted damage, see Q-07 (R-ABIL-35) | Narrowed |
| 5 | Eden can have Glass Cannon as Eternal | Eden's V2 text gives the chosen item eternal; eternal is not destroyed | Holds |
| 6 | Steamy Sale only affects shop items | V2 text says "shop items" | Holds (now in text) |
| 7 | Items recharge before the loot draw; the draw can be answered | R-TURN-02 to R-TURN-04 | Holds |
| 8 | "Prevent damage" lasts until end of turn | Prevent is a replacement effect; duration comes from the card text (R-ABIL-27) | Superseded |
| 9 | Cards can be added to the stack while it resolves | Priority passes after each item (R-STACK-04); no interrupting one item (R-STACK-05) | Holds |
| 10 | Death cancels combat | R-DEATH-12, R-ATK-15 | Holds |
| 11 | Bombs are not combat | Combat damage comes only from attack rolls (R-ATK-12) | Holds |

Tests: `cards/b2/interactions_test.go` covers rows 1–6, 10 and 11; rows 7–9 are engine tests.

## Card rulings from the Russian FAQ (S-RU)

Candidates for situation tests; check each against the English rules.

| # | Card | Ruling (summary) | Status |
|---|---|---|---|
| 1 | Two of Clubs | Activating it twice still only doubles; a replacement applies once per event (R-ABIL-32) | Holds (Q-09, R-ABIL-36); tests |
| 2 | Ambush! | Attacks on the monster deck made earlier this turn count toward its 2; the rest are additional attacks | **Parked** (2026-10-09); the engine forces 2 more |
| 3 | Multi-dice "choose one result" | Only the chosen result affects the game | Consistent; test later (later sets) |
| 4 | Mulliboom, The Lamb | ATK is printed with "!" (4!, 6!); card_db has no meaning for it; the engine uses the number | **Resolved**: "!" is only design (owner, 2026-10-09) |
| 5 | Devil Deal | card_db lists five options; the engine reads three: discard / loot 2, take 1 / take 2, search a Guppy | **Resolved**: three choices (owner, 2026-10-09) |

Tests: FAQ row 1 is in `cards/b2/interactions_test.go`, with Yuggy's "above 6 counts as 6" ruling.

## Card rulings from Yuggy (S-YUG)

Searched 2026-10-12 (8.7): his r/FourSouls comments and tweets that name a Base Game card, 294 questions. Rulings on V1 wording that V2 changed (Placebo, Butter Bean, Evil Twin, Eve, Judas) are left out. Tests: `cards/b2/yuggy_test.go`.

| # | Cards | Ruling (our words) | Status |
|---|---|---|---|
| 1 | Chaos Card, XIII. Death | Death clears only unresolved attack rolls; an item ability still resolves after its controller dies | Holds; test |
| 2 | Chaos Card, Butter Bean | Cancelling "Destroy this. If you do" leaves the item tapped, not destroyed | Holds; test |
| 3 | Chaos Card, Glass Cannon | Destroyed in response, "Destroy this. If you do" fizzles; a destroyed item's other abilities still resolve | Holds; test |
| 4 | Mom's Shovel | Usable with no soul to steal; it is destroyed anyway | Holds; test |
| 5 | Eden, Chaos Card | An eternal "Destroy this. If you do" item does nothing | Holds; test |
| 6 | Guppy's Paw, Fanny Pack | Paying HP is not damage; HP you don't have can't be paid | Holds; tests |
| 7 | The Dead Cat, Fanny Pack | Prevented damage is not taken: "when you take damage" doesn't trigger | Holds; test |
| 8 | Yum Heart | V2 prevents one whole instance of damage | Holds; test |
| 9 | Curved Horn, The D6 | A reroll is the same roll: the first attack roll keeps its +1 | Holds; test |
| 10 | Sacred Heart, Book of Belial | A roll changed down to 1 triggers "would roll a 1" | Holds; test |
| 11 | Sacred Heart, Dad's Lost Coin | Same-time triggers stack in turn order; a "would roll" trigger fizzles if the roll changed (R-DICE-07) | **Fixed** in the engine; test |
| 12 | Sacred Heart, Godhead | A roll changed to 1 triggers Sacred Heart again | Holds; test |
| 13 | Glass Cannon, The D20 | A target gone by resolution: Glass Cannon fizzles, stays tapped, no roll | Holds; test |
| 14 | Sloth, Bomb | The damage belongs to whoever played the loot: that player killed Sloth | Holds; test |
| 15 | Decoy, The Chest, Glass Cannon | A target swapped away is still the target; the soul goes to its new controller (R-ABIL-34) | **Fixed** in the engine; test |
| 16 | Modeling Clay | A copy of an active item stays tapped until recharged | Holds; test |
| 17 | Placebo, The D20 | Shop items count as in play for targeting (R-CARD-34) | **Fixed** in the engine; tests |
| 18 | Soul of Gluttony | A bonus soul is gained once per game, even if destroyed | Holds; test |
| 19 | Any "choose a player" | Dead players can't be chosen while a living one can (R-ABIL-35) | **Fixed** in the engine: Q-07; tests |
| 20 | 3+ player cards | Satan and others are out of 2-player games (R-CARD-32) | **Fixed** in the engine: Q-08; test |
| 21 | Steven, Bomb | Steven dies to a bomb before his "roll of 6" trigger: the roll is gone, no trigger | Gold Box; test |
| 22 | Holy Mantle, Bomb | Death ends combat: unresolved attack rolls are cancelled; a bomb on the stack still resolves | Gold Box; test |
| 23 | Plan C, Eden | An eternal "Destroy this. If you do" item does nothing (as row 5) | Gold Box; test |
| 24 | Euthanasia, 20/20 | "Same result twice in a row" uses final results | Four Souls+; test |

## Card rulings from Discord (S-DIS)

Searched 2026-10-12 (8.7) with each card's Russian name from foursouls.ru (the owner's database migrations, not in the repo): 188 Base Game cards are named in the channel. Read: every question to "@Судья" (the judges) that names one, 87 threads. Medium trust; kept where it agrees with S-OFF or fills a gap.

| # | Cards | Ruling (our words) | Status |
|---|---|---|---|
| 1 | Fanny Pack | A dead player still takes damage: nothing is marked, but "when you take damage" triggers (R-MECH-51, as S-OFF says) | **Fixed** in the engine; test |
| 2 | Two of Clubs | Stacking is disputed | **Resolved**: no stacking (Q-09); tests |
| 3 | Placebo | It uses one ↷ ability; it doesn't become a copy | Holds (V2 text) |
| 4 | Modeling Clay, ? Card | They can't copy cards in a discard or a deck | Holds |
| 5 | Diplopia | It may copy other players' items and shop items | Holds (R-CARD-34) |
| 6 | Any item | Items enter play charged unless they say otherwise | Holds |
| 7 | Butter Bean | It cancels a loot card entirely | Holds |
| 8 | Fruitcake | Each pair of equal rolls in a row recharges it: three in a row is two | Gold Box (Q-10); test |
| 9 | Fruitcake, Mutant Spider | Only Mutant Spider's chosen result counts as a roll | Gold Box (R-DICE-08) |
