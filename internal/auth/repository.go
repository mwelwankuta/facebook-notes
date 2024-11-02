package auth

import "gorm.io/gorm"

type AuthRepository struct {
	db gorm.DB
}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{}
}

func (a *AuthRepository) GetAllUsers() ([]User, error) {
	var users []User

	result := a.db.Find(&users)
	if result.Error != nil {
		return users, result.Error
	}

	return users, nil
}

func (a *AuthRepository) GetUserByID(userId string) (User, error) {
	var user User

	result := a.db.Where("id = ?", userId).Find(&user)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func (a *AuthRepository) CreateUser(userDto FacebookUser) (User, error) {
	var newUser = User{
		FacebookID: userDto.ID,
		Name:       userDto.Name,
		Picture:    userDto.Picture.Data.Url,
	}

	result := a.db.Create(&newUser)
	if result.Error != nil {
		return newUser, result.Error
	}

	return a.GetUserByID(newUser.ID)
}
