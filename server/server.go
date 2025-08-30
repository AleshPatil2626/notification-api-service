package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-redis/redis/v8"
	"github.com/lm-Alesh-Patil/notification-api-service/config"
	"github.com/lm-Alesh-Patil/notification-api-service/notification/repository"
	"github.com/lm-Alesh-Patil/notification-api-service/notification/service"
	"github.com/lm-Alesh-Patil/notification-api-service/processor"
	"github.com/lm-Alesh-Patil/notification-api-service/routes"
	"github.com/robfig/cron/v3"

	_ "github.com/go-sql-driver/mysql"
)

type Server struct {
	Config *config.Config
	DB     *sql.DB
	Redis  *redis.Client
}

func (s *Server) SetupMysqlDatabase() error {
	dbCfg := s.Config.DB.MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		dbCfg.Username,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.Database,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("mysql connect error: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("mysql ping error: %w", err)
	}

	s.DB = db
	return nil
}

func (s *Server) SetupRedis() error {
	redisCfg := s.Config.DB.Redis

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port),
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})

	s.Redis = client
	return nil
}

func (s *Server) Setup() error {
	if err := s.SetupMysqlDatabase(); err != nil {
		log.Println("MYSQL setup error", err)
		return err
	}
	if err := s.SetupRedis(); err != nil {
		log.Println("REDIS setup error", err)
		return err
	}
	return nil
}

// StartCronJobs runs scheduled jobs
func (s *Server) startCronJobs() {
	c := cron.New()

	// Example: run every minute
	_, err := c.AddFunc("@every 1m", func() {
		log.Println("Cron Job: Checking scheduled tasks...")
	})
	if err != nil {
		log.Fatalf("failed to add cron job: %v", err)
	}

	c.Start()
	log.Println("Cron scheduler started...")
}

func (s *Server) Start() error {
	// Setup MySQL + Redis
	if err := s.Setup(); err != nil {
		return err
	}

	// --- Create dependencies ---
	notifRepo := repository.NewNotificationRepository(s.DB, s.Redis, "notification_queue")
	notifService := service.NewNotificationService(notifRepo)

	// Start worker
	worker := processor.NewNotificationWorker(notifRepo, notifService)
	go worker.Start(context.Background())

	// Start cron jobs
	s.startCronJobs()

	// Setup HTTP server
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(30 * time.Second))

	// Register routes with NotificationService
	routes.RegisterRoutes(router, s.DB, s.Redis, notifService)
	httpServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.Config.Connection.HTTP.Host, s.Config.Connection.HTTP.Port),
		ReadTimeout:  time.Duration(s.Config.Connection.HTTP.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.Config.Connection.HTTP.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(s.Config.Connection.HTTP.IdleTimeout) * time.Second,
		Handler:      router,
	}

	fmt.Printf("Server starting at %s:%d...\n",
		s.Config.Connection.HTTP.Host,
		s.Config.Connection.HTTP.Port,
	)
	return httpServer.ListenAndServe()
}
