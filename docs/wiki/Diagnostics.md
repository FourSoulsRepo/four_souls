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
2. **One row per variant**: frames per second on the main menu, how many frames came late (over 34 ms), the longest wait, the screen scale, whether the sky moved, and the GPU name the webview reports.
3. **The verdict**: the fastest variant, how to start the game with it, and any problems found.

The Linux variants:

| Variant | What changes |
|---|---|
| `default` | Nothing: a normal start. |
| `no-dmabuf` | `WEBKIT_DISABLE_DMABUF_RENDERER=1`; often fixes Intel and NVIDIA. |
| `no-compositing` | `WEBKIT_DISABLE_COMPOSITING_MODE=1`; helps on some systems, hurts on others. |
| `gpu-always` | The webview always uses the GPU. |
| `gpu-never` | The webview never uses the GPU; for comparison. |

When the best variant is an environment switch, the report prints the command, for example `WEBKIT_DISABLE_DMABUF_RENDERER=1 ./FourSouls`. The two `gpu-*` variants can't be chosen yet; a later version will let the game pick the best one itself.

## Problems it finds

* **The webview draws on the CPU (`llvmpipe`).** The GPU driver is missing or not working; on Linux check that Mesa is installed and up to date.
* **A GPU does not open.** Your user may not use `/dev/dri/renderD*`. Add yourself to the `render` group (`sudo usermod -aG render $USER`) and log in again.
* **No GPU render node.** No GPU driver is loaded at all.
* **No MP3 playback.** The game has no sound. Install the GStreamer plugins: `gstreamer1.0-plugins-good` and `gstreamer1.0-libav` (Debian, Ubuntu), `gstreamer1-plugins-good` (Fedora), `gst-plugins-good` and `gst-libav` (Arch).
* **Reduced motion.** The system asks apps not to animate, so the sky stands still and the frame rate says less.
* **Every variant is slow.** None reached 40 fps; send the report to the developers.
