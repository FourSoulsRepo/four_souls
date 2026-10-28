# Load test

`cmd/server_load` plays whole games against a running dedicated server and reports how much load it takes. It is a WebSocket client like the app: it sends intents and reads updates (ADR 006). It never changes the server.

## Run it

```sh
go run ./cmd/server_load -url ws://game.example.com:4774
```

Run it against a server you run yourself. The games it opens are real games on that server; they end on their own, or with the server's idle and length limits.

## Stages

Each game seats 4 bots (`-players`) at one table and plays it to the end with random legal moves.

The games run in stages:

* stage 1: 2 games at once (`-start`);
* stage 2: 4 games at once;
* stage 3: 8 games at once, and so on, up to `-max-games` (default 64).

A stage passes when:

* every game ends, and no game fails or times out (`-game-timeout`, default 30 minutes);
* the slowest response is not over `-max-response` (default 1s).

The first stage that fails ends the run. The summary names the highest load that passed.

With `-stop-at 30`, the stage stops its other games once 30% of its games have ended (rounded up, at least one). Those games are counted as stopped, not failed. A stage that stops early measures how the server handles the load before most games end, not whole games.

## Flags

| Flag | Default | What it does |
|------|---------|--------------|
| `-url` | required | `ws://host:port`; the path `/ws` is added |
| `-players` | 4 | players per game, 2 to 4 |
| `-start` | 2 | games at once in the first stage |
| `-max-games` | 64 | most games at once |
| `-sets` | b2 | card sets, comma-separated |
| `-game-timeout` | 30m | a game that has not ended by then fails |
| `-max-response` | 1s | a stage fails when a response takes longer |
| `-seed` | clock | seed for the bots' moves |
| `-stop-at` | 0 | percent of a stage's games that, once ended, stop the others; 0 waits for all |

## What it measures

* **response**: from an intent sent to the next update that seat gets. Min, median, p95 and max.
* **ping**: a WebSocket ping and its pong, once a second per bot.
* **handshake**: connect, hello and welcome.
* **updates**: how many updates the bots received, their size and the rate per second.
* **errors**: error replies from the server by code; `refused` is an intent the engine did not accept.
* **notices**: a notice that ends a game (lost or idle) fails it.
* **agreement**: every seat must see the same last step and winners, or the game fails.

Progress is printed every 2 seconds. With up to 8 games it shows one part per game: its turn and step, then `done` or `FAILED`. With more, it shows counts of running, done and failed games.

## Things to know

* A bot ends its turn after 500 steps of it, as a player would. A game that still does not end fails, and the error names its turn and step.
* A failed game stays on the server until the server's limits end it.
* The bots play random legal moves, so one game can take far longer than a real game of the same length.
