package auth

import (
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

	result := a.db.Find(&users).Offset(dto.Offset).Limit(dto.Limit)
	if result.Error != nil {
		return users, result.Error
	}

	return users, nil
}

// GetUserByID returns a user by ID
func (a *AuthRepository) GetUserByID(userId string) (models.User, error) {
	var user models.User

	result := a.db.Where("id = ?", userId).Find(&user)
	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

// CreateUser creates a new user
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

func (a *AuthRepository) UpdateUserRole(userId string, role string) (models.User, error) {
	result := a.db.Model(&models.User{}).Where("id = ?", userId).Update("role", role)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return a.GetUserByID(userId)
}

func (a *AuthRepository) UpdateUserStatus(userId string, isActive bool) (models.User, error) {
	result := a.db.Model(&models.User{}).Where("id = ?", userId).Update("is_active", isActive)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return a.GetUserByID(userId)
}
