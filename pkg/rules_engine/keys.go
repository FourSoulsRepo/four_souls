package rulesengine

// The engine's own lines (ADR 021, item 9). Their words and notes are in
// frontend/src/locales/en.json; shared option keys are in msg.go.

// Why a player can't do something now (Block, RuleError.Say).
var (
	refuseMustAttack      = NewKey("refuse.must_attack")
	refuseItemsForbidden  = NewKey("refuse.items_forbidden")
	refuseNoAttacks       = NewKey("refuse.no_attacks")
	refuseAttackCost      = NewKey("refuse.attack_cost") // n: the ¢ it costs
	refuseNoPurchases     = NewKey("refuse.no_purchases")
	refuseOneLoot         = NewKey("refuse.one_loot")
	refuseLootForbidden   = NewKey("refuse.loot_forbidden")
	refuseNothingToCopy   = NewKey("refuse.nothing_to_copy")
	refuseOnlyYourTurn    = NewKey("refuse.only_your_turn")
	refuseNextTurn        = NewKey("refuse.next_turn")
	refuseOneAbility      = NewKey("refuse.one_ability")
	refuseNoTarget        = NewKey("refuse.no_target")
	refuseNotUsable       = NewKey("refuse.not_usable")
	refuseAnswerPrompt    = NewKey("refuse.answer_prompt")
	refuseWaitPriority    = NewKey("refuse.wait_priority")
	refuseNotYourTurn     = NewKey("refuse.not_your_turn")
	refuseActionPhase     = NewKey("refuse.action_phase")
	refuseTurnEnding      = NewKey("refuse.turn_ending")
	refuseFinishAttack    = NewKey("refuse.finish_attack")
	refuseFinishPurchase  = NewKey("refuse.finish_purchase")
	refuseWaitStack       = NewKey("refuse.wait_stack")
	refuseLootYourTurn    = NewKey("refuse.loot_your_turn")
	refuseNoLootPlays     = NewKey("refuse.no_loot_plays")
	refuseCostUsed        = NewKey("refuse.cost.used")
	refuseCostCents       = NewKey("refuse.cost.cents")
	refuseCostHP          = NewKey("refuse.cost.hp")
	refuseCostEternal     = NewKey("refuse.cost.eternal")
	refuseCostCounters    = NewKey("refuse.cost.counters")
	refuseCostDiscard     = NewKey("refuse.cost.discard")
	refusePickGlowing     = NewKey("refuse.pick.glowing")
	refusePickGlowingOnly = NewKey("refuse.pick.glowing.only")
)

// What a wrong click shows: what the target may be (Aim.Reason). The
// .only line adds that only the glowing ones fit.
var refusePick = map[TargetKind][2]Key{
	TargetPlayer:          {NewKey("refuse.pick.player"), NewKey("refuse.pick.player.only")},
	TargetMonster:         {NewKey("refuse.pick.monster"), NewKey("refuse.pick.monster.only")},
	TargetMonsterOrPlayer: {NewKey("refuse.pick.monster_or_player"), NewKey("refuse.pick.monster_or_player.only")},
	TargetItem:            {NewKey("refuse.pick.item"), NewKey("refuse.pick.item.only")},
	TargetDiceRoll:        {NewKey("refuse.pick.dice_roll"), NewKey("refuse.pick.dice_roll.only")},
	TargetOtherPlayer:     {NewKey("refuse.pick.other_player"), NewKey("refuse.pick.other_player.only")},
	TargetStackAbility:    {NewKey("refuse.pick.stack_ability"), NewKey("refuse.pick.stack_ability.only")},
	TargetCurse:           {NewKey("refuse.pick.curse"), NewKey("refuse.pick.curse.only")},
	TargetYourItem:        {NewKey("refuse.pick.your_item"), NewKey("refuse.pick.your_item.only")},
	TargetYourHandCard:    {NewKey("refuse.pick.your_hand_card"), NewKey("refuse.pick.your_hand_card.only")},
	TargetItemOrSoul:      {NewKey("refuse.pick.item_or_soul"), NewKey("refuse.pick.item_or_soul.only")},
}

// The engine's own questions (Prompt.Text).
var (
	promptReplacement      = NewKey("prompt.replacement")
	promptReplacementCents = NewKey("prompt.replacement_cents") // n: the ¢ gained
	promptWhichDie         = NewKey("prompt.which_die")
	promptRedirect         = NewKey("prompt.redirect")
	promptPassCurse        = NewKey("prompt.pass_curse")
	promptCopyAbility      = NewKey("prompt.copy_ability")
	promptDiscardRoom      = NewKey("prompt.discard_room")
	promptDamageOrder      = NewKey("prompt.damage_order")
	askRestInSlotOrder     = NewKey("ask.rest_in_slot_order")
	askRollOf              = NewKey("ask.roll_of") // n: the roll
)

// What a stack item is (StackItem.Label), also on its events and on the
// throw prompt of a roll.
var (
	eventRoll            = NewKey("event.roll")
	eventAttackRoll      = NewKey("event.attack_roll")
	eventDamage          = NewKey("event.damage")
	eventCombatDamage    = NewKey("event.combat_damage")
	eventDeath           = NewKey("event.death")
	eventSoul            = NewKey("event.soul")
	eventRewards         = NewKey("event.rewards")
	eventRewardRoll      = NewKey("event.reward_roll")
	eventDeathPenalty    = NewKey("event.death_penalty")
	eventRollResult      = NewKey("event.roll_result")
	eventCopyOf          = NewKey("event.copy_of") // card: the copied card
	eventFirstPlayerRoll = NewKey("event.first_player_roll")
	// eventWords is an ability's own words, from the card: {text}.
	eventWords = NewKey("event.words")
)

// ratioKeys name the deck ratio groups with a shortfall: {have} of
// {need} (GS-07); see RatioShort.Msg.
var ratioKeys = [ratioGroupCount]Key{
	NoRatioGroup:              NewKey("refuse.ratio.no_group"),
	GroupPassiveItems:         NewKey("refuse.ratio.passive_items"),
	GroupActiveItems:          NewKey("refuse.ratio.active_items"),
	GroupPaidItems:            NewKey("refuse.ratio.paid_items"),
	GroupOneUseItems:          NewKey("refuse.ratio.one_use_items"),
	GroupSoulItems:            NewKey("refuse.ratio.soul_items"),
	GroupPennies:              NewKey("refuse.ratio.pennies"),
	GroupTwoCents:             NewKey("refuse.ratio.two_cents"),
	GroupThreeCents:           NewKey("refuse.ratio.three_cents"),
	GroupFourCents:            NewKey("refuse.ratio.four_cents"),
	GroupNickels:              NewKey("refuse.ratio.nickels"),
	GroupButterBeans:          NewKey("refuse.ratio.butter_beans"),
	GroupBombs:                NewKey("refuse.ratio.bombs"),
	GroupBatteries:            NewKey("refuse.ratio.batteries"),
	GroupDiceShardsSoulHearts: NewKey("refuse.ratio.dice_shards_soul_hearts"),
	GroupPillsRunes:           NewKey("refuse.ratio.pills_runes"),
	GroupWildcards:            NewKey("refuse.ratio.wildcards"),
	GroupTrinkets:             NewKey("refuse.ratio.trinkets"),
	GroupLostSouls:            NewKey("refuse.ratio.lost_souls"),
	GroupBasicMonsters:        NewKey("refuse.ratio.basic_monsters"),
	GroupCursedMonsters:       NewKey("refuse.ratio.cursed_monsters"),
	GroupHolyCharmedMonsters:  NewKey("refuse.ratio.holy_charmed_monsters"),
	GroupBosses:               NewKey("refuse.ratio.bosses"),
	GroupEpicBosses:           NewKey("refuse.ratio.epic_bosses"),
	GroupGoodEvents:           NewKey("refuse.ratio.good_events"),
	GroupBadEvents:            NewKey("refuse.ratio.bad_events"),
	GroupCurses:               NewKey("refuse.ratio.curses"),
}

// RefuseDeckRatios heads the reason the deck ratios can't be used:
// {n} players, then the groups (GS-07).
var RefuseDeckRatios = NewKey("refuse.deck_ratios")
