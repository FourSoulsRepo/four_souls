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
| 9 | Spectators | Viewers and judges watch live | 8 |
| 10 | Replays | Records screen and replay player | 8 |
| 11 | Learning | Rules viewer, sandbox, tutorial, bot | 8 |
| 12 | Rules website | Engine in WebAssembly on GitHub Pages | 11 |
| 13 | Web client | Browser client served by the game server | 9 |
| 14 | More sets | Gold Box, Four Souls+, Requiem, promos | 8 |
| — | Future | Not scheduled yet | — |

Step 3 reuses the browser scraping tool from step 2.
Steps 5, 6 and 7 can run in parallel.
Step 7 ends with an internal build on a partial card set.
Step 8 needs all Base Game cards from step 5.
Steps 9–14 can be reordered after the MVP.

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
8. Ideas: N-01 – N-09, CD-02, CD-03, GS-01 – GS-06, GS-09, GS-10, A-12, RP-01, RP-03, RP-05, RP-07 – RP-09, RP-12.

## 7. Client

1. Main menu: host, join, settings.
2. Table for 2–4 players, Hearthstone-style flat cards.
3. Drag and drop, animations, card zoom.
4. Stack panel, history panel, turn indicator, counters.
5. Response UI: allowed actions, skip, timer.
6. Settings: fullscreen, theme, game mat, animation speed.
7. Touch-friendly; tested on Linux.
8. Internal play-test build on the cards done so far.
8. Ideas: M-01, V-01 – V-08, ST-01 – ST-03, ST-05, B-06.

## 8. First release (MVP)

1. Full Base Game for 2–4 players over LAN.
2. Release pipeline: archives for all platforms.
3. `README.txt` and wiki page for OS warnings.
4. Stable app identity; no code signing.
5. Ideas: B-03 – B-05, L-03.

## 9. Spectators

1. Viewers and judges join a live match.
2. Live or delayed view; viewer count for players.
3. Viewers never slow the game.
4. Outdated clients still watch with blank cards.
5. Plain-text judge log.
6. Ideas: SP-01 – SP-03, SP-05, CD-09, V-04.

## 10. Replays

1. Records screen: list, watch, export, delete.
2. Replay player: speed, pause, step, jump to turn.
3. Show or hide hands.
4. Optional check against the current engine.
5. Ideas: RP-02, RP-04, RP-10, RP-11, M-01.

## 11. Learning

1. Rules text viewer in the app.
2. Sandbox on the game table; export payload and link.
3. Tutorial lessons.
4. Simple random bot; seeded and free practice games.
5. FAQ from frequent rules questions, linked to rule IDs.
6. Ideas: LR-01 – LR-06, LR-08.

## 12. Rules website

1. Engine compiled to WebAssembly.
2. Simple sandbox: payload or link in, answer with reason out.
3. Rules text with links by rule ID.
4. Hosted on GitHub Pages.
5. Ideas: A-03.

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

## Future

Not scheduled; revisit after the MVP.

1. Relay server for viewers (SP-04).
2. Online server: accounts, lobby, matchmaking, rating (O-01).
3. Tournaments (TR-01 – TR-06, RP-06).
4. In-game chat (CH-01).
5. Co-op and solitaire modes (GS-11).
6. Phase stops (N-10).
7. Localization (LO-01 – LO-03).
8. Audio (AU-01).
9. Card gallery (CD-05).
10. Keyboard shortcuts (ST-04).
11. Fan-made sets (CD-10).
12. Mobile apps (B-06).
