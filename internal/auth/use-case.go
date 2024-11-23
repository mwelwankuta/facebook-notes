package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/mwelwankuta/facebook-notes/pkg/adapters"
	"github.com/mwelwankuta/facebook-notes/pkg/config"
	"github.com/mwelwankuta/facebook-notes/pkg/models"
	"github.com/mwelwankuta/facebook-notes/pkg/utils"
)

type AuthUseCase struct {
	config config.Config
	repo   AuthRepository
	redis  *adapters.RedisClient
}

func NewAuthUseCase(repo AuthRepository, cfg config.Config, redis *adapters.RedisClient) *AuthUseCase {
	return &AuthUseCase{
		repo:   repo,
		config: cfg,
		redis:  redis,
	}
}

func (a *AuthUseCase) CreateUser(userDto models.FacebookUser) (models.User, error) {
	return a.repo.CreateUser(userDto)
}

// AuthenticateUser authenticates a user using the facebook code
func (a *AuthUseCase) AuthenticateUser(code string) (AuthenticateUserResponse, error) {
	var user models.User

	// get facebook access token
	accessToken, err := adapters.GetFacebookUserAccessToken(code, a.config.OpenGraphClientID, a.config.OpenGraphClientSecret)
	if err != nil {
		return AuthenticateUserResponse{
			User:  user,
			Token: "",
		}, err
	}

	// get facebook user profile
	userDto, err := adapters.FetchUserProfile(accessToken)
	if err != nil {
		return AuthenticateUserResponse{User: user, Token: ""}, err
	}

	// get user from database
	user, err = a.repo.GetUserByID(userDto.ID)
	if err != nil {
		return AuthenticateUserResponse{User: user, Token: ""}, err
	}

	if user.ID == "" {
		// create user if not found
		user, err = a.repo.CreateUser(userDto)
		if err != nil {
			return AuthenticateUserResponse{User: user, Token: ""}, err
		}
	}

	jwtToken, err := utils.GenerateJwtToken(a.config.OpenGraphClientSecret, user, accessToken)
	if err != nil {
		return AuthenticateUserResponse{User: user, Token: ""}, err
	}

	return AuthenticateUserResponse{User: user, Token: jwtToken}, nil
}

func (a *AuthUseCase) GetAllUsers(dto models.PaginateDto) ([]models.User, error) {
	return a.repo.GetAllUsers(dto)
}

func (a *AuthUseCase) GetUserByID(userId string) (models.User, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:%s", userId)

	// Try to get from cache first
	var user models.User
	err := a.redis.Get(ctx, cacheKey, &user)
	if err == nil {
		return user, nil
	}

	// If not in cache, get from database
	user, err = a.repo.GetUserByID(userId)
	if err != nil {
		return models.User{}, err
	}

	// Cache the result
	a.redis.Set(ctx, cacheKey, user, 15*time.Minute)
	return user, nil
}

// UpdateUserRole updates a user's role
func (a *AuthUseCase) UpdateUserRole(userId string, role string) (models.User, error) {
	// Validate role
	validRoles := []string{models.RoleUser, models.RoleModerator, models.RoleAdmin}
	isValidRole := false
	for _, r := range validRoles {
		if r == role {
			isValidRole = true
			break
		}
	}
	if !isValidRole {
		return models.User{}, fmt.Errorf("invalid role: %s", role)
	}

	user, err := a.repo.UpdateUserRole(userId, role)
	if err != nil {
		return models.User{}, err
	}

	// Invalidate cache
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:%s", userId)
	a.redis.Delete(ctx, cacheKey)

	return user, nil
}

// UpdateUserStatus updates a user's active status
func (a *AuthUseCase) UpdateUserStatus(userId string, isActive bool) (models.User, error) {
	return a.repo.UpdateUserStatus(userId, isActive)
}

// GetUserByFacebookID returns a user by their Facebook ID
func (a *AuthUseCase) GetUserByFacebookID(facebookId string) (models.User, error) {
	var user models.User
	result := a.repo.db.Where("facebook_id = ?", facebookId).First(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

// ValidateUserRole checks if a user has the required role
func (a *AuthUseCase) ValidateUserRole(userId string, requiredRole string) (bool, error) {
	user, err := a.GetUserByID(userId)
	if err != nil {
		return false, err
	}

	if user.Role == models.RoleAdmin {
		return true, nil
	}

	return user.Role == requiredRole, nil
}

// GetCurrentUserProfile gets the current user's full profile
func (a *AuthUseCase) GetCurrentUserProfile(userId string) (models.User, error) {
	user, err := a.GetUserByID(userId)
	if err != nil {
		return models.User{}, err
	}

	if !user.IsActive {
		return models.User{}, fmt.Errorf("user account is inactive")
	}

	return user, nil
}

// DeactivateUser deactivates a user account
func (a *AuthUseCase) DeactivateUser(userId string) error {
	_, err := a.UpdateUserStatus(userId, false)
	return err
}

// ReactivateUser reactivates a user account
func (a *AuthUseCase) ReactivateUser(userId string) error {
	_, err := a.UpdateUserStatus(userId, true)
	return err
}
