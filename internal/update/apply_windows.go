//go:build windows

package update

import (
	"context"
	"os/exec"
	"syscall"
)

const (
	createNoWindow         = 0x08000000
	createNewProcessGroup  = 0x00000200
	createBreakawayFromJob = 0x01000000
)

// startDetached runs the helper without a console. Breakaway keeps it
// alive after the app quits; some job objects forbid it, so it retries
// without.
func startDetached(name string, args ...string) error {
	err := startHidden(name, args, createNoWindow|createNewProcessGroup|createBreakawayFromJob)
	if err != nil {
		return startHidden(name, args, createNoWindow|createNewProcessGroup)
	}
	return nil
}

func startHidden(name string, args []string, flags uint32) error {
	cmd := exec.CommandContext(context.Background(), name, args...) //nolint:gosec // our own helper script
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true, CreationFlags: flags}
	return cmd.Start()
}
