package repository

import (
	"context"
	"kazokku/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CreditCardRepository interface {
	Insert(context.Context, pgx.Tx, domain.CreditCard) error
	Update(context.Context, pgx.Tx, domain.CreditCard) error
}

type creditCardRepository struct {
	db *pgxpool.Pool
}

func NewCreditCardRepository(db *pgxpool.Pool) creditCardRepository {
	return creditCardRepository{db}
}

func (repo creditCardRepository) Insert(ctx context.Context, tx pgx.Tx, data domain.CreditCard) error {
	stmt := "INSERT INTO credit_cards(user_id, type, number, name, expired, cvv) VALUES ($1, $2, $3, $4, $5, $6);"

	_, err := tx.Exec(ctx, stmt, data.UserID, data.Type, data.Number, data.Name, data.Expired, data.CVV)
	if err != nil {
		return err
	}

	return nil
}

func (repo creditCardRepository) Update(ctx context.Context, tx pgx.Tx, data domain.CreditCard) error {
	stmt := "UPDATE credit_cards SET type = COALESCE($1, type), number = COALESCE($2, number), name = COALESCE($3, name), expired = COALESCE($4, expired), cvv = COALESCE($5, cvv) WHERE user_id = $6;"

	_, err := tx.Exec(ctx, stmt, data.Type, data.Number, data.Name, data.Expired, data.CVV, data.UserID)
	if err != nil {
		return err
	}

	return nil
}
