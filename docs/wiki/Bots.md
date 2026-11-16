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
* `standard`: tries each candidate on one imagined world and plays the best score (below). The server's lower level.
* `hard`: plays each candidate out on many imagined worlds to the end of the turn and plays the best average (below). The server's upper level.

### Playout policy

The policy (`internal/bot/policy.go`, Decision 6) weighs each allowed intent and draws one by weight. It copies nothing, so later levels can use it for every seat of every playout. It reads kinds and public numbers only: coins, the shop price, HP, ATK, evasion, souls. Every weight is a constant at the top of the file.

* Outside its own action phase a seat mostly passes, so the stack resolves.
* In its own action phase: buy when it has the coins, then attack, play loot, use abilities (less after each use this turn), and end the turn when nothing else draws.
* Attack odds: the chance to kill a monster before it kills the attacker, from the attacker's ATK and HP and the monster's HP, evasion and ATK. A roll of 1 counts as a miss, as abilities may raise the evasion during the fight. The monster deck's top is judged as an average monster (HP 3, evasion 4, ATK 1).
* Attack target: weight from the odds, counted twice, and the target's souls.
* A choice that cancels a play (a target, a mode, a copied ability) gets no weight, so a play is never taken back and tried again for ever.
* A curse goes to another seat, the more souls the likelier.

### Score

`Score(g, p)` (`internal/bot/score.go`, Decision 7) is a table from seat `p`'s side: a won game is the top value, a lost one the bottom. Else it is the seat's own sum less the opponents' average. A sum counts souls, cents, loot cards in hand, items, curses, HP left and being alive this turn; a card's base value comes from `cardValue`. Every weight is a constant in that file. A soul (100) outweighs every other part of a table; HP (0.5 each) and being alive (3) weigh little, as both heal at the end of the turn; an item (12) is worth a little more than the 10¢ it costs; a loot card in hand (0.75) is worth less than a cent, so a card that gives 1¢ is played and one that does nothing now is kept.

A settled copy is scored by `settledScore`: its `Score`, plus `digValue` for each monster the seat killed in the settle (its rewards went to the seat). A kill lets in the next card of the monster deck, which may give a soul: `digValue` is the soul weight times the share of soul cards (card property `Soul`) in the world's monster deck, times `digShare` (0.5). Without it a kill of a monster with no soul rarely beat ending the turn.

### Standard

The standard level (`internal/bot/standard.go`, Decision 8) asks for one world. For each candidate it copies the world (`fork`, the only place the bot copies a game), applies the candidate and settles the copy: every seat passes, the policy answers the other prompts, and the bot's own choices get the policy's top option. A settle runs to the end of what the candidate started: the stack empty with no attack or purchase going on, the end of the turn, or the game's end; at most 150 steps. So an attack runs until the monster dies, the attacker dies and pays its penalty, or the attack is cancelled. The settled copy is scored.

A candidate whose settle threw a die is settled 6 times: on the world's own dice, then on copies whose dice are reseeded. Each settle's policy choices come from a seed of its own, and every candidate of a move meets the same seeds, so scores differ by the move, not by luck.

An action must beat doing nothing (a pass, or an end of turn in the bot's own action phase): its gain over doing nothing, settle by settle, counts less a quarter coin (`wAct`) and less a share of its standard error: one on an empty stack; two on a stack, where the move is judged through the dice of every item below it. A move that gains nothing, or gains only by the luck of its dice, is not made; this also ends the long chains of paid abilities a seat could stack while it has coins. Options of a choice that take a play back are left out, as in the policy.

Some tables no bot level unlocks, by the cards alone: The Harbingers (`r`) makes the active player attack it every turn and heals when it would die (four times, then it flips); with ATK 1 characters against its DC 4 or 5, every turn's attack goes there at a few percent and the game can reach the turn cap (×4 seed 15, 21.5). The card works as its text says.

It plays the best score; on equal scores a pass or an end of turn first. One move in ten it picks by chance among the best three, weighted 3:2:1.

### Hard

The hard level (`internal/bot/hard.go`, Decision 9) lists the candidates as standard does, doing nothing first. Rounds of playouts halve them: each round gives every candidate still in an even share of the budget, world by world. A world is `World(i)`, a new imagined deal; every candidate still in plays on it with the same policy seed, then it is dropped, so one world and one copy live at a time. A playout copies the world (`fork`), applies the candidate and plays on to the end of the turn, at most 120 steps; every seat moves by the policy, so the others may respond, and the bot's own choices get the policy's top option. Each playout is scored by `settledScore`, the function standard settles use. After a round the better half of the moves stays, by `judge`; doing nothing stays to the last round, against the one move left, and a move must beat it as in standard. In that noise a gain counts at most four souls (`noiseClip`), so one world where the game ends does not hide the rest. There are no mistakes on purpose.

Hard's score (Decision 7.6) weighs the opponents by strength, not as an average: each opponent's weight is `exp(sum / 100)` over its full sum, so one soul ahead counts e times as much; equal opponents count as their average. Its answers have a holding value: the loot cards and abilities `Allowed` shows it on another seat's turn (on its own turn, the abilities of its cards but rooms'). Each answer still held at the end of a playout (in hand; in play and charged, or recharging as its turn starts) counts its base value (`cardValue`, the one function that values a loot card in hand and an item in play) times 1 + `wHold` × (1 + `wHoldLead` × lead): an item answer 2.4 to 7.2 more than a spent one, more as one opponent pulls away. So it keeps Godhead when a weak seat kills a monster with no soul, and spends it when the strong seat kills one with a soul.

`Budget{Playouts, Time}`: whichever ends first. `Playouts` 0 means no playout cap, `Time` 0 no time limit; both 0 think 1.5 s. With a time budget the first world measures a playout and the rounds share out what the time leaves. The search checks `ctx` before each playout, yields the core after each (`runtime.Gosched`) and, when stopped, plays the best so far. A playout takes about 4 ms on one core, most of it the engine's `checkMaxHP`; a move of 150 playouts about 0.4 s and 150 MB of garbage.

## Arena

`cmd/bot_arena` plays whole seeded games with every set: bonus souls on, the room deck in odd seeds, deck ratios in every third. Each `-match` is one line-up of levels; game `k` plays seed `-seed`+k and rotates the line-up by `k` seats, so each level plays from each seat. Games run side by side on every core; the results do not depend on it.

```sh
make bots                     # BOT_GAMES=50 per default line-up: standard against random, hard against standard, one hard among three standard, standard×2–4, policy×2–4
go run ./cmd/bot_arena -games 20 -playouts 500 -match hard,standard
go run ./cmd/bot_arena -games 20 -match policy,random -match policy,policy,policy
go run ./cmd/bot_arena -games 20 -cap 400 -match policy,policy,policy,policy
```

It prints, per line-up: games, games with a winner within the turn cap, wins by level against chance, median and max game turns, intents per turn, the share of passes, time per move, player deaths per monster killed; and for every game cut by the cap, its seed, turn, souls by seat and the kinds of its last 300 intents. It exits 1, in any line-up, when a bot intent is refused or a loop guard ends a game. The finish gate, 95% of games with a winner within the cap (Decision 18.4), is only for line-ups of the levels players meet (standard from 21.5, hard from 21.6) with no random seat: `finishGated` in `cmd/bot_arena`. Random and policy line-ups print their share for reference. The win gates are `winGates` there: standard wins 70% of 2-player games against random, hard 60% of 2-player games against standard, and one hard among three standard wins over 30% (Decision 18.5 – 18.7). Hard thinks `-playouts` playouts a move (`arena.DefaultPlayouts`, 150), never a time, so the results repeat.

Each step the arena asks the waiting seat's bot, with `world.For` on the game for that seat. A refused intent is listed and the arena sends the forced answer: the first allowed intent.

The cap, `arena.TurnCap`, is 150 game turns (`Turn.Number`). Two loop guards end a game as a failed gate, printed with its seed: one game turn over 1000 applied intents other than passes (`MaxTurnIntents`), a game over 30000 applied intents (`MaxGameIntents`). The turn guard counts no passes: a long turn of real plays in a 4-player game (a paid ability used 40 times, each a stack item four seats pass) passed 1000 intents with no loop (21.5). A turn is measured, not intents, because about 80% of intents are passes: every stack item waits for a pass from each seat, so a 4-player turn takes about 86 intents and a 2-player turn about 46.

Tests: `TestNoRefusedIntent` (`internal/arena`) plays 30 whole policy games, every set, 2 to 4 players, and fails on any refused intent.
