package restaurant

import (
	"context"
	"dip/models"
	repo "dip/repository"

	"github.com/gofrs/uuid"
)

type RestaurantService struct {
	repo repo.Restaurants
}

func NewRestaurantService(repo repo.Restaurants) *RestaurantService {
	return &RestaurantService{repo: repo}
}

func (s *RestaurantService) GetById(ctx context.Context, id string) (*models.RestaurantSql, error) {
	newTableId, err := uuid.FromString(id)
	if err != nil {
		return nil, err
	}

	return s.repo.GetById(ctx, newTableId)
}

func (s *RestaurantService) GetAll(ctx context.Context) ([]*models.RestaurantSql, error) {
	return s.repo.GetAll(ctx)
}

func (s *RestaurantService) Create(ctx context.Context, res *models.RestaurantSql) error {
	return s.repo.Create(ctx, res)
}

func (s *RestaurantService) DeleteById(ctx context.Context, restId string) error {
	newTableId, err := uuid.FromString(restId)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, newTableId)
}

func (s *RestaurantService) UpdateById(ctx context.Context, upRest *models.UpdateRestaurantInputSql) error {
	newRestId, err := uuid.FromString(upRest.RestaurantId)
	if err != nil {
		return err
	}

	newRestaurant := models.RestaurantSql{
		ID:      newRestId,
		Name:    upRest.Name,
		Address: upRest.Address,
		Contact: upRest.Contact,
	}

	return s.repo.UpdateById(ctx, &newRestaurant)
}
