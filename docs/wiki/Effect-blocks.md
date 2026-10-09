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

## Targets

Targets are always asked, with a "cancel" option before anything is paid. A target that is gone when the ability resolves makes it fizzle.

| Block | Choose … |
|---|---|
| `Choose(TargetPlayer)` | a living player |
| `Choose(TargetMonster)` | a monster in play |
| `Choose(TargetMonsterOrPlayer)` | either |
| `Choose(TargetItem)` | an item a player controls |
| `Choose(TargetDiceRoll)` | a dice roll on the stack |

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
| `EffectFunc(func(c *Ctx){…})` | Anything the blocks do not cover |

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

## Static abilities

| Block | Example card text |
|---|---|
| `YouHave(StatPlayerATK, 1)` | You have +1 ATK. |
| `YouHave(StatShopPrice, -5)` | Shop items you purchase cost 5¢ less. |
| `MonstersHave(StatMonsterDC, 1)` | Monsters have +1 DC. |

Stats: `StatPlayerATK`, `StatPlayerHP`, `StatMonsterDC`, `StatMonsterATK`, `StatMonsterHP`, `StatShopPrice`.

## Replacement effects

A `Replacement` has a `When` (does it apply to this pending action?) and a `Do` (what happens instead). Returning no actions prevents the event. Each replacement applies once per event; if several apply, the affected player chooses the order.
