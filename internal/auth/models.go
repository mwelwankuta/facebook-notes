package auth

import "github.com/mwelwankuta/facebook-notes/pkg/models"

// AuthenticateUserResponse is a struct that represents the data that is returned when a user is authenticated
type AuthenticateUserResponse struct {
	User  models.User `json:"user"`
	Token string      `json:"token"`
}

// GetUserByIDDto is a struct that represents the data that is required to get a user by ID
type GetUserByIDDto struct {
	ID string `json:"id"`
}
