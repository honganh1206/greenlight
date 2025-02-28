package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"

	"greenlight.honganhpham.net/internal/validator"
)

const (
	// TODO: Add more scopes!
	ScopeActivation = "activation"
)

type Token struct {
	Plaintext string
	Hash      []byte
	UserID    int64
	Expiry    time.Time
	Scope     string
}

type TokenModel struct {
	DB *sql.DB
}

type TokenModelInterface interface {
	New(userID int64, ttl time.Duration, scope string) (*Token, error)
	Insert(token *Token) error
	DeleteAllForUser(scope string, userID int64) error
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

func (m TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)

	return token, err
}

func (m TokenModel) Insert(token *Token) error {
	query := `
	INSERT INTO tokens (hash, user_id, expiry, scope)
	VALUES ($1, $2, $3, $4)`

	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)

	return err
}

func (m TokenModel) DeleteAllForUser(scope string, userID int64) error {
	query := `
		DELETE FROM tokens
		WHERE scope = $1 AND user_id = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, scope, userID)
	return err
}

func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl), // Add time-to-live
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)

	// Fill the byte slice with random bytes with CSPRNG
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	// A token will look like this: Y3QMGX3PJ3WLRL2YRTQGQ6KRHU - it will have an underlying randomness of 16 bytes
	// NoPadding to remove the "=" at the end
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	// Sum256() function returns an array, so we convert the result to a slice to work with it easier
	// Since slices are more flexible to be used
	// And we have to do it for the token.Hash field as a slice :)
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil

}
