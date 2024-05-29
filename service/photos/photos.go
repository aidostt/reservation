package photos

import (
	"context"
	"dip/domain"
	repo "dip/repository"
	"github.com/gofrs/uuid"
)

type PhotoService struct {
	repo repo.Photos
}

func NewPhotosService(repo repo.Photos) *PhotoService {
	return &PhotoService{repo: repo}
}

func (s *PhotoService) Upload(ctx context.Context, photos []*domain.PhotoSql) error {
	return s.repo.Upload(ctx, photos)
}

func (s *PhotoService) Delete(ctx context.Context, url, restaurantID string) error {
	convertedID, err := uuid.FromString(restaurantID)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, url, convertedID)
}
