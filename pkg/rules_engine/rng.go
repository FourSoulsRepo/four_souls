package rulesengine

// RNG is a small seeded random generator (PCG32, XSH-RR variant).
//
// It lives inside the game state, so a saved game continues with the same
// numbers (A-08). Its fields are never sent to clients (R-ZONE-04, A-07).
type RNG struct {
	State uint64 `json:"state"`
	Inc   uint64 `json:"inc"`
}

const pcgMultiplier = 6364136223846793005

// NewRNG returns a generator for seed. Equal seeds give equal sequences.
func NewRNG(seed uint64) RNG {
	r := RNG{Inc: seed<<1 | 1}
	r.next()
	r.State += seed
	r.next()
	return r
}

// next returns the next 32 random bits.
func (r *RNG) next() uint32 {
	old := r.State
	r.State = old*pcgMultiplier + r.Inc
	xorShifted := uint32(((old >> 18) ^ old) >> 27) //nolint:gosec // PCG keeps the low 32 bits on purpose
	rot := uint32(old >> 59)                        //nolint:gosec // top 5 bits, always fits
	return xorShifted>>rot | xorShifted<<((-rot)&31)
}

// Intn returns a uniform number in [0, n). It panics if n <= 0.
func (r *RNG) Intn(n int) int {
	if n <= 0 {
		panic("rulesengine: Intn with n <= 0")
	}
	bound := uint32(n) //nolint:gosec // game sizes are far below 2^32
	// Lemire's method without bias: reject the low values that wrap.
	threshold := -bound % bound
	for {
		x := r.next()
		m := uint64(x) * uint64(bound)
		if uint32(m) >= threshold { //nolint:gosec // low 32 bits by design
			return int(m >> 32)
		}
	}
}

// Shuffle puts n elements in random order (Fisher–Yates).
func (r *RNG) Shuffle(n int, swap func(i, j int)) {
	for i := n - 1; i > 0; i-- {
		swap(i, r.Intn(i+1))
	}
}

// D6 rolls a six-sided die (R-DICE-01).
func (r *RNG) D6() int {
	return r.Intn(6) + 1
}
