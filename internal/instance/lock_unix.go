//go:build !windows

package instance

import (
	"os"

	"golang.org/x/sys/unix"
)

// lock takes an exclusive lock on f, or fails at once if another open
// file holds it.
func lock(f *os.File) error {
	return unix.Flock(int(f.Fd()), unix.LOCK_EX|unix.LOCK_NB) //nolint:wrapcheck,gosec // only busy or not matters
}
