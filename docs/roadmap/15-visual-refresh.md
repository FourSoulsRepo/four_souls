# Step 15. Visual refresh

Goal: a clearer table, using fouronlinesouls.net as a second reference.

Ideas: V-09, V-12, V-13, V-15 – V-29.

Order: 15.1, 15.9 – 15.11, 15.2, 15.12 – 15.14, then 15.3 – 15.8.
Images come first, so the new table is built with backs and mats.

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

### [ ] 15.2 Rework the table

1. Goal: apply the chosen ideas.
2. Tasks:
   1. Change layout and panels for 2, 3, 4 players.
   2. Keep V-01 – V-08 and touch support (B-06).
   3. Update `docs/wiki/Client.md`.
   4. Simple but informative: nothing needed is hidden (V-18).
   5. Every action gets instant feedback (V-15).
   6. Empty zones keep their place and size (V-23).
   7. One framed board: decks, shop, monsters, souls, stack (F-01, H-01).
   8. Nothing floats over the board; fixes O-01 and O-02 (F-02).
   9. Drop the green felt and seat boxes (H-05, O-03).
   10. Own character centered at the bottom edge, items beside it (F-06, H-02).
   11. Own hand large and fanned at the bottom, partly off screen (F-05, H-03).
   12. Opponent: name and coins, then character and items (F-07).
   13. Opponent hand: fanned backs with a count (F-08, H-17).
   14. 3–4 players: opponents in side columns; the board grows (F-09).
   15. The active player's name is highlighted (F-10).
   16. Player names small, outside the board (H-21).
   17. Decks show their backs at the board's edge (H-16).
   18. Coins big next to own name, like mana (H-18).
   19. HP and ATK as big numbers on the character card (H-15).
   20. The board has one fixed shape and size, empty or full (H-06).
3. Done when:
   1. `make frontend e2e` passes.
   2. `make screenshots` shows every player count.
   3. e2e test: board size stays the same when a shop slot empties.
   4. e2e test: no panel covers a deck with 2, 3, or 4 players.

### [ ] 15.3 Attack animation

1. Goal: attacks look like Hearthstone (V-12).
2. Tasks:
   1. The attacker's card jumps at the monster, hits, and returns.
   2. Damage shows on the card that was hit.
   3. Only `transform` and `opacity`; follows animation speed (ST-05).
   4. A starburst with "-N" on the card that was hit (H-14).
   5. HP on the card drops as soon as the damage resolves (C-03).
3. Done when:
   1. e2e screenshot of an attack in progress.

### [ ] 15.4 Hearthstone-style card play

1. Goal: cards are played by dragging them to a place or target (V-16).
2. Tasks:
   1. Study Hearthstone screenshots: dragging, aiming, highlights.
   2. Drag a card onto the table or onto its target.
   3. Aiming arrow; valid targets glow, others stay plain.
   4. Click card, then click target, still works.
   5. Works with mouse and touch (B-06).
   6. One-click mode: a click plays the card at once (V-28).
   7. A setting chooses one-click or drag; both send the same intent.
   8. Hover a hand card: it rises large, keyword tips beside it (H-13).
   9. Playable cards glow green; the rest stay plain (H-12).
   10. Aiming arrow: no screenshot found yet; study gameplay videos (H-23).
   11. Zoom opens over the hovered card, not at a screen side (F-16).
3. Done when:
   1. e2e test: a bomb dragged onto a monster damages it.
   2. e2e test: the same play by two clicks.
   3. e2e test: the same play in one-click mode.

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

### [ ] 15.7 Never leave the player lost

1. Goal: the player always knows when the game waits for them (V-20).
2. Tasks:
   1. List every prompt the engine can send: priority, targets, choices, pay, roll.
   2. Each prompt shows a clear banner: what to do, in a few words.
   3. Cards and buttons that answer it glow; the rest dim.
   4. The stack panel shows the same prompt (V-11).
   5. Waiting on others names who and what for.
   6. Optional: a soft sound or window flash when input is needed.
   7. Hover an unusable card: the zoom says why (F-17, C-06).
   8. A refused action shows its reason next to the clicked item (F-18).
   9. Choices: big cards over a dimmed board (F-20, H-11).
   10. A "View board" button hides the choice to look at the table (F-20).
3. Done when:
   1. e2e test: every prompt type shows its banner and highlights.
   2. `make screenshots` covers "your turn to act" and "waiting for Ann".

### [ ] 15.8 Show what the active player is doing

1. Goal: every player knows what is happening right now (V-21).
2. Tasks:
   1. List the activities to show: fight, shop, play a card, pay, think.
   2. Derive each one from game state on the client.
   3. The active seat shows a short status: "Ann attacks Pooter".
   4. In a fight, the player and the monster get a battle mark.
   5. Attack rolls, targets, and damage are visible to all.
   6. "Thinking" shows while the active player holds priority.
   7. Spectators see the same marks.
   8. Share the banner and highlight code with 15.7.
   9. In a fight, the attacker and the monster get a red frame (F-19).
   10. A card another player plays shows large for a moment (H-10).
3. Done when:
   1. e2e test: a second player sees the fight mark and the status line.
   2. `make screenshots` covers "Ann attacks" and "Ann is thinking".

### [ ] 15.9 Mats and dice in WebP

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

### [ ] 15.10 Mat backgrounds

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

### [ ] 15.11 Card backs

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

### [ ] 15.12 Turn action buttons

1. Goal: turn actions as clear icons next to own name (F-11).
2. Tasks:
   1. Buy and attack as icon buttons by own name.
   2. Lit when usable, dim when not (F-12).
   3. Buy turns into Cancel while choosing what to buy (F-13).
   4. Attack turns into Roll once 15.5 adds the throw (F-13).
   5. End turn: yellow on own turn, grey "Their turn" otherwise (H-19).
   6. End turn glows green when nothing is left to do (H-20).
   7. Keys: Space resolves or passes; E ends the turn (F-15).
   8. A tooltip names each button and its key.
   9. No number badges on cards (V-29).
   10. End turn: one big button at the board's right edge (H-04).
3. Done when:
   1. e2e test: each button is lit or dim as the engine allows.
   2. e2e test: Space and E send the same intents as the buttons.

### [ ] 15.13 Stack panel and history strip

1. Goal: stack and history read at a glance, with icons (S-01 – S-11).
2. Tasks:
   1. Stack panel inside the board, at its left side (F-03).
   2. Newest item on top; each row is an icon and one short line (S-01, S-02).
   3. Names in color: player in their color, target in red (S-03).
   4. Icon by kind: card art, die for a roll, mark for damage (S-04).
   5. Hover a row: card large, used ability highlighted, owner below (S-06, S-07).
   6. Resolve or Pass button in the panel; Space does the same (S-08).
   7. Empty stack: a quiet "Nothing on the stack" (S-09).
   8. History strip next to the stack: same icons, newest on top (S-10, F-04).
   9. History icon badge: sword for an attack, bolt for a spell (H-08).
   10. Hover a history icon: the same large card view (S-11, H-09).
   11. The full history panel (V-03) stays one click away.
   12. Keep the "waiting for" lines (V-11); the reference has none.
   13. The history strip runs along the board's left edge (H-07).
3. Done when:
   1. e2e test: hover a stack row shows its card and owner.
   2. `make screenshots` shows a stack of 3 items and the history strip.

### [ ] 15.14 Counters and card states

1. Goal: a card's state reads on the card itself (C-01 – C-06).
2. Tasks:
   1. Counter: a big round badge with a number on the card art (C-01).
   2. The zoomed card shows the same badge (C-02).
   3. Live HP, ATK, and dice values on the card (C-03).
   4. **Owner** decides: values drawn over the image, or a badge beside it.
   5. Text cards without images get the same badges.
   6. Used items stay dimmed, not rotated; already so (C-04).
   7. Two abilities: a dialog with the big card; pick one line (C-05).
   8. Status marks on cards, e.g. not usable yet (H-22).
3. Done when:
   1. e2e test: a counter badge on the card and in its zoom.
   2. e2e test: a damaged monster shows its new HP on the card.
