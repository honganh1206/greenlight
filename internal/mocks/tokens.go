package mocks

import (
	"time"

	"greenlight.honganhpham.net/internal/data"
)

type MockTokenModel struct {
	Tokens        map[int64]*data.Token // Map UserID with the token
	ErrorToReturn error
}

// Mock token for testing
var mockToken = &data.Token{
	Plaintext: "Y3QMGX3PJ3WLRL2YRTQGQ6KRHU",
	Hash:      []byte("mockhash"),
	UserID:    1,
	Expiry:    time.Now().Add(24 * time.Hour),
	Scope:     data.ScopeActivation,
}

func (m MockTokenModel) New(userID int64, ttl time.Duration, scope string) (*data.Token, error) {
	if m.ErrorToReturn != nil {
		return nil, m.ErrorToReturn
	}

	// Create a new token with the provided parameters
	token := &data.Token{
		Plaintext: "MOCKTOKEN123456789ABCDEFGHIJK", // Fixed test token
		Hash:      []byte("mockhash"),
		UserID:    userID,
		Expiry:    time.Now().Add(ttl),
		Scope:     scope,
	}

	// Store the token in our mock storage
	m.Tokens[userID] = token

	return token, nil
}

func (m MockTokenModel) Insert(token *data.Token) error {
	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}

	m.Tokens[token.UserID] = token
	return nil
}

func (m MockTokenModel) DeleteAllForUser(scope string, userID int64) error {
	if m.ErrorToReturn != nil {
		return m.ErrorToReturn
	}
	delete(m.Tokens, userID)
	return nil
}

// Helper methods for testing

// // GetTokenForUser returns the token for a specific user
// func (m *MockTokenModel) GetTokenForUser(userID int64) *data.Token {
// 	return m.Tokens[userID]
// }

// // Reset resets the mock's state
// func (m *MockTokenModel) Reset() {
// 	m.Tokens = make(map[int64]*data.Token)
// 	m.InsertCalled = false
// 	m.DeleteCalled = false
// 	m.LastUserID = 0
// 	m.LastScope = ""
// 	m.LastTTL = 0
// 	m.ErrorToReturn = nil
// }

// // SetError sets an error to be returned by the mock methods
// func (m *MockTokenModel) SetError(err error) {
// 	m.ErrorToReturn = err
// }
