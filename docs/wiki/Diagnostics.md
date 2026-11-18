# Diagnostics

If the game feels slow — the main menu stutters, cards drag behind the mouse — run the built-in check and send its report to the developers. It matters most on Linux, where the game draws with WebKitGTK and its speed depends on how the system is set up.

## Run it

Start the game from a terminal with `--diagnose`:

| System | Command |
|---|---|
| Linux | `./FourSouls --diagnose` |
| macOS | `./FourSouls.app/Contents/MacOS/FourSouls --diagnose` |
| Windows | `FourSouls.exe --diagnose` (the report goes to the file only) |

The check opens a small test window that shows the main menu for about six seconds, then closes it. On Linux it does this five times, once per way of starting the webview, so the whole run takes under a minute. Don't touch the windows while it runs. The game itself does not start.

## The report

The report prints to the terminal and is saved as `diagnose-<date>-<time>.txt` in the app's settings folder; its last line names the file. It has three parts:

1. **The system** (Linux): the WebKitGTK version, Wayland or X11, each GPU in `/dev/dri` with its driver and whether the game may use it, and your `render` and `video` groups.
2. **One row per variant**: frames per second on the main menu, how many frames came late (over 34 ms), the longest wait, the screen scale, whether the sky moved, and the GPU name the webview reports. WebKit (Linux, macOS) hides that name and always says "Apple GPU", so there the column reads "hidden by WebKit"; the GPU lines above it name the real one.
3. **The verdict**: the fastest variant, how to start the game with it, and any problems found.

The Linux variants:

| Variant | What changes |
|---|---|
| `default` | Nothing: a normal start. The webview always uses the GPU. |
| `no-dmabuf` | `WEBKIT_DISABLE_DMABUF_RENDERER=1`; fixes some Intel and NVIDIA setups. |
| `no-compositing` | `WEBKIT_DISABLE_COMPOSITING_MODE=1`; helps on some systems, hurts on others. |
| `gpu-ondemand` | `FOURSOULS_GPU_POLICY=ondemand`: the GPU only when the page asks. WebKitGTK 2.54 and later treat it as always. |
| `gpu-never` | `FOURSOULS_GPU_POLICY=never`: the webview never uses the GPU; every version before 0.6.3 ran like this. |

When the best variant is not `default`, the report prints the command, for example `FOURSOULS_GPU_POLICY=ondemand ./FourSouls`.

## The GPU switch

Since 0.6.3 the game on Linux always draws with the GPU. On one tested laptop the old setting drew the menu at 6 frames a second instead of 60. If your window stays black or looks broken, your GPU driver may not cope; start the game with `FOURSOULS_GPU_POLICY=never ./FourSouls` (`ondemand` also works on WebKitGTK older than 2.54), and send the `--diagnose` report to the developers.

## Problems it finds

* **The webview draws on the CPU (`llvmpipe`).** Windows only, where the webview names its GPU. The GPU driver is missing or not working.
* **`LIBGL_ALWAYS_SOFTWARE=1` is set.** OpenGL draws on the CPU; remove the variable.
* **A GPU does not open.** Your user may not use `/dev/dri/renderD*`. Add yourself to the `render` group (`sudo usermod -aG render $USER`) and log in again.
* **No GPU render node.** No GPU driver is loaded at all.
* **No MP3 playback.** The game has no sound. Install the GStreamer plugins: `gstreamer1.0-plugins-good` and `gstreamer1.0-libav` (Debian, Ubuntu), `gstreamer1-plugins-good` (Fedora), `gst-plugins-good` and `gst-libav` (Arch).
* **Reduced motion.** The system asks apps not to animate, so the sky stands still and the frame rate says less.
* **Every variant is slow.** None reached 40 fps; send the report to the developers.
