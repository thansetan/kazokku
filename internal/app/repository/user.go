package repository

import (
	"context"
	"fmt"
	"kazokku/internal/app/delivery/dto"
	"kazokku/internal/domain"
	"kazokku/internal/helpers"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Insert(context.Context, pgx.Tx, domain.User) (uint, error)
	GetAll(context.Context, dto.UserQuery) ([]domain.User, error)
	GetByID(context.Context, uint) (domain.User, error)
	Update(context.Context, pgx.Tx, uint, domain.User) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) userRepository {
	return userRepository{db}
}

func (repo userRepository) Insert(ctx context.Context, tx pgx.Tx, data domain.User) (uint, error) {
	stmt := "INSERT INTO users(name, address, email, password) VALUES ($1, $2, $3, $4) RETURNING id;"
	var id uint
	err := tx.QueryRow(ctx, stmt, data.Name, data.Address, data.Email, data.Password).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (repo userRepository) GetAll(ctx context.Context, query dto.UserQuery) ([]domain.User, error) {
	stmt := fmt.Sprintf(`SELECT u.id, u.name, u.email, u.address, cc.type, cc.number, cc.name, cc.expired
	FROM users u
	JOIN credit_cards cc ON cc.user_id = u.id
	WHERE u.name ILIKE '%%%s%%'
	OR u.email ILIKE '%%%s%%'
	OR u.address ILIKE '%%%s%%'
	ORDER BY u.%s %s
	OFFSET %d
	LIMIT %d;`, query.Query, query.Query, query.Query, query.OrderBy, query.SortBy, query.Offset, query.Limit)
	var users []domain.User
	var ids []uint
	rows, err := repo.db.Query(ctx, stmt)
	if err != nil {
		return users, err
	}

	for rows.Next() {
		var user domain.User
		var cc domain.CreditCard
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Address, &cc.Type, &cc.Number, &cc.Name, &cc.Expired)
		if err != nil {
			return users, err
		}
		user.CreditCard = cc
		users = append(users, user)
		ids = append(ids, user.ID)
	}

	if users == nil {
		return users, nil
	}

	idsStr := helpers.JoinIDs(ids)
	// get photos
	stmt = fmt.Sprintf(`SELECT user_id, filename FROM photos WHERE user_id IN (%s);`, idsStr)
	rows, err = repo.db.Query(ctx, stmt)
	if err != nil {
		return users, err
	}

	for rows.Next() {
		var photo domain.Photo
		var userID uint
		err = rows.Scan(&userID, &photo.Filepath)
		if err != nil {
			return users, err
		}
		for i, user := range users {
			if user.ID == userID {
				users[i].Photos = append(users[i].Photos, photo)
			}
		}

	}

	return users, nil
}

func (repo userRepository) GetByID(ctx context.Context, userID uint) (domain.User, error) {
	stmt := `SELECT u.id, u.name, email, address, p.filename, cc.type, cc.number, cc.name, cc.expired
			FROM users u
			LEFT JOIN photos p ON p.user_id = u.id
			LEFT JOIN credit_cards cc ON cc.user_id = u.id
			WHERE u.id = $1;`
	var user domain.User
	rows, err := repo.db.Query(ctx, stmt, userID)
	if err != nil {
		return user, err
	}
	for rows.Next() {
		var photo domain.Photo
		var cc domain.CreditCard
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Address, &photo.Filepath, &cc.Type, &cc.Number, &cc.Name, &cc.Expired)
		if err != nil {
			return user, err
		}
		user.Photos = append(user.Photos, photo)
		user.CreditCard = cc
	}
	return user, nil
}

func (repo userRepository) Update(ctx context.Context, tx pgx.Tx, userID uint, data domain.User) error {
	stmt := "UPDATE users SET name = COALESCE($1, name), address = COALESCE($2, address), email = COALESCE($3, email), password = COALESCE($4, password) WHERE id = $5;"

	cmd, err := tx.Exec(ctx, stmt, data.Name, data.Address, data.Email, data.Password, userID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return helpers.ErrUserNotFound
	}

	return nil
}
