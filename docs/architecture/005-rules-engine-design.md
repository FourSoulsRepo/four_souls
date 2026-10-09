# 5. Rules engine design

* **Status:** Proposed
* **Date:** 2026-10-09
* **Authors:** @HardDie

---

## Context

1. The engine implements our rules text (`pkg/rules_engine/docs/rules`).
2. It runs on the server, in tests, and in the browser (A-03, A-09).
   1. No `os`, `net`, `time`, `math/rand` (enforced by depguard).
3. It must be deterministic for replays and checksums (A-08, RP-12).
4. Cards are data built from reusable blocks (A-06).
5. Card effects rewrite pending events (A-05; R-ABIL-29 to R-ABIL-33).
6. Every player may respond while holding priority (R-PRIO, N-05).
7. Clients see only their own view (A-07).
8. Lessons from ADR 001: general rules apart from card logic.

## Considered options

1. **Event loop driven by the clock or goroutines**
   1. Lost. Not deterministic; hard to replay and test.
2. **Pure state machine: input in, events and next prompt out**
   1. Won. Same input gives the same result; easy to save and replay.
3. **Card scripts in a text language (Lua, JSON rules)**
   1. Lost. Needs an interpreter; errors appear only at run time.
4. **Cards as Go values built from block constructors**
   1. Won. Checked by the compiler; reads like data; no parser.
5. **Copy-on-write immutable state**
   1. Lost for now. Slower and more code; KISS.
6. **One mutable state, deep copy on demand**
   1. Won. Simple; copies serve views, sandbox and fuzzing.

## Decision

Use options 2, 4 and 6.

1. Shape
   1. One Go package `rulesengine` (imported as `engine`).
   2. Card definitions live in sub-packages per set, e.g. `engine/cards/b2`.
   3. The engine never imports `card_db`; it knows cards by ref (`the_d6`).
2. State
   1. `Game` holds everything; it marshals to JSON and back.
   2. No Go maps in state; ordered slices only (A-08).
   3. Objects live in a slice indexed by object ID.
   4. A card that changes zone becomes a new object ID (R-ZONE-01).
   5. Zones: decks, discards, hands, slots (with covered cards).
   6. Also: in play, the stack, outside the game.
   7. RNG: our own small seeded generator inside the state.
   8. Its seed and position are never sent to clients.
3. Flow
   1. `Apply(intent)` checks an intent against the current prompt.
   2. Invalid intent: error with a reason and rule ID; state unchanged.
   3. Valid intent: the engine runs until it needs input again.
   4. It returns ordered events and the next prompt.
   5. A prompt names who must answer and what they may do.
   6. Prompt kinds: priority, choose target, choose option.
   7. More prompt kinds: order triggers, discard, vote, end-of-game.
4. Priority and the stack
   1. Stack items: loot, activated ability, triggered ability.
   2. More stack items: dice roll, damage, death.
   3. Attack and purchase declarations open priority windows.
   4. Priority order follows R-PRIO-02 and R-PRIO-03.
   5. The top item resolves after all players pass in a row.
5. Pending actions
   1. Everything that changes state is an action value first.
   2. Examples: gain ¢, loot, deal damage, destroy, move card.
   3. Replacement hooks of objects in play may rewrite the action.
   4. Each hook applies once per action (R-ABIL-32).
   5. If several apply, the affected player orders them (R-ABIL-33).
   6. The final action is applied and emits events.
6. Triggers
   1. After each action, triggered abilities look at its events.
   2. Matches wait in a queue until someone would get priority.
   3. They go on the stack in R-ABIL-15 order.
7. Card definitions
   1. A card is a Go value made from blocks.
   2. Blocks: costs, targets, effects, triggers, static rules.
   3. Replacement hooks are blocks too.
   4. A card that does not fit gets a small custom Go hook.
   5. Game data (stats, rewards, souls, copies) is in the definition.
   6. A generator copies that data from `card_db` (step 5.1).
8. Views and events
   1. `View(game, viewer)` builds what a player, viewer or judge sees.
   2. Hidden: other hands, deck order, RNG state.
   3. Events carry a visibility rule and pass the same filter.
   4. Every applied step has a state checksum (FNV-64 of the state).
9. Situations and tests
   1. A situation payload builds a `Game` directly (LR-04).
   2. Its action is applied; the answer is allowed or not, with reasons.
   3. Fuzz tests drive random legal intents with invariants.

```go
// Sketch of the public API; names may change in step 4.
g, err := engine.NewGame(engine.Setup{Seed: 42, Players: 3, Cards: b2.Set})
p := g.Prompt()                // who must act, and the allowed intents
events, err := g.Apply(intent) // validate, run until the next prompt
v := engine.View(g, viewer)    // filtered state for one viewer
sum := g.Checksum()            // RP-12
```

## Consequences

### Positive

1. The server, tests, sandbox and website share one engine.
2. A record of intents and seed replays a game exactly.
3. Card bugs stay inside card definitions, not the core.
4. Invalid moves always come with a reason and a rule ID.

### Negative and risks

1. Prompts for every choice make the engine larger up front.
2. Replacement ordering and simultaneous triggers are complex.
3. Card data exists twice (card_db display, engine game data).
   1. The generator and a test keep them in sync.
4. Deep copies may cost time in fuzzing; measure before optimizing.

### Neutral

1. Exact type names are settled while coding steps 4.2 to 4.12.
2. Each rule implemented gets tests citing its rule ID.
