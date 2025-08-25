package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/lm-Alesh-Patil/docker-api-service/user_management/models"
)

type NotificationRepository struct {
	db          *sql.DB
	redisClient *redis.Client
	queueKey    string
}

func NewNotificationRepository(db *sql.DB, rdb *redis.Client, queueKey string) *NotificationRepository {
	return &NotificationRepository{
		db:          db,
		redisClient: rdb,
		queueKey:    queueKey,
	}
}

func (r *NotificationRepository) PopUser(ctx context.Context) (*models.User, error) {
	data, err := r.redisClient.LPop(ctx, r.queueKey).Result()
	if err != nil {
		return nil, err
	}
	var user models.User
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil, fmt.Errorf("failed to parse user data: %w", err)
	}
	return &user, nil
}

func (r *NotificationRepository) RetryUser(ctx context.Context, user *models.User) error {
	data, _ := json.Marshal(user)
	return r.redisClient.RPush(ctx, r.queueKey, data).Err()
}

func (r *NotificationRepository) QueueLength(ctx context.Context) (int64, error) {
	length, err := r.redisClient.LLen(ctx, r.queueKey).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get queue length: %w", err)
	}
	return length, nil
}

func (r *NotificationRepository) LogStatus(ctx context.Context, userID int64, status, message string) error {
	query := "INSERT INTO notification_logs (user_id, status, message) VALUES (?, ?, ?)"
	_, err := r.db.ExecContext(ctx, query, userID, status, message)
	return err
}
