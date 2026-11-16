# Fuzzing

The rules engine has two Go fuzz targets in `pkg/rules_engine/fuzz_test.go`, and three more with the real card sets in `pkg/rules_engine/cards`. They play random games and check after every step that the state is still valid. They catch bugs nobody wrote a test for.

## What one fuzz input does

An input is a seed, a player count and a list of bytes. The game is set up from the seed. At each step the engine lists what the waiting player may do (`Allowed`), and the next byte picks one of those intents. After every step the harness checks:

* the intent the engine offered is accepted (`Step` does not refuse it);
* `CheckInvariants()` passes: every object is in exactly one place matching its zone, no negative cents or damage, slot tops are in play, and the game is waiting on someone;
* no player's view contains another player's hand cards (A-07);
* every object whose zone, damage or charge changed is named by one of the step's events (`CheckEvents`, 22.2; see below);
* the engine does not panic, and the run loop settles (the `maxRunSteps` guard turns an endless loop into a panic).

At the end the recorded steps are replayed from the same setup and every checksum must match (A-08).

Then the game is rolled back to before one step (picked from the seed) and the undone steps are played again, both with `ReplayFrom` from the state after setup. The redo must end in the same bytes as the game (ADR 020, item 7). `FuzzSets` does the same.

A game is capped at `maxFuzzSteps` (2000) steps. Every step hashes the whole state, so longer inputs only slow the fuzzer down.

## Every change has an event

`FuzzAllCards`, `FuzzFocused` and every `playSets` target (`FuzzSets`, `FuzzFocusedSets`, `FuzzChallengeSets`) take a snapshot before each step (`Game.TakeSnapshot`) and call `Game.CheckEvents` with the step's events after it. The code is in `pkg/rules_engine/event_check.go`. A change the client is never told about can't be drawn; this check finds it before a play test does (TS-09).

What it compares:

* each object's zone, damage and charge;
* every move: a move makes a new object (R-ZONE-01), and the snapshot turns on a log of moves in `Game.move`, so a card that moves twice in one step is checked twice;
* the order of every deck: the cards off the longest run kept in order count as moved, so a card put on the bottom with `DeckToBottom` is caught;
* shuffles: `ShuffleDeck` logs the deck, and its new order is one "shuffle" change.

An event names an object by its `Object`, `Prev`, `Source` or stack item source. A move is named when an event names the new object, or the old object at or after the move; an event before the move does not count. A hidden move (a deck card, a hand card) still needs its event: a private event names the object too.

Some events name many objects at once (`wideEvents`): `recharged` and `deactivated` with no object name every thing their player controls, `deck_bottom` names the cards that went into its deck, `moved_to_deck` with no object names the cards an empty deck took back from its discard (R-ZONE-05, 19.174), and `turn_ended` names every heal to 0 (R-TURN-13).

Changes that need no event are on the `silentChanges` list, each with its reason. A gap the check found that is not fixed yet is on the list too, with its 19.N item in `Bug`; 22.3 fixes those and removes them. Silent on purpose:

* a slot card covered by the card that enters the slot on top; its `entered_play` names the slot (R-ZONE-10);
* minion slots renumbered after one is removed; no card moves (R-CHAL-11);
* a shuffle: the deck shows no card and no count change;
* damage removed by a lower max HP: not healing, and the HP left stays the same (R-MECH-21).

When the check fails, the message lists each change no event names, then the step's event kinds. Find the code that made the change; give it an event (log a 19.N item first, as for any bug), or, if it is silent on purpose, add it to the list with its reason.

## Targets

* `FuzzAllCards` — every card in the fuzz sets can show up.
* `FuzzFocused` — the given cards are put into every player's play area at the start, so their abilities come up often. Use it for a new card. The focused game is not replayed, because a replay starts without those items.
* `FuzzSets` (in `pkg/rules_engine/cards`) — the same checks with the real card sets, all of them together. The room deck and the deck ratios are on in some games.
* `FuzzFocusedSets` (same package) — `FuzzSets` with only the sets named in `FOUR_SOULS_SETS`, the room deck on in most games. Use it for a new set: `SETS=b2,r` puts the Base Game with Requiem, so the new cards come up often. It keeps its own corpus, so a run fuzzes new inputs at once instead of replaying the big `FuzzSets` corpus first.
* `FuzzChallengeSets` (same package) — `FuzzSets` as a 2-player challenge (14.21): a test challenge puts a final boss into its slot and two monsters into minion slots, and the top monster card into a new minion slot at the end of each turn, so events in minion slots and removed slots come up often. The invariants also check that every slot card's zone names its row and slot number. The room deck is on in some games. Flags 1 to 3 play How the Krampus Stole Christmas instead, at Normal, Hard or Ultra Hard (14.22): presents, the start-of-turn rolls, minions, the forced attack and the forgo questions. Flag bit 3 (8 to 11) plays Trick or Treat instead (14.32), the difficulty from bits 0 and 1 as for Krampus, 0 as Normal: fear, trick counters, the flip and the attack-start rolls.

## Simulations

`TestSimulations` (in `pkg/rules_engine/cards`) plays whole games with random legal moves for 2, 3 and 4 players and checks the invariants after every step. It prefers acting over passing, so most games end with a winner. A quarter as many Krampus challenge games follow, the difficulties in turn (14.22), then as many Trick or Treat games (14.32); random players mostly lose them. CI plays 100 games; locally:

```sh
cd pkg/rules_engine
FOUR_SOULS_SIMS=3000 go test -count=1 -timeout 30m -run TestSimulations ./cards/
```

3000 games take a few minutes.

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
