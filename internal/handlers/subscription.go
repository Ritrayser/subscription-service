package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"subscription-service/internal/models"
	"subscription-service/internal/service"
	"log/slog"
)

type SubscriptionHandler struct {
	svc *service.SubscriptionService
	log *slog.Logger
}

func NewSubscriptionHandler(svc *service.SubscriptionService, log *slog.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{svc: svc, log: log}
}

func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input models.CreateSubscriptionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	sub, err := h.svc.Create(r.Context(), &input)
	if err != nil {
		h.log.Error("Failed to create subscription", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sub)
}

func (h *SubscriptionHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	sub, err := h.svc.Get(r.Context(), id)
	if err != nil {
		h.log.Error("Failed to get subscription", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if sub == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(sub)
}

func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var input models.UpdateSubscriptionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	sub, err := h.svc.Update(r.Context(), id, &input)
	if err != nil {
		h.log.Error("Failed to update subscription", "error", err)
		if err.Error() == "subscription not found" {
			http.Error(w, "Not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}
	json.NewEncoder(w).Encode(sub)
}

func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.svc.Delete(r.Context(), id)
	if err != nil {
		h.log.Error("Failed to delete subscription", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 20
	}
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if offset < 0 {
		offset = 0
	}
	userID := r.URL.Query().Get("user_id")
	serviceName := r.URL.Query().Get("service_name")

	var uidPtr *string
	if userID != "" {
		uidPtr = &userID
	}
	var namePtr *string
	if serviceName != "" {
		namePtr = &serviceName
	}

	subs, err := h.svc.List(r.Context(), limit, offset, uidPtr, namePtr)
	if err != nil {
		h.log.Error("Failed to list subscriptions", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(subs)
}

func (h *SubscriptionHandler) SumByPeriod(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	if start == "" || end == "" {
		http.Error(w, "start and end are required", http.StatusBadRequest)
		return
	}
	userID := r.URL.Query().Get("user_id")
	serviceName := r.URL.Query().Get("service_name")

	var uidPtr *string
	if userID != "" {
		uidPtr = &userID
	}
	var namePtr *string
	if serviceName != "" {
		namePtr = &serviceName
	}

	sum, err := h.svc.SumByPeriod(r.Context(), start, end, uidPtr, namePtr)
	if err != nil {
		h.log.Error("Failed to calculate sum", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(map[string]int{"total": sum})
}