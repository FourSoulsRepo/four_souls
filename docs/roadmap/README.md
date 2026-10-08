# Detailed roadmap

One file per global step from [../roadmap.md](../roadmap.md).
Each sub-step is sized for one agent session.

## Sub-step format

```markdown
### N.M Title

1. Goal: one line.
2. Tasks:
   1. …
3. Done when:
   1. …
```

## Rules for every sub-step

1. Read `CLAUDE.md`, `docs/ideas.md`, and the step file first.
2. Keep it simple (PR-01); one task, done well.
3. Tests come with the code, not after.
4. CI stays green, including the engine `js/wasm` build.
5. Docs in the same change:
   1. New decision: ADR in `docs/architecture/`.
   2. New user-visible behavior: use case in `docs/use-cases/`.
   3. Package notes: `docs/wiki/`.
6. Questions marked **Owner** need the owner's answer first.
7. Mark the sub-step done in its file when finished: `[x]`.
8. Commit after every finished sub-step.
9. Stop after every finished global step.
   1. Wait for the owner's review and approval.
   2. Never start the next step without approval.
10. Keep `NOTICE` current.
   1. Any added, removed or upgraded third-party item updates it.
   2. Same commit as the change.

## Files

| Step | File |
|------|------|
| 1 | [01-foundation.md](01-foundation.md) |
| 2 | [02-card-data.md](02-card-data.md) |
| 3 | [03-rules-research.md](03-rules-research.md) |
| 4 | [04-engine-core.md](04-engine-core.md) |
| 5 | [05-base-game-cards.md](05-base-game-cards.md) |
| 6 | [06-server.md](06-server.md) |
| 7 | [07-client.md](07-client.md) |
| 8 | [08-first-release.md](08-first-release.md) |
| 9 | [09-spectators.md](09-spectators.md) |
| 10 | [10-replays.md](10-replays.md) |
| 11 | [11-learning.md](11-learning.md) |
| 12 | [12-rules-website.md](12-rules-website.md) |
| 13 | [13-web-client.md](13-web-client.md) |
| 14 | [14-more-sets.md](14-more-sets.md) |
