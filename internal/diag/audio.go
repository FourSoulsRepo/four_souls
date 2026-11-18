package diag

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Audio is what GStreamer gives WebKitGTK to play sounds with (19.229).
// Without autoaudiosink the page's web process hangs at the first sound.
type Audio struct {
	// Dirs are the GStreamer plugin folders found; none means unknown.
	Dirs []string
	// AutoSink: autoaudiosink is there (the good plugins, autodetect).
	AutoSink bool
	// Outputs are the sound systems it can pick: pulseaudio, pipewire, alsa.
	Outputs []string
	// Decoders are the MP3 decoders: libav, mpg123.
	Decoders []string
}

// Known tells that GStreamer's plugins were found at all.
func (a Audio) Known() bool { return len(a.Dirs) > 0 }

// Plays tells that WebKitGTK can play the game's sounds.
func (a Audio) Plays() bool {
	return a.AutoSink && len(a.Outputs) > 0 && len(a.Decoders) > 0
}

// gstPlugins names the plugin files the game needs, by what they give.
var (
	gstAutoSink = "libgstautodetect.so"
	gstOutputs  = map[string]string{
		"libgstpulseaudio.so": "pulseaudio",
		"libgstpipewire.so":   "pipewire",
		"libgstalsa.so":       "alsa",
	}
	gstDecoders = map[string]string{
		"libgstlibav.so":  "libav",
		"libgstmpg123.so": "mpg123",
	}
)

// gstDirPatterns are where distributions keep GStreamer 1.0 plugins;
// Debian and Ubuntu put them under a multiarch folder.
var gstDirPatterns = []string{
	"/usr/lib/gstreamer-1.0",
	"/usr/lib64/gstreamer-1.0",
	"/usr/lib/*/gstreamer-1.0",
	"/usr/local/lib/gstreamer-1.0",
	"/usr/local/lib64/gstreamer-1.0",
	"/usr/local/lib/*/gstreamer-1.0",
}

// gstEnv are the variables that move GStreamer's plugin search.
var gstEnv = []string{"GST_PLUGIN_SYSTEM_PATH_1_0", "GST_PLUGIN_SYSTEM_PATH", "GST_PLUGIN_PATH_1_0", "GST_PLUGIN_PATH"}

// ReadAudio looks for the plugins in this computer's GStreamer folders.
func ReadAudio() Audio {
	var dirs []string
	for _, name := range gstEnv {
		for d := range strings.SplitSeq(os.Getenv(name), string(os.PathListSeparator)) {
			if d != "" {
				dirs = append(dirs, d)
			}
		}
	}
	for _, p := range gstDirPatterns {
		found, _ := filepath.Glob(p) //nolint:errcheck // only a bad pattern fails
		dirs = append(dirs, found...)
	}
	return audioIn(dirs)
}

// audioIn checks the plugin files in dirs; folders that don't exist
// are left out.
func audioIn(dirs []string) Audio {
	var a Audio
	for _, d := range dirs {
		if info, err := os.Stat(d); err != nil || !info.IsDir() || slices.Contains(a.Dirs, d) {
			continue
		}
		a.Dirs = append(a.Dirs, d)
		has := func(file string) bool {
			_, err := os.Stat(filepath.Join(d, file))
			return err == nil
		}
		a.AutoSink = a.AutoSink || has(gstAutoSink)
		for file, name := range gstOutputs {
			if has(file) && !slices.Contains(a.Outputs, name) {
				a.Outputs = append(a.Outputs, name)
			}
		}
		for file, name := range gstDecoders {
			if has(file) && !slices.Contains(a.Decoders, name) {
				a.Decoders = append(a.Decoders, name)
			}
		}
	}
	slices.Sort(a.Outputs)
	slices.Sort(a.Decoders)
	return a
}

// Packages are the Linux packages the game needs (wiki First start).
const Packages = "Debian, Ubuntu: libwebkit2gtk-4.1-0 gstreamer1.0-plugins-good gstreamer1.0-libav; " +
	"Fedora: webkit2gtk4.1 gstreamer1-plugins-good gstreamer1-plugin-libav; " +
	"Arch: webkit2gtk-4.1 gst-plugins-good gst-libav"
