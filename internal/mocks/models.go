package mocks

import (
	"greenlight.honganhpham.net/internal/data"
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
