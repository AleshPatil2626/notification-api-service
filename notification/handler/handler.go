package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/lm-Alesh-Patil/notification-api-service/notification/models"
	"github.com/lm-Alesh-Patil/notification-api-service/notification/service"
)

// Handler struct for dependency injection
type NotificationHandler struct {
	Service *service.NotificationService
}

// NewNotificationHandler returns a new handler
func NewNotificationHandler(svc *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{Service: svc}
}

// CreateNotification handles POST /notifications
func (h *NotificationHandler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	var notif models.Notification
	if err := json.NewDecoder(r.Body).Decode(&notif); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.Service.EnqueueNotification(context.Background(), &notif); err != nil {
		http.Error(w, "failed to enqueue notification", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"message": "Notification enqueued successfully",})
}

// ProcessQueue handles POST /notifications/process
func (h *NotificationHandler) ProcessQueue(w http.ResponseWriter, r *http.Request) {
	go h.Service.ProcessQueue(context.Background())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Queue processing started"))
}

// GetNotificationByID handles GET /notifications/{id}
func (h *NotificationHandler) GetNotificationByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	notif, err := h.Service.GetNotificationByID(context.Background(), id)
	if err != nil {
		http.Error(w, "notification not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notif)
}
