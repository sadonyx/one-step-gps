package tokens

import "time"

const (
	// The number of characters in the token string.
	TokenStringLength int = 32
	// The number of minutes between clean-up cycles.
	TokenCleanupTimer time.Duration = 1 * time.Minute
)
