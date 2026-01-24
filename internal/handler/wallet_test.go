package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/simple-wallet/internal/mocks"
	"github.com/simple-wallet/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWalletHandler_GetBalance(t *testing.T) {
	t.Run("returns balance for existing wallet", func(t *testing.T) {
		svc := mocks.NewMockWalletService(t)
		h := NewWalletHandler(svc)

		walletID := uuid.New()
		expected := &model.Wallet{ID: walletID, Balance: 1000}

		svc.EXPECT().
			GetBalance(mock.Anything, walletID).
			Return(expected, nil)

		router := mux.NewRouter()
		h.RegisterRoutes(router)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/"+walletID.String(), nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp model.WalletResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, walletID, resp.ID)
		assert.Equal(t, int64(1000), resp.Balance)
	})

	t.Run("returns 404 for non-existing wallet", func(t *testing.T) {
		svc := mocks.NewMockWalletService(t)
		h := NewWalletHandler(svc)

		walletID := uuid.New()

		svc.EXPECT().
			GetBalance(mock.Anything, walletID).
			Return(nil, model.ErrWalletNotFound)

		router := mux.NewRouter()
		h.RegisterRoutes(router)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/"+walletID.String(), nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("returns 400 for invalid UUID", func(t *testing.T) {
		svc := mocks.NewMockWalletService(t)
		h := NewWalletHandler(svc)

		router := mux.NewRouter()
		h.RegisterRoutes(router)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets/invalid-uuid", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestWalletHandler_ProcessOperation(t *testing.T) {
	t.Run("successful deposit", func(t *testing.T) {
		svc := mocks.NewMockWalletService(t)
		h := NewWalletHandler(svc)

		walletID := uuid.New()
		reqBody := model.WalletOperationRequest{
			WalletID:      walletID,
			OperationType: model.OperationDeposit,
			Amount:        1000,
		}
		expected := &model.Wallet{ID: walletID, Balance: 1000}

		svc.EXPECT().
			ProcessOperation(mock.Anything, mock.MatchedBy(func(r *model.WalletOperationRequest) bool {
				return r.WalletID == walletID &&
					r.OperationType == model.OperationDeposit &&
					r.Amount == 1000
			})).
			Return(expected, nil)

		router := mux.NewRouter()
		h.RegisterRoutes(router)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp model.WalletResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, walletID, resp.ID)
		assert.Equal(t, int64(1000), resp.Balance)
	})

	t.Run("successful withdraw", func(t *testing.T) {
		svc := mocks.NewMockWalletService(t)
		h := NewWalletHandler(svc)

		walletID := uuid.New()
		reqBody := model.WalletOperationRequest{
			WalletID:      walletID,
			OperationType: model.OperationWithdraw,
			Amount:        500,
		}
		expected := &model.Wallet{ID: walletID, Balance: 500}

		svc.EXPECT().
			ProcessOperation(mock.Anything, mock.MatchedBy(func(r *model.WalletOperationRequest) bool {
				return r.WalletID == walletID &&
					r.OperationType == model.OperationWithdraw &&
					r.Amount == 500
			})).
			Return(expected, nil)

		router := mux.NewRouter()
		h.RegisterRoutes(router)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("returns 400 on insufficient funds", func(t *testing.T) {
		svc := mocks.NewMockWalletService(t)
		h := NewWalletHandler(svc)

		walletID := uuid.New()
		reqBody := model.WalletOperationRequest{
			WalletID:      walletID,
			OperationType: model.OperationWithdraw,
			Amount:        10000,
		}

		svc.EXPECT().
			ProcessOperation(mock.Anything, mock.Anything).
			Return(nil, model.ErrInsufficientFunds)

		router := mux.NewRouter()
		h.RegisterRoutes(router)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp model.ErrorResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "insufficient funds", resp.Error)
	})

	t.Run("returns 400 on invalid JSON", func(t *testing.T) {
		svc := mocks.NewMockWalletService(t)
		h := NewWalletHandler(svc)

		router := mux.NewRouter()
		h.RegisterRoutes(router)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}
