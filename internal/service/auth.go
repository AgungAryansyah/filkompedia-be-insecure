package service

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/yogarn/filkompedia-be/entity"
	"github.com/yogarn/filkompedia-be/internal/repository"
	"github.com/yogarn/filkompedia-be/model"
	"github.com/yogarn/filkompedia-be/pkg/bcrypt"
	"github.com/yogarn/filkompedia-be/pkg/jwt"
	"github.com/yogarn/filkompedia-be/pkg/response"
	"github.com/yogarn/filkompedia-be/pkg/smtp"
)

type IAuthService interface {
	Register(registerReq *model.RegisterReq) (user *entity.User, err error)
	SendOTP(email string) error
	VerifyOTP(email, otp string) error
	Login(loginReq *model.LoginReq, ipAddress string, userAgent string, expiry int) (loginRes *model.LoginRes, err error)
	GetSessions(userId uuid.UUID) (*[]model.SessionsRes, error)
	ExchangeToken(token string, expiry int) (jwtToken string, newToken string, err error)

	ClearToken(userId uuid.UUID) error
	DeleteToken(info *model.DeleteToken) error
	ChangePassword(req *model.ChangePassword) error
}

type AuthService struct {
	AuthRepository repository.IAuthRepository
	UserRepository repository.IUserRepository
	Bcrypt         bcrypt.IBcrypt
	Jwt            jwt.IJwt
	Smtp           *smtp.SMTPClient
}

func NewAuthService(authRepository repository.IAuthRepository, userRepository repository.IUserRepository, bcrypt bcrypt.IBcrypt, jwt jwt.IJwt, smtp *smtp.SMTPClient) IAuthService {
	return &AuthService{
		AuthRepository: authRepository,
		UserRepository: userRepository,
		Bcrypt:         bcrypt,
		Jwt:            jwt,
		Smtp:           smtp,
	}
}

func (s *AuthService) Register(registerReq *model.RegisterReq) (user *entity.User, err error) {
	hashedpassword, err := s.Bcrypt.GenerateFromPassword(registerReq.Password)
	if err != nil {
		return nil, err
	}

	user = &entity.User{
		Id:         uuid.New(),
		Username:   registerReq.Username,
		Email:      registerReq.Email,
		Password:   hashedpassword,
		RoleId:     0,
		IsVerified: false,
	}

	err = s.AuthRepository.Register(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) SendOTP(email string) error {
	// can not check user verification status
	// this func might be needed for reset password

	otp := generateOTP()
	err := s.Smtp.SendEmail(email, "FilkomPedia OTP Verification", "Do not share this code with others. Your OTP code is "+otp)
	if err != nil {
		return err
	}

	err = s.AuthRepository.StoreOTP(email, otp)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) VerifyOTP(email, otp string) error {
	status := s.AuthRepository.VerifyOTP(email, otp)
	if !status {
		return &response.InvalidOTP
	}

	err := s.AuthRepository.DeleteOTP(email)
	if err != nil {
		return err
	}

	err = s.AuthRepository.VerifyEmail(email)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Login(loginReq *model.LoginReq, ipAddress string, userAgent string, expiry int) (loginRes *model.LoginRes, err error) {
	user, err := s.UserRepository.GetUserByEmail(loginReq.Email)
	if err != nil {
		// prevent user from guessing that an account is existed or not
		return nil, &response.InvalidCredentials
	}

	if !user.IsVerified {
		return nil, &response.UserUnverified
	}

	if user.Id == uuid.Nil {
		return nil, &response.InvalidCredentials
	}

	err = s.Bcrypt.CompareAndHashPassword(user.Password, loginReq.Password)
	if err != nil {
		return nil, &response.InvalidCredentials
	}

	token, err := s.Jwt.CreateToken(user.Id)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateRandomString(32)
	if err != nil {
		return nil, err
	}

	session := &entity.Session{
		UserId:    user.Id,
		Token:     refreshToken,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		DeviceId:  generateDeviceID(),
		ExpiresAt: time.Now().Add(time.Duration(expiry) * time.Second),
	}

	// before proceed, clear existing sessions
	_ = s.AuthRepository.DeleteExpiredToken(user.Id)

	err = s.AuthRepository.Login(session)
	if err != nil {
		return nil, err
	}

	return &model.LoginRes{
		JwtToken:     token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) GetSessions(userId uuid.UUID) (*[]model.SessionsRes, error) {
	sessions, err := s.AuthRepository.GetSessions(userId)
	if err != nil {
		return nil, err
	}

	var sessionsRes []model.SessionsRes

	for _, session := range *sessions {
		sessionsRes = append(sessionsRes, model.SessionsRes{
			IPAddress: session.IPAddress,
			ExpiresAt: session.ExpiresAt,
			UserAgent: session.UserAgent,
			DeviceId:  session.DeviceId,
		})
	}

	return &sessionsRes, nil
}

func (s *AuthService) ExchangeToken(token string, expiry int) (jwtToken string, newToken string, err error) {
	currentSession, err := s.AuthRepository.CheckUserSession(token)
	if err != nil {
		return "", "", err
	}

	jwtToken, err = s.Jwt.CreateToken(currentSession.UserId)
	if err != nil {
		return "", "", err
	}

	newToken, err = generateRandomString(32)
	if err != nil {
		return "", "", err
	}

	expiresAt := time.Now().Add(time.Duration(expiry) * time.Second)

	err = s.AuthRepository.ReplaceToken(token, newToken, currentSession.UserId, expiresAt)
	if err != nil {
		return "", "", err
	}

	return jwtToken, newToken, nil
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}

func generateDeviceID() string {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func generateOTP() string {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%06d", n.Int64())
}

func (s *AuthService) ClearToken(userId uuid.UUID) error {
	return s.AuthRepository.ClearToken(userId)
}

func (s *AuthService) DeleteToken(info *model.DeleteToken) error {
	return s.AuthRepository.DeleteToken(info.UserId, info.Token)
}

func (s *AuthService) ChangePassword(req *model.ChangePassword) error {
	hashedpassword, err := s.Bcrypt.GenerateFromPassword(req.NewPassword)
	if err != nil {
		return err
	}

	if err := s.AuthRepository.ChangePassword(req.Email, hashedpassword); err != nil {
		return err
	}

	return nil
}
