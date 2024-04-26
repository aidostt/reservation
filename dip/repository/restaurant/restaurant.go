package restaurant

import (
	"context"
	"dip/models"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgconn"

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

func (r *RestaurantRepo) GetById(ctx context.Context, id uuid.UUID) (*models.RestaurantSql, error) {
	query := `SELECT * FROM restaurants WHERE id = $1`
	var table models.RestaurantSql
	err := r.db.QueryRow(ctx, query, id).Scan(
		&table.ID,
		&table.Name,
		&table.Address,
		&table.Contact,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("not found in db")
		}
		return nil, err
	}

	return &table, nil
}

func (r *RestaurantRepo) GetAll(ctx context.Context) ([]*models.RestaurantSql, error) {
	query := "SELECT * FROM restaurants"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	restaurants := make([]*models.RestaurantSql, 0)
	for rows.Next() {
		restaurant := new(models.RestaurantSql)
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

func (r *RestaurantRepo) Create(ctx context.Context, rest *models.RestaurantSql) error {
	query := `
	INSERT INTO restaurants (name, address, contact) VALUES ($1, $2, $3)
	RETURNING id;
	`
	args := []any{rest.Name, rest.Address, rest.Contact}

	err := r.db.QueryRow(ctx, query, args...).Scan(&rest.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return err
		}
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

func (r *RestaurantRepo) UpdateById(ctx context.Context, upRest *models.RestaurantSql) error {
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
