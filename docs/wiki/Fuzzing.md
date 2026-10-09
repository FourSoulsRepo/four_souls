# Fuzzing

The rules engine has two Go fuzz targets in `pkg/rules_engine/fuzz_test.go`. They play random games and check after every step that the state is still valid. They catch bugs nobody wrote a test for.

## What one fuzz input does

An input is a seed, a player count and a list of bytes. The game is set up from the seed. At each step the engine lists what the waiting player may do (`Allowed`), and the next byte picks one of those intents. After every step the harness checks:

* the intent the engine offered is accepted (`Step` does not refuse it);
* `CheckInvariants()` passes: every object is in exactly one place matching its zone, no negative cents or damage, slot tops are in play, and the game is waiting on someone;
* no player's view contains another player's hand cards (A-07);
* the engine does not panic, and the run loop settles (the `maxRunSteps` guard turns an endless loop into a panic).

At the end the recorded steps are replayed from the same setup and every checksum must match (A-08).

A game is capped at `maxFuzzSteps` (2000) steps. Every step hashes the whole state, so longer inputs only slow the fuzzer down.

## Targets

* `FuzzAllCards` — every card in the fuzz sets can show up.
* `FuzzFocused` — the given cards are put into every player's play area at the start, so their abilities come up often. Use it for a new card. The focused game is not replayed, because a replay starts without those items.

## Running

CI runs only the seed corpus, as normal tests (`go test ./...`). Real fuzzing runs on your machine, from `pkg/rules_engine`:

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
