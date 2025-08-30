package models

type Notification struct {
	ID       int64             `json:"id"`
	OrgID    int64             `json:"org_id"` //link to org for SMTP + templates
	Type     string            `json:"type"`   // e.g. "email", "sms"
	To       string            `json:"to"`     // email address, phone, etc.
	Subject  string            `json:"subject"`
	Body     string            `json:"body"`
	Template string            `json:"template"`
	Data     map[string]string `json:"data"`
}

type User struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	OrgID int64  `json:"org_id"`
}

type SMTPConfig struct {
	Host        string
	Port        int
	Username    string
	Password    string
	SenderEmail string
}

type EmailTemplate struct {
	Name    string
	Subject string
	Body    string
}
