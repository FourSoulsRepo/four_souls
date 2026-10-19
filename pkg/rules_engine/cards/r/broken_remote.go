package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Broken Remote (Trinket Card)
//
//	Each time you activate an item, you may put a room or monster not being attacked into discard.
//	-Trinket- This loot becomes an item under your control when it resolves.
var brokenRemote = engine.CardDef{
	Ref:     "broken_remote",
	Kind:    engine.LootCard,
	Copies:  1,
	Trinket: true,
	Abilities: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "Each time you activate an item, you may put a room or monster not being attacked into discard.",
		Trigger: engine.Trigger{On: engine.EvActivated, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
			return e.Player == g.Object(self).Controller && e.Object != 0 && g.Object(e.Object).Role == engine.RoleItem
		}},
		Effects: []engine.Effect{roomOrMonsterToDiscard(true)},
	}},
}
