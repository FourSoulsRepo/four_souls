# Bots

Bots fill free seats of a network game (ADR 030, step 21). This page is for developers: where the code lives, the levels, and the arena that measures them.

## Layout

| Package | What it does |
|---|---|
| `internal/bot` | Picks a seat's move: `Choose(ctx, Input)`. Imports the engine and the standard library only. |
| `internal/world` | Turns the real game into a seat's worlds: `world.For(g, seat, known, seed)`. |
| `internal/arena` | Plays whole games between levels and counts them. |
| `cmd/bot_arena` | Prints the arena's results; `make bots` runs it. |

### A bot never holds the real game

`bot.Input` holds the seat, the level, the bot's own random source, a budget and `World`. It holds no game. `World(i)` returns a new game as the seat may imagine it: `Game.Imagine` redeals every card the seat cannot see (see "Imagined game" in [Rules engine](Rules-engine)). So a bot cannot read another hand or a deck order, even by mistake.

`world.For` is the one place that makes worlds. It copies the game once (`Game.Clone`), at once, and returns a function: world `i` is that copy imagined by the seat with the seed `world.Seed(seed, i)` and the cards the seat remembers (`known`; empty until the memory of step 21.15, `world.Memory`, which lives in `internal/world` because it names cards). The arena and the server's bot driver (step 21.8) both call it, so a bot gets the same thing in both. The copy is made on the caller's goroutine; the game may change after `For` returns. One world function is for one goroutine at a time.

`internal/world` sits outside `internal/bot` on purpose: it touches the real game, the bot never does. Two checks hold this:

* the `bot-honest` rule of `.golangci.yml` (depguard): `internal/bot` imports neither `internal/server`, nor `internal/world`, nor a card set;
* `TestBotRules` (`internal/bot/rules_test.go`) reads the bot's source: only the standard library and the engine are imported; no `CardRef`, no card field (`Card`, `Ref`, `CardOf`), no card labels; no string that is a card's ref, and strings only in constants and error messages; no `NewGame`, `Load` or `Imagine`; `Game.Clone` only in `fork`, which copies a world (step 21.5); no package variable holds a game; `Input` holds no game. `TestBotRulesCatch` checks that these rules catch a piece of bad code.

## Choose

`Choose` asks `World(0)` for `Allowed(seat)`. One allowed intent is the answer at once (a discard or an order gets its cards first). Else the level picks. The answer is always one of `Allowed`, a discard cut to the asked count and an order of the same cards. An error (no world, nothing allowed, an unknown level) is not a move: the caller sends its forced answer (N-07).

Candidates (Decision 5): every allowed intent; a discard as up to 8 different sets of the asked count, picked at random; an order as the order now and up to 5 random orders.

## Levels

* `random`: one of the candidates at random. For tests and the arena.
* `policy`: the playout policy alone. For tests and the arena.
* `standard` and `hard` (the server's levels) come with steps 21.5 and 21.6.

### Playout policy

The policy (`internal/bot/policy.go`, Decision 6) weighs each allowed intent and draws one by weight. It copies nothing, so later levels can use it for every seat of every playout. It reads kinds and public numbers only: coins, the shop price, HP, ATK, evasion, souls. Every weight is a constant at the top of the file.

* Outside its own action phase a seat mostly passes, so the stack resolves.
* In its own action phase: buy when it has the coins, then attack, play loot, use abilities (less after each use this turn), and end the turn when nothing else draws.
* Attack odds: the chance to kill a monster before it kills the attacker, from the attacker's ATK and HP and the monster's HP, evasion and ATK. A roll of 1 counts as a miss, as abilities may raise the evasion during the fight. The monster deck's top is judged as an average monster (HP 3, evasion 4, ATK 1).
* Attack target: weight from the odds, counted twice, and the target's souls.
* A choice that cancels a play (a target, a mode, a copied ability) gets no weight, so a play is never taken back and tried again for ever.
* A curse goes to another seat, the more souls the likelier.

## Arena

`cmd/bot_arena` plays whole seeded games with every set: bonus souls on, the room deck in odd seeds, deck ratios in every third. Each `-match` is one line-up of levels; game `k` plays seed `-seed`+k and rotates the line-up by `k` seats, so each level plays from each seat. Games run side by side on every core; the results do not depend on it.

```sh
make bots                     # BOT_GAMES=50 per default line-up
go run ./cmd/bot_arena -games 20 -match policy,random -match policy,policy,policy
go run ./cmd/bot_arena -games 20 -cap 400 -match policy,policy,policy,policy
```

It prints, per line-up: games, games with a winner within the turn cap, wins by level against chance, median and max game turns, intents per turn, the share of passes, time per move, player deaths per monster killed; and for every game cut by the cap, its seed, turn, souls by seat and the kinds of its last 300 intents. It exits 1, in any line-up, when a bot intent is refused or a loop guard ends a game. The finish gate, 95% of games with a winner within the cap (Decision 18.4), is only for line-ups of the levels players meet (standard from 21.5, hard from 21.6) with no random seat: `finishGated` in `cmd/bot_arena`. Random and policy line-ups print their share for reference.

Each step the arena asks the waiting seat's bot, with `world.For` on the game for that seat. A refused intent is listed and the arena sends the forced answer: the first allowed intent.

The cap, `arena.TurnCap`, is 150 game turns (`Turn.Number`). Two loop guards end a game as a failed gate, printed with its seed: one game turn over 1000 applied intents (`MaxTurnIntents`), a game over 30000 (`MaxGameIntents`). A turn is measured, not intents, because about 80% of intents are passes: every stack item waits for a pass from each seat, so a 4-player turn takes about 86 intents and a 2-player turn about 46.

Tests: `TestNoRefusedIntent` (`internal/arena`) plays 30 whole policy games, every set, 2 to 4 players, and fails on any refused intent.
