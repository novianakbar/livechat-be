package email

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/novianakbar/livechat-be/internal/domain"
	"github.com/novianakbar/livechat-be/pkg/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type sendgridService struct {
	client   *sendgrid.Client
	config   *config.EmailConfig
	fromMail *mail.Email
}

// NewSendGridService creates new SendGrid email service
func NewSendGridService(cfg *config.EmailConfig) domain.EmailService {
	client := sendgrid.NewSendClient(cfg.SendGridAPIKey)
	fromEmail := mail.NewEmail(cfg.FromName, cfg.FromEmail)

	return &sendgridService{
		client:   client,
		config:   cfg,
		fromMail: fromEmail,
	}
}

// SendEmail sends a basic email
func (s *sendgridService) SendEmail(ctx context.Context, request *domain.SendEmailRequest) (*domain.EmailResponse, error) {
	if len(request.To) == 0 {
		return nil, fmt.Errorf("no recipients specified")
	}

	// Create personalization for each recipient
	message := mail.NewV3Mail()
	message.SetFrom(s.fromMail)
	message.Subject = request.Subject

	// Add content based on type
	if request.IsHTML {
		message.AddContent(mail.NewContent("text/html", request.Content))
	} else {
		message.AddContent(mail.NewContent("text/plain", request.Content))
	}

	// Add recipients
	personalization := mail.NewPersonalization()
	for _, email := range request.To {
		personalization.AddTos(mail.NewEmail("", email))
	}

	// Add template variables if provided
	if request.Variables != nil && len(request.Variables) > 0 {
		for key, value := range request.Variables {
			personalization.SetDynamicTemplateData(key, value)
		}
	}

	message.AddPersonalizations(personalization)

	// Send email
	response, err := s.client.Send(message)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return &domain.EmailResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to send email: %v", err),
		}, err
	}

	if response.StatusCode >= 400 {
		log.Printf("SendGrid error response: %d - %s", response.StatusCode, response.Body)
		return &domain.EmailResponse{
			Success: false,
			Message: fmt.Sprintf("SendGrid error: %d", response.StatusCode),
		}, fmt.Errorf("sendgrid error: %d", response.StatusCode)
	}

	return &domain.EmailResponse{
		Success:   true,
		MessageID: response.Headers["X-Message-Id"][0],
		Message:   "Email sent successfully",
	}, nil
}

// SendTemplatedEmail sends email using a template
func (s *sendgridService) SendTemplatedEmail(ctx context.Context, template *domain.EmailTemplate, to []string, variables map[string]string) (*domain.EmailResponse, error) {
	// Replace variables in template content
	content := template.Content
	for key, value := range variables {
		content = strings.ReplaceAll(content, "{{"+key+"}}", value)
	}

	// Replace variables in subject
	subject := template.Subject
	for key, value := range variables {
		subject = strings.ReplaceAll(subject, "{{"+key+"}}", value)
	}

	request := &domain.SendEmailRequest{
		To:      to,
		Subject: subject,
		Content: content,
		IsHTML:  template.IsHTML,
	}

	return s.SendEmail(ctx, request)
}

// SendWelcomeEmail sends welcome email to new users
func (s *sendgridService) SendWelcomeEmail(ctx context.Context, to string, name string) (*domain.EmailResponse, error) {
	template := &domain.EmailTemplate{
		Name:    "welcome",
		Subject: "Welcome to LiveChat System",
		Content: `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Welcome to LiveChat</title>
</head>
<body>
    <div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #333; text-align: center;">Welcome to LiveChat System!</h1>
        <p>Hi {{name}},</p>
        <p>Welcome to our LiveChat system! Your account has been successfully created.</p>
        <p>You can now log in to the system and start using all the features available to you.</p>
        <div style="text-align: center; margin: 30px 0;">
            <a href="#" style="background-color: #007bff; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px;">Get Started</a>
        </div>
        <p>If you have any questions, please don't hesitate to contact our support team.</p>
        <p>Best regards,<br>LiveChat Team</p>
    </div>
</body>
</html>`,
		IsHTML: true,
	}

	variables := map[string]string{
		"name": name,
	}

	return s.SendTemplatedEmail(ctx, template, []string{to}, variables)
}

// SendPasswordResetEmail sends password reset email
func (s *sendgridService) SendPasswordResetEmail(ctx context.Context, to string, resetToken string) (*domain.EmailResponse, error) {
	// You should replace this with your actual reset URL
	resetURL := fmt.Sprintf("https://yourapp.com/reset-password?token=%s", resetToken)

	template := &domain.EmailTemplate{
		Name:    "password_reset",
		Subject: "Reset Your Password - LiveChat System",
		Content: `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Password Reset</title>
</head>
<body>
    <div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #333; text-align: center;">Password Reset Request</h1>
        <p>You have requested to reset your password for your LiveChat account.</p>
        <p>Click the button below to reset your password:</p>
        <div style="text-align: center; margin: 30px 0;">
            <a href="{{reset_url}}" style="background-color: #dc3545; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px;">Reset Password</a>
        </div>
        <p>This link will expire in 1 hour for security reasons.</p>
        <p>If you didn't request this password reset, please ignore this email.</p>
        <p>Best regards,<br>LiveChat Team</p>
    </div>
</body>
</html>`,
		IsHTML: true,
	}

	variables := map[string]string{
		"reset_url": resetURL,
	}

	return s.SendTemplatedEmail(ctx, template, []string{to}, variables)
}

// SendChatTranscriptEmail sends chat transcript to customer
func (s *sendgridService) SendChatTranscriptEmail(ctx context.Context, to string, transcript string, sessionID uuid.UUID) (*domain.EmailResponse, error) {
	template := &domain.EmailTemplate{
		Name:    "chat_transcript",
		Subject: "Your Chat Transcript - LiveChat System",
		Content: `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Chat Transcript</title>
</head>
<body>
    <div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
        <h1 style="color: #333; text-align: center;">Your Chat Transcript</h1>
        <p>Thank you for contacting us! Here is the transcript of your chat session.</p>
        <div style="background-color: #f8f9fa; padding: 20px; border-radius: 5px; margin: 20px 0;">
            <h3>Session ID: {{session_id}}</h3>
            <div style="white-space: pre-wrap; font-family: monospace; font-size: 14px;">{{transcript}}</div>
        </div>
        <p>If you need further assistance, please don't hesitate to contact us again.</p>
        <p>Best regards,<br>LiveChat Support Team</p>
    </div>
</body>
</html>`,
		IsHTML: true,
	}

	variables := map[string]string{
		"session_id": sessionID.String(),
		"transcript": transcript,
	}

	return s.SendTemplatedEmail(ctx, template, []string{to}, variables)
}
