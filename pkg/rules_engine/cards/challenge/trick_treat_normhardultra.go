package challenge

import engine "github.com/FourSoulsRepo/rules_engine"

// Trick (Challenge Card)
//
//	As each turn starts, you may flip this.
//	Each time you kill a monster, remove a trick counter from the final boss.
//	At the start of each attack against a monster, the active player rolls. If you are playing on hard, roll 2x instead. If you are playing on ultra hard, roll 3x instead-
//	1-2: Damage you would take is increased by 1 till end of turn.
//	3-4: Reveal cards from the top of the monster deck until you reveal a non-event monster card. The monster being attacked gains its abilities till end of turn. Put any revealed cards on the bottom of the deck.
//	5-6: The rewards of that monster become 1¢ till end of turn.
var trickTreatNormhardultra = engine.CardDef{
	Ref:    "trick_treat_normhardultra",
	Kind:   engine.ChallengeCard,
	Copies: 1,
	Back:   "treat_normhardultra",
	Abilities: []engine.Ability{
		trickOrTreating(flipToTreatAsk),
		nobodysFool,
		ghostlyTrickery,
		trickGoal,
	},
}

// treatNormhardultra is the back of trick_treat_normhardultra: Treat.
//
//	As each turn starts, you may flip this.
//	Each time you gain rewards, remove a fear counter from a player.
//	At the start of each attack against a monster, the active player rolls-
//	1-2: Add 3¢ to the rewards of that monster till end of turn.
//	3-4: Add loot 1 to the rewards of that monster till end of turn.
//	5-6: Add +1 treasure to the rewards of that monster till end of turn.
var treatNormhardultra = engine.CardDef{
	Ref:     "treat_normhardultra",
	Kind:    engine.ChallengeCard,
	Outside: true,
	Abilities: []engine.Ability{
		trickOrTreating(flipToTrickAsk),
		youDontScareMe,
		hallowsEveBounty,
		trickGoal,
	},
}
