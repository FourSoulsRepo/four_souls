package challenge

import engine "github.com/FourSoulsRepo/rules_engine"

// The challenges whose cards are in Set, so a situation's setup.game
// starts one by its ID (24.34). Set is generated (set_gen.go), so they
// join it here; the challenges name Set, so they can't be in its
// literal (an initialization cycle).
func init() {
	Set.Challenges = []*engine.ChallengeDef{&Krampus, &TrickOrTreat}
}
