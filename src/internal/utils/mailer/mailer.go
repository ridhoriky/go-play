package mailer

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"github.com/rs/zerolog"
)

type Mailer struct {
	log         zerolog.Logger
	host        string
	port        int
	username    string
	password    string
	senderName  string
	senderEmail string
}

func NewMailer(log *zerolog.Logger, host string, port int, username, password, senderName, senderEmail string) *Mailer {
	return &Mailer{
		log:         *log,
		host:        host,
		port:        port,
		username:    username,
		password:    password,
		senderName:  senderName,
		senderEmail: senderEmail,
	}
}

// SendEmail sends a real email if credentials are provided, otherwise logs a simulation.
func (m *Mailer) SendEmail(ctx context.Context, to, subject, body string) error {
	if m.host == "" || m.username == "" || m.password == "" {
		m.log.Info().
			Str("to", to).
			Str("subject", subject).
			Msg("SMTP credentials not configured. Logging simulated email.")
		fmt.Printf("\n--- SIMULATED EMAIL ---\nTo: %s\nSubject: %s\nBody:\n%s\n------------------------\n\n", to, subject, body)
		return nil
	}

	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", m.senderName, m.senderEmail)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	var messageSb52 strings.Builder
	for k, v := range headers {
		fmt.Fprintf(&messageSb52, "%s: %s\r\n", k, v)
	}
	message += messageSb52.String()
	message += "\r\n" + body

	auth := smtp.PlainAuth("", m.username, m.password, m.host)
	addr := fmt.Sprintf("%s:%d", m.host, m.port)

	if m.port == 465 {
		return m.sendEmailSSL(ctx, addr, auth, to, message)
	}

	return smtp.SendMail(addr, auth, m.senderEmail, []string{to}, []byte(message))
}

func (m *Mailer) sendEmailSSL(ctx context.Context, addr string, auth smtp.Auth, to string, message string) error {
	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         m.host,
	}

	dialer := &tls.Dialer{
		Config: tlsconfig,
	}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return err
	}
	defer func() {
		_ = conn.Close()
	}()

	client, err := smtp.NewClient(conn, m.host)
	if err != nil {
		return err
	}
	defer func() {
		_ = client.Close()
	}()

	if err = client.Auth(auth); err != nil {
		return err
	}

	if err = client.Mail(m.senderEmail); err != nil {
		return err
	}

	if err = client.Rcpt(to); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return client.Quit()
}
