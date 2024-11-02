package models

import "gorm.io/gorm"

// FacebookUser is a struct that represents the user data that is returned from the Facebook API
type FacebookUser struct {
	Name    string `json:"name"`
	Picture struct {
		Data struct {
			IsSilhouette bool   `json:"is_silhouette"`
			Url          string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
	ID string `json:"id"`
}

// User is a struct that represents the user data that is stored in the database
type User struct {
	gorm.Model
	CreatedAt  string `json:"created_at" gorm:"autoCreateTime"`
	ID         string `json:"id" gorm:"primarykey;size:16"`
	FacebookID string `json:"facebook_id"`
	Name       string `json:"name"`
	Picture    string `json:"picture"`
}

// PaginateDto is a struct that represents the data that is required to paginate a list of users
type PaginateDto struct {
	Offset string `json:"offset"`
	Limit  string `json:"limit"`
}

// GetUserByIDDto is a struct that represents the data that is required to get a user by ID
type GetUserByIDDto struct {
	ID string `json:"id"`
}
