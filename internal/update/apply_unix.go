//go:build !windows

package update

import (
	"context"
	"os/exec"
	"syscall"
)

// startDetached runs the helper in its own session, so it outlives the app.
// Its context is never canceled: the helper must outlive the app.
func startDetached(name string, args ...string) error {
	cmd := exec.CommandContext(context.Background(), name, args...) //nolint:gosec // our own helper script
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	return cmd.Start()
}
