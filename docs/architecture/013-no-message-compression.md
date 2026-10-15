# 13. No WebSocket message compression

* **Status:** Accepted
* **Date:** 2026-10-15
* **Authors:** @HardDie

---

## Context

1. Each `update` carries the full view of the table (ADR 006).
   1. It also has the new events and the allowed intents.
   2. It never carries the game's history.
   3. A resync sends the full view with no events.
2. Steps come at human speed: a few per second at most.
3. Most games run on a LAN or on one computer.
4. Internet play goes through a VPN or a VPS (ADR 006, 9).
5. Measured on 2026-10-15:
   1. 4 players, base set `b2`, 400 random steps.
   2. One player's `update`, JSON-encoded.
   3. Average 2.8 KB raw; largest 3.9 KB.
   4. Average 1.0 KB with deflate at the fastest level.
6. The game client is the browser `WebSocket` (`frontend/src/net.ts`).
   1. Browsers offer `permessage-deflate` by themselves.
   2. Only the server's `Upgrader` decides whether it is used.

## Considered options

1. **No compression**
   1. Won. About 3 KB per step costs nothing on a LAN.
   2. Over the internet it is still milliseconds per step.
2. **`permessage-deflate` (gorilla `EnableCompression`)**
   1. Lost for now. It cuts messages to about a third.
   2. It costs CPU and memory on every send, per client.
   3. gorilla marks it experimental.
3. **Compression inside the protocol (gzip in the envelope)**
   1. Lost. It breaks readable JSON and needs a protocol bump.
4. **Deltas instead of full views**
   1. Out of scope; ADR 006 already parks it.

## Decision

Use option 1.

1. The server's `Upgrader` leaves `EnableCompression` off.
2. Messages stay plain JSON text frames.
3. Revisit when one of these comes true:
   1. Online play on mobile or metered links becomes a real use case.
   2. The view grows a lot, e.g. it starts carrying text or a log.
   3. Updates are sent much more often than once per step.
4. Turning it on later:
   1. Set `EnableCompression: true` on the server's `Upgrader`.
   2. Clients without support get plain frames; no protocol bump.
   3. Measure again before and after.

## Consequences

### Positive

1. No CPU or memory cost per message.
2. Frames stay readable in browser dev tools and proxies.
3. One less experimental feature in the server.

### Negative and risks

1. About three times more bytes than needed on the wire.
2. A slow link pays for that first, if online play grows.

### Neutral

1. Match records on disk stay gzip-compressed (ADR 007).
2. The decision can flip without touching the protocol.
