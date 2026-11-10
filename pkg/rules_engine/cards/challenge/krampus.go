package challenge

import (
	"errors"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// Krampus is How the Krampus Stole Christmas (ADR 028; R-KRAMPUS).
// Its own setup comes in roadmap 14.22; until then a game can't start.
var Krampus = engine.ChallengeDef{
	ID:  "krampus",
	Set: Set.Code,
	Levels: []engine.ChallengeLevel{
		{Difficulty: engine.DifficultyNormal, Boss: "pilfering_krampus_norm", Card: "how_the_krampus_stole_christmas_norm"},
		{Difficulty: engine.DifficultyHard, Boss: "pilfering_krampus_hard", Card: "how_the_krampus_stole_christmas_hard"},
		{Difficulty: engine.DifficultyUltra, Boss: "pilfering_krampus_ultra", Card: "how_the_krampus_stole_christmas_ultra"},
	},
	Setup: krampusSetup,
}

// errKrampusSetup: the boss and the presents are not set up yet (14.22).
var errKrampusSetup = errors.New("the Krampus setup is not implemented yet (roadmap 14.22)")

func krampusSetup(*engine.Game, string) error { return errKrampusSetup }
