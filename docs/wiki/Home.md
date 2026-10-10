# Four Souls

Developer notes for the Four Souls fan game: a desktop network card game for 2–4 players, built with Wails (Go + React/TypeScript).

These pages describe the current code. Decisions live in `docs/architecture/` (ADRs); user-visible behavior in `docs/use-cases/`; the plan in `docs/roadmap.md`.

## Pages

* [Repo layout](Repo-layout) — where everything lives and why.
* [Build modes](Build-modes) — embedded vs. external files, versions, platforms.
* [Linters](Linters) — strict Go and frontend linting, and how to run it.
* [Screenshots](Screenshots) — capture the UI without a real window.
* [Card data tools](Card-data-tools) — where card data and images come from.
* [Effect blocks](Effect-blocks) — the blocks cards are built from.
* [Situation payload](Situation-payload) — sandbox situations and rules tests.
* [Fuzzing](Fuzzing) — random games that check the engine state.
* [Rules engine](Rules-engine) — map of the engine code and its main rules.
* [Adding a card](Adding-a-card) — write and test a card, or a whole fan set.
* [Playing over the internet](Playing-over-the-internet) — virtual LAN (Tailscale, ZeroTier, Radmin VPN, Hamachi) or port forwarding.
* [Dedicated server on a VPS](Dedicated-server-on-a-VPS) — rent a server and host games for friends.
* [Server](Server) — how the lobby, game rooms and transport fit together.
