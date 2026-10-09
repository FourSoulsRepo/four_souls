# Step 6. Server

Goal: games run on a server; clients only send intents.
Tested with headless Go clients before the UI exists.

Ideas: N-01 – N-09, N-11, CD-02, CD-03, GS-01 – GS-06, GS-09, GS-10, A-12, RP-01, RP-03, RP-05, RP-07 – RP-09, RP-12.

---

### [x] 6.1 Protocol ADR

1. Goal: agree on the wire format.
2. Tasks:
   1. WebSocket transport, JSON messages.
   2. Message kinds: hello, lobby, intent, view, events.
   3. Also: allowed actions, resync request, error with reason.
   4. Protocol version in hello; clear "client outdated" error.
   5. Card refs only; no card text (CD-01).
3. Done when:
   1. **Owner** accepts the ADR.
   2. Types live in `internal/protocol`.

### [x] 6.2 Game room

1. Goal: one game runs safely on the server.
2. Tasks:
   1. One goroutine per game owns the engine state.
   2. Connections send intents into the room.
   3. Room sends each seat its own view and events.
   4. Session tokens for reconnect.
3. Done when:
   1. Race detector tests pass with 4 clients.

### [x] 6.3 Dedicated server

1. Goal: `cmd/server` runs games (N-01).
2. Tasks:
   1. Flags plus an optional JSON config file.
   2. Address, port, records folder, retention.
   3. Prints the notice (L-02) and versions.
   4. Clean shutdown saves running records.
   5. Check the wiki page "Dedicated server on a VPS" against the real flags.
3. Done when:
   1. Server starts, accepts clients, stops cleanly.

### [x] 6.4 Hosting from the app

1. Goal: a player's app runs the server (N-01).
2. Tasks:
   1. Same server package, started in process.
   2. The host joins its own server as a normal client.
3. Done when:
   1. Host and one remote client play over LAN.

### [x] 6.5 Lobby

1. Goal: create and join games.
2. Tasks:
   1. Create a game; join by address.
   2. Seats 2–4 (GS-09); ready check.
   3. Play mode: free-for-all; 2 vs 2 is parked (Q-03).
   4. Card sets; Base Game only for now (CD-02).
   5. Collection check on join (CD-03).
   6. Nicknames from `hello`: checked, made unique, shown per seat (ST-06).
3. Done when:
   1. Use cases: host a game, join a game.

### [x] 6.6 Match setup

1. Goal: simple and detailed setup (GS-10).
2. Tasks:
   1. Simple: defaults only.
   2. Detailed: picking mode, bans, timers.
   3. Picking: random, draft (GS-01).
   4. Ban rounds, snake order, optional ban timer (GS-04).
   5. Host ban list (GS-02); dealing after bans (GS-03).
   6. First player by dice roll (GS-06).
3. Done when:
   1. Every mode has a server test.

### [x] 6.7 Turn flow over the network

1. Goal: responses without stalls.
2. Tasks:
   1. Push allowed actions after every change (N-04).
   2. Resync on request.
   3. Auto-skip from the full allowed list (N-05, N-09).
   4. "Skip all" for visible stack items only (N-06).
   5. Response timer: off or ≥ 60 s, plus animation allowance (N-07).
3. Done when:
   1. Headless games never wait for a player with no options.

### [x] 6.8 Disconnects

1. Goal: lost players can come back (N-08).
2. Tasks:
   1. Game pauses for everyone.
   2. Vote: wait longer or kick.
   3. Kicked seat stays empty.
   4. Reconnect gets the full current view.
3. Done when:
   1. Tests cover drop, reconnect, kick.
   2. Use case: reconnect.

### [x] 6.9 Match records

1. Goal: every match is saved (A-12, RP-01 – RP-09).
2. Tasks:
   1. `pkg/record`: gzip-compressed JSON lines, `.fsrec`.
   2. Header: versions, sets, seats, card text of cards used.
   3. Steps: only applied steps, dice, checksums, all hands.
   4. Autosave on the server, finished and unfinished.
   5. Retention: 30 days dedicated, forever when hosted, configurable.
   6. Players may download only after the match ends.
3. Done when:
   1. Writer and reader round-trip tests pass.
   2. ADR: record format.

### [x] 6.10 Headless integration tests

1. Goal: full games through the real server.
2. Tasks:
   1. Go test clients play random legal moves.
   2. Run 2, 3, 4 player games over WebSocket.
   3. Include drops and reconnects.
3. Done when:
   1. Full games finish and produce valid records.

### [ ] 6.11 Save and continue later

1. Goal: friends finish a game on another day (N-11).
2. Tasks:
   1. Host saves an unfinished game: engine `Save` plus seats.
   2. The save stays with the host or server; players get no copy.
   3. Loading opens a lobby with the saved seats.
   4. Players reclaim their own seats; resume when all are back.
   5. Refuse saves from an incompatible version with a clear message.
3. Done when:
   1. Headless test: save mid-stack, load, finish; checksums match.
   2. Use case: save and continue a game.
