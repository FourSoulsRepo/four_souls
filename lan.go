package main

import (
	"context"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/lan"
	"github.com/FourSoulsRepo/four_souls/internal/protocol"
)

// findWait is how long the join screen listens for answers.
const findWait = time.Second

// FindGames lists the games hosted on the local network that wait for
// players (N-13, ADR 009).
func (a *App) FindGames() ([]lan.Found, error) {
	ctx := a.ctx
	if ctx == nil {
		ctx = context.Background()
	}
	found, err := lan.Find(ctx, protocol.DefaultPort, findWait)
	if found == nil {
		found = []lan.Found{} // an empty list, not null, for the frontend
	}
	return found, playerError(err)
}
