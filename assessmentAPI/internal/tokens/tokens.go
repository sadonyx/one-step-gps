package tokens

import (
	"time"

	"crypto/rand"
	"encoding/base64"
)

// Creates a new instance of TokenManager.
func NewTokenManager() *TokenManager {
	tm := &TokenManager{
		tokens: make(map[string]*Token),
	}

	go tm.startCleanupRoutine()
	return tm
}

// Gets the value of a token using a token key.
func (tm *TokenManager) Get(tokenKey string) *Token {
	return tm.tokens[tokenKey]
}

// Creates a cryptographically secure random token.
func generateTokenString() string {
	b := make([]byte, TokenStringLength)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// Generates a new token with a specified duration.
func (tm *TokenManager) Create(sessionID string, duration time.Duration) string {
    if tm.tokens == nil {
        tm.tokens = make(map[string]*Token)
    }
    
    tm.mu.Lock()
    defer tm.mu.Unlock()
    
    tokenKey := generateTokenString()
    token := &Token{
        SessionID: sessionID,
        ExpiresAt: time.Now().Add(duration),
    }
    
    tm.tokens[tokenKey] = token
    return tokenKey
}

// Checks if a token is valid and not expired.
func (tm *TokenManager) Validate(tokenKey string) bool {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	token, exists := tm.tokens[tokenKey]
	return exists && time.Now().Before(token.ExpiresAt)
}

// Periodically removes expired tokens.
func (tm *TokenManager) startCleanupRoutine() {
	ticker := time.NewTicker(TokenCleanupTimer)
	defer ticker.Stop()

	for range ticker.C {
		tm.cleanupExpiredTokens()
	}
}

// Removes all expired tokens.
func (tm *TokenManager) cleanupExpiredTokens() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	now := time.Now()
	for tokenKey, token := range tm.tokens {
		if now.After(token.ExpiresAt) {
			delete(tm.tokens, tokenKey)
		}
	}
}
