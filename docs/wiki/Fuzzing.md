# Fuzzing

The rules engine has two Go fuzz targets in `pkg/rules_engine/fuzz_test.go`, and three more with the real card sets in `pkg/rules_engine/cards`. They play random games and check after every step that the state is still valid. They catch bugs nobody wrote a test for.

## What one fuzz input does

An input is a seed, a player count and a list of bytes. The game is set up from the seed. At each step the engine lists what the waiting player may do (`Allowed`), and the next byte picks one of those intents. After every step the harness checks:

* the intent the engine offered is accepted (`Step` does not refuse it);
* `CheckInvariants()` passes: every object is in exactly one place matching its zone, no negative cents or damage, slot tops are in play, and the game is waiting on someone;
* no player's view contains another player's hand cards (A-07);
* every object whose zone, damage or charge changed, and every player whose cents, damage or soul value changed, is named by one of the step's events (`CheckEvents`, 22.2 and 22.10; see below);
* the engine does not panic, and the run loop settles (the `maxRunSteps` guard turns an endless loop into a panic).

At the end the recorded steps are replayed from the same setup and every checksum must match (A-08).

Then the game is rolled back to before one step (picked from the seed) and the undone steps are played again, both with `ReplayFrom` from the state after setup. The redo must end in the same bytes as the game (ADR 020, item 7). `FuzzSets` does the same.

A game is capped at `maxFuzzSteps` (2000) steps. Every step hashes the whole state, so longer inputs only slow the fuzzer down.

## Every change has an event

`FuzzAllCards`, `FuzzFocused`, every `playSets` target (`FuzzSets`, `FuzzFocusedSets`, `FuzzChallengeSets`) and `TestSimulations` (`make sims`, 22.10) take a snapshot before each step (`Game.TakeSnapshot`) and call `Game.CheckEvents` with the step's events after it. The code is in `pkg/rules_engine/event_check.go`. A change the client is never told about can't be drawn; this check finds it before a play test does (TS-09).

What it compares:

* each object's zone, damage and charge;
* every move: a move makes a new object (R-ZONE-01), and the snapshot turns on a log of moves in `Game.move`, so a card that moves twice in one step is checked twice;
* the order of every deck: the cards off the longest run kept in order count as moved, so a card put on the bottom with `DeckToBottom` is caught;
* shuffles: `ShuffleDeck` logs the deck, and its new order is one "shuffle" change; its `shuffled` event names it (19.178).
* each player's cents, damage and soul value (22.10). The soul value is what the win counts (`SoulValue`): an item that becomes a soul where it is, or a soul stolen in place, moves no card, so the object check alone misses it.

An event names an object by its `Object`, `Prev`, `Source` or stack item source. Object 0 is a real card, the first deck card made at setup, and an unset field is 0 too; so an event names object 0 only when its `Card` is object 0's card (22.10). Object 0 never leaves its deck under its own ID, and every event that names a deck card carries the card. A move is named when an event names the new object, or the old object at or after the move; an event before the move does not count. A hidden move (a deck card, a hand card) still needs its event: a private event names the object too.

A player's change needs an event of its own kind (`playerKinds`) whose `Player` is that player:

* cents: the step's cents events must add up to the change exactly: `gained_cents`, `lost_cents`, `paid` (not HP), `stole` (its `Text` is the robbed player) and `gave_cents` (`By` is the receiver plus 1). A popup that shows the wrong amount is caught too, as 19.2 was for HP;
* damage: `damaged`, `healed`, `paid`, `prevented` (a death ward), `death_pending` (a kill) or `died`;
* soul value: `gained_soul`, `destroyed`, `discarded`, `gained_treasure` or `character_dealt` (a character with a soul value); a stolen `gained_treasure` also names the player it was taken from in `By` (19.193).

An event of another kind (`passed_priority`, say) does not name a player's change, so a player who acts in the same step does not hide a silent change of their own cents.

Some events name many objects at once (`wideEvents`): `recharged` and `deactivated` with no object name every thing their player controls, `deck_bottom` names the cards that went into its deck, `moved_to_deck` with no object names the cards an empty deck took back from its discard (R-ZONE-05, 19.174), `shuffled` names the whole new order of its deck (19.178), and `turn_ended` names every heal to 0, players' too (R-TURN-13).

Changes that need no event are on the `silentChanges` list, each with its reason. A gap the check found that is not fixed yet is on the list too, with its 19.N item in `Bug`; 22.3 fixes those and removes them. Silent on purpose:

* a slot card covered by the card that enters the slot on top; its `entered_play` names the slot (R-ZONE-10);
* minion slots renumbered after one is removed; no card moves (R-CHAL-11);
* damage removed by a lower max HP: not healing, and the HP left stays the same (R-MECH-21). The snapshot keeps each object's and each player's max HP, so only a drop no larger than the max HP's fall passes (22.3). A max HP can fall and rise again in one step (six items rerolled at once): `checkMaxHP` logs the damage each fall removes while a snapshot is on, and that counts too (22.10).

When the check fails, the message lists each change no event names, then the step's event kinds. Find the code that made the change; give it an event (log a 19.N item first, as for any bug), or, if it is silent on purpose, add it to the list with its reason.

## Saves at a nested question

`playSets` (so `FuzzSets`, `FuzzFocusedSets`, `FuzzChallengeSets`) and `TestSimulations` also save and load the game whenever a nested question waits (`Game.NestedAsk`, 19.177): one asked inside another effect, such as Gimpy's "Loot 1, then discard" or Plan C's `DestroyThis`. The loaded twin must have the same whole-state checksum, then step with the same intents and stay the same until no nested question waits (`loadCheck`, `cards/load_check_test.go`). `TestSimulations` logs how many such saves it made and for how many cards; 100 games reach about 30 cards. The input that found 19.143 and 19.177 is in the `FuzzSets` corpus (`55bb281b0c9ed398`).

## Targets

* `FuzzAllCards` — every card in the fuzz sets can show up.
* `FuzzFocused` — the given cards are put into every player's play area at the start, so their abilities come up often. Use it for a new card. The focused game is not replayed, because a replay starts without those items.
* `FuzzSets` (in `pkg/rules_engine/cards`) — the same checks with the real card sets, all of them together. The room deck and the deck ratios are on in some games.
* `FuzzFocusedSets` (same package) — `FuzzSets` with only the sets named in `FOUR_SOULS_SETS`, the room deck on in most games. Use it for a new set: `SETS=b2,r` puts the Base Game with Requiem, so the new cards come up often. It keeps its own corpus, so a run fuzzes new inputs at once instead of replaying the big `FuzzSets` corpus first.
* `FuzzChallengeSets` (same package) — `FuzzSets` as a 2-player challenge (14.21): a test challenge puts a final boss into its slot and two monsters into minion slots, and the top monster card into a new minion slot at the end of each turn, so events in minion slots and removed slots come up often. The invariants also check that every slot card's zone names its row and slot number. The room deck is on in some games. Flags 1 to 3 play How the Krampus Stole Christmas instead, at Normal, Hard or Ultra Hard (14.22): presents, the start-of-turn rolls, minions, the forced attack and the forgo questions. Flag bit 3 (8 to 11) plays Trick or Treat instead (14.32), the difficulty from bits 0 and 1 as for Krampus, 0 as Normal: fear, trick counters, the flip and the attack-start rolls.

## Simulations

`TestSimulations` (in `pkg/rules_engine/cards`) plays whole games with random legal moves for 2, 3 and 4 players and checks the invariants and the events (`CheckEvents`, 22.10) after every step. Whole games reach late-game states a short fuzz input rarely does: big piles of cents, stolen souls, auctions. It prefers acting over passing, so most games end with a winner. A quarter as many Krampus challenge games follow, the difficulties in turn (14.22), then as many Trick or Treat games (14.32); random players mostly lose them. CI plays 100 games; locally:

```sh
cd pkg/rules_engine
FOUR_SOULS_SIMS=3000 go test -count=1 -timeout 30m -run TestSimulations ./cards/
```

3000 games take a few minutes.

`make sims` then runs `TestMemoryStress` (in `internal/world`): the bot memory's check of `TestMemoryIsTrue` (ADR 033, Decision 19.18) on 800 simulated games of every set and challenge, up to 2000 steps each; every seat's memory is checked against the real game at every step. Plain `go test` plays 100 games of 500 steps; the long run needs `FOUR_SOULS_MEMORY_STRESS` set and never runs with `-short`:

```sh
FOUR_SOULS_MEMORY_STRESS=1 go test -count=1 -timeout 60m -run TestMemoryStress ./internal/world
```

## Running

CI runs only the seed corpus, as normal tests (`go test ./...`). Real fuzzing runs on your machine. From the repo root:

```sh
make fuzz                          # fuzz-engine, fuzz-sets, then fuzz-challenge, 10 minutes each
make fuzz-sets FUZZTIME=2m
make fuzz-challenge FUZZTIME=2m    # a final boss and minions; Krampus; Trick or Treat
make fuzz-focused FOCUS=razor,d6   # cards of the engine's test sets
make fuzz-new SETS=b2,r            # a new set with the Base Game
make sims SIMS=3000
```

Or by hand, from `pkg/rules_engine`:

```sh
go test -run '^$' -fuzz=FuzzAllCards -fuzztime=10m -fuzzminimizetime=5s .
FOUR_SOULS_FOCUS=razor,d6 go test -run '^$' -fuzz=FuzzFocused -fuzztime=5m -fuzzminimizetime=5s .
```

`FOUR_SOULS_FOCUS` takes card refs separated by commas.

The log sometimes shows `0/sec` for a while, right after `new interesting` grows. That is the fuzzer shrinking the new input; those runs are not counted. `-fuzzminimizetime=5s` keeps those pauses short (the default is a minute).

## When it fails

Go writes the failing input to `testdata/fuzz/<Target>/<hash>` and prints the command to rerun it:

```sh
go test -run 'FuzzAllCards/<hash>' .
```

Fix the bug, then either keep the file, which makes it part of the seed corpus that CI runs, or better, turn it into a small rules test or a [situation payload](Situation-payload) that shows the bug clearly.
