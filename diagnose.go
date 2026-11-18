package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options/linux"

	"github.com/FourSoulsRepo/four_souls/internal/diag"
	"github.com/FourSoulsRepo/four_souls/internal/version"
)

// runDiagnose is `--diagnose` (ADR 035): the system, then one probe
// window per variant, then the report on the terminal and in a file.
// It returns the exit code.
func runDiagnose() int {
	sys := diag.ReadSystem(version.App)
	exe, err := os.Executable()
	if err != nil {
		log.Print(err)
		return 1
	}
	variants := diag.Variants(runtime.GOOS)
	fmt.Printf("Four Souls diagnostics: test windows open one after another (%d in all); please wait.\n", len(variants))
	outs := diag.ProbeAll(variants, func(v diag.Variant) ([]byte, error) {
		ctx, cancel := context.WithTimeout(context.Background(), diag.ProbeTimeout)
		defer cancel()
		return diag.Command(ctx, exe, v).Output()
	}, func(o diag.Outcome) {
		if o.Result != nil {
			fmt.Printf("  %s: %.1f fps\n", o.Variant.ID, o.Result.FPS())
		} else {
			fmt.Printf("  %s: failed: %s\n", o.Variant.ID, o.Err)
		}
	})
	var report bytes.Buffer
	diag.Write(&report, sys, outs)
	fmt.Println()
	fmt.Print(report.String())
	// Players often start the app without a terminal (ADR 035, item 4).
	path := appPath("diagnose-" + time.Now().Format("20060102-150405") + ".txt")
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		log.Print(err)
		return 1
	}
	if err := os.WriteFile(path, report.Bytes(), 0o600); err != nil {
		log.Print(err)
		return 1
	}
	fmt.Println()
	fmt.Println("Saved:", path)
	return 0
}

// runProbe is one window of `--diagnose`: the main menu, measured by the
// page, which reports through the probe binding and quits. It holds no
// instance slot and only reads the first copy's settings. Its switches
// come in the environment, as a player would set them.
func runProbe(id string) {
	v, ok := diag.VariantByID(runtime.GOOS, id)
	if !ok {
		log.Fatal("unknown diagnose variant; see --diagnose")
	}
	title := "Four Souls diagnostics: " + v.ID
	if err := wails.Run(appOptions(title, appPath("settings.json"), false, v)); err != nil {
		log.Fatal(err)
	}
}

// linuxOptions sets the webview's GPU policy (ADR 036): always on,
// unless FOURSOULS_GPU_POLICY asks for "ondemand" or "never". With no
// options Wails v2 picks Never, and WebKitGTK drew the menu on the CPU
// at 6 fps (25.2). Other systems ignore it.
func linuxOptions(value string) *linux.Options {
	policy, ok := diag.GPUPolicy(value)
	if !ok {
		log.Print(diag.GPUPolicyEnv + ": unknown value; using " + diag.GPUAlways)
	}
	switch policy {
	case diag.GPUOnDemand:
		return &linux.Options{WebviewGpuPolicy: linux.WebviewGpuPolicyOnDemand}
	case diag.GPUNever:
		return &linux.Options{WebviewGpuPolicy: linux.WebviewGpuPolicyNever}
	}
	return &linux.Options{WebviewGpuPolicy: linux.WebviewGpuPolicyAlways}
}
