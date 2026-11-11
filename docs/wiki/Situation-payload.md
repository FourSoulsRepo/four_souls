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
* `sets` — card sets the situation needs, by code (`b2`) or name (`Base Game V2`).
* `seed` — optional; fixes shuffles and dice.
* `setup.players[]` — `character`, `items`, `hand`, `cents`, `souls`, `damage`, `dead`, and `deactivated` (cards, or `"character"`, that start deactivated).
* `setup.active` — the active player; the turn is in the open action phase with one loot play, attack and purchase left (`loot_plays` overrides the loot plays).
* `setup.monsters`, `setup.shop` — the top card of each slot.
* `setup.bosses`, `setup.minions` — a challenge's final boss and minion slots, one card each (R-CHAL-02, R-CHAL-07).
* `setup.stack` — items already on the stack, bottom first: dice rolls, loot cards, abilities, triggers and attack rolls; see [Stacked items](#stacked-items).
* `action.kind` — `play`, `activate`, `attack`, `purchase`, `end_turn` or `pass`; `card` and `ability` pick the card.
* `action.choices` — answers to the prompts that follow, by option label: an option's plain form, with no language (`Game.Plain`, 16.4). A card is its ref (`"bomb"`), a player `"player 2 (hero_b)"`, a number its digits, a card's own words (a "choose one-" option) the words; any other option is its key, then its values (`"ask.deck.monster"`, `"ask.cancel"`, `"ask.yes"`, `"ask.slot_card 1 fly"`).
* `expect` — optional: `allowed`, the refusal `rule`, and `after` checks per player (`cents`, `hand`, `damage`, `souls`, `dead`).

## Stacked items

`setup.stack` lists what is already on the stack, bottom first (LR-09, 12.5). Each item has an `owner` (a seat, from 0, like `action.player`) and a `kind`:

| `kind` | Item | Fields |
|--------|------|--------|
| none | a dice roll | `dice_roll` (1–6): `{"dice_roll": 4, "owner": 1}`, the first form |
| `play` | a loot card | `card` (a loot card), `choices` |
| `activate` | an activated ability | `card` (in play, the owner's), `ability` (as in `action`), `choices` |
| `trigger` | a triggered ability | `card` (the owner's, or a monster or room in a slot when the owner is the active player), `ability` |
| `attack_roll` | an attack roll | `dice_roll` (1–6), `card` (a monster in play: a monster, final boss or minion slot); the owner is the active player |

`choices` are the mode and targets chosen when the item went on the stack, in the order the prompts ask them, by the same labels as `action.choices`: `"player 2 (maggy)"`, `"bomb"`, a dice roll below it `"ask.roll_of 2"`, a mode its words. A target can be an item lower on the stack. Answers the item asks for while it resolves (an amount, a yes) go in `action.choices`.

Each item is as it is just after it went on the stack: its costs are paid and nothing is paid again. The setup shows what the costs left: a stacked loot card is in no hand, a tapped item is in `deactivated`, spent cents are not in `cents`. Nothing triggers as the items go on the stack. A trigger's event is the owner's ("when you …"). An attack roll means the active player attacks that monster: the attack is declared and used, and every attack roll attacks the same monster.

The items are built with the engine's own stack items for a played loot, an activation, a trigger and an attack roll, so they resolve as in a game. Priority goes to the top item's owner; `{"kind": "pass"}` by them lets the whole stack resolve, top first, and `expect.after` checks the end.

```json
{
  "title": "A roll of 2, Bomb! and Lucky Foot on the stack: which resolves first?",
  "sets": ["b2"],
  "setup": {
    "players": [
      {"character": "isaac", "items": ["lucky_foot"], "deactivated": ["lucky_foot"]},
      {"character": "maggy"}
    ],
    "active": 0,
    "stack": [
      {"dice_roll": 2, "owner": 0},
      {"kind": "play", "card": "bomb", "owner": 0, "choices": ["player 2 (maggy)"]},
      {"kind": "activate", "card": "lucky_foot", "owner": 0, "choices": ["ask.roll_of 2"]}
    ]
  },
  "action": {"player": 0, "kind": "pass", "choices": ["2"]},
  "expect": {"allowed": true, "after": [{"player": 0, "hand": 0}, {"player": 1, "damage": 1}]}
}
```

Lucky Foot resolves first and makes the roll a 4, then Bomb! deals 1 damage to Maggy, then the roll resolves as 4. An unknown field, kind, card, owner or target is an error that names the item and lists the options. A loot card that copies an ability as it is played (? Card) can't start on the stack.

## What the runner does

1. Builds the table; decks hold every card of the sets, shuffled with `seed`.
2. Tries the action. A refusal answers `allowed: false` with a rule ID, an English `reason` for developers, and `say`: the player's words as a message key with values (ADR 021), which each client shows in its language.
3. If allowed, answers the choices and lets every player pass until the stack is empty.
4. With `expect`, compares and lists any mismatch.

```go
ans, err := engine.RunSituation(data, b2.Set)
```

## As tests

Every file in `pkg/rules_engine/testdata/situations/` runs in `go test` and must match its expectation. A disputed case, once resolved, becomes such a file (R-03).
