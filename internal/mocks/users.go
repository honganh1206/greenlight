package mocks

import (
	"time"

	"greenlight.honganhpham.net/internal/data"
)

type MockUserModel struct {
	users map[string]*data.User
}

func (m MockUserModel) Insert(user *data.User) error {
	if _, exists := m.users[user.Email]; exists {
		return data.ErrDuplicateEmail
	}

	user.ID = int64(len(m.users) + 1)
	user.CreatedAt = time.Now()
	user.Version = 1

	m.users[user.Email] = user
	return nil
}

// GetByEmail simulates fetching a user by email
func (m MockUserModel) GetByEmail(email string) (*data.User, error) {
	if user, exists := m.users[email]; exists {
		// Return a copy to prevent modifications to the original
		userCopy := *user
		return &userCopy, nil
	}
	return nil, data.ErrRecordNotFound
}

func (m MockUserModel) Update(user *data.User) error {
	return nil
}

func (m MockUserModel) GetForToken(tokenScope, tokenPlaintext string) (*data.User, error) {
	return m.users["mock@example.com"], nil
}
