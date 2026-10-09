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

## Localization

1. LO-01. First version: English only.
2. LO-02. Later: user picks the app language.
   1. Cards are shown in the chosen language.
3. LO-03. Later: view a card in other languages right in the game.
   1. How exactly: TBD.
   2. Proposal: hover or zoom a card, a hotkey cycles its languages.

## Settings

1. ST-01. Fullscreen or windowed mode.
2. ST-02. Theme: dark or light.
   1. Changes menus and panels only for now.
3. ST-03. Game mat: choose one of the official mats.
   1. Mat images are collected by the scraping script (CD-07).
   2. Bad images are replaced by hand.
   3. Maybe reuse TTS table themes too (Basement, Cathedral).
4. ST-04. Keyboard shortcuts (future).
   1. Not in the first version.
   2. Configurable in the settings.
5. ST-05. Animation speed setting.

## Audio (later)

1. AU-01. Sound effects and music.
   1. Taken from any Isaac game: the Flash original or the C++ remake.

## Menu

1. M-01. Main menu start options:
   1. Host a game.
   2. Join a game.
   3. Records: list, watch, export, and delete saved matches.

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
11. CD-11. Alt art (later).
   1. Same stats and text: same card, different image.
   2. Different stats or text: a new card or a new version (CD-08).

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
