# Step 8. First release (MVP)

Goal: the full Base Game for 2–4 players over LAN.
Public downloads for all platforms.

Ideas: B-03 – B-05, L-03, V-10, V-11, DO-01, V-14, A-16, N-13, N-14, V-19.

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
