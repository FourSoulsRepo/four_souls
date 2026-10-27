# Ideas

Living list of ideas for the game.
Grouped by topic.
Each idea has a stable ID; new ideas are appended to the end of their group.

---

## Principles

1. PR-01. KISS: keep it simple.
   1. Each feature does one task and does it well.
   2. Prefer the simplest design that works.
2. PR-02. Big self-contained modules go to `/pkg`.
   1. Each is its own Go module, ready to move to a separate repo.
   2. Candidates are proposed to the owner before moving.
3. PR-03. Strict linters with security rules.
   1. Go: golangci-lint, strict, incl. `gosec`.
   2. Frontend: ESLint, strict type-checked, with security plugins.
   3. Both run locally and in CI; CI fails on any finding.

## Visual

1. V-01. Visual reference: Hearthstone.
   1. Cards are plain and flat, not 3D-ish.
   2. Animations are welcome; how to make them is still open.
   3. Hero portrait with health: optional, only if it looks good.
2. V-02. Opponent placement depends on player count.
   1. 2 players: opponent at the top.
   2. 3–4 players: opponents on the left, top, and right.
   3. Own area is always at the bottom.
3. V-03. Match history panel, like Hearthstone.
   1. Shows what happened, so nobody gets lost in the game.
   2. Available to players, judges, and viewers.
   3. Placement decided during UI design; left side is the first try.
   4. A button hides and shows the panel.
   5. Each entry: small icon; full details on hover.
   6. Entries show only what that person may see.
   7. Whole match history, with scrolling.
   8. One resolved stack is one entry at first; may change if it looks bad.
4. V-04. Separate plain-text game log for judges.
   1. Judges can copy text from it.
5. V-05. Turn indicator built into our UI.
   1. Shows whose turn it is and the turn number.
6. V-06. Player counters in our own style.
   1. Hand size, attack, dice modifiers, shop price, etc.
   2. Values are computed by the engine.
7. V-07. "Look at cards" effects show the cards Hearthstone-style.
   1. No extra buttons like in Tabletop Simulator.
8. V-08. Animations: cute, but not overcomplicated.
   1. Drag and drop for playing cards.
   2. Card moves, flips, zoom, damage, dice.
   3. No particle effects.
9. V-09. Second visual reference: https://fouronlinesouls.net/
   1. The table shows only after starting a game against a bot.
   2. Borrow layout and readability ideas only.
   3. Never copy its images, styles, or code.
10. V-10. The app does not behave like a web page.
   1. Screens never scroll as a whole page.
   2. Only lists with their own scroll bar may scroll.
   3. No text or element selection by mouse drag.
   4. Text stays selectable where copying helps (logs, payloads).
11. V-11. The stack shows what the game is waiting for.
   1. Lists the players who must still respond or pass.
   2. Shows who already passed.
   3. Shows a pending prompt: who must choose and what.
12. V-12. Attack animation, like Hearthstone.
   1. The attacker's card jumps at the monster and hits it.
   2. Then it returns to its place.
   3. Damage shows on the card that was hit.
13. V-13. Check the Hearthstone table again.
   1. Today our table looks more like a poker table.
   2. Compare layout, board shape, and card placement.
14. V-14. History shows which card was stolen.
   1. The thief and the victim see the card.
   2. Everyone else sees only "A stole a card from B".
15. V-15. The UI feels responsive.
   1. Every click, hover, and drag gets instant visual feedback.
   2. Feedback shows before the server answers.
   3. The game state still changes only on server events.
   4. A slow answer shows a small busy mark, not a frozen screen.
16. V-16. Play cards like in Hearthstone.
   1. Drag a trinket or item from hand onto the table.
   2. Drag a card onto its target, e.g. a bomb onto a monster.
   3. An arrow follows the pointer while aiming; valid targets glow.
   4. Simple way still works: click the card, then click the target.
   5. Both ways send the same intent.
   6. Study Hearthstone screenshots first; maybe other card games too.
17. V-17. Dice look and feel like real dice.
   1. A die falls on the table, bounces, and rolls to a stop.
   2. The value is still decided by the engine, as now.
   3. The animation only shows that value; it never picks one.
   4. The player presses a "Throw" button; no automatic roll.
   5. Open: what happens if the player does not press it (timer, bot).
   6. Open: whether the button needs a new intent in the protocol.
   7. The player always knows the game waits for their throw.
   8. The Throw button is big, glows, and sits where the die will land.
   9. The stack panel says "Waiting for you to throw" (V-11).
   10. Others see "Ann is throwing…" and a die shaking in her hand.
   11. Proposal for 5: no timer by default, as at a real table.
   12. The response timer (UC-06) throws for a player who stays idle.
   13. Bots and disconnected players throw at once.
   14. Proposal for 6: yes, a "throw" intent from the roller.
   15. The engine draws the value only when the intent arrives.
   16. So nobody, the roller included, can see it early.
   17. The record keeps the intent, so replays stay exact.
   18. Items 5–17 are open; the owner decides them in roadmap step 15.5.
   19. Decided (owner, 2026-10-15): proposals 11–17 accepted (15.5).
18. V-18. Simple but informative interface.
   1. Do not overload the screen with numbers and panels.
   2. Do not hide what players need to know.
   3. Key facts are always visible: health, coins, souls, turn, stack.
   4. Details show on hover or zoom, not as extra text on the table.
19. V-19. Players choose the look of their dice.
   1. Dice textures live with the private card images.
   2. Tabletop Simulator layout: 3×3 faces, top row unused.
   3. Row 2: faces 2, 4, 5; row 3: faces 1, 3, 6.
   4. The setting keeps the chosen look; the first one is the default.
   5. No texture: a plain die with the number, as in a public clone.
   6. First version: one look ("Tear"), no picker yet.
   7. Later: a picker in settings once there are more looks.
   8. Later: others see each player's own dice; needs a protocol field.
   9. Everyone sees the roller's dice, never their own look.
   10. The look a player picks is what all players see on that player's rolls.
   11. Each player's look reaches the others with their seat in the game.
   12. Today (8.15) every roll shows the viewer's own look; step 15.6 changes that.
   13. A player without that texture sees a plain die.
20. V-20. The player never feels lost.
   1. When the game waits for this player, it says so clearly.
   2. A clear banner names what to do: "Choose a target", "Pay 3¢".
   3. The cards or buttons that answer it glow; others dim.
   4. The stack panel shows the same prompt (V-11).
   5. When the game waits for others, it names who and for what.
   6. No silent waits: a screen never looks idle while input is needed.
   7. Covers every prompt: priority, targets, choices, dice (V-17).
21. V-21. Everyone sees what the active player is doing.
   1. V-20 covers prompts to the viewer; this covers watching others.
   2. Every player can tell the current activity at a glance.
   3. Fight: the player and the monster are marked as in battle.
   4. Others see the attack roll, the target, and the damage.
   5. Thinking: "Ann is choosing what to do" when she holds priority.
   6. Also: buying from the shop, playing a card, paying a cost.
   7. One status line or badge on the active seat says it in a few words.
   8. Same marks on every screen, spectators included.
   9. Built from game state; no new intents needed.
22. V-22. Compare with real Hearthstone screenshots.
   1. Find Hearthstone screenshots on the internet.
   2. Keep them in `.cache/screenshots/`, never in the repo.
   3. Compare board shape, hand, hero, and buttons with our table.
23. V-23. Table areas keep their size.
   1. An empty zone keeps its place and size.
   2. Nothing on the table grows or shrinks when a zone empties.
   3. As on fouronlinesouls.net today.
24. V-24. Game mat image as the table background.
   1. First mat: `pkg/card_db/images/assets/mats/basement.jpg`.
   2. Every image in `assets/mats/` becomes a choice in the mat setting (ST-03).
   3. A new mat needs only a new file in that folder.
   4. A build without images uses plain colors, as now.
25. V-25. Mats and dice in WebP, like the cards.
   1. Mats are JPG today; dice are JPG and PNG.
   2. Convert them to WebP for size.
   3. The conversion lives in `FourSoulsRepo/tools`, not here.
   4. The result amends the images repo; it stays one commit.
26. V-26. Card backs for decks and opponent hands.
   1. Each deck shows its own back: loot, treasure, monster.
   2. An opponent's hand shows backs, not blank boxes.
   3. No back images exist yet; the tools repo must collect them.
   4. Without images: plain colored backs, as now.
   5. Owner: backs come from the official foursouls.com site.
   6. The tools repo downloads them, like the card images.
27. V-27. Study more of fouronlinesouls.net.
   1. Counters: how they look and change.
   2. Stack window with card icons; it looks good.
   3. Stack history strip with icons; it looks good.
   4. Their room has "Allow cheat options" (on by default).
   5. Cheat menu sits in the game menu ("Cheat view").
   6. It gives any card, which makes testing a card quick.
   7. Later: study it as a model for our sandbox (LR-01).
28. V-28. Two ways to play a card; the player picks.
   1. One click: the card plays at once, as on fouronlinesouls.net.
   2. Drag and drop: as in Hearthstone (V-16).
   3. A setting chooses the mode.
   4. Both modes send the same intent.
   5. Dropped (2026-10-15): one click, then a click on the target (15.26).
29. V-29. No number badges on cards.
   1. Owner rejected 1–9 badges on valid targets (15.1, F-14).
   2. Targets glow instead (V-20).
30. V-30. Souls look like mana crystals in Hearthstone.
   1. A row of slots; filled ones are gained souls.
   2. Shows at a glance how close a player is to winning.
31. V-31. Gained souls sit left of the character, not right.
32. V-32. Switch character or eternal item art between reprints.
   1. Right click, or long tap on touch, cycles the art.
   2. Only cards that have reprints in other sets.
   3. Open question: only reprints from enabled sets, or always?
   4. It is only an image, so "always" may be fine.
   5. Related: alt art (CD-11), per-copy art (CD-12).
33. V-33. Title screen uses the Isaac fonts.
   1. The fonts in `assets/fonts` of the private images repo.
   2. Without images: the plain UI font, as now.
34. V-34. Top image of foursouls.ru on the home screen.
   1. The banner at the top of the foursouls.ru main page.
   2. The tools repo downloads it; it lives with the private images.
   3. Without images: the home screen as now.
35. V-35. Background animation of foursouls.ru.
   1. On the home screen and the game setup screens.
   2. Same look as the site; our own code.
   3. Without images: a plain background.
36. V-36. Declarations show on the stack.
   1. "Buy" and "Attack" show as stack entries.
   2. Other declarations too, wherever the rules use the stack.
   3. Each entry names the player and the target.
37. V-37. After the game, everyone may look at the table.
   1. The winner is announced; the table stays open.
   2. All hands are revealed to everyone.
   3. Players can discuss how the game was played.
   4. The end overlay can be closed and opened again.
38. V-38. Failed effects show in the history.
   1. A fizzled ability gets its own history entry.
   2. Example: a purchase fails because the player cannot pay.
   3. The entry says what failed and why.
39. V-39. Stack and history rows move with animation.
   1. A new row fades in.
   2. Other rows slide to their new place, not jump.
   3. A leaving row fades out before it is gone.
   4. Same for the stack panel and the history strip.
40. V-40. Played cards fly to the stack and on to the discard.
   1. The card flies from the hand to the stack.
   2. There it turns into its stack row.
   3. Shrink, or another effect; the owner picks.
   4. Resolved, the card flies from the stack to the discard pile.
   5. Damage and death marks play first; then the card flies.
41. V-41. Stat boosts show what gives them.
   1. A character with extra HP or ATK shows small icons on its card.
   2. One icon per card or effect that changes the stat.
   3. Hover an icon to see its source card.
   4. Owner (2026-10-22): icons only for temporary effects (15.54).
   5. Lasting sources show in the badge's hover sum.
   6. A fight shows what the roll needs (15.94).
42. V-42. Curses sit left of the player.
   1. Not on the right with the items.
   2. Each player's curses form their own group.
43. V-43. Soul cards hide behind the soul crystals.
   1. Gained soul cards are not on the table.
   2. Hovering a filled crystal zooms its soul card.
   3. A killed monster shows with its soul.
44. V-44. History shows what a player bought.
   1. A purchase gets its own history entry.
   2. The entry names the buyer and the item.
   3. Hovering the entry shows the item card.
45. V-45. "Choose one-" cards ask like The Bone.
   1. The mode dialog shows the card large, as the ability dialog does.
   2. One parchment line per mode, in the card's own words.
   3. Same look as The Bone's ability lines (15.24).
46. V-46. Ask before ending the turn while plays are left (owner, 2026-10-20).
   1. Stops an End turn pressed by mistake.
   2. Asks whenever End turn is not glowing (H-20).
   3. The question lists what is left: attack, buy, loot, abilities.
   4. A setting turns it off; on by default.
47. V-47. Each animation states its reason and direction (owner, 2026-10-20).
   1. The animation list says what each one tells the player.
   2. It says where the motion goes: from where to where.
   3. The owner checks every animation against that idea.
48. V-48. A die at rest is seen from above (owner, 2026-10-20).
   1. The value is on the top face, as on a real table.
   2. The top face fills most of the die.
   3. Two sides still show, so the die looks solid.
49. V-49. The turn counters say "Shop", not "Buys" (owner, 2026-10-20).
   1. The line reads "Loot · Attacks · Shop".
   2. It reads better next to the other counters.
50. V-50. Texts say "dice", never "die", for a roll (owner, 2026-10-20).
   1. "Die" also means a character's death.
   2. Players could not tell which one the game meant.
   3. Card texts keep "die" for death.
51. V-51. Stack and history are the players' source of truth (owner, 2026-10-21).
   1. Players can read every step before: who did what.
   2. Who played which card, bought which item, killed which monster.
   3. The stack is now; the history is what happened.
   4. Simple, with no repeats between them.
   5. History: flat lines under turn headers.
   6. History strip: headlines only, Hearthstone-style tiles.
   7. Hover a tile: the card, an arrow, and what it did.
   8. An open "Attack" or "Buy": a top row in the stack panel.
   9. One design for V-36, V-38, and V-44.
52. V-52. The table shows who holds priority (owner, 2026-10-22).
   1. Players still miss that the game waits for them.
   2. The banner sits at the edge; eyes stay on the middle.
   3. Holding priority: the center board glows green, any seat.
   4. It pulses once when priority comes, then glows softly.
   5. No yellow "waiting" glow: green and yellow look alike to many.
   6. Instead the asked player's seat is marked, for everyone.
   7. Not holding priority: the board stays plain.

## Rules

1. R-01. No single source of truth for the rules.
   1. The rules are complicated and spread across sources.
   2. We clarify them step by step during development.
2. R-02. Rule sources:
   1. Official site: https://foursouls.com/rules/
   2. Russian translation: https://foursouls.ru/rules
   3. Reddit: https://www.reddit.com/r/FourSouls/
   4. Fan Discord server "Турнирный сервер".
   5. Twitter threads by people who know the rules best:
      1. Charlie "Yuggy", foursouls.com Rules Tzar: https://twitter.com/YuggyHD
      2. Charlie on Reddit: https://www.reddit.com/user/Yuggy
      3. Sean "Kizzycocoa", foursouls.com developer: https://twitter.com/Kizzycocoa
      4. Edmund McMillen, designer: https://x.com/edmundmcmillen
   6. BoardGameGeek rules forum:
      https://boardgamegeek.com/boardgame/255664/the-binding-of-isaac-four-souls/forums/66
3. R-03. Disputed cases are resolved one by one.
   1. We discuss each case together.
   2. The decision becomes a rules engine test.
   3. The test links to the source of the decision.
   4. The set of resolved cases forms our final rules.
   5. Each test is a situation payload (LR-04) with `expect` filled in.
4. R-04. Our rules text lives inside the rules engine.
   1. Each rule has a stable ID, e.g. `R-COMBAT-03`.
   2. A text copy ships in the client app and the website.
   3. Links point only to our local rules, never to external sites.
5. R-05. First version supports only the "Base Game" set.
   1. The engine design keeps all later sets in mind.
   2. Their mechanics must fit without a full engine rework.
6. R-06. Mulligan: check in the rules whether it exists.
7. R-07. Known rulings by Edmund (from the TTS "Ed's Notes" notebook).
   1. Monster Manual can force a second attack.
   2. You cannot kill dead players.
   3. Copying an item is the same as getting it new.
      Items 1–3: https://twitter.com/edmundmcmillen/status/1069419513421000704
   4. You can deal damage to dead players.
      https://twitter.com/edmundmcmillen/status/1069125603494813696
   5. Eden can have Glass Cannon as Eternal.
      https://twitter.com/edmundmcmillen/status/1069276355932508162
   6. Steam Sale only affects the shop item, not the top deck.
      https://twitter.com/edmundmcmillen/status/1070040592258756608
   7. Start phase: items recharge before the loot draw; you can respond to the draw.
      https://twitter.com/edmundmcmillen/status/1074370086964584448
   8. "Prevent damage" effects last until end of turn.
      https://twitter.com/edmundmcmillen/status/1076372423040106496
   9. You can add cards to a stack while it resolves.
      https://twitter.com/edmundmcmillen/status/1068935710252576768
   10. Death cancels combat.
      https://twitter.com/edmundmcmillen/status/1082176191283396609
   11. Bombs are not combat.
      https://twitter.com/edmundmcmillen/status/1070033655412645888
   12. Old rulings: check each against the current official rules.
   13. Each confirmed ruling becomes a test case (R-03).

## Architecture

1. A-01. Rules engine is a separate module.
   1. It has its own version.
   2. It lives in `/pkg` for now, to keep things simple.
   3. Separate Go module with its own `go.mod`.
   4. Its own `docs/` folder.
   5. Treated as a separate repo, ready to move out later.
   6. Folder: `pkg/rules_engine`.
   7. Imported in code under the alias `engine`.
   8. Module path prefix: `github.com/FourSoulsRepo` (may change later).
2. A-02. The app shows two versions:
   1. App version.
   2. Rules engine version.
   3. Both are shown in a corner of the main menu.
3. A-03. The rules engine may later power a website.
   1. Compiled to WebAssembly.
   2. Users check rules and game situations there.
   3. A simple sandbox, not a full table view.
   4. User adds cards and effects.
   5. Site says whether the action is possible, and why not.
   6. A situation is shareable as a link.
   7. The whole situation is encoded in the URL query.
   8. A pasted payload (LR-04) is accepted too.
   9. Both forms, so anyone can host the same sandbox elsewhere.
   10. Every answer explains why, not just yes or no.
   11. The answer links to the rule it is based on, when possible.
4. A-04. Several entry points in one repo, one shared codebase.
   1. `main.go` in the root: Wails client app.
   2. `cmd/server/main.go`: dedicated server, no UI.
   3. `cmd/website/main.go`: future WebAssembly site.
   4. The site has no web server; rules run in the browser.
   5. Hosted on GitHub Pages; no server, not downloadable.
   6. `cmd/relay/main.go`: relay server for viewers (SP-04, future).
5. A-05. General rules are separate from card logic (MTG Arena style).
   1. Engine knows only general rules: phases, priority, stack, death.
   2. Pending actions are written down first ("deal 1 damage to X").
   3. Card effects may rewrite them: prevent, replace, cancel.
   4. Then the engine applies the final list.
6. A-06. Cards are data built from reusable effect blocks.
   1. Blocks like: loot X, gain X¢, deal X damage, roll and choose.
   2. Only truly unique cards get their own code.
7. A-07. One view-filter function decides who sees what.
   1. Builds the view for a player, a viewer, or a judge.
   2. Nothing else sends state to clients.
8. A-08. Engine is deterministic.
   1. The RNG seed never leaves the server.
   2. Never decide anything by Go map iteration order.
   3. A test runs the same game twice and compares results.
9. A-09. CI builds the engine for WebAssembly from day one.
   1. Catches dependencies that break the website build.
10. A-10. Stay on Wails v2.
   1. Wails v3 is in beta (Aug 2026); v2 is the stable line.
11. A-11. `pkg/card_db`: card display data, a separate Go module.
   1. Text, artists, translators, translations, versions, images.
   2. Engine needs only card IDs and effects, not this data.
   3. Card images live in a separate private repository.
   4. Private because the images are licensed.
   5. CI reads the private repo with a deploy key or token secret.
   6. The main repo is public; the image repo protects source files only.
   7. Without image access the app still builds and runs.
   8. Missing images fall back to blank cards with text (like SP-05).
   9. Plain git, no Git LFS (LFS has a monthly download quota).
   10. Images are WebP, small but with readable card text.
   11. One stored size, readable when a card is zoomed in.
   12. The image repo always has a single commit.
   13. New images are added by amending that commit; old images are never updated.
   14. CI clones with `--depth 1`.
   15. Fallback if it grows too big: image packs as GitHub Release assets.
   16. One image repo, one folder per set.
   17. Split into a repo per set only when we hit the limits.
   18. The image repo is `FourSoulsRepo/card_db` (private).
   19. Mounted as a git submodule at `pkg/card_db/images`.
   20. Release builds embed images with `-tags cardimages`.
12. A-12. `pkg/record`: match record format, a separate Go module.
   1. Stays in this repo for now; not moved out.
13. A-13. Network protocol stays inside the app.
   1. Location: `internal/protocol`.
   2. App-specific; never split out.
14. A-14. Frontend: React + TypeScript, table drawn with the DOM.
   1. No canvas, no particles; keeps the app fast.
   2. Animate only `transform` and `opacity`.
   3. Frontend talks to Wails through a thin bridge layer.
   4. Eases a later move to Wails v3 (needed for mobile).
   5. Decision: ADR 002.
15. A-15. Thin web client.
   1. The same React frontend, built for a normal browser.
   2. Joins games only; cannot host.
   3. Useful for mobile without a native app.
   4. Served by the game server itself (dedicated or hosting app).
   5. Example: open `http://192.168.1.10:7777` on a phone in the same Wi-Fi.
   6. Same origin: plain HTTP works, no certificate needed.
   7. The page always matches the server version; never outdated.
   8. Card images: served only if an image folder sits next to the server.
   9. Otherwise blank cards with text (like SP-05).
   10. Used for playing and watching; not for replays.
   11. Later maybe: a GitHub Pages copy for public servers with TLS.
16. A-16. No leaks through card lists.
   1. Lists sent to a client hold only cards that are public.
   2. Public: on the table, played, discarded, revealed, or seen by that client.
   3. Never cards still hidden in a deck.
   4. The collection check compares whole sets, not the dealt decks.
   5. With deck ratios (GS-07) the decks are a random subset.
   6. Sending that subset would tell players what is in the decks.
   7. Exception: cards the host banned are public and may be left out.

## Network

1. N-01. Two ways to run a server:
   1. Dedicated server: headless Go console app.
   2. Hosted: one player's app runs the server; others join.
2. N-02. Server is authoritative.
   1. Only the server applies the rules.
   2. Clients send intents ("I want to do X").
   3. Server applies the intent or rejects it with a reason.
   4. This leaves no room for cheating.
3. N-03. Clients do not depend on the rules engine version.
   1. Only the server checks the rules.
   2. A client may connect with a different engine version.
4. N-04. Server sends the allowed actions to each client.
   1. Sent after every change of the game state.
   2. The UI shows what the player can do right now.
   3. A client may re-request them at any time.
   4. So a lost message never leaves the client stuck.
5. N-05. Responses to the stack by any player.
   1. Not only the current player may respond.
   2. A player with no possible response is skipped automatically.
   3. The server does not wait for that player.
   4. This leaks nothing: charges and empty hands are public.
   5. A player with a charged character or item is never auto-skipped.
6. N-06. "Skip all" button.
   1. The player passes on everything automatically.
   2. Control returns only when a new card or event hits the stack.
   3. It skips only the events visible when the button was pressed.
   4. Lasts until the stack is empty.
   5. The player can cancel it at any time.
   6. Removed (owner, 2026-10-20): play tests showed it unused.
7. N-07. Response timer.
   1. Host sets it when creating a game.
   2. Can be disabled.
   3. Minimum value: 60 seconds for now, so matches stay playable.
   4. When it runs out, the server passes for the player.
   5. Includes a fixed allowance for client animations.
8. N-08. Reconnect after a lost connection.
   1. The player rejoins the same game.
   2. The server sends the full current state.
   3. While a player is away, the game pauses for everyone.
   4. Other players vote: wait longer or kick.
   5. A kicked player's seat stays empty; never a bot.
9. N-09. Auto-skip uses the full list of allowed actions.
   1. Includes item and character abilities, not only cards in hand.
   2. Lesson from MTG Arena players missing response windows.
10. N-10. Phase stops (to consider).
   1. Player marks phases where the game always stops for them.
   2. Example: always stop at the opponent's end phase.
11. N-11. Save a game and continue it later with the same friends.
   1. The host saves an unfinished game at any moment.
   2. The save file stays with the host or server.
   3. Players get no copy: it holds every hand (RP-05).
   4. Loading opens a lobby with the saved seats.
   5. Each player takes back their own seat.
   6. The game resumes when every seat is taken.
   7. It continues exactly as saved: stack, prompts, dice state.
   8. A save made by an older game version may not load.
   9. Builds on the engine save (step 4) and records (RP-08).
12. N-12. Rollback of an event.
   1. Undoes a mistake, e.g. a misclick.
   2. Soft: the player who did the event cannot roll it back.
   3. Soft: only the other players can roll it back.
   4. Hidden info: the rollback mode decides (item 15).
   5. 3–4 player games: a majority of the other players (item 16).
   6. The rollback is shown in the history and kept in the record.
   7. Decided: after a rollback, the next dice roll gives the same value.
   8. Same value: the RNG is restored, so the player knows the next roll.
   9. That holds whatever is played first, since it is one RNG stream.
   10. New value: a rollback becomes a free re-roll.
   11. Rejected as the only rule: a seen dice roll blocks rollback.
   12. Accepted risk: the player knows the next roll after a rollback.
   13. Only other players can roll back, which limits the abuse.
   14. The author may ask for a rollback; the others decide (owner, 2026-10-26).
   15. Three rollback modes (owner, 2026-10-26; raw idea).
      1. Disabled: no rollback at all.
      2. Soft: steps go back until a barrier event.
      3. Soft barriers: a dice roll, damage dealt, hidden info seen.
      4. Soft barrier too: the turn ends (owner, 2026-10-26).
      5. Godmode: every step can be rolled back.
      6. Godmode keeps item 7: the restored RNG repeats the roll.
      7. A death follows damage; a revealed card is hidden info seen.
      8. Not a barrier: a shuffle; the restored RNG repeats its order.
      9. Godmode: no vote; the rollback happens at once (owner, 2026-10-26).
   16. Majority: more than half of the other connected players (owner, 2026-10-26).
      1. 2 players: the opponent.
      2. 3 players: both others.
      3. 4 players: two of the three others.
   17. The host picks the mode in Detailed setup (owner, 2026-10-26).
   18. The default mode is Disabled (owner, 2026-10-26).
   19. Godmode: anyone may roll back, the author too (owner, 2026-10-26).
   20. Rejected: the RNG advances on every step (owner, 2026-10-26).
      1. The step count would pick the next roll.
      2. Same steps bring back a good roll; one more pass re-rolls a bad one.
      3. Every engine change that adds a step would shift every roll.
   21. The record keeps the undone steps and a rollback line (owner, 2026-10-26).
      1. Evidence for bug reports stays in the record.
      2. The record stays append-only; nothing is rewritten.
      3. Replays show what really happened.
   22. Old states come from a replay, not from a stored state per step.
      1. The server keeps one base state: the game start or the last restart.
      2. Rollback to step k: load the base, replay its intents up to k.
   23. A server restart keeps the rollback history if it checks out (owner, 2026-10-26).
      1. The restart file also keeps the old base state.
      2. On load, the new engine replays from the old base to the end.
      3. Same bytes as the kept state: rollback reaches past the restart.
      4. Different: the restart becomes a barrier; the history shows it.
   24. Soft: a played card can be rolled back (owner, 2026-10-26).
      1. Played card: a loot card or an activated item.
      2. Not a purchase; it keeps the barriers it has.
      3. Allowed while no new card came from the treasure, monster or loot deck.
      4. And no dice was rolled or changed, and no damage was dealt.
      5. A turn end, a death, and a seen hand card still block it (owner, 2026-10-26).
      6. Mostly true already: playing or activating is no barrier.
      7. Change: a card put privately on the loot deck no longer blocks.
   25. One Undo button instead of a button per history row (owner, 2026-10-26).
      1. It sits under the History title.
      2. One press undoes one move: back to before the latest target.
      3. Press again to go further back; Godmode has no jump to any row.
      4. Greyed at a barrier, during a vote, at game end; the tooltip says why.
      5. Hidden in Disabled mode.
      6. The tooltip names the move it undoes.
      7. Soft: the author's press asks; another player's opens the vote.
      8. The history stays clean: no play rows (17.6 rows removed).
13. N-13. Find LAN games without typing an IP.
   1. A hosting app announces its waiting game on the LAN.
   2. The join screen lists the games it finds.
   3. Typing an address still works.
   4. A started game keeps announcing, marked "In progress" (owner, 2026-10-22).
      1. Only a player who held a seat can rejoin it.
      2. So a crashed player finds the game in the list again.
   5. A dedicated server never announces.
   6. Players join it by address and pick from its list of games.
14. N-14. A game can have a password.
   1. Set by whoever creates the game.
   2. The games list shows a lock; joining asks for the password.
15. N-15. Live aiming and hover reach the other players.
   1. Others see the aiming arrow while it moves (15.15, T-08).
   2. Others see which card is hovered or aimed at (T-10, T-21).
   3. A hidden card shows only as a back from its hand slot (T-09, A-16).
   4. A separate kind of message: not an intent, not a game event.
   5. The engine never sees these messages.
   6. Never recorded, so replays never show them.
   7. Lost or late ones do no harm; the next one replaces them.
   8. The server limits how often a player may send them.
   9. Viewers see them too (owner, 2026-10-20).
   10. A waiting player sends only aims, no hover (owner, 2026-10-20).
16. N-16. "Wait for my turn" button (owner, 2026-10-13).
   1. Offered to a player on another player's turn.
   2. The server passes for them on everything.
   3. Choices and discards still ask them.
   4. It ends when the player before them ends the turn.
   5. That end window is their last chance to act.
   6. The player can cancel it at any time.
   7. It also ends when another player targets them (owner, 2026-10-20).
   8. Or targets their items, souls, curses, or hand cards.
   9. So they can answer before cards are stolen from them.
   10. Cards that pick the victim on resolution end it too (owner, 2026-10-20).
17. N-17. Create a game on a dedicated server (owner, 2026-10-14).
   1. A player connected to a dedicated server opens a new game there.
   2. The Games list shows "Create a game" only on a dedicated server.
   3. A game hosted in an app: only that app creates or loads games.
   4. Friends there join; they never open games on the host's app.
18. N-18. Server names and a Recent list the player edits (owner, 2026-10-20).
   1. A dedicated server's operator sets its name in the setup.
   2. The app shows that name instead of the bare address.
   3. The player removes a server from the Recent list.
   4. A match record names the server it ran on (owner, 2026-10-20).
19. N-19. Easy server setup: systemd, Docker, compose with Caddy (owner, 2026-10-21).
   1. The operator picks one of three ready setups.
   2. ENV variables configure the server, for Docker.
   3. One Dockerfile, downloaded alone, builds the latest release.
   4. Server archives have no version in their names.
   5. Desktop app archives keep their versions.
   6. Caddy gives HTTPS; the client speaks `wss://`.
   7. A bare domain tries `wss://` first, then plain `ws://` on 4774.
   8. Compose tracks the server with a healthcheck.
20. N-20. Seamless server restart and idle game limits (owner, 2026-10-21).
   1. A server update or restart does not end the games.
   2. Players reconnect by themselves and play on.
   3. Same engine version: the game goes on silently.
   4. Newer engine that still reads the game: a notice to the players.
   5. The notice says cards may now behave differently.
   6. A game the new server cannot read ends with an apology.
   7. "Can read": the state checksum after loading matches the saved one.
   8. No drain mode: games can last an hour, or never end.
   9. A game with no player connected for 30 minutes ends.
   10. A game ends 24 hours after it started, even if still played.

## Learning

1. LR-01. Situation sandbox inside the client app.
   1. Build a situation and check how the rules handle it.
   2. Uses the same table view as the game, not a web-style form.
   3. Every answer explains why, with a link to the rule when possible.
2. LR-02. Rules text available inside the client app.
3. LR-03. Tutorial stage in the client app.
   1. Teaches the rules in action, not by reading text.
   2. Lessons: short fixed scenarios on one topic each.
   3. Practice game against a bot.
4. LR-04. Situations can be exported as a shareable payload.
   1. Not a link: a payload pasted into the website or the app.
   2. Holds the situation setup only, no step history.
   3. Holds what the situation is meant to verify (which rule).
   4. The app offers both: copy payload and copy link.
   5. The link's base URL is added later.
   6. Payload contents: setup, one action, optional expected result.
   7. Expected result: action allowed or not, plus optional checks.
   8. A short free-text title states the question.
   9. The action itself is the question.
   10. Without `expect`, the engine just answers.
   11. With `expect`, the engine says "matches" or "doesn't match".
   12. Example (format is a draft, designed later):

```json
{
  "title": "Can a dead player respond to a dice roll?",
  "setup": {
    "players": [
      { "character": "isaac", "dead": true, "hand": ["a_penny"] },
      { "character": "maggy", "items": ["d6"] }
    ],
    "stack": [ { "dice_roll": 4, "owner": 1 } ]
  },
  "action": { "player": 0, "play": "a_penny" },
  "expect": { "allowed": false }
}
```
5. LR-05. Two practice game modes against a bot.
   1. Seeded: bot plays and player draws known cards.
   2. Seeded: only the planned moves are allowed.
   3. Free: like a real game, with hints from the app.
6. LR-06. The bot is simple.
   1. Plays random legal moves.
   2. Used only for practice, never in network games.
7. LR-07. Payloads with `expect` double as engine tests.
   1. Stored as files in the rules engine module.
   2. The test runner loads each file and checks `expect`.
   3. A case shared from the app or site becomes a test by saving it.
8. LR-08. FAQ built from frequent rules questions.
   1. Questions come from the collected sources (Reddit, Yuggy, Discord).
   2. The most frequent ones become FAQ entries.
   3. Each answer explains how it works and why.
   4. Each answer links to the rule IDs it rests on.
   5. Where possible, an entry has a situation payload to try it.
   6. Shown in the app's rules viewer and on the rules website.
   7. BoardGameGeek rules forum is a source too (R-02).
   8. Its answered questions also become test cases (R-03).
9. LR-09. Payload setup holds any event on the stack.
   1. Not only dice rolls: loot cards, abilities, attacks, triggers.
   2. Each event has its owner, targets, and choices made.
   3. The stack is ready to resolve, bottom first.
10. LR-10. Step-by-step stack resolution on the rules website.
   1. Not only a yes/no answer to one action.
   2. The user resolves the stack one step at a time.
   3. Each step shows what happened and why, with rule links.
   4. Players may respond between steps.
   5. A step back is possible.
   6. Goal: understand how the stack works.

## Testing

1. TS-01. Each card is written test-first.
2. TS-02. Card status report.
   1. Per card: not started, implemented, tested.
   2. Shows the share of done cards per set.
3. TS-03. Card stub generator.
   1. Reads scraped card data (CD-07).
   2. Creates a stub per card with its text and a TODO marker.
4. TS-04. Interaction tests for card combinations.
   1. Most bugs come from combinations, not single cards.
5. TS-05. Fuzz tests.
   1. Generate random situations from random cards.
   2. Check invariants: no crash, valid state, game always progresses.
   3. Same seed gives the same result.
   4. Views never leak hidden cards.
   5. Mode "all cards": fully random situations from the whole collection.
   6. Mode "focused": given card(s) always in play, others random.
   7. Focused mode is run for new or changed cards.
   8. Real fuzzing runs locally only, not in CI.
   9. CI only checks the fuzz tests still work (short smoke run).
6. TS-06. Every e2e test fails on a console error.
   1. Today the shared fixture checks only animation clashes.
   2. React warnings and page errors pass silently.
   3. `games.spec.ts` already fails on them; 19.13 was found that way.
   4. Move that check into `e2e/tests/fixtures.ts`.
   5. Harmless errors get a fix or a named allowance.
7. TS-07. The client reads every field the server sends.
   1. 19.29 and 19.30: the engine sent a field; the client never drew it.
   2. Engine tests pass; store-injected e2e tests cannot see it.
   3. Go test: every JSON field of the view types is in `protocol.ts`.
   4. Vitest: every view field in `protocol.ts` is read in `frontend/src`.
   5. A field left unread on purpose is listed with its reason.
   6. The golden file covers sample messages only, not empty fields.

## Localization

1. LO-01. First version: English only.
2. LO-02. Later: user picks the app language.
   1. Cards are shown in the chosen language.
3. LO-03. Later: view a card in other languages right in the game.
   1. How exactly: TBD.
   2. Proposal: hover or zoom a card, a hotkey cycles its languages.
4. LO-04. Everything is translatable.
   1. Menus and screens.
   2. Messages: errors, refusal reasons, history, prompts.
   3. Cards: name and text, from card_db translations (A-11).
   4. The server sends message keys and values, not ready text.
   5. So each client shows the text in its own language (N-03).
5. LO-05. Translating is easy for anyone (owner, 2026-10-26).
   1. One interface file per language.
   2. Each line has a note: where it shows.
   3. A check lists untranslated lines; another lists unused keys.
   4. A screenshot shows a line in place: one key or all of them.
   5. Wiki: one guide for the interface, one for cards.
   6. The same tools help an agent translate to a new language.
   7. Item 4 per key is paused: too costly for now (owner, 2026-10-27).
      1. Step 16 keeps one screenshot per screen.
      2. Per-key pages may come back later.
6. LO-06. English is always the fallback (owner, 2026-10-26).
   1. Interface lines, messages, card names, card images.
7. LO-07. Card images per language (owner, 2026-10-26).
   1. One folder per language: `images/en/b2/`, `images/ru/b2/`.
8. LO-08. Russian card names from the old site's migrations.
   1. The importer lives in `FourSoulsRepo/tools`.
   2. Card texts come later.
9. LO-09. Work-in-progress translations are welcome (owner, 2026-10-27).
   1. A card-only language is a choice in Settings.
   2. A half-done interface file ships; English fills the rest.
   3. Translators may keep card images in their own repo.

## Settings

1. ST-01. Fullscreen or windowed mode.
2. ST-02. Theme: dark or light.
   1. Changes menus and panels only for now.
   2. Replaced by ST-07 (owner, 2026-10-20).
3. ST-03. Game mat: choose one of the official mats.
   1. Mat images are collected by the scraping script (CD-07).
   2. Bad images are replaced by hand.
   3. Maybe reuse TTS table themes too (Basement, Cathedral).
4. ST-04. Keyboard shortcuts (future).
   1. Not in the first version.
   2. Configurable in the settings.
5. ST-05. Animation speed setting.
6. ST-06. Nickname.
   1. Each player sets a nickname; asked on the first start.
   2. Sent when joining a game; shown in the lobby, at the table and in the history.
   3. 1–20 characters; the server makes equal names in one game unique.
   4. Can be changed in the settings; a running game keeps the old one.
7. ST-07. One dark theme; no theme choice (owner, 2026-10-20).
   1. Replaces ST-02.
   2. The light theme and the Theme field go away.
8. ST-08. No Green felt mat (owner, 2026-10-20).
   1. It read as a poker table (H-05).
   2. A saved Green felt falls back to the default mat.

## Audio (later)

1. AU-01. Sound effects and music.
   1. Taken from any Isaac game: the Flash original or the C++ remake.
2. AU-02. Sounds for game actions.
   1. First list: dice roll, attack, damage dealt.
   2. Sources follow AU-01.
3. AU-03. Damage sounds (owner, 2026-10-13).
   1. A card hitting a card when anyone deals damage.
   2. Hearthstone's Jaina fireball for bombs and other effects.
   3. No background music in the files.
   4. Small files.
4. AU-04. Every sound at one loudness (owner, 2026-10-27).
   1. Styles cannot match, so at least loudness must.
   2. No sound feels louder or quieter than the rest.
   3. A tool checks every sound, the existing ones and every new one.
   4. The tool lives in `FourSoulsRepo/tools`, not here.
   5. Planned in 18.5 and 18.6; decision in ADR 022.
5. AU-05. Loot draw sound (owner, 2026-10-27).
   1. A looted card makes a sound, from the deck or the discard.
   2. One per card: start of turn, monster reward, other effects.
   3. Planned in 18.7.

## Menu

1. M-01. Main menu start options:
   1. Host a game.
   2. Join a game.
   3. Records: list, watch, export, and delete saved matches.
2. M-02. Exit button in the main menu (owner, 2026-10-20).
   1. Last button of the menu.
   2. Closes the app.
3. M-03. Menu split by where you play (owner, 2026-10-22).
   1. Main menu: Play, Records, Settings, Exit.
   2. Play: LAN, Online, Learn to play.
   3. Learn to play is disabled, marked "Soon" (LR-03 later).
   4. LAN: Host button on top.
   5. Below it: LAN games, searched all the time; no "Search again".
   6. Below them: an address field; it keeps the last typed value only.
   7. LAN has no Recent list.
   8. Online: an address field and the Recent list (N-18).
   9. Online opens the server's games; create a game there (N-17).
   10. LAN Host always hosts in this app; Online always on the server.
   11. Back and Esc go up one level.
   12. The setup screen goes back to the path that opened it.
4. M-04. Recent servers show live status (owner, 2026-10-22).
   1. A quick check when the Online screen opens.
   2. Shows whether the server is up and how many games are open.
5. M-05. "Rejoin" button in the main menu (owner, 2026-10-22).
   1. Shows when the app closed during a game: "Rejoin: <game> · LAN/Online".
   2. Uses the saved seat (`rejoin` in the settings, N-08).
   3. The saved seat is cleared when the game ends.
   4. Also cleared when the player leaves on purpose.
   5. A crash or a force-quit keeps it.
   6. On menu open the app asks that server once: is the seat still live?
   7. Yes: the button shows. No: it hides and the seat is cleared.
   8. A crashed LAN host has nothing to rejoin; it re-hosts from a save.
6. M-06. The Play group opens in place (owner, 2026-10-22).
   1. Play opens its buttons under itself, indented.
   2. The buttons below slide down; the ones above stay put.
   3. A second click on Play or Esc closes the group.
   4. No separate Play screen and no Back button there.
7. M-07. Records shows "Soon" while disabled (owner, 2026-10-22).
   1. Same look as Learn to play.

## Game setup

1. GS-01. Host picks how characters are chosen:
   1. Random: each player gets one random character (original rule).
   2. Draft: each player gets X characters and picks one.
   3. Ban phase: players ban characters, then get random ones.
   4. Ban rounds: each player bans one character per round.
   5. Host sets the number of ban rounds.
2. GS-02. Host can ban characters for the whole game.
   1. Banned characters are never dealt.
   2. Players see the list of banned characters.
   3. Can be combined with any picking mode.
3. GS-03. After the ban phase, characters are dealt:
   1. One random character each, or
   2. X random characters each; the player picks one.
4. GS-04. Ban order follows esports drafts (e.g. Dota 2 Captains Mode).
   1. Snake order: 1-2-3-4, then 4-3-2-1 in the next round.
   2. Ban timer: optional; host picks short timer or none.
   3. Bans are visible to everyone as they happen.
   4. When the ban timer runs out, that player makes no ban.
5. GS-05. Play modes (normal games and tournaments):
   1. Free-for-all: one against all others.
   2. Two vs two: teams of two.
6. GS-06. First player is chosen by a dice roll.
   1. The roll is shown to everyone.
   2. Check the rules: some characters may change this.
   3. Today the roll is silent; only the result shows.
   4. Everyone sees each player's die, re-rolls on ties included.
7. GS-07. Optional: deck ratios.
   1. When sets are mixed, decks follow Ed's recommended ratios.
   2. The official TTS table has the same option.
8. GS-08. Optional: remove multiplayer-only cards in 2-player games.
9. GS-09. Player count stays 2–4; 5–6 players not planned.
10. GS-10. Match setup has two levels.
   1. Main screen: only simple options with safe defaults.
   2. Goal: do not scare new players.
   3. Detailed setup: every extra option.
   4. Examples: Bonus Souls, Rooms, deck ratios, ban modes.
11. GS-11. Co-op and solitaire modes (future).
   1. Players together against the game, or one player alone.
   2. Official challenges, each against a boss (Rag Man, Mom, Greed, ...).
   3. Difficulty: Normal, Hard, Ultra Hard.
   4. Not part of the first main goal.
12. GS-12. Optional house rule: no damage loops at 0 HP.
   1. Official rules: damage at 0 HP still fires "when you take damage" (R-MECH-51).
   2. Example: Forever Alone with Gurdy Jr. loops while the death waits.
   3. With the option on, damage at 0 HP fires no "take damage" triggers.
   4. Off by default; the host turns it on per match.
   5. Detailed setup shows a hint that explains it.
   6. Research: Q-13 in `pkg/rules_engine/docs/rules/open-questions.md`.
13. GS-13. Host bans cards before hosting (owner, 2026-10-13).
   1. Banned cards are left out of every deck and never dealt.
   2. Players see the ban list (A-16: banned cards are public).
   3. Tests use the same bans for fixed deals.

## Spectators

1. SP-01. Viewers can join a match.
   1. Viewers do not see players' hands.
   2. Viewers see only steps that happened, never attempts.
   3. Host picks live view or delayed view.
   4. Players do not see who watches; only the viewer count.
   5. No chat for viewers for now.
   6. Viewers may connect to the game server directly; relay is optional.
   7. Without a relay, the game server applies the delay.
2. SP-02. A judge can join a match.
   1. The judge sees all players' hands.
   2. For now the judge only watches; no extra powers.
3. SP-03. Viewers never slow the game down.
   1. A slow viewer falls behind; the game does not wait.
   2. Many viewers must not overload the game server.
4. SP-04. Relay server for viewers (future).
   1. Players connect to the game server.
   2. Viewers connect to a separate relay server.
   3. The relay copies the public event stream to its viewers.
   4. Its own entry point (A-04).
   5. The relay applies the viewer delay.
   6. The judge never connects through a relay.
   7. No relay chains: a relay reads only from a game server.
   8. A new viewer first gets the current state, then live events.
   9. One relay can serve several games at once.
5. SP-05. Viewers and judges with an outdated client can still watch.
   1. Unknown cards are drawn as blank cards of the same size.
   2. The blank card shows title, text, and stats.
   3. A message says the client is outdated.
   4. Watching is never blocked.

## Cards

1. CD-01. Each client ships the full card collection.
   1. Server sends only card IDs.
2. CD-02. Host picks the card sets for a game.
   1. Examples: base, Gold Box, Requiem.
   2. A client missing a chosen set cannot join.
   3. Promo packs are selectable sets too; all packs, incl. NSFW.
   4. Start with Base Game only; other sets and packs come later.
3. CD-03. Collection check on join.
   1. Server sends the list of card IDs in the game.
   2. Client checks it has all of them.
   3. If not, it shows "Your client is outdated".
4. CD-04. Each card stores its artist, when known.
   1. Most artists are listed on the official site.
   2. Name only, no links.
5. CD-05. Card gallery in the app (future).
   1. Browse every card.
   2. Shows the artist's name for each card.
   3. Simple display for now.
   4. Later: search, filters, collection features.
6. CD-06. Each card stores its fan translators, when translated.
   1. Name only, no links.
7. CD-07. Card data is collected automatically.
   1. A tool in the separate `FourSoulsRepo/tools` repo gathers data from official sources.
   2. Names, text, artists, images.
   3. The official site is behind a Cloudflare check.
   4. The script opens a real browser; the user passes the check.
   5. The script then reuses that browser session to collect data.
   6. Second source: the official TTS table (Steam Workshop 2501791757).
   7. It gives card names and images (about 600×900 px per card).
   8. It has no card text; text comes from the official site.
8. CD-08. A changed card gets a new ID.
   1. Applies to text fixes and balance changes (nerfs).
   2. The old card version stays in the collection.
   3. Old replays keep working and show the old text.
   4. ID stays the same; a separate `version` field tells versions apart.
   5. First version omits the field; later ones are 2, 3, ...
   6. A card reference is ID plus version, e.g. `a_penny` or `a_penny@2`.
9. CD-09. Server sends card text for unknown cards on request.
   1. Title, text, and stats; no images.
   2. Used by outdated viewers and judges (SP-05).
10. CD-10. Fan-made sets: maybe later.
   1. Only after all official cards are done.
   2. Cards are Go code; adding one must be easy after cloning (ADR 005).
11. CD-11. Alt art (later).
   1. Same stats and text: same card, different image.
   2. Different stats or text: a new card or a new version (CD-08).
12. CD-12. Every copy keeps its own art.
   1. Coins and similar cards have a different picture per copy.
   2. A reprint in a newer set adds its copies to the same card.
   3. With Base Game and Gold Box, all copies of both go into the deck.
   4. Sets are ordered by release date, oldest first.
13. CD-13. Magic Marker changes any number in a card's text (future).
   1. Each card lists the numbers it may change.
   2. Until then only numbers in common blocks count (Q-14).
   3. Owner, 2026-10-19: a task for later.

## Build

1. B-01. Two build types:
   1. Embedded: cards and other files are packed into the binary.
   2. External: files are read from a folder next to the app.
   3. External mode is for development: no re-embedding on every change.
   4. Switched with the Go build tag `-tags embed`.
   5. External is the default for `wails dev`; releases are embedded.
   6. Changed files need an app restart; no hot reload.
2. B-02. Target platforms: Windows, Linux, macOS.
   1. Both x64 and ARM for each.
   2. Built and released by GitHub Actions.
   3. Test on Linux early: WebKitGTK is the weakest webview.
3. B-03. Release binaries are not code-signed.
   1. App identity (name, bundle ID) stays stable between versions.
   2. Best effort only; the OS may still warn on every update.
4. B-04. Wiki page: how to run the app when the OS blocks it.
5. B-05. Release archive includes `README.txt`.
   1. Explains how to get past the OS warning.
   2. macOS: right-click → Open.
   3. Windows: "More info" → "Run anyway".
6. B-06. Mobile (iOS, Android): maybe in the future.
   1. Only Wails v3 supports mobile, and it is experimental.
   2. Keep the UI touch-friendly from the start.
   3. Every hover action also works by tap or long press.
   4. Drag and drop works with touch (pointer events).
7. B-07. Setup scripts for developers and translators (owner, 2026-10-26).
   1. One script per OS installs every tool `make` needs.
   2. The wiki says what each `make` target needs, per OS.
   3. Only macOS and Linux for now; Windows later (owner, 2026-10-27).

## Legal

1. L-01. Start screen states clearly that this is a fan game.
   1. All rights belong to the game's owners.
   2. It names who has the rights to distribute the game.
   3. Shown for a few seconds at every app start.
   4. Original developers' names link to official resources.
   5. The links lead to where to buy the physical game.
2. L-02. Dedicated server prints the same notice to the console.
   1. Includes the links.
3. L-03. Card images ship with the app.
   1. The app is free; nothing is sold.

## Replay

1. RP-01. Matches are recorded.
2. RP-02. Recorded matches can be replayed.
   1. Replay may skip the rules engine and just show what happened.
   2. Replay can also check moves against a newer engine.
   3. It marks the moves that would not work there.
   4. This check is information only; steps are never rejected.
   5. The viewer can turn these marks off.
3. RP-03. Record format.
   1. Stores only steps that really happened.
   2. Rejected attempts are not recorded.
   3. Dice results are recorded.
   4. Replay validates the file format only, not the rules.
4. RP-04. Replays differ from live viewing.
   1. The record stores all hands.
   2. The replay viewer can show or hide hands.
5. RP-05. Who gets the record.
   1. Every player can save it, but only after the match is over.
   2. Never available mid-game: it would reveal hands.
6. RP-06. Tournament records are also kept on the server.
   1. They can be recovered even if no player saved them.
   2. Kept for 90 days.
7. RP-07. Every match is autosaved on the server.
   1. Not only tournaments.
   2. Dedicated server: records are kept for 30 days.
   3. Client-hosted server: kept until the user deletes them.
   4. Retention periods are set in the server config.
   5. Retention can be disabled to keep all records forever.
8. RP-08. Unfinished matches are saved too.
   1. Example: all players left before the end.
   2. May be disabled in the future.
9. RP-09. Records are stored compressed to keep them small.
   1. One file per match, own extension (e.g. `.fsrec`).
   2. Inside: readable JSON events.
   3. Compressed with gzip or zstd.
   4. Unpacks for reading by eye when debugging.
10. RP-10. Outdated clients can still watch replays.
   1. Same as live viewing (SP-05): blank cards with text.
   2. The record includes the text of every card used in the match.
11. RP-11. Replay playback speed is adjustable.
   1. Full controls: play, pause, speed.
   2. Step one event forward or back.
   3. Jump to any turn.
12. RP-12. Each recorded step stores a state checksum.
   1. Shows exactly where a replay and the engine disagree.

## Online services (future)

1. O-01. Separate online server.
   1. Registration with username and password.
   2. Lobby.
   3. Matchmaking.
   4. Rating.
   5. Leaderboard.

## Tournaments (future)

1. TR-01. Closed tournaments.
   1. Creator sends a link; players join through it.
   2. Joined players appear in the tournament bracket.
2. TR-02. Creator can disqualify players.
   1. For technical or other reasons.
3. TR-03. Open tournaments.
   1. A separate menu entry lists current tournaments.
4. TR-04. Creator sets the tournament rules.
   1. Card sets.
   2. How characters are picked (GS-01).
   3. Format: double elimination, and others.
   4. Host picks one format: single elim, double elim, Swiss, round robin.
   5. Play mode (GS-05).
5. TR-05. Who advances to the next round:
   1. Free-for-all: only the winner.
   2. Two vs two: both players of the winning team.
6. TR-06. Requires the online server (O-01).

## Chat (future)

1. CH-01. In-game chat.
   1. Messages are stored in the match record.
   2. Messages are shown in replays.

## Docs

1. DO-01. The wiki has three groups.
   1. Players: install, first start, play over LAN or internet.
   2. Server admins: how to run a dedicated server.
   3. Developers: layout, builds, engine, tests.
   4. `_Sidebar.md` and `Home.md` show the three groups.
   5. A fourth group: Translators (owner, 2026-10-26).
