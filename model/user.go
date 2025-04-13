package model

import (
	"github.com/google/uuid"
	"github.com/yogarn/filkompedia-be/entity"
)

type ProfilesReq struct {
	Page     int `json:"page" validate:"required,min=1"`
	PageSize int `json:"page_size" validate:"required,min=1"`
}

type Profile struct {
	Id             uuid.UUID `json:"id" db:"id"`
	Username       string    `json:"username" db:"username"`
	Email          string    `json:"email" db:"email"`
	RoleId         int       `json:"roleId" db:"role_id"`
	ProfilePicture string    `json:"profilePicture" db:"profile_picture"`
}

type RoleUpdate struct {
	Id     uuid.UUID `json:"id" db:"id" validate:"required"`
	RoleId int       `json:"roleId" db:"role_id" validate:"required,min=0,max=1"`
}

type EditProfile struct {
	Id             uuid.UUID `json:"id" db:"id" validate:"required,uuid"`
	Username       string    `json:"username" db:"username" validate:"omitempty,lte=32"`
	ProfilePicture string    `json:"profilePicture" db:"profile_picture" validate:"omitempty,url"`
}

func UserToProfile(user entity.User) Profile {
	return Profile{
		Id:             user.Id,
		Username:       user.Username,
		Email:          user.Email,
		RoleId:         user.RoleId,
		ProfilePicture: user.ProfilePicture,
	}
}
