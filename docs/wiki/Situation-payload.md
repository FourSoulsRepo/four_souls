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
* `source` — optional for a payload, required in a test file: a list of what the `expect` values come from (ADR 033, CLAUDE.md 5.25); see [As tests](#as-tests). The engine does not read it.

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

A situation file with `expect` and `source` is a rules test (LR-07, ADR 033). Each package with tests keeps them in its own `testdata/situations/` tree, and its `go test` runs them: one subtest per file, named by its path in the tree without `.json`, e.g. `TestSituations/razor_blade/deals-1-damage-to-a-player`. A disputed case, once resolved, becomes such a file (R-03); it goes where its card or rule goes.

### Layout

Paths are relative to `pkg/rules_engine/`:

| Tree | Folder | Holds |
|------|--------|-------|
| `cards/<set>/testdata/situations/` | `<card>/` | one card's tests; `<card>` is its ref, as in `<card>_test.go` |
| `cards/<set>/testdata/situations/` | `interactions/` | several cards of the set together, named after them: `lucky-foot-bomb.json` |
| `cards/testdata/situations/` | `interactions/` | cards of several sets |
| `testdata/situations/` | `<topic>/` | engine rules with the test sets; `<topic>` is the rule topic in lower case (`R-ATK-01` is `atk/`) |

A folder holds only `.json` files named in kebab case after what they check (`pays-5-cents.json`); no deeper levels. The runner fails a file outside this layout, a card folder that is no card ref of the package's sets, and a topic folder that is no rule topic (the topics in `pkg/rules_engine/docs/rules/README.md`).

### `source`

A non-empty list; each item names one basis for the expected values:

| Form | Example | The runner checks |
|------|---------|-------------------|
| `card <ref>` | `card razor_blade` | the card is in the package's sets; its own text is the basis |
| `R-<TOPIC>-<NN>` | `R-ABIL-07` | the rule is defined in `pkg/rules_engine/docs/rules/` |
| `Q-<NN>` | `Q-07` | the question is a heading of `open-questions.md` |
| `R-07 #<N>` | `R-07 #4` | the row is in the old rulings table of `open-questions.md` |
| `https://…` | a ruling outside the repo | only the form |

Any other form fails with the list of forms. The `expect` values come from these sources, never from running the code (CLAUDE.md 5.25).

### The runner

`enginetest.RunSituations(t, files, sets...)` runs a card tree; `enginetest.RunRuleSituations` the engine's topic tree (`situation_files_test.go`). For each file it checks the layout, then that the file has `expect` and `source` and that each source exists, then runs it with `engine.RunSituation`. A mismatch names the file, the step and every failed check:

```
razor_blade/hits-cain.json, step 1 (activate razor_blade): player 1 damage is 1, expected 2
```

The rule and question IDs come from `pkg/rules_engine/docs` (package `docs`), which embeds the rules text; only tests import it, so the website's wasm stays the same size.

A set package that has files adds one test (`cards/b2/situations_test.go` is the model):

```go
//go:embed testdata/situations
var situations embed.FS

// TestSituations runs this set's situation files (ADR 033).
func TestSituations(t *testing.T) {
	enginetest.RunSituations(t, situations, Set)
}
```

`go:embed` fails on a folder with no files, so a package adds this test with its first file, not before. A cross-set file in `cards/testdata/situations/` passes every set (`cards.Sets()` and the challenge set) the same way.
