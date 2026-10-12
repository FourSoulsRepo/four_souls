# Step 15. Visual refresh

Goal: a clearer table, using fouronlinesouls.net as a second reference.

Ideas: V-09, V-12, V-13, V-15 – V-29.

Order: 15.1, 15.15, 15.9 – 15.11, 15.16, 15.17, 15.38, 15.18 – 15.24,
15.3, 15.37, 15.25, 15.32 – 15.34, 15.26, 15.27, 15.5, 15.6,
15.28 – 15.31, 15.35, 15.36.
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

### [ ] 15.3 Attack animation

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

### [~] 15.4 Hearthstone-style card play

1. Goal: cards are played by dragging them to a place or target (V-16).
2. Split into 15.25 – 15.27; tasks moved there.

### [ ] 15.5 Dice throw

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

### [ ] 15.6 Each player's own dice

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

### [ ] 15.20 End turn button and keys

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

### [ ] 15.21 Stack panel

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

### [ ] 15.22 History strip

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

### [ ] 15.23 Counter badges and live stats

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

### [ ] 15.24 Ability dialog and status marks

1. Goal: clear choice of ability and clear card states (C-04, C-05, H-22).
2. Tasks:
   1. Two abilities: a dialog with the big card; pick one line (C-05).
   2. Status marks on cards, e.g. not usable yet (H-22).
   3. Used items stay dimmed, not rotated; already so (C-04).
   4. A character that can attack has a green frame; not yet: "zzz" (T-11).
3. Done when:
   1. e2e test: a two-ability card opens the dialog and uses the picked line.

### [ ] 15.25 Drag and aim

1. Goal: play a card by dragging it to the table or a target (V-16).
2. Tasks:
   1. Apply the targeting ideas the owner picked in 15.15.
   2. Drag a card onto the table or onto its target.
   3. Aiming arrow; valid targets glow, others stay plain.
   4. Works with mouse and touch (B-06).
   5. Touch: drag from the source, release on the target (T-17).
3. Done when:
   1. e2e test: a bomb dragged onto a monster damages it.

### [ ] 15.26 Click modes

1. Goal: the player picks one-click or drag play (V-28).
2. Tasks:
   1. Click card, then click target, still works.
   2. One-click mode: a click plays the card at once.
   3. A setting chooses the mode; both send the same intent.
3. Done when:
   1. e2e test: the bomb play by two clicks.
   2. e2e test: the same play in one-click mode.

### [ ] 15.27 Hand hover and glow

1. Goal: the hand shows what can be played and what it does.
2. Tasks:
   1. Hover a hand card: it rises large, keyword tips beside it (H-13).
   2. Zoom pops up right beside the hovered card, not at a screen side (F-16, T-27).
   3. Playable cards glow green; the rest stay plain (H-12).
3. Done when:
   1. e2e test: playable cards glow, others do not.

### [ ] 15.28 Prompt banners

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

### [ ] 15.29 Reasons and choice dialogs

1. Goal: the player sees why something is not possible, and chooses calmly.
2. Tasks:
   1. Hover an unusable card: the zoom says why (F-17, C-06).
   2. A refused action shows its reason next to the clicked item (F-18).
   3. Choices: big cards over a dimmed board (F-20, H-11).
   4. A "View board" button hides the choice to look at the table (F-20).
3. Done when:
   1. e2e test: an unusable card's zoom shows the reason.
   2. e2e test: "View board" hides and restores a choice.

### [ ] 15.30 Activity status

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

### [ ] 15.31 Fight marks and played cards

1. Goal: fights and played cards are easy to follow for everyone.
2. Tasks:
   1. In a fight, the attacker and the monster get a red frame (F-19).
   2. Attack rolls, targets, and damage are visible to all.
   3. A card another player plays shows large for a moment (H-10).
   4. It stays large while its effect resolves (T-16).
3. Done when:
   1. e2e test: a second player sees the fight mark.

### [ ] 15.32 Aiming arrow

1. Goal: a clear Hearthstone-style arrow while aiming (T-02, T-18 – T-20).
2. Tasks:
   1. A chain of thick red blocks with a big head, source to pointer (T-02, T-18).
   2. A little see-through; drawn over every card (T-19).
   3. Off target: plain head; on a valid target: a red ring (T-03, T-20).
   4. Only `transform` and `opacity`; follows animation speed (ST-05).
   5. Without images the arrow looks the same.
3. Done when:
   1. `make screenshots` shows the arrow off and on a target.

### [ ] 15.33 Start, finish, and cancel aiming

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

### [ ] 15.34 Plays without aiming

1. Goal: cards that need no target, or place then aim, feel natural.
2. Tasks:
   1. No target needed: drag onto the board; no arrow (T-12).
   2. Place first, then aim; cancelling the aim undoes the play (T-13).
   3. No valid target: follow what the engine allows (T-14).
3. Done when:
   1. e2e test: a no-target loot card dragged onto the board plays.
   2. e2e test: a cancelled aim leaves the card in hand.

### [ ] 15.35 Presence messages

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

### [ ] 15.36 Others' aim and hover

1. Goal: everyone sees what the active player is aiming at (T-08 – T-10, T-21, T-22).
2. Tasks:
   1. Others see the arrow live while it moves (T-08).
   2. A hidden card shows as a back lifted from its hand slot (T-09, T-22).
   3. The aimed-at card glows red on the target's screen (T-21).
   4. A hovered card is marked red for the others (T-10).
   5. A late or missing message never breaks the table.
3. Done when:
   1. e2e test: a second player sees the arrow and the red mark.

### [ ] 15.37 Hit and death marks

1. Goal: damage and deaths read at a glance (T-15, T-24 – T-26).
2. Tasks:
   1. After the target locks, the effect flies from source to target (T-15).
   2. A yellow starburst "-N" on every damaged card at once (T-24).
   3. An HP number turns red while below full (T-25).
   4. A dying card gets a big skull before it leaves (T-26).
   5. Only `transform` and `opacity`; follows animation speed (ST-05).
3. Done when:
   1. `make screenshots` shows a starburst, red HP, and a skull.

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
