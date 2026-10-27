# Banned cards

Cards that no game deals, for every host and every game. This page is the list of record; keep it in step with the code.

## Banned in every game

| Card | Ref | Why | Bug | Since |
|------|-----|-----|-----|-------|
| Magic Marker | `magic_marker` | Does not work properly yet. | 19.96 | 2026-10-28 |

## Where the list lives

* `alwaysBanned` in `pkg/rules_engine/ratios.go`: the cards no deck gets.
* `inDecks` checks it with the host bans, so both deck paths skip the card.
* `TestMagicMarkerIsLeftOut` in `pkg/rules_engine/cards/r/magic_marker_test.go` checks it.

## Host bans (GS-13)

* A host can also leave cards out of one game. Those bans are per game and are not listed here.
* A banned character is never dealt. A starting item cannot be banned; ban its character instead.

## Lifting a ban

1. Fix the card and its bug item in `docs/roadmap/19-bugs.md`.
2. Remove the card from `alwaysBanned`.
3. Remove its row above.
4. Remove or update `TestMagicMarkerIsLeftOut` for the card.
5. Run `make test` and `make lint`.
