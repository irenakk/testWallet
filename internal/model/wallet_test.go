package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestWalletOperationRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     WalletOperationRequest
		wantErr error
	}{
		{
			name: "valid deposit",
			req: WalletOperationRequest{
				WalletID:      uuid.New(),
				OperationType: OperationDeposit,
				Amount:        1000,
			},
			wantErr: nil,
		},
		{
			name: "valid withdraw",
			req: WalletOperationRequest{
				WalletID:      uuid.New(),
				OperationType: OperationWithdraw,
				Amount:        500,
			},
			wantErr: nil,
		},
		{
			name: "invalid wallet ID",
			req: WalletOperationRequest{
				WalletID:      uuid.Nil,
				OperationType: OperationDeposit,
				Amount:        1000,
			},
			wantErr: ErrInvalidWalletID,
		},
		{
			name: "invalid operation type",
			req: WalletOperationRequest{
				WalletID:      uuid.New(),
				OperationType: "INVALID",
				Amount:        1000,
			},
			wantErr: ErrInvalidOperation,
		},
		{
			name: "zero amount",
			req: WalletOperationRequest{
				WalletID:      uuid.New(),
				OperationType: OperationDeposit,
				Amount:        0,
			},
			wantErr: ErrInvalidAmount,
		},
		{
			name: "negative amount",
			req: WalletOperationRequest{
				WalletID:      uuid.New(),
				OperationType: OperationDeposit,
				Amount:        -100,
			},
			wantErr: ErrInvalidAmount,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
