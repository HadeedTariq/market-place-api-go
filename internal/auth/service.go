package auth

import (
	"context"

	repo "github.com/HadeedTariq/market-place-api-go/internal/adapters/postgresql/sqlc"
)

type Service interface {
	FindExistingUserByEmail(ctx context.Context, email string) (int32, error)
	InsertUser(ctx context.Context, arg repo.InsertUserParams) error
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) FindExistingUserByEmail(ctx context.Context, email string) (int32, error) {
	return s.repo.FindExistingUserByEmail(ctx, email)
}

func (s *svc) InsertUser(ctx context.Context, arg repo.InsertUserParams) error {
	return s.repo.InsertUser(ctx, arg)
}
