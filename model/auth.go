package model

import (
	"time"

	"github.com/google/uuid"
)

type RegisterReq struct {
	Username string `json:"username" db:"username" validate:"required,lte=32"`
	Email    string `json:"email" db:"email" validate:"required,email"`
	Password string `json:"password" db:"password" validate:"required,gte=8"`
}

type OtpReq struct {
	Email string `json:"email" db:"email" validate:"required,email"`
}

type OtpVerifyReq struct {
	Email string `json:"email" db:"email" validate:"required,email"`
	Otp   string `json:"otp" validate:"required"`
}

type LoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8"`
}

type LoginRes struct {
	JwtToken     string `json:"jwtToken"`
	RefreshToken string `json:"refreshToken"`
}

type SessionsRes struct {
	IPAddress string    `json:"ip_address"`
	ExpiresAt time.Time `json:"expires_at"`
	UserAgent string    `json:"user_agent"`
	DeviceId  string    `json:"device_id"`
}

type DeleteToken struct {
	UserId uuid.UUID `json:"userId"`
	Token  string    `json:"token"`
}

type ChangePassword struct {
	Email       string `json:"email" db:"email" validate:"required,email"`
	NewPassword string `json:"password" db:"password" validate:"required,gte=8"`
}
