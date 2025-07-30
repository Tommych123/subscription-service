package repository

import (
	"database/sql"
	"errors"
	"github.com/Tommych123/subscription-service/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SubscriptionRepository struct {
	db *sqlx.DB
}

func NewSubscriptionRepository(db *sqlx.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(sub *model.Subscription) (string, error) {
	id := uuid.New().String()
	query := `
		INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date)
		VALUES (:id, :service_name, :price, :user_id, :start_date, :end_date)
	`
	subWithID := *sub
	subWithID.ID = id
	_, err := r.db.NamedExec(query, subWithID)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *SubscriptionRepository) GetByID(id string) (*model.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE id = $1
	`
	var sub model.Subscription
	err := r.db.Get(&sub, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &sub, nil
}

func (r *SubscriptionRepository) Update(sub *model.Subscription) error {
	query := `
		UPDATE subscriptions
		SET service_name = :service_name,
			price = :price,
			user_id = :user_id,
			start_date = :start_date,
			end_date = :end_date
		WHERE id = :id
	`
	_, err := r.db.NamedExec(query, sub)
	return err
}

func (r *SubscriptionRepository) Delete(id string) error {
	_, err := r.db.Exec("DELETE FROM subscriptions WHERE id = $1", id)
	return err
}

func (r *SubscriptionRepository) List() ([]model.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
	`
	var subs []model.Subscription
	err := r.db.Select(&subs, query)
	if err != nil {
		return nil, err
	}
	return subs, nil
}
