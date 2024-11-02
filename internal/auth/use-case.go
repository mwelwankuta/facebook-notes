package auth

import (
	"github.com/mwelwankuta/facebook-notes/pkg/adapters"
	"github.com/mwelwankuta/facebook-notes/pkg/models"
)

type AuthUseCase struct {
	OpenGraphClientID     string
	OpenGraphClientSecret string
	repo                  AuthRepository
}

func NewAuthUseCase(repo AuthRepository, clientId string, clientSecret string) *AuthUseCase {
	return &AuthUseCase{repo: repo, OpenGraphClientID: clientId, OpenGraphClientSecret: clientSecret}
}

func (a *AuthUseCase) CreateUser(userDto models.FacebookUser) (models.User, error) {
	return a.repo.CreateUser(userDto)
}

// AuthenticateUser authenticates a user using the facebook code
func (a *AuthUseCase) AuthenticateUser(code string) (models.User, error) {
	var user models.User

	// get facebook access token
	accessToken, err := adapters.GetFacebookUserAccessToken(code, a.OpenGraphClientID, a.OpenGraphClientSecret)
	if err != nil {
		return user, err
	}

	// get facebook user profile
	userDto, err := adapters.FetchUserProfile(accessToken)
	if err != nil {
		return user, err
	}

	// get user from database
	user, err = a.repo.GetUserByID(userDto.ID)
	if err != nil {
		return user, err
	}

	if user.ID == "" {
		// create user if not found
		user, err = a.repo.CreateUser(userDto)
		if err != nil {
			return user, err
		}
	}

	return user, nil
}

func (a *AuthUseCase) GetAllUsers(dto models.PaginateDto) ([]models.User, error) {
	return a.repo.GetAllUsers(dto)
}

func (a *AuthUseCase) GetUserByID(userId string) (models.User, error) {
	return a.repo.GetUserByID(userId)
}
