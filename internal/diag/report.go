package diag

import (
	"fmt"
	"io"
	"slices"
	"strings"
	"text/tabwriter"
)

// lowFPS is below what the drifting sky needs to look smooth.
const lowFPS = 40

// softwareRenderers mark a webview that draws without the GPU.
var softwareRenderers = []string{"llvmpipe", "softpipe", "swiftshader", "software"}

// Write prints the report (ADR 035, item 4): the system, one row per
// probe, the best variant, and what looks wrong.
func Write(w io.Writer, s System, outs []Outcome) {
	p := &printer{w: w}
	p.line("Four Souls diagnostics")
	p.line("")
	p.line("App:        %s (%s/%s)", s.App, s.OS, s.Arch)
	if s.OS == "linux" {
		p.line("WebKitGTK:  %s", orUnknown(s.WebKit))
		p.line("Session:    %s, desktop %s", orUnknown(s.Session), orUnknown(s.Desktop))
		p.line("Display:    DISPLAY=%s WAYLAND_DISPLAY=%s", orNone(s.Display), orNone(s.Wayland))
		writeGPUs(p, s)
	}
	p.line("Switches:   %s", orNone(strings.Join(s.Switches, " ")))
	p.line("")
	writeProbes(p, s.OS, outs)
	p.line("")
	writeVerdict(p, s, outs)
}

func writeGPUs(p *printer, s System) {
	if len(s.GPUs) == 0 {
		p.line("GPU:        no render node in /dev/dri")
	}
	for _, g := range s.GPUs {
		access := "opens"
		if !g.Opens {
			access = "does not open: " + g.Err
		}
		p.line("GPU:        %s, driver %s, %s, %s", g.Node, orUnknown(g.Driver), orUnknown(g.Vendor), access)
	}
	if s.GroupsKnown {
		p.line("Groups:     %s", orNone(strings.Join(s.Groups, ", ")))
	}
}

func writeProbes(p *printer, goos string, outs []Outcome) {
	p.line("Main menu, %dx%d window, frames over 5 s:", ProbeWidth, ProbeHeight)
	tw := tabwriter.NewWriter(p.w, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprintln(tw, "variant\tfps\tslow\tworst\tscale\tsky\trenderer") //nolint:errcheck // see printer
	for _, o := range outs {
		if o.Result == nil {
			_, _ = fmt.Fprintf(tw, "%s\t-\t-\t-\t-\t-\tfailed: %s\n", o.Variant.ID, o.Err) //nolint:errcheck // see printer
			continue
		}
		r := o.Result
		slow := 0
		if r.Frames > 0 {
			slow = r.Slow * 100 / r.Frames
		}
		_, _ = fmt.Fprintf(tw, "%s\t%.1f\t%d%%\t%.0f ms\t%.2g\t%s\t%s\n", //nolint:errcheck // see printer
			o.Variant.ID, r.FPS(), slow, r.WorstMs, r.PixelRatio, sky(r), renderer(goos, r))
	}
	_ = tw.Flush() //nolint:errcheck // see printer
}

func writeVerdict(p *printer, s System, outs []Outcome) {
	best, ok := Best(outs)
	if ok {
		p.line("Best: %s, %.1f fps.", best.Variant.ID, best.Result.FPS())
		if hint := tryHint(best.Variant); hint != "" {
			p.line("Try:  %s", hint)
		}
	} else {
		p.line("Best: none; no probe sent a result.")
	}
	flags := Flags(s, outs)
	if len(flags) == 0 {
		return
	}
	p.line("")
	p.line("Problems found:")
	for _, f := range flags {
		p.line("- %s", f)
	}
}

// Best is the probe with the most frames per second.
func Best(outs []Outcome) (Outcome, bool) {
	var best Outcome
	found := false
	for _, o := range outs {
		if o.Result != nil && (!found || o.Result.FPS() > best.Result.FPS()) {
			best, found = o, true
		}
	}
	return best, found
}

// tryHint is how to start the game with a variant today; "" when it is
// the default or needs a setting the game does not have yet (25.3).
func tryHint(v Variant) string {
	switch {
	case v.ID == "default":
		return ""
	case len(v.Env) > 0:
		return strings.Join(v.Env, " ") + " ./FourSouls"
	}
	return ""
}

// Masked tells that the webview hides the GPU from pages (19.228): WebKit
// (Linux, macOS) answers "Apple GPU" on every computer.
func Masked(goos string) bool {
	return goos == "linux" || goos == "darwin"
}

// renderer is a probe's GPU as the report shows it.
func renderer(goos string, r *Result) string {
	if Masked(goos) {
		return "hidden by WebKit"
	}
	return orUnknown(r.Renderer)
}

// Flags lists what looks wrong (ADR 035, item 4.4).
func Flags(s System, outs []Outcome) []string {
	var flags []string
	if s.OS == "linux" {
		if len(s.GPUs) == 0 {
			flags = append(flags, "No GPU render node in /dev/dri: the webview draws on the CPU.")
		}
		if slices.Contains(s.Switches, "LIBGL_ALWAYS_SOFTWARE=1") {
			flags = append(flags, "LIBGL_ALWAYS_SOFTWARE=1 is set: OpenGL draws on the CPU.")
		}
		for _, g := range s.GPUs {
			if !g.Opens {
				flags = append(flags, fmt.Sprintf("%s does not open (%s): add the user to the render group and log in again.", g.Node, g.Err))
			}
		}
	}
	if d, ok := defaultResult(outs); ok {
		flags = append(flags, resultFlags(s.OS, d)...)
	}
	if best, ok := Best(outs); ok && best.Result.FPS() < lowFPS {
		flags = append(flags, fmt.Sprintf("Every variant is slow: under %d fps.", lowFPS))
	}
	return flags
}

// resultFlags are the problems the default start shows. A masked
// renderer says nothing (19.228); the GPU lines and frame rates do.
func resultFlags(goos string, r *Result) []string {
	var flags []string
	renderer := strings.ToLower(r.Renderer)
	if !Masked(goos) && slices.ContainsFunc(softwareRenderers, func(s string) bool { return strings.Contains(renderer, s) }) {
		flags = append(flags, fmt.Sprintf("The webview draws on the CPU (%s): check the GPU driver (Mesa).", r.Renderer))
	}
	if !slices.Contains(r.Sounds, "mp3") {
		flags = append(flags, "No MP3 playback: install the GStreamer plugins (good, and libav).")
	}
	if r.ReducedMotion {
		flags = append(flags, "The system asks for reduced motion: the sky stands still, so the frame rate means less.")
	}
	return flags
}

// defaultResult is the default start's result, or else the first one.
func defaultResult(outs []Outcome) (*Result, bool) {
	for _, o := range outs {
		if o.Result != nil && o.Variant.ID == "default" {
			return o.Result, true
		}
	}
	for _, o := range outs {
		if o.Result != nil {
			return o.Result, true
		}
	}
	return nil, false
}

func sky(r *Result) string {
	if r.Sky {
		return "moving"
	}
	return "off"
}

func orUnknown(s string) string {
	if s == "" {
		return "unknown"
	}
	return s
}

func orNone(s string) string {
	if s == "" {
		return "none"
	}
	return s
}

// printer writes lines; the report goes to a terminal and a file the
// caller checks, so a single write error is not worth stopping for.
type printer struct{ w io.Writer }

func (p *printer) line(format string, args ...any) {
	_, _ = fmt.Fprintf(p.w, format+"\n", args...) //nolint:errcheck // see printer
}
