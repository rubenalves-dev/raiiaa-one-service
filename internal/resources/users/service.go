package users

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	repo "github.com/rubenalves-dev/raiiaa-one-service/internal/adapters/postgres/sqlc"
	"github.com/rubenalves-dev/raiiaa-one-service/internal/common/validators"
)

var (
	ErrEmailAlreadyInUse = errors.New("email already in use")
	ErrUserNotFound      = errors.New("user not found")
)

type Service interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (repo.User, error)
	CreateUser(ctx context.Context, args repo.CreateUserParams) (repo.User, error)
	GetUserByEmail(ctx context.Context, email string) (repo.User, error)
}

type service struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &service{repo: repo}
}

func (s *service) CreateUser(ctx context.Context, args repo.CreateUserParams) (repo.User, error) {
	_, err := s.GetUserByEmail(ctx, args.Email)
	if err == nil {
		return repo.User{}, ErrEmailAlreadyInUse
	}
	if !errors.Is(err, pgx.ErrNoRows) && !errors.Is(err, ErrUserNotFound) {
		return repo.User{}, err
	}

	_, err = UserFullname(args.FullName.String)
	if err != nil {
		return repo.User{}, err
	}

	password, err := Password(args.PasswordHash)
	if err != nil {
		return repo.User{}, err
	}
	args.PasswordHash = string(password)

	user, err := s.repo.CreateUser(ctx, args)
	if err != nil {
		return repo.User{}, err
	}

	return user, nil
}

func (s *service) GetUserByID(ctx context.Context, id uuid.UUID) (repo.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return repo.User{}, ErrUserNotFound
	}
	return user, nil
}

func (s *service) GetUserByEmail(ctx context.Context, email string) (repo.User, error) {
	_, err := validators.Email(email)
	if err != nil {
		return repo.User{}, validators.ErrInvalidEmail
	}
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return repo.User{}, ErrUserNotFound
	}
	return user, nil
}
