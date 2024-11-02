package auth

import (
	"strconv"

	"github.com/mwelwankuta/facebook-notes/pkg/models"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

// GetAllUsers returns all users
func (a *AuthRepository) GetAllUsers(dto models.PaginateDto) ([]models.User, error) {
	var users []models.User

	// convert string params to int for pagination
	offset, err := strconv.Atoi(dto.Offset)
	if err != nil {
		return nil, err
	}
	limit, err := strconv.Atoi(dto.Limit)
	if err != nil {
		return nil, err
	}

	result := a.db.Find(&users).Offset(offset).Limit(limit)
	if result.Error != nil {
		return users, result.Error
	}

	return users, nil
}

func (a *AuthRepository) GetUserByID(userId string) (models.User, error) {
	var user models.User

	result := a.db.Where("id = ?", userId).Find(&user)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func (a *AuthRepository) CreateUser(userDto models.FacebookUser) (models.User, error) {
	var newUser = models.User{
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
