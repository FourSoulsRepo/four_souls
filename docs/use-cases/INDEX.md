# Use cases

| ID | Name | Module | Status | File |
|----|------|--------|--------|------|
| UC-01 | Start the app with the fan-game notice | `internal/legal`, `frontend` | Implemented | [app/uc-01-start-app-with-notice.md](app/uc-01-start-app-with-notice.md) |
| UC-02 | Start the dedicated server | `cmd/server` | In progress | [server/uc-02-start-dedicated-server.md](server/uc-02-start-dedicated-server.md) |
| UC-03 | Host a game | `internal/server`, `host.go` | In progress | [server/uc-03-host-a-game.md](server/uc-03-host-a-game.md) |
| UC-04 | Join a game | `internal/server` | In progress | [server/uc-04-join-a-game.md](server/uc-04-join-a-game.md) |
| UC-05 | Detailed match setup | `internal/server` | In progress | [server/uc-05-detailed-match-setup.md](server/uc-05-detailed-match-setup.md) |
| UC-06 | Respond without stalls | `internal/server` | In progress | [server/uc-06-respond-without-stalls.md](server/uc-06-respond-without-stalls.md) |
| UC-07 | Reconnect after a lost connection | `internal/server` | In progress | [server/uc-07-reconnect.md](server/uc-07-reconnect.md) |
| UC-08 | Get a match record | `internal/server`, `pkg/record` | In progress | [server/uc-08-get-a-match-record.md](server/uc-08-get-a-match-record.md) |
| UC-09 | Save and continue a game | `internal/server`, `pkg/record` | In progress | [server/uc-09-save-and-continue.md](server/uc-09-save-and-continue.md) |
| UC-10 | Change settings | `internal/settings`, `frontend` | In progress | [app/uc-10-change-settings.md](app/uc-10-change-settings.md) |
| UC-11 | Look at a card | `frontend`, `cards.go` | Implemented | [app/uc-11-look-at-a-card.md](app/uc-11-look-at-a-card.md) |
| UC-12 | See the table | `frontend` | In progress | [app/uc-12-see-the-table.md](app/uc-12-see-the-table.md) |
| UC-13 | Take actions | `frontend` | Implemented | [app/uc-13-take-actions.md](app/uc-13-take-actions.md) |
| UC-14 | Follow the stack | `frontend` | Implemented | [app/uc-14-follow-the-stack.md](app/uc-14-follow-the-stack.md) |
| UC-15 | Read the history | `frontend` | Implemented | [app/uc-15-read-the-history.md](app/uc-15-read-the-history.md) |
| UC-16 | Run several copies on one computer | `internal/instance`, `main.go` | Implemented | [app/uc-16-run-several-copies.md](app/uc-16-run-several-copies.md) |
| UC-17 | Join a game found on the LAN | `internal/lan`, `frontend` | Implemented | [server/uc-17-join-a-lan-game.md](server/uc-17-join-a-lan-game.md) |
| UC-18 | See a hidden card | `frontend`, `rules_engine` | Implemented | [app/uc-18-see-a-hidden-card.md](app/uc-18-see-a-hidden-card.md) |
| UC-19 | Create a game on a dedicated server | `internal/server`, `frontend` | Implemented | [server/uc-19-create-on-a-dedicated-server.md](server/uc-19-create-on-a-dedicated-server.md) |
| UC-20 | Play with deck ratios | `internal/server`, `rules_engine`, `frontend` | Implemented | [server/uc-20-deck-ratios.md](server/uc-20-deck-ratios.md) |
| UC-21 | Know when the game waits for me | `frontend` | Implemented | [app/uc-21-know-when-the-game-waits.md](app/uc-21-know-when-the-game-waits.md) |
| UC-22 | See why I can't act | `frontend`, `rules_engine` | Implemented | [app/uc-22-see-why-i-cant-act.md](app/uc-22-see-why-i-cant-act.md) |
| UC-23 | See what the active player does | `frontend`, `rules_engine` | Implemented | [app/uc-23-see-what-the-active-player-does.md](app/uc-23-see-what-the-active-player-does.md) |
| UC-24 | Follow fights and played cards | `frontend` | Implemented | [app/uc-24-follow-fights-and-played-cards.md](app/uc-24-follow-fights-and-played-cards.md) |
| UC-25 | See what others aim at | `frontend`, `internal/server` | Implemented | [app/uc-25-see-what-others-aim-at.md](app/uc-25-see-what-others-aim-at.md) |
| UC-26 | Roll for the first player | `rules_engine`, `internal/server`, `frontend` | Implemented | [app/uc-26-roll-for-the-first-player.md](app/uc-26-roll-for-the-first-player.md) |
| UC-27 | Join a server by its domain | `frontend` | Implemented | [server/uc-27-join-by-domain.md](server/uc-27-join-by-domain.md) |
| UC-28 | Play on after a server update | `internal/server`, `rules_engine`, `frontend` | Implemented | [server/uc-28-play-on-after-a-server-update.md](server/uc-28-play-on-after-a-server-update.md) |
| UC-29 | A game ends by an idle limit | `internal/server`, `frontend` | Implemented | [server/uc-29-idle-game-limits.md](server/uc-29-idle-game-limits.md) |
| UC-30 | Roll back a mistake (Godmode) | `internal/server`, `pkg/record`, `frontend` | Implemented | [app/uc-30-roll-back-a-mistake.md](app/uc-30-roll-back-a-mistake.md) |
