package service

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"

	"github.com/lm-Alesh-Patil/notification-api-service/notification/models"
	"github.com/lm-Alesh-Patil/notification-api-service/notification/repository"
	"github.com/lm-Alesh-Patil/notification-api-service/utils"
)

type NotificationService struct {
	repo *repository.NotificationRepository
}

func NewNotificationService(repo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}
func (s *NotificationService) SendNotification(ctx context.Context, notif *models.Notification) error {
	// 1. Get org SMTP config
	org, err := s.repo.GetSMTPConfig(ctx, notif.OrgID)
	if err != nil {
		return fmt.Errorf("failed to get org config: %w", err)
	}

	// 2. Load template
	templateData, err := s.repo.GetTemplate(ctx, notif.OrgID, notif.Template)
	if err != nil {
		return fmt.Errorf("failed to get template: %w", err)
	}

	// 3. Render subject
	subject, err := renderTemplate(templateData.Subject, notif.Data)
	if err != nil {
		return fmt.Errorf("failed to render subject: %w", err)
	}

	// 4. Render body
	body, err := renderTemplate(templateData.Body, notif.Data)
	if err != nil {
		return fmt.Errorf("failed to render body: %w", err)
	}

	// 5. Prepare SMTP config
	emailCfg := utils.EmailConfig{
		SMTPHost: org.Host,
		SMTPPort: org.Port,
		Username: org.Username,
		Password: org.Password,
		Sender:   org.SenderEmail,
	}

	// 6. Send Email
	err = utils.SendEmail(emailCfg, notif.To, subject, body)
	if err != nil {
		s.repo.LogStatus(ctx, notif.OrgID, notif.Template, notif.To, "FAILED", err.Error())
		return err
	}

	// 7. Log success
	s.repo.LogStatus(ctx, notif.OrgID, notif.Template, notif.To, "SUCCESS", "Email sent successfully")
	log.Printf("Email sent successfully to %s\n", notif.To)
	return nil
}

// renderTemplate replaces placeholders using Go's template engine
func renderTemplate(tmpl string, data map[string]string) (string, error) {
	t, err := template.New("email").Parse(tmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// ProcessQueue pulls notifications from Redis queue and sends them
func (s *NotificationService) ProcessQueue(ctx context.Context) {
	for {
		notif, err := s.repo.PopNotification(ctx) // Pop notification, not user
		if err != nil {
			log.Println("Queue error:", err)
			return
		}

		// Send the notification
		if err := s.SendNotification(ctx, notif); err != nil {
			log.Printf("Failed to send notification to %s: %v\n", notif.To, err)
			// Retry logic
			_ = s.repo.RetryNotification(ctx, notif)
		} else {
			log.Printf("Notification sent successfully to %s\n", notif.To)
		}
	}
}

// EnqueueNotification pushes a notification into Redis queue
func (s *NotificationService) EnqueueNotification(ctx context.Context, notif *models.Notification) error {
	return s.repo.PushNotification(ctx, notif)
}

// GetNotificationByID fetches a notification by ID
func (s *NotificationService) GetNotificationByID(ctx context.Context, id string) (*models.Notification, error) {
	return s.repo.GetNotificationByID(ctx, id)
}
