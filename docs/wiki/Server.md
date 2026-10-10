# Server

Notes on `internal/server` and `internal/protocol` as they are now. The decisions are in ADR 006 (protocol) and ADR 007 (records); user guides are [Playing over the internet](Playing-over-the-internet) and [Dedicated server on a VPS](Dedicated-server-on-a-VPS).

## Shape

```
WebSocket (ws.go) ──Conn──▶ Client ──▶ Hub (lobby goroutine)
                                  └──▶ Room (one goroutine per game)
```

* **Conn** is all a transport implements: `Send` and `Close`. The WebSocket adapter adds a send queue and a writer goroutine per client; a client that falls behind is disconnected and resyncs on reconnect. Another transport (WebRTC, a relay) only adds a Conn.
* **Hub** owns the lobby: hello, the game list, tables, ready checks, saved games. Only its goroutine touches that state; everything reaches it as a request on a channel.
* **Room** owns one game: the ban and pick phase, then the engine game. Only its goroutine touches the engine.
* **Client.Receive** sends a message to the hub until the client's game starts, then straight to its room.
* **Server** (`server.go`) is the HTTP server with `/ws` and `/record`; `cmd/server` and the app's Host binding both start it.

## A game's life

1. Hello: protocol version, nickname, card sets; the server keeps or makes a token.
2. Lobby: create (seats, sets, options) or join; ready. A full, ready table becomes a Room.
3. Setup: host bans out; ban rounds in snake order with an optional timer; random or draft picks.
4. Play: each intent is applied for the client's own seat; then `after` runs automatic passes, restarts the response timer, sends every client its update and records the steps.
5. Disconnects pause the game; a vote can kick; a token brings a player back.
6. The host may save: the game stops and can be loaded again from the lobby.
7. At the end the record is closed and its players may download it.

## Automatic answers

The server answers for a player only in these cases, never choosing among real options otherwise:

* the player can only pass (auto-skip, N-05);
* the player turned on "skip all" for the stack they saw (N-06);
* the response timer ran out, or the seat was kicked: pass or end the turn, the first cards for a discard, the first option of a choice (N-07, N-08).

## Tests

* `hub_test.go`, `setup_test.go`, `flow_test.go`, `pause_test.go`, `records_test.go`, `saves_test.go`: in-memory connections.
* `integration_test.go`: whole games of bots over real WebSockets, a drop and a reconnect, a save and a load; every record must replay in the engine. Under `-race` only the 2-player game runs.
* `internal/protocol/testdata/golden.jsonl` pins the JSON of every message: rename or remove a field, bump `protocol.Version`; then `go test ./internal/protocol -update`.
* Timers count `Hub.second`, so tests make a second one millisecond.
