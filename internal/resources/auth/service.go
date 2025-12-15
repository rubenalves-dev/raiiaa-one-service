package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	repo "github.com/rubenalves-dev/raiiaa-one-service/internal/adapters/postgres/sqlc"
	"github.com/rubenalves-dev/raiiaa-one-service/internal/env"
	refreshtokens "github.com/rubenalves-dev/raiiaa-one-service/internal/resources/refresh-tokens"
	"github.com/rubenalves-dev/raiiaa-one-service/internal/resources/users"
)

const (
	ACCESS_TOKEN_TTL time.Duration = 15 * time.Minute
)

var jwtSecret = env.GetString("JWT_TOKEN", "5346cd290c0ffa229eca4b65a92e79b3")

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken       = errors.New("invalid token")
	ErrExpiredToken       = errors.New("token has expired")
)

type Service interface {
	Register(ctx context.Context, args repo.CreateUserParams) (repo.User, error)
	Login(ctx context.Context, args LoginParams) (string, error)
	ValidateToken(tokenString string) (jwt.MapClaims, error)
	GenerateAccessToken(user repo.User, secret string, accessTokenTTL time.Duration) (string, error)
}

type service struct {
	usersService  users.Service
	tokensService refreshtokens.Service
}

func NewService(userService users.Service, tokensService refreshtokens.Service) *service {
	return &service{
		usersService:  userService,
		tokensService: tokensService,
	}
}

func (s *service) Register(ctx context.Context, args repo.CreateUserParams) (repo.User, error) {
	return s.usersService.CreateUser(ctx, args)
}

func (s *service) Login(ctx context.Context, args LoginParams) (string, error) {
	user, err := s.usersService.GetUserByEmail(ctx, args.Email)
	if err != nil {
		return "", err
	}
	if err := users.VerifyPassword(user.PasswordHash, args.Password); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := s.GenerateAccessToken(user, jwtSecret, ACCESS_TOKEN_TTL)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *service) ValidateToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, jwt.ErrTokenExpired
		}
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

func (s *service) GenerateAccessToken(user repo.User, secret string, accessTokenTTL time.Duration) (string, error) {
	expirationTime := time.Now().Add(accessTokenTTL)

	claims := jwt.MapClaims{
		"sub":      user.ID.String(),
		"username": user.FullName,
		"email":    user.Email,
		"exp":      expirationTime.Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
