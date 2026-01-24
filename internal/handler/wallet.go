package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/simple-wallet/internal/model"
	"github.com/simple-wallet/internal/service"
)

type WalletHandler struct {
	service service.WalletService
}

func NewWalletHandler(service service.WalletService) *WalletHandler {
	return &WalletHandler{service: service}
}

func (h *WalletHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/api/v1/wallet", h.ProcessOperation).Methods(http.MethodPost)
	r.HandleFunc("/api/v1/wallets/{walletId}", h.GetBalance).Methods(http.MethodGet)
}

// ProcessOperation handles POST /api/v1/wallet
func (h *WalletHandler) ProcessOperation(w http.ResponseWriter, r *http.Request) {
	var req model.WalletOperationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	wallet, err := h.service.ProcessOperation(r.Context(), &req)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, model.WalletResponse{
		ID:      wallet.ID,
		Balance: wallet.Balance,
	})
}

// GetBalance handles GET /api/v1/wallets/{walletId}
func (h *WalletHandler) GetBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	walletIDStr := vars["walletId"]

	walletID, err := uuid.Parse(walletIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid wallet ID format")
		return
	}

	wallet, err := h.service.GetBalance(r.Context(), walletID)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	respondWithJSON(w, http.StatusOK, model.WalletResponse{
		ID:      wallet.ID,
		Balance: wallet.Balance,
	})
}

func handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, model.ErrWalletNotFound):
		respondWithError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, model.ErrInsufficientFunds):
		respondWithError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, model.ErrInvalidOperation):
		respondWithError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, model.ErrInvalidAmount):
		respondWithError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, model.ErrInvalidWalletID):
		respondWithError(w, http.StatusBadRequest, err.Error())
	default:
		respondWithError(w, http.StatusInternalServerError, "internal server error")
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, model.ErrorResponse{Error: message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
