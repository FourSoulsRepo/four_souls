# Situation payload

A situation is a small JSON document: a table position, one action, and optionally what should happen. The same format is used by the in-app sandbox, the rules website and rules tests (LR-04, LR-07).

```json
{
  "title": "Golden Razor Blade-like item: pay 5 cents, deal 1 damage to a player",
  "sets": ["test", "abilities"],
  "setup": {
    "players": [
      {"character": "hero_a", "items": ["razor"], "cents": 7},
      {"character": "hero_b"}
    ],
    "active": 0,
    "monsters": ["gaper", "gaper"]
  },
  "action": {"player": 0, "kind": "activate", "card": "razor", "choices": ["player 2 (hero_b)"]},
  "expect": {"allowed": true, "after": [{"player": 0, "cents": 2}, {"player": 1, "damage": 1}]}
}
```

## Fields

* `title` — the question in plain words.
* `sets` — card sets the situation needs (e.g. `b2` once the Base Game is implemented).
* `seed` — optional; fixes shuffles and dice.
* `setup.players[]` — `character`, `items`, `hand`, `cents`, `souls`, `damage`, `dead`, and `deactivated` (cards, or `"character"`, that start deactivated).
* `setup.active` — the active player; the turn is in the open action phase with one loot play, attack and purchase left (`loot_plays` overrides the loot plays).
* `setup.monsters`, `setup.shop` — the top card of each slot.
* `setup.stack` — dice rolls already on the stack, bottom first: `{"dice_roll": 4, "owner": 1}`.
* `action.kind` — `play`, `activate`, `attack`, `purchase`, `end_turn` or `pass`; `card` and `ability` pick the card.
* `action.choices` — answers to the prompts that follow, by option label (a target's card name, `"player 2 (hero_b)"`, `"monster deck"`, `"cancel"`).
* `expect` — optional: `allowed`, the refusal `rule`, and `after` checks per player (`cents`, `hand`, `damage`, `souls`, `dead`).

## What the runner does

1. Builds the table; decks hold every card of the sets, shuffled with `seed`.
2. Tries the action. A refusal answers `allowed: false` with a reason and a rule ID.
3. If allowed, answers the choices and lets every player pass until the stack is empty.
4. With `expect`, compares and lists any mismatch.

```go
ans, err := engine.RunSituation(data, b2.Set)
```

## As tests

Every file in `pkg/rules_engine/testdata/situations/` runs in `go test` and must match its expectation. A disputed case, once resolved, becomes such a file (R-03).
