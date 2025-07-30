package service

import (
	"time"
	"github.com/Tommych123/subscription-service/models"
	"github.com/Tommych123/subscription-service/repository"
	"go.uber.org/zap"
)

type SubscriptionService struct {
	repo   *repository.SubscriptionRepository
	logger *zap.Logger
}

func NewSubscriptionService(repo *repository.SubscriptionRepository, logger *zap.Logger) *SubscriptionService {
	return &SubscriptionService{
		repo:   repo,
		logger: logger,
	}
}

func (s *SubscriptionService) Create(sub *models.Subscription) (string, error) {
	id, err := s.repo.Create(sub)
	if err != nil {
		s.logger.Error("Failed to create subscription", zap.Error(err), zap.Any("subscription", sub))
		return "", err
	}
	s.logger.Info("Subscription created", zap.String("id", id))
	return id, nil
}

func (s *SubscriptionService) GetByID(id string) (*models.Subscription, error) {
	sub, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("Failed to get subscription by ID", zap.Error(err), zap.String("id", id))
		return nil, err
	}
	return sub, nil
}

func (s *SubscriptionService) Update(sub *models.Subscription) error {
	err := s.repo.Update(sub)
	if err != nil {
		s.logger.Error("Failed to update subscription", zap.Error(err), zap.Any("subscription", sub))
		return err
	}
	s.logger.Info("Subscription updated", zap.String("id", sub.ID))
	return nil
}

func (s *SubscriptionService) Delete(id string) error {
	err := s.repo.Delete(id)
	if err != nil {
		s.logger.Error("Failed to delete subscription", zap.Error(err), zap.String("id", id))
		return err
	}
	s.logger.Info("Subscription deleted", zap.String("id", id))
	return nil
}

func (s *SubscriptionService) List() ([]models.Subscription, error) {
	subs, err := s.repo.List()
	if err != nil {
		s.logger.Error("Failed to list subscriptions", zap.Error(err))
		return nil, err
	}
	return subs, nil
}

func (s *SubscriptionService) GetTotalCost(userID string, serviceName string, from, to time.Time) (int, error) {
	subs, err := s.repo.List()
	if err != nil {
		s.logger.Error("Failed to list subscriptions for total cost", zap.Error(err))
		return 0, err
	}
	total := 0
	for _, sub := range subs {
		if userID != "" && sub.UserID != userID {
			continue
		}
		if serviceName != "" && sub.ServiceName != serviceName {
			continue
		}
		start := sub.StartDate.Time
		end := time.Now()
		if sub.EndDate != nil {
			end = sub.EndDate.Time
		}
		if end.Before(from) || start.After(to) {
			continue
		}
		actualStart := maxDate(start, from)
		actualEnd := minDate(end, to)
		months := diffMonths(actualStart, actualEnd)

		total += sub.Price * months
	}
	s.logger.Info("Calculated total cost", zap.String("user_id", userID), zap.String("service_name", serviceName), zap.Time("from", from), zap.Time("to", to), zap.Int("total_cost", total))
	return total, nil
}


func maxDate(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func minDate(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func diffMonths(from, to time.Time) int {
	from = time.Date(from.Year(), from.Month(), 1, 0, 0, 0, 0, time.UTC)
	to = time.Date(to.Year(), to.Month(), 1, 0, 0, 0, 0, time.UTC)
	months := int(to.Month()) - int(from.Month()) + 12*(to.Year()-from.Year()) + 1
	if months < 0 {
		return 0
	}
	return months
}
