# Effect blocks

Cards are Go values built from small blocks (ADR 005). This page lists every block with a one-line example. The full guide to adding a card comes with roadmap step 5.10.

```go
{Ref: "razor", Kind: TreasureCard, Abilities: []Ability{
    {Kind: Activated, Text: "Pay 5¢: Deal 1 damage to a monster or player.",
        Costs:   []Cost{PayCents(5)},
        Targets: []TargetSpec{Choose(TargetMonsterOrPlayer)},
        Effects: []Effect{DealDamage(1, 0)}},
}}
```

## Card definition fields

| Field | Meaning |
|---|---|
| `Ref` | Card ID, as in `card_db` (`the_d6`, `the_d6@2`) |
| `Kind` | `CharacterCard`, `TreasureCard`, `LootCard`, `MonsterCard`, `EventCard`, `BonusSoulCard`, `RoomCard` |
| `Copies` | Copies in the set |
| `HP`, `ATK`, `DC`, `Soul` | Stat box and soul value |
| `StartingItem` | A character's starting item |
| `Eternal` | Keyword: cannot be destroyed or discarded |
| `Outside` | Starts outside the game (starting items) |
| `Tap` | Has a ↷ ability (the death penalty deactivates it) |
| `GoesFirst` | A character whose player goes first (Cain) |
| `StartingChoice` | Look at this many top treasures at the start, pick an eternal starting item (Eden) |
| `Rewards` | Reward box, e.g. `[]Reward{{Kind: RewardCents, Amount: 3}}` |
| `Abilities` | Activated, loot and triggered abilities |
| `Statics` | Static abilities that change numbers |
| `Replacements` | Replacement effects ("instead", "prevent") |

## Ability kinds

| Kind | When it happens | Example |
|---|---|---|
| `Activated` | The controller uses it with priority, paying costs | ↷: Gain 3¢. |
| `LootAbility` | When the loot card resolves; targets are picked when it is played | Deal 1 damage to a monster or player. |
| `Triggered` | When its `Trigger` matches; goes on the stack at the next priority | At the start of your turn, gain 1¢. |

## Costs

| Block | Example |
|---|---|
| `Tap()` | ↷: deactivate the item (it must be charged) |
| `PayCents(n)` | Pay 5¢: … |
| `DestroySelf()` | Destroy this: … (not on eternal objects) |
| `RemoveCounters(n)` | Remove 2 counters from this: … |

## Targets

Targets are always asked, with a "cancel" option before anything is paid. A target that is gone when the ability resolves makes it fizzle.

| Block | Choose … |
|---|---|
| `Choose(TargetPlayer)` | a living player |
| `Choose(TargetMonster)` | a monster in play |
| `Choose(TargetMonsterOrPlayer)` | either |
| `Choose(TargetItem)` | an item a player controls |
| `Choose(TargetDiceRoll)` | a dice roll on the stack |
| `Choose(TargetOtherPlayer)` | a living player other than you |

Effects refer to targets by number: `DealDamage(1, 0)` hits the first target.

## Effects

| Block | Example card text |
|---|---|
| `Loot(n)` | Loot 2. |
| `GainCents(n)` | Gain 3¢. |
| `LoseCents(n)` | Lose 4¢. |
| `GainTreasure(n)` | Gain +1 treasure. |
| `DealDamage(n, t)` | Deal 1 damage to a monster or player. |
| `AddLootPlays(n)` | Play an additional loot card this turn. |
| `RerollRoll(t)` | Choose a dice roll. Its controller rerolls it. |
| `Roll(table)` | Roll- 1-2: Loot 1. 3-4: Gain 3¢. 5-6: Lose 4¢. |
| `RechargeSelf()` | Recharge this. |
| `AddCounters(n)` | Put a counter on this. |
| `ModifyRoll(n, t)` | Add +1 to a dice roll. (stays between 1 and 6) |
| `BecomeSoul()` | This becomes a soul and loses all abilities. |
| `StealCents(n, t)` | Steal 1¢ from another player. |
| `GainATKThisTurn(n, t)` | They gain +1 ATK till end of turn. |
| `PreventNextDamage(t)` | Prevent the next instance of damage they would take this turn. |
| `Ask(do, questions…)` | Anything decided on resolution: "a deck", "a card from your hand" |
| `EffectFunc(func(c *Ctx){…})` | Anything the blocks do not cover |

## Choose one

"Choose one-" options are `Modes`. The player picks the mode when activating, before targets (R-ABIL-04); only modes with valid targets are offered, plus "cancel".

```go
{Kind: Activated, Costs: []Cost{Tap()}, Modes: []Mode{
    {Text: "Add 1 to a roll.", Targets: []TargetSpec{Choose(TargetDiceRoll)}, Effects: []Effect{ModifyRoll(1, 0)}},
    {Text: "Subtract 1 from a roll.", Targets: []TargetSpec{Choose(TargetDiceRoll)}, Effects: []Effect{ModifyRoll(-1, 0)}},
}}
```

## Questions on resolution

What is not a target is decided when the ability resolves (R-ABIL-05). `Ask` asks the controller its questions in order, then calls `do` with the answers. A question whose `Options` are empty is skipped and its answer is `-1`. The answers are saved in the game state, so a saved game continues mid-question.

```go
Ask(func(c *Ctx, a []int) {
    if a[0] >= 0 {
        c.G.MillTop(c.Deck(a[0]))
    }
}, DeckQuestion("Put the top card of which deck into discard?"))
```

Ready questions: `DeckQuestion(text)` (read with `c.Deck(i)`) and `HandQuestion(text)` (read with `c.HandCard(i)`). Game methods for `do`: `LookAt`, `DeckTop`, `SetDeckTop`, `MillTop`, `DiscardTopToDeck`, `HandToDeckTop`, `GiveHandCard`, `DiscardFromHand`. `c.Do(effects…)` runs ordinary blocks.

A roll table is built with `Results`:

```go
Roll(RollTable{}.Results(1, 2, Loot(1)).Results(3, 4, GainCents(3)).Results(5, 6, LoseCents(4)))
```

## Triggers

| Block | Triggers … |
|---|---|
| `AtStartOfYourTurn()` | at the start of the controller's turn |
| `AtEndOfYourTurn()` | at the end of the controller's turn |
| `WhenAMonsterDies()` | when any monster dies |
| `OnRollOf(n)` | when any roll resolves as n (circled number) |
| `WhenThisIsDestroyed()` | when this object is destroyed |
| `WhenYouDie()` | when the controller dies, after the penalties |
| `WhenYouTakeDamage()` | each time the controller takes damage |

## Static abilities

| Block | Example card text |
|---|---|
| `YouHave(StatPlayerATK, 1)` | You have +1 ATK. |
| `YouHave(StatShopPrice, -5)` | Shop items you purchase cost 5¢ less. |
| `MonstersHave(StatMonsterDC, 1)` | Monsters have +1 DC. |

Stats: `StatPlayerATK`, `StatPlayerHP`, `StatMonsterDC`, `StatMonsterATK`, `StatMonsterHP`, `StatShopPrice`.

## Replacement effects

A `Replacement` has a `When` (does it apply to this pending action?) and a `Do` (what happens instead). Returning no actions prevents the event. Each replacement applies once per event; if several apply, the affected player chooses the order.

## Testing a card

Each card gets `<card>_test.go` next to it; the [status report](Rules-engine) counts a card as tested when that file exists. The `enginetest` package builds a table and plays it:

```go
tb := enginetest.New(t, b2.Set, enginetest.Seat("isaac", "the_d6"), enginetest.Seat("cain"))
tb.Activate(0, "the_d6", 0, "roll of 2") // answer prompts by label
tb.EndTurn()                             // run to the next player's action phase
```

`NewSetup` takes a full `SituationSetup` (hands, monsters, a roll on the stack). `ForceRolls` fixes dice results, `Attack` attacks with given rolls, `Start` and `Pass` leave something on the stack so another player can respond.
