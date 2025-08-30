package routes

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"github.com/lm-Alesh-Patil/notification-api-service/notification/handler"
	"github.com/lm-Alesh-Patil/notification-api-service/notification/service"
)

// RegisterRoutes wires handlers to routes
func RegisterRoutes(router *chi.Mux, db *sql.DB, rdb *redis.Client, notifService *service.NotificationService) {
	handler := handler.NewNotificationHandler(notifService)
	router.Post("/notifications", handler.CreateNotification)
	router.Post("/notifications/process", handler.ProcessQueue)
	router.Get("/notifications/{id}", handler.GetNotificationByID)
}
