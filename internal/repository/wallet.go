package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/simple-wallet/internal/model"
)

type WalletRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*model.Wallet, error)
	ProcessOperation(ctx context.Context, req *model.WalletOperationRequest) (*model.Wallet, error)
}

type walletRepository struct {
	pool *pgxpool.Pool
}

func NewWalletRepository(pool *pgxpool.Pool) WalletRepository {
	return &walletRepository{pool: pool}
}

func (r *walletRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Wallet, error) {
	query := `SELECT id, balance FROM wallets WHERE id = $1`

	var wallet model.Wallet
	err := r.pool.QueryRow(ctx, query, id).Scan(&wallet.ID, &wallet.Balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrWalletNotFound
		}
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	return &wallet, nil
}

// ProcessOperation handles deposit/withdraw with row-level locking for concurrency safety
func (r *walletRepository) ProcessOperation(ctx context.Context, req *model.WalletOperationRequest) (*model.Wallet, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var wallet model.Wallet
	lockQuery := `SELECT id, balance FROM wallets WHERE id = $1 FOR UPDATE`

	err = tx.QueryRow(ctx, lockQuery, req.WalletID).Scan(&wallet.ID, &wallet.Balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			if req.OperationType == model.OperationWithdraw {
				return nil, model.ErrWalletNotFound
			}
			wallet, err = r.createWalletInTx(ctx, tx, req.WalletID, req.Amount)
			if err != nil {
				return nil, err
			}
			if err := tx.Commit(ctx); err != nil {
				return nil, fmt.Errorf("failed to commit transaction: %w", err)
			}
			return &wallet, nil
		}
		return nil, fmt.Errorf("failed to lock wallet: %w", err)
	}

	var newBalance int64
	switch req.OperationType {
	case model.OperationDeposit:
		newBalance = wallet.Balance + req.Amount
	case model.OperationWithdraw:
		newBalance = wallet.Balance - req.Amount
		if newBalance < 0 {
			return nil, model.ErrInsufficientFunds
		}
	default:
		return nil, model.ErrInvalidOperation
	}

	updateQuery := `UPDATE wallets SET balance = $1, updated_at = NOW() WHERE id = $2 RETURNING id, balance`
	err = tx.QueryRow(ctx, updateQuery, newBalance, req.WalletID).Scan(&wallet.ID, &wallet.Balance)
	if err != nil {
		return nil, fmt.Errorf("failed to update wallet: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &wallet, nil
}

func (r *walletRepository) createWalletInTx(ctx context.Context, tx pgx.Tx, id uuid.UUID, amount int64) (model.Wallet, error) {
	query := `
		INSERT INTO wallets (id, balance, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		ON CONFLICT (id) DO UPDATE SET
			balance = wallets.balance + EXCLUDED.balance,
			updated_at = NOW()
		RETURNING id, balance
	`

	var wallet model.Wallet
	err := tx.QueryRow(ctx, query, id, amount).Scan(&wallet.ID, &wallet.Balance)
	if err != nil {
		return wallet, fmt.Errorf("failed to create wallet: %w", err)
	}

	return wallet, nil
}
