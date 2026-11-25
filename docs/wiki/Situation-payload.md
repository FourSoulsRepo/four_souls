# Situation payload

A situation is a small JSON document: a table position, one action or several [steps](#steps), and optionally what should happen. The same format is used by the in-app sandbox, the rules website and rules tests (LR-04, LR-07).

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
* `setup.players[].counters` — counters on the player's cards: `[{"card": "tech_x", "count": 3}]`; `card` is an item or soul ref, or `"character"`; `name` names the counters (`{"card": "bum_friend", "name": "gold", "count": 1}`), none for plain counters. See [Counters](#counters).
* `setup.players[].eternal` — items or souls that are eternal, as a starting item (R-ABIL-18): `["the_bone"]`.
* `setup.players[].extra_loot_plays` — loot plays from abilities that work on any turn, as Cain's ↷ gives on another player's turn: `1`.
* `setup.active` — the active player; the turn is in the open action phase with one loot play, attack and purchase left (`loot_plays`, `attacks` and `purchases` override them; bonus ones, as Champion Belt's, come on top), unless `setup.phase` or `setup.declared` names another window.
* `setup.monsters`, `setup.shop` — the top card of each slot.
* `setup.rooms` — the room in each room slot (R-SETUP-12); the active player uses its abilities (R-CARD-33, R-CARD-36).
* `setup.monster_damage` — damage on a card in a monster, final boss, minion or room slot, by its ref: `{"gemini": 2}`. It must leave the card at least 1 HP.
* `setup.bonus_souls` — the active bonus souls, by ref (R-SETUP-06): `["soul_of_greed"]`. They start with no owner; `setup.counters` and `expect.counters` name them after the slot cards.
* `setup.counters` — counters on a card no player controls: a card in a slot or an active bonus soul, in the same form as a player's: `[{"card": "soul_of_wrath", "count": 5}]`.
* `setup.decks` — cards put on top of decks, by deck name (`treasure`, `loot`, `monster`, `room`), top first: `{"monster": ["chest_2"], "loot": ["bomb", "a_penny"]}`. Each card must belong to that deck. The rest of the deck is every card of the sets, shuffled with `seed`.
* `setup.rolls` — the next dice results, in order, each 1–6: `[6, 1]`. Any die rolled after them is random (from `seed`). In a test file every listed roll must be used: a roll left over at the end is an error (24.32); see [Leftover rolls and picks](#leftover-rolls-and-picks).
* `setup.random` — the next random picks, in order, each named by its option label as in `choices`: `["player 2 (cain)"]`. Any pick after them is random (from `seed`). In a test file every listed pick must be used, as `rolls` (24.32). See [Random picks](#random-picks).
* `setup.bosses`, `setup.minions` — final boss and minion slots, one card each (R-CHAL-02, R-CHAL-07). In a regular game `bosses` is only a slot; a challenge (`setup.game`) puts its own boss out, so `bosses` is then an error.
* `setup.game` — the game type (24.34): `{"type": "regular"}`, the default, or a challenge started on the table, `{"type": "challenge", "id": "krampus", "difficulty": "normal"}`. See [Game type](#game-type).
* `setup.stack` — items already on the stack, bottom first: dice rolls, loot cards, abilities, triggers and attack rolls; see [Stacked items](#stacked-items).
* `setup.priority` — the seat that holds priority at the start, from 0: `1`. See [Priority](#priority).
* `setup.phase` — the turn's window the situation starts in: `start`, `loot`, `action` (the default) or `end` (24.10). See [Windows](#windows).
* `setup.declared` — `attack` or `purchase`: the active player has declared it in the action phase and not chosen the target or item yet (R-ATK-02, R-SHOP-02). See [Windows](#windows).
* `setup.must_attack`, `setup.must_attacks`, `setup.must_attack_deck`, `setup.must_attack_one_of` — the attack duties in force this turn (19.235): a card the active player must attack (`"krampus"`), how many more attacks they must make (`1`), and how many of their next attacks must go to the monster deck (`1`). `setup.must_attack_one_of` names the cards an extra attack is owed on one of (`["monstro", "widow"]`, 19.241). A card in a slot does not set its duty; write the duty here if it is in force. See [Attack duties](#attack-duties).
* `action.kind` — `play`, `activate`, `attack`, `purchase`, `end_turn` or `pass`; `card` and `ability` pick the card. `activate` names a card the player controls, or a room in a room slot (24.6): only the active player may use a room, so another player's `activate` is refused by R-CARD-33, as in a game. Also `throw`, `choose` and `reveal`, which mostly make sense as [steps](#steps).
* `action.choices` — answers to the prompts that follow, by option label: an option's plain form, with no language (`Game.Plain`, 16.4). A card is its ref (`"bomb"`), a player `"player 2 (hero_b)"`, a number its digits, a card's own words (a "choose one-" option) the words; any other option is its key, then its values (`"ask.deck.monster"`, `"ask.cancel"`, `"ask.yes"`, `"ask.slot_card 1 fly"`).
* `steps` — several actions in a row, each with its own `choices` and an optional `expect`; a situation has `action` or `steps`, never both. See [Steps](#steps).
* `expect` — optional: `allowed`, the refusal `rule`, `after` checks per player, and checks of monsters, decks, the turn, the game's end, the open question, the resolve order and events; see [Checks](#checks). `allowed` is a plain boolean: leaving it out expects a refusal.
* `source` — optional for a payload, required in a test file: a list of what the `expect` values come from (ADR 034, CLAUDE.md 5.25); see [As tests](#as-tests). The engine does not read it.

## Game type

`setup.game` says what kind of game the table is (24.34). Leaving it out is a regular game, the same as `{"type": "regular"}`.

| `type` | Game | Fields |
|---|---|---|
| `regular` | a regular game (the default) | none |
| `challenge` | a challenge, started on the built table as a real game starts it (R-CHAL-01) | `id`, `difficulty`; for its state: `presents`, `holidays`, `rule_card` |
| `team` | refused: team games are not in the engine yet (Q-03, parked) | — |

```json
{
  "title": "Krampus challenge, Normal: the players hold 4 souls between them; may Isaac attack Pilfering Krampus?",
  "sets": ["b2", "challenge"],
  "setup": {
    "game": {"type": "challenge", "id": "krampus", "difficulty": "normal"},
    "players": [{"character": "isaac", "souls": ["monstro", "gemini", "chub"]}, {"character": "cain", "souls": ["carrion_queen"]}],
    "active": 0,
    "monsters": ["leech", "fly"],
    "decks": {"treasure": ["boomerang", "breakfast", "box", "blank_card", "book_of_sin"]}
  },
  "steps": [{"player": 0, "kind": "attack", "expect": {"allowed": true, "prompt": {"player": 0, "options": ["leech", "fly", "ask.deck.monster"]}}}],
  "expect": {"allowed": true}
}
```

A challenge starts after the table is built, as `Game.StartChallenge` does: the timer at 8 (R-COOP-02), the rule card of the difficulty outside the game, then the challenge's own setup. Rounds count from the active player (R-TURN-15).

* `id`: `krampus` (How the Krampus Stole Christmas) or `trick_or_treat`; `difficulty`: `normal`, `hard` or `ultra`. Each is an error that lists the options. The challenges come with the card set `challenge` (`CardSet.Challenges`), so `sets` names it.
* Krampus: Pilfering Krampus of the difficulty in a final boss slot (R-KRAMPUS-04); the top 4 treasure cards become the ¢, loot, treasure and soul presents, in that order, with their counters (R-KRAMPUS-05). Name them with `setup.decks`: `{"treasure": ["boomerang", "breakfast", "box", "blank_card"]}`.
* Trick or Treat: Mischievous The Haunt of the difficulty in a final boss slot with its trick counters; the Trick/Treat card Trick side up (R-TRICK-04).
* Exactly 2 players (R-COOP-01). No `bosses` (the challenge puts its boss out, R-CHAL-02) and no `bonus_souls` (R-COOP-02).
* `setup.counters` on the final boss sets the counters it has left, in place of its start: `{"card": "mischievous_the_haunt_norm", "name": "trick", "count": 3}`; `count` 0 leaves none.
* A field of one type only is an error with another: `id fits only type challenge, not regular`.

The challenge's state, when the test needs it:

| Field | Says | Example |
|---|---|---|
| `presents` | a present by `kind` (`cents`, `loot`, `treasure`, `soul`): `counters` left (1 or more), `destroyed` (in the treasure discard, R-KRAMPUS-08), or `to`, the seat that gained it at Happy Holidays! (R-KRAMPUS-10). A present not named keeps its start | `[{"kind": "cents", "counters": 2}, {"kind": "loot", "destroyed": true}]` |
| `holidays` | Happy Holidays! has happened (R-KRAMPUS-10): each present not destroyed is a charged item of the player its `to` names, so every one is named | `true` |
| `rule_card` | the rule card's side up, by ref: `treat_normhardultra` is the Trick/Treat card flipped (R-TRICK-11) | `"treat_normhardultra"` |

Nothing triggers as the state is set up, as with every setup field. A Krampus challenge has presents and Trick or Treat none, so `presents` and `holidays` there are errors.

## Random picks

"A player chosen at random" is a pick everyone at a table sees, so a file names it, as `rolls` names the dice (24.30). `setup.random` lists the next picks in order; each is the plain label of the picked option, the form `choices` uses (`"player 3 (samson)"`, a card's ref for a random attack's target).

```json
{
  "title": "Cursed Mulliboom in a slot: Isaac kills the Fly; the player chosen at random is Cain",
  "sets": ["b2", "r"],
  "setup": {
    "players": [{"character": "isaac"}, {"character": "cain"}, {"character": "samson"}],
    "active": 0,
    "monsters": ["cursed_mulliboom", "fly"],
    "rolls": [6],
    "random": ["player 2 (cain)"]
  },
  "steps": [{"player": 0, "kind": "attack", "choices": ["fly"]}],
  "expect": {"allowed": true, "after": [{"player": 1, "dead": true}, {"player": 2, "damage": 0}]}
}
```

* Every random pick of the engine reads the list, in one place (`Game.randomPick`): a question that picks at random (`Question.Random`: Monstro's Tooth, Ring of Flies, Boomerang, Betrayal, Corny, Cursed Mulliboom, Hornfel) and a random attack's target (Pilfering Krampus, Q-31).
* A pick that is no option fails the run with the options: `situation: step 1 (attack fly): random: pick 1 "player 4 (eden)" is not an option of ["player 1 (isaac)" "player 2 (cain)" "player 3 (samson)"]`. An empty pick is a setup error.
* A pick of a hidden card (Boomerang, Ring of Flies: "steal a loot card at random") has an empty label, so the list can't name it; it stays random from `seed`.
* Shuffles and random deck order are no picks: decks come from `seed` and `setup.decks`, cards put on the bottom "in random order" are checked with `deck_bottoms` (Q-22).
* Only a setup names picks; a game never does, so it draws every pick from its RNG as before.

## Leftover rolls and picks

A roll or pick listed in the setup but never used means the setup is likely wrong: a roll that was meant to happen did not (24.32). So a test file fails when, after its last step, `setup.rolls` or `setup.random` still holds entries. One check covers both fields and names each with what was left:

```
conjoined_fatty/kills-isaac-and-ends-the-attack.json, end: setup.rolls: 4 left unused: [6, 6, 6, 6]
razor_blade/hits-cain.json, end: setup.random: 1 left unused: ["player 2 (cain)"]
```

* It fires once, at the end, after every step and the top-level `expect`; it comes after the mismatches and does not stop them.
* A refused or stopped step uses no dice, so a roll it would have needed is left over too: list only what the run throws.
* To show a roll does *not* happen (an attack that ends), list only the rolls that do and check their count with `resolved_rolls`.
* Only the test runner (`enginetest.CheckSituation`) fails on leftovers, as it alone reads a file strictly. A payload from the sandbox or the website is only being explored, so leftovers are no error there: the engine keeps them in `Answer.Unused`, which is never sent (`json:"-"`).

## Counters

A counter entry is `{"card": …, "name": …, "count": …}`; `count` is 1 or more and `name` is optional. A ref names the first such card: in seat order for a player's cards, then in slot order (monster, final boss, minion, room slots) and the bonus souls for `setup.counters`. An unknown field, a card the player has not in play, or a card in no slot is an error that lists the options.

In `setup.players[].counters` and `after[].counters`, `"card": "player"` puts the counters on the player, not a card (R-MECH-13): Trick or Treat's fear, `{"card": "player", "name": "fear", "count": 2}` (R-TRICK-08). They stay when the character changes or dies.

## Checks

`expect` checks what is true once the action, or the step it sits on, has settled (24.5). A field left out is not checked, so old payloads keep their meaning. Names follow `setup`, so a setup and an expect read alike.

| Field | Checks | Example |
|-------|--------|---------|
| `after[].cents`, `damage`, `dead` | the player's cents, damage, dead | `{"player": 1, "damage": 1}` |
| `after[].hand`, `souls` | counts: cards in hand, soul value | `{"player": 0, "hand": 2, "souls": 1}` |
| `after[].hp`, `atk` | HP left (max HP minus damage), ATK | `{"player": 0, "hp": 1, "atk": 2}` |
| `after[].items` | every item the player has in play | `{"player": 1, "items": []}` |
| `after[].hand_cards` | every card in hand | `{"player": 0, "hand_cards": ["bomb"]}` |
| `after[].soul_cards` | every soul card in play | `{"player": 0, "soul_cards": ["soul_of_wrath"]}` |
| `after[].deactivated` | every deactivated item, and `"character"` | `{"player": 0, "deactivated": ["jawbone"]}` |
| `after[].counters` | counters on the player's cards, as in setup; `"card": "player"` on the player | `{"player": 0, "counters": [{"card": "sumptorium", "count": 0}]}` |
| `monsters[]` | a card in a slot: `card`, and `slot`, `damage`, `hp`, `atk`, `evasion`, `dead` | `{"slot": 1, "card": "fly", "damage": 1, "hp": 1}` |
| `counters` | counters on cards no player controls, as `setup.counters` | `[{"card": "gaper", "name": "level", "count": 1}]` |
| `decks` | the cards on top of a deck, top first, as `setup.decks` | `{"monster": ["fly"]}` |
| `discards` | every card in a discard pile | `{"loot": ["a_dime"]}` |
| `deck_bottoms` | the cards at the bottom of a deck, in any order | `{"loot": ["a_penny", "bomb"]}` |
| `active`, `loot_plays`, `attacks` | the active seat, and the loot plays and attacks it can still use, bonus ones too | `"attacks": 2` |
| `over`, `winners` | the game has ended; the winning seats | `"over": true, "winners": [0]` |
| `prompt` | the open question: `player` who is asked, every option in `options` | `{"player": 0, "options": ["the_d6", "ask.cancel"]}` |
| `resolved` | stack items in the order they resolved; others may come between | `["lucky_foot", "bomb", "dice_roll"]` |
| `events` | event kinds that happened, with `player`; `none: true` that did not | `[{"kind": "damaged", "player": 1}, {"kind": "died", "none": true}]` |
| `resolved_rolls` | the value each roll resolved as, in order | `[5]` |
| `presents` | a challenge's presents by `kind`: `counters` (0 once gone), `destroyed`, `to` (the seat that has it as an item) | `[{"kind": "cents", "counters": 9}, {"kind": "loot", "destroyed": true}]` |

* A list of cards (`items`, `hand_cards`, `soul_cards`, `deactivated`, `discards`, `winners`) names every card there, in any order: a card more or less is a mismatch, and `[]` expects none. The deck in `decks` is the exception: it holds the rest of the sets, so only its top cards are compared, in order.
* `hand` and `souls` stay counts, as before 24.5; the cards are `hand_cards` and `soul_cards`.
* A counter check names its card as in setup (`"character"` for the character); a `count` of 0 expects none.
* A `monsters` check needs `card`. Without `slot`, `card` names the first such card in slot order (monster, final boss, minion, room slots), as `setup.monster_damage` does. With `slot` (from 1, a monster slot), that slot's top card must be `card`. Every entry without `dead: true` checks that the card is on top of a slot (that slot, with `slot`), even `{"card": "fly"}` alone; otherwise the mismatch says it is in no slot and lists the slot cards (19.234). A covered card (R-ZONE-10) is in no slot: it is not in play. `dead: false` therefore says the same as leaving `dead` out. `dead: true` means no slot holds the card and it is in a discard or a player's souls.
* Nothing checks nothing (19.234). An `after` entry with only `player` is a mismatch; so is a list or map of checks that is there but empty: `after`, `monsters`, `counters` (top-level or a player's), `events`, `decks`, `discards`, `deck_bottoms`. Leave such a field out instead. Inside `decks` and `deck_bottoms`, `[]` for one deck expects that deck to be empty; in `discards` it expects an empty discard pile.
* An `events` entry's `player` must be a seat in the game; any other seat is a mismatch, since `none: true` would always pass for it.
* `hp` is what is left: max HP minus damage, for players and monsters.
* `attacks` and `loot_plays` (24.28) count what the active player can still do, as the engine does (`attacksLeft`, `lootPlaysFor`). `attacks` is the turn's attacks left plus the unused bonus attacks from cards in effect (Champion Belt, Might For The Meek, Skeleton Key). `loot_plays` is the turn's loot plays left, plus the unused bonus loot plays (Polydactyly, Belly Button), plus the active player's `extra_loot_plays`. So Champion Belt checks `"attacks": 1` after the first attack and `0` after the second. The setup fields keep their meaning: `setup.attacks` and `setup.loot_plays` are the turn's own counts, and the bonus ones come on top; with Champion Belt, `"attacks": 0` in setup still leaves one attack, so `expect.attacks` reads 1.
* `prompt` (24.6) checks the question open once the step settled, a choose or order question. `player` is the seat asked; `options` names every option by its label, as `choices` do, in any order. A target or mode question ends with `ask.cancel`: the player may back out before paying (R-ABIL-04). A field left out is not checked, so `"prompt": {}` only says a question is open. A prompt check also stops the step at that question; see [When a step settles](#when-a-step-settles).
* `resolved` (24.6) names stack items in the order they resolved (`stack_resolved` events, as `enginetest.ResolvedAt` reads them). It is an ordered subsequence: each name must resolve after the one before it, and other items may resolve between them, so a check of two cards ignores the damage a Bomb! puts on the stack. A loot card, an activated or a triggered ability is named by its card ref; any other item by its kind: `dice_roll`, `attack_roll`, `damage`, `death`, `death_penalty`, `death_step`. `[]` expects nothing resolved. Any other name is an error that lists the kinds.
* `events` (24.6) checks event kinds, by the engine's names (`damaged`, `looked_at`, `turn_ended_early`, …; `stack_resolved` and the rest). Each entry needs `kind`; `player` (a seat) narrows it to that player's events, and without it any player's count. `none: true` expects no such event; without it, at least one. An unknown kind is an error that lists every kind.
* A step's `resolved` and `events` read that step's events; the top-level `expect` reads every step's (the answer's `events`).
* `deck_bottoms` (24.11) checks the bottom of a deck, by deck name as `decks`: as many cards as the list names, from the bottom up, compared in any order. Several cards put on the bottom go there in random order (Q-22), so no rule says their order; one card at the bottom is simply that card. The rest of the deck is the shuffled set cards, so the list names only cards an effect put there. `decks` still compares the top, top first.
* `resolved_rolls` (24.11) lists the value every dice roll resolved as, in order: the final result after every change (Lucky Foot, I. The Magician, Book of Belial, continuous roll bonuses) and the 1 to 6 bounds (R-DICE-01, R-DICE-06), as its `roll_resolved` event says. The list holds every roll, attack rolls too, so its length counts them; `[]` expects none. The die thrown is `setup.rolls`; this is what that die became. A step's list reads that step's rolls; the top-level one every step's.

Each mismatch says what it found and what was expected:

```
player 0 items are [razor sword coin_bag], expected [razor sword]
monster slot 1 is "gaper", expected "boss"
monster gaper evasion is 3, expected 4
monster fly is in no slot, options [gaper boss]
player 0 counters on razor are 2, expected 3
monster deck top is [gaper], expected [boss]
loot deck bottom is [bomb_loot], expected [penny] in any order
resolved_rolls are [4], expected [5]
winners are [], expected [0]
prompt options are ["the_d6" "ask.cancel"], expected ["the_d6"]
resolved [lucky_foot bomb damage dice_roll], expected [bomb lucky_foot] in this order
events: no "looked_at" for player 1, expected one
events: 1 "damaged" for player 1, expected none
```

## Stacked items

`setup.stack` lists what is already on the stack, bottom first (LR-09, 12.5). Each item has an `owner` (a seat, from 0, like `action.player`) and a `kind`:

| `kind` | Item | Fields |
|--------|------|--------|
| none | a dice roll | `dice_roll` (1–6): `{"dice_roll": 4, "owner": 1}`, the first form |
| `play` | a loot card | `card` (a loot card), `choices` |
| `activate` | an activated ability | `card` (in play, the owner's, or a room in a room slot when the owner is the active player), `ability` (as in `action`), `choices` |
| `trigger` | a triggered ability | `card` (the owner's, or a monster or room in a slot when the owner is the active player), `ability` |
| `attack_roll` | an attack roll | `dice_roll` (1–6), `card` (a monster in play: a monster, final boss or minion slot); the owner is the active player. `from_deck` and `slot`: the monster was revealed by an attack on the monster deck (see [A deck attack's roll](#a-deck-attacks-roll)) |

`choices` are the mode and targets chosen when the item went on the stack, in the order the prompts ask them, by the same labels as `action.choices`: `"player 2 (maggy)"`, `"bomb"`, a dice roll below it `"ask.roll_of 2"`, a mode its words. A target can be an item lower on the stack. Answers the item asks for while it resolves (an amount, a yes) go in `action.choices`.

A room's ability works as an item's (24.6): the room's "you" is the active player (R-CARD-36), and only the active player uses its abilities (R-CARD-33), so its owner must be the active player, who activated it; another owner is an error that names the rule.

Each item is as it is just after it went on the stack: its costs are paid and nothing is paid again. The setup shows what the costs left: a stacked loot card is in no hand, a tapped item is in `deactivated`, spent cents are not in `cents`. Nothing triggers as the items go on the stack. A trigger's event is the owner's ("when you …"). An attack roll means the active player attacks that monster: the attack is declared and used, and every attack roll attacks the same monster. The first one spends the attack as a declaration does (`spendAttack`, 24.29): the turn's own attack first, then a bonus one (Champion Belt, Might For The Meek). So with `"attacks": 0` and Champion Belt the roll uses the Belt's attack and `expect.attacks` reads 0; with no attack left at all it is an error naming R-TURN-07. The attack duties are what is left after it (see [Attack duties](#attack-duties)). A roll of an attack on the monster deck says `from_deck` (see below).

The items are built with the engine's own stack items for a played loot, an activation, a trigger and an attack roll, so they resolve as in a game. Priority goes to the top item's owner (or the seat in `setup.priority`); `{"kind": "pass"}` by them lets the whole stack resolve, top first, and `expect.after` checks the end.

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
  "expect": {"allowed": true, "resolved": ["lucky_foot", "bomb", "dice_roll"], "after": [{"player": 0, "hand": 0}, {"player": 1, "damage": 1}]}
}
```

Lucky Foot resolves first and makes the roll a 4, then Bomb! deals 1 damage to Maggy, then the roll resolves as 4. An unknown field, kind, card, owner or target is an error that names the item and lists the options. A loot card that copies an ability as it is played (? Card) can't start on the stack.

### A deck attack's roll

An `attack_roll` with `"from_deck": true` is the roll of an attack on the monster deck (24.33). Its `card` was revealed from the deck and put on monster slot `slot` (from 1), covering the card there, as the attacker chose (R-ATK-08, R-ZONE-10). The setup builds what the game builds once the revealed card is placed and the first roll waits:

* The revealed card is a new card on top of the slot; the slot's `monsters` card is covered under it. When the revealed card leaves (killed, put on the deck bottom), the covered card comes back on top (R-ZONE-11).
* The attack is on the revealed card, which counts as an attack on a monster (R-ATK-05, Q-24): "cancel an attack on a monster" (Emergency Contact, IBS) works on it. An ability that looks for an attack *declared* on a monster does not (R-ATK-05): the attack was declared on the deck.
* It spends one attack as any stacked attack roll (`spendAttack`, 24.29). A setup has no field for deck-only attacks ("attack the monster deck an additional time", `Turn.DeckAttacks`): such an attack is one more in `attacks`, and the roll spends it like a regular one.
* `must_attack_deck` is what is left after this attack, as `must_attacks` is: write 0 (or leave it out) when the duty was for this one attack.
* Nothing triggers as it is placed: what its entry triggered ("when an attack is declared on this", R-ATK-06; Angelic Intervention) has resolved before the first roll (R-ATK-25).
* A later roll of the same attack names the revealed card without `from_deck` and `slot`; it attacks the same copy, even with another copy in a slot.

```json
"setup": {
  "players": [{"character": "the_soiled", "items": ["ibs"], "hand": ["bomb"]}, {"character": "cain"}],
  "active": 0,
  "monsters": ["leech", "fly"],
  "stack": [{"kind": "attack_roll", "dice_roll": 1, "card": "pooter", "owner": 0, "from_deck": true, "slot": 1}]
}
```

The Soiled attacked the deck; the Pooter covers the Leech in slot 1, and her roll of 1 waits (`cards/r/testdata/situations/ibs/cancels-an-attack-on-a-monster-revealed-from-the-deck.json`).

Each of these is an error that names the item and the rule:

* `from_deck` without `slot`, or `slot` without `from_deck`; either on a kind other than `attack_roll`.
* A `slot` that is no monster slot: the error lists 1 to the slot count.
* A `card` that is no card of the monster deck (R-ATK-08); an event, which ends the attack when revealed (R-ATK-09); a card that can't be attacked, which cancels it (R-ATK-10).
* `from_deck` on a later roll of the attack: a deck attack reveals one card (R-ATK-08).
* No attack left (R-TURN-07), as for any attack roll.

## Priority

`setup.priority` names the seat that holds priority when the situation starts (24.4). Without it, priority goes to the top stacked item's owner, or to the active player.

* With an empty stack, the players before that seat in turn order, counted from the active player, have already passed (R-PRIO-02). The engine counts those passes, so when the seat and the players after it pass too, the window closes as in a game (R-PRIO-05) and the game moves on (R-STACK-07); see [Windows](#windows).
* In the open action phase (no `phase` or `"phase": "action"`, nothing declared, an empty stack) only the active player may hold priority: they can't pass there (R-TURN-09), so no other player ever gets it. Another seat is an error that names R-TURN-09; the active player's seat is fine (it is the default). Start the situation in a real window instead: `phase` `start`, `loot` or `end`, a `declared` attack or purchase, or a stack.
* With a stack, it overrides the top item's owner; nobody has passed yet. This works in every window.
* A seat outside the players is an error.

So a player who may act on another player's turn needs no dice roll on the stack just to get priority:

```json
{
  "title": "Isaac's loot step: Cain used his tap, so with priority he plays A Penny!",
  "source": ["card cain", "card a_penny", "R-TURN-04", "R-PRIO-02", "R-CARD-08"],
  "sets": ["b2"],
  "setup": {
    "players": [
      {"character": "isaac"},
      {"character": "cain", "hand": ["a_penny"], "deactivated": ["character"], "extra_loot_plays": 1}
    ],
    "active": 0,
    "phase": "loot",
    "priority": 1
  },
  "action": {"player": 1, "kind": "play", "card": "a_penny"},
  "expect": {"allowed": true, "after": [{"player": 1, "cents": 1, "hand": 0, "deactivated": ["character"]}]}
}
```

## Windows

A real game has an empty stack with priority passing around in several moments (24.10). `setup.phase` and `setup.declared` start the situation in one of them; the step's own work is already done, and the setup shows the table as it is then.

| Setup | The moment | Turn counts | When everyone passes |
|-------|------------|-------------|----------------------|
| `"phase": "start"` | the start phase after the recharge step and the start-of-turn triggers (R-TURN-02, R-TURN-03) | no loot play, attack or purchase yet | the loot step: the active player loots 1 (R-TURN-04), then its window |
| `"phase": "loot"` | the loot step after the active player looted 1 (R-TURN-04) | none yet | the action phase begins: 1 loot play, attack and purchase (R-TURN-05, R-TURN-07) |
| `"phase": "action"` or none | the open action phase (R-TURN-06) | 1 of each, or `loot_plays`, `attacks`, `purchases` | the active player can't pass (R-TURN-09) |
| `"phase": "end"` | the end phase after the end-of-turn triggers (R-TURN-10); the turn's end is declared | as in the action phase | the hand size step, the room step, then the next player's turn (R-TURN-11 to R-TURN-13) |
| `"declared": "attack"` | the action phase: the attack is declared, the target not chosen (R-ATK-02) | the attack is used (R-TURN-07) | the active player picks the target; a `setup.stack` of responses resolves first |
| `"declared": "purchase"` | the action phase: the purchase is declared, the item not chosen (R-SHOP-02) | the purchase is used (R-SHOP-05) | the active player picks the shop item or the treasure deck; a `setup.stack` of responses resolves first |

* Each window opens with the active player (R-PRIO-02), or with `setup.priority`'s seat after the players before it passed.
* As with stacked items, nothing triggers and nothing is paid on the way in: the start-of-turn or end-of-turn triggers and the declaration's triggers are done, and a cost such as an attack tax shows in the setup's cents.
* The game moves on as in a game (R-STACK-07). An `action` stops at the first window with an empty stack, as before; a step settles to the open action phase, so from `start` or `loot` it runs into this turn's action phase, from `end` into the next player's, and from a declaration through the target or item question (answer it with the step's `choices` or a following `choose` step).
* In `start` and `loot`, `loot_plays`, `attacks` and `purchases` are errors: the action phase gives them (R-TURN-05, R-TURN-07). The active player's `extra_loot_plays` still count.
* `declared` needs the action phase (R-ATK-01, R-SHOP-01) and an attack or purchase left (R-TURN-07, R-SHOP-05). An `attack_roll` stack item needs the action phase (R-ATK-01) and an attack left, bonus ones too (R-TURN-07). Each is an error that names the rule; an unknown `phase` or `declared` lists the options.
* `declared` may come with a `setup.stack`: the responses made to the declaration, waiting on the stack (24.10 follow-up). The declaration itself needed an empty stack; the responses came after it. Priority goes to the top item's owner, or to `setup.priority`'s seat, as with any stack. The stack resolves first (R-STACK-03), priority passes again after each item (R-STACK-04), and once everyone passes on the empty stack the active player picks the target or item (R-ATK-02, R-SHOP-02), as in a game.
* An `attack_roll` stack item with `declared` is an error: with `attack`, the target is chosen only after the declaration's window (R-ATK-02) and attack rolls go against the target (R-ATK-11); with `purchase`, a purchase is not declared during an attack (R-SHOP-01).

```json
{
  "title": "Isaac declared an attack: before he picks the target, Maggy uses Monster Manual on the Fly, so he must attack it",
  "source": ["card monster_manual", "card fly", "R-ATK-02", "R-ATK-03", "R-ATK-04", "R-ATK-12", "R-TURN-07", "R-PRIO-02"],
  "sets": ["b2"],
  "setup": {
    "players": [{"character": "isaac"}, {"character": "maggy", "items": ["monster_manual"]}, {"character": "cain"}],
    "active": 0,
    "monsters": ["clotty", "fly"],
    "declared": "attack",
    "priority": 1,
    "rolls": [6]
  },
  "steps": [
    {"player": 1, "kind": "activate", "card": "monster_manual", "choices": ["fly"],
     "expect": {"allowed": true, "attacks": 0, "prompt": {"player": 0, "options": ["fly"]}}},
    {"kind": "choose", "choices": ["fly"],
     "expect": {"allowed": true, "monsters": [{"card": "fly", "dead": true}], "after": [{"player": 0, "cents": 1}]}}
  ],
  "expect": {"allowed": true}
}
```

A response already on the stack (`cards/b2/testdata/situations/bomb/kills-the-fly-before-the-attack-target.json`): Cain played Bomb! on the Fly after Isaac declared. Cain's pass lets it resolve; the Fly dies (Isaac, the active player, gets its 1¢, R-CARD-13), the Dip refills the slot (R-SHOP-07), and only then does Isaac pick the target:

```json
{
  "setup": {
    "players": [{"character": "isaac"}, {"character": "cain"}],
    "active": 0,
    "monsters": ["clotty", "fly"],
    "decks": {"monster": ["dip"]},
    "declared": "attack",
    "stack": [{"kind": "play", "card": "bomb", "owner": 1, "choices": ["fly"]}],
    "rolls": [6, 6]
  },
  "steps": [
    {"player": 1, "kind": "pass",
     "expect": {"allowed": true, "resolved": ["bomb"], "prompt": {"player": 0, "options": ["clotty", "dip", "ask.deck.monster"]}}},
    {"kind": "choose", "choices": ["clotty"], "expect": {"allowed": true, "after": [{"player": 0, "cents": 5}]}}
  ]
}
```

## Attack duties

A card in a slot does not set its duty; write the duty in the setup if it is in force (owner, 2026-11-23). Krampus in a monster slot says "The active player must attack this once each turn if able", but a setup runs no start-of-turn trigger, and a player may have stopped the duty (with an eternal card's ability, say). So without these fields no duty is in force, as before 19.235. They are the engine's own turn state (`Turn.MustAttack`, `MustAttacks`, `MustAttackDeck`, `MustAttackOneOfCards`), so the game treats them as if a card had set them:

| Field | Means | As set by |
|-------|-------|-----------|
| `must_attack` | a card the active player must attack this turn if able (R-ATK-03, R-ATK-22): the attack's target question offers only it, and the turn can't end while it is in play and an attack is left | Krampus, The Harbingers, Monster Manual; Friendly Ball and Lust for Blood before 19.241 |
| `must_attacks` | how many more attacks the active player must make; the turn can't end while one is owed and an attack is left (R-ATK-03). It adds no attack: `attacks` says how many are left | Portal, Conquest, Envy, Ambush, Gutted Fatty (each also gives the attacks, so `attacks` goes up too), Blood Lust (`r`), Curse of Blood Lust, Blind Rage |
| `must_attack_one_of` | cards the active player owes one more attack on one of (19.241): while a card with one of these refs can be attacked, the attack's target question offers only those; a covered one counts once it comes up (owner, 2026-11-24) (with `must_attack`'s card too), and the turn can't end while an attack is left. An attack on one meets it; with none left, the duty ends. It adds no attack: `attacks` says how many are left | Boss Rush! (both bosses), Friendly Ball, Lust for Blood! (a player's character); each also gives the attack (`ForceAttackOn`) |
| `must_attack_deck` | how many of the active player's next attacks must go to the monster deck: each one's target question offers only the deck. It forces no attack by itself; with Blind Rage or Ambush `must_attacks` says that | Blind Rage (1, with `must_attacks` 1), Ambush (2, with `must_attacks` 2) |

```json
"setup": {
  "players": [{"character": "isaac"}, {"character": "cain"}],
  "active": 0,
  "monsters": ["krampus", "fly"],
  "must_attack": "krampus"
}
```

Isaac can't end his turn (refused by R-ATK-03), and his attack's only target is Krampus (`cards/fsp2/testdata/situations/krampus/isaac-must-attack-it-this-turn.json`).

* `must_attack` names a card by ref, the first with that ref the active player may attack now, in the engine's order: monster, final boss and minion slots, attackable items, attackable rooms, then the other players' characters (only while an effect lets the active player attack players, R-ATK-20). Any other ref is an error that lists those options. Lust for Blood's own DC for its target is no setup field, so a player is a target only with a card in play that makes players attackable.
* `must_attack_one_of` names cards, not objects (19.241): each ref is a card the active player may attack now, found the same way. Any other ref, or a ref named twice, is an error. In a game a covered card named this way counts once it comes up (R-ZONE-01, owner 2026-11-24), but a setup covers a card only under a deck attack's revealed card (`from_deck`), which can't be attacked now, so a dump leaves one out, as it does one gone to a discard, a deck or souls (`cards/fsp2/testdata/situations/boss_rush/the-duty-as-a-setup-field.json`, `the-covered-boss-still-counts.json`).
* The values are the turn as the situation starts: with `declared: attack` or an `attack_roll` stack item the attack has already counted, so `must_attacks` is what is left after it. An `attack_roll` stack item on the `must_attack` card is an error naming R-ATK-22: that attack already met the duty.
* `phase` `start` and `loot` may carry them: the start-of-turn triggers that set such duties are done (R-TURN-03), and they hold into the action phase. In `phase: end` they are an error naming R-ATK-03: the turn's end is declared, which R-ATK-03 refuses while an attack is owed.
* A negative `must_attacks` or `must_attack_deck` is an error.
* No field for Pilfering Krampus's random minion attack (R-KRAMPUS-13, Q-31): it is a "while a minion is in play" effect, not a duty a trigger sets, and the engine keeps only whether it was made this turn (`Turn.RandomAttackMade`); a challenge test that needs that stays Go (24.3).

## Steps

`steps` plays several actions in a row (24.4). Each step has the fields of `action` (`player`, `kind`, `card`, `ability`, `choices`) and an optional `expect`, checked right after that step. The top-level `expect` is checked at the end.

```json
{
  "title": "Eden's Blessing: Isaac kills a fly for 1¢, then ends his turn and gains nothing",
  "source": ["card edens_blessing", "card fly", "R-ATK-12"],
  "sets": ["b2"],
  "setup": {
    "players": [{"character": "isaac", "items": ["edens_blessing"]}, {"character": "cain"}],
    "active": 0,
    "monsters": ["fly"],
    "decks": {"monster": ["leech"]},
    "rolls": [2]
  },
  "steps": [
    {"player": 0, "kind": "attack", "choices": ["fly"], "expect": {"allowed": true, "after": [{"player": 0, "cents": 0}]}},
    {"kind": "throw", "expect": {"allowed": true, "after": [{"player": 0, "cents": 1, "damage": 0}]}},
    {"player": 0, "kind": "end_turn"}
  ],
  "expect": {"allowed": true, "after": [{"player": 0, "cents": 1}]}
}
```

Step kinds, besides those of `action`:

| `kind` | Does | As in `enginetest` |
|--------|------|--------------------|
| `throw` | throws the die of the roll that waits to be thrown | `ThrowPending` (one roll) |
| `choose` | answers the open question with `choices[0]`; the rest answer the questions after it | `Choose` |
| `reveal` | puts `card` (a monster deck card) on top of the monster deck, and `player` attacks the deck, placing it in monster slot 1; `choices` answer what follows | `RevealFromDeck` |

`throw` and `choose` act for the player the game waits on; they do not read `player`. A `throw` with no roll waiting, or a `choose` with no open question or no `choices`, is an error, as is an unknown kind.

### When a step settles

After each allowed step the runner settles: it answers the questions with that step's `choices`, in order, throws every roll, and lets every player pass until the active player's action phase is open again: an empty stack, no attack or purchase going on, the turn's end not declared. Whoever holds priority then keeps it, so a `pass` step hands priority to the next player. So an `end_turn` step runs through the end phase into the next player's action phase, as `enginetest`'s `EndTurn` does; a step in a [window](#windows) runs on the same way.

Two exceptions look one step ahead, and only one:

* If the next step is `choose`, settling stops at a question when the step's `choices` are used up. Without a following `choose`, that is an error that lists the options.
* If the next step is `throw`, settling stops at the first roll waiting to be thrown.

A `prompt` check stops at its question too (24.6). A step whose `expect` has `prompt` settles as if a `choose` step followed: it answers its `choices`, then stops at the first question left, and the check reads it. The top-level `expect.prompt` does the same for the last step, and for `action`. A following step must then answer the question (`choose`) or the file ends there. So a file can check what a player is asked, then answer it:

```json
"steps": [
  {"player": 0, "kind": "activate", "card": "the_battery",
   "expect": {"allowed": true, "prompt": {"player": 0, "options": ["the_d6", "ask.cancel"]}}},
  {"kind": "choose", "choices": ["the_d6"]}
]
```

If the step settles with no question open, the check fails and names what the game waits for: `prompt: no question is open, the game waits for a priority action or a pass`.

A game that ends stops the settling too. A hand-size discard is not answered: the step stops there.

`action` keeps its old rule, so old payloads (website links) answer as before: it stops at the first priority window with an empty stack and no attack or purchase, so an `end_turn` action stops once the end of the turn is declared.

### Refusals and the answer

A refused step changes nothing, and the next steps still run. A step that should be refused says so in its own `expect`: `{"allowed": false, "rule": "R-ABIL-07"}`.

The answer's `allowed`, `rule`, `reason` and `say` are those of the first refused step; `allowed` is true only when no step was refused. The top-level `expect.allowed` and `rule` are compared with that, so a file with a refused step also says `allowed: false` at the top. `events` hold every step's events. A mismatch names its step, or the end:

```
step 2 (choose): player 1 damage is 1, expected 2
end, first refusal at step 3 (activate razor), refused by R-ABIL-07: …: allowed is false, expected true
end: player 1 damage is 0, expected 1
```

An error (an unknown card, no open question) names the step too: `situation: step 1 (throw): throw: no roll waits to be thrown; the game waits for a priority action or a pass`.

## What the runner does

1. Builds the table; decks hold every card of the sets, shuffled with `seed`.
2. Tries the action, or each step in turn. A refusal answers `allowed: false` with a rule ID, an English `reason` for developers, and `say`: the player's words as a message key with values (ADR 021), which each client shows in its language.
3. If allowed, answers the choices and lets every player pass until the stack is empty; for steps, see [When a step settles](#when-a-step-settles).
4. With `expect` (top level or on a step), compares and lists any mismatch.

```go
ans, err := engine.RunSituation(data, b2.Set)
```

## As tests

A situation file with `expect` and `source` is a rules test (LR-07, ADR 034). Each package with tests keeps them in its own `testdata/situations/` tree, and its `go test` runs them: one subtest per file, named by its path in the tree without `.json`, e.g. `TestSituations/razor_blade/deals-1-damage-to-a-player`. A disputed case, once resolved, becomes such a file (R-03); it goes where its card or rule goes.

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

`enginetest.RunSituations(t, files, sets...)` runs a card tree; `enginetest.RunRuleSituations` the engine's topic tree (`situation_files_test.go`). For each file it checks the layout, then reads it strictly (an unknown field fails with its line, its path and the nearest field name, e.g. `line 7, setup: unknown field "rols"; did you mean "rolls"?`; a payload outside a tree is read as before), then checks that the file has `expect` and `source` and that each source exists, then runs it with `engine.RunSituation` and fails on rolls or picks left unused ([Leftover rolls and picks](#leftover-rolls-and-picks)). A mismatch names the file, the step and every failed check; with `steps`, a check of the top-level `expect` names the end:

```
razor_blade/hits-cain.json, step 1 (activate razor_blade): player 1 damage is 1, expected 2
razor_blade/hits-cain.json, end: player 1 damage is 0, expected 1
```

The strict reading covers steps too: `line 9, steps[1]: unknown field "choises"; did you mean "choices"?`.

The rule and question IDs come from `pkg/rules_engine/docs` (package `docs`), which embeds the rules text; only tests import it, so the website's wasm stays the same size.

A set package that has files adds one test (`cards/b2/situations_test.go` is the model):

```go
//go:embed testdata/situations
var situations embed.FS

// TestSituations runs this set's situation files (ADR 034).
func TestSituations(t *testing.T) {
	enginetest.RunSituations(t, situations, Set)
}
```

`go:embed` fails on a folder with no files, so a package adds this test with its first file, not before. The challenge set's runner (`cards/challenge/situations_test.go`) passes every set and the Go-only test sets `krampus_trial` (Mixed Bag: 2¢ and a loot, no real monster has both) and `trick_trial` (a bare mob: no rewards); a file there that needs them names them in `sets`. A cross-set file in `cards/testdata/situations/` passes every set (`cards.Sets()` and the challenge set) through `enginetest.RunInteractionSituations`; that tree holds only `interactions/`, so a card's own files stay in its set's tree.

## Dump a Go test

A Go card test can write its own draft file (24.7): `enginetest.Table` records its setup and every step it takes through the `Table` helpers, and with the `SITUATION_DUMP` env var set it writes them as JSON when the test ends. Run it for one test:

```sh
cd pkg/rules_engine
SITUATION_DUMP=.cache/situations go test ./cards/b2 -run 'TestDinga$' -count=1
```

A relative folder is from the repository root (`go test` runs in the package folder), and `.cache/` is not tracked. Each table writes `<folder>/<package folder>/<TestName>.json`, e.g. `.cache/situations/b2/TestDinga.json`; a subtest is a subfolder (`TestX/case.json`), and a second table in one test gets `-2`. Without the env var a table records nothing and writes nothing.

```json
{
  "title": "Dinga",
  "source": [],
  "sets": ["b2"],
  "setup": {"players": [{"character": "isaac"}, {"character": "cain"}], "active": 0, "monsters": ["dinga"], "rolls": [6, 6, 6, 3]},
  "steps": [{"player": 0, "kind": "attack", "choices": ["dinga"]}],
  "todo": ["a draft of TestDinga (24.7): …", "setup.rolls are every die the test threw, forced or not"]
}
```

(The file is indented one value per line; it is shown compact here.)

What it records:

* `title` from the test name, `sets` by code, and the setup the test built (`NewSetup`, `New` and the set helpers on them).
* State the test set by hand before its first step, where a setup field says it: cents, damage, items, hand, souls, deactivated cards, counters, `eternal`, `extra_loot_plays`, slot cards, `monster_damage`, slot `counters`, `loot_plays`, `attacks`, `purchases`, the attack duties (`must_attack`, `must_attacks`, `must_attack_deck`, 19.235; `must_attack_one_of`, 19.241), and cards put on top of a deck (`G.PutOnDeck`) as `decks`. A `must_attack` card the active player can't attack has no ref the setup takes, so it stays a `todo`.
* `rolls`: every die the test threw, in order, forced (`G.ForceRolls`, `Attack`'s rolls) or random, so the file throws the same dice.
* `game`: a challenge the test started (`G.StartChallenge`), with what it changed: presents, Happy Holidays!, the rule card's side, the boss's counters; a present or counter as the start leaves it is left out. The boss is no `bosses` entry. Fear and other player counters go to `counters` with `"card": "player"`.
* `random`: every random pick the test made, in order, by its option label, so the file picks the same (24.30). A pick with no label (a hidden card) leaves `random` out and adds a `todo` line instead.
* A step per outermost helper, with its choices as the plain labels the test gave: `Activate`, `Start` and `Refused` are `activate`, `Play` is `play`, `Attack` and `Do` with an attack are `attack`, `EndTurn` is `end_turn`, `Pass` is `pass`, `Choose` is `choose`, `RevealFromDeck` is `reveal`, and `ThrowPending` a `throw` per die. `Settle`'s choices go to the step before it. `UntilChoose` adds nothing: the step before it stops at the question when a `choose` follows.
* Intents the test applied itself (`tb.G.Apply`, `tb.G.Settle`) become steps too, read from the game's applied-intent log; passes and throws after a step are left out, since the runner settles each step.
* State set by hand before a direct `tb.G.Apply` that is the test's first step (24.11): the engine keeps the game as it was just before its first applied intent (its save origin, 19.177), and the recorder reads the setup from that state, so `tb.G.Object(…).Eternal = true` before the first `G.Apply` lands in `eternal`. Should that state be missing (a `Save` moved it), a `todo` line says so.

It never writes `expect`: an expected value comes from the rules, never from a run (CLAUDE.md 5.25, ADR 034). `source` is an empty list to fill.

`todo` lists what the file can't say yet, and the work left:

* State set by hand that no setup field says (`Players[1].SkipTurns`, a boost, a pushed stack item), found by rebuilding the table from the dumped setup and comparing it with the test's table before its first step.
* State changed between two steps outside a helper (`tb.G.Players[0].Cents = 3`, a `PutOnDeck` mid-test, `tb.G.Run()`).
* A choice made by a direct `G.Apply` (its label is unknown: `#<index>`), a response in the middle of a stack (owner answer 7: start from `setup.stack`), a helper that does not settle (`Start`, `Pass`), and a `Refused` step whose `expect` needs `allowed: false`.
* A response to an attack roll of a deck attack also names that roll as a `setup.stack` item with `from_deck` and `slot` (24.33), to copy. The dump writes steps, not `setup.stack`, so the steps and rolls before it are left out by hand.

`todo` is no situation field, so the runner rejects the file until it is gone. To turn a draft into a test file:

1. Check each `todo` line; a test with state no field says stays Go (ADR 034).
2. Write `title` as the question, and fill `source`.
3. Write `expect` from the Go test's checks and their source comments or card text; more checks from the card text are fine, never fewer (owner answer 8).
4. Remove `todo` and save it as `cards/<set>/testdata/situations/<card>/<case>.json`.
5. Run the set's `go test`, and delete the Go test it replaces in the same change (owner answer 1).

`cards/b2/testdata/situations/dinga/doubles-rewards-on-a-6.json` is the draft above, finished: `TestDinga` is gone.
