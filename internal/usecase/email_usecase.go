package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/novianakbar/livechat-be/internal/domain"
)

type EmailUseCase interface {
	SendEmail(ctx context.Context, request *domain.SendEmailRequest) (*domain.EmailResponse, error)
	SendWelcomeEmail(ctx context.Context, email, name string) (*domain.EmailResponse, error)
	SendPasswordResetEmail(ctx context.Context, email, resetToken string) (*domain.EmailResponse, error)
	SendChatTranscriptEmail(ctx context.Context, email, transcript string, sessionID uuid.UUID) (*domain.EmailResponse, error)
	SendCustomEmail(ctx context.Context, template *domain.EmailTemplate, to []string, variables map[string]string) (*domain.EmailResponse, error)
}

type emailUseCase struct {
	emailService domain.EmailService
}

// NewEmailUseCase creates new email use case
func NewEmailUseCase(emailService domain.EmailService) EmailUseCase {
	return &emailUseCase{
		emailService: emailService,
	}
}

// SendEmail sends a custom email
func (uc *emailUseCase) SendEmail(ctx context.Context, request *domain.SendEmailRequest) (*domain.EmailResponse, error) {
	if len(request.To) == 0 {
		return &domain.EmailResponse{
			Success: false,
			Message: "No recipients specified",
		}, fmt.Errorf("no recipients specified")
	}

	if request.Subject == "" {
		return &domain.EmailResponse{
			Success: false,
			Message: "Email subject is required",
		}, fmt.Errorf("email subject is required")
	}

	if request.Content == "" {
		return &domain.EmailResponse{
			Success: false,
			Message: "Email content is required",
		}, fmt.Errorf("email content is required")
	}

	return uc.emailService.SendEmail(ctx, request)
}

// SendWelcomeEmail sends welcome email to new users
func (uc *emailUseCase) SendWelcomeEmail(ctx context.Context, email, name string) (*domain.EmailResponse, error) {
	if email == "" {
		return &domain.EmailResponse{
			Success: false,
			Message: "Email address is required",
		}, fmt.Errorf("email address is required")
	}

	if name == "" {
		return &domain.EmailResponse{
			Success: false,
			Message: "Name is required",
		}, fmt.Errorf("name is required")
	}

	return uc.emailService.SendWelcomeEmail(ctx, email, name)
}

// SendPasswordResetEmail sends password reset email
func (uc *emailUseCase) SendPasswordResetEmail(ctx context.Context, email, resetToken string) (*domain.EmailResponse, error) {
	if email == "" {
		return &domain.EmailResponse{
			Success: false,
			Message: "Email address is required",
		}, fmt.Errorf("email address is required")
	}

	if resetToken == "" {
		return &domain.EmailResponse{
			Success: false,
			Message: "Reset token is required",
		}, fmt.Errorf("reset token is required")
	}

	return uc.emailService.SendPasswordResetEmail(ctx, email, resetToken)
}

// SendChatTranscriptEmail sends chat transcript to customer
func (uc *emailUseCase) SendChatTranscriptEmail(ctx context.Context, email, transcript string, sessionID uuid.UUID) (*domain.EmailResponse, error) {
	if email == "" {
		return &domain.EmailResponse{
			Success: false,
			Message: "Email address is required",
		}, fmt.Errorf("email address is required")
	}

	if transcript == "" {
		return &domain.EmailResponse{
			Success: false,
			Message: "Transcript is required",
		}, fmt.Errorf("transcript is required")
	}

	if sessionID == uuid.Nil {
		return &domain.EmailResponse{
			Success: false,
			Message: "Session ID is required",
		}, fmt.Errorf("session ID is required")
	}

	return uc.emailService.SendChatTranscriptEmail(ctx, email, transcript, sessionID)
}

// SendCustomEmail sends email using custom template
func (uc *emailUseCase) SendCustomEmail(ctx context.Context, template *domain.EmailTemplate, to []string, variables map[string]string) (*domain.EmailResponse, error) {
	if template == nil {
		return &domain.EmailResponse{
			Success: false,
			Message: "Email template is required",
		}, fmt.Errorf("email template is required")
	}

	if len(to) == 0 {
		return &domain.EmailResponse{
			Success: false,
			Message: "No recipients specified",
		}, fmt.Errorf("no recipients specified")
	}

	if template.Subject == "" {
		return &domain.EmailResponse{
			Success: false,
			Message: "Email subject is required",
		}, fmt.Errorf("email subject is required")
	}

	if template.Content == "" {
		return &domain.EmailResponse{
			Success: false,
			Message: "Email content is required",
		}, fmt.Errorf("email content is required")
	}

	return uc.emailService.SendTemplatedEmail(ctx, template, to, variables)
}
