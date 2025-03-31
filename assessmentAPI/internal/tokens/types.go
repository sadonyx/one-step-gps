package tokens

import (
	"sync"
	"time"
)

// Represents a short-lived authentication token
type Token struct {
	SessionID string
	ExpiresAt time.Time
}

// Handles the creation, storage, and expiration of tokens
type TokenManager struct {
	tokens map[string]*Token
	mu     sync.RWMutex
}
