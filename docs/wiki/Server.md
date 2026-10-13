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
* the player waits for their turn: until the seat before them declares the end of its turn, dies, or reaches its end phase; the end window then asks them (N-16);
* the response timer ran out, or the seat was kicked: pass or end the turn, the first cards for a discard, the first option of a choice (N-07, N-08).

## LAN discovery

`internal/lan` (ADR 009). A joining app sends a UDP probe to port 4774: to 255.255.255.255, to each interface's broadcast address, and to 127.0.0.1 (several copies on one computer, UC-16). A server started with `Config.Discovery` answers the sender by unicast with its waiting games; with none waiting it stays silent. Only the app's `Host` sets `Discovery`; it has no JSON name, so `cmd/server` never answers. Answers carry a random server ID, so the same server reached over loopback and the LAN shows once, with its LAN address.

## Leak audit (A-16)

Step 8.12 checked every message the server sends for hidden cards. Hidden means: in a deck, in another player's hand, or seen only by another player.

| Message | What it carries | Verdict |
|---|---|---|
| `welcome` | versions, your token, your seat | Tokens never reach other clients; `Seat` has no token. |
| `games` | host nickname, seats, set codes; saves with nicknames | No cards. |
| `table` | seats, set codes, options | Host bans are public by design (GS-02). |
| `setup` | pool, bans, who picked | `Offers` go only to their own seat. |
| `update` → `view` | `engine.View` for the seat | Other hands are sizes only; deck contents and order never leave; choice labels and `Options` only for the chooser; the top treasure only for a player who may peek. |
| `update` → `events` | `engine.FilterEvents` | Private events keep `Card` and `Object` only for `Player` and `SeenBy`. Every non-private event with a card was checked: each moves a card from or to a public place (stack, slot, discard, play), or is a public reveal. |
| `update` → `allowed` | this seat's intents | Only the seat's own objects and choice indexes. |
| `update` → `aims` | first targets of allowed plays and uses, and a `reason` saying what they may be | Targets are public objects, players and stack items; a "your hand card" target goes only to its own seat. |
| `error` | code, message, rule | Engine refusals name only cards in play; an activation checks the zone before it names the card. |
| `/record` | the match record (all hidden info) | Only to the game's players, only after the game is over (RP-05); `TestDownloadRecord`. |

Rules that keep it so:

* A random pick over hidden cards (Boomerang, Ring of Flies) leaves its option labels empty; `EvRandomPick` sends the label to everyone.
* Questions (`Prompt.Text`) are fixed card text, so they are public.
* The collection check compares whole set codes, never the dealt decks (CD-03).
* Card text for outdated clients (CD-09) is not built yet; when it is, it must answer only for cards the client may see.
* The fuzz check `checkNoLeaks` fails when any view names an object that sits in a deck, or a stack item from a deck.

Found and fixed: an `activate` intent with an unknown object ID made the engine panic, which stopped the whole server. It is now refused (`TestUnknownObjectsAreRefused`).

## Tests

* `hub_test.go`, `setup_test.go`, `flow_test.go`, `pause_test.go`, `records_test.go`, `saves_test.go`: in-memory connections.
* `integration_test.go`: whole games of bots over real WebSockets, a drop and a reconnect, a save and a load; every record must replay in the engine. Under `-race` only the 2-player game runs.
* `internal/protocol/testdata/golden.jsonl` pins the JSON of every message: rename or remove a field, bump `protocol.Version`; then `go test ./internal/protocol -update`.
* Timers count `Hub.second`, so tests make a second one millisecond.
