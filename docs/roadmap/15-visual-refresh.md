# Step 15. Visual refresh

Goal: a clearer table, using fouronlinesouls.net as a second reference.

Ideas: V-09, V-12, V-13, V-15 – V-31, V-33 – V-46, GS-06.

Order: 15.1, 15.15, 15.9 – 15.11, 15.16, 15.17, 15.38, 15.18 – 15.24,
15.3, 15.37, 15.25, 15.32 – 15.34, 15.26, 15.27, 15.59 – 15.63, 15.6, 15.64,
15.28 – 15.31, 15.35, 15.36, 15.39 – 15.46, 15.51, 15.52,
15.54 – 15.58.
Images come first, so the new table is built with backs and mats.

## Rules for this step (owner, 2026-10-12)

1. One sub-step at a time, in the order above.
2. Stop after every sub-step for the owner's review.
   1. Show screenshots of the work (`make screenshots`, `.cache/`).
   2. Never commit screenshots to the repo.
3. Start the next sub-step only after the owner approves.
4. These rules replace "stop after the global step" for step 15.
5. Screenshots show how the interface changed.
   1. Before and after shots of every changed screen.
   2. Owner reviews the changes and the screenshots.
   3. The review message links each screenshot file; no bare paths.

---

### [x] 15.1 Study the reference

1. Goal: a list of what to borrow (V-09).
2. Tasks:
   1. Open https://fouronlinesouls.net/ and start a game with a bot.
   2. Screenshots go to `.cache/screenshots/`, never into the repo.
   3. Compare with our table: layout, card sizes, panels, stack, prompts.
   4. List ideas to borrow; ideas only, no images, styles, or code.
   5. Compare with Hearthstone again; avoid the poker-table look (V-13).
   6. Find Hearthstone screenshots on the internet; compare (V-22).
   7. Study the stack window and its icon history strip (V-27).
   8. Study counters; use the cheat menu to get any card (V-27).
   9. Note how the cheat menu works, for the sandbox later (LR-01).
3. Done when:
   1. **Owner** picks the ideas to apply.
   2. Hearthstone screenshots compared; findings added below.
4. How the study went (2026-10-12)
   1. The site has no bot; "Add copy" gives a hot-seat second player.
   2. Played 2, 3 and 4 seats: loot, play, buy, attack, roll.
   3. Screenshots: `.cache/screenshots/reference/`.
   4. Hearthstone compared from memory; no game client here.
   5. Then web screenshots (V-22): `.cache/screenshots/reference/hearthstone/`.
   6. Sources: Wikipedia, interfaceingame.com, image search; ideas only.
   7. Talk to watch: GDC 2015, "Hearthstone: How to Create an Immersive UI".
5. Problems seen in our table
   1. O-01. 4 players: History panel covers the monster and treasure decks.
   2. O-02. 4 players: Stack panel covers the loot deck.
   3. O-03. Seat boxes with text stats on green felt read as poker (V-13).
6. Ideas to borrow: layout
   1. F-01. One framed central board: decks, shop, monsters, souls, stack.
   2. F-02. Nothing floats over the board; panels sit beside it.
   3. F-03. Stack sits inside the board, at its left side.
   4. F-04. Resolved stack items leave as icons in a strip by the stack.
   5. F-05. Own hand large at the bottom edge.
   6. F-06. Own character and items in one row above the hand.
   7. F-07. Opponent: name and coins, then character and items.
   8. F-08. Opponent hand: one card back with a count.
   9. F-09. Board grows for 3 players; opponents take side columns.
   10. F-10. Active player's name is highlighted.
7. Ideas to borrow: actions and feedback
   1. F-11. Turn actions as icons by own name: end turn, buy, attack.
   2. F-12. Icons light up when usable, dim when not.
   3. F-13. Icons change mode: "Attack" becomes "Roll"; "Buy" becomes "Cancel".
   4. F-14. Valid targets get number badges; keys 1–9 pick them.
   5. F-15. Space resolves or passes; E ends the turn.
   6. F-16. Hover zooms the card in place.
   7. F-17. Zoom of an unplayable card says why: "No loot plays left".
   8. F-18. A refused action shows its reason next to the clicked item.
   9. F-19. In a fight, attacker and monster both get a red frame.
   10. F-20. Choices open a dialog with big cards and a "View board" button.
   11. F-21. Last die value stays visible as a small die by the stack.
8. Things not to borrow
   1. N-01. One click plays a card at once; we keep drag and two clicks (V-16).
   2. N-02. No dice animation; we want V-17.
   3. N-03. No turn number, phase, or "waiting for" list; ours stay (V-05, V-11).
   4. N-04. HP and souls only on card art; V-18 wants them always visible.
   5. N-05. Their art, wood table, and fonts; never copy assets (V-09).
9. Hearthstone again (V-13)
   1. H-01. One shaped play area in the middle; no boxes around the edges.
   2. H-02. Hero portrait centered at own and opponent edge.
   3. H-03. Own hand fanned at the bottom, partly off screen.
   4. H-04. One big "End turn" button at the board's right edge.
   5. H-05. Drop the green felt and seat boxes; this removes the poker look.
   6. H-06. The board has one fixed shape and size, empty or full (V-23).
   7. H-07. History: a thin icon column on the board's left edge.
   8. H-08. Each history icon has a badge: sword for attack, bolt for a spell.
   9. H-09. Hover a history icon: the full card shows large beside it.
   10. H-10. An opponent's played card shows large at the left, then resolves.
   11. H-11. Choices: big cards over a dimmed board, nothing else on screen.
   12. H-12. Playable hand cards and ready attackers glow green.
   13. H-13. Hover a hand card: it rises large, keyword tips beside it.
   14. H-14. Damage: a starburst with "-3" on the hit target.
   15. H-15. Health and attack are big numbers on the portrait, not text.
   16. H-16. Decks sit at the right edge, seen from the side, with backs.
   17. H-17. Opponent hand: fanned card backs at the top, partly off screen.
   18. H-18. Mana: "3/8" and crystals at the bottom right; for us, coins.
   19. H-19. End turn: yellow on your turn, grey "Enemy turn" otherwise.
   20. H-20. End turn glows green when no actions are left.
   21. H-21. Player names sit small in the corners, outside the board.
   22. H-22. Status marks on cards: "Zzz" for not ready, a shield for Taunt.
   23. H-23. Aiming arrow: no screenshot found; study it in 15.4.
11. Stack on the reference (V-27)
   1. S-01. Stack panel inside the board; newest item on top.
   2. S-02. Each row: an icon and one short line of text.
   3. S-03. Names in color: player in their color, target in red.
   4. S-04. Icon by kind: card art, die for a roll, knife for damage.
   5. S-05. Each step is its own item: attack roll, then damage.
   6. S-06. Hover a row: the card shows large, owner name below.
   7. S-07. The used ability is highlighted on that large card.
   8. S-08. Resolve button sits in the panel; Space does the same.
   9. S-09. Empty stack: "Nothing on the stack…" and a grey button.
   10. S-10. History strip: resolved items as the same icons, newest on top.
   11. S-11. Hover a history icon: the same large card view.
12. Counters on the reference (V-27)
   1. C-01. A counter is a big green circle with a number, on the card art.
   2. C-02. The zoomed card shows the same circle.
   3. C-03. Card stats are live: monster HP on the card drops after damage.
   4. C-04. A tapped item turns grey; it is not rotated.
   5. C-05. A card with two abilities opens a dialog to pick one line.
   6. C-06. Zoom of an unusable card says why: "No effects usable now".
   7. Screenshots: `.cache/screenshots/reference/counters/`.
13. Cheat menu on the reference (V-27)
   1. K-01. Room option "Allow cheat options", on by default.
   2. K-02. Game menu, "Cheat view", adds cheat controls to the table.
   3. K-03. Decks get "Draw" and "Select"; monster deck gets "Put in slot".
   4. K-04. Every card gets an ×; coins get a +.
   5. K-05. "Select" opens a searchable picker of the whole deck.
   6. K-06. Search also matches card text, not only names.
   7. K-07. Every cheat is announced: "Ann used a cheat to gain 1 treasures".
   8. K-08. For us: a model for the sandbox (LR-01), not for normal games.
10. Owner notes (2026-10-12)
   1. F-14 rejected: no 1–9 badges on cards (V-29).
   2. New ideas from the review: V-22 – V-29.
   3. Owner approves the Hearthstone ideas, H-01 – H-23.
   4. Owner picks F (without F-14), S, C, and K; planned in 15.2 – 15.14.

### [~] 15.2 Rework the table

1. Goal: apply the chosen ideas.
2. Split into 15.16 – 15.18; tasks moved there.

### [x] 15.3 Attack animation

1. Goal: attacks look like Hearthstone (V-12).
2. Tasks:
   1. The attacker's card jumps at the monster, hits, and returns.
   2. Damage shows on the card that was hit.
   3. Only `transform` and `opacity`; follows animation speed (ST-05).
   4. A starburst with "-N" on the card that was hit (H-14).
   5. HP on the card drops as soon as the damage resolves (C-03).
   6. The lunge goes into the target and back, as in T-23.
3. Done when:
   1. e2e screenshot of an attack in progress.
4. Result (2026-10-13)
   1. A combat `damaged` event with a source plays a lunge.
   2. A copy of the source card flies; the card waits hidden at home.
   3. The copy winds up, goes into the target, and returns.
   4. At the hit, a yellow starburst "−N" pops on the target.
   5. The effect ends at the hit: the target's HP drops at once.
   6. A monster's hit on a player lunges at the character card.
   7. The copy plays in its own layer; the table never clips it.
   8. Web Animations, `transform` and `opacity` only; speed setting applies.
   9. Other damage has no lunge yet; 15.37 adds its marks.
   10. The "−N" popup when damage goes on the stack stays.
   11. Unit test: combat damage maps to lunges, other damage does not.
   12. e2e test: the copy flies, the card hides, the burst, HP drops.
   13. Screenshots: `.cache/screenshots/15.3/before/` and `after/`.
5. Owner notes (2026-10-13)
   1. Owner: a rising "−1" shows after a hit roll, before the lunge.
   2. Owner: drop it for attacks; the engine marks the damage.
6. Result of the fixes (2026-10-13)
   1. The engine's `damage_pending` says "combat" for attack damage.
   2. Same text as on `damaged`; no trigger reads it.
   3. The client shows no popup for it; the lunge's burst is the one "−N".
   4. Other pending damage keeps its popup until 15.37.
   5. Engine test: pending attack damage says "combat".
   6. Unit test: no popup for pending attack damage.
   7. Screenshots: `.cache/screenshots/15.3/fix/`.
7. Owner notes (2026-10-13)
   1. Owner marks 15.3 solved.

### [~] 15.4 Hearthstone-style card play

1. Goal: cards are played by dragging them to a place or target (V-16).
2. Split into 15.25 – 15.27; tasks moved there.

### [~] 15.5 Dice throw

1. Goal: dice look real and the player throws them (V-17).
2. Tasks:
   1. **Owner** decides what happens if nobody presses "Throw".
   2. "Throw" button for the player who rolls; others watch.
   3. Fall, bounce, and roll animation ends on the engine's value.
   4. Only `transform` and `opacity`; follows animation speed (ST-05).
   5. The roller can never miss that the game waits for them (V-17).
   6. Use the player's dice look (V-19).
   7. **Owner** decides the open V-17 questions: timer, bots, throw intent.
   8. The last roll stays visible as a small die by the stack (F-21).
3. Done when:
   1. e2e screenshot of a die in the air.
   2. The rolled value matches the engine in every test.
4. Owner decisions (2026-10-15)
   1. A new throw intent from the roller.
   2. The engine draws the value only when the intent arrives.
   3. No timer by default; the response timer throws for an idle player.
   4. Kicked seats throw at once; bots will too (step 11).
   5. A disconnected player pauses the game, as now.
5. Split into 15.59 – 15.63; tasks moved there.

### [x] 15.6 Each player's own dice

1. Goal: everyone sees the roller's dice, not their own (V-19).
2. Tasks:
   1. Hello sends the player's dice look; the seat keeps it.
   2. Lobby and game views list each seat's look.
   3. A roll is drawn with the roller's look on every screen.
   4. Unknown or missing texture: a plain die.
   5. Bots use the default look.
   6. Update protocol docs and golden files.
3. Done when:
   1. Server test: each seat's look reaches the other players.
   2. e2e test: two players see the same look for one roll.
4. Result (2026-10-15)
   1. `hello` and `Seat` gain `dice`; added fields keep protocol 3 (ADR 006).
   2. `CleanDice`: lowercase letters, digits, `_`, `-`; else the default.
   3. The table seat shows its client's look; the game seat keeps it.
   4. Away, the seat keeps its look; a return takes the new hello's.
   5. `diceImage` (`dice.ts`) finds the roller's texture in the seats.
   6. Throw, landing die, stack roll icons and last roll use it.
   7. No look (a bot): the default; an unknown look: a plain die.
   8. Golden files: hello and seats carry a look.
   9. Tests: `TestSeatsShowDice`, `TestCleanDice`, `diceImage` unit tests.
   10. e2e: the other player picks gold; the roller's tear shows on both.
   11. `make e2e` hit 19.21 (Eden in the soul test); fixed first.
   12. UC-12, ADR 006, Client and Server wiki, CHANGELOG updated.
   13. Screenshots: `.cache/screenshots/15.6/before/` and `after/`.

### [~] 15.7 Never leave the player lost

1. Goal: the player always knows when the game waits for them (V-20).
2. Split into 15.28 – 15.29; tasks moved there.

### [~] 15.8 Show what the active player is doing

1. Goal: every player knows what is happening right now (V-21).
2. Split into 15.30 – 15.31; tasks moved there.

### [x] 15.9 Mats and dice in WebP

1. Goal: mat and dice images as small as the cards (V-25).
2. Tasks:
   1. Work in `FourSoulsRepo/tools`, never in this repo.
   2. Convert `assets/mats/` and `assets/dices/` to WebP.
   3. Pick a size that still looks sharp full-screen.
   4. Amend the single commit of the images repo, as for cards.
   5. Update the code paths here that load JPG or PNG.
3. Done when:
   1. No JPG or PNG left under `assets/mats/` and `assets/dices/`.
   2. Mats and dice still show in `make internal`.
   3. A build without the images still works.
4. Result (2026-10-12)
   1. `scraper assets` in the tools repo: mats 3840 px, dice 1200 px, quality 85.
   2. 12 files, 1.6 MB to 0.4 MB; the images repo is still one commit.
   3. The app serves only WebP under `/cardimg/` now.
   4. Originals kept locally in `.cache/assets-originals-2026-10-12/`.
   5. Not converted: `coin.png` and two small Steam icons; owner decides.
   6. Owner: keep those three as they are for now; owner reviewed and pushed.

### [x] 15.10 Mat backgrounds

1. Goal: a game mat image under the table (V-24, ST-03).
2. Tasks:
   1. First mat: `basement` from `assets/mats/`.
   2. The app lists every image in `assets/mats/` as a mat choice.
   3. A new mat needs only a new file; no code change.
   4. Build without images: plain colors, as now.
   5. A saved mat that is missing falls back to colors.
   6. Cards and text stay readable on every mat.
   7. Update `docs/wiki/Client.md` and `docs/wiki/Build-modes.md`.
3. Done when:
   1. Go test: the mat list matches the folder.
   2. e2e test: the setting changes the background.
   3. A public-clone build shows the color table.
4. Result (2026-10-12)
   1. `App.Mats` lists every WebP in `assets/mats/`.
   2. The menu keeps the three color mats; an image takes its namesake's place.
   3. A new image name is added to the menu after them.
   4. The image lies under a dark shade, over its colors.
   5. The mat setting keeps any file name; a missing image shows colors.
   6. Default: the first mat image; Green felt without images.
   7. Screenshots: `.cache/screenshots/15.10/before/` and `after/`.
5. Owner notes (2026-10-12)
   1. Owner: the mat image is the default when it exists.
   2. Owner marks 15.10 resolved.
   3. Owner reviewed the default-mat change; 15.10 solved.

### [x] 15.11 Card backs

1. Goal: decks and opponent hands show real card backs (V-26).
2. Tasks:
   1. **Owner** confirms the source of back images.
   2. Tools repo collects loot, treasure, and monster backs.
   3. Decks show their own back; opponent hands show loot backs.
   4. No images: plain colored backs, as now.
   5. A back never reveals the card under it.
   6. Owner answer to 1: backs come from foursouls.com (tools repo).
   7. Amend the single commit of the images repo.
3. Done when:
   1. `make screenshots` shows backs on decks and opponent hands.
   2. A public-clone build shows the plain backs.
4. Result (2026-10-12)
   1. `scraper backs` in the tools repo: three site backs, 600 px, quality 85.
   2. Files: `assets/backs/loot.webp`, `treasure.webp`, `monster.webp`; 73 KB.
   3. `App.Backs` lists them; a missing file gives a plain back.
   4. An opponent's hand is one loot back with its count (F-08).
   5. Drawn cards fly to that back; it never carries a card ref.
   6. Screenshots: `.cache/screenshots/15.11/before/` and `after/`.
5. Owner notes (2026-10-12)
   1. Owner merged the backs into the images repo.
   2. Owner marks 15.11 solved.

### [~] 15.12 Turn action buttons

1. Goal: turn actions as clear icons next to own name (F-11).
2. Split into 15.19 – 15.20; tasks moved there.

### [~] 15.13 Stack panel and history strip

1. Goal: stack and history read at a glance, with icons (S-01 – S-11).
2. Split into 15.21 – 15.22; tasks moved there.

### [~] 15.14 Counters and card states

1. Goal: a card's state reads on the card itself (C-01 – C-06).
2. Split into 15.23 – 15.24; tasks moved there.

### [x] 15.15 Study Hearthstone targeting

1. Goal: know exactly how aiming at a target works in Hearthstone (H-23).
2. Tasks:
   1. Find gameplay videos and screenshots of aiming; keep them in `.cache/`.
   2. How the arrow starts, follows the pointer, and snaps to a target.
   3. How valid targets glow and invalid ones look.
   4. How aiming is cancelled: right click, drop on nothing, Esc.
   5. Differences: spell, minion attack, hero power, battlecry.
   6. How it works on touch (mobile Hearthstone).
   7. Write findings below as T-01, T-02, …; ideas only.
3. Done when:
   1. **Owner** picks the ideas for 15.25.
4. How the study went (2026-10-12)
   1. Main source: the Hearthstone wiki pages Target, Attack, Hero Power, Gameplay.
   2. Video: YouTube stops headless playback; storyboard frames are too small.
   3. Then the owner's Short, `YhKh9Eybil0`, via yt-dlp: 1080p, 2 frames a second.
   4. Frames: `.cache/screenshots/reference/hearthstone/targeting/`.
   5. Touch is still unconfirmed; marked "memory".
5. Findings
   1. T-01. Only a playable card or ability can start aiming.
   2. T-02. A bold red arrow runs from the source to the pointer.
   3. T-03. Over a valid target a red crosshair marks it.
   4. T-04. Click a valid target: the action completes, the arrow goes.
   5. T-05. Click an invalid target: a short reason, e.g. "A minion with Taunt is in the way".
   6. T-06. Cancel: click empty space or right-click.
   7. T-07. A targeted ability is cancelled by clicking its source again.
   8. T-08. The opponent sees the arrow live, while it moves.
   9. T-09. The opponent sees which hand slot aims, not which card.
   10. T-10. Hover a card: it is marked red on the opponent's screen.
   11. T-11. Ready attackers have a bright green frame; "zzz" when not ready.
   12. T-12. No target needed: drag the card onto the board; no arrow.
   13. T-13. Place first, then aim: cancelling the aim undoes the whole play.
   14. T-14. No valid target: the card plays without its targeted effect.
   15. T-15. After the target locks, the effect flies from source to target.
   16. T-16. The played card shows large at the left meanwhile (H-10).
   17. T-17. Touch: drag from the source, release on the target (memory).
   18. T-18. The arrow: a chain of thick red blocks with a big head.
   19. T-19. The arrow is a little see-through and draws over every card.
   20. T-20. Off target: a plain head; on a target: a red ring around the head.
   21. T-21. On the defender's screen, the aimed-at card glows red.
   22. T-22. An opponent's aimed card shows as a card back lifted from the hand.
   23. T-23. The attacker lunges into the defender, then returns (V-12).
   24. T-24. Damage: a yellow starburst "-N" on every damaged card at once.
   25. T-25. An HP number turns red while it is below full.
   26. T-26. A dying card gets a big skull over it before it leaves.
   27. T-27. Hover a board card: its zoom pops up right beside it (F-16).
6. What it means for us
   1. Our sources: hand cards, items, the character (attack).
   2. Our targets: monsters, players, items, stack items, the monster deck.
   3. T-08 and T-09 need a new server message for live aiming.
   4. **Owner** decides: share live aiming with others, or keep it local.
   5. T-09 must never leak the card (A-16).
   6. T-14 depends on the engine; it already decides valid targets.
7. Owner notes (2026-10-12)
   1. Owner picks all T ideas; planned in 15.32 – 15.37.
   2. Live aiming: a new kind of message, never recorded (N-15).

### [x] 15.16 Board frame

1. Goal: one framed board in the middle that never changes size.
2. Tasks:
   1. One framed board: decks, shop, monsters, souls, stack (F-01, H-01).
   2. Nothing floats over the board; fixes O-01 and O-02 (F-02).
   3. Drop the green felt and seat boxes (H-05, O-03).
   4. One fixed shape and size, empty or full (H-06, V-23).
   5. Decks show their backs at the board's edge (H-16).
   6. Keep V-01 – V-08 and touch (B-06); update `docs/wiki/Client.md`.
3. Done when:
   1. e2e test: board size stays the same when a shop slot empties.
   2. e2e test: no panel covers a deck with 2, 3, or 4 players.
4. Result (2026-10-12)
   1. One framed board: info column, slots, bonus souls, decks.
   2. Info column: turn badge, stack, history; nothing floats over the board.
   3. Decks sit at the right edge, each with its discard (H-16).
   4. Fixed shape; it scales with the window, never with its content.
   5. Up to four slots show full size; more share the row.
   6. No seat boxes; the active player's name is in the accent color.
   7. Green felt is no longer the default; Basement is, without images.
   8. Choices open at the bottom edge, so the board stays in view.
   9. Screenshots: `.cache/screenshots/15.16/before/` and `after/`.
   10. e2e tests: 3, 4, and 6 monster and shop slots on the same board.
   11. Screenshots now show real card and dice images, before and after.
5. Owner notes (2026-10-12)
   1. Owner agrees: Basement, not Green felt, is the default without images.
   2. Owner agrees: choices open at the bottom edge.
   3. Owner agrees: the info column keeps its width with the history hidden.
   4. Owner marks 15.16 solved.

### [x] 15.17 Own area

1. Goal: own character, items, hand, and coins read at a glance.
2. Tasks:
   1. Own character centered at the bottom edge, items beside it (F-06, H-02).
   2. Own hand large and fanned at the bottom, partly off screen (F-05, H-03).
   3. Coins big next to own name, like mana (H-18).
   4. HP and ATK as big numbers on the character card (H-15).
   5. The active player's name is highlighted (F-10).
   6. Key facts always visible; details on hover (V-18).
3. Done when:
   1. `make screenshots` shows the new own area.
   2. `make frontend e2e` passes.
4. Result (2026-10-12)
   1. Status at the left: name, big coins with a coin mark, souls.
   2. The active player's name is in the accent color; "what is left" below.
   3. Character centered under the board; items to its right.
   4. Live HP and ATK as big round badges on the character card.
   5. Hand fanned at the bottom edge; each card hangs about 40 px off screen.
   6. The fan narrows and overlaps more as the hand grows.
   7. A hovered hand card comes to the front; rising is 15.27.
   8. e2e test: character centered, badges and coins shown, hand cut off.
   9. The board is a little smaller: the own area takes about 35 px more.
   10. Screenshots: `.cache/screenshots/15.17/before/` and `after/`.
5. Owner notes (2026-10-12)
   1. Owner marks 15.17 solved.
   2. New sub-step: the real coin image as the coin icon; planned as 15.38.

### [x] 15.18 Opponents

1. Goal: opponents placed well for 2, 3, and 4 players.
2. Tasks:
   1. Opponent: name and coins, then character and items (F-07).
   2. Opponent hand: fanned backs with a count (F-08, H-17).
   3. 3–4 players: opponents in side columns; the board grows (F-09).
   4. Player names small, outside the board (H-21).
   5. Opponent placement still follows V-02.
3. Done when:
   1. `make screenshots` shows every player count.
   2. `make frontend e2e` passes.
4. Result (2026-10-12)
   1. No text header; a small name, coins with the coin icon, and souls.
   2. Live HP and ATK as badges on the character, as for the own area.
   3. Hand: one loot back per card, at most eight, fanned, with the count.
   4. Top: the own area mirrored; the character centered over the board.
   5. Top hand hangs past the top edge, backs upside down.
   6. Sides: 150 px instead of 200 px; name first, clear of the menu.
   7. Side cards are small, two to a row.
   8. 3 players: the board is about 12% wider than before.
   9. 4 players: about 2% wider; 2 players: about 1% smaller.
   10. The allowance wraps only between its parts.
   11. e2e test: coins, badges, one back per card, top hand off screen.
   12. Screenshots: `.cache/screenshots/15.18/before/` and `after/`.
5. Owner notes (2026-10-12)
   1. Owner marks 15.18 solved.

### [x] 15.19 Action icons

1. Goal: buy and attack as clear icons next to own name (F-11).
2. Tasks:
   1. Buy and attack as icon buttons by own name.
   2. Lit when usable, dim when not (F-12).
   3. Buy turns into Cancel while choosing what to buy (F-13).
   4. Attack turns into Roll once 15.5 adds the throw (F-13).
   5. A tooltip names each button and its key.
   6. Instant feedback on press (V-15); no number badges (V-29).
3. Done when:
   1. e2e test: each icon is lit or dim as the engine allows.
4. Result (2026-10-12)
   1. Buy (a bag) and Attack (a sword): round icons beside own status.
   2. A column right of the coins, so the own area keeps its height.
   3. Lit gold when the engine allows; dim grey otherwise.
   4. The action bar keeps only Pass, End turn, and Skip all.
   5. Keys: B buys, A attacks; the tooltip says "Buy (B)".
   6. A press shows at once, before the server answers (V-15).
   7. While choosing what to buy or attack, that icon stays pressed.
   8. No Cancel: R-SHOP-02 passes priority once a purchase is declared.
   9. Roll on Attack waits for 15.5.
   10. e2e test: icons follow the allowed intents; B opens the buy choice.
   11. Screenshots: `.cache/screenshots/15.19/before/` and `after/`.
5. Owner notes (2026-10-12)
   1. Owner marks 15.19 solved.

### [x] 15.20 End turn button and keys

1. Goal: one big End turn button, as in Hearthstone (H-04).
2. Tasks:
   1. One big button at the board's right edge (H-04).
   2. Yellow on own turn, grey "Their turn" otherwise (H-19).
   3. Glows green when nothing is left to do (H-20).
   4. Keys: Space resolves or passes; E ends the turn (F-15).
3. Done when:
   1. e2e test: Space and E send the same intents as the buttons.
   2. `make screenshots` shows own turn and their turn.
4. Result (2026-10-12)
   1. End turn left the action bar for the board's new last column.
   2. The column is part of the fixed board shape (V-23).
   3. Cards keep their size; the board only grows wider.
   4. Own turn: yellow "End turn"; dim while not allowed yet.
   5. Other turns: grey "Their turn".
   6. Green pulsing glow when only Pass or End turn is left (H-20).
   7. The glow holds still with reduced motion.
   8. E ends the turn; Space passes, like the Pass button.
   9. Keys in a text field, dialog, or menu are left alone.
   10. Space never clicks a focused button, so it cannot end the turn.
   11. Tooltips: "End turn (E)", "Pass (Space)".
   12. B and A now share the key helper (`keys.ts`).
   13. e2e test: both labels; Space sends Pass, E sends End turn.
   14. The green glow has a unit test; no screenshot of it yet.
   15. Screenshots: `.cache/screenshots/15.20/before/` and `after/`.
5. Owner notes (2026-10-12)
   1. Owner: the button looks generic; style it for the game.
   2. Owner: use the game's fonts; the system font without the images.
   3. Fonts: font-souls (MIT, GitHub) and the foursouls.ru card font.
   4. Owner amended and pushed the images repo with the fonts.
   5. Owner marks 15.20 solved.
6. Result of the restyle (2026-10-12)
   1. `scraper fonts` in the tools repo downloads both fonts.
   2. font-souls is pinned to one commit; its license is kept.
   3. foursouls.com fonts sit behind Cloudflare; not used.
   4. Files: `assets/fonts/` in the images repo, staged for the owner.
   5. `App.Fonts` gives their URLs, or none without the images.
   6. Stack: font-souls for Latin, the RU font for Cyrillic, then system.
   7. End turn: cut-out look with ink outline, tilt, and hard shadow.
   8. Pass, Skip all, Buy, and Attack share the ink outline and shadow.
   9. End turn is bigger: label about 21 px at 2 players.
   10. e2e test: both font faces are registered.
   11. Screenshots: `.cache/screenshots/15.20/styled/`.

### [x] 15.21 Stack panel

1. Goal: the stack reads at a glance, with icons (S-01 – S-09).
2. Tasks:
   1. Stack panel inside the board, at its left side (F-03).
   2. Newest on top; each row is an icon and one short line (S-01, S-02).
   3. Names in color: player in their color, target in red (S-03).
   4. Icon by kind: card art, die for a roll, mark for damage (S-04).
   5. Hover a row: card large, used ability highlighted, owner below (S-06, S-07).
   6. Resolve or Pass in the panel; "Nothing on the stack" when empty (S-08, S-09).
   7. Keep the "waiting for" lines (V-11); the reference has none.
3. Done when:
   1. e2e test: hover a stack row shows its card and owner.
   2. `make screenshots` shows a stack of 3 items.
4. Result (2026-10-12)
   1. The panel always shows, in the board's info column (F-03).
   2. Empty: "Nothing on the stack" and a grey Pass (S-09).
   3. Rows: newest on top; a square icon and one line (S-01, S-02).
   4. Icons: card art, a die with the roll, a blood drop, a skull (S-04).
   5. Each seat has its own color; a dot by the name shows it (S-03).
   6. Players on the stack in their color; targets in red.
   7. Hover or tap a row: its card large, a caption below (S-06).
   8. The caption: the used ability's words, highlighted, then the owner.
   9. Highlight is a caption, not on the image; images have no ability map.
   10. The engine view adds `text` to stack items: ability or mode words.
   11. Pass left the action bar for the panel header; Space still passes.
   12. The label stays "Pass": passing resolves only once all pass.
   13. "Waiting for" and "Passed" lines stay (V-11).
   14. e2e test: a 3-item stack, hover shows card, ability, and owner.
   15. Screenshots: `.cache/screenshots/15.21/before/` and `after/`.
5. Owner notes (2026-10-12)
   1. Owner likes the caption under the card and the "Pass" label.
   2. Owner: center the card art in the row icon.
   3. Owner: a roll shows the player's die, not only the number.
   4. Owner: compare history with Hearthstone; theirs is one thin line.
   5. Owner marks 15.21 solved.
6. Result of the fixes (2026-10-12)
   1. The card art is centered in its square.
   2. A roll shows its face in the viewer's dice look (V-19).
   3. Each player's own dice come with 15.6.
   4. Hearthstone: no stack; a ~45 px icon column at the board's left edge.
   5. Ours: the info column is ~25% of the board, history included.
   6. 15.22 brings the thin strip; see the proposal there.
   7. Screenshots: `.cache/screenshots/15.21/fix/`.

### [x] 15.22 History strip

1. Goal: a thin icon history along the board's left edge (H-07).
2. Tasks:
   1. Resolved items as the stack's icons, newest on top (S-10, F-04).
   2. The strip runs along the board's left edge (H-07).
   3. Icon badge: sword for an attack, bolt for a spell (H-08).
   4. Hover an icon: the same large card view (S-11, H-09).
   5. The full history panel (V-03) stays one click away.
3. Done when:
   1. e2e test: hover a history icon shows its card.
   2. `make screenshots` shows the history strip.
4. Proposal after the Hearthstone check (2026-10-12)
   1. The strip replaces the History panel in the info column.
   2. The strip is ~0.45 u wide, at the board's left edge (H-07).
   3. The info column keeps only the turn badge and the stack.
   4. **Owner** decides: narrow the info column, or let the stack use the room.
   5. Owner picks B: the info column keeps its width.
   6. The strip sits inside the board, left of the info column.
   7. The stack takes the room History used: ~7 rows before scrolling.
5. Result (2026-10-12)
   1. A new first board column, 0.45 u wide, holds the strip (H-07).
   2. The board grows 12.1 u → 12.75 u wide; cards shrink ~5% to fit.
   3. Icons: the stack's own, from one shared component (S-10).
   4. Newest on top; a gold line under the first item of each turn.
   5. Badges: a red sword for an attack, a blue bolt for a loot play (H-08).
   6. Other items (abilities, triggers) have no badge.
   7. Hover or tap: the stack's large card, ability, and owner (S-11, H-09).
   8. A scroll button tops the strip; it opens or hides the full History.
   9. The full History starts hidden; open, it sits under the stack.
   10. The engine's `stack_resolved` event now carries the stack item.
   11. The item has its ability's words, like the view's stack items.
   12. Only stack items show in the strip; loots and turns stay in History.
   13. e2e test: a loot played in a real game shows with a bolt badge.
   14. e2e test: 3 items, badges, order, and hover with card and owner.
   15. Screenshots: `.cache/screenshots/15.22/before/` and `after/`.
6. Owner notes (2026-10-12)
   1. Owner: attack rolls and attack damage do not say whom they hit.
   2. Owner: Blank Rune shows the whole card text, not the rolled line.
   3. Owner: Blank Rune shows twice; the first did nothing yet.
   4. Owner: show only the effect that was applied; split the items.
7. Result of the fixes (2026-10-12)
   1. Stack items name their source and target cards (`source_card`, `target_card`).
   2. Views and the resolved event fill them, while the card is public.
   3. A dead monster stays named; its object is gone by then.
   4. Attack roll: "Attack: Ann rolls 4 against Fly"; hover shows Fly.
   5. Damage names its dealer: "Ann deals 2 damage to Fly".
   6. A monster hitting back: "Fly deals 1 damage to Ann".
   7. An item without ability words shows this line under the card.
   8. A card that rolls resolves twice: the play, then its result.
   9. The play now says only "Roll." (the words up to the roll).
   10. The result says only its line: "1: Each player gains 1¢."
   11. A result without a line says "6: Nothing happens."
   12. The stack says "Roll." too (owner, see below).
   13. A card's roll names the card: "Ann rolls 1 for Blank Rune".
   14. The full History panel keeps its old lines for now.
   15. Screenshots: `.cache/screenshots/15.22/fix/`.
8. Owner notes (2026-10-12)
   1. Owner: show "Roll." on the stack too; the card has the full text.
   2. Owner asks: does a card with several abilities show only the used one?
9. Result (2026-10-12)
   1. The stack says "Roll." for an item that rolls.
   2. Without images the zoomed text card lists every result too.
   3. Several abilities (The Bone): the item knows its ability's index.
   4. Stack and history show only that ability's words.
   5. "Choose one-" cards (Dagaz) show only the chosen mode.
   6. Engine tests check The Bone and Blank Rune on the stack and in history.
10. Owner notes (2026-10-12)
   1. Owner marks 15.22 solved.

### [x] 15.23 Counter badges and live stats

1. Goal: counters and current stats read on the card (C-01 – C-03).
2. Tasks:
   1. Counter: a big round badge with a number on the card art (C-01).
   2. The zoomed card shows the same badge (C-02).
   3. Live HP, ATK, and dice values on the card (C-03).
   4. **Owner** decides: values drawn over the image, or a badge beside it.
   5. Text cards without images get the same badges.
3. Done when:
   1. e2e test: a counter badge on the card and in its zoom.
   2. e2e test: a damaged monster shows its new HP on the card.
4. Owner decision (2026-10-13)
   1. Monster stats as badges beside the print, only when changed.
5. Result (2026-10-13)
   1. Counter: a big green circle with the count, on the card art (C-01).
   2. A text card has it in the top corner; the name keeps clear of it.
   3. The engine view gives monsters `stats`: live and printed HP, dice, ATK.
   4. A badge shows only a value that differs from the print (C-03).
   5. Red health, a white die for the dice, gold attack.
   6. A green ring marks a value above the print.
   7. A damaged monster's health badge replaces its red damage mark.
   8. Characters keep their always-on health and attack badges (15.17).
   9. The zoom shows the same badges as the card (C-02).
   10. Text cards get the same badges.
   11. e2e test: counters, health, dice and attack on the card and in the zoom.
   12. Engine test: a damaged, boosted monster's view stats.
   13. Fix: screenshots had text cards; the fake bridge read a wrong field.
   14. Screenshots: `.cache/screenshots/15.23/before/` and `after/`.
6. Owner notes (2026-10-13)
   1. Owner: use the game's fonts for the badges.
   2. Owner: no badges on the zoomed card; they hide its text.
   3. Owner: style the badges to suit the game.
7. Result of the fixes (2026-10-13)
   1. Numbers in the game's hand-drawn font, with an ink edge.
   2. Shapes in the End turn button's style: ink outline, hard shadow, tilt.
   3. Health is a red heart, the dice a bone-white die, attack a gold blob.
   4. Counters are a green blob; damage on a character a dark-red one.
   5. Characters' health and attack badges get the same look.
   6. The zoomed card shows no badges; this replaces C-02.
   7. The green "higher" ring is gone; a badge's tooltip names the print.
   8. e2e test: badges on the card, none in the zoom.
   9. Screenshots: `.cache/screenshots/15.23/fix/`.
8. Owner notes (2026-10-13)
   1. Owner marks 15.23 solved.

### [x] 15.24 Ability dialog and status marks

1. Goal: clear choice of ability and clear card states (C-04, C-05, H-22).
2. Tasks:
   1. Two abilities: a dialog with the big card; pick one line (C-05).
   2. Status marks on cards, e.g. not usable yet (H-22).
   3. Used items stay dimmed, not rotated; already so (C-04).
   4. A character that can attack has a green frame; not yet: "zzz" (T-11).
3. Done when:
   1. e2e test: a two-ability card opens the dialog and uses the picked line.
4. Result (2026-10-13)
   1. The engine view gives played cards `abilities`: each ability's words.
   2. Only on cards with two or more activated abilities.
   3. Indexed like the activate intent's choice, copied abilities too.
   4. Such a card opens a dialog: the card large, one button per line.
   5. Lines not usable now stay greyed, with "Not usable now".
   6. Cancel or a click outside closes it; one ability keeps the menu.
   7. Own character, own turn: a bright green frame while it may attack.
   8. "zzz" in the game's font once no attack is left this turn.
   9. Nothing while an attack waits (stack) or on others' turns.
   10. Used items stay dimmed; no change (C-04).
   11. The character box no longer clips the glow.
   12. Engine test: The Bone lists its four lines; Cain has none.
   13. e2e test: The Bone's dialog, a greyed line, the picked line is sent.
   14. e2e test: green frame with an attack, "zzz" without.
   15. Screenshots: `.cache/screenshots/15.24/before/` and `after/`.
5. Owner notes (2026-10-13)
   1. Owner: use the Four Souls font for dialogs.
   2. Owner: style each ability line to suit the game.
6. Result of the fixes (2026-10-13)
   1. The app also loads font-souls Body, the card text font.
   2. Dialog titles in the game's title font; options in the body font.
   3. Each line: parchment, ink outline, hard shadow, a slight tilt.
   4. The cost icon boxed at the left: a drawn tap arrow, or "$".
   5. The cost in red ink, then the effect.
   6. Fix: "zzz" showed before the action step; now only in it.
   7. Screenshots: `.cache/screenshots/15.24/fix/`.
7. Owner notes (2026-10-13)
   1. Owner: remove "zzz"; it does not suit this game.
   2. Owner: no green "ready to attack" frame either.
   3. Owner: rework the Loot, Attacks, Buys line; badges on the icons.
8. Result of the fixes (2026-10-13)
   1. "zzz" and the green frame are gone (T-11, H-22 dropped).
   2. Buy and Attack icons carry a paper badge: what is left this turn.
   3. A badge turns dark at 0; badges show only on own turn.
   4. Own name shows "Loot N"; opponents keep the full line.
   5. e2e test: badges on own icons, the full line for others.
   6. Screenshots: `.cache/screenshots/15.24/fix2/`.
9. Owner notes (2026-10-13)
   1. Owner: the badge digit touches the circle's top.
   2. Owner: check 2, 3, 99 and 100.
10. Result of the fixes (2026-10-13)
    1. A top pad centers the digit; the game font sits high.
    2. One digit: a round disc; two or three: a pill with smaller digits.
    3. Past 99 the badge says "99+".
    4. The pill grows left: the character card is close on the right.
    5. e2e test: 100 reads "99+" and stays clear of the card.
    6. Screenshots: `.cache/screenshots/15.24/badges/`.
11. Owner notes (2026-10-13)
    1. Owner: use the sword from the images' `assets` folder for Attack.
    2. Owner: convert it to WebP and move it to the icons folder first.
12. Result (2026-10-13)
    1. `assets/sword.png` to `assets/icons/sword.webp`: lossless, 65 px, 1.5 KB.
    2. Pixel-identical to the PNG; `sword.png` is kept, as `coin.png` was.
    3. The WebP is staged in the images repo; the owner amends and pushes.
    4. `App.Icons` gives its URL; the Attack icon shows it, pixels crisp.
    5. Dim, it is greyed and faded like the drawn icons.
    6. A missing or broken image shows the drawn sword.
    7. e2e test: one page gets the image, the other the drawn sword.
    8. Screenshots: `.cache/screenshots/15.24/sword/`.
13. Owner notes (2026-10-13)
    1. Owner marks 15.24 solved.

### [x] 15.25 Drag and aim

1. Goal: play a card by dragging it to the table or a target (V-16).
2. Tasks:
   1. Apply the targeting ideas the owner picked in 15.15.
   2. Drag a card onto the table or onto its target.
   3. Aiming arrow; valid targets glow, others stay plain.
   4. Works with mouse and touch (B-06).
   5. Touch: drag from the source, release on the target (T-17).
3. Done when:
   1. e2e test: a bomb dragged onto a monster damages it.
4. Result (2026-10-13)
   1. Owner: the engine sends the targets, so they glow while dragging.
   2. New `update.aims`: first target options of each allowed play or use.
   3. Keyed by the intent's index in `allowed`; built by `Game.Aims`.
   4. No aim for a card with a mode, no target, or a copied ability.
   5. A card that aims stays lifted in the fan while dragged.
   6. A plain red arrow runs from the card to the pointer; 15.32 styles it.
   7. Valid targets glow red: cards, player areas, stack rows.
   8. A drop on a target sends the play; its target prompt picks it at once.
   9. A drop elsewhere on the board plays it; the dialog asks, as before.
   10. Cards that need no target move with the pointer, as before.
   11. Hand cards only; items and the character aim in 15.33.
   12. Touch shares the pointer code; e2e drives the mouse only.
   13. Engine test: every aim matches the target prompt its intent opens.
   14. Unit tests: aims by index, drop target lookup, prompt auto-pick.
   15. e2e test: a save gives a Bomb!; dragged onto a monster, it hits it.
   16. Screenshots: `.cache/screenshots/15.25/before/` and `after/`.
5. Owner notes (2026-10-13)
   1. Owner: aim by clicks: one click starts aiming, a second picks the target.
   2. Owner: no press, drag, and release to aim.
   3. Owner: the arrow's geometry is bad.
6. Result of the fixes (2026-10-13)
   1. A click on a card whose one action needs a target starts aiming.
   2. The next click on a target plays it there; no dialog asks.
   3. A click elsewhere, a right click, Esc, or a new update cancels.
   4. Such a card opens no menu and is not dragged.
   5. Items with one targeted ability aim the same way.
   6. The arrow starts at the card's middle, measured as it moves.
   7. It used the card's place before the lift, so it started off the card.
   8. The shaft stops at the head's base; the head's tip is at the pointer.
   9. The head is a plain triangle; the shaft no longer shows through it.
   10. The arrow shows only for a mouse; touch taps the card, then the target.
   11. e2e helper `playLoot` aims when a card needs a target.
   12. e2e test: a right click cancels; the arrow's ends; two clicks hit.
   13. Screenshots: `.cache/screenshots/15.25/fix/`.
7. Owner notes (2026-10-13)
   1. Owner marks 15.25 solved.

### [-] 15.26 Click modes

1. Dropped by the owner, 2026-10-15.
2. Every play is one click or tap; a second one picks the target.
3. Drag play is not needed, so no mode setting.

### [x] 15.27 Hand hover and glow

1. Goal: the hand shows what can be played and what it does.
2. Tasks:
   1. Hover a hand card: it rises large, keyword tips beside it (H-13).
   2. Zoom pops up right beside the hovered card, not at a screen side (F-16, T-27).
   3. Playable cards glow green; the rest stay plain (H-12).
3. Done when:
   1. e2e test: playable cards glow, others do not.
4. Result (2026-10-15)
   1. A hand card the player may play or use glows green, not yellow.
   2. The rest of the hand stays plain; table cards keep yellow.
   3. The class stays `hl-act`, so e2e helpers find it as before.
   4. A hovered hand card's zoom rises from it.
   5. It stands on the window's bottom edge, centered over the card.
   6. Every other zoom opens right beside its card or row.
   7. It opens toward the window's middle and stays inside the window.
   8. Stack rows, history icons, and history entries zoom the same way.
   9. A touch zoom stays centered, as before.
   10. Keyword tips beside the zoom: the card's footnotes.
   11. Today: Eternal, Trinket, Curse, Guppy; name and meaning.
   12. They stand on the side away from the card.
   13. The blank zoom face no longer repeats them.
   14. `zoom.ts`: `placeZoom` is pure; the zoom measures itself first.
   15. Unit tests: beside, both sides, window edges, rise; keyword names.
   16. e2e test: playable cards glow green; the other player's do not.
   17. e2e test: the hand zoom rises over its card; a monster's is beside it.
   18. e2e test: The D6's zoom shows the Eternal tip; Isaac's has none.
   19. Screenshots: `.cache/screenshots/15.27/before/` and `after/`.
5. Owner notes (2026-10-15)
   1. Owner: the green glow looks cheap; make it look like a serious game.
6. Result of the fixes (2026-10-15)
   1. No hard 3 px ring any more.
   2. A thin light rim on the card's edge; light falls a little inward.
   3. Behind the card, a soft aura in layers: bright near, fading out.
   4. The wide aura breathes slowly; only its opacity moves.
   5. Reduced motion stands it still.
   6. New animation `hand-glow`, listed in the Animations wiki page.
   7. The card's fan slot draws the aura: class `usable`, beside `hl-act`.
   8. e2e test reads the aura from the slot, not a ring on the card.
   9. Screenshots: `.cache/screenshots/15.27/fix/`.
7. Owner notes (2026-10-15)
   1. Owner marks 15.27 solved.

### [x] 15.28 Prompt banners

1. Goal: the player always knows when the game waits for them (V-20).
2. Tasks:
   1. List every prompt the engine can send: priority, targets, choices, pay, roll.
   2. Each prompt shows a clear banner: what to do, in a few words.
   3. Cards and buttons that answer it glow; the rest dim.
   4. The stack panel shows the same prompt (V-11).
   5. Waiting on others names who and what for.
   6. Optional: a soft sound or window flash when input is needed.
3. Done when:
   1. e2e test: every prompt type shows its banner and highlights.
   2. `make screenshots` covers "your turn to act" and "waiting for Ann".
4. Result (2026-10-19)
   1. The status line, bottom right, became the prompt banner.
   2. Asked player: a gold plate, the game's letters, a breathing glow.
   3. Others: a dark plate, "Waiting for Ann", her name in her color.
   4. A hint under the title says more.
   5. Prompts covered: turn, respond, act out of turn.
   6. Also: choice, discard, put back in order, throw.
   7. The player's own aim counts too: "Choose a target", Esc cancels.
   8. Own turn lists what is left: "Play a card, attack, buy, or end the turn".
   9. Words live in `prompt.ts` (`promptBanner`); `waitText` is gone.
   10. The stack panel repeats prompts other than priority (V-11).
   11. Its words match the banner; "Waiting for you to throw" is gone.
   12. Pass glows when passing answers the prompt.
   13. Picking a target: answers glow, every other card dims.
   14. While aiming, other playable cards no longer glow.
   15. New animation `prompt-glow`, in the Animations wiki page.
   16. Optional sound or window flash: not done; left for the owner.
   17. Use case UC-21; UC-14 and the Client wiki updated.
   18. e2e `prompt.spec.ts`: every prompt, for both players.
   19. e2e: Pass glow, pick dimming, own aim.
   20. `make e2e`: only 19.18 failed (Dagaz in the deal), twice.
   21. Screenshots: `.cache/screenshots/15.28/before/` and `after/`.
5. Owner notes (2026-10-19)
   1. Owner marks 15.28 solved.

### [x] 15.29 Reasons and choice dialogs

1. Goal: the player sees why something is not possible, and chooses calmly.
2. Tasks:
   1. Hover an unusable card: the zoom says why (F-17, C-06).
   2. A refused action shows its reason next to the clicked item (F-18).
   3. Choices: big cards over a dimmed board (F-20, H-11).
   4. A "View board" button hides the choice to look at the table (F-20).
3. Done when:
   1. e2e test: an unusable card's zoom shows the reason.
   2. e2e test: "View board" hides and restores a choice.
4. Result (2026-10-19)
   1. Engine: `Game.Blocked(p)` says why each card or turn action is not allowed.
   2. It checks each intent as if p had priority; `RuleError.Say` holds the words.
   3. No refusal then: "Wait for your priority" or "Answer the prompt first".
   4. Server sends it as `update.blocked`, only to its seat; ADR 006 updated.
   5. Zoom: a red line under the card, e.g. "No loot plays left".
   6. A click on that card shows the same note where it was clicked.
   7. Buy, Attack and own End turn: dim, but a click says why.
   8. One shared note for all of these and the aim's wrong click (`WhyNote`).
   9. Choice, discard and order: big cards over the dimmed table.
   10. "View board" hides the window; a bar brings it back.
   11. Meanwhile table options glow and answer as before.
   12. Order dialog keeps its smaller cards, so long rows fit.
   13. Rooms get reasons from the engine; the client draws no rooms yet.
   14. Use case UC-22; UC-13, Client, engine, Layers, Animations wikis updated.
   15. e2e `reasons.spec.ts`: zoom, click note, icon note, View board.
   16. Screenshots: `.cache/screenshots/15.29/before/` and `after/`.
   17. The Menu button now draws above a choice: leave and save stay reachable.
   18. Board-layout tests ban Eden: its starting choice now covers the decks.
   19. The soul pick test presses "View board" before it picks on the table.
   20. 19.28: the Eden ban made the 4-player test check too early; fixed.
   21. `make e2e`: 194 passed; earlier failures under load are kept as evidence.
5. Owner notes (2026-10-19)
   1. Card options repeat the card's name in a button under it; see 15.66.
   2. Owner marks 15.29 solved.

### [x] 15.30 Activity status

1. Goal: every player knows what the active player is doing (V-21).
2. Tasks:
   1. List the activities to show: fight, shop, play a card, pay, think.
   2. Derive each one from game state on the client.
   3. The active seat shows a short status: "Ann attacks Pooter".
   4. "Thinking" shows while the active player holds priority.
   5. Spectators see the same marks.
   6. Share the banner and highlight code with 15.28.
3. Done when:
   1. e2e test: a second player sees the status line.
   2. `make screenshots` covers "Ann attacks" and "Ann is thinking".
4. Result (2026-10-19)
   1. Activities: play, use, death penalty, attack, buy, think.
   2. The active seat shows one line under the name: "Ann attacks Pooter".
   3. Names in seat color, targets in red, as on the stack.
   4. Engine: `GameView.attack` (target, deck) and `buying`, public.
   5. The rest comes from the stack top and the open prompt.
   6. "Thinking": any prompt that waits for the active player.
   7. No line while the game waits for others with nothing going on.
   8. The viewer's own seat shows none; the banner covers it.
   9. Spectators: public state only; no spectator screen yet (step 9).
   10. Shared with 15.28: `activity` in `prompt.ts`, the stack's `Part`.
   11. Engine tests: attack and purchase in every viewer's view.
   12. Vitest: every activity; e2e `prompt.spec.ts`: the other seat.
   13. Use case UC-23; wiki `Client`, `Rules-engine`.
   14. Screenshots: `.cache/screenshots/15.30/before/` and `after/`.
5. Owner notes (2026-10-19)
   1. Owner marks 15.30 solved.

### [x] 15.31 Fight marks and played cards

1. Goal: fights and played cards are easy to follow for everyone.
2. Tasks:
   1. In a fight, the attacker and the monster get a red frame (F-19).
   2. Attack rolls, targets, and damage are visible to all.
   3. A card another player plays shows large for a moment (H-10).
   4. It stays large while its effect resolves (T-16).
3. Done when:
   1. e2e test: a second player sees the fight mark.
4. Result (2026-10-19)
   1. Attacker's character and target: red frame on every screen.
   2. The frame comes with the target; a deck attack waits for its card.
   3. Rolls, targets and damage were already public; now tested.
   4. Owner picked the place for played cards: pop, then by the stack.
   5. Another player's loot card pops large in the middle, "Ann plays".
   6. It then stays under the stack panel until the item resolves.
   7. It is still there while its own resolution effects play.
   8. Only loot cards; the viewer's own plays show neither.
   9. New animation `played`; claims the card's `show` (ADR 012).
   10. Vitest: `fighters`, `playedCard`, the `played` effect.
   11. e2e `fight.spec.ts`: the frame and roll, the pop and the kept card.
   12. Use case UC-24; wiki `Client`, `Animations`.
   13. Screenshots: `.cache/screenshots/15.31/before/` and `after/`.
5. Owner notes (2026-10-19)
   1. Owner marks 15.31 solved.

### [x] 15.32 Aiming arrow

1. Goal: a clear Hearthstone-style arrow while aiming (T-02, T-18 – T-20).
2. Tasks:
   1. A chain of thick red blocks with a big head, source to pointer (T-02, T-18).
   2. A little see-through; drawn over every card (T-19).
   3. Off target: plain head; on a valid target: a red ring (T-03, T-20).
   4. Only `transform` and `opacity`; follows animation speed (ST-05).
   5. Without images the arrow looks the same.
3. Done when:
   1. `make screenshots` shows the arrow off and on a target.
4. Result (2026-10-13)
   1. The arrow is a chain of thick red blocks with a notched head.
   2. Blocks: 28 by 26 px, one every 36 px, with a light top edge.
   3. The chain is placed once at the card's middle; it never stretches.
   4. The blocks march toward the head; the ends fade, so no seam shows.
   5. A little see-through, outlined, drawn over every card and the zoom.
   6. Over a valid target, a red ring with cross ticks marks the tip.
   7. The ring fades and shrinks in; off target it fades out.
   8. Only `transform` and `opacity` move.
   9. March speed follows the animation setting; off stands it still.
   10. Reduced motion also stands it still.
   11. Pure CSS: the same without the card images.
   12. Unit test: `blocks` counts the blocks short of the head.
   13. e2e test: plain head off a target, ring on it; the chain's origin.
   14. Screenshots: `.cache/screenshots/15.32/before/` and `after/`.
5. Owner notes (2026-10-13)
   1. Owner marks 15.32 solved.

### [x] 15.33 Start, finish, and cancel aiming

1. Goal: aiming starts, ends, and cancels as in Hearthstone.
2. Tasks:
   1. Only a playable card or ability can start aiming (T-01).
   2. A valid target completes the action; the arrow goes (T-04).
   3. An invalid target shows the engine's short reason (T-05).
   4. Cancel: click empty space, right-click, or Esc (T-06).
   5. Clicking the source again also cancels (T-07).
3. Done when:
   1. e2e test: right-click cancels and sends nothing.
   2. e2e test: an invalid target shows the reason.
4. Result (2026-10-13)
   1. T-01, T-04, T-06 came with the 15.25 fixes; kept as they are.
   2. Each aim now carries a `reason` from the engine (`targetReason`).
   3. It says what the card may target: "Pick a monster or a player."
   4. A filtered target adds "Only the glowing ones fit."
   5. A click on a card, player, pile, or stack row that is no target keeps aiming.
   6. The reason shows in a small note above the click; it fades in 2.5 s.
   7. A click on the aiming card cancels (T-07).
   8. The card wins over the player area around it.
   9. Bomb! may target its own player; a click on it used to pick you.
   10. A source that is its own target still picks itself.
   11. Empty space, a right click, or Esc still cancel.
   12. Protocol: `aims` gains the optional `reason`; golden files updated.
   13. Engine test: every aim has a reason; every target kind has words.
   14. Unit tests: `aimClick` picks, cancels, or calls a click wrong.
   15. e2e test: right click and the source cancel; nothing is sent.
   16. e2e test: a shop item shows "Pick a monster or a player." above it.
   17. Screenshots: `.cache/screenshots/15.33/before/` and `after/`.
5. Owner notes (2026-10-14)
   1. Owner: aim at a player by the character card only, as at a monster.
   2. Owner asks when a source can be its own target.
6. Result of the fixes (2026-10-14)
   1. A player is aimed at by the character card; the area no longer counts.
   2. The character glows red while aiming; the area stays plain.
   3. A click on empty space in a player area now cancels.
   4. Own target: The D20 and Battery Bum may pick any item, themselves too.
   5. No character card has a targeted ability today.
   6. e2e test: the opponent's character glows, the area does not.
   7. Screenshots: `.cache/screenshots/15.33/fix/before/` and `after/`.
   8. e2e helper `playLoot` takes any glowing target, e.g. your own item.
   9. It missed XV. The Devil's target, an item you control.
7. Owner notes (2026-10-14)
   1. Owner: a click on the aiming card no longer cancels; it does nothing.
   2. Owner: only a click on empty space stops aiming.
   3. Owner: a card that may target itself picks itself on that click.
8. Result of the fixes (2026-10-14)
   1. `aimClick` gives "none" for the source; aiming goes on.
   2. A source that is its own target picks itself, as before.
   3. Empty space, a right click, or Esc still cancel.
   4. Unit test: the source does nothing, even over a player it may target.
   5. e2e test: the bomb clicked again stays aimed; empty space cancels.
   6. No screenshots: nothing on screen changed.
9. Owner notes (2026-10-14)
   1. Owner marks 15.33 solved.

### [x] 15.34 Plays without aiming

1. Goal: cards that need no target, or place then aim, feel natural.
2. Tasks:
   1. No target needed: drag onto the board; no arrow (T-12).
   2. Place first, then aim; cancelling the aim undoes the play (T-13).
   3. No valid target: follow what the engine allows (T-14).
3. Done when:
   1. e2e test: a no-target loot card dragged onto the board plays.
   2. e2e test: a cancelled aim leaves the card in hand.
4. Result (2026-10-14)
   1. T-12 came with 15.25: a card with no target is dragged; no arrow.
   2. T-13: a card with a mode is dragged; the mode dialog asks first.
   3. Then its target is aimed from the card; no target dialog.
   4. The engine names that card in `view.activating`, only for the chooser.
   5. The card stays lifted in the hand while it aims.
   6. Esc, a right click, or empty space sends the prompt's cancel.
   7. The engine then leaves the card in hand; nothing is paid.
   8. Items with a mode, as Chaos Card, aim the same way.
   9. A target not on the table keeps the dialog.
   10. T-14: the engine never offers a card with no valid target (R-ABIL-06).
   11. While aiming, player areas no longer show "Choose …" picks.
   12. Engine test: the chooser sees the card; others do not.
   13. Unit tests: `promptAim` and `cancelOption`.
   14. e2e test: A Penny! dragged onto the board plays; no arrow, no dialog.
   15. e2e test: Dagaz aims after its mode; a right click leaves it in hand.
   16. Screenshots: `.cache/screenshots/15.34/before/` and `after/`.
5. Owner notes (2026-10-14)
   1. Owner marks 15.34 solved.

### [x] 15.35 Presence messages

1. Goal: a light message kind for live aiming and hover (N-15).
2. Tasks:
   1. ADR: presence messages; update ADR 006 and ADR 007.
   2. New message kind: not an intent, not an event; the engine never sees it.
   3. The server relays it to the other players of that game only.
   4. Never recorded; replays never show it.
   5. A hidden card travels as its hand slot only, never its identity (A-16).
   6. The server limits the rate per player and drops the rest.
   7. Update protocol docs and golden files.
3. Done when:
   1. Server test: an aim message reaches the others, not the sender.
   2. Server test: the record has no presence messages.
   3. Server test: a hidden card's aim carries no card ID.
4. Result (2026-10-20)
   1. ADR 016: presence messages; ADR 006 and ADR 007 point to it.
   2. `protocol.Presence`: `from`, `at`, `point`; the server sets `seat`.
   3. `protocol.Spot`: one of object, hand slot, player, deck, stack item.
   4. `server/presence.go`: relayed to the other seated players only.
   5. Spots are checked against the spectator view of this step.
   6. A hidden object becomes no spot; a hand card travels as its slot.
   7. 20 a second per seat, a burst of as many; the rest are dropped.
   8. Never answered, never recorded; the engine never sees it.
   9. A dropped player's arrow goes: the others get an empty presence.
   10. Protocol stays version 3: an older client ignores the new type.
   11. Golden file: one sent and three relayed presence lines.
   12. The client does not read it yet (15.36); TS test lists it as not yet.
   13. Tests in `presence_test.go`: reach, hidden card, pace, drop, record.
   14. Wiki `Server`: a Presence section and a leak audit row.
   15. No screen changes, so no screenshots.
5. Owner notes (2026-10-20)
   1. Spectators should see presence too.
   2. 20 a second is fine.
   3. A dedicated server may turn presence off, to save resources.
6. Result of the fixes (2026-10-20)
   1. Relayed to every other client of the game, spectators too (N-15).
   2. No spectator joins yet; step 9 brings them.
   3. `Config.Presence`, flag `-presence`, JSON `presence`; on by default.
   4. Off: the room drops every presence and relays none.
   5. `welcome.presence` tells the client; it sends none when unset.
   6. ADR 016, wikis `Server` and `Dedicated-server-on-a-VPS` updated.
   7. Tests: `TestPresenceOff`; `TestParse` reads the file and the flag.
   8. Open for 15.36: should waiting players send hover, or only aims?
7. Owner notes (2026-10-20)
   1. Owner marks 15.35 solved.

### [x] 15.36 Others' aim and hover

1. Goal: everyone sees what the active player is aiming at (T-08 – T-10, T-21, T-22).
2. Tasks:
   1. Others see the arrow live while it moves (T-08).
   2. A hidden card shows as a back lifted from its hand slot (T-09, T-22).
   3. The aimed-at card glows red on the target's screen (T-21).
   4. A hovered card is marked red for the others (T-10).
   5. A late or missing message never breaks the table.
3. Done when:
   1. e2e test: a second player sees the arrow and the red mark.
4. Owner decision (2026-10-20)
   1. A waiting player sends only aims, never hover: less noise.
   2. The active player sends aims and hover.
5. Result (2026-10-20)
   1. `PresenceSender` shares the aim: source, and the spot or board point.
   2. On the own turn it also shares the hovered card.
   3. An own hand card goes as its slot only (A-16).
   4. At most every 70 ms; a still one is renewed every 2 s.
   5. `PresenceLayer` draws the other seats' arrows live (T-08).
   6. A hidden source is a back lifted from its slot (T-09, T-22).
   7. The spot under the head or pointer gets a red frame (T-10, T-21).
   8. Not renewed in 5 s, a presence goes; a lost message heals.
   9. Others' arrows are a little fainter, under every window.
   10. New layer `presence`, under choices, over the hand.
   11. ADR 016 gains the client rules; UC-25; wiki `Client`, `Layers`.
   12. TS protocol test models the presence lines now.
   13. Unit tests in `presence.test.ts`: spots, what is shared, frame points.
   14. e2e test: the other sees the arrow, the lifted back, and the marks.
   15. e2e test: a waiting player's hover reaches nobody.
   16. e2e test: another player's arrow stays under a choice.
   17. Screenshots: `.cache/screenshots/15.36/before/` and `after/`.
6. Owner notes (2026-10-20)
   1. The aim shot framed the monster, but the head was on a shop item.
7. Result of the fixes (2026-10-20)
   1. Cause: the frame showed at once, the arrow one frame later.
   2. The pointer crossed the shop on its way; the shot caught both.
   3. Marks and arrows now change together, in one pass per frame.
   4. e2e test: the head ends on the marked monster's middle.
   5. Screenshots: `.cache/screenshots/15.36/after/`, retaken.
8. Owner notes (2026-10-20)
   1. Owner marks 15.36 solved.

### [x] 15.37 Hit and death marks

1. Goal: damage and deaths read at a glance (T-15, T-24 – T-26).
2. Tasks:
   1. After the target locks, the effect flies from source to target (T-15).
   2. A yellow starburst "-N" on every damaged card at once (T-24).
   3. An HP number turns red while below full (T-25).
   4. A dying card gets a big skull before it leaves (T-26).
   5. Only `transform` and `opacity`; follows animation speed (ST-05).
3. Done when:
   1. `make screenshots` shows a starburst, red HP, and a skull.
4. Result (2026-10-13)
   1. Damage that is not an attack lands as a glowing shot.
   2. The shot flies from the source card, else from its stack row.
   3. All damage of one update lands at once: a "−N" burst on each card.
   4. "Target locks" is when the damage resolves, as for attacks.
   5. Pending damage has no popup now; damage shows when it lands.
   6. Hurt health: a pale heart with a red number (T-25).
   7. A red number on the red heart did not read; hence the pale heart.
   8. A hurt monster shows health even when it matches the print.
   9. A dying card gets a big bone-white skull, the stack's icon (T-26).
   10. The card leaves while the skull fades.
   11. A character kept its top-corner damage mark at first.
   12. Unit tests: damage joins one `hits` effect; deaths get skulls.
   13. e2e tests: shots, bursts at once, red health; a skull, then gone.
   14. Screenshots: `.cache/screenshots/15.37/before/` and `after/`.
5. Owner notes (2026-10-13)
   1. Owner: the pale heart with a red number is fine.
   2. Owner: drop the character's corner damage mark; red health shows it.
6. Result of the fixes (2026-10-13)
   1. Cards with a health badge have no corner damage mark.
   2. A card without health keeps the mark: its only sign of damage.
   3. e2e test: a hurt character has no corner mark.
   4. Screenshots: `.cache/screenshots/15.37/fix/`.
7. Owner notes (2026-10-13)
   1. Owner marks 15.37 solved.

### [x] 15.38 Coin icon

1. Goal: coins show the game's own coin, not a drawn one.
2. Tasks:
   1. Work in `FourSoulsRepo/tools`, never in this repo (as 15.9).
   2. Crop the ¢ face from `assets/coin.png`; drop the pentagram side.
   3. Convert it to WebP, small and sharp at icon size.
   4. Amend the single commit of the images repo, as for cards.
   5. Serve it like the backs; the app gets its URL or none.
   6. Own coins use the image; the drawn ¢ coin is the fallback.
   7. Opponents get the same icon in 15.18.
   8. **Owner** decides: keep or drop the pentagram face and `coin.png`.
3. Done when:
   1. `make screenshots` shows the coin image by own name.
   2. A build without the images shows the drawn coin.
   3. `make frontend e2e` passes.
4. Result (2026-10-12)
   1. `scraper coin` in the tools repo: 128 px WebP, quality 85, 6 KB.
   2. It takes the biggest opaque blob in the left half: the ¢ face.
   3. File: `assets/icons/coin.webp`; `coin.png` is kept.
   4. The icon is staged in the images repo; the owner amends and pushes.
   5. `App.Icons` gives its URL; a missing file gives the drawn coin.
   6. Own coins show the image at 30 px; a failed load shows the drawn coin.
   7. e2e test: one page gets the image, the other the drawn coin.
   8. Screenshots: `.cache/screenshots/15.38/before/` and `after/`.
5. Owner notes (2026-10-12)
   1. Owner: the faint ¢ at 30 px is fine.
   2. Owner: keep the pentagram face and `coin.png`; may be useful later.
   3. Owner pushed the tools and images repos.
   4. Owner marks 15.38 solved.

### [x] 15.39 Souls as crystals

1. Goal: a player's souls read at a glance (V-30, V-31).
2. Tasks:
   1. A row of slots, like mana crystals in Hearthstone (V-30).
   2. Filled slots are gained souls; empty ones show what is left to win.
   3. The row sits left of the character, not right (V-31).
   4. Opponents get the same row, smaller.
   5. A card worth two souls fills two slots.
   6. **Owner** decides: slot count from the win target, or a fixed four.
3. Done when:
   1. e2e test: a gained soul fills one more slot.
   2. `make screenshots` shows the row left of the character.
   3. UC-12 updated.
4. Result (2026-10-13)
   1. Done out of order, at the owner's request.
   2. Owner: slot count from the win target.
   3. `GameView.win_souls` sends `Game.WinSouls`; always 4 today.
   4. A view without it counts four (older server or record).
   5. A column of six-sided gems left of the character, bottom up.
   6. Filled: a glowing blue gem; empty: a dark gem with a pale rim.
   7. Own 16 px wide, opponents 12 px; the same in every seat.
   8. The cell pads its right side, so the card stays centered.
   9. Souls past the target add filled slots.
   10. The "N souls" text is gone; hover reads "Souls: N of 4".
   11. A dead player's status still says "dead".
   12. Pure CSS: the same without the card images.
   13. Unit test: `soulSlots`; engine test: the view's win target.
   14. e2e test: a gained soul fills one more; the row left of the card.
   15. Screenshots: `.cache/screenshots/15.39/before/` and `after/`.
5. Owner notes (2026-10-13)
   1. Owner: the vertical column is fine.
   2. Owner marks 15.39 solved.

### [x] 15.40 Isaac fonts on the title screen

1. Goal: the title screen looks like the game (V-33).
2. Tasks:
   1. Use the fonts from `assets/fonts` for the title.
   2. **Owner** decides which font for which text.
   3. Without images: the plain UI font, as now.
3. Done when:
   1. `make screenshots` shows the new title.
   2. `wails build -tags embed` still shows a readable title.
4. Result (2026-10-14)
   1. Done out of order, at the owner's request.
   2. Owner: the title in tboi4s, the Russian card font.
   3. Owner: menu buttons in font-souls Title, like the game buttons.
   4. Owner: the same fonts on every menu-like screen.
   5. New `--font-title` in `style.css`: tboi4s, then the system font.
   6. Splash and menu title: `--font-title`, 56 px.
   7. Screen headings (Welcome, Host, Join, Games, Table, Settings) too.
   8. Every button on those screens: `--font-game`.
   9. Body text, labels, inputs and the legal notice stay plain.
   10. `font-synthesis: none`: no fake bold on the game font.
   11. Without images: the bold system font, as before.
   12. e2e test: title, buttons and a heading ask for the game fonts.
   13. Shell, welcome and lobby shots now load the real fonts.
   14. Screenshots: `.cache/screenshots/15.40/before/` and `after/`.
5. Owner notes (2026-10-14)
   1. Owner: buttons like the ability lines of The Bone.
   2. Owner: settings and game setup should fit the game too.
6. Rework (2026-10-14)
   1. Menu, screen and overlay buttons: parchment, ink outline, hard shadow.
   2. They tilt a little in turn, like the ability lines.
   3. The main action (Host, Save, Ready) is the gold of End turn.
   4. Disabled: greyed, like a line that is not allowed.
   5. The in-game Menu button gets the same look, smaller.
   6. Panels: a dark board with an ink frame, in both themes.
   7. Overlay boxes (pause, game over, menu) get the same frame.
   8. Field names in the game font, gold; option names light.
   9. Inputs and selects: parchment slips; typed text stays plain.
   10. Checkboxes: parchment boxes with an ink tick.
   11. The look lives in `style.css` (`--ink`, `--paper`, `--gold`).
   12. Screenshots: `.cache/screenshots/15.40/paper-before/` and `paper-after/`.
7. Owner notes after the rework (2026-10-14)
   1. Owner marks 15.40 solved.

### [ ] 15.41 foursouls.ru banner on the home screen

1. Goal: the home screen greets like foursouls.ru (V-34).
2. Tasks:
   1. The tools repo downloads the top image of foursouls.ru.
   2. It goes into the private images repo; amend its one commit.
   3. The home screen shows it at the top.
   4. Without images: the home screen as now.
   5. Credit it in `NOTICE`.
3. Done when:
   1. `make screenshots` shows the banner.
   2. e2e passes without images.

### [ ] 15.42 Animated background

1. Goal: home and setup screens feel alive (V-35).
2. Tasks:
   1. Study the background animation of foursouls.ru.
   2. Rebuild it with our own CSS; copy no code.
   3. Use it on the home screen and the game setup screens.
   4. Respect the animation speed setting (ST-05).
   5. Without images: a plain background.
3. Done when:
   1. `make screenshots` shows both screens.
   2. Owner compares it with the site.

### [ ] 15.43 First player roll for everyone

1. Goal: everyone sees who goes first and why (GS-06).
2. Tasks:
   1. Show each player's first-player roll as real dice (V-17).
   2. Ties roll again; those rolls show too.
   3. The winner of the roll is named; then the game starts.
   4. A character that goes first (Cain) says so instead.
   5. The history keeps the rolls.
3. Done when:
   1. e2e test: every player's roll shows before the first turn.
   2. Use case.

### [ ] 15.44 Declarations on the stack

1. Goal: everyone sees a "Buy" or "Attack" before it resolves (V-36).
2. Tasks:
   1. List every declaration the rules put on the stack.
   2. Each one shows as a stack entry with player and target.
   3. Same icons in the history strip (15.13).
3. Done when:
   1. e2e test: a purchase and an attack show on the stack.
   2. UC for the stack panel updated.

### [ ] 15.45 Look at the table after the game

1. Goal: players discuss the finished game (V-37).
2. Tasks:
   1. After the game, the view filter reveals all hands (A-07).
   2. The end overlay can be closed and opened again.
   3. The table stays live until the player leaves.
   4. Spectators see the same.
3. Done when:
   1. Engine test: a finished game's view shows every hand.
   2. e2e test: close the overlay, see an opponent's hand.
   3. Use case.

### [ ] 15.46 Failed effects in the history

1. Goal: nothing fails silently (V-38).
2. Tasks:
   1. A fizzled ability gets a history entry.
   2. A failed purchase too: the player could not pay.
   3. The entry says what failed and why.
   4. Check which other events end silently today.
3. Done when:
   1. Vitest: a fizzle event becomes a history entry.
   2. e2e test: a failed purchase shows in the history.

### [x] 15.47 Who deals damage a player takes

1. Goal: the history names the real dealer (play test, 2026-10-12).
2. Bug: Boom Fly died; the line said "Harddie deals 1 damage to Harddie".
3. Owner: fix it now, before 15.37; write the decision down.
4. Tasks:
   1. Find the cause in the match record.
   2. Name the source card when a player takes damage.
   3. ADR for the choice.
5. Done when:
   1. Vitest: Boom Fly's item names the fly; no card says "takes".
   2. e2e test: the history strip names Boom Fly.
   3. Use case UC-14 updated.
6. Result (2026-10-13)
   1. Record: each damage item had controller = target = the player.
   2. `EachPlayer(DealDamage(1, You))` makes each player the controller.
   3. The engine is right; `stackParts` named the controller first.
   4. Now damage to its own controller names the source card.
   5. No source card: "Ann takes 1 damage".
   6. Engine, protocol and records unchanged; ADR 011.
   7. The hover's owner line still shows the controller.
   8. Screenshots: `.cache/screenshots/15.47/before/` and `after/`.
7. Owner notes (2026-10-13)
   1. Owner marks 15.47 solved.

### [x] 15.48 Player names in record file names

1. Goal: a record file says who played (owner, 2026-10-13).
2. Owner: fix it now; clean the names for security and file systems.
3. Tasks:
   1. The server adds seat nicknames to the record file name.
   2. Nicknames keep letters and digits; other runs become `-`.
   3. The app suggests the server's file name, cleaned again.
   4. ADR 007 and UC-08 updated.
4. Done when:
   1. Go tests: hostile nicknames give a safe file name.
   2. Go test: the app's save dialog suggests the cleaned name.
5. Result (2026-10-13)
   1. Name: `2026-10-13-120000-g3_Ann_Bob.fsrec`.
   2. `record.FileName` and `SafeName` in `pkg/record`.
   3. Kept: Unicode letters and digits; at most 20 per nickname.
   4. Dropped: dots, slashes, spaces, control and bidi runes.
   5. Empty after cleaning: `p1`, `p2`, ….
   6. `Content-Disposition` is built by `mime`; Cyrillic names survive.
   7. The app downloads first, then the dialog suggests that name.
   8. The app cleans it again with `record.SafeFileName`.
   9. No screen changed; no screenshots.
6. Owner notes (2026-10-13)
   1. Owner marks 15.48 solved.

### [x] 15.49 Wait for my turn

1. Goal: a player skips a whole turn of others in one press (N-16).
2. Owner: fix it now; follow the step's rules.
3. Tasks:
   1. Button "Wait for my turn" for a player on another's turn.
   2. The server passes for them, as for "Skip all".
   3. It ends when the player before them ends the turn.
   4. Then they are asked once more: the last chance to act.
4. Done when:
   1. Go test: passes while waiting; the end of the turn before stops it.
   2. Vitest: when the button is offered.
   3. e2e test: on, cancel, and asked again at the end of the turn.
5. Result (2026-10-13)
   1. New message `wait_turn`; `update` carries `wait_turn` (ADR 006).
   2. Protocol version unchanged: only added.
   3. "The turn before": the seat before, in seat order.
   4. It ends at the declared end, a death, or the end phase.
   5. It ends on the player's own turn too, whatever came first.
   6. Not offered while the turn before is ending.
   7. Choices and discards still ask the player.
   8. Leaving the table turns it off, as "Skip all".
   9. Button sits by "Skip all"; it says "Cancel wait" while on.
   10. Screenshots: `.cache/screenshots/15.49/before/` and `after/`.
6. Owner notes (2026-10-13)
   1. Owner marks 15.49 solved.

### [x] 15.50 Damage sounds

1. Goal: damage is heard as well as seen (AU-03).
2. Owner: do it now; follow the step's rules.
3. Tasks:
   1. Find a card-hits-card sound; download it into the assets.
   2. Find Hearthstone's Jaina fireball for bombs and effects.
   3. No background music; convert the files small.
   4. Play the hit on attack damage, the fireball on other damage.
4. Done when:
   1. Go test: the sounds are served; none without the private images.
   2. Vitest: the fallback when a fireball file is missing.
   3. e2e test: an attack plays the hit; card damage the fireball.
5. Result (2026-10-13)
   1. Card hit: Kenney "Casino Audio", `card-place-1`, CC0.
   2. It ships with the frontend: `src/assets/sounds/card-hit.mp3`.
   3. Fireball: Hearthstone's Fireball spell effects, from hearthstone.wiki.gg.
   4. Files `FX_FireballEvent01_BirthStart_01` and `…04_SpellImpact_01`.
   5. Sound effects only: no voice, no music.
   6. They belong to Blizzard: private images only, `assets/sounds/`.
   7. Served by `App.Sounds()`, like the fonts.
   8. The cast plays as shots leave; the impact as they land.
   9. All shots of one update share one fireball.
   10. Without the private files the impact is the card hit; no cast.
   11. All mono MP3, 64 kbit/s: 5 KB, 7 KB, 13 KB.
   12. Sounds follow animations; with them off, nothing plays.
   13. No volume or mute setting yet; volume is fixed at 60 %.
   14. No screenshots: the screens do not change.
6. Owner notes (2026-10-13)
   1. Owner: the cast is fine.
   2. Owner: the card hit on attacks sounds bad and has a wrong name.
   3. Owner: take it from sounds.spriters-resource.com, Hearthstone.
   4. Owner picks `WoW_unarmedparrymetala` for attacks, card to card.
   5. Owner: keep the effect sounds; give the files proper names.
7. Result of the fixes (2026-10-13)
   1. The site has no file named "kick"; 11 impacts went to the owner.
   2. Names by use: `attack-hit`, `effect-cast`, `effect-impact`.
   3. All three are Hearthstone's: private images only.
   4. Kenney's card hit is gone; a public clone is silent.
   5. `attack-hit.mp3`: 4 KB.
8. Owner notes (2026-10-13)
   1. Owner merged the sounds into the private images.
   2. Owner marks 15.50 solved.

### [ ] 15.51 Stack and history rows move

1. Goal: rows of the stack and history never jump (V-39).
2. Tasks:
   1. A new row fades in.
   2. Other rows slide to their new place.
   3. A leaving row fades out, then the rest close the gap.
   4. Stack panel and history strip both.
   5. Follow the animation speed setting (ST-05); off means none.
   6. Rows keep a stable key, so React moves them, not redraws.
3. Done when:
   1. Vitest: rows keep their keys when the stack grows and shrinks.
   2. e2e test: a resolved row fades out; the next slides up.
   3. e2e test: with animations off, rows change at once.

### [ ] 15.52 Played cards fly through the stack

1. Goal: everyone follows a played card to the end (V-40).
2. Tasks:
   1. A played card flies from its hand to the stack panel.
   2. Opponents' cards fly from their hand, face up.
   3. At the stack the card becomes its row.
   4. **Owner** picks the effect: shrink into the row, or another.
   5. Resolved, the row turns back into the card.
   6. The card flies to its discard pile.
   7. Damage and death marks (15.37) play first; then the flight.
   8. Cards that stay in play (items, abilities): the row fades out (15.51).
   9. Use the playback queue (`fly` effect); follow ST-05.
3. Done when:
   1. Vitest: a resolved loot gives hits, then the discard flight.
   2. e2e test: a played loot flies to the stack, then to the discard.
   3. `make screenshots` shows the card mid-flight.

### [x] 15.53 Create a game on a dedicated server

1. Goal: a player opens a game on a dedicated server from the app (N-17).
2. Owner: do it now; follow the step's rules.
3. Tasks:
   1. `welcome` says whether this client may create a game.
   2. A dedicated server lets every client create.
   3. A server hosted in the app lets only the host's own app create.
   4. The same for loading a saved game there; friends still join it.
   5. Games list: "Create a game" only when the server allows it.
   6. It opens the host setup; the game opens on that server.
   7. No server starts in the app; Back returns to the Games list.
   8. Fix the VPS wiki page, which says players create games there.
4. Done when:
   1. Go test: a hosted server refuses a friend's create and load.
   2. Go test: a dedicated server lets anyone create.
   3. Vitest: when the button shows.
   4. e2e test: create on the test server; another player joins.
5. Result (2026-10-14)
   1. `welcome` carries `create`; protocol version unchanged: only added.
   2. `Config.HostOnly`: only loopback clients create or load.
   3. Only the app's `Host` sets it; `cmd/server` never does.
   4. Others get `not_host`; they still join and continue loaded saves.
   5. The host's app always connects over `127.0.0.1`, so it may.
   6. A second copy on the host's computer counts as local too (UC-16).
   7. The setup screen says "Create a game" and "On the server …".
   8. No nickname field there: it went with hello.
   9. No "Continue a saved game" there: saves are in the Games list.
   10. ADR 006, UC-03, new UC-19, Server and VPS wiki pages updated.
   11. Screenshots: `.cache/screenshots/15.53/before/` and `after/`.
6. Owner notes (2026-10-14)
   1. Owner marks 15.53 solved.

### [ ] 15.54 Sources of stat boosts

1. Goal: a player sees why a character's HP or ATK changed (V-41).
2. Tasks:
   1. The engine view lists each source of a character's HP and ATK change.
   2. A source: an item, a curse, a soul, a played effect "till end of turn".
   3. Small icons on the character card, beside its HP or ATK badge.
   4. One icon per source; hover names the card and its change.
   5. Hovering an icon zooms the source card.
   6. **Owner** decides: icon from the card art, or a plain mark.
   7. **Owner** decides: also on monsters, or only characters.
   8. Text cards without images get the same icons.
3. Done when:
   1. Engine test: a boosted character's view lists its sources.
   2. e2e test: an item with +1 ATK shows one icon by the ATK badge.
   3. `make screenshots` shows a character with two sources.
   4. UC-12 updated.
4. Owner decisions (2026-10-14)
   1. Icons use the source card's art.
   2. Monsters get the icons too, not only characters.
   3. Reference: the stat icons on fouronlinesouls.net characters.
   4. Study them first; borrow ideas only, no assets (15.1).

### [x] 15.55 Curses left of the player

1. Goal: curses read apart from items (V-42).
2. Tasks:
   1. Curses leave the items row on the right.
   2. They form their own group left of the player.
   3. Own area and opponents alike.
   4. **Owner** decides: left of the soul crystals, or between them and the card.
   5. The table areas keep their size when there is no curse (V-23).
3. Done when:
   1. e2e test: a curse sits left of the character, not with the items.
   2. `make screenshots` shows a cursed player in every seat.
   3. UC-12 updated.
4. Owner decisions (2026-10-14)
   1. Curses go at the far left of the player area.
   2. Not next to the soul crystals.
5. Result (2026-10-15)
   1. Done before 15.54, at the owner's request.
   2. The items row skips curses (`isCurse`, `table.ts`).
   3. `Curses` in `PlayerArea.tsx`: their own group, same card size.
   4. Own area and top: far left of the status cell.
   5. Side columns: a row above the character, at its left.
   6. No curse: no group; the areas keep their size.
   7. e2e: a curse in every seat, left of the crystals; sizes unchanged.
   8. UC-12, Client wiki, CHANGELOG updated.
   9. A flaky aim test, also on the base commit, logged as 19.20.
   10. Screenshots: `.cache/screenshots/15.55/before/` and `after/`.
6. Owner notes (2026-10-15)
   1. Owner marks 15.55 solved.

### [x] 15.56 Soul cards behind the crystals

1. Goal: the table shows no soul cards; the crystals hold them (V-43).
2. Tasks:
   1. Gained soul cards leave the items row.
   2. Each filled crystal knows its soul card.
   3. A card worth two souls ties both of its crystals to it.
   4. Hovering a filled crystal zooms its card.
   5. A killed monster zooms with the soul it gave.
   6. Touch: a long tap does the same.
   7. Effects that aim at a soul still can aim at the crystal.
   8. **Owner** decides: the zoom for a bonus soul, which is no monster.
3. Done when:
   1. Vitest: crystals map to soul cards, two-soul cards included.
   2. e2e test: no soul card in the items row; hover zooms it.
   3. `make screenshots` shows the zoom over a crystal.
   4. UC-12 updated.
4. Owner decisions (2026-10-14)
   1. A bonus soul zooms from its crystal too, like a monster soul.
5. Result (2026-10-15)
   1. Done out of order, at the owner's request.
   2. Engine: `CardView.souls` is what a soul or a character is worth.
   3. `soulsOf` and `characterSouls` (`turn.go`) also sum `SoulValue`.
   4. The items row skips soul cards (`isSoul`, `table.ts`).
   5. `soulCards` ties each filled crystal to its card, bottom up.
   6. A character's own soul value (The Lost) comes first.
   7. A two-soul card holds two crystals; a view without worth counts one.
   8. Hover or long tap zooms the card: `useZoom`, shared with `Card`.
   9. A killed monster is its own soul, so it zooms as the monster.
   10. A soul's crystal has the card's zone and pick handle.
   11. A highlighted crystal grows and glows in the highlight color.
   12. The choice dialog covers the own crystals; it lists the souls too.
   13. Engine test: each soul's worth; Vitest: `soulCards`.
   14. e2e: no soul card in the items row; hover zooms; a pick by crystal.
   15. UC-12, Client wiki, CHANGELOG updated.
   16. Screenshots: `.cache/screenshots/15.56/before/` and `after/`.
6. Owner notes (2026-10-15)
   1. Owner marks 15.56 solved.

### [ ] 15.57 Purchases in the history

1. Goal: everyone sees what a player bought (V-44).
2. Tasks:
   1. Check which event a purchase sends today.
   2. The engine sends a purchase event if none fits.
   3. The full History gets a line: "A buys <item>".
   4. A bought shop top card and a bought deck top card both count.
   5. Hovering the line zooms the item card.
   6. **Owner** decides: also an icon in the history strip, or the full History only.
3. Done when:
   1. Vitest: a purchase event becomes a history entry.
   2. e2e test: buy an item; the History names it.
   3. UC-12 updated.

### [ ] 15.58 "Choose one-" cards ask like The Bone

1. Goal: picking a mode looks like picking an ability (V-45).
2. Tasks:
   1. The engine names the card in the mode prompt's `subject`.
   2. Only the chooser sees it, as with other subjects.
   3. The mode dialog shows that card large, like the ability dialog.
   4. One line per mode, styled like The Bone's lines.
   5. Cancel and a click outside still send the prompt's cancel.
   6. Played from hand (Dagaz) and from an item (Chaos Card) alike.
   7. A dragged card still asks its mode first, then aims (15.26).
   8. **Owner** decides: modes with no valid target hidden, or greyed.
3. Done when:
   1. Engine test: a mode prompt names its card as the subject.
   2. e2e test: Dagaz opens the dialog with the card; the picked mode is sent.
   3. `make screenshots` shows the mode dialog.
   4. UC-13 updated.

### [x] 15.59 Throw in the engine

1. Goal: every roll waits for the roller's throw (V-17, 15.5).
2. Tasks:
   1. A throw prompt and a throw intent.
   2. The value is drawn when the intent arrives, never before.
   3. Covers attack, reward, effect, and many-dice rolls.
   4. A reroll asks the roll's own player to throw again.
   5. The first-player roll stays automatic until 15.43.
   6. Test runner and fuzz harness throw for the roller.
   7. ADR for the throw intent.
   8. Old records cannot replay: bump the record format, or refuse them clearly.
3. Done when:
   1. Engine test: no value exists before the throw.
   2. Engine test: forced rolls still decide the value.
   3. All engine tests and fuzz runs pass.
4. Result (2026-10-15)
   1. `PromptThrow` and `IntentThrow`; ADR 014.
   2. A roll goes on the stack with no value and queues `ActThrow`.
   3. Queued actions run first: nobody responds before the throw.
   4. The throw draws the value; then `dice_rolled` with the stack ID.
   5. Attack, reward, `Roll` effect rolls; many dice are drawn with the throw.
   6. `Reroll` and `RerollRoll` clear the value; the roll's player throws.
   7. A roll gone before its throw is never thrown.
   8. A game-controlled roll is thrown at once; the first-player roll unchanged.
   9. Test helpers and `settle` throw; fuzz throws through `Allowed`.
   10. Engine 0.2.0; `CanReplay` refuses older records clearly.
   11. Record format stays 1, so old records stay readable (ADR 007).
   12. Server `forcedIntent` throws, so the build stays exhaustive; tests in 15.60.
   13. The client has no throw yet: the app waits at every roll until 15.61.
   14. ADR 005, ADR 007, Adding-a-card wiki, CHANGELOG updated.
   15. No screen change; no screenshots.
5. Owner notes (2026-10-15)
   1. Owner marks 15.59 solved.
   2. 15.59 – 15.61 stay together and land together.
   3. The app and e2e wait at rolls until 15.61; e2e gates it then.

### [x] 15.60 Throw on the server

1. Goal: the server carries the throw and answers for idle seats.
2. Tasks:
   1. Throw intent in the protocol; protocol docs and golden files.
   2. A kicked seat throws at once.
   3. The response timer throws for an idle player.
   4. "Skip all" and "wait for my turn" never throw for a player.
   5. Records keep the throw intent.
   6. UC-06 and UC-07 updated.
3. Done when:
   1. Server test: the timer throws; a kicked seat throws.
   2. Server test: a record with throws replays to the same game.
4. Result (2026-10-15)
   1. Protocol 3: the throw prompt and intent; golden line added; ADR 006.
   2. The server needed no new code: `forcedIntent` throws since 15.59.
   3. Kicked seats and the response timer throw through it.
   4. `autoIntent` passes only at a priority prompt, so skip and wait never throw.
   5. Records keep every intent, the throw too.
   6. Tests: `TestKickedSeatThrows`, `TestTimerThrows`, `TestSkipAndWaitNeverThrow`.
   7. `TestRecordReplays` now attacks, so its record holds throws and replays.
   8. Client: protocol 3 and the throw enums; no button yet (15.61).
   9. UC-06, UC-07, Server wiki, CHANGELOG updated.
   10. No screen change; no screenshots.
5. Owner notes (2026-10-15)
   1. Owner marks 15.60 solved.

### [x] 15.61 Throw button and waiting

1. Goal: the roller always sees the game waits for their throw (V-17).
2. Tasks:
   1. A big glowing Throw button where the die will land.
   2. The stack panel says "Waiting for you to throw" (V-11).
   3. Others see "Ann is throwing…" and a shaking die.
   4. The die uses the player's dice look (V-19).
   5. e2e helpers throw when a test waits for a roll.
   6. UC-13 and UC-14 updated.
3. Done when:
   1. e2e test: attack, throw, the result matches the engine.
   2. `make screenshots` shows the button and the others' view.
4. Result (2026-10-15)
   1. `Throw.tsx`: the roller's big glowing Throw button, mid-screen.
   2. The button holds the player's die; the roll's name shows below.
   3. Space throws too; it never passes during a throw prompt.
   4. Others see "Ann is throwing…" and the die shaking (`throw-shake`).
   5. The halo breathes (`throw-glow`); both only `transform` and `opacity`.
   6. Hidden while an effect plays, so only the landing die shows.
   7. Stack panel: "Waiting for you to throw", or "Ann is throwing…".
   8. Action bar: "Throw the die", or "Ann is throwing".
   9. The die is the viewer's dice look until 15.6 sends each seat's.
   10. e2e `answer` and the random-games driver throw.
   11. e2e test: attack, throw; both screens get the shown value.
   12. Full `make e2e` passes: the app no longer waits at rolls.
   13. UC-13, UC-14, Client and Animations wiki, CHANGELOG updated.
   14. Screenshots: `screen-throw-button.png`, `screen-throw-others.png`.
5. Owner notes (2026-10-15)
   1. Owner marks 15.61 solved.

### [x] 15.62 Throw animation

1. Goal: the die falls, bounces, and rolls to the engine's value (V-17).
2. Tasks:
   1. Only `transform` and `opacity`.
   2. Follows the animation speed setting (ST-05).
   3. Tag and row in `docs/wiki/Animations.md` (ADR 012).
   4. Plain die without a texture, as now.
3. Done when:
   1. e2e screenshot of a die in the air.
   2. The rolled value matches the engine in every test.
4. Result (2026-10-15)
   1. `ThrownDie` in `Effects.tsx`: `dice_rolled` (fresh roll or reroll).
   2. Falls from above, spinning; bounces twice; lies still; fades.
   3. Web Animations on `transform` and `opacity`; tag `dice-throw`.
   4. It tumbles through two other faces and lands on the value (`tumbleFaces`).
   5. Its label says "Rolled N" all along: the engine's value.
   6. Length scales with the animation speed; Off skips it, as every effect.
   7. A plain die without a texture, as before.
   8. `roll_changed` keeps the old spin-in (`dice`).
   9. e2e test: held in the air, then landed; both say the engine's value.
   10. The 15.61 test still checks the value on both screens.
   11. UC-13, Animations wiki, CHANGELOG updated.
   12. Screenshots: `screen-dice-in-air.png`, `screen-dice-landed.png`.
   13. The fast-burst e2e test looks for its fly on every frame.
   14. The longer throw left the fly too short for `toBeVisible`'s polling.
5. Owner notes (2026-10-15)
   1. Owner marks 15.62 solved.

### [x] 15.63 Last roll by the stack

1. Goal: the last roll stays visible as a small die by the stack (F-21).
2. Tasks:
   1. Shows the final value, after changes and rerolls.
   2. Hovering it names the roller and what the roll was for.
3. Done when:
   1. e2e test: after a roll the small die shows its value.
   2. `make screenshots` shows it.
4. Result (2026-10-15)
   1. A small die right of the "Stack" title, in the viewer's dice look.
   2. `lastRollAfter` (`dice.ts`) follows the shown events.
   3. A throw starts a new last roll; reroll, change and result set its value.
   4. Hover or focus: "Ann rolled 4 — attack roll" (`lastRollText`).
   5. None before the first roll; reset with the history.
   6. Unit tests: the fold and the text.
   7. e2e: after the attack, both screens show the final value; hover names it.
   8. UC-14, Client wiki, CHANGELOG updated.
   9. Screenshot: `screen-last-roll.png`.
5. Owner notes (2026-10-15)
   1. Owner marks 15.63 solved.

### [x] 15.64 Dice look in the settings

1. Goal: the player picks their dice look in the settings (V-19).
2. Asked by the owner after 15.6 (2026-10-15).
3. Tasks:
   1. Every d6 of the holy and unholy sets in the images.
   2. A dice field in Settings, grouped by set.
   3. A preview of the picked look's faces.
   4. Without the images: no field; dice stay plain.
   5. The pick goes out in the next hello (15.6).
4. Done when:
   1. Go test: every look's texture is found and served.
   2. e2e test: pick a look, save, the preview and setting follow.
5. Result (2026-10-15)
   1. `diceSkins` and `settings.Dice`: ten looks; Tear stays the default.
   2. Holy: Tear, Blue Fly, Cent, Mom's Knife, The Bible.
   3. Unholy: Brimstone, Guppy's Paw, Pentagram, Poop, Tinted Rock.
   4. All ten share the tear's 3×3 layout; `Die` needs no change.
   5. `DiceSkin.Set` names the group; `diceChoices` (`dice.ts`) groups them.
   6. Settings: a Dice field with Holy and Unholy groups.
   7. Under the mat preview, the six faces of the picked look.
   8. No texture: no field and no preview.
   9. Fake bridge: all ten looks; `?nodice` has no textures.
   10. Tests: `TestDiceTexturesServed` covers all ten; `diceChoices` unit tests.
   11. e2e: pick Poop, save, reload; without textures, no field.
   12. UC-10, Client wiki, CHANGELOG updated.
   13. Screenshots: `.cache/screenshots/15.64/before/` and `after/`.
6. Owner notes (2026-10-15)
   1. Every dice texture is upside down; turn it 180°.
   2. Fixed: `faceStyle` turns the face about its own middle.
   3. Plain dice and the throw animation are unchanged.
   4. Screenshot: `.cache/screenshots/15.64/rotated/`.
   5. Owner marks 15.64 solved.

### [x] 15.65 Cathedral mat image

1. Goal: the owner's new Cathedral picture is a mat image (V-24).
2. Asked by the owner after 15.64 (2026-10-19).
3. Tasks:
   1. Convert `assets/mats/cathedral.png` with `scraper assets` (tools repo).
   2. The fake bridge lists it, so screenshots show it.
4. Done when:
   1. Settings and the table show the Cathedral image.
5. Result (2026-10-19)
   1. 4000×2000 PNG to 3840×1920 WebP: 4.5 MB to 0.2 MB.
   2. Original kept in `.cache/assets-originals-2026-10-15/mats/`.
   3. No code change: an image mat named `cathedral` replaces its color mat.
   4. Basement stays the default: it is the first image.
   5. Owner amended the images repo's single commit (`3ed5075`).
   6. Screenshots: `.cache/screenshots/mat-cathedral/before/` and `after/`.

### [ ] 15.66 Choice cards without name buttons

1. Goal: a card option is only its card, with no button repeating its name.
2. Asked by the owner after 15.29 (2026-10-19).
3. Tasks:
   1. Card options: the big card is the button; no name button under it.
   2. Text options ("yes", "keep", "monster deck") keep their buttons.
   3. Ability lines and outcomes under an option stay.
   4. Discard and order dialogs: check they repeat no names either.
   5. e2e helpers that click `.choice-option button` also click the card.
4. Done when:
   1. e2e test: a card option has no name button and its card answers.
   2. `make screenshots` covers a card choice and a text choice.

### [x] 15.67 Kept die numbers on the card

1. Goal: a card that keeps a die number shows it as a die badge.
2. Asked by the owner (2026-10-19).
3. Cards:
   1. Telepathy for Dummies: the roll stored at the start of the turn.
   2. Red Candle: the chosen number, till the end of the turn.
   3. Crystal Ball: the guessed number, till the next roll.
4. Tasks:
   1. Engine: `CardDef.Dice` gives the numbers a card keeps.
   2. The card view sends them as `dice`.
   3. Crystal Ball keeps its guess outside counters; no counter badge.
   4. Client: a die badge per number, in the controller's dice look.
   5. Plain die with the number without a dice image.
   6. No badge on the zoomed card, as for the other badges.
5. Done when:
   1. Engine tests: the three cards' view `dice`.
   2. e2e test: a die badge on an item that keeps a number.
   3. `make screenshots` covers a card with a die badge.
6. Result (2026-10-19)
   1. Engine: `CardDef.Dice`; the card view sends `dice`.
   2. Telepathy reads its stored roll; Red Candle its numbers this turn.
   3. Crystal Ball's guess moved to `Object.Chosen`; it was a counter.
   4. So it no longer shows a green counter or counts as one.
   5. Client: a die per number on the card art, right side.
   6. Drawn in the owner's dice look, with the badges' ink edge.
   7. Fix: `.card img` stretched the die texture; now `.card > img`.
   8. Engine tests: view `dice` of the three cards, gone after use.
   9. e2e test: `table.spec.ts`, "a card that keeps a die number".
   10. UC-12 2f; wiki `Client`, `Effect-blocks`.
   11. Screenshots: `.cache/screenshots/15.67/before/` and `after/`.
7. Owner notes (2026-10-19)
   1. Owner marks 15.67 solved.

### [ ] 15.68 Every view field is read by the client

1. Goal: a field the server sends but the client ignores fails a test (TS-07).
2. Asked by the owner after 19.30 (2026-10-19).
3. Tasks:
   1. Go test: reflect the JSON fields of `GameView` and the types under it.
   2. It fails when a field is missing from `frontend/src/protocol.ts`.
   3. Vitest: each such field is read outside `protocol.ts`.
   4. An allow list names fields unread on purpose, each with a reason.
   5. Fix or list every field the first run finds.
   6. Wiki `Frontend-tests.md`: how to add a field.
4. Done when:
   1. Removing the `treasure_top` read from `Board.tsx` fails the Vitest.
   2. Adding a field to `PlayerView` in Go fails the Go test.
   3. `make lint test frontend` pass.

### [x] 15.69 Dice in 3D, with a sound

1. Goal: dice look like solid cubes, and a throw is heard (V-17).
2. Asked by the owner (2026-10-19).
3. Tasks:
   1. Every die gets light and shade on its face: a bevel.
   2. The thrown die is a cube of six faces that tumbles in 3D.
   3. It lands on the engine's value, tilted to show two more faces.
   4. A shadow on the table grows as it falls.
   5. A plastic die on wood plays at the first landing.
   6. The sound's source is unknown: private images only.
   7. Without the texture: a plain cube, numbers on its faces.
4. Done when:
   1. e2e: the thrown die is a cube, says the value, plays the sound.
   2. Go test: the sound is served; none without the private images.
   3. Screenshots before and after.
5. Result (2026-10-19)
   1. `Die.css`: every die gets a bevel and a gloss, sized in `cqi`.
   2. `DieCube` in `Die.tsx`: six faces with `preserve-3d`.
   3. `cubeFaces` (`dice.ts`): value in front; opposite faces add up to 7.
   4. It replaces `tumbleFaces`: the cube really turns now.
   5. The throw: fall, tumble, shadow and label share one length.
   6. At rest: `rotateX(-20deg) rotateY(-24deg)`, top and right showing.
   7. Faces are shaded by side; light comes from above.
   8. Sound: `dice-roll.mp3`, mono MP3 64 kbit/s, 3 KB.
   9. It plays at the first landing; served by `App.Sounds()`.
   10. The file sits in the submodule, not committed there yet.
   11. e2e: six faces, "Rolled 4", and the sound in `playback.spec.ts`.
   12. `actions.spec.ts` reads the look from the cube's front face.
   13. NOTICE, UC-13, wikis `Client` and `Animations`, CHANGELOG updated.
   14. Screenshots: `.cache/screenshots/15.69/before/` and `after/`.
   15. Playwright's headless WebKit draws any CSS cube flat: one face shows.
   16. A textbook cube does the same there; Chromium draws it right.
   17. Owner saw a test cube in real Safari: it draws in 3D.
6. Owner notes (2026-10-20)
   1. Owner merged the dice sound into the private images.
   2. Owner marks 15.69 solved.

### [x] 15.70 Ask before ending the turn

1. Goal: End turn asks first while the player can still play (V-46).
2. Asked by the owner (2026-10-20).
3. Owner answers (2026-10-20)
   1. Ask whenever End turn is not glowing: anything allowed counts.
   2. The question lists what is left; a setting turns it off.
4. Tasks:
   1. `actions.ts`: list what is left besides passing and ending.
   2. A dialog: "End turn?" with that list.
   3. Buttons: End turn and Keep playing.
   4. E or Enter confirms; Esc or a click outside cancels.
   5. A glowing End turn ends at once, as now.
   6. Setting "Ask before ending the turn", on by default.
   7. The dialog gets a row in `docs/wiki/Layers.md`.
5. Done when:
   1. Vitest: the list for attack, buy, loot and abilities.
   2. e2e: End turn with plays left asks; Keep playing keeps the turn.
   3. e2e: with the setting off, End turn ends at once.
   4. `make screenshots` covers the dialog.
6. Result (2026-10-20)
   1. `actionsFor` fills `TurnActions.left` when End turn does not glow.
   2. Order: Attack, Buy, Play loot, Use abilities.
   3. `EndTurnQuestion` in `TurnActions.tsx`: "End turn?" and that list.
   4. It opens through a portal to `body`, so the board does not clip it.
   5. A new update closes it; nothing is sent until End turn.
   6. Setting `no_end_turn_ask`: false by default, so old files ask.
   7. Settings checkbox: "Ask before ending the turn while plays are left".
   8. Layers: it joins the `overlay` row; no new layer was needed.
   9. e2e helper `endTurn` confirms the question; `games.spec` answers it.
   10. Fake bridge `?noendask` starts with the question off.
   11. Vitest `actions.test.ts`; e2e `actions.spec.ts`, `settings.spec.ts`.
   12. Go test `TestNoEndTurnAsk`.
   13. UC-13 7e, UC-10 3g; wikis `Client` and `Layers`; CHANGELOG.
   14. Screenshots: `.cache/screenshots/15.70/before/` and `after/`.
7. Owner notes (2026-10-20)
   1. Owner marks 15.70 solved.
