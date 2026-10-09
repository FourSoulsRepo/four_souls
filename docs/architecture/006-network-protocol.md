# 6. Network protocol

* **Status:** Proposed
* **Date:** 2026-10-10
* **Authors:** @HardDie

---

## Context

1. The server is authoritative; clients only send intents (N-02).
2. Clients do not depend on the engine version (N-03).
3. Each seat gets its allowed actions after every change (N-04).
4. Clients ship all cards; the wire carries card refs only (CD-01).
5. Hidden information never leaves the server (A-07).
6. Players reconnect into the same seat (N-08).
7. The same server runs dedicated or inside a player's app (N-01).
8. A browser client comes later (step 13); spectators too (step 9).

## Considered options

1. **WebSocket, JSON text messages**
   1. Won. Works from Go and the browser; readable when debugging.
2. **gRPC or a binary format**
   1. Lost. Needs a proxy in the browser; harder to debug; more deps.
3. **gorilla/websocket**
   1. Won. Already in the dependency tree through Wails; BSD-2.
4. **coder/websocket**
   1. Lost. A new dependency for the same result.
5. **Full view after every change**
   1. Won. A lost message can never leave a client wrong; small on LAN.
6. **Deltas only**
   1. Lost for now. Smaller, but needs ordering and replay logic.

## Decision

Use options 1, 3 and 5.

1. Transport
   1. One WebSocket per client at `ws://host:port/ws`.
   2. Text frames, one JSON message each.
   3. Client messages are at most 64 KiB.
2. Envelope
   1. Every message: `{"type": "...", "id": 7, "data": {...}}`.
   2. `id` is set by the client on requests; replies carry it back.
   3. Types live in `internal/protocol`.
3. Handshake
   1. The client sends `hello`: protocol version, app version, name, role, card sets.
   2. On reconnect `hello` also carries the session token.
   3. The server answers `welcome`: versions, a new or kept token, the seat.
   4. A different protocol version gets `error` `client_outdated` or `server_outdated`, then close.
4. Playing
   1. The client sends `intent`: an engine intent; the server fills in the seat.
   2. The server answers a refused intent with `error` `refused`, with the rule and reason.
   3. After every applied step each client gets `update`.
   4. `update` holds a step number, the new events, the full view and the allowed intents.
   5. Views and events are filtered per seat by the engine (A-07).
   6. A client that misses a step number sends `resync` and gets a fresh `update`.
5. Compatibility
   1. The view, event and intent JSON is part of the protocol.
   2. Enum values are numbers and only ever appended, never reordered.
   3. A golden test pins the JSON of each message.
   4. A change to it bumps the protocol version.
6. Cards
   1. Messages carry card refs (`the_d6`, `the_d6@2`), never card text.
   2. A client shows an unknown ref as a blank card with the ref (SP-05).
   3. `hello` lists the client's card sets for the collection check (CD-03).
7. Trust
   1. LAN first: no accounts; the token only proves "same seat as before".
   2. Tokens are 128 random bits.
8. Later messages use the same envelope.
   1. Lobby (6.5), match setup (6.6), pause and votes (6.8), saves (6.11).
   2. Spectators send `hello` with role `spectator` or `judge`.

```json
{"type": "hello", "id": 1, "data": {"protocol": 1, "app": "0.1.0", "name": "Oleg", "role": "player", "sets": ["b2"]}}
{"type": "welcome", "id": 1, "data": {"protocol": 1, "app": "0.1.0", "engine": "0.1.0", "token": "…", "seat": 0}}
{"type": "intent", "id": 2, "data": {"intent": {"kind": 3, "objects": [41]}}}
{"type": "update", "data": {"step": 18, "events": […], "view": {…}, "allowed": […]}}
{"type": "error", "id": 2, "data": {"code": "refused", "rule": "R-CARD-08", "message": "no loot play available"}}
```

## Consequences

### Positive

1. Any language with WebSocket and JSON can be a client.
2. Debugging is reading JSON.
3. Resync makes lost or reordered messages harmless.

### Negative and risks

1. Full views cost bandwidth; fine on LAN, may need deltas online.
2. Engine JSON is now a public format; renames need care.

### Neutral

1. No accounts or encryption until online play (future).
