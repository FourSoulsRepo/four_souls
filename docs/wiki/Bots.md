# Bots

Bots fill free seats of a network game (ADR 030, step 21). This page is for developers: where the code lives, the levels, and the arena that measures them.

## Layout

| Package | What it does |
|---|---|
| `internal/bot` | Picks a seat's move: `Choose(ctx, Input)`. Imports the engine and the standard library only. |
| `internal/world` | Turns the real game into a seat's worlds: `world.For(g, seat, known, seed)`; a hard seat's memory, `world.Memory`. |
| `internal/arena` | Plays whole games between levels and counts them. |
| `cmd/bot_arena` | Prints the arena's results; `make bots` runs it. |
| `cmd/bot_values` | Measures each card's value for the score into `internal/bot/values.json`; `make bot-values` runs it. |

### A bot never holds the real game

`bot.Input` holds the seat, the level, the bot's own random source, a budget and `World`. It holds no game. `World(i)` returns a new game as the seat may imagine it: `Game.Imagine` redeals every card the seat cannot see (see "Imagined game" in [Rules engine](Rules-engine)). So a bot cannot read another hand or a deck order, even by mistake.

`world.For` is the one place that makes worlds. It copies the game once (`Game.Clone`), at once, and returns a function: world `i` is that copy imagined by the seat with the seed `world.Seed(seed, i)` and the cards the seat remembers (`known`: a hard seat's `world.Memory`, below, which lives in `internal/world` because it names cards; empty for the other levels). The arena and the server's bot driver (step 21.8) both call it, so a bot gets the same thing in both. The copy is made on the caller's goroutine; the game may change after `For` returns. One world function is for one goroutine at a time.

`internal/world` sits outside `internal/bot` on purpose: it touches the real game, the bot never does. Two checks hold this:

* the `bot-honest` rule of `.golangci.yml` (depguard): `internal/bot` imports neither `internal/server`, nor `internal/world`, nor a card set;
* `TestBotRules` (`internal/bot/rules_test.go`) reads the bot's source: only the standard library and the engine are imported; no `CardRef`, no card field (`Card`, `Ref`, `CardOf`), no card labels; no string that is a card's ref, and strings only in constants and error messages; no `NewGame`, `Load` or `Imagine`; `Game.Clone` only in `fork`, which copies a world (step 21.5); no package variable holds a game; `Input` holds no game. `TestBotRulesCatch` checks that these rules catch a piece of bad code.

## Choose

`Choose` asks `World(0)` for `Allowed(seat)`. One allowed intent is the answer at once (a discard or an order gets its cards first); but a discard of part of the hand goes to standard and hard, which weigh its sets as any move, so a forced discard drops a cheap card. Else the level picks. The answer is always one of `Allowed`, a discard cut to the asked count and an order of the same cards. An error (no world, nothing allowed, an unknown level) is not a move: the caller sends its forced answer (N-07).

Candidates (Decision 5): every allowed intent; a discard as up to 8 different sets of the asked count, picked at random; an order as the order now and up to 5 random orders.

## Levels

* `random`: one of the candidates at random. For tests and the arena.
* `policy`: the playout policy alone. For tests and the arena.
* `standard`: tries each candidate on one imagined world and plays the best score (below). The server's lower level.
* `hard`: plays each candidate out on many imagined worlds to the end of the turn and plays the best average (below). The server's upper level.
* `flat`: standard with every loot card and item worth alike, as before the card values (below). For the arena and `make bot-values`.

### Playout policy

The policy (`internal/bot/policy.go`, Decision 6) weighs each allowed intent and draws one by weight. It copies nothing, so later levels can use it for every seat of every playout. It reads kinds and public numbers only: coins, the shop price, HP, ATK, evasion, souls. Every weight is a constant at the top of the file.

* Outside its own action phase a seat mostly passes, so the stack resolves.
* In its own action phase: buy when it has the coins, then attack, play loot, use abilities (less after each use this turn), and end the turn when nothing else draws.
* Attack odds: the chance to kill a monster before it kills the attacker, from the attacker's ATK and HP and the monster's HP, evasion and ATK. A roll of 1 counts as a miss, as abilities may raise the evasion during the fight. The monster deck's top is judged as an average monster (HP 3, evasion 4, ATK 1; `TestDeckMonsterMean` keeps them within 1 of the mean monster of every set), with the world's share of soul cards there as its souls (`soulShare`, which `digValue` uses too).
* Attack target: weight from the odds, counted twice, and the target's souls.
* A choice that cancels a play (a target, a mode, a copied ability) gets no weight, so a play is never taken back and tried again for ever.
* A curse goes to another seat, the more souls the likelier.

### Score

`Score(g, p)` (`internal/bot/score.go`, Decision 7) is a table from seat `p`'s side: a won game is the top value, a lost one the bottom. Else it is the seat's own sum less the opponents' average. A sum counts souls, cents, loot cards in hand, items, curses, HP left and being alive this turn; a card's base value comes from `cardValue`, which reads the measured card values (below). Every weight is a constant in that file. A soul (100) outweighs every other part of a table; HP (0.5 each) and being alive (3) weigh little, as both heal at the end of the turn; the flat values, which the `flat` level and the measure keep, are an item 12, a little more than the 10¢ it costs, and a loot card in hand 0.75, less than a cent, so a card that gives 1¢ is played and one that does nothing now is kept.

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

`Budget{Playouts, Time}`: whichever ends first. `Playouts` 0 means no playout cap, `Time` 0 no time limit; both 0 think 1.5 s. With a time budget the first world measures a playout and the rounds share out what the time leaves. The search checks `ctx` before each playout, yields the core after each (`runtime.Gosched`) and, when stopped, plays the best so far. A playout takes about 4 ms on one core, most of it the engine's `checkMaxHP`; a move of 150 playouts about 0.4 – 0.6 s and 200 MB of garbage (Decision 9.8).

## Card values

A flat value misprices cards: a bomb or a cancel is worth more than a penny. So each loot card and item has a measured value (ADR 030, Decision 20), in `internal/bot/values.json`: card → value, samples, standard error. It is generated by `make bot-values` (`cmd/bot_values`), never written by hand, and embedded in the bot.

The tool plays arena games of the `flat` level (2, 3 and 4 players), keeps positions of them, and measures each card on positions where its card lies in its deck: up to 32 for a loot card, 64 for an item. It puts the card into a copy with `enginetest.PlaceFromDeck` (a quiet move from the deck, `CheckInvariants` after; the server and the bot never import `enginetest`, a depguard rule):

* a loot card into the hand of the seat to move: where the seat may play it, stack positions too, its best play and doing nothing are settled as standard settles a candidate, on the same draws; the value is the mean of max(0, play − nothing). Nothing keeps the card at 0.75, so the value is what playing adds over holding it;
* an item into the seat's play area: the copy with it and the copy without are played out by the policy to the end of the seat's next turn on the same seeds; the value is the difference. An active item is also measured as a loot card, by max(0, use − nothing) where the seat may use it, and the larger measure counts: the policy uses items at random, so its playouts miss what an answer is worth.

Each paired difference counts at most one soul either way (`ValueClip`): a game won in one copy and not the other would else be the whole mean.

Both score with the flat values: one pass, no iterations. The same seed gives the same file; `BOT_VALUES=-missing` measures only the cards the table lacks, `BOT_VALUES='-cards a,b'` only those.

`cardValue` reads the table: each card's value weighed against its kind's mean by precision, `(v/se² + mean/τ²) / (1/se² + 1/τ²)`, where τ² is the spread of the kind's values less their noise (at least 0.01). Items are then scaled so the average item is worth 12 again: they were measured over about one round, but stay for the game. A loot card in hand counts half its value (`wHand`): it is not played yet, it can be discarded or stolen, and a turn has few loot plays; so playing a card that does its usual work beats holding it. The kinds are loot, active item (a ↷ ability) and passive item; a card the table lacks (a starting item, a challenge's card) gets its kind's mean. Only the seat's own hand and cards in play count by card: a card in another seat's hand counts half its kind's mean, as the bot cannot tell a remembered card from one its world dealt. The bot reads a card's identity in one function, `tableKey`, to look it up; `TestBotRules` allows it there alone. The hard level's holding value multiplies this base value.

**After cards are added: `make bot-values` (or `BOT_VALUES=-missing`), then `make bots`.** `TestValuesFile` fails while a card a deck holds is missing from the table.

## Memory

The hard level remembers cards it saw go out of sight (ADR 030, Decision 19); standard does not. `world.Memory` (`internal/world/memory.go`) is one seat's memory. It is fed only what the seat's client gets, so it cannot hold a card the seat never saw:

```go
m := world.NewMemory(seat)                 // empty, at the game's start
m.Observe(g.ViewFor(base, viewer), engine.FilterEvents(events, viewer)) // after every step, in order
worlds, _ := world.For(g, seat, m.Known(), seed)                         // the remembered cards go into Imagine
m.Reset()                                  // after a load or a rollback
```

`Observe` must see every step, in order; its first call only takes the view. `Known` returns a copy. A memory is for one goroutine at a time and is not saved.

In the server the room keeps one memory per hard seat (`internal/server/bots.go`, ADR 030 Decision 12.13). `Room.apply`, which every step goes through, feeds it, so the server's own passes and forced answers count too; `resetMemory` starts it over at the game start and after a load, a restart or a rollback. `TestBotsMemoryLikeArena` replays a room's game step by step and checks the room's memories equal ones fed the arena's way; the memory counts its `Observe` calls, so even a skipped step that changed nothing shows.

It remembers the cards of other hands and one place per card of each deck, bottom up. It learns from a look (`looked_at`; a deck card's look names its deck and depth, `Event.Deck`), its own top choice, a deck top shown to all, its treasure peek, a hand its view shows, a card taken from a discard into a hand, a card it gave or got, a card put on, into or under a deck in sight, and a remembered top card that a loot draw (`looted`, then `loot_done`) takes into a hand. The full list is in Decision 19.14.

When unsure, it forgets (Decision 19.15): a whole hand when a card leaves it unseen or its size does not add up with the step's events; every remembered copy of a card that leaves in sight; a whole deck on `deck_order` (a shuffle, a refill from the discard, a reorder, a card moved within it), on `deck_bottom`, on a card taken from it unseen, on a look or draw that shows another card than remembered, or when its size does not add up. Draws it does not follow (slots, gains, attacks on the deck) leave the size off, so the deck is forgotten.

A new engine move of a deck card out of sight that keeps the deck's size needs a `deck_order` event (`emitDeckOrder`), or the memory goes wrong. `TestMemoryIsTrue` (`internal/world`) catches it: 100 simulated games of every set and challenge, every seat's memory checked against the real game at every step; `make sims` runs it long, on 800 games of 2000 steps (`TestMemoryStress`, see [Fuzzing](Fuzzing)). `TestMemoryNamesOnlyShownCards` checks that a memory names only cards its seat's views and events showed.

## Arena

`cmd/bot_arena` plays whole seeded games with every set: bonus souls on, the room deck in odd seeds, deck ratios in every third. Each `-match` is one line-up of levels; game `k` plays seed `-seed`+k and rotates the line-up by `k` seats, so each level plays from each seat. Games run side by side on every core; the results do not depend on it.

```sh
make bots                     # BOT_GAMES=50 per default line-up: standard against random, standard against flat, hard against standard, one hard among three standard, standard×2–4, policy×2–4
go run ./cmd/bot_arena -games 20 -playouts 500 -match hard,standard
go run ./cmd/bot_arena -games 20 -match policy,random -match policy,policy,policy
go run ./cmd/bot_arena -games 20 -cap 400 -match policy,policy,policy,policy
```

It prints, per line-up: games, games with a winner within the turn cap, wins by level against chance, median and max game turns, intents per turn, the share of passes, time per move, player deaths per monster killed; per level, in the same games: its time per move, the share of its turns with no attack, and its loot plays and purchases per seat and game (`Result.Levels`); and for every game cut by the cap, its seed, turn, souls by seat and the kinds of its last 300 intents. It exits 1, in any line-up, when a bot intent is refused or a loop guard ends a game. The finish gate, 95% of games with a winner within the cap (Decision 18.4), is only for line-ups of the levels players meet (standard from 21.5, hard from 21.6) with no random seat: `finishGated` in `cmd/bot_arena`. Random and policy line-ups print their share for reference. The win gates are `winGates` there: standard wins 70% of 2-player games against random, 55% of 2-player games against flat (Decision 20.11), hard 60% of 2-player games against standard, and one hard among three standard wins over 30% (Decision 18.5 – 18.7). Hard thinks `-playouts` playouts a move (`arena.DefaultPlayouts`, 150), never a time, so the results repeat.

Each step the arena asks the waiting seat's bot, with `world.For` on the game for that seat. Each hard seat keeps a memory, fed after every step with its view and filtered events; `-memory=false` (`Match.NoMemory`) plays hard without it. A refused intent is listed and the arena sends the forced answer: the first allowed intent.

The cap, `arena.TurnCap`, is 150 game turns (`Turn.Number`). Two loop guards end a game as a failed gate, printed with its seed: one game turn over 1000 applied intents other than passes (`MaxTurnIntents`), a game over 30000 applied intents (`MaxGameIntents`). The turn guard counts no passes: a long turn of real plays in a 4-player game (a paid ability used 40 times, each a stack item four seats pass) passed 1000 intents with no loop (21.5). A turn is measured, not intents, because about 80% of intents are passes: every stack item waits for a pass from each seat, so a 4-player turn takes about 86 intents and a 2-player turn about 46.

Tests: `TestNoRefusedIntent` (`internal/arena`) plays 30 whole policy games, every set, 2 to 4 players, and fails on any refused intent. `make bots` then runs `TestStandardNoLoop`: the whole standard games where seats once used the same items again and again in a turn (Decision 8.6), each to its end or the cap. It needs `FOUR_SOULS_ARENA_GAMES` and never runs in plain `go test`: one of them, 3 players seed 28, lasts 129 turns and 15 minutes. Plain `go test` checks the rule on built tables (`TestJudge`, `TestStandardKeepsUselessLoot`, `TestStandardNoHopelessAttack`).
