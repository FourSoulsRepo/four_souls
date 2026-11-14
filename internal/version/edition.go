package version

// The two release editions (ADR 029). Edition is one of them, set by
// build tags in edition_full.go and edition_for_poor.go.
const (
	// EditionFull has the card images: tags embed and cardimages.
	EditionFull = "full"
	// EditionForPoor shows text cards: any other build.
	EditionForPoor = "for-poor"
)
