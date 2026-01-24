package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/simple-wallet/internal/model"
	"github.com/simple-wallet/internal/repository"
)

type WalletService interface {
	GetBalance(ctx context.Context, walletID uuid.UUID) (*model.Wallet, error)
	ProcessOperation(ctx context.Context, req *model.WalletOperationRequest) (*model.Wallet, error)
}

type walletService struct {
	repo repository.WalletRepository
}

func NewWalletService(repo repository.WalletRepository) WalletService {
	return &walletService{repo: repo}
}

func (s *walletService) GetBalance(ctx context.Context, walletID uuid.UUID) (*model.Wallet, error) {
	return s.repo.GetByID(ctx, walletID)
}

func (s *walletService) ProcessOperation(ctx context.Context, req *model.WalletOperationRequest) (*model.Wallet, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return s.repo.ProcessOperation(ctx, req)
}
