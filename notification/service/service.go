package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/lm-Alesh-Patil/notification-api-service/notification/repository"
	"github.com/lm-Alesh-Patil/notification-api-service/utils"
)

type NotificationService struct {
	repo     *repository.NotificationRepository
	emailCfg utils.EmailConfig
}

func NewNotificationService(repo *repository.NotificationRepository, emailCfg utils.EmailConfig) *NotificationService {
	return &NotificationService{
		repo:     repo,
		emailCfg: emailCfg,
	}
}

// ProcessQueue continuously polls Redis and sends emails
func (s *NotificationService) ProcessQueue(ctx context.Context) {
	log.Println("NotificationService: Started processing Redis queue...")

	length, _ := s.repo.QueueLength(ctx)
	log.Printf("NotificationService: Queue length=%d\n", length)

	for {
		// Pop next user from Redis queue
		user, err := s.repo.PopUser(ctx)
		if err != nil {
			if err.Error() == "queue empty" {
				log.Println("NotificationService: Queue is empty, waiting 2 seconds...")
				time.Sleep(2 * time.Second)
				continue
			}
			log.Printf("NotificationService: Error popping user from queue: %v\n", err)
			time.Sleep(2 * time.Second)
			continue
		}

		log.Printf("NotificationService: Processing user ID=%d, Email=%s\n", user.ID, user.Email)

		subject := "Welcome to AP Solutions Pvt. Ltd, Pune!"
		body := fmt.Sprintf(`
			Hi %s,

			Welcome to AP Solutions Pvt. Ltd, Pune! 🎉

			We are excited to have you on board. Your registration has been successfully completed.

			Here are some key points to get started:
			- Explore our services and solutions tailored for you.
			- Stay updated with our latest news and offers.
			- Feel free to reach out to our support team anytime at support@apsolutions.com.

			Once again, welcome to the AP Solutions family!

			Best regards,
			The AP Solutions Team
			Pune, India
			`, user.Name)

		err = utils.SendEmail(s.emailCfg, user.Email, subject, body)
		if err != nil {
			log.Printf("NotificationService: Failed to send email to %s: %v. Retrying...\n", user.Email, err)
			s.repo.RetryUser(ctx, user)
			s.repo.LogStatus(ctx, user.ID, "FAILED", err.Error())
			continue
		}

		// Email sent successfully
		log.Printf("NotificationService: Email sent successfully to %s\n", user.Email)
		s.repo.LogStatus(ctx, user.ID, "SUCCESS", "Email sent successfully")
	}
}
