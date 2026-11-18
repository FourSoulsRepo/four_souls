# Step 25. Diagnostics

Goal: find why the webview is slow on a player's computer; later fix it in the game.

Ideas: B-09.

## Why (2026-11-18)

1. The owner's Linux laptop: the main menu is slow.
   1. A strong CPU with an integrated Intel GPU.
   2. The menu has buttons over the drifting sky (15.42).
2. Wails draws with the system webview (B-02).
   1. Linux: WebKitGTK; macOS: WKWebView; Windows: WebView2.
   2. WebKitGTK is the slowest; often a setup problem, not a limit.
3. The sky is six window-sized layers moved by `transform`.
   1. Cheap on the GPU.
   2. Heavy when WebKitGTK composites on the CPU.
4. Known WebKitGTK switches, set before the window opens:
   1. `WEBKIT_DISABLE_DMABUF_RENDERER=1`: a common fix on Intel and NVIDIA.
   2. `WEBKIT_DISABLE_COMPOSITING_MODE=1`: helps on some, hurts on others.
   3. Wails `linux.Options.WebviewGpuPolicy`; `main.go` sets none (`OnDemand`).
5. A move to Electron waits for these results (owner, 2026-11-18).

## Rules for this step (owner, 2026-11-18)

1. Branch: `step-25-diagnostics`, from `main`.
2. 25.2 first; the owner runs it on the Linux laptop.
3. 25.3 waits for those results and the owner's go.
4. Design: ADR 035.
5. Gates: `make lint test`, `make frontend e2e`, `wails build -m -tags embed`.

---

### [x] 25.1 Plan and ADR

1. Goal: this file, idea B-09, ADR 035.
2. Done when:
   1. The ADR is Accepted and its row is in `INDEX.md`.
3. Result (2026-11-18)
   1. ADR 035; roadmap row and section 25; idea B-09.

### [ ] 25.2 The `--diagnose` flag

1. Goal: one command tells what the webview runs on and how fast it draws.
2. Tasks:
   1. `FourSouls --diagnose` runs the report, then exits; no game window.
   2. System part in Go, package `internal/diag` (ADR 035 item 2).
   3. Probe part: the app starts itself once per variant (ADR 035 item 3).
   4. Binding `bindings/diag`: the variant to the frontend, the result back.
   5. Frontend `probe.ts`: menu on the sky, frames measured, result sent.
   6. Report to the terminal and a file in the app folder (ADR 035 item 4).
   7. Wiki page Diagnostics, in the players group.
   8. Use case: run the diagnostics.
3. Done when:
   1. The flag runs on macOS here; the owner runs it on Linux.
   2. A normal start is unchanged; every e2e test passes.
   3. The report names the best variant by frame rate.

### [ ] 25.3 Self-repair in the game (later)

1. Goal: the game finds a slow webview and fixes itself.
2. Waits for the owner's 25.2 results; not to start before the owner's go.
3. Draft, to settle then:
   1. First start: the menu measures its own frame rate.
   2. Low: offer "compatibility mode"; it saves the best switch.
   3. The app restarts with it; `main.go` sets it before the window.
   4. A setting turns it off again.
   5. Still low: show what the report found, with a fix per distro.
4. **Owner** questions:
   1. Which switches helped on the laptop?
   2. Ask the player first, or switch by itself?
   3. Still slow after the switch: Electron, or lighter effects?
5. Done when:
   1. A new ADR or an update of ADR 035 holds the decision.
