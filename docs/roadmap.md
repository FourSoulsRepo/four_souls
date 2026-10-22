# Roadmap

Global steps; each has a detailed file in [roadmap/](roadmap/).
Idea IDs refer to [ideas.md](ideas.md).
Each step ends in something that runs and is tested.

---

## Overview

| Step | Name | Result | Depends on |
|------|------|--------|------------|
| 1 | Foundation | Repo layout, CI, versions, legal splash | — |
| 2 | Card data | Base Game data and images, scraped automatically | 1 |
| 3 | Rules research | Rules downloaded from all sources; our rules text with stable IDs | 2 |
| 4 | Engine core | General rules, stack, effect blocks, test tools | 3 |
| 5 | Base Game cards | Every Base Game card implemented and tested | 2, 4 |
| 6 | Server | Dedicated and hosted server, lobby, records | 4 |
| 7 | Client | Playable table in the desktop app; internal play-test build | 4, 6 |
| 8 | First release | **MVP:** 2–4 players, Base Game, over LAN | 7 |
| 15 | Visual refresh | Table layout reworked after fouronlinesouls.net | 8 |
| 16 | Localization | Menus, messages, and cards in several languages | 8 |
| 17 | Rollback | Other players can undo a mistaken event | 8 |
| 14 | More sets | Gold Box, Four Souls+, Requiem, promos | 8 |
| 18 | Play-test | New-user check, play-test, docs check | 14 |
| 9 | Spectators | Viewers watch live | 8 |
| 10 | Replays | Records screen and replay player | 8 |
| 11 | Learning | Rules viewer, sandbox, tutorial, bot | 8 |
| 12 | Rules website | Engine in WebAssembly on GitHub Pages | 11 |
| 13 | Web client | Browser client served by the game server | 9 |
| 19 | Bugs | Every bug report, logged and fixed under the roadmap rules | — |
| — | Future | Not scheduled yet | — |

Step 3 reuses the browser scraping tool from step 2.
Steps 5, 6 and 7 can run in parallel.
Step 7 ends with an internal build on a partial card set.
Step 8 needs all Base Game cards from step 5.
After the MVP the order is 15, 16, 17, 14, 18, 9, 10, 11, 12, 13.
Step numbers stay as IDs; the table lists steps in work order.

---

## 1. Foundation

1. Repo layout
   1. `pkg/rules_engine`, `pkg/card_db`, `pkg/record`: own Go modules.
   2. `internal/protocol`.
   3. Entry points: `main.go`, `cmd/server`, `cmd/website`.
2. CI on GitHub Actions
   1. Go tests, frontend build.
   2. Engine built for `js/wasm` from day one.
   3. App built for Windows, Linux, macOS (x64 + ARM).
3. Build tags: embedded and external files.
4. App and engine versions in the main menu corner.
5. Fan-game notice on start; same notice on the server console.
6. Ideas: PR-01, PR-02, A-01, A-02, A-04, A-09, A-10, A-13, A-14, B-01, B-02, L-01, L-02.

## 2. Card data

1. Scraper in the separate `FourSoulsRepo/tools` repo
   1. Official site via a real browser (Cloudflare check).
   2. Official TTS table for names and images.
2. `pkg/card_db`: text, stats, artists, versions.
3. Private image repo: WebP, one commit, one folder per set.
4. Fallback: blank cards with text when images are missing.
5. Card stub generator and card status report.
6. Scope: Base Game only.
7. Ideas: CD-01, CD-04, CD-06, CD-07, CD-08, A-11, TS-02, TS-03.

## 3. Rules research

1. Download rules from all sources with scripts.
   1. Solve Cloudflare checks with the browser tool from step 2.
2. Collect Base Game rules from all sources.
3. Write our rules text with stable rule IDs.
4. Check Ed's rulings against current rules.
5. List open questions: mulligan, first player, others.
6. Ideas: R-01 – R-07.

## 4. Engine core

1. Game state, turn phases, priority, the stack.
2. Pending actions that card effects can rewrite.
3. Reusable effect blocks.
4. Deterministic RNG; seed stays on the server.
5. One view filter for player, viewer, judge.
6. Allowed actions per player.
7. Situation payload format and test runner.
8. Fuzz test harness: all-cards and focused modes.
9. Ideas: A-05 – A-08, N-04, LR-04, LR-07, TS-05.

## 5. Base Game cards

1. Implement every Base Game card, test-first.
2. Interaction tests for card combinations.
3. Card status report shows 100% for Base Game.
4. Ideas: R-05, TS-01, TS-04.

## 6. Server

1. Protocol between client and server.
2. Dedicated server and hosting from the app.
3. Lobby: create, join, sets, collection check.
4. Character picking modes and bans; two-level match setup.
5. Auto-skip, "Skip all", response timer.
6. Reconnect, pause, vote to kick.
7. Match records: autosave, retention, checksums.
8. Ideas: N-01 – N-09, N-11, CD-02, CD-03, GS-01 – GS-06, GS-09, GS-10, A-12, RP-01, RP-03, RP-05, RP-07 – RP-09, RP-12.
9. Save a game and continue later with the same players.

## 7. Client

1. Main menu: host, join, settings.
2. Table for 2–4 players, Hearthstone-style flat cards.
3. Drag and drop, animations, card zoom.
4. Stack panel, history panel, turn indicator, counters.
5. Response UI: allowed actions, skip, timer.
6. Settings: fullscreen, theme, game mat, animation speed.
7. Touch-friendly; tested on Linux.
8. Internal play-test build on the cards done so far.
8. Ideas: M-01, V-01 – V-08, ST-01 – ST-03, ST-05, ST-06, B-06, N-11.

## 8. First release (MVP)

1. Full Base Game for 2–4 players over LAN.
2. Release pipeline: archives for all platforms.
3. `README.txt` and wiki page for OS warnings.
4. Stable app identity; no code signing.
5. Ideas: B-03 – B-05, L-03.
6. Card rulings from Yuggy and Discord become tests, before the play-test.
7. No page scroll and no drag selection, before the play-test (V-10).
8. Stack shows whom it waits for, before the play-test (V-11).
9. Wiki in three groups: players, server admins, developers (DO-01).
10. Stolen card shown to the thief and the victim (V-14).
11. Audit card lists for leaks of hidden cards (A-16).
12. LAN games found without typing an IP (N-13).
13. Game passwords (N-14).
14. New-user check, play-test and docs check: moved to step 18.
15. Dice looks: a real dice texture; core for choosing one (V-19).
16. Seamless server restart and idle game limits (N-20).

## 9. Spectators

1. Viewers join a live match.
2. Live view; viewer count for players.
3. Viewers never slow the game.
4. Outdated clients still watch with blank cards.
5. Judges and viewer delay: moved to Future.
6. Ideas: SP-01, SP-03, SP-05, CD-09.

## 10. Replays

1. Records screen: list, watch, export, delete.
2. Replay player: speed, pause, step, jump to turn.
3. Show or hide hands.
4. Optional check against the current engine.
5. Ideas: RP-02, RP-04, RP-10, RP-11, M-01.
6. Replays never show live aiming or hover (N-15).

## 11. Learning

1. Rules text viewer in the app.
2. Sandbox on the game table; export payload and link.
3. Tutorial lessons.
4. Simple random bot; seeded and free practice games.
5. FAQ from frequent rules questions, linked to rule IDs.
6. Ideas: LR-01 – LR-06, LR-08.
7. Any event on the stack in payloads (LR-09).
8. Sandbox tools modelled on the fouronlinesouls.net cheat menu.
9. BoardGameGeek rules forum: a source for the FAQ and tests (R-02, 11.7).

## 12. Rules website

1. Engine compiled to WebAssembly.
2. Simple sandbox: payload or link in, answer with reason out.
3. Rules text with links by rule ID.
4. Hosted on GitHub Pages.
5. Ideas: A-03.
6. Step-by-step stack resolution (LR-10).

## 13. Web client

1. React client built for the browser.
2. Served by the game server; play and watch.
3. Images only if the image folder is present.
4. Ideas: A-15.

## 14. More sets

1. Gold Box, Four Souls+, Requiem, Warp Zone, promos.
2. Detailed setup: Bonus Souls, Rooms, deck ratios.
3. Optional removal of multiplayer-only cards.
4. Ideas: CD-02, GS-07, GS-08, CD-11.
5. Deck ratios never reveal the dealt decks (A-16).
6. Reprints add copies; every copy keeps its own art (CD-12, ADR 010).
7. The host picks the card sets; the Base Game is always in (14.6).
8. Switch character and eternal item art between reprints (V-32, 14.7).
9. House rule option: no damage loops at 0 HP (GS-12, 14.8).
10. Host bans cards before hosting (GS-13, 14.9).
11. Tests ban stat-changing monsters for fixed deals (14.10).
12. Draft pick screen for normal games, no bans (GS-01, 14.11).
13. Ban list in the host screen (GS-13, 14.12).
14. Requiem: rooms on the table (14.13).
15. Requiem: attack a player or an attackable item in the app (14.14).
16. Requiem: two-sided cards (14.15).

## 15. Visual refresh

1. Study fouronlinesouls.net: start a game with a bot.
2. Borrow layout and readability ideas; copy no assets.
3. Rework the table; keep e2e tests and screenshots green.
4. Ideas: V-09.
5. Check Hearthstone again; the table must not feel like poker (V-13).
6. Attack animation: the attacker hits the monster (V-12).
7. Responsive, simple, informative UI (V-15, V-18).
8. Drag cards onto the table or a target, Hearthstone-style (V-16).
9. Real-looking dice, thrown by a button (V-17).
10. Everyone sees the roller's own dice look (V-19).
11. The game always shows clearly when it waits for the player (V-20).
12. Everyone sees what the active player is doing now (V-21).
13. Compare with Hearthstone screenshots from the internet (V-22).
14. Table areas keep their size when a zone empties (V-23).
15. Game mat images as the background; colors without images (V-24).
16. Mats and dice in WebP, like the cards (V-25).
17. Card backs on decks and opponent hands (V-26).
18. Study counters, stack window, and cheat menu of the reference (V-27).
19. Play a card by one click or by drag; a setting picks (V-28).
20. No number badges on cards (V-29).
21. Turn action icons with hotkeys (15.12).
22. Stack panel and history strip with icons (15.13).
23. Counters and live stats on the cards (15.14).
24. First study how Hearthstone aims at targets (15.15).
25. Big sub-steps split into one-session ones (15.16 – 15.31).
26. Aiming like Hearthstone: arrow, start, cancel (15.32 – 15.34).
27. Others see live aiming; never recorded (N-15, 15.35, 15.36).
28. Hit and death marks: starbursts, red HP, skull (15.37).
29. Souls as Hearthstone-style crystals, left of the character (V-30, V-31, 15.39).
30. Isaac fonts, foursouls.ru banner and background (V-33 – V-35, 15.40 – 15.42).
31. Everyone sees the first-player roll (GS-06, 15.43).
32. "Buy" and "Attack" declarations on the stack (V-36, 15.44).
33. After the game, all hands are revealed (V-37, 15.45).
34. Failed effects show in the history (V-38, 15.46).
35. Damage a player takes names its card (play test, 15.47).
36. Stack and history rows fade and slide (V-39, 15.51).
37. Played cards fly hand → stack → discard (V-40, 15.52).
38. Icons for the sources of HP and ATK boosts (V-41, 15.54).
39. Curses left of the player (V-42, 15.55).
40. Soul cards shown by hovering a soul crystal (V-43, 15.56).
41. Purchases shown in the history (V-44, 15.57).
42. "Choose one-" cards ask in a dialog like The Bone (V-45, 15.58).
43. Dice throw split into engine, server, button, animation, last roll (15.59 – 15.63).
44. Stack and history tell every play: who did what (V-51, 15.88 – 15.91).
45. Real stats: breakdown, temporary icons, fight odds (V-41, 15.92 – 15.94).

## 16. Localization

1. Translation layer for the frontend; ADR picks the library.
2. Server sends message keys and values, not text.
3. Card names and text from card_db translations.
4. Language setting; first extra language: Russian.
5. Ideas: LO-01 – LO-04.

## 17. Rollback

1. Owner decides limits: hidden info, dice, who must agree.
2. Engine undoes the last event from a saved state.
3. Only other players can roll back; never the author.
4. History and records show the rollback.
5. Ideas: N-12.
6. Dice keep their value after a rollback (N-12).

## 18. Play-test

1. A new user starts the app from the README alone.
2. Full games with 2, 3, 4 players on all sets so far.
3. Docs match the app: use cases, README, ADRs.
4. Moved from step 8 (8.4–8.6): more cards make a better test.
5. Ideas: B-04, B-05.
6. Review the multi-instance solution (UC-16).

## 19. Bugs

1. Every bug report from play gets an item in `roadmap/19-bugs.md`.
2. Each fix follows the roadmap rules: tests, docs, screenshots, review.
3. Open the whole time; never closed as a step.

## Future

Not scheduled; revisit after the MVP.

1. Relay server for viewers (SP-04).
2. Online server: accounts, lobby, matchmaking, rating (O-01).
3. Tournaments (TR-01 – TR-06, RP-06).
4. In-game chat (CH-01).
5. Co-op and solitaire modes (GS-11).
6. Phase stops (N-10).
7. Localization: moved to step 16.
8. Audio (AU-01, AU-02).
9. Card gallery (CD-05).
10. Keyboard shortcuts (ST-04).
11. Fan-made sets (CD-10).
12. Mobile apps (B-06).
13. Tournament tools, moved from steps 7 and 9.
   1. Ban and pick screen (GS-01 – GS-04 UI); server part is done.
   2. Judges: role, all-hands view, plain-text log (SP-02, V-04).
   3. Delayed view for viewers (SP-01).
   4. The plain draft pick screen moved to 14.11.
14. Magic Marker: every card lists the numbers it may change (CD-13).
