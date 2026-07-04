package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	ServiceName string     `db:"service_name" json:"service_name"`
	Price       int        `db:"price" json:"price"`
	UserID      uuid.UUID  `db:"user_id" json:"user_id"`
	StartDate   time.Time  `db:"start_date" json:"start_date"`
	EndDate     *time.Time `db:"end_date" json:"end_date,omitempty"`
}

type CreateSubscriptionInput struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserID      string `json:"user_id"`
	StartDate   string `json:"start_date"`   // "MM-YYYY"
	EndDate     string `json:"end_date"`     // "MM-YYYY" (optional)
}

type UpdateSubscriptionInput struct {
	ServiceName *string `json:"service_name"`
	Price       *int    `json:"price"`
	StartDate   *string `json:"start_date"`
	EndDate     *string `json:"end_date"`
}

func ParseMonthYear(s string) (time.Time, error) {
	t, err := time.Parse("01-2006", s)
	if err != nil {
		return time.Time{}, errors.New("invalid date format, expected MM-YYYY")
	}
	return t, nil
}

func (i *CreateSubscriptionInput) Validate() error {
	if i.ServiceName == "" {
		return errors.New("service_name is required")
	}
	if i.Price <= 0 {
		return errors.New("price must be positive")
	}
	if _, err := uuid.Parse(i.UserID); err != nil {
		return errors.New("invalid user_id format (UUID expected)")
	}
	if _, err := ParseMonthYear(i.StartDate); err != nil {
		return err
	}
	if i.EndDate != "" {
		if _, err := ParseMonthYear(i.EndDate); err != nil {
			return errors.New("invalid end_date format")
		}
	}
	return nil
}