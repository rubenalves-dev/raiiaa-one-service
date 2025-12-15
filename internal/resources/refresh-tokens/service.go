package refreshtokens

import (
	"context"

	repo "github.com/rubenalves-dev/raiiaa-one-service/internal/adapters/postgres/sqlc"
)

type Service interface {
	CreateRefreshToken(ctx context.Context, args repo.CreateRefreshTokenParams) (repo.RefreshToken, error)
	GetResfreshToken(ctx context.Context, token string) (repo.RefreshToken, error)
}

type service struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &service{repo: repo}
}

func (s *service) CreateRefreshToken(ctx context.Context, args repo.CreateRefreshTokenParams) (repo.RefreshToken, error) {
	rtoken, err := s.repo.CreateRefreshToken(ctx, args)
	if err != nil {
		return repo.RefreshToken{}, err
	}
	return rtoken, nil
}

func (s *service) GetResfreshToken(ctx context.Context, token string) (repo.RefreshToken, error) {
	rtoken, err := s.repo.GetRefreshToken(ctx, token)
	if err != nil {
		return repo.RefreshToken{}, err
	}
	return rtoken, nil
}

func (s *service) RevokeRefreshToken(ctx context.Context, token string) error {
	return s.repo.RevokeRefreshToken(ctx, token)
}
