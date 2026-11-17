package join

import (
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/lan"
	"github.com/FourSoulsRepo/four_souls/internal/protocol"
)

// findWait is how long the join screen listens for answers.
const findWait = time.Second

// FindGames lists the games hosted on the local network that wait for
// players (N-13, ADR 009).
func (j *Join) FindGames() ([]lan.Found, error) {
	found, err := lan.Find(j.api.Context(), protocol.DefaultPort, findWait)
	if found == nil {
		found = []lan.Found{} // an empty list, not null, for the frontend
	}
	return found, lanError(err)
}
