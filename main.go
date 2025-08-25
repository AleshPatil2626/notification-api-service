package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis/v8"

	"github.com/lm-Alesh-Patil/notification-api-service/config"
	"github.com/lm-Alesh-Patil/notification-api-service/notification/repository"
	"github.com/lm-Alesh-Patil/notification-api-service/notification/service"
	"github.com/lm-Alesh-Patil/notification-api-service/routes"

	"github.com/lm-Alesh-Patil/notification-api-service/utils"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.DB.MySQL.Username,
		cfg.DB.MySQL.Password,
		cfg.DB.MySQL.Host,
		cfg.DB.MySQL.Port,
		cfg.DB.MySQL.Database,
	))
	if err != nil {
		log.Fatalf("failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.DB.Redis.Host, cfg.DB.Redis.Port),
		Password: cfg.DB.Redis.Password,
		DB:       cfg.DB.Redis.DB,
	})
	defer rdb.Close()

	queueKey := "notification_queue"
	repo := repository.NewNotificationRepository(db, rdb, queueKey)
	emailCfg := utils.EmailConfig{
		SMTPHost:    cfg.Notification.SMTPHost,
		SMTPPort:    cfg.Notification.SMTPPort,
		Username:    cfg.Notification.Username,
		Password:    cfg.Notification.Password,
		SenderEmail: cfg.Notification.SenderEmail,
	}

	notifService := service.NewNotificationService(repo, emailCfg)

	// Start worker in background
	go notifService.ProcessQueue(context.Background())

	// Setup router
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	routes.RegisterRoutes(router, db, rdb, notifService)

	addr := fmt.Sprintf("%s:%d", cfg.Connection.HTTP.Host, cfg.Connection.HTTP.Port)
	fmt.Printf("Notification service running at http://%s\n", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
