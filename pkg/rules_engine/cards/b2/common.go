package b2

import "github.com/FourSoulsRepo/rules_engine/cards/cardkit"

// Pieces many Base Game cards share; they live in cardkit for every set.
var (
	allItems                = cardkit.AllItems
	attackRollOnThis        = cardkit.AttackRollOnThis
	changeRollBy            = cardkit.ChangeRollBy
	damageAMonster          = cardkit.DamageAMonster
	damageAPlayer           = cardkit.DamageAPlayer
	damageAnything          = cardkit.DamageAnything
	damageEqualToRoll       = cardkit.DamageEqualToRoll
	discardOne              = cardkit.DiscardOne
	eachMonsterGains        = cardkit.EachMonsterGains
	endOfTurnIf             = cardkit.EndOfTurnIf
	endYourTurn             = cardkit.EndYourTurn
	extraLoot               = cardkit.ExtraLoot
	firstAttackRollATK      = cardkit.FirstAttackRollATK
	forRoller               = cardkit.ForRoller
	forcedDiscard           = cardkit.ForcedDiscard
	forcedDiscardQuestions  = cardkit.ForcedDiscardQuestions
	giveThisAway            = cardkit.GiveThisAway
	guppyItems              = cardkit.GuppyItems
	killAMonster            = cardkit.KillAMonster
	killAPlayer             = cardkit.KillAPlayer
	livingPlayers           = cardkit.LivingPlayers
	lookMayBottom           = cardkit.LookMayBottom
	may                     = cardkit.May
	monstersInPlay          = cardkit.MonstersInPlay
	onYourAttackRollOf      = cardkit.OnYourAttackRollOf
	oneOnTop                = cardkit.OneOnTop
	otherPlayers            = cardkit.OtherPlayers
	pickAPlayer             = cardkit.PickAPlayer
	playerLabel             = cardkit.PlayerLabel
	playerRollOf            = cardkit.PlayerRollOf
	putBackInOrder          = cardkit.PutBackInOrder
	rechargeAnItem          = cardkit.RechargeAnItem
	rechargeAtEndOfTurn     = cardkit.RechargeAtEndOfTurn
	rollerDiscards          = cardkit.RollerDiscards
	searchGuppy             = cardkit.SearchGuppy
	whenThisDealsDamage     = cardkit.WhenThisDealsDamage
	whenThisTakesDamage     = cardkit.WhenThisTakesDamage
	whenYouDealCombatDamage = cardkit.WhenYouDealCombatDamage
)
