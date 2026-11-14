package update

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Apply starts a helper that waits for this app to quit, swaps the
// install for the downloaded one and starts it. The caller quits next.
func (c *Client) Apply() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := c.identity(); err != nil {
		return err
	}
	if c.payload == "" || c.workDir == "" {
		return ErrNoDownload
	}
	dest, err := c.installDest()
	if err != nil {
		return err
	}
	if !dirWritable(filepath.Dir(dest)) {
		return fmt.Errorf("update: cannot write %s", filepath.Dir(dest))
	}
	script, err := c.writeHelper(dest, os.Getpid())
	if err != nil {
		return err
	}
	start := c.Start
	if start == nil {
		start = startDetached
	}
	if c.cfg.GOOS == "windows" {
		return start("powershell", "-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-File", script)
	}
	return start("/bin/sh", script)
}

// installDest is what the update replaces: the .app on macOS, else the
// executable itself.
func (c *Client) installDest() (string, error) {
	exe := c.Executable
	if exe == "" {
		var err error
		if exe, err = os.Executable(); err != nil {
			return "", fmt.Errorf("update: cannot find this app: %w", err)
		}
	}
	exe, err := filepath.EvalSymlinks(exe)
	if err != nil {
		return "", fmt.Errorf("update: cannot find this app: %w", err)
	}
	if c.cfg.GOOS == "darwin" {
		marker := ".app" + string(filepath.Separator) + "Contents" + string(filepath.Separator) + "MacOS" + string(filepath.Separator)
		if i := strings.LastIndex(exe, marker); i > 0 {
			return exe[:i+len(".app")], nil
		}
	}
	return exe, nil
}

func (c *Client) installWritable() bool {
	dest, err := c.installDest()
	return err == nil && dirWritable(filepath.Dir(dest))
}

func dirWritable(dir string) bool {
	if dir == "" || dir == "." {
		return false
	}
	f, err := os.CreateTemp(dir, ".write-check-*")
	if err != nil {
		return false
	}
	name := f.Name()
	return errors.Join(f.Close(), os.Remove(name)) == nil
}

func (c *Client) writeHelper(dest string, pid int) (string, error) {
	old := dest + "." + c.cfg.Name + "-old"
	var p string
	var body []byte
	if c.cfg.GOOS == "windows" {
		p = filepath.Join(c.workDir, "apply.ps1")
		body = windowsHelper(dest, c.payload, old, pid)
	} else {
		p = filepath.Join(c.workDir, "apply.sh")
		body = unixHelper(dest, c.payload, old, pid, c.cfg.GOOS == "darwin" && strings.HasSuffix(dest, ".app"))
	}
	if err := os.WriteFile(p, body, 0o700); err != nil { //nolint:gosec // the helper must run
		return "", fmt.Errorf("update: %w", err)
	}
	return p, nil
}

// unixHelper waits up to 30 s for the app to quit, swaps the install,
// puts the old one back if the move fails, and starts the new one.
func unixHelper(dest, payload, old string, pid int, mac bool) []byte {
	run := `nohup "$DEST" >/dev/null 2>&1 &`
	if mac {
		run = `open "$DEST"`
	}
	return []byte(strings.Join([]string{
		"#!/bin/sh",
		"DEST=" + shQuote(dest),
		"NEW=" + shQuote(payload),
		"OLD=" + shQuote(old),
		"PID=" + strconv.Itoa(pid),
		"left=60",
		`while kill -0 "$PID" 2>/dev/null && [ "$left" -gt 0 ]; do sleep 0.5; left=$((left - 1)); done`,
		`[ -e "$NEW" ] || exit 1`,
		`rm -rf "$OLD"`,
		`if [ -e "$DEST" ]; then mv "$DEST" "$OLD" || exit 1; fi`,
		`if ! mv "$NEW" "$DEST"; then mv "$OLD" "$DEST"; exit 1; fi`,
		`chmod u+wx "$DEST" 2>/dev/null || true`,
		run,
		`rm -rf "$OLD"`,
		"",
	}, "\n"))
}

// shQuote quotes s for sh: single quotes, each ' written as '\”.
func shQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

// windowsHelper is UTF-8 with a BOM so PowerShell keeps Unicode paths.
// It retries while Windows still holds the old files.
func windowsHelper(dest, payload, old string, pid int) []byte {
	body := strings.Join([]string{
		"$ErrorActionPreference = 'Continue'",
		"$dest = " + psQuote(dest),
		"$new = " + psQuote(payload),
		"$old = " + psQuote(old),
		"Wait-Process -Id " + strconv.Itoa(pid) + " -Timeout 30 -ErrorAction SilentlyContinue",
		"if (-not (Test-Path -LiteralPath $new)) { return }",
		"$ok = $false",
		"$left = 30",
		"while ($left -gt 0) {",
		"  $left--",
		"  try {",
		"    if (Test-Path -LiteralPath $old) { Remove-Item -LiteralPath $old -Recurse -Force -ErrorAction Stop }",
		"    if (Test-Path -LiteralPath $dest) { Move-Item -LiteralPath $dest -Destination $old -Force -ErrorAction Stop }",
		"    Move-Item -LiteralPath $new -Destination $dest -Force -ErrorAction Stop",
		"    $ok = $true",
		"    break",
		"  } catch {",
		"    if ((Test-Path -LiteralPath $old) -and -not (Test-Path -LiteralPath $dest)) {",
		"      try { Move-Item -LiteralPath $old -Destination $dest -Force -ErrorAction Stop } catch {}",
		"    }",
		"    Start-Sleep -Seconds 1",
		"  }",
		"}",
		"if (-not $ok) { return }",
		"Start-Process -FilePath $dest -WorkingDirectory (Split-Path -LiteralPath $dest)",
		"if (Test-Path -LiteralPath $old) {",
		"  try { Remove-Item -LiteralPath $old -Recurse -Force -ErrorAction Stop } catch {}",
		"}",
		"",
	}, "\r\n")
	out := make([]byte, 0, 3+len(body))
	out = append(out, 0xEF, 0xBB, 0xBF)
	return append(out, body...)
}

// psQuote quotes s for PowerShell: single quotes, each quote doubled.
// PowerShell also ends a string at the curly quotes.
func psQuote(s string) string {
	var b strings.Builder
	b.WriteByte('\'')
	for _, r := range s {
		b.WriteRune(r)
		switch r {
		case '\'', '\u2018', '\u2019', '\u201a', '\u201b':
			b.WriteRune(r)
		}
	}
	b.WriteByte('\'')
	return b.String()
}
