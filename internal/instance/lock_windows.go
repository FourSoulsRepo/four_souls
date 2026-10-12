//go:build windows

package instance

import (
	"os"

	"golang.org/x/sys/windows"
)

// lock takes an exclusive lock on f, or fails at once if another open
// file holds it.
func lock(f *os.File) error {
	flags := uint32(windows.LOCKFILE_EXCLUSIVE_LOCK | windows.LOCKFILE_FAIL_IMMEDIATELY)
	return windows.LockFileEx(windows.Handle(f.Fd()), flags, 0, 1, 0, new(windows.Overlapped)) //nolint:wrapcheck // only busy or not matters
}
