package mail

import (
	"crypto/tls"
	"fmt"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

// loginAuth implements the LOGIN SMTP auth mechanism (required by some servers like Beget).
type loginAuth struct{ user, pass string }

func newAuth(user, pass string) smtp.Auth { return &loginAuth{user, pass} }

func (a *loginAuth) Start(_ *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", nil, nil
}
func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if !more {
		return nil, nil
	}
	switch string(fromServer) {
	case "Username:":
		return []byte(a.user), nil
	case "Password:":
		return []byte(a.pass), nil
	default:
		return nil, fmt.Errorf("unexpected SMTP challenge: %s", fromServer)
	}
}

// extractEnvelopeAddress извлекает чистый email (без имени и скобок) для команд SMTP (MAIL FROM / RCPT TO).
func extractEnvelopeAddress(raw string) string {
	raw = strings.TrimSpace(raw)
	if addr, err := mail.ParseAddress(raw); err == nil && addr.Address != "" {
		return addr.Address
	}
	if start := strings.Index(raw, "<"); start >= 0 {
		if end := strings.Index(raw[start:], ">"); end > 0 {
			return strings.TrimSpace(raw[start+1 : start+end])
		}
	}
	return strings.Trim(raw, "<> \t\r\n")
}

// formatFromHeader форматирует заголовок From: с RFC 2047 кодированием имени (для поддержки кириллицы).
func formatFromHeader(raw string) string {
	raw = strings.TrimSpace(raw)
	if addr, err := mail.ParseAddress(raw); err == nil {
		if addr.Name != "" {
			return fmt.Sprintf("%s <%s>", mime.QEncoding.Encode("utf-8", addr.Name), addr.Address)
		}
		return addr.Address
	}
	return raw
}

// formatSubjectHeader кодирует тему письма по RFC 2047.
func formatSubjectHeader(subj string) string {
	return mime.QEncoding.Encode("utf-8", subj)
}

func mailBytes(fromAddr, toAddr, subject, body string, html bool) []byte {
	ctype := "text/plain; charset=UTF-8"
	if html {
		ctype = "text/html; charset=UTF-8"
	}
	hdr := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: %s\r\n\r\n",
		formatFromHeader(fromAddr), toAddr, formatSubjectHeader(subject), ctype)
	return []byte(hdr + body)
}

// SendText — простое письмо: тема + текст.
func SendText(host string, port int, user, password, fromAddr, toAddr, subject, body string, useTLS, tlsInsecure bool) error {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	msg := mailBytes(fromAddr, toAddr, subject, body, false)
	envFrom := extractEnvelopeAddress(fromAddr)
	envTo := extractEnvelopeAddress(toAddr)

	if useTLS || port == 465 {
		return sendSMTPS(addr, user, password, envFrom, []string{envTo}, msg, tlsInsecure)
	}
	return smtp.SendMail(addr, newAuth(user, password), envFrom, []string{envTo}, msg)
}

// SendHTML — письмо с телом text/html (как Nest + hbs).
func SendHTML(host string, port int, user, password, fromAddr, toAddr, subject, htmlBody string, useTLS, tlsInsecure bool) error {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	msg := mailBytes(fromAddr, toAddr, subject, htmlBody, true)
	envFrom := extractEnvelopeAddress(fromAddr)
	envTo := extractEnvelopeAddress(toAddr)

	if useTLS || port == 465 {
		return sendSMTPS(addr, user, password, envFrom, []string{envTo}, msg, tlsInsecure)
	}
	return smtp.SendMail(addr, newAuth(user, password), envFrom, []string{envTo}, msg)
}

func sendSMTPS(addr, user, pass, from string, to []string, msg []byte, tlsInsecure bool) error {
	host, _, _ := net.SplitHostPort(addr)
	rawConn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return fmt.Errorf("smtp dial %s: %w", addr, err)
	}
	conn := tls.Client(rawConn, &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: tlsInsecure,
		MinVersion:         tls.VersionTLS12,
	})
	if err := conn.SetDeadline(time.Now().Add(15 * time.Second)); err != nil {
		rawConn.Close()
		return err
	}
	if err := conn.Handshake(); err != nil {
		rawConn.Close()
		return fmt.Errorf("smtp tls handshake: %w", err)
	}
	defer conn.Close()
	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer c.Close()
	if ok, _ := c.Extension("AUTH"); ok {
		if err := c.Auth(newAuth(user, pass)); err != nil {
			return err
		}
	}
	if err := c.Mail(extractEnvelopeAddress(from)); err != nil {
		return err
	}
	for _, a := range to {
		if err := c.Rcpt(extractEnvelopeAddress(a)); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		return err
	}
	return w.Close()
}

// SendTextSmart — если SMTP_SECURE=true или порт 465, шлём через TLS. tlsInsecure — SMTP_TLS_INSECURE (как Nest rejectUnauthorized: false).
func SendTextSmart(host string, port int, user, password, fromAddr, toAddr, subject, body string, smtpSecure, tlsInsecure bool) error {
	if port == 465 {
		return SendText(host, port, user, password, fromAddr, toAddr, subject, body, true, tlsInsecure)
	}
	if smtpSecure {
		return sendSTARTTLS(host, port, user, password, fromAddr, toAddr, subject, body, false, tlsInsecure)
	}
	return SendText(host, port, user, password, fromAddr, toAddr, subject, body, false, tlsInsecure)
}

// SendHTMLSmart — то же для HTML (verify-email / forgot-password).
func SendHTMLSmart(host string, port int, user, password, fromAddr, toAddr, subject, htmlBody string, smtpSecure, tlsInsecure bool) error {
	if port == 465 {
		return SendHTML(host, port, user, password, fromAddr, toAddr, subject, htmlBody, true, tlsInsecure)
	}
	if smtpSecure {
		return sendSTARTTLS(host, port, user, password, fromAddr, toAddr, subject, htmlBody, true, tlsInsecure)
	}
	return SendHTML(host, port, user, password, fromAddr, toAddr, subject, htmlBody, false, tlsInsecure)
}

func sendSTARTTLS(host string, port int, user, password, fromAddr, toAddr, subject, body string, html, tlsInsecure bool) error {
	addr := net.JoinHostPort(host, strconv.Itoa(port))
	msg := mailBytes(fromAddr, toAddr, subject, body, html)
	envFrom := extractEnvelopeAddress(fromAddr)
	envTo := extractEnvelopeAddress(toAddr)

	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return fmt.Errorf("smtp dial %s: %w", addr, err)
	}
	_ = conn.SetDeadline(time.Now().Add(20 * time.Second))
	defer conn.Close()
	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer c.Close()
	if err := c.Hello("localhost"); err != nil {
		return err
	}
	if ok, _ := c.Extension("STARTTLS"); ok {
		cfg := &tls.Config{
			ServerName:         host,
			InsecureSkipVerify: tlsInsecure,
			MinVersion:         tls.VersionTLS12,
		}
		if err := c.StartTLS(cfg); err != nil {
			return err
		}
	}
	if ok, _ := c.Extension("AUTH"); ok {
		if err := c.Auth(newAuth(user, password)); err != nil {
			return err
		}
	}
	if err := c.Mail(envFrom); err != nil {
		return err
	}
	if err := c.Rcpt(envTo); err != nil {
		return err
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return c.Quit()
}
