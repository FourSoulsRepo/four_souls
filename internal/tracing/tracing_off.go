//go:build !trace

package tracing

import "context"

// start refuses: this build has no recording code.
func start(context.Context, Options) (func() error, error) {
	return nil, ErrNotBuilt
}
