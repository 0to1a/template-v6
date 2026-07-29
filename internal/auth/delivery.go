package auth

import (
	"bytes"
	"context"
	_ "embed"
	"html/template"
	"log"
)

// Implementations must never log the code.
type LoginCodeSender interface {
	SendLoginCode(ctx context.Context, email, code string) error
}

// Lets the server run before real delivery is wired up
type NoopLoginCodeSender struct{}

func (NoopLoginCodeSender) SendLoginCode(context.Context, string, string) error {
	return nil
}

//go:embed otp.html
var otpTemplateSource string

var otpTemplate = template.Must(template.New("otp").Parse(otpTemplateSource))

// Decouples this package from SMTP specifics
type EmailSender interface {
	Send(ctx context.Context, to, subject, html string) error
}

type EmailLoginCodeSender struct {
	sender EmailSender
}

func NewEmailLoginCodeSender(sender EmailSender) *EmailLoginCodeSender {
	return &EmailLoginCodeSender{sender: sender}
}

// Logs failures only; caller discards to avoid oracle
func (e *EmailLoginCodeSender) SendLoginCode(ctx context.Context, email, code string) error {
	var body bytes.Buffer
	if err := otpTemplate.Execute(&body, struct{ Code string }{Code: code}); err != nil {
		log.Printf("auth: rendering login code email failed: %v", err)
		return err
	}

	if err := e.sender.Send(ctx, email, "Your login code", body.String()); err != nil {
		log.Printf("auth: sending login code email failed: %v", err)
		return err
	}
	return nil
}
