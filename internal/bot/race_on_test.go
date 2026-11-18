//go:build race

package bot

// raceSlow is how much slower the race detector makes the code: timing
// bounds of the tests stretch by it (TestHardYields).
const raceSlow = 4
