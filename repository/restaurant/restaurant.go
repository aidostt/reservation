package restaurant

import (
	"context"
	"dip/domain"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RestaurantRepo struct {
	db *pgxpool.Pool
}

func NewRestaurantRepo(db *pgxpool.Pool) *RestaurantRepo {
	return &RestaurantRepo{
		db: db,
	}
}

func (r *RestaurantRepo) GetById(ctx context.Context, id uuid.UUID) (*domain.RestaurantSql, error) {
	query := `SELECT * FROM restaurants WHERE id = $1`
	var table domain.RestaurantSql
	err := r.db.QueryRow(ctx, query, id).Scan(
		&table.ID,
		&table.Name,
		&table.Address,
		&table.Contact,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFoundInDB
		}
		return nil, err
	}

	return &table, nil
}

func (r *RestaurantRepo) GetAll(ctx context.Context) ([]*domain.RestaurantSql, error) {
	query := "SELECT * FROM restaurants"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	restaurants := make([]*domain.RestaurantSql, 0)
	for rows.Next() {
		restaurant := new(domain.RestaurantSql)
		err := rows.Scan(&restaurant.ID, &restaurant.Name, &restaurant.Address, &restaurant.Contact)
		if err != nil {
			return nil, err
		}
		restaurants = append(restaurants, restaurant)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return restaurants, nil
}

func (r *RestaurantRepo) Create(ctx context.Context, rest *domain.RestaurantSql) error {
	query := `
	INSERT INTO restaurants (name, address, contact) VALUES ($1, $2, $3)
	RETURNING id;
	`
	args := []any{rest.Name, rest.Address, rest.Contact}

	err := r.db.QueryRow(ctx, query, args...).Scan(&rest.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *RestaurantRepo) Delete(ctx context.Context, restId uuid.UUID) error {
	query := `DELETE FROM restaurants WHERE id = $1`

	_, err := r.db.Exec(ctx, query, restId)
	if err != nil {
		return err
	}
	return nil
}

func (r *RestaurantRepo) UpdateById(ctx context.Context, upRest *domain.RestaurantSql) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	query := "UPDATE restaurants SET name = $1, address = $2, contact = $3 WHERE id = $4"
	_, err = tx.Exec(ctx, query, upRest.Name, upRest.Address, upRest.Contact, upRest.ID)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
