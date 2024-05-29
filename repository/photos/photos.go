package photos

import (
	"context"
	"dip/domain"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
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
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}

	for _, photo := range photos {
		_, err = tx.Exec(ctx, `INSERT INTO photos (restaurantID, url) VALUES ($1, $2)`, photo.RestaurantID, photo.URl)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *PhotoRepo) Delete(ctx context.Context, url string, restaurantID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM photos WHERE url = $1 AND restaurantID = $2`, url, restaurantID)
	return err
}
