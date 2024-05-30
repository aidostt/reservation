package photos

import (
	"context"
	"dip/domain"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

//TODO: implement delete func

type PhotoRepo struct {
	db *pgxpool.Pool
}

func NewPhotoRepo(db *pgxpool.Pool) *PhotoRepo {
	return &PhotoRepo{
		db: db,
	}
}

func (r *PhotoRepo) Upload(ctx context.Context, photos []*domain.PhotoSql) error {
	for _, photo := range photos {
		_, err := r.db.Exec(ctx, `INSERT INTO photos (restaurantID, url) VALUES ($1, $2)`, photo.RestaurantID, photo.URl)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if strings.Contains(pgErr.ConstraintName, "photos_url_key") {
					continue
				}
			} else {
				return err
			}
		}
	}
	return nil
}

func (r *PhotoRepo) Delete(ctx context.Context, url string, restaurantID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM photos WHERE url = $1 AND restaurantID = $2`, url, restaurantID)
	return err
}
