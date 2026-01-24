package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/simple-wallet/internal/mocks"
	"github.com/simple-wallet/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWalletService_GetBalance(t *testing.T) {
	t.Run("returns wallet when found", func(t *testing.T) {
		repo := mocks.NewMockWalletRepository(t)
		svc := NewWalletService(repo)

		walletID := uuid.New()
		expected := &model.Wallet{ID: walletID, Balance: 1000}

		repo.EXPECT().
			GetByID(mock.Anything, walletID).
			Return(expected, nil)

		result, err := svc.GetBalance(context.Background(), walletID)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("returns error when wallet not found", func(t *testing.T) {
		repo := mocks.NewMockWalletRepository(t)
		svc := NewWalletService(repo)

		walletID := uuid.New()

		repo.EXPECT().
			GetByID(mock.Anything, walletID).
			Return(nil, model.ErrWalletNotFound)

		result, err := svc.GetBalance(context.Background(), walletID)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, model.ErrWalletNotFound)
	})
}

func TestWalletService_ProcessOperation(t *testing.T) {
	t.Run("successful deposit", func(t *testing.T) {
		repo := mocks.NewMockWalletRepository(t)
		svc := NewWalletService(repo)

		walletID := uuid.New()
		req := &model.WalletOperationRequest{
			WalletID:      walletID,
			OperationType: model.OperationDeposit,
			Amount:        1000,
		}
		expected := &model.Wallet{ID: walletID, Balance: 1000}

		repo.EXPECT().
			ProcessOperation(mock.Anything, req).
			Return(expected, nil)

		result, err := svc.ProcessOperation(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})

	t.Run("returns validation error for zero amount", func(t *testing.T) {
		repo := mocks.NewMockWalletRepository(t)
		svc := NewWalletService(repo)

		req := &model.WalletOperationRequest{
			WalletID:      uuid.New(),
			OperationType: model.OperationDeposit,
			Amount:        0,
		}

		result, err := svc.ProcessOperation(context.Background(), req)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, model.ErrInvalidAmount)
	})

	t.Run("returns validation error for invalid operation", func(t *testing.T) {
		repo := mocks.NewMockWalletRepository(t)
		svc := NewWalletService(repo)

		req := &model.WalletOperationRequest{
			WalletID:      uuid.New(),
			OperationType: "INVALID",
			Amount:        1000,
		}

		result, err := svc.ProcessOperation(context.Background(), req)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, model.ErrInvalidOperation)
	})

	t.Run("returns error on insufficient funds", func(t *testing.T) {
		repo := mocks.NewMockWalletRepository(t)
		svc := NewWalletService(repo)

		walletID := uuid.New()
		req := &model.WalletOperationRequest{
			WalletID:      walletID,
			OperationType: model.OperationWithdraw,
			Amount:        1000,
		}

		repo.EXPECT().
			ProcessOperation(mock.Anything, req).
			Return(nil, model.ErrInsufficientFunds)

		result, err := svc.ProcessOperation(context.Background(), req)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, model.ErrInsufficientFunds)
	})
}
