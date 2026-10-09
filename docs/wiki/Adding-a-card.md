# Adding a card

This guide takes you from nothing to a working, tested card. You need Go (see `go.mod` for the version) and a clone of this repository. Card images are not needed.

Cards are Go values. Each card is one short file built from the blocks listed on [Effect blocks](Effect-blocks); a card the blocks cannot express gets a small Go function.

## 1. Find or make the stub

Base Game cards live in `pkg/rules_engine/cards/b2/`, one file per card, named after the card ID: `the_d6.go`, `a_nickel.go`.

If the card exists in `card_db`, the generator has already written a stub:

```sh
go run ./cmd/cardgen -set b2
```

A stub has the printed text as a comment, the game data (kind, copies, HP, ATK, DC, soul, rewards, keywords) and a `TODO(card):` line:

```go
// Bomb! (Bomb Card)
//
//	Deal 1 damage to a monster or player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var bomb = engine.CardDef{
	Ref:    "bomb",
	Kind:   engine.LootCard,
	Copies: 4,
}
```

While the `TODO(card):` line is there the generator rewrites the file. Delete it when you start: from then on the file is yours. Keep the variable name; the set list (`set_gen.go`) refers to it.

## 2. Write the card

Add what the text does. Most cards are a few blocks.

A loot card: its loot ability happens when it resolves; targets are chosen when it is played.

```go
var bomb = engine.CardDef{
	Ref:    "bomb",
	Kind:   engine.LootCard,
	Copies: 4,
	Abilities: []engine.Ability{{
		Kind:    engine.LootAbility,
		Text:    "Deal 1 damage to a monster or player.",
		Targets: []engine.TargetSpec{engine.Choose(engine.TargetMonsterOrPlayer)},
		Effects: []engine.Effect{engine.DealDamage(1, 0)},
	}},
}
```

An item with a ↷ ability and a trigger:

```go
var theD6 = engine.CardDef{
	Ref: "the_d6", Kind: engine.TreasureCard, Copies: 1, Eternal: true, Outside: true, Tap: true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Choose a dice roll. Its controller rerolls it.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetDiceRoll)},
			Effects: []engine.Effect{engine.RerollRoll(0)},
		},
		{
			Kind:    engine.Triggered,
			Text:    "At the end of your turn, recharge this.",
			Trigger: engine.AtEndOfYourTurn(),
			Effects: []engine.Effect{engine.RechargeSelf()},
		},
	},
}
```

A monster: stats and rewards come from the stub; "when this dies" is a trigger. On a monster card, "you" is the active player.

```go
var boomFly = engine.CardDef{
	Ref: "boom_fly", Kind: engine.MonsterCard, Copies: 1, HP: 1, DC: 4, ATK: 1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 4}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "When this dies, it deals 1 damage to each player.",
		Trigger: engine.WhenThisDies(),
		Effects: []engine.Effect{engine.EachPlayer(engine.DealDamage(1, engine.You))},
	}},
}
```

A custom hook, when no block fits: `EffectFunc` gets the context (`c.G` is the game, `c.Controller` the player, `c.Targets` the chosen targets) and may call the game's exported methods.

```go
Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
	c.G.Players[c.Controller].Cents = 0
})},
```

Things decided while the card resolves ("choose a deck", "a card from your hand") are questions: see "Questions on resolution" on [Effect blocks](Effect-blocks). Targets and "choose one-" are chosen when the card is played, never on resolution (R-ABIL-04).

Rules follow our rules text in `pkg/rules_engine/docs/rules/`; each rule has an ID such as `R-DEATH-14`. When the card text and a rule seem to disagree, add the case to `open-questions.md` instead of guessing.

## 3. Test it

Every card gets its own test file next to it, `<card>_test.go`. The `enginetest` package builds a table, plays, and answers prompts by their labels:

```go
package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBomb(t *testing.T) {
	tb := lootTable(t, "bomb") // Isaac holds the bomb; Cain sits next to him
	fly, _ := tb.G.Monsters[0].TopOf()
	tb.Play(0, "bomb", "fly")
	if tb.G.Object(fly).Zone.Kind == engine.ZoneInPlay {
		t.Error("the fly survived 1 damage")
	}
}
```

Useful tools, all in `enginetest` and `cards/b2/common_test.go`:

* `enginetest.NewSetup(t, engine.SituationSetup{…}, Set)`: any table — items, hands, cents, souls, monsters, a dice roll on the stack.
* `Activate`, `Play`, `Attack`, `EndTurn`: act and settle; the strings answer the prompts in order.
* `Start`, `Choose`, `Pass`: step by step, so another player can respond.
* `tb.G.ForceRolls(6, 1, …)`: fix the next dice results.
* `RevealFromDeck(card)`: put a monster deck card into play through an attack, for events.

Run the tests and update the status report:

```sh
cd pkg/rules_engine && go test ./cards/...
cd ../.. && go run ./cmd/cardgen -report
```

A test fails while the report is out of date, so commit the regenerated `card-status.md` with the card. Then let the fuzzer try your card with everything else (see [Fuzzing](Fuzzing)):

```sh
cd pkg/rules_engine
go test -run '^$' -fuzz=FuzzSets -fuzztime=2m -fuzzminimizetime=5s ./cards/
```

## 4. A fan-made set

A set is a Go package under `pkg/rules_engine/cards/` with one variable, `Set`.

1. Make the folder: `pkg/rules_engine/cards/myset/`.
2. Add `doc.go` with the package comment, and one file per card as above.
3. List the cards in `set.go`:

   ```go
   package myset

   import engine "github.com/FourSoulsRepo/rules_engine"

   // Set is every card of My Set.
   var Set = engine.CardSet{Code: "myset", Name: "My Set", Cards: []engine.CardDef{goldenFly}}
   ```

4. Register it in `pkg/rules_engine/cards/cards.go`, in `Sets()`.
5. Card IDs must be unique across all sets; pick names that no official card uses.

A fan set is played together with the Base Game, so it may be small: one card is a set. Cards of a fan set may refer to Base Game cards (a starting item, a card to search for) when the game is played with both sets.

If the set also has display data (names, text, images), add `pkg/card_db/data/myset.json` in the card data format (ADR 004); then `go run ./cmd/cardgen -set myset` writes the stubs for you and the status report counts the set. Without card data the set plays, but the client shows its cards as their IDs.

## Checklist

* The file has no `TODO(card):` line.
* `<card>_test.go` exists and passes.
* The status report is regenerated.
* `make lint` passes (see [Linters](Linters)).
