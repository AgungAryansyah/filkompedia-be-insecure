package entity

import "github.com/google/uuid"

type User struct {
	Id             uuid.UUID `json:"id" db:"id"`
	Username       string    `json:"username" db:"username"`
	Email          string    `json:"email" db:"email"`
	Password       string    `json:"password" db:"password"`
	RoleId         int       `json:"roleId" db:"role_id"`
	IsVerified     bool      `json:"isVerified" db:"is_verified"`
	ProfilePicture string    `json:"profilePicture" db:"profile_picture"`
}
