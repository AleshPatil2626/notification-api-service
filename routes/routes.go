package routes

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/lm-Alesh-Patil/notification-api-service/notification/service"
)

func RegisterRoutes(router *chi.Mux, db *sql.DB, rdb *redis.Client, notifService *service.NotificationService) {
	// Example endpoint to trigger processing manually
	router.Post("/process/queue", func(w http.ResponseWriter, r *http.Request) {
		go notifService.ProcessQueue(context.Background())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Queue processing started"))
	})
}
