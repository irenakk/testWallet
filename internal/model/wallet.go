package model

import (
	"errors"

	"github.com/google/uuid"
)

type OperationType string

const (
	OperationDeposit  OperationType = "DEPOSIT"
	OperationWithdraw OperationType = "WITHDRAW"
)

var (
	ErrWalletNotFound     = errors.New("wallet not found")
	ErrInsufficientFunds  = errors.New("insufficient funds")
	ErrInvalidOperation   = errors.New("invalid operation type")
	ErrInvalidAmount      = errors.New("amount must be positive")
	ErrInvalidWalletID    = errors.New("invalid wallet ID")
)

type Wallet struct {
	ID      uuid.UUID `json:"id"`
	Balance int64     `json:"balance"`
}

type WalletOperationRequest struct {
	WalletID      uuid.UUID     `json:"walletId"`
	OperationType OperationType `json:"operationType"`
	Amount        int64         `json:"amount"`
}

type WalletResponse struct {
	ID      uuid.UUID `json:"id"`
	Balance int64     `json:"balance"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (r *WalletOperationRequest) Validate() error {
	if r.WalletID == uuid.Nil {
		return ErrInvalidWalletID
	}
	if r.Amount <= 0 {
		return ErrInvalidAmount
	}
	if r.OperationType != OperationDeposit && r.OperationType != OperationWithdraw {
		return ErrInvalidOperation
	}
	return nil
}
