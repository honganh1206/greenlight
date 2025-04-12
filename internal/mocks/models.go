package mocks

import (
	"time"

	"greenlight.honganhpham.net/internal/data"
	"greenlight.honganhpham.net/internal/mailer"
)

func NewMockModels() *data.Models {
	return &data.Models{
		Movies: MockMovieModel{},
		Users:  newMockUserModel(),
		Token:  newMockTokenModel(),
	}
}

func newMockUserModel() *MockUserModel {
	users := make(map[string]*data.User)

	// Add test users with different states
	users["mock@example.com"] = &data.User{
		ID:        1,
		CreatedAt: time.Now(),
		Name:      "Mock User",
		Email:     "mock@example.com",
		Activated: true,
		Version:   1,
	}

	users["not_activated@example.com"] = &data.User{
		ID:        2,
		CreatedAt: time.Now(),
		Name:      "Activated User",
		Email:     "not_activated@example.com",
		Activated: false,
		Version:   1,
	}

	return &MockUserModel{
		users: users,
	}
}

// NewMockTokenModel creates a new instance of MockTokenModel with initialized fields
func newMockTokenModel() *MockTokenModel {
	return &MockTokenModel{
		Tokens: map[int64]*data.Token{
			1: mockToken,
		},
	}
}

func NewMockMailer() *mailer.Mailer {
	return &mailer.Mailer{
		Dialer: mailer.NewDialer("localhost", 25, "username@example.com", "password"),
		Sender: "sender@example.com",
	}
}
