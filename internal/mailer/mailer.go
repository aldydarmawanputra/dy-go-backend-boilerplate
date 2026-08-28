package mailer

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"

	"go-backend-boilerplate/internal/config"
)

type Mailer interface {
	Send(ctx context.Context, to, subject, htmlBody string) error
}

func New(cfg *config.Config) Mailer {
	if cfg.SMTPHost == "" {
		return noop{}
	}
	return &smtpMailer{
		host:     cfg.SMTPHost,
		port:     cfg.SMTPPort,
		user:     cfg.SMTPUser,
		password: cfg.SMTPPassword,
		from:     cfg.SMTPFrom,
	}
}

type smtpMailer struct {
	host, port, user, password, from string
}

func (m *smtpMailer) Send(_ context.Context, to, subject, htmlBody string) error {
	var b strings.Builder
	fmt.Fprintf(&b, "From: %s\r\n", m.from)
	fmt.Fprintf(&b, "To: %s\r\n", to)
	fmt.Fprintf(&b, "Subject: %s\r\n", subject)
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	b.WriteString("\r\n")
	b.WriteString(htmlBody)

	var auth smtp.Auth
	if m.user != "" {
		// net/smtp negotiates STARTTLS automatically when the server supports it
		// (e.g. Gmail on :587). Local Mailpit on :1025 needs no auth/TLS.
		auth = smtp.PlainAuth("", m.user, m.password, m.host)
	}
	return smtp.SendMail(m.host+":"+m.port, auth, m.from, []string{to}, []byte(b.String()))
}

type noop struct{}

func (noop) Send(context.Context, string, string, string) error { return nil }
