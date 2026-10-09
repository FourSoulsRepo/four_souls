package rulesengine

// Stat is a number static abilities can change.
type Stat int

// The stats.
const (
	StatPlayerATK  Stat = iota // a player's attack
	StatPlayerHP               // a player's max HP
	StatMonsterDC              // a monster's evasion
	StatMonsterATK             // a monster's attack
	StatMonsterHP              // a monster's max HP
	StatShopPrice              // the price of a shop item for a player
)

// Static is a static ability that changes a stat while its object is in
// play (R-ABIL-12, R-ABIL-28). Applies says whom it affects: for player
// stats the subject is the player, for monsters the monster object.
type Static struct {
	Stat    Stat
	Amount  int
	Applies func(g *Game, self ObjectID, player PlayerID, subject ObjectID) bool
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
	for _, id := range g.inPlay() {
		o := g.Object(id)
		d, ok := g.cards.find(o.Card)
		if !ok {
			continue
		}
		for _, s := range d.Statics {
			if s.Stat == stat && s.Applies(g, o.ID, p, subject) {
				total += s.Amount
			}
		}
	}
	return total
}
