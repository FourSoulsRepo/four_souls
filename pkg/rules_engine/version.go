package rulesengine

// Version is the rules engine version.
// Bump it with every change to rules behavior; see docs/roadmap/08-first-release.md.
const Version = "0.5.0"

// StateMeaning numbers what the saved state's fields mean (N-20). Raise
// it when a field keeps its JSON name and type but means something new:
// Resume cannot see that, and games kept over a server restart from
// before the change then end instead of going on wrong.
// 2: prompts, stack labels and reasons are messages, not English (16.4).
const StateMeaning = 2

// throwVersion is the first engine version where every roll waits for
// its roller's throw, the first-player roll too (15.43, ADR 014). Older
// records cannot replay.
const throwVersion = "0.3.0"

// msgVersion is the first engine version whose state holds messages
// (key and values) instead of English text (16.4, ADR 021). Older
// records' checksums differ: they cannot replay.
const msgVersion = "0.4.0"

// chainVersion is the first engine version whose step checksums are a
// hash chain over the history, with the whole state checked every
// anchorEvery steps (RP-12). Older records' checksums cover the whole
// state at every step: they cannot replay.
const chainVersion = "0.5.0"

// orderVersion is the first engine version where a dead monster's slot
// refills after its roll rewards (19.104) and a roll's throw waits for
// the triggers already waiting (19.105). Older records' checksums differ:
// they cannot replay.
const orderVersion = "0.4.2"
