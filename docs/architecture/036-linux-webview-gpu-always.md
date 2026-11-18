# 36. Linux webview always on the GPU

* **Status:** Accepted
* **Date:** 2026-11-18
* **Authors:** @HardDie

---

## Context

1. Wails sets the WebKitGTK GPU policy.
   1. `main.go` set none until 0.6.2.
   2. Without Linux options Wails v2 uses `Never` (`window.go`, v2.17.0).
2. The owner's Linux laptop, `--diagnose` of 0.6.2 (25.2, ADR 035):
   1. WebKitGTK 2.54.2, Wayland, Intel `i915`; the GPU opens.
   2. No options, so `Never`: 6.3 fps on the main menu, every frame slow.
   3. `Always`: 60.0 fps, no slow frame.
   4. `Never`, and the two WebKitGTK switches: 6.2–6.3 fps.
3. So every start until 0.6.2 drew on the CPU.
4. A broken GPU driver may show a black window with `Always`.
5. Owner, 2026-11-18: `Always` is the Linux default; keep `--diagnose`.
6. 0.6.3 on the same laptop (2026-11-19):
   1. `default` (`Always`) 58.2 fps; `gpu-ondemand` 60.0; `gpu-never` 6.3.
   2. WebKitGTK 2.54 logs that `OnDemand` now equals `Always`.
   3. `no-dmabuf` with `Always`: 5.9 fps; `no-compositing`: 60.0.

## Considered options

1. **`Always` by default, with an environment switch back**
   1. Won. Fixes the measured case; a broken driver has a way out.
2. **Measure on the first start, switch when slow (self-repair)**
   1. Lost for now. More code for a case no one has reported yet.
   2. Parked as 25.4.
3. **A setting on the Settings screen**
   1. Lost. A player whose window is black can't reach the screen.
4. **Keep no options (Wails: `Never`)**
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
   2. On WebKitGTK 2.54 `ondemand` equals `always`; only `never` is a way out.
2. Only one computer measured so far.

### Neutral

1. Self-repair waits for a report that needs it (25.4).
