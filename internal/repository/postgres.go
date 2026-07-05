package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"subscription-service/internal/models"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, sub *models.Subscription) error {
	query := `INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date)
              VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, sub.ID, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate)
	return err
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	var sub models.Subscription
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id=$1`
	err := r.db.GetContext(ctx, &sub, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &sub, err
}

func (r *Repository) Update(ctx context.Context, sub *models.Subscription) error {
	query := `UPDATE subscriptions SET service_name=$1, price=$2, user_id=$3, start_date=$4, end_date=$5 WHERE id=$6`
	_, err := r.db.ExecContext(ctx, query, sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate, sub.ID)
	return err
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE id=$1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *Repository) List(ctx context.Context, limit, offset int, userID *uuid.UUID, serviceName *string) ([]models.Subscription, error) {
	var subs []models.Subscription
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE 1=1`
	args := []interface{}{}
	argIdx := 1

	if userID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argIdx)
		args = append(args, *userID)
		argIdx++
	}
	if serviceName != nil {
		query += fmt.Sprintf(" AND service_name = $%d", argIdx)
		args = append(args, *serviceName)
		argIdx++
	}
	query += fmt.Sprintf(" ORDER BY start_date LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, limit, offset)

	err := r.db.SelectContext(ctx, &subs, query, args...)
	return subs, err
}

func (r *Repository) SumByPeriod(ctx context.Context, start, end time.Time, userID *uuid.UUID, serviceName *string) (int, error) {
	query := `SELECT COALESCE(SUM(price), 0) FROM subscriptions
              WHERE start_date <= $1 AND (end_date IS NULL OR end_date >= $2)`
	args := []interface{}{end, start}
	argIdx := 3

	if userID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argIdx)
		args = append(args, *userID)
		argIdx++
	}
	if serviceName != nil {
		query += fmt.Sprintf(" AND service_name = $%d", argIdx)
		args = append(args, *serviceName)
		argIdx++
	}

	var sum int
	err := r.db.GetContext(ctx, &sum, query, args...)
	return sum, err
}