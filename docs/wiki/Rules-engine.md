# Rules engine

Notes on `pkg/rules_engine` as it is now. The design and its reasons are in ADR 005 (`docs/architecture/005-rules-engine-design.md`); this page is a map of the code.

## The shape

The engine is a pure state machine. It has no files, network, clock or `math/rand` (enforced by `depguard`), so the same code runs on the server and in the browser (WASM).

```go
g, events, err := rulesengine.NewGame(setup) // deals, shuffles, rolls for first
p := g.Prompt()                              // who must answer what
allowed := g.Allowed(p.Player)               // every intent that would be accepted
events, err = g.Apply(intent)                // validate, then run until input is needed
```

`Setup.Banned` lists cards the host left out (GS-13): no copy goes into a deck, a banned character is never dealt, and the list stays in `Game.Banned`. A starting item cannot be banned; ban its character instead.

`Apply` either refuses the intent with a `RuleError` (rule ID plus reason, state unchanged) or applies it and runs automatic work until the game waits on a player again. The engine never picks for a player, even when only one option exists.

## State

`Game` (`state.go`) is plain data that encodes to JSON:

* players, objects, decks and discards, shop / monster / room slots;
* the action queue, the stack, pending triggers, the open priority window;
* turn, attack and purchase state, the current `Waiting` prompt or `Choice`;
* the RNG (PCG32) itself, so a saved game continues with the same rolls.

There are no maps in the saved state (map order is random). The card index is not saved; `Load` rebuilds it from the card sets.

Every card is an object with an ID. Moving a card to another zone gives it a new ID, as the rules treat it as a new object. So `g.Objects` only grows. Code that looks at "what is in play" must use `g.inPlay()` (characters, play areas, slot tops), not scan all objects.

## Where things happen

| File | What |
|------|------|
| `turn.go` | turn steps, priority windows, the `run()` loop, win check |
| `stack.go` | stack items, resolving, dice rolls |
| `actions.go` | the action queue and replacement effects |
| `combat.go` | attacks, damage, death, rewards |
| `shop.go` | purchasing |
| `ability.go`, `activate.go` | effect blocks, activating abilities, targets |
| `triggers.go`, `statics.go` | triggered and static abilities |
| `choice.go` | `Choose` prompts (targets, order of replacements and triggers, penalties) |
| `allowed.go`, `aims.go`, `blocked.go`, `view.go` | allowed intents, their first targets, and why the rest are not; per-viewer hidden information |
| `record.go`, `situation.go` | step records, replay, save/load; situation payloads |
| `invariants.go` | state checks used by tests and fuzzing |

Cards are Go values built from effect blocks; see [Effect blocks](Effect-blocks).

## Card sets

Each set is a package under `cards/`, e.g. `cards/b2`, with one file per card. `cards.Sets()` lists every set and `cards.Find("b2")` returns one.

`go run ./cmd/cardgen -set b2` (from the repo root) writes the files from `card_db` data:

* each card file holds the printed text as a comment, the game data (kind, copies, stats, soul, rewards, keywords) and a `TODO(card):` line;
* a file with that line is a stub and is rewritten on every run;
* when you implement a card, delete the line: the file is yours and the generator never touches it again;
* `set_gen.go` lists every card of the set and is always rewritten, so keep the variable name the generator chose;
* data it cannot read (e.g. a printed `X` or a rolled reward) is listed as `Not generated:` in the comment.

The generator lives in the app module because the engine never imports `card_db`.

## Card status report

`go run ./cmd/cardgen -report` rewrites `pkg/rules_engine/docs/card-status.md`: progress per set, per engine kind, per printed type and per card.

* Not started: the card file still has the `TODO(card):` line.
* Implemented: the line is gone.
* Tested: the card also has `<card>_test.go` next to it.

A test fails when the committed report is out of date, so run the command after finishing a card.

## The run loop must always settle

`run()` keeps doing automatic work (queued actions, triggers, step changes) until the game waits on a prompt. Any path that leaves no prompt and changes nothing spins forever. Rules for new code:

* every step either opens a priority window, asks something, or moves to the next step;
* a window step with no open window opens one again instead of returning;
* priority given while actions are queued is deferred (`Priority.Deferred`) and opened once the queue is empty;
* `maxRunSteps` turns an endless loop into a panic with a clear message, so tests and fuzzing fail instead of hanging.

## Testing

* Rules tests name the rule they check (`R-…` IDs from the rules document).
* [Situation payload](Situation-payload): build a position from JSON, run intents, compare the result.
* [Fuzzing](Fuzzing): random games with invariant, leak and replay checks.
* Determinism: replaying the same setup and intents must give the same checksum after every step.

## Choice options in views

`GameView.Options` tells the chooser what each option of a choose prompt points at (`OptionView`: an object, a player, a deck, or a stack item), one per `Waiting.Options`. Clients use it to let the player pick on the table. It is built from `g.Choice` at view time, so it is not saved state and does not change checksums.

`GameView.Activating` names the card whose mode or targets the open choice asks for, while a play or use is being put on the stack (`g.Activating`). Only the chooser gets it, like `Options`, so a hand card is never named to others. Clients aim from that card (15.34).

## Why an intent is not allowed

`Game.Blocked(p)` (`blocked.go`, 15.29) lists what player p can't do now, each with a short reason for the player: every loot card in their hand, every card of theirs with an activated ability, every room with one, and attack, purchase and end turn. Allowed intents are left out. Each intent is run through `check` as if p held priority, with `Waiting` put back after: a refusal then gives the real reason, from `RuleError.Say`; no refusal means p only waits ("Wait for your priority", or "Answer the prompt first" while p has another prompt open). The turn actions first say why the turn itself stops them ("Not your turn", "Wait for the stack to resolve"). `Say` is set with `deny` at the refusals a player meets; a cost that can't be paid says which (`costSay`). A refusal without words reads "Not usable now". Reasons never name a card, so they leak nothing. The server sends them as `update.blocked`; the engine itself never reads them.

## Counts in views

`PlayerView.LootPlays`, `Attacks` and `Purchases` are what a player may still do now, bonuses included (attacks and purchases only for the active player). `CardView.Counters` copies the object's counters. `CardView.Stats` is set on monsters only: live `HP`, `DC` and `ATK` from `HP`, `Evasion` and `MonsterATK`, and the printed ones from the card definition (the printed DC clamped like `Evasion`). Like `Options`, they are computed at view time.

## Who the stack waits for

In an open priority window, `GameView.Passed` lists the players who passed in a row (first passer first) and `GameView.Responding` lists those who may still respond, the waiting player first (V-11). Priority goes round once in turn order, so both come from `Priority.Holder` and `Priority.Passes`. Outside a window (a choice, a discard) both are empty. The server's automatic passes (N-05, N-06) go through the engine, so skipped players show up in `Passed`.

`GameView.WinSouls` is the soul value that wins (`Game.WinSouls`, R-WIN-01), so clients can show how close each player is (V-30).

## Private events

An event with `Private` shows its card only to `Player` and to the players in `SeenBy`; `FilterEvents` blanks `Card` and `Object` for everyone else (judges see all). `GiveHandCard` sets `SeenBy` to the giver, so a steal shows the card to the thief and the victim only (V-14). `Text` keeps the giver for all viewers.
