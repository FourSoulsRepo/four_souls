# 30. Bots in network games

* **Status:** Accepted
* **Date:** 2026-11-16
* **Authors:** @HardDie

---

## Context

1. Owner, 2026-11-16: bots fill the free seats of a network game (N-21).
   1. Two levels: `standard` and `hard`; the host picks one per bot.
   2. A table has 0 to 3 bots and at least one human.
   3. A bot sees only what a player in its seat sees.
   4. Bots sit down only in the lobby; a kicked seat stays empty (N-08).
   5. No offline mode: the app's own server runs a solo game.
   6. Bots play challenges too (ADR 028).
2. This replaces LR-06: a random bot, never in network games.
3. What the engine gives today
   1. `Game.Allowed(p)`: every intent the engine accepts now.
   2. `Game.Clone()`: a deep copy through JSON.
   3. `Game.ViewFor`: the table as one player sees it.
4. Measured on a 4-player game with every set, 687 objects (2026-11-16)
   1. `Clone`: 1.9 ms.
   2. One applied intent: 0.1 ms.
   3. `Allowed`: 0.01 ms.
   4. Random legal moves: no winner after 6000 steps.
5. So a bot cannot play a game to its end to judge a move.
   1. A random playout never ends.
   2. A 2 ms copy allows about 500 tries a second.
6. The cards: five sets, hundreds of effects, more to come.
   1. A table of per-card advice would never be complete.
   2. Rule 5.24 of `CLAUDE.md`: fix game-wide, never one card.
7. The room (`internal/server`) is one goroutine.
   1. A request that holds it over 100 ms is logged as a stall (18.17).
   2. A seat without a client pauses the game (N-08).
8. A public server runs many games on a few cores.

## Considered options

1. **Bot code in `internal/bot`, on the engine only**
   1. Won. No server import; tests and the arena run it without a network.
2. **Bot as a network client of its own server**
   1. Lost. A view has no game to try moves on.
3. **Bot inside the engine module**
   1. Lost. The engine is the rules; a bot is a player.
4. **Honesty by an imagined game: `Game.Imagine`**
   1. Won. The engine redeals what the seat cannot see; the bot gets only that.
   2. One test proves it: the seat's view of both games is the same.
5. **Honesty by promise: the bot gets the real game and looks away**
   1. Lost. One careless read of a hand makes a cheating bot.
6. **Judge moves by trying them on copies**
   1. Won. Works for every card, old and new, with no card list.
7. **Hand-written advice per card**
   1. Lost. Context 6.
8. **A neural network trained by self-play**
   1. Lost. Training, model files and their upkeep are a project of their own.
9. **Short playouts with a position score**
   1. Won. Context 5: the game end is out of reach.
10. **Full Monte Carlo tree search to the game end**
   1. Lost. Context 4 and 5.
11. **One score, one playout policy; levels differ in depth**
   1. Won. Hard is standard with more worlds, a longer look and no mistakes.
12. **Two separate bots**
   1. Lost. Twice the code to keep right for every new card.
13. **Thinking outside the room goroutine**
   1. Won. Context 7; the answer comes back as a room request.
14. **Thinking inside the room goroutine**
   1. Lost. A 1.5 s search would freeze the table.
15. **Faster `Clone` by hand-written copy**
   1. Won. The search is copies; a test keeps it equal to the JSON copy.
16. **Bots fill seats after a kick**
   1. Lost. Owner: N-08 stays.
17. **Bots vote on rollbacks and kicks**
   1. Lost. Owner, 2026-11-16: only humans vote.

## Decision

Use options 1, 4, 6, 9, 11, 13 and 15.

1. Levels
   1. `standard` and `hard`, as strings in the protocol, saves and records.
   2. `random` exists in `internal/bot` for tests and the arena only.
      1. So do `policy` (Decision 6) and `flat`: standard with every loot card and item worth alike, as before the card values (Decision 20.6).
   3. The server refuses any other level.
2. Engine: `Game.Imagine(p, seed, known)`
   1. Returns a copy of the game as seat `p` may imagine it.
   2. Hidden from `p`: other hands, deck orders, any card `p`'s view hides.
   3. Hidden cards of one kind are pooled, shuffled and dealt back.
   4. Each zone keeps its size.
   5. A card `p` looks at now (a deck look, a shown card) stays.
   6. The copy's RNG is reseeded from `seed`: no real dice, no real draws.
   7. `known` is what the seat remembers (Decision 19); those cards are dealt where it remembers them.
   8. A remembered card that is not among the hidden cards of its pool is skipped.
      1. So is one whose remembered place is not hidden, or is taken.
      2. The memory is never checked against the card's real place: a wrong memory tells nothing.
   9. Empty `known`: every card `p` cannot see now is unknown.
   10. Test: `ViewFor(p)` of the copy equals `ViewFor(p)` of the game.
   11. Test: two seeds give different hidden cards.
   12. The real game is never changed.
   13. The copy names no hidden card elsewhere either.
      1. Past objects nothing points at get the redealt card's name.
      2. So do the options of another player's prompt.
      3. A past object the state points at keeps its card: the stack shows it.
      4. Each new name is drawn from the pool: it says nothing of the old one.
   14. The copy depends only on what `p` knows: the hidden cards of each pool, `seed` and `known`.
      1. Each pool is sorted before the shuffle; the real layout never shows.
      2. Test: an imagined game, imagined again, equals the game imagined, for one seed.
   15. Discard piles stay real, below the top too: open by R-ZONE-07, every card went there in sight.
3. Engine: `Game.Clone`
   1. A hand-written deep copy replaces the JSON round trip.
   2. Test: the copy's `Checksum` and `Save` equal the game's, on simulated games of every set.
   3. Test: a new field of `Game` without a copy rule fails the build's tests.
   4. Goal: under 0.3 ms on the game of Context 4.
4. Package `internal/bot`
   1. Imports the engine; never `internal/server`.
   2. `Choose(ctx, Input) (engine.Intent, error)`; an error is no move: the caller sends its forced answer.
   3. `Input`: the seat, the level, a random source, a budget, `World(i)`.
      1. No memory: the seat's memory reaches the bot only as `known` inside its worlds (Decision 19.2).
   4. `World(i)` returns a new imagined game; the bot never holds the real one.
   5. The answer is always one of `Allowed`, completed for discards and orders.
   6. With one allowed intent it answers at once.
      1. But a discard of part of the hand: standard and hard weigh its sets as any move, so a forced discard drops a cheap card (21.16).
   7. The bot's random source is its own; never the game's RNG.
   8. The match record needs no bot data: it keeps intents (ADR 007).
   9. `internal/world` makes `World` for the arena and the server alike: one `Clone` at once, then `Imagine` per world.
      1. The bot never imports it; a lint rule and `TestBotRules` check the bot's imports and code.
5. Candidates
   1. Every intent of `Allowed(p)`.
   2. A discard: up to 8 sets of `Count` cards, picked at random.
   3. An order: the order now, plus up to 5 random orders.
6. Playout policy: how every seat moves inside a playout
   1. No copies; weights on the allowed intents, then a weighted random pick.
   2. Weights look at kinds and public numbers only: coins, costs, HP, attack odds.
   3. No card names in the policy.
   4. It ends its turn; it never loops: each playout has a step cap.
7. Score: one function, `Score(g, p)`, from seat `p`'s side
   1. A win is the top value; a loss the bottom one.
   2. Then souls, far above everything else.
   3. Then coins, loot cards, items, HP, being alive.
      1. HP and being alive heal at the end of the turn (R-TURN-13): they weigh little, 0.5 an HP and 3 alive, against 1 a cent and 100 a soul.
   4. The same sum for each opponent, taken away from the seat's own.
   5. Standard: the opponents' average.
   6. Hard: the leader counts most; a leader one soul from a win, most of all.
      1. The leader is the opponent with the best full sum: souls, then items, coins, loot, HP (owner via lead, 2026-11-16); not souls alone.
      2. Each opponent's weight is `exp(sum / 100)`, shared out over the opponents: one soul ahead counts e times as much (`wStrength`).
      3. No leader threshold: the weights move smoothly with every part of each sum; equal opponents count as their average, as standard's.
      4. A card or charge that can answer has a holding value: spent now, it is gone against a later threat.
      5. "Can answer": what `Allowed` shows the seat on another seat's turn: a loot card or a card's ability (`answers`). No card names.
         1. On its own turn: the abilities of cards it controls; a room's work on the active player's turn only (R-CARD-36), and a loot card needs a loot play another seat's turn does not give (R-CARD-08).
      6. It is a multiplier on the card's base value, not a term of its own (lead, 2026-11-16): each answer the seat still holds (in hand; in play and charged, or recharging as its turn starts) counts `cardValue` × (1 + `wHold` × (1 + `wHoldLead` × lead)), with 0.2 and 2.
         1. `cardValue` is the one function that gives a loot card in hand and an item in play their base value; 21.16 makes it a value per card, with items scaled to a mean of `wItem` (15; 12 before the tuning of 2026-11-17). The flat values, 0.75 and 12, stay for `flat` and the measure (Decision 20.6).
         2. So a held item answer is worth 3 to 9 more than a spent one; a character has no base value, so its ↷ adds nothing.
      7. Lead: the strongest opponent's share of the weights of every seat, past an even share; 0 with no one ahead, towards 1 as one pulls away.
      8. So the multiplier is never 0 and is up to three times higher with one opponent clearly the strongest.
      9. Test (21.6): 3 seats; hard holds Godhead, which can turn a killing roll into a miss. The weak seat (0 souls, no items) kills a monster with no soul while the strong seat (3 souls, items) moves later: hard passes. The strong seat kills a soul monster: hard answers. It does so with one turn of playouts and the holding value; no plan over several turns.
   7. Weights are constants in one file; the arena tunes them.
   8. A challenge: one team score (Decision 14).
   9. Dig: a monster the seat kills in a settle or playout adds `wDig` to that copy's score.
      1. A kill lets in the next card of the monster deck, which may give a soul.
      2. `wDig` = soul weight × share of soul cards in the monster deck × 0.5 (`digShare`).
      3. The share is of the card property "gives a soul" (`Soul` above 0), over the monster deck of the seat's world: every card there is unseen, redealt by `Imagine` from the real unseen cards; no card names.
      4. Kills are counted from the settle's events: the seat got a monster's rewards. One function scores every settle and playout (`settledScore`); hard uses it too.
      5. Without it a kill of a monster with no soul rarely beats ending the turn: its reward is a few cents, a death costs an item.
8. Standard
   1. One imagined world.
   2. For each candidate: copy, apply, settle, score.
   3. Settle: every seat passes; the policy answers the other prompts; the bot's own choices get the policy's top option.
   4. It runs to the end of what the candidate started: the stack empty and no attack or purchase going on, the turn's end, or the game's end; cap 150 steps.
      1. So an attack runs to its end: the monster dead, the attacker dead and its penalty paid, or the attack cancelled.
      2. A cap of 40 (2026-11-16) cut off a fifth of 2-player fights and half of 4-player ones.
      3. A settle the cap cuts off, the table not quiet yet, judges no move (bug 19.162): a move does not beat doing nothing when any settle of it, or of doing nothing, was cut (`uncut`).
         1. The two then stop at different items of a deep stack; the gain was what the items left unresolved would do, not the move's.
         2. Standard×4 seed 28, a stack of 20–36 items: a paid roll ability measured +6 a use cut off, 0.0 settled to its end; it was used 77 times on one turn.
         3. A choice with no way to do nothing keeps its means.
   5. A candidate whose settle threw a die is settled 12 times and scored by the mean (6 before the tuning of 2026-11-17: an attack that wins a soul in one settle of several was judged on one or two draws).
      1. The first on the world's own dice, the others with the copy's RNG reseeded.
      2. Each settle's policy choices come from a seed of its own.
      3. Every candidate of a move gets the same seeds: scores differ by the move, not by luck.
      4. A candidate with no die keeps one settle.
   6. An action must beat doing nothing: a pass, or an end of turn in the bot's own action phase.
      1. Its gain over doing nothing, settle by settle on the same seeds, counts less 0.25 (a quarter coin) and less a share of its standard error.
      2. On an empty stack one standard error; an attack none (tuning, 2026-11-17).
         1. The spread of an attack's settles is the attack's real risk, not a measuring error: draws like [-10 -10 -10 123 -10 -10] beat ending the turn on the mean, and one standard error cut every such attack; standard left about half its turns without one.
         2. With no bar for any move, standard×3 seed 28 hit the turn guard (items used again and again), so every other move keeps one standard error.
         3. An attack is told apart by the intent's kind alone (`standardNoise`), never by a card. Hard keeps one standard error for every move.
         4. So with 12 draws a hopeless attack whose world kills in 1 settle of 12 is made: `TestStandardNoHopelessAttack` and `TestStandardDigs` allow 1 world of 8.
      3. On a stack two: the move is judged through the dice of every item below it, which its own die shifts.
      4. So a move that gains nothing, or gains only by luck of the dice, is not made; no repeat counter.
   7. Options of a choice that take a play back are left out, as in the policy (Decision 6).
   8. It picks the best score; on equal scores a pass or an end of turn.
   9. One move in ten it picks by weighted chance among the top three, 3:2:1 by rank.
   10. No time budget: under 100 ms per move on average (lead, 2026-11-16).
9. Hard
   1. Many worlds, each a new `Imagine`.
   2. For each candidate and world: copy, apply, play out, score.
   3. The playout runs to the end of the current turn; cap 120 steps.
   4. Every seat follows the policy, so the others may respond.
   5. Candidates are halved round by round; the tries go to the best ones.
      1. Each round gives every candidate still in an even share of the budget, world by world: every candidate meets the same worlds, with the same policy seeds.
      2. Doing nothing (a pass, or an end of turn) is the baseline and stays to the last round, against the one move left.
      3. One world at a time is kept: it is dropped once every candidate still in played on it.
   6. The average score picks the move; no mistakes on purpose.
      1. Every playout is scored by `settledScore`, as a standard settle, with the opponents weighted by strength (Decision 7.6).
      2. A move must beat doing nothing as in Decision 8.6, world by world.
      3. In the noise a gain counts at most four souls (`noiseClip`): a game won or lost in one world of fifty is part of the mean, but as noise it hid what every other world showed.
      4. Known bias of a flat search (lead review, 2026-11-17): the end-of-turn candidate is scored at once, while every other candidate's playout includes the policy's continuation of the turn.
   7. It answers the stack the same way: "pass" is one more candidate.
   8. Budget: `Budget{Playouts, Time}`, whichever ends first.
      1. `Playouts` 0: no cap on playouts; `Time` 0: no time limit; both 0: 1.5 s. The server passes the time only.
      2. With a time budget the first world measures a playout; the rounds share out the playouts that time leaves.
         1. With a time budget and no playout cap the last candidates meet further worlds until the time runs out (639e595): a move uses its whole budget, even with two candidates. The share of such moves and their time wait for the check phase (21.10).
      3. Measured (21.6, i7-8750H, one core, `BenchmarkPlayout`, `BenchmarkHardMove`): a playout 2.2 – 5.3 ms by load, 3.9 ms on a quiet machine, about 260 a second; a move of 150 playouts with its worlds 0.34 – 0.40 s, 375 – 440 playouts a second (playouts on the stack end sooner).
      4. So 1.5 s gives about 400 – 650 playouts on this core, perhaps half that on a weak server core; "3000 playouts" was wrong. Half of a playout was the engine's `checkMaxHP`; 21.10 computes it in one pass, and its effect waits for the check phase.
      5. A move of 150 playouts allocates 142 – 153 MB in 2.4 – 2.6 million allocations, all garbage when it returns: each playout copies its world (about 240 KB) and the engine allocates as it plays (a playout 0.8 – 1.3 MB). A standard move: 4.9 MB, 79 000 allocations.
      6. Again at the end of line 1 (2026-11-17, after 21.16, machine loaded, `go test -bench . -benchmem`, 3 runs): a playout 2.7 – 2.9 ms, 1.0 MB, 22 000 allocations; a hard move of 150 playouts 0.43 – 0.61 s, 196 – 198 MB, 3.8 – 3.9 million allocations; a standard move 6.7 – 7.5 ms, 2.9 – 3.3 MB, 42 000 – 46 000 allocations. In the arena, 12 games side by side: hard 0.56 s a move at 150 playouts, 0.19 s at 50, 1.67 s at 500; standard 10 – 25 ms.
   9. `ctx` ends it early; it then answers with the best so far.
      1. It checks `ctx` before each playout and yields the core after each and every 16 steps within one (`runtime.Gosched`, `yieldSteps`): one core holds many tables.
      2. Tests: a cancelled `ctx` answers an allowed intent within 50 ms; with `GOMAXPROCS=1` another goroutine gets the core at least every 20 ms while hard thinks.
   10. The arena uses the playout count only, so its results repeat.
10. Lobby and protocol
   1. `Seat` gets `bot` (the level) and `bot_n` (its number, from 1).
   2. `Seat.Name` of a bot is `Bot N`, in English, for records and logs.
   3. The client shows a translated name with the level, and a robot icon.
   4. New message `bot`: `{seat, level}`; an empty level takes the bot away.
      1. A new level for a bot's seat keeps its number and name.
      2. A new bot takes the smallest number no seat uses, as a number or as a name.
   5. Only the host sends it; only for a free seat or a bot's seat.
      1. The host sits in seat 0, so a bot never does.
      2. Not at a loaded save's table: its seats wait for its players (N-11).
      3. Refusals are translated messages (`error.host_bots`, `error.bot_seat`, `error.bot_saved`, `error.bot_level`, `error.bots_off`).
   6. A joining player gets a free seat only; never a bot's.
   7. The game starts when every seat is taken and every human is ready.
      1. A bot's seat is always ready.
   8. A table with no human closes.
   9. The game list counts bot seats as taken.
   10. `welcome` → `bots` lists the levels the server allows (Decision 17); the app hides the others.
   11. The protocol version goes up by one (ADR 006): version 7.
      1. v7: bot seats, deck_order (the event of Decision 19.16.2).
      2. The client and the website skip an event kind they do not know (checked at the merge, 2026-11-17).
   12. A host who leaves a table with bots: the next human in seat order moves into seat 0 and hosts.
      1. The game list names the new host; bots keep their seats.
      2. With no human left the table closes (10.8).
      3. So the room's host stays seat 0, as for saves and the end of a game (19.84).
      4. A table without bots keeps seat 0 free for the next player, as before (Rule 6 of step 21).
      5. Its first `bot` message moves the next human into seat 0 first.
   13. A restart file of a lobby table keeps its bots (`savedSeat` → `bot`, `bot_n`).
      1. A level the restarted server does not allow leaves its seat free.
11. Room
   1. `roomSeat` gets the level; a bot seat has no client and no token.
   2. A bot seat is never away: it does not pause the game.
   3. Kick votes and rollback votes count clients, so bots never vote (owner).
   4. A bot never asks for a rollback; its moves are rolled back like any.
   5. A bot sends no presence and has the default dice.
   6. Random seating shuffles bots with the humans.
   7. "No player connected for 30 minutes" counts humans only (N-20).
12. Bot driver: `internal/server/bots.go`
   1. Runs after every broadcast, when the game is not over, paused or voting.
   2. First the old rule: a seat that can only pass is passed at once (N-05).
   3. For each bot seat the prompt waits for, one thinker starts.
   4. The room makes one `Clone` and gives it to the thinker: `world.For` (Decision 4.9).
   5. The thinker's `World(i)` is `Imagine` of that copy.
   6. The answer returns as a room request with the room's step number.
   7. The room applies it like a client's intent: record, rollback history, auto passes.
   8. An answer for an older step is dropped.
   9. Any applied step, rollback, pause or vote cancels the open thinkers.
   10. An intent the engine refuses is logged; the fallback move is used.
      1. The fallback: the forced answer of the response timer (N-07), else any allowed intent (a must-attack forbids pass and end turn).
      2. An allowed intent the engine refuses is an engine bug: logged with seed, step and intent.
   11. The response timer may be off, so a bot never relies on it.
   12. A bot panic never ends the server (21.8).
      1. Bot code and the engine run on imagined states the real game never reached.
      2. The thinker recovers a panic in `Choose` or its `World(i)`; the room recovers one in `world.For`.
      3. The log line `bot panic` has the game, seed, bot seed, step, seat and the stack.
      4. The thinker sends an error answer; the room makes the fallback move (12.10).
   13. The room keeps one `world.Memory` per hard seat (Decision 19; standard none, 19.1).
      1. `Room.apply` feeds it after every applied step: intents, bot moves, fallbacks, the server's passes and forced answers.
      2. With the seat's `ViewFor` and `FilterEvents`, as the arena does; `TestBotsMemoryLikeArena` checks the two agree.
      3. It starts over (`resetMemory`) at the game start, after a load, a restart and a rollback.
      4. Its `Known` goes into `world.For` for the seat's thinker.
13. Pace (owner, 2026-11-16)
   1. A pass or an answer to a prompt: 0.7 s after the prompt came.
   2. An action in the bot's own turn: 1 to 2 s, at random.
   3. Thinking time counts inside the pause.
   4. Times scale with the room's `second`, so tests run fast.
14. Bans, picks and start choices
   1. A ban turn: the bot bans one random character from the pool.
   2. A draft: the bot picks one of its offer at random.
   3. Both levels; nothing tells yet which character is strong.
   4. Start-of-game choices are prompts; `Choose` answers them.
15. Challenges (owner, 2026-11-16)
   1. A challenge table takes a bot in its second seat.
   2. The score is the team's: lost is the bottom, won the top.
   3. Then the final boss's HP, the players' souls, HP, loot and items, summed.
   4. No opponents' part; no leader.
      1. `teamScore` (21.13): every seat's sum as seat `p` sees it (the partner's hand by its kind's mean), less 2 (`wBossHP`) for each HP a final boss has left: four times a player's HP, as the boss is the team's goal; its damage heals at the end of the turn, and an attack on it settles to its end, where only a kill lasts.
      2. Hard scores a challenge by the same team sum: no holding value, which rests on a leader.
      3. The policy never prefers to curse the partner: in a challenge every seat counts as its own.
   5. The arena plays each challenge at each difficulty.
      1. `bot_arena -challenge a,b`: the line-up at every challenge and difficulty; no bonus souls, rooms or deck ratios. It prints won, lost and unfinished; no finish or win gate, only refused intents and loop guards. `make bots` plays none.
16. Saves and restarts
   1. `savedSeat` gets the level and the number.
   2. A loaded table seats its bots at once; humans come back as today.
      1. A bot's seat is nobody's to claim, by token or by name.
      2. A save with a level the server does not allow is not loaded: `error.bot_level`, `error.bots_off` (Decision 17.1.3).
   3. A restart file keeps bot seats of tables and games (N-20).
      1. A kept game keeps its bots, even of a level `-bots` no longer allows: the flag is for seating, and the game needs them.
   4. After a restart nobody waits for a bot; open thinking starts over.
   5. `record.Seat` gets the level and the number: `bot`, `bot_n`, left out for a human.
      1. The record format stays 1: old readers skip the fields, and old records read with none (ADR 007, item 2.20).
17. Server limits (lead, 2026-11-16: the owner's server has 1 CPU and 1 GB)
   1. Flag `-bots=off|standard|all` on the dedicated server, default `all`: the levels a host may seat.
      1. ENV `FOUR_SOULS_BOTS`, config `bots`; any other value stops the start.
      2. `welcome` → `bots` lists the allowed levels; the app hides the others.
      3. A bot of another level is refused with a translated message (Decision 10.5).
      4. The app's own server always allows every level.
   2. At most `GOMAXPROCS / 2` hard thinkers at once, at least 1; the rest wait.
      1. Flag `-bot_think`: the hard budget per move, default 1.5 s (Decision 9.8); ENV `FOUR_SOULS_BOT_THINK`, config `bot_think`, a Go duration from 100 ms to 10 s; any other value stops the start (21.10).
      2. The app's own server always uses the default.
      3. One limit for the whole hub (`thinkSlots`); standard thinks without a slot.
   3. A waiting thinker still answers before the response timer.
      1. A hard thinker that gets no slot within its budget plays that move as standard (21.10): it answers in about twice the budget, far under the shortest response timer (60 s).
   4. Load test: `cmd/server_load` gets tables with bots (21.10).
      1. Check phase (2026-11-17), `GOMAXPROCS=1 GOMEMLIMIT=512MiB`, 3 tables of 1 client and 3 bots, on a machine shared with other jobs.
      2. Hard (load average 12.0 → 5.4): response p99 32 ms, max 40 ms; longest move 3.1 s; peak RSS 35 MB; GC 5% of the CPU; the server took 98% of the core. 54% of hard moves found no slot and played standard; no search ended with 0 playouts (at least 179).
      3. Standard (load average 9.2 → 1.5): response p99 20 ms, max 38 ms; longest move 192 ms; peak RSS 34 MB; GC under 1%; the server took 4% of the core.
18. Arena: how strength is checked
   1. `cmd/bot_arena`: whole games between levels, seeded, seats rotated.
   2. It prints wins by level, unfinished games, steps and time per move.
   3. `make bots` runs it; `BOT_GAMES` sets the count (50 a line-up). Then it replays the standard games that once looped (`TestStandardNoLoop`, Decision 8.6.2.2).
   4. Gate: 95% of games end with a winner within the turn cap.
      1. The cap: 150 game turns (`Turn.Number`); a game past it is unfinished.
      2. For line-ups of the levels players meet (standard from 21.5, hard from 21.6) with no `random` seat, each player count on its own: standard×2, ×3 and ×4 each finish 95%.
      3. `random` and `policy` line-ups print their share for reference only.
      4. Loop guards fail the gate in any line-up, printed with the seed: a game turn over 1000 applied intents other than passes; a game over 30000 applied intents. So does a refused intent.
         1. Passes do not count in a turn (lead, 2026-11-16): 4-player turns of real plays with combo cards passed 1000 intents, nearly all passes, with no loop.
      5. Per line-up the arena prints median and max game turns, intents per turn, the share of passes; per level, its time per move, the share of its turns with no attack, its loot plays and purchases per seat and game (lead, 2026-11-17: to compare levels in the same games).
      6. Policy against policy, 50 games each (2026-11-16): 2 players median 42 turns, max 150, 46 intents a turn, 77% passes; 3 players median 50, max 138, 65, 83%; 4 players median 58, max 150, 86, 87%. Finished within the cap: 96%, 100%, 92%.
   5. Gate: standard wins at least 70% of 2-player games against random.
   6. Gate: hard wins at least 60% of 2-player games against standard.
      1. Hard's strength is measured at the playouts a weak server gives: `-playouts`, 150 by default; never by time.
      2. Hard against standard is also reported at 50 and 500 playouts a move; the gate must hold at 150 (owner, 2026-11-16). End of line 1 (2026-11-17, 50 games, 2 players): 64% at 50, 68% at 150, 82% at 500; again in the check phase, after the tuning.
   7. Gate: one hard among three standard wins over 30%; chance is 25%.
   8. Gate: no bot intent is refused, on every set, 2 to 4 players.
   9. CI plays a few games for the last gate only; win rates are a local check.
   10. A missed gate is tuned in the score and policy, never by a peek at hidden cards.
19. Memory (owner, 2026-11-16)
   1. The hard bot remembers cards it saw; standard does not.
   2. `world.Memory` in `internal/world`: one per hard seat, empty at the game start.
      1. It names cards (`CardRef`), so it is not in `internal/bot`; the bot never imports it.
      2. In `internal/world`, next to `world.For`, which takes its `Known`: one package for the arena and the server.
   3. It is fed only what the seat's client gets: `ViewFor` and `FilterEvents` of its seat.
   4. So it cannot hold a card the seat never saw.
   5. It keeps the known cards of each other hand.
   6. It keeps the known cards on the top and the bottom of each deck.
   7. Learned: a card the seat sees go into a hidden place.
      1. A card it looked at or was shown.
      2. A card taken from a public zone into a hand.
      3. A card put on or under a deck in its sight.
   8. Forgotten: a card seen leaving; a shuffled deck; any doubt.
   9. When unsure, forget: a wrong memory is worse than none.
   10. It gives `known` to `Imagine` (Decision 2.7).
   11. An event that lacks what a rule needs gets the field in the engine, filtered like any other.
   12. Not saved: after a save or a restart the memory starts empty; after a rollback too (`Reset`).
   13. The arena keeps a memory per hard seat the same way; `-memory=false` plays hard without it.
   14. The ways a seat sees a card that then hides (21.15):
      1. A look (`looked_at`, the looker's only): at deck cards, with their place, or at another hand's cards (`By`).
      2. Its own top choice (`top_choice`): the card stays on top or goes to the bottom.
      3. A deck top shown to all (`card_revealed`, "<deck> deck").
      4. The treasure peek its view shows (`treasure_top`).
      5. A hand its view shows: a revealed hand, every hand in a co-op challenge.
      6. A card taken from a discard into a hand (`looted`, "from the discard"), in sight of all.
      7. A hand card it gave or got (`gave_card`; the giver is in `SeenBy`).
      8. A card put on, into or under a deck in sight (`moved_to_deck` with its card); its own hand card put on the loot deck.
      9. A remembered top card drawn by a loot draw (`looted`, then `loot_done` of as many): it is in the drawer's hand.
      10. Not memory: cards in play, in slots and in discards stay in sight (Decision 2.15).
   15. What it forgets (19.8 – 19.9):
      1. A hand: all of it when a card leaves unseen (`gave_card` or `moved_to_deck` with no card); every remembered copy of a card that leaves in sight; all of it when the hand size does not add up with the step's events.
      2. A deck: all of it on `deck_order` or `deck_bottom`; on a card taken from it unseen (`looted` with no `loot_done`); when a look or a draw shows another card than remembered; when the deck size does not add up.
      3. It keeps one place per card of each deck: a loot draw, a mill, a card put on top, under the top cards or on the bottom moves the places. Other draws (slots, gains, attacks on the deck) are not followed: the size is off and the deck is forgotten.
      4. No reorder keeps every count and sends no event (lead review, 2026-11-17): every card of every set moves hand and deck cards only through the engine's movers (grep of `cards/` for each, and for direct writes to `Decks` and `Hand`: none), and each mover sends an event: a hand swap (Incubus, Chaos, Clog) is one `gave_card` per card, a hand card put on the loot deck `moved_to_deck`, a reorder `deck_order` or `top_choice`; a card taken from a deck into a slot (`PlaceFromDeck`, Red Key) sends none but leaves the size off, so the deck is forgotten.
   16. Engine fields for it (19.11):
      1. `Event.Deck` on `looked_at` of a deck card: the deck and the depth from the top. `FilterEvents` hides it with the card.
      2. Event `deck_order` ("<deck> deck"), public, no card, no barrier: `ShuffleDeck`, a refill of an empty deck from its discard, `SetDeckTop` and `DeckToBottom`; these moved deck cards with no event. Setup shuffles tell nothing: no one remembers yet. The client's history skips it.
      3. The texts that name a deck or a discard are the engine's: `DeckText`, `DeckOfText`, `DeckTextTop`, `DeckTextBottom`, `TextFromDiscard`; the engine writes them and the memory reads them, so the two cannot drift (lead review, 2026-11-17).
   17. API, for the arena and the room goroutine: `world.NewMemory(seat)`; `Observe(view, events)` after every step, in order, with the seat's `ViewFor` and `FilterEvents`; `Known()` for `world.For`; `Reset()`. One goroutine at a time.
   18. Test (21.15): 100 simulated games of every set and challenge, 500 steps each, every seat checked against the real game at every step: no wrong memory.
      1. A long run of the same, 800 games of 2000 steps (`TestMemoryStress`), runs from `make sims` only.
20. Card values (owner, 2026-11-16; sub-step 21.16)
   1. A flat value for every loot card (0.75) and item (12) misprices cards: a bomb or a cancel is worth more than 1¢.
   2. Cards are added often, so values are measured, never written by hand.
   3. `make bot-values` (`cmd/bot_values`) measures each card on sampled positions from saved arena games, not on whole games.
      1. Whole games are too slow and too biased: a card is played when it helps, and a strong item sits with a seat that is ahead anyway.
   4. Paired measurement on common seeds, in `settledScore` units, as in Decision 8.6.
      1. Loot: in positions where the card can be played, stack positions included, so answers are not zero: mean of max(0, play − nothing).
      2. The positive part assumes a card is played only when it helps.
      3. "Nothing" keeps the card in hand at the flat 0.75: the value is what playing adds over holding it. So a card worth 5¢ now is worth 4.25 held, and the bot plays it; a card whose play never helps is worth 0, or its kind's mean with few samples.
      4. Each play and doing nothing are settled as standard settles a candidate (Decision 8.4 – 8.5): 6 draws when dice roll (`lootDraws`), the same for both; the best play's mean gain counts.
         1. Standard settles 12 since the tuning of 2026-11-17; the measure keeps 6, as it keeps the flat values (20.6.1), so the table stays as it was measured.
      5. Item: a clone with the item and one without it; a policy playout to the end of the seat's next turn; the difference.
      6. The item is added, never removed: no destroy is needed.
      7. The playout runs to the end of the seat's next turn that starts after the position, at most 1500 steps; every seat moves by the policy, the seat's own choices get its best option; dice and choices come from the same seeds in both copies. With the item the item counts its flat 12 there (`wItemFlat`, Decision 20.6.1), so the value is 12 plus what the item does in that time.
      8. Positions: the tool plays arena games of the `flat` level, 12 of each table size by default, side by side, each to at most 60 game turns, and keeps about one priority position in 23 steps, at most 60 a game, with the seat to move. They stay in memory: the same seed plays the same games, so a run needs no saved file.
         1. A card is measured on the positions where its card lies in its deck (and, for loot, the seat may play it), in an order of the card's own, up to 32 for a loot card and 64 for an item; an item on 4 seeds a position.
         2. Playing the games is most of a run: 24 games to the 150-turn cap took 16 minutes on 12 cores, the slowest game alone most of it; the cap of 60 turns keeps the early and middle game, where most positions are.
      9. Each paired difference (a draw of a loot play or a use, a seed of an item) counts at most one soul either way (`ValueClip`): a game won or lost in one copy and not the other is a million points; unclipped, one such sample in a hundred made an item's standard error ten thousand.
      10. An active item is also measured as a loot card is (lead, 2026-11-17): on positions where the seat may use it, stack positions included, the mean of max(0, use − nothing). Its value is the larger of this and its playouts, before the scaling of 20.8.5.
         1. The policy plays answers badly: in a playout it uses an item at random, rarely when it counts, so the playouts miss what an answer is worth (Godhead measured −1.2 by them). Hence the larger of the two.
   5. A card is put into a clone by an existing test path, or by one `enginetest` helper; `CheckInvariants` after it.
      1. Not in the `Game` API; `internal/server` and `internal/bot` never import it (depguard rule `no-enginetest`).
      2. No test path put a card into a hand of a running game; the engine's `TakeFromDeck` gains an item, but with events and triggers outside `Apply`. So the helper is `enginetest.PlaceFromDeck`: it moves the card from its deck by the exported state fields, as a draw or a gain moves it, as a new object, with no event and no trigger. An item's "when this enters play" is not measured.
   6. One pass, with the table off (the flat values); no iterations, so the measure does not feed on itself.
      1. The flat values are constants of their own: `wLootFlat` 0.75 and `wItemFlat` 12, as before 21.16 (lead, 2026-11-17).
      2. The `flat` level and the measure use them; a tuning of `wItem` or `wLoot` moves neither, so `flat` stays the level before the card values and the gate of 20.11 compares against it.
   7. Output: a generated data file, card id → value, n, standard error. A fixed seed gives the same file for the same cards.
      1. A mode measures only cards missing from the table, for a new set.
      2. A full run takes at most 30 minutes on all cores of a work machine; else fewer samples a card, with the error printed.
         1. Measured (2026-11-17, 12 cores, i7-8750H): 11 minutes, 5 of them the games; 32 samples a loot card, 64 an item. Median standard error: loot 0.79, active items 3.05, passive items 2.70. With 32 item samples it was 8 minutes and 4.4 and 3.8.
   8. `cardValue` (Decision 7.6) reads the table, each value weighed against the mean of the card's kind by its precision.
      1. Kinds: loot, active item (a ↷ ability), passive item; an unknown card gets its kind's mean. A starting item is in no deck, so it is never measured: it gets its kind's mean.
      2. Own hand by card; other hands by the mean, as before.
         1. A card in another seat's hand counts its kind's mean, not its table value (lead, 2026-11-16): in an imagined world it is a guess, and its price would only add noise between worlds.
         2. A card the seat remembers there could count its table value; but `World(i)` does not tell a remembered card from a dealt one, and the API is not widened for it: every card of another hand counts the mean.
      3. The answer-holding bonus stays a multiplier on top (Decision 7.6); it is not added twice.
         1. A base below 0 holds nothing: the multiplier counts max(0, base), so a held answer is never a loss. Godhead measured −1.2 ± 3.6 by its playouts alone; with its uses (20.4.10) it is 1.4 ± 1.1 on the 5 positions where it could be used, 3.4 after the shrink and scaling, against an active mean of 10.3.
      4. Shrink by precision (lead, 2026-11-17): v' = (v/se² + m/τ²) / (1/se² + 1/τ²), m the kind's mean over its measured cards, τ² the variance of their values less their mean se², at least 0.01 (`tauFloor2`). A card with no sample gets m; one with se 0 keeps its value. A kind with no measured card keeps the flat value.
         1. It replaced (n × v + 8 × m) / (n + 8), which weighed a card by its samples, not by its noise: items have five times the error of loot.
      5. Items are scaled by `wItem` over the mean of every measured item, active and passive (lead, 2026-11-17): the average item is worth `wItem` again (15 since the tuning of 2026-11-17: coins are spent almost only on items, so an average item is worth clearly more than its 10¢), the order stays, a bad item stays below 0.
         1. A correction for the measurement's horizon: an item is measured over about one round (to the end of the seat's next turn), but stays for the game. Unscaled, the means were 4.4 and 7.3 against a price of 10¢, so every purchase lost.
      6. A loot card in hand counts `wHand` = 0.5 of its value; another seat's hand card 0.5 of the kind's mean (lead, 2026-11-17).
         1. A card in hand is not played yet: it can be discarded or stolen, and a turn has few loot plays.
         2. The measure is what playing adds; at full value a card in hand was worth as much as playing it, and standard stopped playing loot (a dime: +10¢ against −10, below `wAct`).
   9. The bot code still names no card: the table is generated data; a test checks every card of every set is in it or falls back.
      1. The table is `internal/bot/values.json`, embedded (`go:embed`). To look a card up the bot reads its identity in one function, `tableKey`; `TestBotRules` allows `Card` there alone, and only for the seat's own hand and cards in play.
   10. After cards are added: `make bot-values` (`BOT_VALUES=-missing` for the new cards only), then `make bots` (21.14: in `Bots.md` and `Adding-a-card.md`).
      1. `TestValuesFile` fails while a card a deck holds is missing from the table.
   11. Gate: 200 games of 2 players, standard with the table wins at least 55% against standard without it (`flat`, with the flat values of 20.6.1); time per move does not grow.
      1. `make bots` plays `standard,flat` and gates it (`winGates` in `cmd/bot_arena`).

## Consequences

### Positive

1. One player can play a full table.
2. New cards work with bots without bot code.
3. A cheating bot is a failed test, not a matter of care.
4. `Clone` gets faster for rollback and checks too.
5. The arena finds games that never end.

### Negative and risks

1. A bot cannot plan over several turns.
   1. The score stands in for the future.
2. With deck ratios a human does not know which cards the decks got (A-16).
   1. `Imagine` deals the real hidden cards, so a bot guesses from the right pool.
   2. It never knows where a card is; only that it is somewhere hidden.
   3. A pool of every card the sets could deal is a later step.
3. Policy weights are guesses until the arena tunes them.
4. Hard thinking costs CPU on a public server.
   1. The thinker limit and the `-bots` flag bound it.
5. Soft rollback with one human: bot moves can be undone, the human's own cannot.
   1. The human is the only voter and never the author of a bot's move.
6. A hand-written `Clone` can miss a new field.
   1. Decision 3.3 catches it.
7. A random draft pick can give a bot a weak character.

### Neutral

1. LR-05 and LR-06 change: practice games use these bots through the app's server.
2. Sub-step 11.4 moves to step 21.
3. Bot moves look like any player's in records and replays.
4. A later step may give bots a pick by character win rate.
5. The standard bot has no memory; that is one more gap between the levels.
