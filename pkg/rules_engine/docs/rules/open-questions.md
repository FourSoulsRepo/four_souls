# Open questions

Cases the sources do not settle, or old rulings to re-check (R-03).
Each case: question, sources, options, status.
Resolved cases later become situation payload tests (step 4.11).

Status: **Owner** = needs the owner's decision; **Resolved**; **Parked**.

Basis of an answer (see "Rulings to revisit" at the end):
1. **Official**: the rulebook (S-OFF), the designers or the rules team.
2. **Community**: only players answered (Discord, Reddit); we follow them.
3. **House**: no ruling found, or the sources disagree; the owner chose.

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
7. Basis: **Official**: Yuggy and Edmund.

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
7. Basis: **Official**: S-OFF and Yuggy.

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
7. Basis: **House**: the sources disagree (Ed and S-RU against Yuggy as quoted).

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
6. Basis: **House**: no ruling found; the engine keeps (a) meanwhile.

### Q-11 PHD: once per roll

1. Question: "You may add or subtract 1 from any of your non-attack rolls."
   1. A changed roll lets "would roll" effects look again (R-DICE-05).
   2. May PHD change the same roll again, 2 → 3 → 4?
2. Sources: no ruling in S-OFF, S-YUG, S-DIS, Reddit, X, or the web (2026-10-13).
3. Options: (a) once per roll; (b) any number of times.
4. Status: **Resolved** (2026-10-13): (a); (b) looks like a bug (R-DICE-10).
5. Basis: **House**: no ruling found.

### Q-12 Euthanasia: which rolls are "in a row"

1. Question: "Each time you roll the same result twice in a row on an attack roll on the same turn".
   1. Must both rolls be attack rolls, or any roll then an attack roll?
2. Sources:
   1. No ruling on this point in S-OFF, S-YUG, Reddit, X, or the web (2026-10-13).
   2. S-DIS: it works on bosses; a rerolled second roll counts by its final value.
   3. S-DIS: it triggers before combat damage; you must control it as the roll resolves.
3. Options: (a) two attack rolls in a row; (b) any roll, then an attack roll (as R-DICE-08).
4. Engine now: (a), as R-DICE-09.
5. Status: **Owner**, open (2026-10-13); keep (a) until a ruling is found.
6. FAQ: explain the difference from Fruitcake (R-DICE-08) in the rules FAQ (11.7).
7. Basis: **House**: no ruling found; the engine keeps (a) meanwhile.

### Q-13 Damage loops at 0 HP

1. Question: may a player at 0 HP loop "when you take damage" triggers?
   1. Example: Forever Alone with Gurdy Jr., while their death waits on the stack.
   2. Each activation deals 1 damage; the damage recharges Forever Alone.
2. Sources:
   1. S-OFF §559: an object at 0 HP still takes damage.
   2. S-X: Edmund, dead players still take damage ("beating a carcass").
   3. S-YUG: "take 1 damage" costs could be used over and over when dead.
   4. BGG 2019: Blue Baby may respond to lethal damage with Forever Alone.
   5. BGG 2023, The Maestro Team: damage to 0 HP still fires damage triggers.
   6. BGG 2019, designer Daniel Cochran: breaking the game with loops is fine.
   7. No ruling on this exact combo in S-RED, S-DIS, S-X, or BGG (2026-10-13).
3. fouronlinesouls.net (2026-10-13):
   1. The death goes on the stack above the recharge trigger.
   2. The death resolves first; the penalty deactivates Forever Alone.
4. Options: (a) official, the loop is legal; (b) death above the triggers;
   (c) house rule: damage at 0 HP fires no "take damage" triggers.
5. Engine now: (a).
6. Status: **Resolved** (2026-10-13): (a) by default; (c) as a match option (GS-12, 14.8).
7. Basis: **Official** for (a); (c) is a house option.

### Q-14 Magic Marker: which numbers

1. Question: "Change a number in the effect text of a card in play or loot being played by 1".
   1. Cards are Go code; the engine does not know every number in the text.
2. Options: (a) only numbers in common blocks; (b) every card lists its numbers.
3. Status: **Resolved** (2026-10-19): (a) now (R-ABIL-39); (b) is future work (CD-13).
4. Basis: **House**: an engine limit, not a rules question.

### Q-15 Lemegeton used more than once

1. Question: after two uses, must two different items be destroyed?
   1. "The next time your turn ends, destroy a non-eternal item you control."
2. Sources:
   1. S-DIS (2024–2026): the same item may be chosen each time; the rest do nothing.
   2. Reddit: "every activation destroys an item".
3. Our rules: the item is a target, chosen as each trigger goes on the stack (R-ABIL-05).
   1. All triggers may target one item; the first destroys it, the others fizzle.
4. Options: (a) as our rules and S-DIS; (b) each trigger destroys a different item.
5. Status: **Resolved** (2026-10-19): (a), rules as written.
6. Basis: **Community**: S-DIS and the rules as written.

### Q-16 Spindown Dice on a Bag Lunch without counters

1. Question: "Destroy an item. If it was controlled by a player, they may steal an item from the shop."
   1. Bag Lunch replaces its destruction with a counter.
   2. May its controller still steal from the shop?
2. Sources: S-DIS split; no ruling in S-OFF, S-YUG, Reddit, X (2026-10-19).
3. Options: (a) yes, the second sentence is not "if you do" (R-MECH-02); (b) no, nothing was destroyed.
4. The English print says "Destroy an item."; an earlier text said "a non-eternal item".
5. The Russian print: "a non-eternal item"; it has no "if you do" either.
6. Status: **Resolved** (2026-10-19): (b); the engine follows the Russian print.
   1. Only a non-eternal item can be chosen; Bag Lunch can.
   2. The steal needs the destruction; Bag Lunch's counter stops it.
   3. The same as Berkano (Q-19).
7. Basis: **House**: no ruling found; the Russian print decides.

### Q-17 Rib of Greed and stealing

1. Question: "Abilities and the death penalty can't make you discard loot cards or lose ¢."
   1. Does it stop another player stealing your ¢ (Magnet)?
2. Sources: Reddit asked Yuggy; no answer found (2026-10-19).
3. Our rules: stealing is giving, not losing (R-MECH-35, R-MECH-38).
4. S-OFF: "lose" moves ¢ to the game's pool; a steal is the other player giving ¢.
5. Options: (a) steals still work; (b) Rib of Greed stops them too.
6. Status: **Resolved** (2026-10-19): (a).
7. Basis: **House**: no ruling found; follows our rules on giving.

### Q-18 The Clicker: the new character's state

1. Question: The Clicker replaces your character; is the new one charged?
2. Sources:
   1. S-DIS: cards enter play charged; the player is not healed.
   2. Damage stays with the player (R-MECH-19).
3. Options: (a) charged, as any card entering play; (b) keeps the old character's state.
4. Status: **Resolved** (2026-10-19): (a).
5. Basis: **Community**: S-DIS.

### Q-19 Berkano: a steal without a destruction

1. Question: "Destroy an item. If that item was controlled by a player, they steal an item from the shop."
   1. If the item is not destroyed (eternal, Bag Lunch), do they still steal?
2. Sources:
   1. English and Russian prints: no "if you do"; an earlier English text said "non-eternal".
   2. S-DIS (2025-04-11): an eternal item may be chosen; nothing is destroyed, nobody steals.
   3. No ruling in S-OFF, S-YUG, Reddit, X, or the web (2026-10-19).
3. Options: (a) they steal anyway; (b) only if the item was destroyed.
4. Status: **Resolved** (2026-10-19): (b), as Spindown Dice (Q-16).
5. Basis: **Community**: one S-DIS answer.

### Q-20 Social Goals: whose progress counts

1. Question: "Note each goal as players complete them."
   1. Does every player add to one tally, or must one player reach a goal alone?
2. Sources:
   1. English and Russian prints say the same.
   2. No ruling in S-OFF, S-YUG, S-DIS, or the web (2026-10-19).
3. Options: (a) one shared tally; (b) one player alone; (c) one player in one turn.
4. Status: **Resolved** (2026-10-19): (a), while the room is in play.
5. Basis: **House**: no ruling found.

### Q-21 Pity for the Poor: how many free purchases

1. Question: "may purchase a shop item for 0¢ on their turn."
2. Sources:
   1. English and Russian prints say the same.
   2. No ruling found (2026-10-19).
3. Options: (a) every purchase costs 0¢; (b) the first each turn; (c) an extra purchase.
4. Status: **Resolved** (2026-10-19): (a); the number of purchases is unchanged.
5. Basis: **House**: no ruling found.

### Q-22 Cards put on the bottom: in what order

1. Question: who orders several cards put on a deck's bottom?
2. Sources:
   1. Most cards don't say; Black Candle says "in any order".
   2. Ouija Board says "in a random order".
   3. No card looks at or takes from the bottom.
   4. The order matters only when the deck runs out.
   5. No ruling found (2026-10-19).
3. Options: (a) the player orders them; (b) always random.
4. Status: **Resolved** (2026-10-19): (b), even for Black Candle (19.27).
5. Basis: **House**: questionable; the owner saw no use for the prompt.

### Q-23 A soul that can't be destroyed

1. Question: what does "destroy a soul you control" do to an Eternal soul?
   1. Example: The Bone became a soul; it was a starting item, so Eternal.
   2. XX Judgement offered it; picking it destroyed nothing (19.36).
2. Sources:
   1. R-ABIL-18: an Eternal object can't be destroyed.
   2. The Bone's soul "loses all abilities".
   3. No ruling found on Eternal souls (2026-10-20).
3. Options:
   1. (a) A soul with no abilities is not Eternal.
   2. (b) Souls that can't be destroyed are not offered.
4. Status: **Resolved** (2026-10-20): both (a) and (b).
   1. The Lost's own soul is a character soul; it is never offered.
5. Basis: **House**: no ruling found.

### Q-24 "An attack on a monster" after a deck attack

1. Question: does "cancel an attack on a monster" work on a deck attack?
   1. Cards: Emergency Contact, IBS.
   2. Case: the deck attack revealed The Fallen (19.62).
2. Sources:
   1. S-OFF Attacking: abilities about an attack *declared* on a monster.
   2. They ignore an attack declared on the deck, even after a reveal.
   3. Neither card says "declared".
   4. No ruling found: Yuggy, Reddit, Discord, X (2026-10-26).
3. Options:
   1. (a) No: any attack on the deck is not an attack on a monster.
   2. (b) Yes: once a monster is revealed, the attack is on it.
4. Status: **Resolved** (2026-10-26): (b).
   1. R-ATK-05 now names "declared".
5. Basis: **House**: no ruling found.

### Q-25 Deactivating a passive item

1. Question: can "deactivate an item" pick a passive item?
   1. Cards: Hourglass, Curse of Impulse.
   2. Case: Hourglass deactivated Dad's Lost Coin (19.64).
2. Sources:
   1. R-ABIL-10: a ↷ ability costs deactivating its object.
   2. R-DEATH-14 deactivates only what has a ↷ ability.
   3. No rule says a passive item can't be deactivated.
   4. No ruling found: rulebook, Yuggy, Reddit (2026-10-26).
3. Options:
   1. (a) Any charged item; a passive one loses nothing.
   2. (b) Only items with a ↷ ability.
4. Status: **Resolved** (2026-10-26): (b), the owner's reading.
   1. R-MECH-54 states it.
5. Basis: **House**: no ruling found.

### Q-26 A roll changed to the value it has

1. Question: is a roll set to its own value changed?
   1. Cards: Godhead, I. The Magician, Libra, PHD, Telepathy for Dummies.
   2. Case: a 1 set to 1 while Dad's Lost Coin was declined (19.66).
2. Sources:
   1. R-DICE-05: a change sends the roll back to R-DICE-04.
   2. No rule says what a change to the same value is.
3. Options:
   1. (a) No change: "would roll" triggers do not look again.
   2. (b) A change: they look again and ask again.
4. Status: **Resolved** (2026-10-26): (a), the owner's choice.
   1. R-DICE-13 states it.
   2. Rewards use "when you roll" and fire once either way (R-DICE-06).
5. Basis: **House**: no ruling found.

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

Requiem searched 2026-10-19 (14.2): 96 of his posts and tweets name a Requiem card; rows 36 on. Tests come with the cards in 14.3.

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
| 25 | Mom's Heart, It Lives! | It triggers "when a player gains this soul" | Four Souls+; test |
| 26 | Mom's Eye Shadow | Only attacks declared on a monster; attacks on the monster deck go through | Four Souls+; test |
| 27 | Mutant Spider | The dice are rolled at once; only the chosen one is a roll | Four Souls+ (R-DICE-11); test |
| 28 | Telepathy for Dummies | Roll at the start of the turn; the next roll may become that result | Four Souls+; test |
| 29 | Distant Admiration | The given loot is played for free: no loot play is used | Four Souls+; test |
| 30 | 1-Up, Eden | An eternal 1-Up can't be destroyed, so it does nothing | Four Souls+; test |
| 31 | Krampus, a forced deck attack | The active player picks which forced attack to make (R-ATK-03) | Four Souls+; test |
| 32 | Cursed Eye, Bomb | A bomb in response to the attack roll ends combat before Cursed Eye triggers | Four Souls+; test |
| 33 | Headless Horseman | It prevents its death and heals; the attack goes on | Four Souls+; test |
| 34 | Hourglass and circled numbers | A circled number triggers on any roll that resolves as it | Four Souls+ (R-DICE-06) |
| 35 | Abaddon | Printed +3 ATK; the site text lacks the number (fixed by an override) | Four Souls+ |
| 36 | Experimental Treatment, Pandora's Box | An item that becomes a soul keeps its abilities | Requiem (R-CARD-38) |
| 37 | Experimental Treatment, Butter Bean | "Destroy this:" is a cost: a cancelled use still destroys it | Requiem |
| 38 | Punching Bag | Killed means destroyed; the active player gets its rewards | Requiem (R-MECH-05) |
| 39 | Rag Mega | "Up to 3" may be none | Requiem (R-ABIL-37) |
| 40 | Rag Mega, Bomb | A loot already on the stack is not discarded by it | Requiem |
| 41 | The Harbingers | One attack each turn must go there; an attack already used means no attack on it | Requiem (R-ATK-22) |
| 42 | Blind Rage, Krampus | The attacker picks which forced attack to make | Requiem (R-ATK-22) |
| 43 | Soul of Wrath | A dying active player gains it before the turn ends | Requiem |
| 44 | Cursed Mulligan | It triggers on its own death; slots have no maximum | Requiem (R-ATK-23) |
| 45 | Jar of Flies | Shop items' triggers don't work in the shop | Requiem (R-CARD-05) |
| 46 | Holy Brain, Lodestone | The -1 DC comes after the roll resolved; it does not save that roll | Requiem |
| 47 | Dice Room | The room step is once per turn, after end-of-turn triggers | Requiem (R-TURN-17, R-TURN-18) |
| 48 | Dice Room | Additional attacks gained in the end phase are lost | Requiem (R-TURN-19) |
| 49 | Devil Beggar | Only the active player uses a room's activated abilities | Requiem (R-CARD-36) |
| 50 | Battle Royale!!!!!, Curse of the Suspicious | The newest attack-a-player ability sets the evasion | Requiem (R-MECH-06) |
| 51 | Eternal D6, Nightmare Tick | A replacement on a roll applies after "would roll" triggers | Requiem (R-DICE-12) |
| 52 | Glitch, Magic Skin | Gained abilities keep Glitch eternal; it can still be given away | Requiem (R-ABIL-38) |
| 53 | Glitch, Mama Mega | A gained "Destroy this. If you do" does nothing | Requiem (R-ABIL-38) |
| 54 | Glitch, Keeper's Sack | Gaining abilities is not gaining the item: no 14¢ | Requiem |
| 55 | Keeper's Sack, Keeper's Penny | Giving is not gaining: no ¢ when passed on | Requiem (R-MECH-35) |
| 56 | Glitch, Sumptorium | Counters stay on the object when the copied ability ends | Requiem (R-MECH-12) |
| 57 | Damocles | A bought shop item is not a treasure; the deck top is | Requiem (R-CARD-37) |
| 58 | The Clicker, Curse of the Soulless | Becoming The Lost is not gaining a soul | Requiem |
| 59 | Soul of Envy, Curse of the Soulless | The cursed player can't gain it; it is then never gained | Requiem |
| 60 | Ultra Greed, Sisters Vis | Extra attack rolls come one at a time; a helper's death doesn't end the attack | Requiem (R-ATK-21) |
| 61 | Sibling Rivalry | It may answer the attack declaration; the target is chosen after | Requiem |
| 62 | Magic Skin | A destroyed item's ability on the stack still resolves | Requiem (R-ABIL-03) |
| 63 | Magic Skin | Its start-of-turn trigger is on the stack before priority; giving it away then doesn't change the target | Requiem (R-ABIL-34) |
| 64 | Anima Sola | It must be recharged before it flips again | Requiem |
| 65 | The Mirror | A player at 1 HP pays 0 HP and still gets the bonus | Requiem |
| 66 | MOTHER! | Evasion never goes above 6 | Requiem (R-ATK-18) |
| 67 | The Fettered (Dead Weight), Resurrected Rag Man | Replacement effects are not triggers; Dead Weight can't cancel them | Requiem |
| 68 | Friendly Ball | It may hold an Indomitable monster | Requiem |
| 69 | Halo of Flies | Its 6 also stops damage from a player attacking you | Requiem (R-ATK-24) |

## Card rulings from Discord (S-DIS)

Searched 2026-10-12 (8.7) with each card's Russian name from foursouls.ru (the owner's database migrations, not in the repo): 188 Base Game cards are named in the channel. Read: every question to "@Судья" (the judges) that names one, 87 threads. Medium trust; kept where it agrees with S-OFF or fills a gap.

Requiem searched 2026-10-19 (14.2) with the Russian names of the owner's Requiem migrations: 77 judge threads; rows 12 on. Reddit (digest, 63 threads) agreed with these rows or with S-OFF.

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
| 10 | Euthanasia | It works on bosses; a rerolled second roll counts by its final value | Four Souls+; test |
| 11 | Euthanasia | It triggers before combat damage; you must control it as the roll resolves | Four Souls+; test |
| 12 | Lemegeton | Used twice, the same item may be chosen twice | Requiem (Q-15) |
| 13 | Blind Rage, The Revenant | The monster deck must be attacked first; The Revenant gives an extra attack after | Requiem |
| 14 | Betrayal | Only a non-eternal item can be stolen (English text wins) | Requiem |
| 15 | All Hallows' Eve, The Lamb | Stealing a soul is not gaining it | Requiem (R-MECH-35) |
| 16 | Gello | It can't choose a monster at 0 HP; a bomb in response makes it fizzle | Requiem |
| 17 | Damocles, Dice Room | Rerolling 6 items gains 12 treasures: each reroll is its own gain | Requiem |
| 18 | Spindown Dice, Bag Lunch | Split answers; the owner follows the Russian print: no steal | Requiem (Q-16) |
| 20 | Berkano | An eternal item may be chosen; not destroyed, so nobody steals | Requiem (Q-19) |
| 19 | The Clicker | Damage stays with the player; the new character is not healed | Requiem (Q-18) |

## Rulings to revisit

House and community answers, made for this game.
They may change when an official ruling appears.
Check here first when a card plays oddly.
Every row of "Card rulings from Discord (S-DIS)" is Community too.

| Q | Cards | The engine now | Basis |
|---|---|---|---|
| Q-09 | Two of Clubs | A second copy does not double again | House |
| Q-10 | Fruitcake | Each player's own rolls in a row | House (open) |
| Q-11 | PHD | Changes a roll once | House |
| Q-12 | Euthanasia | Two attack rolls in a row | House (open) |
| Q-14 | Magic Marker | Only loot, ¢, damage and roll-table numbers | House |
| Q-15 | Lemegeton | Every trigger may pick the same item | Community |
| Q-16 | Spindown Dice | Non-eternal items only; steal only if destroyed | House |
| Q-17 | Rib of Greed | Steals still work | House |
| Q-18 | The Clicker | The new character enters charged | Community |
| Q-19 | Berkano | Steal only if the item was destroyed | Community |
| Q-20 | Social Goals | One shared tally for all players | House |
| Q-21 | Pity for the Poor | Every purchase costs 0¢ | House |
| Q-22 | Black Candle, The Moon, Flush!, Eden, … | Cards go to the bottom in random order | House (questionable) |
| Q-23 | The Bone, XX Judgement, Wizoob, … | No Eternal on a soul without abilities; Eternal souls not offered | House |
| Q-24 | Emergency Contact, IBS | Cancel a deck attack once a monster is revealed | House |
| Q-25 | Hourglass, Curse of Impulse | Only items with a ↷ ability are deactivated | House |
| Q-26 | Godhead, I. The Magician, Libra, PHD, … | A roll set to its own value is not changed | House |
