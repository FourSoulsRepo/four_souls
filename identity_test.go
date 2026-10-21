package main

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// bundleID is the app's identity on macOS and Windows; it never changes
// between releases (B-03, 8.3).
const bundleID = "io.github.foursoulsrepo.foursouls"

func TestStableIdentity(t *testing.T) {
	var cfg struct {
		Name           string `json:"name"`
		OutputFilename string `json:"outputfilename"`
		Info           struct {
			CompanyName string `json:"companyName"`
			ProductName string `json:"productName"`
		} `json:"info"`
	}
	data, err := os.ReadFile("wails.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.Name != "FourSouls" || cfg.Info.ProductName != "Four Souls" || cfg.Info.CompanyName != "FourSoulsRepo" {
		t.Errorf("wails.json identity changed: %+v", cfg)
	}
	// Release archives hold FourSouls.app (from name), FourSouls.exe, or a
	// bare FourSouls (8.27).
	if cfg.OutputFilename != "FourSouls" {
		t.Errorf("wails.json outputfilename = %q, want FourSouls", cfg.OutputFilename)
	}
	for _, f := range []string{"build/darwin/Info.plist", "build/windows/wails.exe.manifest"} {
		data, err := os.ReadFile(f) //nolint:gosec // fixed repo files
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(data), bundleID) {
			t.Errorf("%s does not name %s", f, bundleID)
		}
	}
}

// TestVersionFields checks that every version field takes the stamped
// release version (19.47); Wails fills the template from wails.json.
func TestVersionFields(t *testing.T) {
	const want = "{{.Info.ProductVersion}}"
	plist, err := os.ReadFile("build/darwin/Info.plist")
	if err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"CFBundleVersion", "CFBundleShortVersionString"} {
		if !strings.Contains(string(plist), "<key>"+key+"</key>\n        <string>"+want+"</string>") {
			t.Errorf("Info.plist %s is not %s", key, want)
		}
	}
	var info struct {
		Fixed map[string]string            `json:"fixed"`
		Info  map[string]map[string]string `json:"info"`
	}
	data, err := os.ReadFile("build/windows/info.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, &info); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"file_version", "product_version"} {
		if info.Fixed[key] != want {
			t.Errorf("info.json fixed %s = %q, want %s", key, info.Fixed[key], want)
		}
	}
	for _, key := range []string{"FileVersion", "ProductVersion"} {
		if info.Info["0000"][key] != want {
			t.Errorf("info.json %s = %q, want %s", key, info.Info["0000"][key], want)
		}
	}
}
