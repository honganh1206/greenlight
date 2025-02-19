package data

import (
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name      string
		userID    int64
		ttl       time.Duration
		scope     string
		wantError bool
	}{
		{
			name:      "Valid token generation",
			userID:    1,
			ttl:       24 * time.Hour,
			scope:     ScopeActivation,
			wantError: false,
		},
		{
			name:      "Zero TTL",
			userID:    1,
			ttl:       0,
			scope:     ScopeActivation,
			wantError: false,
		},
		{
			name:      "Negative TTL",
			userID:    1,
			ttl:       -24 * time.Hour,
			scope:     ScopeActivation,
			wantError: false,
		},
		{
			name:      "Zero UserID",
			userID:    0,
			ttl:       24 * time.Hour,
			scope:     ScopeActivation,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := generateToken(tt.userID, tt.ttl, tt.scope)

			// Check error
			if (err != nil) != tt.wantError {
				t.Errorf("generateToken() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				// Check if token is not nil
				if token == nil {
					t.Error("generateToken() returned nil token")
					return
				}

				// Check token fields
				if token.UserID != tt.userID {
					t.Errorf("token.UserID = %v, want %v", token.UserID, tt.userID)
				}

				if token.Scope != tt.scope {
					t.Errorf("token.Scope = %v, want %v", token.Scope, tt.scope)
				}

				// Check plaintext is not empty
				if token.Plaintext == "" {
					t.Error("token.Plaintext is empty")
				}

				// Check hash is not empty
				if len(token.Hash) == 0 {
					t.Error("token.Hash is empty")
				}

				// Check if hash length is correct (SHA-256 produces 32 bytes)
				if len(token.Hash) != 32 {
					t.Errorf("token.Hash length = %v, want 32", len(token.Hash))
				}

				// Check expiry time
				expectedExpiry := time.Now().Add(tt.ttl)
				timeDiff := token.Expiry.Sub(expectedExpiry)
				if timeDiff < -time.Second || timeDiff > time.Second {
					t.Errorf("token.Expiry differs from expected by %v", timeDiff)
				}

				// Check if plaintext tokens are unique
				token2, _ := generateToken(tt.userID, tt.ttl, tt.scope)
				if token.Plaintext == token2.Plaintext {
					t.Error("Generated tokens should be unique")
				}
			}
		})
	}
}

func TestTokenPlaintextFormat(t *testing.T) {
	token, err := generateToken(1, 24*time.Hour, ScopeActivation)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Check if the plaintext token follows base32 format
	if len(token.Plaintext) != 26 { // Base32 encoding of 16 random bytes results in 26 characters
		t.Errorf("token.Plaintext length = %v, want 26", len(token.Plaintext))
	}

	// Check if the plaintext contains only valid base32 characters
	validChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"
	for _, char := range token.Plaintext {
		if !contains(validChars, char) {
			t.Errorf("token.Plaintext contains invalid character: %c", char)
		}
	}
}

// Helper function to check if a string contains a rune
func contains(s string, r rune) bool {
	for _, c := range s {
		if c == r {
			return true
		}
	}
	return false
}
