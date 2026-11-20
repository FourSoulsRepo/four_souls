// Package diag finds why the webview draws slowly on a computer (ADR 035):
// a report of the system, and the frame rate of the real window measured
// once per known WebKitGTK switch.
package diag

import (
	"encoding/json"
	"strings"
)

// The command-line names: Flag runs the whole report; ProbeName starts
// the app as one probe, `--diagnose-probe=<variant>`.
const (
	Flag      = "diagnose"
	ProbeName = "diagnose-probe"
)

// The probe window's size: the same on every computer, so frame rates
// compare (ADR 035, item 3.3).
const (
	ProbeWidth  = 1280
	ProbeHeight = 800
)

// resultPrefix starts the one stdout line a probe sends its parent;
// other lines (Wails, WebKitGTK) are skipped.
const resultPrefix = "four-souls-probe: "

// Args reads the command line: whether it asks for the report, and the
// variant when the app runs as a probe. One or two dashes both work.
func Args(args []string) (diagnose bool, probe string) {
	for _, a := range args {
		if !strings.HasPrefix(a, "-") {
			continue
		}
		name, value, _ := strings.Cut(strings.TrimLeft(a, "-"), "=")
		switch name {
		case Flag:
			diagnose = true
		case ProbeName:
			probe = value
		}
	}
	return diagnose, probe
}

// GPUPolicyEnv sets the Linux webview's GPU policy (ADR 036): "always"
// (the default), "ondemand" or "never".
const GPUPolicyEnv = "FOURSOULS_GPU_POLICY"

// The GPU policies GPUPolicyEnv takes.
const (
	GPUAlways   = "always"
	GPUOnDemand = "ondemand"
	GPUNever    = "never"
)

// GPUPolicy reads a GPUPolicyEnv value: unset is "always"; an unknown
// value is "always" too, and ok is false so the app can say so.
func GPUPolicy(value string) (policy string, ok bool) {
	switch v := strings.ToLower(strings.TrimSpace(value)); v {
	case "":
		return GPUAlways, true
	case GPUAlways, GPUOnDemand, GPUNever:
		return v, true
	}
	return GPUAlways, false
}

// Variant is one way to start the webview: switches added to the
// environment, so the report can print how to start the game with it.
type Variant struct {
	ID string
	// Env holds NAME=value pairs added to the probe's environment.
	Env []string
}

// Variants lists what the report tries on this system (ADR 035, item
// 3.1): the switches only mean something to WebKitGTK. The default is
// the GPU always on (ADR 036); gpu-never is how every start ran before.
func Variants(goos string) []Variant {
	if goos != "linux" {
		return []Variant{{ID: "default"}}
	}
	return []Variant{
		{ID: "default"},
		{ID: "no-dmabuf", Env: []string{"WEBKIT_DISABLE_DMABUF_RENDERER=1"}},
		{ID: "no-compositing", Env: []string{"WEBKIT_DISABLE_COMPOSITING_MODE=1"}},
		{ID: "gpu-ondemand", Env: []string{GPUPolicyEnv + "=" + GPUOnDemand}},
		{ID: "gpu-never", Env: []string{GPUPolicyEnv + "=" + GPUNever}},
	}
}

// VariantByID finds a variant of this system by its ID.
func VariantByID(goos, id string) (Variant, bool) {
	for _, v := range Variants(goos) {
		if v.ID == id {
			return v, true
		}
	}
	return Variant{}, false
}

// Result is what a probe's window measured and saw (ADR 035, item 3.7).
type Result struct {
	Variant string `json:"variant"`
	// Frames drawn in Seconds; Slow of them took over SlowMs.
	Frames  int     `json:"frames"`
	Seconds float64 `json:"seconds"`
	Slow    int     `json:"slow"`
	WorstMs float64 `json:"worstMs"`
	// The WebGL renderer and vendor; "llvmpipe" means no GPU. WebKit
	// hides both (Masked).
	Renderer string `json:"renderer"`
	Vendor   string `json:"vendor"`
	// The window as the page sees it.
	PixelRatio float64 `json:"pixelRatio"`
	Width      int     `json:"width"`
	Height     int     `json:"height"`
	// Sky: the drifting sky was drawn and moving; ReducedMotion: the
	// system asks for no motion, which stops it.
	Sky           bool `json:"sky"`
	ReducedMotion bool `json:"reducedMotion"`
	// Sounds lists the sound types the webview can play.
	Sounds []string `json:"sounds"`
}

// SlowMs is a frame that missed two refreshes of a 60 Hz screen.
const SlowMs = 34

// FPS is the frames per second; 0 when nothing was measured.
func (r Result) FPS() float64 {
	if r.Seconds <= 0 {
		return 0
	}
	return float64(r.Frames) / r.Seconds
}

// ResultLine is the line a probe prints for its parent.
func ResultLine(r Result) string {
	b, err := json.Marshal(r)
	if err != nil {
		return resultPrefix + "{}"
	}
	return resultPrefix + string(b)
}

// parseResult finds a probe's result in its output.
func parseResult(out []byte) (Result, bool) {
	for line := range strings.Lines(string(out)) {
		rest, ok := strings.CutPrefix(strings.TrimSpace(line), resultPrefix)
		if !ok {
			continue
		}
		var r Result
		if json.Unmarshal([]byte(rest), &r) == nil {
			return r, true
		}
	}
	return Result{}, false
}
