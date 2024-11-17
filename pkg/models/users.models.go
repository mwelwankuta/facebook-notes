package models

import "gorm.io/gorm"

const (
	RoleUser      = "user"
	RoleModerator = "moderator"
	RoleAdmin     = "admin"
)

// FacebookUser is a struct that represents the user data that is returned from the Facebook API
type FacebookUserPicture struct {
	Data struct {
		Url          string `json:"url"`
		IsShilouette bool   `json:"is_silhouette"`
	} `json:"data"`
}

type FacebookUser struct {
	ID      string              `json:"id"`
	Name    string              `json:"name"`
	Picture FacebookUserPicture `json:"picture"`
}

// User is a struct that represents the user data that is stored in the database

type User struct {
	gorm.Model
	ID         string `json:"id" gorm:"primarykey"`
	FacebookID string `json:"facebook_id"`
	Name       string `json:"name"`
	Picture    string `json:"picture"`
	Role       string `json:"role" gorm:"default:user"`
	IsActive   bool   `json:"is_active" gorm:"default:true"`
}

func (u *User) IsModerator() bool {
	return u.Role == RoleModerator || u.Role == RoleAdmin
}
