package service

import (
	"context"

	"github.com/google/uuid"
	"wallet/internal/repository"
)

type Service struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) UpdateBalance(
	ctx context.Context,
	walletID uuid.UUID,
	amount int64,
	op string,
) error {
	return s.repo.UpdateBalance(ctx, walletID, amount, op)
}

func (s *Service) GetBalance(
	ctx context.Context,
	walletID uuid.UUID,
) (int64, error) {
	return s.repo.GetBalance(ctx, walletID)
}

