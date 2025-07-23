package usecase

import (
	"context"
	"time"

	"github.com/novianakbar/livechat-be/internal/domain"
)

type AnalyticsUsecase struct {
	sessionRepo  domain.ChatSessionRepository
	messageRepo  domain.ChatMessageRepository
	userRepo     domain.UserRepository
	customerRepo domain.CustomerRepository
}

func NewAnalyticsUsecase(
	sessionRepo domain.ChatSessionRepository,
	messageRepo domain.ChatMessageRepository,
	userRepo domain.UserRepository,
	customerRepo domain.CustomerRepository,
) *AnalyticsUsecase {
	return &AnalyticsUsecase{
		sessionRepo:  sessionRepo,
		messageRepo:  messageRepo,
		userRepo:     userRepo,
		customerRepo: customerRepo,
	}
}

func (u *AnalyticsUsecase) GetDashboardStats() (*domain.DashboardStats, error) {
	ctx := context.Background()
	stats := &domain.DashboardStats{}

	// Get active sessions count
	activeSessions, err := u.sessionRepo.CountByStatus(ctx, "active")
	if err != nil {
		return nil, err
	}
	stats.ActiveSessions = int(activeSessions)

	// Get waiting sessions count
	waitingSessions, err := u.sessionRepo.CountByStatus(ctx, "waiting")
	if err != nil {
		return nil, err
	}
	stats.WaitingSessions = int(waitingSessions)

	// Get completed sessions today
	today := time.Now().Truncate(24 * time.Hour)
	completedToday, err := u.sessionRepo.CountCompletedSince(ctx, today)
	if err != nil {
		return nil, err
	}
	stats.CompletedToday = int(completedToday)

	// Get average response time
	avgResponseTime, err := u.sessionRepo.GetAverageResponseTime(ctx)
	if err != nil {
		return nil, err
	}
	stats.AverageResponseTime = int(avgResponseTime)

	// Get total agents
	totalAgents, err := u.userRepo.CountByRole(ctx, "agent")
	if err != nil {
		return nil, err
	}
	stats.TotalAgents = int(totalAgents)

	// Get online agents
	onlineAgents, err := u.userRepo.CountOnlineAgents(ctx)
	if err != nil {
		return nil, err
	}
	stats.OnlineAgents = int(onlineAgents)

	// Get top questions (based on frequent topics/messages)
	topQuestions, err := u.messageRepo.GetTopQuestions(ctx, 5)
	if err != nil {
		return nil, err
	}
	stats.TopQuestions = topQuestions

	// Get OSS categories statistics
	ossCategories, err := u.sessionRepo.GetOSSCategoriesStats(ctx)
	if err != nil {
		return nil, err
	}
	stats.OSSCategories = ossCategories

	return stats, nil
}

func (u *AnalyticsUsecase) GetAnalytics(req *domain.GetAnalyticsRequest) ([]*domain.ChatAnalytics, error) {
	// Implementation for detailed analytics
	// This would get analytics data based on the request parameters
	// For now, return empty array
	return []*domain.ChatAnalytics{}, nil
}
