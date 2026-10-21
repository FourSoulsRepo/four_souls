package rulesengine

// Version is the rules engine version.
// Bump it with every change to rules behavior; see docs/roadmap/08-first-release.md.
const Version = "0.3.0"

// StateMeaning numbers what the saved state's fields mean (N-20). Raise
// it when a field keeps its JSON name and type but means something new:
// Resume cannot see that, and games kept over a server restart from
// before the change then end instead of going on wrong.
const StateMeaning = 1

// throwVersion is the first engine version where every roll waits for
// its roller's throw, the first-player roll too (15.43, ADR 014). Older
// records cannot replay.
const throwVersion = "0.3.0"
