# 19. Seamless server restart

* **Status:** Accepted
* **Date:** 2026-10-21
* **Authors:** @HardDie

---

## Context

1. A server update or restart ended every game (N-20).
2. Games can last an hour; some are left open and never end.
3. Players should notice a restart only as a short pause.
4. The engine already saves and loads a whole game (N-11).
5. A seat answers to its token; reconnecting takes it back (N-08).

## Considered options

1. **Restart files written on stop**
   1. Won. The N-11 save, written for every table and game.
   2. The next start loads them; players reconnect by token.
2. **Drain: no new games, restart when the last one ends**
   1. Lost. Long and abandoned games would hold an update for hours.
3. **Pass the listening socket to the new process**
   1. Lost. The state must move anyway; it saves only a second or two.
4. **Replay the record on the new engine**
   1. Lost. A rules change makes the replay differ from what was played.

## Decision

Use option 1.

1. Only the dedicated server keeps games.
   1. `Config.KeepGames`; `cmd/server` sets it.
   2. An app's own host ends its games when it stops.
2. On stop
   1. HTTP stops taking connections.
   2. Every client gets `restarting`; 300 ms later the rooms stop.
   3. Each table and game goes to `<saves>/restart/<game>.fsrestart`.
   4. A game's record gets no `end` line (`record.Writer.Stop`).
   5. A finished game and a saved game's table are not kept.
3. A restart file
   1. The save file fields, plus kind, number, host, start time.
   2. Kinds: `table`, `setup`, `game`.
   3. Seats keep name, token, kicked, dice.
   4. A table keeps its password; the file is deleted on load.
   5. A game keeps the engine state and its checksum.
4. On start
   1. Tables and rooms come back before the hub takes requests.
   2. Game numbers go on after the highest kept one.
   3. A game waits paused until every seat is back (N-08).
   4. Its record goes on with a `resume` part (ADR 007).
   5. Bans and picks start over once every player is back.
   6. They use the same seed, so a draft offers the same cards.
   7. A loaded restart file is deleted.
5. A game the new server cannot go on with
   1. It ends; its record gets `end` with finished false.
   2. Its file is renamed to `.fsrestart.lost` for the operator.
6. The protocol version stays.
   1. `restarting` and `notice` are new messages.
   2. An older app ignores them.

## Consequences

### Positive

1. An update keeps every game; players wait a few seconds.
2. No new storage: the saves folder and the N-11 save format.

### Negative and risks

1. The saves folder must survive the restart (a Docker volume).
2. A response timer starts again in full after the pause.
3. Bans and picks made before the restart are lost.
4. A crash or `kill -9` writes no restart files.

### Neutral

1. The protocol bump rules decide which apps can rejoin.
