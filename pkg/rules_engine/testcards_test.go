package rulesengine

// testSet is a tiny fake set for engine tests (no real cards needed).
var testSet = CardSet{Name: "test", Cards: []CardDef{
	{Ref: "hero_a", Kind: CharacterCard, HP: 2, ATK: 1, StartingItem: "item_a"},
	{Ref: "hero_b", Kind: CharacterCard, HP: 2, ATK: 1, StartingItem: "item_b"},
	{Ref: "hero_c", Kind: CharacterCard, HP: 3, ATK: 1},
	{Ref: "hero_d", Kind: CharacterCard, HP: 2, ATK: 2},
	{Ref: "item_a", Kind: TreasureCard, Eternal: true, Outside: true},
	{Ref: "item_b", Kind: TreasureCard, Eternal: true, Outside: true},
	{Ref: "trinket", Kind: TreasureCard, Copies: 6},
	{Ref: "penny", Kind: LootCard, Copies: 30},
	{Ref: "gaper", Kind: MonsterCard, HP: 2, DC: 3, ATK: 1, Copies: 4},
	{Ref: "boss", Kind: MonsterCard, HP: 4, DC: 4, ATK: 1, Soul: 1, Copies: 2},
	{Ref: "ambush_event", Kind: EventCard, Copies: 3},
	{Ref: "soul_x", Kind: BonusSoulCard},
	{Ref: "soul_y", Kind: BonusSoulCard},
	{Ref: "soul_z", Kind: BonusSoulCard},
	{Ref: "soul_w", Kind: BonusSoulCard},
}}
