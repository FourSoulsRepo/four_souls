package diag

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
)

// System is what the report knows before any window opens (ADR 035,
// item 2). The Linux fields stay empty on other systems.
type System struct {
	OS, Arch, App string
	// WebKit is the WebKitGTK version; "" when unknown.
	WebKit string
	// The desktop session: XDG_SESSION_TYPE, XDG_CURRENT_DESKTOP,
	// DISPLAY, WAYLAND_DISPLAY.
	Session, Desktop, Display, Wayland string
	GPUs                               []GPU
	// Groups lists which of the GPU groups the user is in;
	// GroupsKnown is false when they could not be read.
	Groups      []string
	GroupsKnown bool
	// Switches are the webview switches already set in the environment.
	Switches []string
}

// GPU is one render node: the device the webview draws with.
type GPU struct {
	Node   string
	Driver string
	Vendor string
	// Opens: this user may use it; Err says why not.
	Opens bool
	Err   string
}

// gpuGroups are the groups that give access to the GPU on most distros.
var gpuGroups = []string{"render", "video"}

// switchNames are the environment variables that change how the
// webview draws.
var switchNames = []string{
	"WEBKIT_DISABLE_DMABUF_RENDERER",
	"WEBKIT_DISABLE_COMPOSITING_MODE",
	"LIBGL_ALWAYS_SOFTWARE",
	"GDK_BACKEND",
	"GSK_RENDERER",
}

// vendors names the PCI vendor IDs of the usual GPUs.
var vendors = map[string]string{
	"0x8086": "Intel",
	"0x1002": "AMD",
	"0x10de": "NVIDIA",
}

// ReadSystem reads this computer; app is the app's version.
func ReadSystem(app string) System {
	s := System{OS: runtime.GOOS, Arch: runtime.GOARCH, App: app}
	for _, name := range switchNames {
		if v, ok := os.LookupEnv(name); ok {
			s.Switches = append(s.Switches, name+"="+v)
		}
	}
	if s.OS != "linux" {
		return s
	}
	s.WebKit = webkitVersion()
	s.Session = os.Getenv("XDG_SESSION_TYPE")
	s.Desktop = os.Getenv("XDG_CURRENT_DESKTOP")
	s.Display = os.Getenv("DISPLAY")
	s.Wayland = os.Getenv("WAYLAND_DISPLAY")
	s.GPUs = gpus("/sys/class/drm", "/dev/dri")
	s.Groups, s.GroupsKnown = userGroups()
	return s
}

// gpus lists the render nodes under dev, with their driver and vendor
// from sys (the kernel's /sys/class/drm).
func gpus(sys, dev string) []GPU {
	nodes, _ := filepath.Glob(filepath.Join(dev, "renderD*")) //nolint:errcheck // only a bad pattern fails
	list := make([]GPU, 0, len(nodes))
	for _, path := range nodes {
		name := filepath.Base(path)
		g := GPU{Node: path}
		if target, err := os.Readlink(filepath.Join(sys, name, "device", "driver")); err == nil {
			g.Driver = filepath.Base(target)
		}
		if b, err := os.ReadFile(filepath.Join(sys, name, "device", "vendor")); err == nil { //nolint:gosec // the kernel's /sys folder
			id := strings.TrimSpace(string(b))
			g.Vendor = vendors[id]
			if g.Vendor == "" {
				g.Vendor = id
			}
		}
		f, err := os.OpenFile(path, os.O_RDWR, 0) //nolint:gosec // a fixed device folder
		if err != nil {
			g.Err = err.Error()
		} else {
			g.Opens = true
			_ = f.Close() //nolint:errcheck // only opened to test access
		}
		list = append(list, g)
	}
	return list
}

// userGroups tells which GPU groups the user is in.
func userGroups() ([]string, bool) {
	u, err := user.Current()
	if err != nil {
		return nil, false
	}
	ids, err := u.GroupIds()
	if err != nil {
		return nil, false
	}
	var in []string
	for _, id := range ids {
		g, err := user.LookupGroupId(id)
		if err == nil && slices.Contains(gpuGroups, g.Name) {
			in = append(in, g.Name)
		}
	}
	slices.Sort(in)
	return in, true
}
