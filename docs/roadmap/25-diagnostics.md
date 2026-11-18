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
   3. Wails `linux.Options.WebviewGpuPolicy`; `main.go` sets none (Wails then uses `Never`).
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

### [x] 25.2 The `--diagnose` flag

1. Goal: one command tells what the webview runs on and how fast it draws.
2. Tasks:
   1. `FourSouls --diagnose` runs the report, then exits; no game window.
   2. System part in Go, package `internal/diag` (ADR 035 item 2).
   3. Probe part: the app starts itself once per variant (ADR 035 item 3).
   4. Binding `bindings/probe`: the variant to the frontend, the result back.
   5. Frontend `probe.ts`: menu on the sky, frames measured, result sent.
   6. Report to the terminal and a file in the app folder (ADR 035 item 4).
   7. Wiki page Diagnostics, in the players group.
   8. Use case: run the diagnostics.
3. Done when:
   1. The flag runs on macOS here; the owner runs it on Linux.
   2. A normal start is unchanged; every e2e test passes.
   3. The report names the best variant by frame rate.
4. Result (2026-11-18); waits for the owner's Linux run.
   1. `internal/diag`: `Args`, `Variants`, `ReadSystem`, `ProbeAll`, `Write`, `Flags`.
      1. WebKitGTK version by cgo: `webkit_linux_41.go`, `webkit_linux_40.go`.
      2. Without Linux or cgo: `webkit_other.go`, version unknown.
   2. `diagnose.go` in `main`: `runDiagnose`, `runProbe`, `linuxOptions`.
   3. `main.go`: `appOptions` builds a normal start and a probe alike.
      1. A normal start passes no variant: `Linux` stays nil, as before.
   4. `bindings/probe`: `Variant`, `Report`; always bound.
      1. Departure: not `bindings/diag`; ADR 035 item 5.3.
   5. Frontend: `probe.ts`, `probe.test.ts`; `App.tsx` starts the probe.
      1. Bridge: `probeVariant`, `sendProbe`; `fakeBridge.ts` fakes `window.go.probe.Probe`.
   6. Generated: `wailsjs/go/probe/Probe`; `models.ts` namespace `diag` only.
      1. Wails also reordered `records.Seat` and file modes; not kept.
   7. Docs: wiki Diagnostics (players), First-start link, Client, Repo-layout.
      1. UC-48 and its index row; `CLAUDE.md` 4.28.
   8. Checks on macOS:
      1. With images: `--diagnose` 60.0 fps, sky moving, report saved.
      2. `wails build -m -tags embed`: 59.9 fps, sky off; a normal start runs.
      3. Linux cgo files not built here; CI lint and build cover them.
   9. Gates: `make lint`, `make test`, `make frontend` (Vitest 757 passed, 1 skipped).
      1. `make e2e`: 633 passed, 1 failed, 22 skipped.
      2. Failure: `arrow.spec.ts:62`, chromium; logged under 19.153.
      3. Rerun: 5 separate runs of the spec, 30 of 30 passed.
5. Owner's Linux run (0.6.2, 2026-11-18)
   1. WebKitGTK 2.54.2, Wayland, Intel `i915`; the GPU opens; scale 1.
   2. `default`, `no-dmabuf`, `no-compositing`, `gpu-never`: 6.2–6.3 fps, every frame slow.
   3. `gpu-always`: 60.0 fps, no slow frame.
   4. So no options (Wails: `Never`) draws on the CPU here; the switches change nothing.
   5. Renderer "Apple GPU" everywhere: WebKit masks it; logged as 19.228.
   6. "No MP3 playback": GStreamer can't play `audio/mpeg`; owner to confirm.

### [ ] 25.3 Linux webview always on the GPU

1. Goal: the game draws at full speed on the owner's Linux laptop.
2. Owner answers (2026-11-18)
   1. Which switches helped? `gpu-always` only (25.2 item 5).
   2. `Always` becomes the Linux default; no self-repair now.
   3. Keep `--diagnose`; it may help later.
3. Tasks:
   1. ADR 036: the default and `FOURSOULS_GPU_POLICY`.
   2. `--diagnose`: `gpu-always` becomes `gpu-ondemand`; variants are env only.
   3. Fix 19.228 in the same change.
   4. Wiki Diagnostics, First start; UC-48; changelog.
4. Done when:
   1. The owner's laptop runs the menu smoothly in a normal start.
   2. The owner checks the sound (the MP3 warning, 25.2 item 5.6).

### [ ] 25.4 Self-repair in the game (parked)

1. Goal: the game finds a slow webview and fixes itself.
2. Parked (owner, 2026-11-18): waits for a report `Always` does not fix.
3. Draft, to settle then:
   1. First start: the menu measures its own frame rate.
   2. Low: offer a switch; save the best one.
   3. The app restarts with it; `main.go` sets it before the window.
   4. A setting turns it off again.
   5. Still low: show what the report found, with a fix per distro.
4. Done when:
   1. A new ADR holds the decision.
