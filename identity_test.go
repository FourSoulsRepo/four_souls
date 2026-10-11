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
		Name string `json:"name"`
		Info struct {
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
	if cfg.Name != "four_souls" || cfg.Info.ProductName != "Four Souls" || cfg.Info.CompanyName != "FourSoulsRepo" {
		t.Errorf("wails.json identity changed: %+v", cfg)
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
