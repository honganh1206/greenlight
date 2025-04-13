package mocks

import "greenlight.honganhpham.net/internal/data"

type MockPermissionModel struct {
	permissions map[int64]data.Permissions
}

func (m MockPermissionModel) GetAllForUser(userID int64) (data.Permissions, error) {
	if permissions, exists := m.permissions[userID]; exists {
		// Return a copy to prevent modifications to the original
		permissionsCopy := make(data.Permissions, len(permissions))
		copy(permissionsCopy, permissions)
		return permissionsCopy, nil
	}
	return data.Permissions{}, nil // Return empty permissions if user not found
}

func (m MockPermissionModel) AddForUser(userID int64, codes ...string) error {
	return nil
}
