package rulesengine

// Stat is a number static abilities can change.
type Stat int

// The stats.
const (
	StatPlayerATK          Stat = iota // a player's attack
	StatPlayerHP                       // a player's max HP
	StatMonsterDC                      // a monster's evasion
	StatMonsterATK                     // a monster's attack
	StatMonsterHP                      // a monster's max HP
	StatShopPrice                      // the price of a shop item for a player
	StatRoll                           // added to a player's dice rolls as they resolve (R-DICE-06)
	StatAttackRoll                     // added to a player's attack rolls
	StatLootPlays                      // more loot plays on your turn
	StatAttacks                        // more attacks on your turn
	StatPurchases                      // more purchases on your turn
	StatLootStep                       // more cards in your loot step
	StatLockOthers                     // others can't play loot or activate items on your turn
	StatLootDouble                     // a player loots double (Two of Clubs)
	StatDeathWard                      // the player's next death this turn is prevented (Holy Mantle)
	StatAttackedDC                     // the evasion of monsters the player attacks
	StatRechargeLimit                  // the player recharges at most 1 item in their recharge step (Curse of Fatigue)
	StatHandCap                        // the player's max hand size is at most Amount (Curse of Tiny Hands)
	StatNextAttackRoll                 // a boost: added to the next attack roll this turn (R-ATK-19)
	StatLockTurn                       // a boost: other players can't play loot or activate items this turn (Get Out of Jail Card)
	StatNoCharRecharge                 // the player's character doesn't recharge in the recharge step (Mama Haunt)
	StatRevealHand                     // the player's hand is revealed to everyone (Guppy's Eye; Applies picks whose)
	StatChooseAttack                   // the player may choose which monster another player attacks (Mom's Eye Shadow)
	StatUsed                           // a boost marking an object's "once each turn" ability as used
	StatDoubleRewards                  // rewards are doubled (A Lucky Penny!, Blessing of the Sack)
	StatEternal                        // a boost: the object is eternal till end of turn (Algiz)
	StatNoForcedLoss                   // abilities and the death penalty can't make the player discard or lose ¢ (Rib of Greed)
	StatRollToOne                      // a boost: a player's roll of 2 or 3 becomes a 1 as it resolves (Eternal D6, R-DICE-12)
	StatPaidDiscount                   // $ abilities of the player's items cost this much less (Trick Penny)
	StatNoSouls                        // the player can't gain souls (Curse of the Soulless, All Hallows' Eve)
	StatTwoDiceLeftChooses             // players roll 2 dice for each roll; the player to their left picks the result (Cursed Lil Haunt)
	StatPenaltyItemTo                  // a boost: the player's death-penalty item goes to player Amount-1 this turn instead (Lust for Blood)
	StatPlayerAttackDC                 // the active player may attack this player, who has this DC (R-ATK-20); the newest effect wins (R-MECH-06)
)

var statNames = []string{"ATK", "HP", "DC", "ATK", "HP", "price", "roll", "attack roll", "loot plays", "attacks", "purchases", "loot step", "lock", "double loot", "death ward", "DC of attacked monsters", "recharge limit", "max hand size", "next attack roll", "lock this turn", "no character recharge", "revealed hand", "choose attacks", "used", "double rewards", "eternal", "no forced loss", "2 or 3 becomes 1", "cheaper $ abilities", "no souls", "two dice, left chooses", "penalty item given", "DC when attacked"}

func (s Stat) String() string { return statNames[s] }

// Static is a static ability that changes a stat while its object is in
// play (R-ABIL-12, R-ABIL-28). Applies says whom it affects: for player
// stats the subject is the player, for monsters the monster object.
type Static struct {
	Stat    Stat
	Amount  int
	Applies func(g *Game, self ObjectID, player PlayerID, subject ObjectID) bool
	// Value, if set, is the amount instead of Amount: "You have +1 HP for
	// each counter on this".
	Value func(g *Game, self ObjectID) int
}

// amount is the static's amount for its object.
func (s Static) amount(g *Game, self ObjectID) int {
	if s.Value != nil {
		return s.Value(g, self)
	}
	return s.Amount
}

// YouHave changes a stat of the object's controller: YouHave(StatPlayerATK, 1).
func YouHave(s Stat, n int) Static {
	return Static{Stat: s, Amount: n, Applies: func(g *Game, self ObjectID, p PlayerID, _ ObjectID) bool {
		return g.Object(self).Controller == p
	}}
}

// MonstersHave changes a stat of every monster: MonstersHave(StatMonsterDC, 1).
func MonstersHave(s Stat, n int) Static {
	return Static{Stat: s, Amount: n, Applies: func(*Game, ObjectID, PlayerID, ObjectID) bool { return true }}
}

// bonus sums the statics in play that change stat for the subject.
func (g *Game) bonus(stat Stat, p PlayerID, subject ObjectID) int {
	total := 0
	for _, b := range g.Boosts {
		if b.Stat == stat && ((b.Object != 0 && b.Object == subject) || (b.Object == 0 && subject == 0 && b.Player == p)) {
			total += b.Amount
		}
	}
	for _, id := range g.inPlay() {
		for _, face := range g.facesOf(id) {
			d, ok := g.cards.find(face)
			if !ok {
				continue
			}
			for _, s := range d.Statics {
				if s.Stat == stat && s.Applies(g, id, p, subject) {
					total += s.amount(g, id)
				}
			}
		}
	}
	return total
}
