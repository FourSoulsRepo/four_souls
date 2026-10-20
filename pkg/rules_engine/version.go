package rulesengine

// Version is the rules engine version.
// Bump it with every change to rules behavior; see docs/roadmap/08-first-release.md.
const Version = "0.3.0"

// throwVersion is the first engine version where every roll waits for
// its roller's throw, the first-player roll too (15.43, ADR 014). Older
// records cannot replay.
const throwVersion = "0.3.0"
