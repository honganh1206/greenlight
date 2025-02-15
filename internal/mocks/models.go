package mocks

import (
	"greenlight.honganhpham.net/internal/data"
	"greenlight.honganhpham.net/internal/mailer"
)

func NewMockModels() *data.Models {
	return &data.Models{
		Movies: MockMovieModel{},
		Users:  newMockUserModel(),
	}
}

func newMockUserModel() *MockUserModel {
	return &MockUserModel{
		users: map[string]*data.User{
			"mock@example.com": mockUser,
		},
	}
}

func NewMockMailer() *mailer.Mailer {
	return &mailer.Mailer{
		Dialer: mailer.NewDialer("localhost", 25, "username@example.com", "password"),
		Sender: "sender@example.com",
	}
}
