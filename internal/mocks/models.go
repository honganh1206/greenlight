package mocks

import (
	"time"

	"greenlight.honganhpham.net/internal/data"
	"greenlight.honganhpham.net/internal/mailer"
)

func NewMockModels() *data.Models {
	return &data.Models{
		Movies:      MockMovieModel{},
		Users:       newMockUserModel(),
		Tokens:      newMockTokenModel(),
		Permissions: newMockPermissionModel(),
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

func newMockPermissionModel() *MockPermissionModel {
	permissions := make(map[int64]data.Permissions)

	// Add some mock permissions for test users
	permissions[1] = data.Permissions{"movies:read", "movies:write"} // for mock@example.com user
	permissions[2] = data.Permissions{"movies:read"}                 // for not_activated@example.com user

	return &MockPermissionModel{
		permissions: permissions,
	}
}
