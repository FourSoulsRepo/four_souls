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
5. Status: **Owner**.

### Q-02 First player

1. Question: how is the first player chosen?
2. Sources: S-OFF suggests "the saddest player", or the lowest dice roll, or any fair method.
3. Proposal: every player rolls a D6; lowest goes first; reroll ties (GS-06).
4. Status: **Owner**.

### Q-03 Two vs two teams (GS-05)

1. Question: how do teams play? The official rules have no team mode.
2. Open points:
   1. Turn order: alternate teams (A1, B1, A2, B2)?
   2. Win: combined souls of the team, and how many?
   3. Can teammates give each other items or loot (bartering forbids it)?
   4. Who gets rewards and souls when a teammate kills a monster?
3. Options: (a) design house rules; (b) park 2 vs 2 until later.
4. Status: **Owner**.

### Q-04 Co-op timer on death

1. Question: does the co-op timer drop on any death, or only a death during that character's own turn?
2. Sources: S-OFF says "during their turn"; S-RU says on any death.
3. Proposal: follow S-OFF (trust order).
4. Status: **Parked** (co-op is future, GS-11).

### Q-05 Bartering in a digital game

1. Question: how do players barter (R-BARTER-01)?
2. Proposal: a "give ¢" action usable any time, even without priority; promises stay in chat or voice, unenforced (R-BARTER-04).
3. Status: **Owner**.

### Q-06 Official variants as host options

1. Question: offer the official variants Mini-draft (R-SETUP-14) and Eden Only (R-SETUP-15) in detailed setup (GS-10)?
2. Note: S-OFF recommends Mini-draft for 2-player games.
3. Status: **Owner**.

## Old rulings re-checked (R-07)

| # | Ruling (Ed's Notes, 2018–2019) | Check against current rules and V2 text | Status |
|---|---|---|---|
| 1 | Monster Manual can force a second attack | V2 text: "must attack that monster this turn if able"; it grants no extra attack | Superseded |
| 2 | You cannot kill dead players | A player dies at most once per turn (R-DEATH-17) | Holds |
| 3 | Copying an item is like getting it new | Copies take abilities but not counters or state (R-MECH-08) | Holds |
| 4 | You can damage dead players | 0-HP objects take damage that is not marked (R-MECH-16) | Holds |
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
| 1 | Two of Clubs | Activating it twice still only doubles; a replacement applies once per event (R-ABIL-32) | Consistent; test later |
| 2 | Ambush! | Attacks on the monster deck made earlier this turn count toward its 2; the rest are additional attacks | **Owner** |
| 3 | Multi-dice "choose one result" | Only the chosen result affects the game | Consistent; test later (later sets) |
| 4 | Mulliboom, The Lamb | ATK is printed with "!" (4!, 6!); card_db has no meaning for it; the engine uses the number | **Owner** |
| 5 | Devil Deal | card_db lists five options; the engine reads three: discard / loot 2, take 1 / take 2, search a Guppy | **Owner** |

Tests: FAQ row 1 is in `cards/b2/interactions_test.go`, with Yuggy's "above 6 counts as 6" ruling.
