# 36. Linux webview always on the GPU

* **Status:** Accepted
* **Date:** 2026-11-18
* **Authors:** @HardDie

---

## Context

1. Wails sets the WebKitGTK GPU policy; its default is `OnDemand`.
   1. `main.go` set none until 0.6.2.
2. The owner's Linux laptop, `--diagnose` of 0.6.2 (25.2, ADR 035):
   1. WebKitGTK 2.54.2, Wayland, Intel `i915`; the GPU opens.
   2. `OnDemand`: 6.3 fps on the main menu, every frame slow.
   3. `Always`: 60.0 fps, no slow frame.
   4. `Never`, and the two WebKitGTK switches: 6.2–6.3 fps.
3. So `OnDemand` drew on the CPU, as `Never` does.
4. A broken GPU driver may show a black window with `Always`.
5. Owner, 2026-11-18: `Always` is the Linux default; keep `--diagnose`.

## Considered options

1. **`Always` by default, with an environment switch back**
   1. Won. Fixes the measured case; a broken driver has a way out.
2. **Measure on the first start, switch when slow (self-repair)**
   1. Lost for now. More code for a case no one has reported yet.
   2. Parked as 25.4.
3. **A setting on the Settings screen**
   1. Lost. A player whose window is black can't reach the screen.
4. **Keep `OnDemand`**
   1. Lost. Context 2.

## Decision

Use option 1.

1. Linux starts the webview with `WebviewGpuPolicyAlways`.
2. `FOURSOULS_GPU_POLICY` changes it: `always`, `ondemand`, `never`.
   1. Unset: `always`.
   2. An unknown value: `always`, and a log line says so.
3. Other systems ignore the policy.
4. `--diagnose` measures this default as `default`.
   1. `gpu-ondemand` and `gpu-never` set the variable.
   2. So the report prints how to start with any variant.

## Consequences

### Positive

1. The menu runs at full speed on the measured laptop.
2. One variable fixes a broken driver, without a new build.

### Negative and risks

1. A driver that fails with `Always` shows a black or broken window.
   1. The player must find the variable (wiki Diagnostics, First start).
2. Only one computer measured so far.

### Neutral

1. Self-repair waits for a report that needs it (25.4).
