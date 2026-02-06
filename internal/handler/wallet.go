package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"wallet/internal/service"
)

type Handler struct {
	service *service.Service
}

func New(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/api/v1/wallet", h.handleWallet)
	r.Get("/api/v1/wallets/{id}", h.getBalance)

	return r
}

func (h *Handler) handleWallet(w http.ResponseWriter, r *http.Request) {
	var req struct {
		WalletID      uuid.UUID `json:"walletId"`
		OperationType string    `json:"operationType"`
		Amount        int64     `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.service.UpdateBalance(
		r.Context(),
		req.WalletID,
		req.Amount,
		req.OperationType,
	)

	if err != nil {
		if err.Error() == "insufficient funds" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getBalance(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	balance, err := h.service.GetBalance(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]int64{
		"balance": balance,
	})
}

