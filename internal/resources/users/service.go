package users

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	repo "github.com/rubenalves-dev/raiiaa-one-service/internal/adapters/postgres/sqlc"
	"github.com/rubenalves-dev/raiiaa-one-service/internal/common/validators"
)

var (
	ErrEmailInUser  = errors.New("email already in use")
	ErrUserNotFound = errors.New("user not found")
)

type Service interface {
	CreateUser(ctx context.Context, args repo.CreateUserParams) (*repo.User, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

func (s *svc) CreateUser(ctx context.Context, args repo.CreateUserParams) (*repo.User, error) {
	_, err := validators.Email(args.Email)
	if err != nil {
		return nil, err
	}

	_, err = s.repo.GetUserByEmail(ctx, args.Email)
	if err == nil {
		return nil, ErrEmailInUser
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	_, err = UserFullname(args.FullName.String)
	if err != nil {
		return nil, err
	}

	password, err := Password(args.PasswordHash)
	if err != nil {
		return nil, err
	}
	args.PasswordHash = string(password)

	user, err := s.repo.CreateUser(ctx, args)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
