package rulesengine

// Version is the rules engine version.
// Bump it with every change to rules behavior; see docs/roadmap/08-first-release.md.
const Version = "0.2.0"

// throwVersion is the first engine version where every roll waits for
// its roller's throw (15.59, ADR 014). Older records cannot replay.
const throwVersion = "0.2.0"
