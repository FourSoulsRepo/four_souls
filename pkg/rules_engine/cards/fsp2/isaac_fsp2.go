package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// The card's own lines; their words are in en.json (ADR 021).
var (
	isaacFsp2ThisDeals1 = engine.NewKey("card.isaac_fsp2.this_deals_1_damage_to")
)

// Isaac! (Epic Boss Card)
//
//	Each time this takes damage, the active player chooses a living player. This deals 1 damage to that player.
var isaacFsp2 = engine.CardDef{
	Ref:       "isaac_fsp2",
	Kind:      engine.MonsterCard,
	Copies:    1,
	HP:        7,
	DC:        3,
	ATK:       1,
	Soul:      2,
	Rewards:   []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	ThreePlus: true,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time this takes damage, the active player chooses a living player. This deals 1 damage to that player.",
		Trigger: whenThisTakesDamage(false),
		Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
			if ps := livingPlayers(c.G, engine.NoPlayer); a[0] >= 0 && a[0] < len(ps) {
				c.G.DealDamageTo(engine.Target{Player: ps[a[0]], IsPlayer: true}, 1, engine.NoPlayer, c.Source)
			}
		}, pickAPlayer(isaacFsp2ThisDeals1))},
	}},
}
