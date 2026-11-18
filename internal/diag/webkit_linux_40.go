//go:build linux && cgo && !webkit2_41

package diag

/*
#cgo pkg-config: webkit2gtk-4.0
#include <webkit2/webkit2.h>
*/
import "C"

import "fmt"

// webkitVersion asks the linked WebKitGTK 4.0 for its version; the user
// agent does not tell it (ADR 035, item 2.3).
func webkitVersion() string {
	return fmt.Sprintf("%d.%d.%d (webkit2gtk-4.0)",
		uint(C.webkit_get_major_version()), uint(C.webkit_get_minor_version()), uint(C.webkit_get_micro_version()))
}
