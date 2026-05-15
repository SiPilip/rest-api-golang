package helpers

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	mail "github.com/wneessen/go-mail"
)

type EmailData struct {
	To      string
	Subject string
	Body    string
}

func SendEmail(data EmailData) error {
	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")

	m := mail.NewMsg()
	if err := m.From(from); err != nil {
		return fmt.Errorf("fauled to set from %w", err)
	}
	if err := m.To(data.To); err != nil {
		return fmt.Errorf("failed to set to: %w", err)
	}

	m.Subject(data.Subject)
	m.SetBodyString(mail.TypeTextHTML, data.Body)

	c, err := mail.NewClient(host,
		mail.WithPort(port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(username),
		mail.WithPassword(password))

	if err != nil {
		return fmt.Errorf("failed to create mait client: %w", err)
	}

	if err := c.DialAndSend(m); err != nil {
		slog.Error("Fauled to send email", "error", err, "to", data.To)
		return err
	}

	slog.Info("Email sent", "to", data.To, "subject", data.Subject)
	return nil
}