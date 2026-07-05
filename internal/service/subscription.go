package service

import (
	"context"
	"errors"
	"time"
	"github.com/google/uuid"
	"subscription-service/internal/models"
	"subscription-service/internal/repository"
)

type SubscriptionService struct {
	repo *repository.Repository
}

func NewSubscriptionService(repo *repository.Repository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) Create(ctx context.Context, input *models.CreateSubscriptionInput) (*models.Subscription, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	start, _ := models.ParseMonthYear(input.StartDate)
	var end *time.Time
	if input.EndDate != "" {
		t, _ := models.ParseMonthYear(input.EndDate)
		end = &t
	}
	userID, _ := uuid.Parse(input.UserID)

	sub := &models.Subscription{
		ID:          uuid.New(),
		ServiceName: input.ServiceName,
		Price:       input.Price,
		UserID:      userID,
		StartDate:   start,
		EndDate:     end,
	}

	err := s.repo.Create(ctx, sub)
	if err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *SubscriptionService) Get(ctx context.Context, id string) (*models.Subscription, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid id format")
	}
	return s.repo.GetByID(ctx, uid)
}

func (s *SubscriptionService) Update(ctx context.Context, id string, input *models.UpdateSubscriptionInput) (*models.Subscription, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid id format")
	}
	existing, err := s.repo.GetByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("subscription not found")
	}

	if input.ServiceName != nil {
		existing.ServiceName = *input.ServiceName
	}
	if input.Price != nil {
		existing.Price = *input.Price
	}
	if input.StartDate != nil {
		t, err := models.ParseMonthYear(*input.StartDate)
		if err != nil {
			return nil, err
		}
		existing.StartDate = t
	}
	if input.EndDate != nil {
		if *input.EndDate == "" {
			existing.EndDate = nil
		} else {
			t, err := models.ParseMonthYear(*input.EndDate)
			if err != nil {
				return nil, err
			}
			existing.EndDate = &t
		}
	}

	err = s.repo.Update(ctx, existing)
	return existing, err
}

func (s *SubscriptionService) Delete(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid id format")
	}
	return s.repo.Delete(ctx, uid)
}

func (s *SubscriptionService) List(ctx context.Context, limit, offset int, userID, serviceName *string) ([]models.Subscription, error) {
	var uid *uuid.UUID
	if userID != nil {
		u, err := uuid.Parse(*userID)
		if err != nil {
			return nil, errors.New("invalid user_id format")
		}
		uid = &u
	}
	return s.repo.List(ctx, limit, offset, uid, serviceName)
}

func (s *SubscriptionService) SumByPeriod(ctx context.Context, startStr, endStr string, userID, serviceName *string) (int, error) {
	start, err := models.ParseMonthYear(startStr)
	if err != nil {
		return 0, err
	}
	end, err := models.ParseMonthYear(endStr)
	if err != nil {
		return 0, err
	}
	if start.After(end) {
		return 0, errors.New("start date must be before end date")
	}

	var uid *uuid.UUID
	if userID != nil {
		u, err := uuid.Parse(*userID)
		if err != nil {
			return 0, errors.New("invalid user_id format")
		}
		uid = &u
	}

	return s.repo.SumByPeriod(ctx, start, end, uid, serviceName)
}