package mocks

import (
	"time"

	"greenlight.honganhpham.net/internal/data"
)

type MockUserModel struct {
	users map[string]*data.User
}

// Mock user for testing
var mockUser = &data.User{
	ID:        1,
	CreatedAt: time.Now(),
	Name:      "Mock User",
	Email:     "mock@example.com",
	Activated: true,
	Version:   1,
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
	switch email {
	case "mock@example.com":
		return mockUser, nil
	default:
		return nil, data.ErrRecordNotFound
	}
}

func (m MockUserModel) Update(user *data.User) error {
	return nil
}
