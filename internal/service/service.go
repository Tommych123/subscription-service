package service

import (
	"time"

	"github.com/Tommych123/subscription-service/internal/model"
	"github.com/Tommych123/subscription-service/internal/repository"
)

type SubscriptionService struct {
	repo *repository.SubscriptionRepository
}

func NewSubscriptionService(repo *repository.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) Create(sub *model.Subscription) (string, error) {
	return s.repo.Create(sub)
}

func (s *SubscriptionService) GetByID(id string) (*model.Subscription, error) {
	return s.repo.GetByID(id)
}

func (s *SubscriptionService) Update(sub *model.Subscription) error {
	return s.repo.Update(sub)
}

func (s *SubscriptionService) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s *SubscriptionService) List() ([]model.Subscription, error) {
	return s.repo.List()
}

func (s *SubscriptionService) GetTotalCost(userID string, serviceName string, from, to time.Time) (int, error) {
	subs, err := s.repo.List()
	if err != nil {
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
		start := sub.StartDate
		end := time.Now()
		if sub.EndDate != nil {
			end = *sub.EndDate
		}
		if end.Before(from) || start.After(to) {
			continue
		}
		actualStart := maxDate(start, from)
		actualEnd := minDate(end, to)
		months := diffMonths(actualStart, actualEnd)

		total += sub.Price * months
	}

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
