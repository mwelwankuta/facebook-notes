package auth

import (
	"github.com/mwelwankuta/facebook-notes/pkg/adapters"
	"github.com/mwelwankuta/facebook-notes/pkg/config"
	"github.com/mwelwankuta/facebook-notes/pkg/models"
	"github.com/mwelwankuta/facebook-notes/pkg/utils"
)

type AuthUseCase struct {
	config config.Config
	repo   AuthRepository
}

func NewAuthUseCase(repo AuthRepository, cfg config.Config) *AuthUseCase {
	return &AuthUseCase{
		repo:   repo,
		config: cfg,
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
	return a.repo.GetUserByID(userId)
}
