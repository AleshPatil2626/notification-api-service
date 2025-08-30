package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/lm-Alesh-Patil/notification-api-service/notification/models"
)

type NotificationRepository struct {
	db          *sql.DB
	redisClient *redis.Client
	queueKey    string
}

func NewNotificationRepository(db *sql.DB, redisClient *redis.Client, queueKey string) *NotificationRepository {
	return &NotificationRepository{
		db:          db,
		redisClient: redisClient,
		queueKey:    queueKey,
	}
}

// Fetch email template by org and name
func (r *NotificationRepository) GetTemplate(ctx context.Context, orgID int64, name string) (*models.EmailTemplate, error) {
	var tmpl models.EmailTemplate
	query := `SELECT subject, body FROM email_templates WHERE org_id = ? AND name = ? LIMIT 1`
	err := r.db.QueryRowContext(ctx, query, orgID, name).Scan(&tmpl.Subject, &tmpl.Body)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("template %s not found for org %d", name, orgID)
		}
		return nil, err
	}
	tmpl.Name = name
	return &tmpl, nil
}

// Save notification logs
func (r *NotificationRepository) LogStatus(ctx context.Context, orgID int64, templateName, recipient, status, message string) error {
	query := `INSERT INTO notification_logs (org_id, template_name, recipient, status, message) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, orgID, templateName, recipient, status, message)
	return err
}

// Get SMTP config from organizations
func (r *NotificationRepository) GetSMTPConfig(ctx context.Context, orgID int64) (*models.SMTPConfig, error) {
	query := `SELECT smtp_host, smtp_port, smtp_username, smtp_password, sender_email FROM organizations WHERE id = ? LIMIT 1`
	var cfg models.SMTPConfig
	err := r.db.QueryRowContext(ctx, query, orgID).Scan(
		&cfg.Host,
		&cfg.Port,
		&cfg.Username,
		&cfg.Password,
		&cfg.SenderEmail,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no SMTP config found for org %d", orgID)
		}
		return nil, err
	}
	return &cfg, nil
}

// GetNotificationByID fetches notification from DB
func (r *NotificationRepository) GetNotificationByID(ctx context.Context, id string) (*models.Notification, error) {
	query := "SELECT id, type, to, subject, body, template FROM notifications WHERE id = ?"
	row := r.db.QueryRowContext(ctx, query, id)

	var notif models.Notification
	if err := row.Scan(&notif.ID, &notif.Type, &notif.To, &notif.Subject, &notif.Body, &notif.Template); err != nil {
		return nil, fmt.Errorf("failed to fetch notification: %w", err)
	}
	return &notif, nil
}
