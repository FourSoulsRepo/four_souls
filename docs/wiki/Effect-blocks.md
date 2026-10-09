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
| `Trinket` | Loot that becomes an item when it resolves (R-ABIL-19) |
| `Curse` | An event given to a player when it enters play (R-ABIL-20) |
| `Guppy` | The Guppy tag (R-ABIL-22) |
| `Unattackable` | "This can't be attacked." |
| `CombatMod` | Changes combat damage by attack roll: "This takes no combat damage on attack rolls of 6." |
| `BonusSoul` | A bonus soul's condition: "The first player to have 25¢ or more gains this soul." |
| `DamageMod` | Changes damage about to be marked: "Damage you would take is reduced to 1." |
| `EntersDeactivated`, `EntersWithCounters` | "This enters play deactivated." / "starts with 9 counters" |
| `CopiesTapAbilities` | Placebo: may use any ↷ ability of another non-eternal item (`AbilitiesOf`) |
| `TakesPenalties` | Shadow: you choose the item and gain the loot and ¢ of others' death penalties |
| `PeeksTreasure` | "You may look at the top card of the treasure deck at any time on your turn" (in your view) |
| `SoulWhenDestroyed` | "If this would be destroyed, it becomes a soul instead." |
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
| `DiscardChosen(t)` | Discard a loot card: … (with `Choose(TargetYourHandCard)`) |
| `GiveChosen(item, player)` | Give an item you control to another player: … |
| `DestroyChosen(t…)` | Destroy 2 items you control: … (targets filtered with `NotChosen`) |

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
| `Choose(TargetStackAbility)` | an item's ↷ or $ ability, or a loot being played |
| `Choose(TargetCurse)` | a curse a player has |
| `Choose(TargetYourItem)` | an item you control |
| `Choose(TargetItemOrSoul)` | an item or soul a player controls |
| `ChooseWhere(kind, filter)` | e.g. the player with the most souls; `NotThis` is "another item" |

Effects refer to targets by number: `DealDamage(1, 0)` hits the first target. `This` is the ability's own object (a monster shielding itself). `You` means the controller, for text without a target: `DealDamage(1, You)` is "Take 1 damage." A roll ability keeps its targets for the result: "Choose a player, then roll- deal damage equal to the result".

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
| `PreventDamage(n, t)` | Prevent the next 1 damage they would take this turn. |
| `GainHPThisTurn(n, t)` | They gain +2 HP till end of turn. |
| `GainRollBonusThisTurn(n, t)` | They gain +1 to dice rolls till end of turn. |
| `RechargeTarget(t)` | Recharge an item. |
| `RechargeItemsOf(t)` | Choose a player. Recharge each item they control. |
| `Kill(t)` | Kill a player. |
| `CapNextDamage(n, t)` | The next instance of damage they take this turn is reduced to 1. |
| `EachOtherPlayer(effects…)` | Deal 1 damage to each other player. |
| `DestroyThis(then…)` | ↷: Destroy this. If you do, … |
| `RerollTarget(t)` | Reroll an item. (R-MECH-48) |
| `CancelTarget(t)` | Cancel the ↷ or $ ability of an item or a loot being played. |
| `DestroyTarget(t)` | Destroy a curse. |
| `EachPlayer(effects…)` | Each player gains 1¢. / Each player takes 3 damage. |
| `EachMonsterTakesDamage(n)` | Each monster takes 1 damage. (slot order) |
| `AddAttacks(n, t)` | They may attack an additional time this turn. |
| `PreventYourDeath()` | Prevent death. (heals to 1 HP) |
| `ThisToLootBottomExtraTurn()` | Put this on the bottom of the loot deck … take an extra turn (The Sun) |
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

A question with `Random: true` is answered by the game at random ("choose a player at random"). A triggered ability's `c.EventPlayer` and `c.EventAmount` tell who rolled or how much damage was taken. A question's `Player` func picks who answers, e.g. the chosen player (Judgement). Ready questions: `DeckQuestion(text)` (read with `c.Deck(i)`) and `HandQuestion(text)` (read with `c.HandCard(i)`). Game methods for `do` and `EffectFunc`: `LookAt`, `DeckTop`, `SetDeckTop`, `DeckToBottom`, `MillTop`, `DiscardTopToDeck`, `HandToDeckTop`, `GiveHandCard`, `DiscardFromHand`, `DestroyObject`, `GainControl`, `ShopItems`, `DiscardMonster`, `RefillSlots`, `SetRoll`, `CancelStackItem`, `EndTurnNow`. `c.Do(effects…)` runs ordinary blocks.

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
| `WhenAPlayerDies()` | when any player dies, before the penalty (R-DEATH-13) |
| `AfterYourDeathPenalty()` | each time you die, after paying penalties |
| `WhenYouWouldTakeDamage()` | when damage aimed at you goes on the stack (resolves first) |
| `WhenYouWouldDie()` | when your death goes on the stack (resolves first) |
| `WhenThisEntersPlay()` | when this object enters play |
| `WhenThisDies()` | when this monster dies, before its rewards (R-DEATH-05); `c.EventObject` is the dead card |
| `AfterThisRewards()` | after this dead monster's rewards (R-DEATH-07) |
| `WhenARollWouldBe(n)` | each time a player would roll n (R-DICE-05); `c.EventStack` is the roll |

## Static abilities

| Block | Example card text |
|---|---|
| `YouHave(StatPlayerATK, 1)` | You have +1 ATK. |
| `YouHave(StatShopPrice, -5)` | Shop items you purchase cost 5¢ less. |
| `MonstersHave(StatMonsterDC, 1)` | Monsters have +1 DC. |

Stats: `StatPlayerATK`, `StatPlayerHP`, `StatMonsterDC`, `StatMonsterATK`, `StatMonsterHP`, `StatShopPrice`, `StatRoll`, `StatAttackRoll`, `StatLootPlays`, `StatAttacks`, `StatPurchases`, `StatLootStep`.

`YouHave(StatLootPlays, 1)` is "You may play an additional loot card on your turn"; `YouHave(StatLootStep, 1)` is "Loot +1 during your loot step"; `YouHave(StatAttackRoll, 1)` is "+1 to attack rolls"; `YouHave(StatLockOthers, 1)` is Trinity Shield; a `StatLootDouble` boost doubles a player's loot. A roll that would resolve lets "would roll" triggers act first and tries again only if its value changed. Shop items' abilities do not work in the shop (R-CARD-05).

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

## Events and the monster deck

An event's abilities are `Triggered` with `WhenThisEntersPlay()`; "you" is the active player (R-CARD-10). The event stays in its slot until they are done, then goes to the monster discard (R-CARD-15). A `Curse` goes to a player the active player picks, and to discard when that player dies. A reward with `Roll: true` gives as many as a roll's result ("Roll- gain X¢"). Helpers: `ExpandShop`, `ExpandMonsters`, `ForceAttacks(n, deck)`, `AddDeckAttack`, `TakeFromDeck`, `PutIntoDeck`, `PlaceFromDeck`, `KillObject`, `HealPlayer`, `HealObject`.

A monster's death runs in steps: "when this dies" triggers, rewards, "after rewards" triggers, then the soul or discard and the refill (R-DEATH-04 to R-DEATH-09). A step waits under its triggers on the stack. An object remembers who killed it (`Killer()`) and on which attack roll (`KilledOn`).

## Copies

An object's `CopyOf` makes it act as another card: abilities, statics and stats come from that card (`CardOf`). Where it goes when it leaves play still follows its own card. `CopyThisTurn` ends the copy at the end of the turn (Diplopia). A player's `CopyNextLoot` puts a copy of their next non-trinket loot on the stack above it (Blank Card).

## Known simplifications

* "Each monster takes damage" goes on the stack in slot order; the rules let the active player pick (R-MECH-29).
* "Put the rest on the bottom" keeps the cards' order; the rules let the player pick it.
