# 35. Webview diagnostics

* **Status:** Accepted
* **Date:** 2026-11-18
* **Authors:** @HardDie

---

## Context

1. The app draws with the system webview (Wails v2, B-02).
   1. Linux: WebKitGTK; macOS: WKWebView; Windows: WebView2.
2. Owner, 2026-11-18: the main menu is slow on a Linux laptop.
   1. Integrated Intel GPU; the drifting sky (15.42) behind the buttons.
3. WebKitGTK speed depends on its setup.
   1. Its version; Wayland or X11; the GPU driver; access to `/dev/dri`.
   2. Switches read at start: `WEBKIT_DISABLE_DMABUF_RENDERER`, `WEBKIT_DISABLE_COMPOSITING_MODE`.
   3. Wails `WebviewGpuPolicy`: `OnDemand`, `Always` (the default since ADR 036), `Never`.
4. No switch is right for every computer; it must be measured.
5. A switch only works before the webview starts.
   1. One process cannot try several.
6. Players often start the app without a terminal.
7. A move to Electron waits for these results (owner, 2026-11-18).

## Considered options

1. **A `--diagnose` flag in the game binary**
   1. Won. It tests the exact webview and build players run.
2. **A separate diagnostic app**
   1. Lost. Another build may link another WebKitGTK (`webkit2_41` tag).
   2. Twice the release work.
3. **Measure in child processes, one per variant**
   1. Won. Context 5: each switch needs a fresh webview.
4. **Only read the system, no measuring**
   1. Lost. The same setup can be fast or slow; only frames tell.
5. **Read the WebKitGTK version through cgo**
   1. Won. The user agent hides it; Wails already links the library.
6. **Read GPU facts from `/sys` and `/dev`**
   1. Won. Tools like `glxinfo` are often missing.

## Decision

Use options 1, 3, 5 and 6.

1. `FourSouls --diagnose` writes a report and exits.
   1. The flag is read from `os.Args` before Wails starts.
   2. No instance slot, no tracing, no saved settings.
2. System part: `internal/diag`, plain Go.
   1. OS, architecture, app version.
   2. Linux: session type, desktop, `DISPLAY`, `WAYLAND_DISPLAY`.
   3. Linux: WebKitGTK version (cgo; `webkit2_41` picks 4.1 or 4.0).
   4. Linux: each `/dev/dri/renderD*`, its driver, vendor, and if it opens.
   5. Linux: the user's `render` and `video` groups.
   6. Switches already set in the environment.
3. Probe part: the app starts itself with `--diagnose-probe=<variant>`.
   1. Linux variants: `default`, `no-dmabuf`, `no-compositing`, `gpu-ondemand`, `gpu-never`.
   2. Other systems: `default` only.
   3. A probe opens a 1280×800 window on the main menu over the sky.
   4. It skips the splash and the nickname screen.
   5. It forces normal animation speed; nothing is saved.
   6. It waits 1 s, then counts frames for 5 s.
   7. It sends back: frames per second, slow frames, WebGL renderer, pixel ratio, sound formats.
   8. The parent waits at most 30 s per probe.
   9. Each variant is environment switches only (ADR 036 item 4).
   10. WebKit (Linux, macOS) hides the WebGL renderer: "Apple GPU" (19.228).
       1. The report says "hidden by WebKit" there; only Windows shows it.
4. Report
   1. Plain text, to stdout and to `diagnose-<time>.txt` in the app folder.
   2. The last line names the file.
   3. It ends with the variant with the most frames per second.
   4. Red flags: software rendering (`llvmpipe`), no render node, no group access.
      1. `llvmpipe` only where the renderer is not hidden (19.228).
      2. Also `LIBGL_ALWAYS_SOFTWARE=1` in the environment.
5. The binding `bindings/probe` is always bound (ADR 030).
   1. `Variant()`: `""` in a normal start.
   2. `Report(result)`: only in a probe; prints and quits.
   3. Not `bindings/diag`: the `diag` model namespace would merge (ADR 030).

## Consequences

### Positive

1. The owner and players get facts in one command.
2. The same probe can drive the later self-repair (25.3).
3. A normal start is unchanged.

### Negative and risks

1. The run opens up to five windows, one after another (about 40 s).
2. The probe measures the menu only, not a full table.
3. cgo ties `internal/diag` on Linux to the WebKitGTK headers.
   1. Builds without cgo report the version as unknown.

### Neutral

1. Self-repair (25.3) gets its own decision later.
