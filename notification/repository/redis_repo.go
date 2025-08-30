package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/lm-Alesh-Patil/notification-api-service/notification/models"
)

func (r *NotificationRepository) PushNotification(ctx context.Context, notif *models.Notification) error {
	data, err := json.Marshal(notif)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}
	return r.redisClient.LPush(ctx, r.queueKey, data).Err()
}

func (r *NotificationRepository) PopNotification(ctx context.Context) (*models.Notification, error) {
	data, err := r.redisClient.LPop(ctx, r.queueKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("queue empty")
		}
		return nil, fmt.Errorf("failed to pop from queue: %w", err)
	}

	var notif models.Notification
	if err := json.Unmarshal([]byte(data), &notif); err != nil {
		return nil, fmt.Errorf("failed to parse notification: %w", err)
	}
	return &notif, nil
}

func (r *NotificationRepository) RetryNotification(ctx context.Context, notif *models.Notification) error {
	data, err := json.Marshal(notif)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}
	return r.redisClient.RPush(ctx, r.queueKey, data).Err()
}

func (r *NotificationRepository) QueueLength(ctx context.Context) (int64, error) {
	length, err := r.redisClient.LLen(ctx, r.queueKey).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get queue length: %w", err)
	}
	return length, nil
}
