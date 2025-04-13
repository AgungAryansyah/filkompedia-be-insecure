package jwt

import (
	"errors"
	"os"
	"strconv"
	"time"

	lib_jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/yogarn/filkompedia-be/pkg/response"
)

type IJwt interface {
	CreateToken(userId uuid.UUID) (string, error)
	ValidateToken(tokenString string) (uuid.UUID, error)
}

type jwt struct {
	SecretKey   string
	ExpiredTime time.Duration
}

type Claims struct {
	UserId uuid.UUID
	lib_jwt.RegisteredClaims
}

func Init() IJwt {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	expTime, err := strconv.Atoi(os.Getenv("JWT_EXPIRED_TIME"))
	if err != nil {
		panic(err)
	}

	return &jwt{
		SecretKey:   secretKey,
		ExpiredTime: time.Duration(expTime) * time.Minute,
	}
}

func (j *jwt) CreateToken(userId uuid.UUID) (string, error) {
	claim := &Claims{
		UserId: userId,
		RegisteredClaims: lib_jwt.RegisteredClaims{
			ExpiresAt: lib_jwt.NewNumericDate(time.Now().Add(j.ExpiredTime)),
		},
	}

	token := lib_jwt.NewWithClaims(lib_jwt.SigningMethodHS256, claim)
	tokenString, err := token.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *jwt) ValidateToken(tokenString string) (uuid.UUID, error) {
	var claim Claims
	var userId uuid.UUID

	token, err := lib_jwt.ParseWithClaims(tokenString, &claim, func(t *lib_jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		if errors.Is(err, lib_jwt.ErrTokenExpired) {
			return userId, &response.ExpiredToken
		}
		return userId, err
	}

	if !token.Valid {
		return userId, errors.New("invalid token")
	}

	userId = claim.UserId
	return userId, nil
}
