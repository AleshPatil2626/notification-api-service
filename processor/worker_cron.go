package processor

import (
	"context"
	"log"
	"time"

	"github.com/lm-Alesh-Patil/notification-api-service/notification/repository"
	"github.com/lm-Alesh-Patil/notification-api-service/notification/service"
)

type NotificationWorker struct {
	repo    *repository.NotificationRepository
	service *service.NotificationService
}

func NewNotificationWorker(repo *repository.NotificationRepository, svc *service.NotificationService) *NotificationWorker {
	return &NotificationWorker{repo: repo, service: svc}
}

// Start runs a background worker that polls queue
func (worker *NotificationWorker) Start(ctx context.Context) {
	log.Println("NotificationWorker: Started polling queue...")

	for {
		notif, err := worker.repo.PopNotification(ctx) // now we pop a Notification, not just a User
		if err != nil {
			if err.Error() == "queue empty" {
				time.Sleep(2 * time.Second)
				continue
			}
			log.Printf("Error fetching from queue: %v\n", err)
			time.Sleep(2 * time.Second)
			continue
		}
		// Send notification via service
		err = worker.service.SendNotification(ctx, notif)
		if err != nil {
			log.Printf("NotificationWorker: Failed to send notif=%d, retrying...\n", notif.ID)
			worker.repo.RetryNotification(ctx, notif) // retry whole notif, not just user
		}
	}
}
