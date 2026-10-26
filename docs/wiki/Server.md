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
* **Server** (`server.go`) is the HTTP server with `/ws`, `/record`, `/healthz`, `/seat` and `/status`; `cmd/server` and the app's Host binding both start it.
* **`/healthz`** (`health.go`, 8.19) answers 200 while the hub goroutine runs a request within two hub seconds, else 503. Docker's healthcheck asks it; it tells nothing about games.
* **`/seat?token=…`** (`seat.go`, 15.99) answers a `protocol.SeatInfo`: whether the token holds a seat at a table or in a game that has not finished, with the game's id and host. A restarted app asks it before the menu offers Rejoin. A malformed token is never looked up. It sends `Access-Control-Allow-Origin: *`: only the token's holder learns anything, and the e2e tests' fake bridge asks it from a browser page.
* **`/status`** (`status.go`, 15.100) answers a `protocol.ServerStatus`: `open`, the tables with a free seat; started games are not counted. The Online screen asks every Recent server when it opens, to show which are up. It tells no more than the lobby's game list and sends `Access-Control-Allow-Origin: *` like `/seat`. Both wait for the hub through `askHub`, which answers 503 when the lobby is stuck or stopping.

## A game's life

1. Hello: protocol version, nickname, card sets, dice look; the server keeps or makes a token. The seat keeps the look (`CleanDice`: a bad ID is the default), even while its player is away, and every table and update lists it (V-19).
2. Lobby: create (seats, sets, options) or join; ready. A full, ready table becomes a Room.
3. Setup: host bans out; ban rounds in snake order with an optional timer; random or draft picks. Deck ratios are on unless the host turns them off; `checkOptions` turns them off when the sets are too small (GS-07).
4. Play starts with the first-player roll; the first update carries only `first_player` and `game_started` of the setup events, for Cain (`startEvents`, 15.43). Each intent is applied for the client's own seat; then `after` runs automatic passes, restarts the response timer, sends every client its update and records the steps.
5. Disconnects pause the game; a vote can kick; a token brings a player back. A finished game never takes a token back: the hello goes to the lobby, so the player can start a new game (19.32).
6. The host may save: the game stops and can be loaded again from the lobby.
7. At the end the record is closed and its players may download it.
8. A dedicated server that stops keeps every table and game in `<saves>/restart/`; the next start brings them back (`restart.go`, ADR 019). A game from another engine must load and save back to the same bytes (`engine.Resume`); else it ends and its players hear why when they return.
9. Idle limits end a game nobody is connected to, or one past its longest time (`limits.go`, 8.25). Each room checks once a minute on the hub's clock; tests pass a fake one.

## Automatic answers

The server answers for a player only in these cases, never choosing among real options otherwise:

* the player can only pass (auto-skip, N-05);
* the player waits for their turn: until the seat before them declares the end of its turn, dies, or reaches its end phase; the end window then asks them (N-16). A stack item of another player that targets them, or their character, items, souls, curses or hand cards, or deals them damage, ends the wait at once, if it came after the press (15.84). So does one whose ability may take from them with no target, e.g. XV. The Devil, which picks the item it steals as it resolves (`Ability.Takes`, 15.85). The engine decides what threatens a player (`Game.Threatens` in `threat.go`); the server only keeps the stack ID of the press (`targeted` in `flow.go`, 15.86);
* the response timer ran out, or the seat was kicked: pass or end the turn, the first cards for a discard, the first option of a choice, a throw for a roll (N-07, N-08, 15.60).

"Wait for my turn" only passes: a roll always waits for its own player's throw (ADR 014).

## Presence

`presence.go` (15.35, N-15, ADR 016). A seated player sends `presence` while aiming or hovering: `from` (what aims), `at` (what the arrow or pointer is over), `point` (the pointer on the board frame, 0 to 1). The room fills in `seat` and sends it to every other client of the game, spectators too once step 9 seats them; the sender gets nothing back. Each one replaces the sender's last, so a lost or late one does no harm. When a player's connection drops, the others get an empty presence of that seat, so its arrow goes.

The engine never sees it and no record or save keeps it. It is never answered: a message that does not decode, comes before the game starts, or comes too fast is dropped. `pace` allows `protocol.MaxPresence` a second per seat, with a burst of as many.

An operator may turn it off (`Config.Presence`, `-presence=false`): the room then drops every presence and relays none, and `welcome.presence` is unset so clients send none. It is on by default, and always for the app's own host.

A spot names exactly one thing: a card every player can see (`object`), a slot of the sender's own hand (`hand`, from 1), a player, a deck, or a stack item. The room checks it against the table as a spectator sees it, built once per step (`publicTable`). Anything else becomes no spot, so a hand card's object never leaves the server (A-16).

## LAN discovery

`internal/lan` (ADR 009). A joining app sends a UDP probe to port 4774: to 255.255.255.255, to each interface's broadcast address, and to 127.0.0.1 (several copies on one computer, UC-16). A server started with `Config.Discovery` answers the sender by unicast with its games (`lanGames`): every table, full or not, and every started game that has not finished, each with its id and `started` (15.97); with none it stays silent. A started game shows so a player whose app closed finds it again. Only the app's `Host` sets `Discovery`; it has no JSON name, so `cmd/server` never answers. Answers carry a random server ID, so the same server reached over loopback and the LAN shows once, with its LAN address.

## Who opens games

Any client of a dedicated server may create a game or load a saved one. A server started with `Config.HostOnly` lets only clients on its own computer do that (loopback address, `wsConn.Local`); everyone else gets `not_host` and only joins (N-17). Only the app's `Host` sets it, so friends of an app's host never open games on that app. `welcome.create` tells each client which case it is; the app shows **Create a game** only then.

## Server name

An operator may name the server (`Config.Name`, `-name`, N-18). `Config.Check` refuses a name that `protocol.CleanServerName` rejects: empty after trimming, over `protocol.MaxServerName` (40) characters, or with control characters. Every `welcome` carries it in `name`; the app shows it instead of the address. Each match record's header keeps it in `server` too. The app's own host has none.

## Leak audit (A-16)

Step 8.12 checked every message the server sends for hidden cards. Hidden means: in a deck, in another player's hand, or seen only by another player.

| Message | What it carries | Verdict |
|---|---|---|
| `welcome` | versions, your token, your seat | Tokens never reach other clients; `Seat` has no token. |
| `games` | host nickname, seats, set codes; saves with nicknames | No cards. |
| `table` | seats, set codes, options | Host bans are public by design (GS-02, GS-13). |
| `setup` | pool, bans, who picked | `Offers` go only to their own seat. |
| `update` → `view` | `engine.View` for the seat | Other hands are sizes only; deck contents and order never leave; choice labels, `Options` and `Activating` only for the chooser; the top treasure only for a player who may peek. |
| `update` → `events` | `engine.FilterEvents` | Private events keep `Card` and `Object` only for `Player` and `SeenBy`. Every non-private event with a card was checked: each moves a card from or to a public place (stack, slot, discard, play), or is a public reveal. |
| `update` → `allowed` | this seat's intents | Only the seat's own objects and choice indexes. |
| `update` → `aims` | first targets of allowed plays and uses, and a `reason` saying what they may be | Targets are public objects, players and stack items; a "your hand card" target goes only to its own seat. |
| `update` → `blocked` | why the seat can't use a card or turn action now: kind, object, `reason` (15.29) | Only the seat's own hand cards, its own cards in play, and rooms. Reasons are fixed words; none names a card. |
| `presence` | another seat's aim and hover (15.35) | Spots are public cards, players, decks, stack items, or the sender's hand slot; a hidden object is dropped by the server (`TestPresenceHidesHandCards`). |
| `error` | code, message, rule | Engine refusals name only cards in play; an activation checks the zone before it names the card. |
| `/seat` | whether a token holds a live seat, its game and host (15.99) | Only the token's holder knows the token; any other token is "not live". |
| `/status` | how many tables have a free seat (15.100) | A number only; any lobby client sees the games themselves. |
| `/record` | the match record (all hidden info) | Only to the game's players, only after the game is over (RP-05); `TestDownloadRecord`. |

Rules that keep it so:

* A random pick over hidden cards (Boomerang, Ring of Flies) leaves its option labels empty; `EvRandomPick` sends the label to everyone.
* Questions (`Prompt.Text`) are fixed card text, so they are public.
* The collection check compares whole set codes, never the dealt decks (CD-03).
* Card text for outdated clients (CD-09) is not built yet; when it is, it must answer only for cards the client may see.
* The fuzz check `checkNoLeaks` fails when any view names an object that sits in a deck, or a stack item from a deck.

Found and fixed: an `activate` intent with an unknown object ID made the engine panic, which stopped the whole server. It is now refused (`TestUnknownObjectsAreRefused`).

## Rollback

`rollback.go` (ADR 020, 17.3). A game whose `Options.Rollback` is set keeps a base state (`Room.base`, the engine's `Save`) and the live steps since it (`Room.live`: step number, intent and checksum; no events). The base is the game's start, a loaded save, or, until 17.5, the state a restart kept (`setBase`). A game without the mode keeps nothing.

A seated player sends `rollback` with a step number. In Godmode (`god`) any player may, the author too, for any live step, until the game is over; not while the game is paused. The room rebuilds the game with `engine.ReplayFrom(base, steps before the target)`; a checksum that differs refuses the rollback and the game goes on. The steps from the target on leave the live list. Refusals are error `rollback`.

The rollback takes the next step number itself (`Room.step` and `Room.steps` both count it), so step numbers never repeat and the server's numbers stay the record's. Every client gets one update with `rollback` (who, the target, how many live steps it undid) and no events; the app jumps to it. "Wait for my turn" is turned off for everyone, since it would pass the restored prompt away. Passes the server makes for a player who can only pass come after, as new steps in their own update. Each update of such a game carries `rollback_mode` and `rollback_from`, the first step a rollback may target.

The record keeps every step; a `rollback` line follows (`record.Writer.Rollback`). `record.Record.Live` drops the undone steps, so a replay from the setup with them rebuilds the game (`TestGodmodeRollback`). The game list (`GameInfo.rollback`), the table's options, saves and restart files carry the mode. Soft (votes and barriers) is 17.4; `checkOptions` refuses it until then.

## Tests

* `hub_test.go`, `setup_test.go`, `flow_test.go`, `pause_test.go`, `records_test.go`, `saves_test.go`, `presence_test.go`, `rollback_test.go`: in-memory connections.
* `integration_test.go`: whole games of bots over real WebSockets, a drop and a reconnect, a save and a load; every record must replay in the engine. Under `-race` only the 2-player game runs.
* `internal/protocol/testdata/golden.jsonl` pins the JSON of every message: rename or remove a field, bump `protocol.Version`; then `go test ./internal/protocol -update`.
* Timers count `Hub.second`, so tests make a second one millisecond.
