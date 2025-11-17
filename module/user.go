// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

import (
	"errors"
	"time"

	"github.com/clivern/mut/db"
	"github.com/clivern/mut/service"
	"github.com/google/uuid"
)

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrUserEmailAlreadyExists = errors.New("user with this email already exists")
)

// User handles user management operations.
type User struct {
	UserRepository *db.UserRepository
}

// NewUser creates a new user module instance.
func NewUser(repo *db.UserRepository) *User {
	return &User{UserRepository: repo}
}

// CreateUserOptions contains options for creating a user.
type CreateUserOptions struct {
	Email    string
	Password string
	Role     string
	IsActive bool
}

// CreateUser creates a new user.
func (u *User) CreateUser(options *CreateUserOptions) (*db.User, error) {
	// Check if user with email already exists
	existingUser, err := u.UserRepository.GetByEmail(options.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserEmailAlreadyExists
	}

	hashedPassword, err := service.HashPassword(options.Password)
	if err != nil {
		return nil, err
	}

	user := &db.User{
		Email:       options.Email,
		Password:    hashedPassword,
		Role:        options.Role,
		APIKey:      uuid.New().String(),
		IsActive:    options.IsActive,
		LastLoginAt: time.Time{},
	}

	if err := u.UserRepository.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetUser retrieves a user by ID.
func (u *User) GetUser(userID int64) (*db.User, error) {
	user, err := u.UserRepository.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// UpdateUserOptions contains options for updating a user.
type UpdateUserOptions struct {
	UserID   int64
	Email    string
	Password string
	Role     string
	IsActive bool
}

// UpdateUser updates an existing user.
func (u *User) UpdateUser(options *UpdateUserOptions) (*db.User, error) {
	// Get existing user
	user, err := u.UserRepository.GetByID(options.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Check if email is being changed and if it's already taken
	if options.Email != user.Email {
		existingUser, err := u.UserRepository.GetByEmail(options.Email)
		if err != nil {
			return nil, err
		}
		if existingUser != nil && existingUser.ID != options.UserID {
			return nil, ErrUserEmailAlreadyExists
		}
	}

	user.Email = options.Email
	user.Role = options.Role
	user.IsActive = options.IsActive

	// Update password only if provided
	if options.Password != "" {
		hashedPassword, err := service.HashPassword(options.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashedPassword
	}

	if err := u.UserRepository.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// ListUsersOptions contains options for listing users.
type ListUsersOptions struct {
	Limit  int
	Offset int
}

// ListUsersResult contains the result of listing users.
type ListUsersResult struct {
	Users []*db.User
	Total int64
}

// ListUsers retrieves a list of users with pagination.
func (u *User) ListUsers(options *ListUsersOptions) (*ListUsersResult, error) {
	users, err := u.UserRepository.List(options.Limit, options.Offset)
	if err != nil {
		return nil, err
	}

	total, err := u.UserRepository.Count()
	if err != nil {
		return nil, err
	}

	return &ListUsersResult{
		Users: users,
		Total: total,
	}, nil
}

// DeleteUser deletes a user by ID.
func (u *User) DeleteUser(userID int64) error {
	// Check if user exists
	user, err := u.UserRepository.GetByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	// Delete user
	return u.UserRepository.Delete(userID)
}
