//go:build !trace

package tracing

import (
	"context"
	"errors"
	"testing"
)

// Without the tag, any request for tracing is refused, not ignored.
func TestStartRefusesWithoutTag(t *testing.T) {
	for _, o := range []Options{{Dir: "traces"}, {Addr: "127.0.0.1:6060"}} {
		if _, err := Start(context.Background(), o); !errors.Is(err, ErrNotBuilt) {
			t.Errorf("Start(%+v) = %v, want ErrNotBuilt", o, err)
		}
	}
}
