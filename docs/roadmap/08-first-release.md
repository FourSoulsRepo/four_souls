# Step 8. First release (MVP)

Goal: the full Base Game for 2–4 players over LAN.
Public downloads for all platforms.

Ideas: B-03 – B-05, L-03, V-10, V-11, DO-01, V-14, A-16, N-13, N-14, V-19, N-19, N-20.

---

### [x] 8.1 Release workflow

1. Goal: a tag builds and publishes a release.
2. Tasks:
   1. Trigger on version tags.
   2. Build Windows, Linux, macOS × x64, ARM.
   3. Check out the images with `.github/actions/card-images` (deploy key).
   4. Build with `-tags "embed cardimages"` (2.7).
   5. Dedicated server binaries too.
   6. Archives include `README.txt`.
3. Done when:
   1. A test tag produces all archives.

### [x] 8.2 Versioning

1. Goal: clear versions for app and engine.
2. Tasks:
   1. Semantic versions for the app.
   2. Engine tags in the form `rules_engine/vX.Y.Z`.
   3. Changelog file.
3. Done when:
   1. Release notes show both versions.

### [x] 8.3 Stable app identity

1. Goal: same identity every release (B-03).
2. Tasks:
   1. Fixed name and bundle ID on macOS.
   2. Fixed product info on Windows.
   3. No code signing.
3. Done when:
   1. Two releases share the same identity.

### [~] 8.4 OS warning help

1. Goal: users can start the app (B-04, B-05).
2. Tasks:
   1. `README.txt`: macOS right-click → Open.
   2. `README.txt`: Windows "More info" → "Run anyway".
   3. Linux notes if needed.
   4. Wiki page; text only, no images while the visuals change.
3. Done when:
   1. A new user starts the app with only the README.
4. Moved: the new-user check is 18.1, after more sets.

### [~] 8.5 Play-test and fixes

1. Goal: real games with real players.
2. Tasks:
   1. Play full games with 2, 3, 4 players.
   2. Collect issues; fix blockers.
   3. New rule disputes go to the open list.
3. Done when:
   1. **Owner** approves the release.
4. Moved to 18.2: few cards so far; play-test after more sets.

### [~] 8.6 Docs check

1. Goal: docs match the released app.
2. Tasks:
   1. Every MVP feature has a use case.
   2. README is short and current.
   3. Accepted ADRs match the code.
3. Done when:
   1. Use case index has no "Planned" item for MVP features.
4. Moved to 18.3, after the play-test.

### [x] 8.7 Rulings from Yuggy and Discord

1. Goal: card rulings beyond R-07 and the FAQ become tests (R-07, TS-04).
2. When: before the play-test (8.5).
3. Tasks:
   1. Search the step 3 digests for every Base Game card name, English and Russian.
   2. Yuggy first: Reddit comments and X posts, high trust.
   3. Then the Discord channel, threads about specific cards only.
   4. Write each ruling in our own words.
   5. Matching rulings become interaction tests.
   6. Contradictions are fixed, or go to `open-questions.md` for the **Owner**.
4. Done when:
   1. Every Base Game card was searched in both sources.
   2. Found rulings have tests or open questions.

### [x] 8.8 Not a web page

1. Goal: the app does not scroll or select like a web page (V-10).
2. When: before the play-test (8.5).
3. Tasks:
   1. No whole-page scroll on any screen.
   2. Lists keep their own scroll bars.
   3. Disable drag selection of text and elements.
   4. Keep selection where copying helps: logs, payloads, addresses.
   5. e2e test: a drag across the table selects nothing.
4. Done when:
   1. No screen scrolls or selects at any window size.

### [x] 8.9 Stack shows whom it waits for

1. Goal: everyone sees who must act now (V-11).
2. When: before the play-test (8.5).
3. Tasks:
   1. Server sends who still has to respond and who passed.
   2. Stack panel lists them by nickname.
   3. Pending prompt: who must choose and what.
   4. Auto-skipped players are marked as passed.
4. Done when:
   1. In a 4-player game, every client names the same waiting players.
   2. Use case updated: respond to the stack.

### [x] 8.10 Wiki in three groups

1. Goal: each reader finds their pages fast (DO-01).
2. Tasks:
   1. Groups: players, server admins, developers.
   2. Sort every page into one group.
   3. `Home.md` and `_Sidebar.md` show the groups.
   4. Page slugs stay the same, so links keep working.
3. Done when:
   1. Every wiki page sits in exactly one group.

### [x] 8.11 Stolen card in the history

1. Goal: both sides of a steal see the card (V-14).
2. When: before the play-test (8.5).
3. Tasks:
   1. Today `EvGaveCard` is private to the receiver only.
   2. Let a private event name every player who sees the card.
   3. Thief and victim see the card; others see "A stole a card from B".
   4. Engine view test for all three seats.
4. Done when:
   1. Dead Bird steal shows the card to exactly two players.

### [x] 8.12 Card list leak audit

1. Goal: no message reveals hidden cards (A-16).
2. When: before the play-test (8.5).
3. Tasks:
   1. Review every server message: views, setup, games list, errors.
   2. Review card text sent to outdated clients (CD-09).
   3. Records go only to players after the match (RP-05).
   4. Collection check stays by whole sets.
   5. Fuzz invariant: views never name a card hidden in a deck.
4. Done when:
   1. Audit notes in `docs/wiki/Server.md`; found leaks fixed with tests.

### [x] 8.13 LAN game discovery

1. Goal: join a friend's game without typing an IP (N-13).
2. When: before the play-test (8.5).
3. Tasks:
   1. ADR: UDP broadcast or mDNS; port; message format.
   2. Hosting app announces while the game waits for players.
   3. Join screen lists found games; manual address stays.
   4. Dedicated server never announces.
   5. Check macOS local network and Windows firewall prompts.
   6. Use case: join a game found on the LAN.
4. Done when:
   1. Two machines on one Wi-Fi find each other's game.
4. **Owner** checked LAN search on real machines (2026-10-12).

### [x] 8.14 Game password

1. Goal: only invited players join (N-14).
2. Tasks:
   1. Optional password when creating a game.
   2. Games list shows a lock.
   3. Join asks for it; a wrong one is refused with a reason.
   4. Password never appears in records or logs.
3. Done when:
   1. Server tests for right, wrong, and missing passwords.
   2. Use case updated: host a game, join a game.

### [x] 8.15 Dice looks

1. Goal: dice show a real texture; the core for choosing one (V-19).
2. Tasks:
   1. Go lists dice looks; `App.Dice` gives their image URLs.
   2. Textures are served with card images, PNG allowed.
   3. Setting `dice` keeps the chosen look; default "tear".
   4. `Die` cuts one face out of the 3×3 texture.
   5. Plain die when the texture is missing or the value is not 1–6.
   6. Textures must be added to the private submodule via the tools repo.
3. Done when:
   1. Go and Vitest tests for the catalog, setting, and faces.
   2. All six faces checked on two real textures.

### [x] 8.16 Server ENV config

1. Goal: Docker sets every server option by ENV (N-19).
2. Tasks:
   1. `FOUR_SOULS_<OPTION>` for every flag, and `FOUR_SOULS_CONFIG`.
   2. Order: defaults, then JSON file, then ENV, then flags.
   3. A bad value stops the server with a clear error.
3. Done when:
   1. Tests for each variable, the order, and bad values.

### [x] 8.17 Server release archives

1. Goal: server downloads apart from the app (N-19).
2. Tasks:
   1. A release job builds only the server, with `CGO_ENABLED=0`.
   2. Names without a version: `four-souls-server-<os>-<arch>`.
   3. Each has the server, its README, `LICENSE`, `NOTICE`.
   4. App archives keep their versions and lose the server.
3. Done when:
   1. CI builds the server archives.
   2. The Linux server binaries are statically linked.

### [x] 8.18 Secure WebSocket in the client

1. Goal: join a server behind HTTPS by its domain (N-19).
2. Tasks:
   1. `wss://` and `ws://` typed by the player are kept.
   2. A bare domain tries `wss://` first, then `ws://` on 4774.
   3. IPs, `localhost`, LAN names, and a port keep plain `ws://`.
   4. The record link follows: `https://` for `wss://`.
   5. ADR for `wss://` through a reverse proxy.
   6. Use case: join a server by its domain.
3. Done when:
   1. Vitest for every address form; e2e for the fallback.

### [x] 8.19 Dockerfile

1. Goal: one file builds a server image from the latest release (N-19).
2. Tasks:
   1. Alpine image; the server runs as a non-root user.
   2. Downloads `releases/latest/download/…` for the build's CPU.
   3. Records and saves in the `/data` volume.
   4. `/healthz` on the server; `HEALTHCHECK` uses it.
3. Done when:
   1. The image builds and the healthcheck passes.

### [x] 8.20 Docker Compose with Caddy

1. Goal: HTTPS server with one command (N-19).
2. Tasks:
   1. `docker-compose.yml`, `Caddyfile`, `.env.example`.
   2. Caddy gets the certificate; only ports 80 and 443 open.
   3. A healthcheck tracks the server; Caddy waits for it.
3. Done when:
   1. `docker compose up` serves `wss://` on a test domain.

### [x] 8.21 Server setup pages

1. Goal: the operator picks systemd, Docker, or compose (N-19).
2. Tasks:
   1. VPS page keeps the shared parts and the options table.
   2. Options table gets an ENV column.
   3. One wiki page per setup; how to update each.
   4. ADR: server deployment from releases.
   5. `NOTICE` lists Alpine and Caddy.
3. Done when:
   1. Each page's commands were run once.

### [x] 8.22 Server keeps its games over a restart

1. Goal: a stopped server leaves its tables for the next start (N-20).
2. Tasks:
   1. On stop, every client gets a new `restarting` message.
   2. Every table goes to a handoff file in the saves folder.
   3. That covers started games and lobby tables.
   4. A file keeps seats, tokens, and options.
   5. It keeps the state checksum and the record's path.
   6. The record gets no `end` line on such a stop.
   7. On start, the tables come back; seats wait for their tokens.
   8. The record goes on with a `resume` part (ADR 007).
   9. A loaded handoff file is deleted.
   10. Protocol version stays: older apps ignore the new messages.
   11. ADR: seamless server restart.
   12. Bans and picks start over once every player is back.
   13. Only the dedicated server keeps games; an app's host does not.
3. Done when:
   1. Go test: stop, start on the same folder, rejoin by token.
   2. The state checksum is the same after the restart.
4. Owner notes (2026-10-21)
   1. Owner marks 8.22 solved.

### [x] 8.23 Games from another engine version

1. Goal: a newer engine goes on with every game it can read (N-20).
2. Tasks:
   1. Same engine version: the game goes on silently.
   2. Other version: load, save again; the bytes must be equal.
   3. Equal bytes and meaning number: the game goes on.
   4. Its players see "The engine was updated; cards may behave differently."
   5. Else the game ends unfinished; the record gets `end`.
   6. Its players see "Sorry, the server was updated; please start a new game."
   7. The server logs why; the restart file stays, renamed `.lost`.
   8. The ended table answers its tokens with that apology for 30 minutes.
   9. Meaning number: a constant in the engine, raised by hand.
   10. Raise it when a field keeps its name but changes meaning.
   11. New state fields use `omitempty`, so old games still match.
   12. A test checks every state field for `omitempty` or an allow list.
   13. Use case: a game goes on after a server update.
3. Done when:
   1. Go tests for the three cases.
   2. A renamed field and a new field each have a test.
4. Owner notes (2026-10-21)
   1. Owner marks 8.23 solved.

### [x] 8.24 Client reconnects by itself

1. Goal: players do not notice a short server restart (N-20).
2. Tasks:
   1. A lost game connection reconnects by itself.
   2. Retries: at once, then 1, 2, 4 … seconds, at most 10.
   3. First a small "Reconnecting…" badge, no dialog.
   4. The "Connection lost" dialog comes after 30 seconds.
   5. Its Reconnect button still works.
   6. The engine notice and the apology show as dialogs.
   7. No kick vote for 60 seconds after the server starts.
   8. Animations, history, and the board are kept.
   9. A failed retry's error shows in the dialog, not in the action bar.
   10. `cmd/server` stops its games after the warning, not on the signal.
3. Done when:
   1. Vitest for the retry delays.
   2. e2e: restart the server mid-game; the game goes on, no dialog.
   3. Screenshots of the badge and both dialogs.
4. Owner notes (2026-10-21)
   1. Owner marks 8.24 solved.

### [x] 8.25 Idle game limits

1. Goal: stuck and endless games do not live forever (N-20).
2. Tasks:
   1. No player connected for 30 minutes: the game ends.
   2. Server downtime counts toward those 30 minutes.
   3. A game ends 24 hours after it started, even if still played.
   4. It ends unfinished; the record gets `end`.
   5. Connected players see why it ended.
   6. Both limits are server options with ENV variables.
   7. 0 turns a limit off.
   8. Use case: a game ends by an idle limit.
   9. Flags `-empty_minutes`, `-game_hours`; ENV `FOUR_SOULS_EMPTY_MINUTES`, `FOUR_SOULS_GAME_HOURS`.
   10. Players away when it ended hear why when they come back.
   11. A table kept over a restart closes after the empty limit too.
3. Done when:
   1. Go tests with a fake clock for both limits.
4. Owner notes (2026-10-21)
   1. Owner marks 8.25 solved.

### [x] 8.26 Restart docs for server admins

1. Goal: admins know an update keeps the games (N-20).
2. Tasks:
   1. Each setup page: an update keeps games; the version cases.
   2. Options table: both idle limits.
   3. Saves folder must stay on disk between restarts.
   4. Docker: `stop_grace_period` is enough for the handoff.
3. Done when:
   1. A real server was stopped and started mid-game, without Docker (owner's choice).
4. Owner notes (2026-10-21)
   1. Owner marks 8.26 solved.

### [x] 8.27 Binary named FourSouls

1. Goal: release binaries are named `FourSouls`, not `four_souls`.
2. Tasks:
   1. `wails.json`: `outputfilename` becomes `FourSouls`.
   2. Archives hold `FourSouls.app`, `FourSouls.exe`, `FourSouls`.
   3. Bundle ID and the config folder `FourSouls` stay the same (8.3).
   4. `TestStableIdentity` checks the new file name.
   5. `release.yml` and `ci.yml` steps that name the binary.
   6. `README.md`, `release/README.txt`, `.gitignore`.
   7. Wiki: First start, Internal builds, Releases.
   8. Use case UC-16: `open -n FourSouls.app`.
   9. `wails.json` `name` becomes `FourSouls` too.
   10. On macOS Wails names the `.app` after `name`.
3. Done when:
   1. A test tag's archives hold `FourSouls` on every platform.
   2. An old release's settings still load after the update.
4. Owner notes (2026-10-22)
   1. Server binary keeps `four_souls_server`.
   2. Archive names `four-souls-<version>-…` stay.
   3. `name` changes too, so macOS gets `FourSouls.app`.
   4. Owner marks 8.27 solved.
5. Notes
   1. The workflows never name the app binary; they copy `build/bin/`.
   2. `.gitignore` keeps `/four_souls`: a plain `go build .` makes it.
   3. A local `wails build` gives `FourSouls.app`, same bundle ID.
   4. No test tag was pushed; the release job copies the same `build/bin/`.
