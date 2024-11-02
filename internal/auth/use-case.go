package auth

type AuthUseCase struct {
	repo AuthRepository
}

func NewAuthUseCase(repo AuthRepository) *AuthUseCase {
	return &AuthUseCase{repo: repo}
}

func (a *AuthUseCase) CreateUser(userDto FacebookUser) (User, error) {
	return a.repo.CreateUser(userDto)
}

func (a *AuthUseCase) AuthenticateUser(userDto FacebookUser) (User, error) {
	var user User
	// get user from database
	user, err := a.repo.GetUserByID(userDto.ID)
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

func (a *AuthUseCase) GetAllUsers() ([]User, error) {
	return a.repo.GetAllUsers()
}

func (a *AuthUseCase) GetUserByID(userId string) (User, error) {
	return a.repo.GetUserByID(userId)
}
