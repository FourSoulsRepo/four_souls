# 32. Fixed deals in e2e tests

* **Status:** Accepted
* **Date:** 2026-11-16
* **Authors:** @HardDie

---

## Context

1. Step 19 had 163 items by 2026-11-16; 71 were e2e test faults.
   1. About 30 broke because the random deal gave an unplanned card.
   2. Pride, Bob's Brain, Eden, Delirium, Cursed Mom's Hand, and more.
2. Host bans (14.10) fixed them one card at a time.
   1. Each new set adds new cards that can break any test.
3. No card bug was ever found by a random e2e deal.
   1. Card bugs came from play tests, records and the load tool.
4. The server seeds a game from the clock.
   1. `Config.Seed` gives game n the seed `Seed+n`.
   2. In e2e, n depends on test order: not stable per test.
5. A client must not pick its seed on a public server.
   1. It could rig its own deal.

## Considered options

1. **A seed per test, from a test-only create option**
   1. Won. Each test sees the same deal every run.
   2. A failure reproduces on the first rerun.
2. **Keep random deals; ban more cards**
   1. Lost. Endless; each set adds risk.
3. **A fixed seed per test, plus a random-seed env switch**
   1. Lost (owner, 2026-11-16). Random deals belong to fuzz and full games.
4. **Scripted start states from save files for every test**
   1. Lost. A lot of setup per test; saves already serve tests that need it.

## Decision

Use option 1.

1. Server
   1. Flag `-test_seeds` sets `Config.TestSeeds`; default false.
   2. Create option `seed` (`protocol.Options.Seed`, uint64).
   3. With `TestSeeds`, a non-zero `seed` is the game's seed.
   4. Without it, a create with a non-zero `seed` is refused.
   5. The app's own host never sets `TestSeeds`.
2. e2e
   1. The Playwright server command passes `-test_seeds`.
   2. `startGame` sends a seed by default.
   3. The seed is a hash of the test's title path: stable, unique per test.
   4. An option `seed` overrides it.
   5. `seed: 0` asks for a clock seed.
3. Random deals stay only where finding surprises is the point.
   1. `games.spec.ts` (`make e2e-games`) uses `seed: 0` or `GAME_SEED`.
   2. Fuzz and `make sims` are unchanged.
4. Bans stay where a test already has them, with their reasons.

## Consequences

### Positive

1. A test's deal never changes between runs of one commit.
2. A failed test fails again on rerun; its cause can be found.
3. New sets no longer break unrelated e2e tests.

### Negative and risks

1. An engine change that draws randomness in another order changes every deal.
   1. Some tests may then fail at once, all on one commit.
   2. Each is fixed like any test fault: the test assumed its deal.
2. e2e no longer tries many deals by chance.
   1. `make e2e-games`, fuzz and sims keep that job.

### Neutral

1. The match record already stores the seed in its header.
