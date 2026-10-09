# Step 8. First release (MVP)

Goal: the full Base Game for 2–4 players over LAN.
Public downloads for all platforms.

Ideas: B-03 – B-05, L-03.

---

### [ ] 8.1 Release workflow

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

### [ ] 8.2 Versioning

1. Goal: clear versions for app and engine.
2. Tasks:
   1. Semantic versions for the app.
   2. Engine tags in the form `rules_engine/vX.Y.Z`.
   3. Changelog file.
3. Done when:
   1. Release notes show both versions.

### [ ] 8.3 Stable app identity

1. Goal: same identity every release (B-03).
2. Tasks:
   1. Fixed name and bundle ID on macOS.
   2. Fixed product info on Windows.
   3. No code signing.
3. Done when:
   1. Two releases share the same identity.

### [ ] 8.4 OS warning help

1. Goal: users can start the app (B-04, B-05).
2. Tasks:
   1. `README.txt`: macOS right-click → Open.
   2. `README.txt`: Windows "More info" → "Run anyway".
   3. Linux notes if needed.
   4. Wiki page with screenshots.
3. Done when:
   1. A new user starts the app with only the README.

### [ ] 8.5 Play-test and fixes

1. Goal: real games with real players.
2. Tasks:
   1. Play full games with 2, 3, 4 players.
   2. Collect issues; fix blockers.
   3. New rule disputes go to the open list.
3. Done when:
   1. **Owner** approves the release.

### [ ] 8.6 Docs check

1. Goal: docs match the released app.
2. Tasks:
   1. Every MVP feature has a use case.
   2. README is short and current.
   3. Accepted ADRs match the code.
3. Done when:
   1. Use case index has no "Planned" item for MVP features.
