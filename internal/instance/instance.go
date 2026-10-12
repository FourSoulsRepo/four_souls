// Package instance lets several copies of the app run on one computer.
// Each running copy holds a numbered slot; slot 1 is the first copy.
package instance

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// MaxSlots is how many copies can hold their own slot at once.
const MaxSlots = 16

// ErrBusy means every slot is held by a running copy.
var ErrBusy = errors.New("instance: every slot is in use")

// Slot is held until Release or until the process ends; the OS drops
// the lock of a crashed copy.
type Slot struct {
	N   int // 1 for the first copy
	dir string
	f   *os.File
}

// Acquire takes the lowest free slot in dir.
func Acquire(dir string) (*Slot, error) {
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return nil, fmt.Errorf("instance: %w", err)
	}
	for n := 1; n <= MaxSlots; n++ {
		path := filepath.Join(dir, "instance-"+strconv.Itoa(n)+".lock")
		f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o600) //nolint:gosec // the app names the file
		if err != nil {
			return nil, fmt.Errorf("instance: %w", err)
		}
		if err := lock(f); err != nil {
			_ = f.Close() //nolint:errcheck // the slot is busy; nothing was written
			continue
		}
		return &Slot{N: n, dir: dir, f: f}, nil
	}
	return nil, ErrBusy
}

// File names a per-copy file in the slot's folder: name itself for
// slot 1, so the first copy keeps its old files, else name-N.ext.
func (s *Slot) File(name string) string {
	if s.N > 1 {
		ext := filepath.Ext(name)
		name = name[:len(name)-len(ext)] + "-" + strconv.Itoa(s.N) + ext
	}
	return filepath.Join(s.dir, name)
}

// Release frees the slot for the next copy.
func (s *Slot) Release() error {
	return s.f.Close() //nolint:wrapcheck // closing drops the lock
}
