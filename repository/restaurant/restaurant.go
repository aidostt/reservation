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
	query := `SELECT id, name, address, contact FROM restaurants WHERE id = $1`
	var restaurant domain.RestaurantSql
	err := r.db.QueryRow(ctx, query, id).Scan(
		&restaurant.ID,
		&restaurant.Name,
		&restaurant.Address,
		&restaurant.Contact,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFoundInDB
		}
		return nil, err
	}

	photosQuery := `SELECT id, restaurantid, url FROM photos WHERE restaurantid = $1`
	rows, err := r.db.Query(ctx, photosQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var photo domain.PhotoSql
		if err = rows.Scan(&photo.ID, &photo.RestaurantID, &photo.URl); err != nil {
			return nil, err
		}
		restaurant.Photos = append(restaurant.Photos, photo)
	}

	return &restaurant, nil
}

func (r *RestaurantRepo) GetAll(ctx context.Context) ([]*domain.RestaurantSql, error) {
	query := "SELECT id, name, address, contact FROM restaurants"
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var restaurants []*domain.RestaurantSql
	for rows.Next() {
		restaurant := new(domain.RestaurantSql)
		err := rows.Scan(&restaurant.ID, &restaurant.Name, &restaurant.Address, &restaurant.Contact)
		if err != nil {
			return nil, err
		}

		photosQuery := `SELECT id, restaurantid, url FROM photos WHERE restaurantid = $1`
		photosRows, err := r.db.Query(ctx, photosQuery, restaurant.ID)
		if err != nil {
			return nil, err
		}
		defer photosRows.Close()

		for photosRows.Next() {
			var photo domain.PhotoSql
			if err := photosRows.Scan(&photo.ID, &photo.RestaurantID, &photo.URl); err != nil {
				return nil, err
			}
			restaurant.Photos = append(restaurant.Photos, photo)
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

func (r *RestaurantRepo) Search(ctx context.Context, query string, limit, offset int) ([]*domain.RestaurantSql, int, error) {
	querySQL := "SELECT id, name, address, contact FROM restaurants WHERE LOWER(name) LIKE $1 LIMIT $2 OFFSET $3"
	rows, err := r.db.Query(ctx, querySQL, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var restaurants []*domain.RestaurantSql
	for rows.Next() {
		restaurant := new(domain.RestaurantSql)
		err = rows.Scan(&restaurant.ID, &restaurant.Name, &restaurant.Address, &restaurant.Contact)
		if err != nil {
			return nil, 0, err
		}

		photosQuery := `SELECT id, restaurantid, url FROM photos WHERE restaurantid = $1`
		photosRows, err := r.db.Query(ctx, photosQuery, restaurant.ID)
		if err != nil {
			return nil, 0, err
		}
		defer photosRows.Close()

		for photosRows.Next() {
			var photo domain.PhotoSql
			if err = photosRows.Scan(&photo.ID, &photo.RestaurantID, &photo.URl); err != nil {
				return nil, 0, err
			}
			restaurant.Photos = append(restaurant.Photos, photo)
		}

		restaurants = append(restaurants, restaurant)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	// Perform a count query
	countQuery := "SELECT COUNT(*) FROM restaurants WHERE LOWER(name) LIKE $1"
	var total int
	err = r.db.QueryRow(ctx, countQuery, "%"+query+"%").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return restaurants, total, nil
}

func (r *RestaurantRepo) GetSuggestions(ctx context.Context, query string) ([]*domain.RestaurantSql, error) {
	querySQL := "SELECT id, name, address, contact FROM restaurants WHERE LOWER(name) LIKE $1 LIMIT 10"
	rows, err := r.db.Query(ctx, querySQL, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var restaurants []*domain.RestaurantSql
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
