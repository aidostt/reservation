package restaurant

import (
	"context"
	"dip/models"
	repo "dip/repository"

	"github.com/jackc/pgx/v5/pgtype"
)

type RestaurantService struct {
	repo repo.Restaurants
}

func NewRestaurantService(repo repo.Restaurants) *RestaurantService {
	return &RestaurantService{repo: repo}
}

func (s *RestaurantService) GetById(ctx context.Context, id pgtype.UUID) (*models.RestaurantSql, error) {
	return s.repo.GetById(ctx, id)
}

func (s *RestaurantService) GetAll(ctx context.Context) ([]*models.RestaurantSql, error) {
	return s.repo.GetAll(ctx)
}

func (s *RestaurantService) Create(ctx context.Context, res *models.RestaurantSql) error {
	return s.repo.Create(ctx, res)
}

func (s *RestaurantService) DeleteById(ctx context.Context, restId pgtype.UUID) error {
	return s.repo.Delete(ctx, restId)
}

func (s *RestaurantService) UpdateById(ctx context.Context, upRest *models.UpdateRestaurantInputSql) error {
	return s.repo.UpdateById(ctx, upRest)
}
